package check

import (
	"github.com/craigsmitham/dotspec/internal/specfile"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
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
	files := make(map[string]*SourceFile)
	r := project(path, files, types...)
	result := SourceReport{SchemaVersion: r.SchemaVersion, Scope: r.Scope, Complete: r.Complete, Files: []SourceFile{}, Diagnostics: r.Diagnostics}
	for _, file := range files {
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
		if _, p, err := specfile.ProjectContext(filepath.Dir(path)); err == nil && p.Name != "" {
			file.Project = &p
		}
		r.files[path] = file
	}
	if reference != nil {
		file.ReferencedBy = append(file.ReferencedBy, *reference)
	}
}
