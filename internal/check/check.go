// Package check inspects supported specifications without changing their files.
package check

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
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
func Project(path string) Report {
	r := Report{SchemaVersion: 1, Scope: Scope{Path: path, Types: []string{"system", "job", "use-case"}, Discovery: "unavailable"}, Complete: true, ContextFiles: []string{}, Diagnostics: []diagnostic.Diagnostic{}}
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
	paths, err := specfile.Discover(dir, supported)
	if err != nil {
		r.record(err, "scan.discovery", dir)
	}
	r.Summary.FilesDiscovered = len(paths)
	discovered := make(map[string]bool)
	var contexts []string
	for _, path := range paths {
		discovered[path] = true
		r.Summary.FilesChecked++
		if filepath.Base(path) == system.Filename {
			if _, err := system.Read(path); err != nil {
				r.record(err, "scan.read", path)
			}
			continue
		}
		collection := "jobs"
		if strings.HasSuffix(path, usecase.Suffix) {
			collection = usecase.Collection
			spec, err := usecase.Read(path)
			if err != nil {
				r.record(err, "scan.read", path)
				continue
			}
			r.Diagnostics = append(r.Diagnostics, usecase.Check(spec)...)
		} else {
			spec, err := job.Read(path)
			if err != nil {
				r.record(err, "scan.read", path)
				continue
			}
			r.Diagnostics = append(r.Diagnostics, job.Check(spec)...)
		}
		owner, err := specfile.Owner(path, collection)
		if err != nil {
			r.record(err, "scan.context", path)
			continue
		}
		contexts = append(contexts, owner)
	}
	// Match job and use-case show/list association semantics, including systems
	// outside the selected subtree. Report that supplemental scope explicitly.
	owners := make(map[string]bool)
	context := make(map[string]bool)
	for _, owner := range contexts {
		if owners[owner] {
			continue
		}
		owners[owner] = true
		associated, err := system.Enclosing(owner)
		contextPath := ""
		if err != nil {
			r.record(err, "scan.context", owner)
			var finding *diagnostic.Diagnostic
			var pathError *os.PathError
			if errors.As(err, &finding) && supported(finding.Location.Path) {
				contextPath = finding.Location.Path
			} else if errors.As(err, &pathError) && supported(pathError.Path) {
				contextPath = pathError.Path
			}
		} else if associated != nil {
			contextPath = associated.Path
		}
		if contextPath != "" && !discovered[contextPath] && !context[contextPath] {
			context[contextPath] = true
			r.ContextFiles = append(r.ContextFiles, contextPath)
			r.Summary.ContextFilesChecked++
		}
	}
	return r.finish()
}

func supported(path string) bool {
	name := filepath.Base(path)
	return name == system.Filename || strings.HasSuffix(name, job.Suffix) || strings.HasSuffix(name, usecase.Suffix)
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
