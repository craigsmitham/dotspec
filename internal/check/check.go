// Package check inspects supported specifications without changing their files.
package check

import (
	"errors"
	"os"
	"sort"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

type Scope struct {
	Path      string   `json:"path"`
	Types     []string `json:"types"`
	Discovery string   `json:"discovery"`
}

type Summary struct {
	FilesDiscovered     int `json:"files_discovered"`
	FilesChecked        int `json:"files_checked"`
	ContextFilesChecked int `json:"context_files_checked"`
	Errors              int `json:"errors"`
	Advice              int `json:"advice"`
}

type Report struct {
	files         map[string]*SourceFile
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
func project(path string, files map[string]*SourceFile, types ...Type) Report {
	r := newReport(path, "unavailable", types)
	r.files = files
	match := func(name string) bool { return findType(name, types) != nil }
	dir, err := specfile.Scope(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
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
		kind := findType(path, types) // Discovery returned only matching paths.
		findings, err := kind.Inspect(path)
		if err != nil {
			r.record(err, "scan.read", path)
			continue
		}
		r.Diagnostics = append(r.Diagnostics, findings...)
		referenceSources = append(referenceSources, path)
		if kind.Context != nil {
			owner, err := kind.Layout.Owner(path)
			if err != nil {
				r.record(err, "scan.context", path)
				continue
			}
			contexts = append(contexts, contextRequest{contextKey{provider: kind.Context, owner: owner}, path, kind.InspectContext})
		}
	}
	r.inspectContexts(contexts, discovered)
	r.inspectReferences(referenceSources, discovered, types)
	return r.finish()
}

func findType(path string, types []Type) *Type {
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
	return Report{SchemaVersion: 1, Scope: Scope{Path: path, Types: names, Discovery: discovery}, Complete: true, ContextFiles: []string{}, Diagnostics: []diagnostic.Diagnostic{}}
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
