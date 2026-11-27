// Package check inspects supported specifications without changing their files.
package check

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

type Scope struct {
	Project   *specfile.Project `json:"project,omitempty"`
	Path      string            `json:"path"`
	Types     []string          `json:"types"`
	Discovery string            `json:"discovery"`
}

type Summary struct {
	FilesDiscovered     int `json:"files_discovered"`
	FilesChecked        int `json:"files_checked"`
	ContextFilesChecked int `json:"context_files_checked"`
	Errors              int `json:"errors"`
	Advice              int `json:"advice"`
}

type Report struct {
	sources       *specfile.SourceSet
	files         map[string]*SourceFile
	contexts      map[string]ContextResult
	GeneratedSpec string                  `json:"generated_spec,omitempty"`
	SchemaVersion int                     `json:"schema_version"`
	Scope         Scope                   `json:"scope"`
	Complete      bool                    `json:"complete"`
	ContextFiles  []string                `json:"context_files"`
	Summary       Summary                 `json:"summary"`
	Diagnostics   []diagnostic.Diagnostic `json:"diagnostics"`
}

func (r Report) ExitCode() int {
	if !r.Complete {
		return 2
	}
	if r.Summary.Errors > 0 {
		return 1
	}
	return 0
}

// Project uses the same discovery policy as list: Git's ignore rules when in a
// repository, no ignore rules otherwise, and no symlinks or nested repositories.
// A malformed file blocks only its dependent checks, not the rest of the scan.
func Project(path string, types ...Type) Report {
	return project(path, nil, types...)
}

// The optional inventory collects files during the same scan used by check.
func project(path string, inventory *Inventory, types ...Type) Report {
	r := newReport(path, "unavailable", types)
	if inventory != nil {
		r.files = make(map[string]*SourceFile)
		r.contexts = make(map[string]ContextResult)
		inventory.Sources = r.sources
	}
	match := func(name string) bool { return MatchType(name, types) != nil }
	dir, err := specfile.Directory(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
	r.collectManifest()
	root, err := specfile.RepositoryRoot(dir)
	if err != nil {
		r.record(err, "scan.discovery", dir)
		return r.finish()
	}
	r.Scope.Discovery = "filesystem"
	if root != "" {
		r.Scope.Discovery = "git"
	}
	paths, err := specfile.Discover(dir, match)
	if inventory != nil {
		inventory.Selected = paths
		inventory.Repository = root
		inventory.DiscoveryComplete = err == nil
	}
	if err != nil {
		r.record(err, "scan.discovery", dir)
	}
	r.Summary.FilesDiscovered = len(paths)
	discovered := make(map[string]bool)
	var contexts []contextRequest
	var referenceSources []string
	for _, path := range paths {
		r.collectFile(path, nil)
		discovered[path] = true
		r.Summary.FilesChecked++
		kind := MatchType(path, types) // Discovery returned only matching paths.
		findings, err := kind.Inspect(r.sources, path)
		if err != nil {
			r.record(err, "scan.read", path)
			continue
		}
		r.Diagnostics = append(r.Diagnostics, findings...)
		referenceSources = append(referenceSources, path)
		if kind.Context != nil {
			owner := kind.Layout.ContextDir(path)
			contexts = append(contexts, contextRequest{contextKey{provider: kind.Context, owner: owner}, path, kind.InspectContext})
		}
	}
	r.inspectContexts(contexts, discovered)
	r.inspectReferences(referenceSources, discovered, types)
	r.inspectUnlistedCompanions(dir, paths, types)
	return r.finish()
}

// Formal models and Gherkin features nearly always carry intent. Each one in
// a whole-project scan should be some specification's companion. Listings come
// from every discovered specification's frontmatter, including a source with
// unrelated findings, and a listed directory covers the files beneath it.
var companionFormats = []string{".qnt", ".als", ".tla", ".feature"}

func (r *Report) inspectUnlistedCompanions(dir string, sources []string, types []Type) {
	if r.Scope.Project == nil || filepath.Clean(dir) != filepath.Clean(r.Scope.Project.Root) {
		return
	}
	var listed []string
	for _, source := range sources {
		kind := MatchType(source, types) // Discovery returned only matching paths.
		md, err := r.sources.Markdown(source)
		if err != nil {
			continue
		}
		// An invalid list is reported by the source's own inspection.
		companions, err := md.Companions(kind.Layout.Kind)
		if err != nil {
			continue
		}
		for _, c := range companions {
			u, err := url.Parse(c.Target)
			if err != nil || u.IsAbs() || u.Host != "" || u.Path == "" || strings.HasPrefix(u.Path, "/") {
				continue
			}
			listed = append(listed, filepath.Clean(filepath.Join(filepath.Dir(source), filepath.FromSlash(u.Path))))
		}
	}
	files, err := specfile.Discover(dir, func(name string) bool {
		return slices.Contains(companionFormats, strings.ToLower(filepath.Ext(name)))
	})
	if err != nil {
		r.record(err, "scan.discovery", dir)
		return
	}
	for _, file := range files {
		if !slices.ContainsFunc(listed, func(target string) bool {
			return file == target || strings.HasPrefix(file, target+string(filepath.Separator))
		}) {
			r.Diagnostics = append(r.Diagnostics, diagnostic.Advice("companion.unlisted", "No specification lists this file as a companion. Add it to the frontmatter companions of the specification whose intent it carries.", diagnostic.Location{Path: file}))
		}
	}
}

// MatchType returns the first registered layout that recognizes the source path.
func MatchType(path string, types []Type) *Type {
	for i := range types {
		if types[i].Layout.Match(path) {
			return &types[i]
		}
	}
	return nil
}

func newReport(path, discovery string, types []Type) Report {
	names := make([]string, 0, len(types))
	for _, kind := range types {
		names = append(names, kind.Layout.Kind)
	}
	scope := Scope{Path: path, Types: names, Discovery: discovery}
	return Report{sources: &specfile.SourceSet{}, SchemaVersion: 1, Scope: scope, Complete: true, ContextFiles: []string{}, Diagnostics: []diagnostic.Diagnostic{}}
}

// record and finish keep report accounting, error identity, and deterministic
// ordering identical for project-wide and focused checks.
func (r *Report) record(err error, rule, path string) {
	var finding *diagnostic.Diagnostic
	if !errors.As(err, &finding) {
		r.Complete = false
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			path = pathError.Path
		}
		finding = diagnostic.Error(rule, err.Error(), diagnostic.Location{Path: path})
	}
	r.Diagnostics = append(r.Diagnostics, *finding)
}

