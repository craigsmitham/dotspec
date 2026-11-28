package check

import (
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

// SourceFile records incoming occurrences, including explicit catalog declarations.
type SourceFile struct {
	Project      *specfile.Project     `json:"project,omitempty"`
	Path         string                `json:"path"`
	ReferencedBy []diagnostic.Location `json:"referenced_by"`
	Declarations []SourceDeclaration   `json:"declarations"`
}

// SourceDeclaration identifies the named source and enclosing specification.
// A verification procedure may assess a named specification within that document.
type SourceDeclaration struct {
	Name                  string              `json:"name"`
	Defines               string              `json:"defines,omitempty"`
	VerifiedSpecification string              `json:"verified_specification,omitempty"`
	Source                string              `json:"source"`
	Location              diagnostic.Location `json:"location"`
}

type SourceReport struct {
	SchemaVersion int                     `json:"schema_version"`
	SelectedFile  string                  `json:"selected_file,omitempty"`
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
// Contexts holds each inspected source's context resolution, keyed by source path.
type Inventory struct {
	Types             []Type
	Sources           *specfile.SourceSet
	Report            Report
	Files             []SourceFile
	Contexts          map[string]ContextResult
	Selected          []string
	Repository        string
	DiscoveryComplete bool
}

// ContextResult is what a source's ContextProvider returned during check:
// the resolved paths, including failed reads, and the resolution error.
type ContextResult struct {
	Paths []string
	Err   error
}

func InspectSources(path string, types ...Type) Inventory {
	result := Inventory{Types: slices.Clone(types), Files: []SourceFile{}}
	result.Report = project(path, &result, result.Types...)
	for _, file := range result.Report.files {
		sort.Slice(file.ReferencedBy, func(i, j int) bool { return locationLess(file.ReferencedBy[i], file.ReferencedBy[j]) })
		file.ReferencedBy = slices.Compact(file.ReferencedBy)
		sort.Slice(file.Declarations, func(i, j int) bool {
			a, b := file.Declarations[i], file.Declarations[j]
			if a.Location != b.Location {
				return locationLess(a.Location, b.Location)
			}
			return a.Name < b.Name
		})
		file.Declarations = slices.Compact(file.Declarations)
		result.Files = append(result.Files, *file)
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	result.Contexts = result.Report.contexts
	result.Report.files, result.Report.contexts = nil, nil
	return result
}

// locationLess orders locations by path and source position.
func locationLess(a, b diagnostic.Location) bool {
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
}

// collectFile records an existing file and, when reference is set, the link
// that reaches it. A catalog occurrence also identifies its declaration.
func (r *Report) collectFile(path string, reference *specfile.Reference) {
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
		file = &SourceFile{Path: path, ReferencedBy: []diagnostic.Location{}, Declarations: []SourceDeclaration{}}
		if p, err := specfile.FindProject(filepath.Dir(path)); err != nil {
			r.record(err, "scan.context", path)
		} else {
			file.Project = p
		}
		r.files[path] = file
	}
	if reference == nil {
		return
	}
	file.ReferencedBy = append(file.ReferencedBy, reference.Location)
	if reference.Name != "" {
		file.Declarations = append(file.Declarations, SourceDeclaration{Name: reference.Name, Defines: reference.Defines, VerifiedSpecification: reference.VerifiedSpecification, Source: reference.Destination, Location: reference.Location})
	}
}
