package check

import (
	"errors"
	"os"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

// FocusedJob checks one local job and the system needed to interpret it. It uses
// the same report and exit semantics as a project scan, without scanning siblings.
func FocusedJob(path, slug string) Report {
	r := Report{SchemaVersion: 1, Scope: Scope{Path: path, Types: []string{"job"}, Discovery: "direct"}, Complete: true, ContextFiles: []string{}, Diagnostics: []diagnostic.Diagnostic{}}
	dir, err := specfile.Scope(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
	if err := job.ValidateSlug(slug); err != nil {
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	// Match Show's placement checks before reading the selected file.
	if err := specfile.ValidateLocal(dir, "jobs", func(name string) bool {
		return strings.HasSuffix(name, job.Suffix)
	}); err != nil {
		r.record(err, "scan.read", dir)
		return r.finish()
	}
	source := specfile.Path(dir, "jobs", slug+job.Suffix)
	r.Summary.FilesDiscovered = 1
	r.Summary.FilesChecked = 1
	spec, err := job.Read(source)
	if err != nil {
		r.record(err, "scan.read", source)
		return r.finish()
	}
	r.Diagnostics = append(r.Diagnostics, job.Check(spec)...)
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