// WithFindings incorporates another inspection phase into the same sorted report.
func (r Report) WithFindings(ds ...diagnostic.Diagnostic) Report {
	r.Diagnostics = append(r.Diagnostics, ds...)
	return r.finish()
}

func (r Report) finish() Report {
	sort.Strings(r.ContextFiles)
	sort.SliceStable(r.Diagnostics, func(i, j int) bool {
		a, b := r.Diagnostics[i], r.Diagnostics[j]
		if a.Location.Path != b.Location.Path {
			return a.Location.Path < b.Location.Path
		}
		if a.Location.Line != b.Location.Line {
			return a.Location.Line < b.Location.Line
		}
		if a.Location.Column != b.Location.Column {
			return a.Location.Column < b.Location.Column
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		return a.Message < b.Message
	})
	seen := make(map[string]bool)
	findings := r.Diagnostics[:0]
	for _, d := range r.Diagnostics {
		key := d.Rule + "\x00" + d.Location.String() + "\x00" + d.Message
		if seen[key] {
			continue
		}
		seen[key] = true
		findings = append(findings, d)
	}
	r.Diagnostics = findings
	r.Summary.Errors, r.Summary.Advice = 0, 0
	for _, d := range r.Diagnostics {
		if d.Severity == "error" {
			r.Summary.Errors++
		} else {
			r.Summary.Advice++
		}
	}
	return r
}

func (r *Report) collectManifest() {
	p, err := specfile.FindProject(r.Scope.Path)
	if err != nil {
		r.record(err, "scan.scope", r.Scope.Path)
		return
	}
	r.Scope.Project = p
	if p != nil && p.Report.Enabled && (r.Scope.Path != p.Root || r.Scope.Discovery == "direct") {
		r.GeneratedSpec = "out-of-scope"
	}
	if p != nil {
		manifest := p.Manifest()
		r.ContextFiles = append(r.ContextFiles, manifest)
		r.Summary.ContextFilesChecked++
		r.collectFile(manifest, nil)
	}
}
