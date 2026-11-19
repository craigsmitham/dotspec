package check

import (
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// SourceFile records a file and its authored incoming links, without assigning
// authority or interpreting the meaning of its contents.
type SourceFile struct {
	Project      *specfile.Project     `json:"project,omitempty"`
	Path         string                `json:"path"`
	ReferencedBy []diagnostic.Location `json:"referenced_by"`
}

type SourceReport struct {
	SchemaVersion int                     `json:"schema_version"`
	Scope         Scope                   `json:"scope"`
	Complete      bool                    `json:"complete"`
	Files         []SourceFile            `json:"files"`
	Diagnostics   []diagnostic.Diagnostic `json:"diagnostics"`
}

func (r SourceReport) ExitCode() int {
	if !r.Complete {
		return 2
	}
	for _, finding := range r.Diagnostics {
		if finding.Severity == "error" {
			return 1
		}
	}
	return 0
}

// Sources inventories discovered specs, consulted context, and direct file
// references. Validation, boundaries, and diagnostics are identical to Project.
func Sources(path string, types ...Type) SourceReport {
	inspection := InspectSources(path, types...)
	r := inspection.Report
	return SourceReport{SchemaVersion: r.SchemaVersion, Scope: r.Scope, Complete: r.Complete, Files: inspection.Files, Diagnostics: r.Diagnostics}
}

// Inventory retains the observed project boundary, selected specs and source report
// for subsequent model mapping without a second discovery or source check.
type Inventory struct {
	Report            Report
	Files             []SourceFile
	Selected          []string
	Repository        string
	DiscoveryComplete bool
}

func InspectSources(path string, types ...Type) Inventory {
	result := Inventory{Files: []SourceFile{}}
	result.Report = project(path, &result, types...)
	for _, file := range result.Report.files {
		sort.Slice(file.ReferencedBy, func(i, j int) bool {
			a, b := file.ReferencedBy[i], file.ReferencedBy[j]
			if a.Path != b.Path {
				return a.Path < b.Path
			}
			if a.Line != b.Line {
				return a.Line < b.Line
			}
			if a.Column != b.Column {
				return a.Column < b.Column
			}
			return a.Element < b.Element
		})
		file.ReferencedBy = slices.Compact(file.ReferencedBy)
		result.Files = append(result.Files, *file)
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	result.Report.files = nil
	return result
}

func (r *Report) collectFile(path string, reference *diagnostic.Location) {
	if r.files == nil || path == "" {
		return
	}
	// Failed lookups stay in diagnostics, not the inventory of existing files.
	// In particular, context providers can report missing paths or directories.
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return
	}
	file := r.files[path]
	if file == nil {
		file = &SourceFile{Path: path, ReferencedBy: []diagnostic.Location{}}
		if p, err := specfile.FindProject(filepath.Dir(path)); err != nil {
			r.record(err, "scan.context", path)
		} else {
			file.Project = p
		}
		r.files[path] = file
	}
	if reference != nil {
		file.ReferencedBy = append(file.ReferencedBy, *reference)
	}
}
