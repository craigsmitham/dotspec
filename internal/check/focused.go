package check

import (
	"errors"
	"os"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

// FocusedJob checks one local job and the system needed to interpret it. It uses
// the same report and exit semantics as a project scan, without scanning siblings.
func FocusedJob(path, slug string) Report {
	return focused(path, slug, focusedType{"job", "jobs", job.Suffix, job.ValidateSlug, func(source string) ([]diagnostic.Diagnostic, error) {
		spec, err := job.Read(source)
		if err != nil {
			return nil, err
		}
		return job.Check(spec), nil
	}})
}

// FocusedUseCase checks one local use case and its enclosing system.
func FocusedUseCase(path, slug string) Report {
	return focused(path, slug, focusedType{"use-case", usecase.Collection, usecase.Suffix, usecase.ValidateSlug, func(source string) ([]diagnostic.Diagnostic, error) {
		spec, err := usecase.Read(source)
		if err != nil {
			return nil, err
		}
		return usecase.Check(spec), nil
	}})
}

type focusedType struct {
	name, collection, suffix string
	validateSlug             func(string) error
	inspect                  func(string) ([]diagnostic.Diagnostic, error)
}

func focused(path, slug string, kind focusedType) Report {
	r := Report{SchemaVersion: 1, Scope: Scope{Path: path, Types: []string{kind.name}, Discovery: "direct"}, Complete: true, ContextFiles: []string{}, Diagnostics: []diagnostic.Diagnostic{}}
	dir, err := specfile.Scope(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
	if err := kind.validateSlug(slug); err != nil {
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	// Match Show's placement checks before reading the selected file.
	if err := specfile.ValidateLocal(dir, kind.collection, func(name string) bool {
		return strings.HasSuffix(name, kind.suffix)
	}); err != nil {
		r.record(err, "scan.read", dir)
		return r.finish()
	}
	source := specfile.Path(dir, kind.collection, slug+kind.suffix)
	r.Summary.FilesDiscovered = 1
	r.Summary.FilesChecked = 1
	findings, err := kind.inspect(source)
	if err != nil {
		r.record(err, "scan.read", source)
		return r.finish()
	}
	r.Diagnostics = append(r.Diagnostics, findings...)
	associated, err := system.Enclosing(dir)
	if err != nil {
		r.record(err, "scan.context", dir)
		var d *diagnostic.Diagnostic
		contextPath := ""
		var pathError *os.PathError
		if errors.As(err, &d) {
			contextPath = d.Location.Path
		} else if errors.As(err, &pathError) {
			contextPath = pathError.Path
		}
		if contextPath != "" {
			r.ContextFiles = append(r.ContextFiles, contextPath)
			r.Summary.ContextFilesChecked++
		}
		return r.finish()
	}
	if associated != nil {
		r.ContextFiles = append(r.ContextFiles, associated.Path)
		r.Summary.ContextFilesChecked++
	}
	return r.finish()
}
