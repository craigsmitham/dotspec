package check

import (
	"errors"
	"io/fs"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Focused checks one local specification and the context requested by its type.
// It uses the same reporting rules as Project, without discovering siblings.
func Focused(path, slug string, kind Type, referenceTypes ...Type) Report {
	r := newReport(path, "direct", []Type{kind})
	dir, err := specfile.Directory(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
	r.collectManifest()
	source, err := kind.Layout.LocalPath(dir, slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			r.Summary.FilesDiscovered = 1
			r.Summary.FilesChecked = 1
			r.record(err, "scan.read", dir)
		} else {
			r.record(err, "scan.scope", dir)
		}
		return r.finish()
	}
	r.Summary.FilesDiscovered = 1
	r.Summary.FilesChecked = 1
	findings, err := kind.Inspect(r.sources, source)
	if err != nil {
		r.record(err, "scan.read", source)
		return r.finish()
	}
	r.Diagnostics = append(r.Diagnostics, findings...)
	if kind.Context != nil {
		r.inspectContexts([]contextRequest{{contextKey{provider: kind.Context, owner: kind.Layout.ContextDir(source)}, source, kind.InspectContext}}, map[string]bool{source: true})
	}
	r.inspectReferences([]string{source}, map[string]bool{source: true}, append([]Type{kind}, referenceTypes...))
	return r.finish()
}
