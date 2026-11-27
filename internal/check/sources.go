package check

import (
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

// SourceFile records a file and its authored incoming links, without assigning
// authority or interpreting the meaning of its contents. CompanionOf repeats
// each frontmatter companion entry that targets the file with the role its
// author stated; the role is reported, not inferred or verified.
type SourceFile struct {
	Project      *specfile.Project     `json:"project,omitempty"`
	Path         string                `json:"path"`
	ReferencedBy []diagnostic.Location `json:"referenced_by"`
	CompanionOf  []SourceCompanion     `json:"companion_of"`
}

// SourceCompanion is one specification's frontmatter entry listing a file as
// its companion. The location's path is the declaring specification.
type SourceCompanion struct {
	Role     string              `json:"role"`
	Target   string              `json:"target"`
	Location diagnostic.Location `json:"location"`
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
		sort.Slice(file.CompanionOf, func(i, j int) bool {
			a, b := file.CompanionOf[i], file.CompanionOf[j]
			if a.Location != b.Location {
				return locationLess(a.Location, b.Location)
			}
			return a.Role < b.Role
		})
		file.CompanionOf = slices.Compact(file.CompanionOf)
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
// that reaches it. A frontmatter companion link also records its role.
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
		file = &SourceFile{Path: path, ReferencedBy: []diagnostic.Location{}, CompanionOf: []SourceCompanion{}}
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
	if reference.Role != "" {
		file.CompanionOf = append(file.CompanionOf, SourceCompanion{Role: reference.Role, Target: reference.Destination, Location: reference.Location})
	}
}
