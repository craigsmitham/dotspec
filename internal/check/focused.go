package check

import "github.com/craigsmitham/dotspec/internal/specfile"

// Focused checks one local specification and the context requested by its type.
// It uses the same reporting rules as Project, without inspecting siblings.
func Focused(path, slug string, kind Type) Report {
	r := newReport(path, "direct", []Type{kind})
	dir, err := specfile.Scope(path)
	if err != nil {
		r.Complete = false
		r.record(err, "scan.scope", path)
		return r.finish()
	}
	r.Scope.Path = dir
	if kind.ValidateSlug != nil {
		if err := kind.ValidateSlug(slug); err != nil {
			r.record(err, "scan.scope", path)
			return r.finish()
		}
	}
	if err := kind.Layout.ValidateLocal(dir); err != nil {
		r.record(err, "scan.read", dir)
		return r.finish()
	}
	source := kind.Layout.Path(dir, slug)
	r.Summary.FilesDiscovered = 1
	r.Summary.FilesChecked = 1
	findings, err := kind.Inspect(source)
	if err != nil {
		r.record(err, "scan.read", source)
		return r.finish()
	}
	r.Diagnostics = append(r.Diagnostics, findings...)
	if kind.Context != nil {
		r.inspectContexts([]contextRequest{{kind.Context, dir}}, map[string]bool{source: true})
	}
	return r.finish()
}
