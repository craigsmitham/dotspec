package check

import (
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Type supplies the policy needed to inspect one specification kind. Callers
// explicitly compose the supported types; the first matching layout wins.
// Inspect validates placement, reads the file, and returns writing advice only
// after structural validation succeeds. Context is optional.
type Type struct {
	// References reads validated targets and outgoing links without following them.
	References   func(path string) (*specfile.ReferenceDocument, error)
	Layout       specfile.Layout
	ValidateSlug func(string) error
	Inspect      func(path string) ([]diagnostic.Diagnostic, error)
	Context      *ContextProvider
	// InspectContext checks cross-source relationships after successful context
	// resolution. Sources are exactly those returned by Context, never ancestors
	// inferred by the checker. It is skipped when context cannot be read.
	InspectContext func(path string, sources []string) ([]diagnostic.Diagnostic, error)
}

// ContextProvider inspects additional sources needed by a type. Resolve returns
// their absolute paths, including attempted reads that failed, alongside any
// error. Paths identify files, not directories where lookup failed.
// Share a provider between types to resolve it once per owner in a project scan.
type ContextProvider struct {
	Resolve func(owner string) ([]string, error)
	// Inspect optionally contributes findings from each resolved context file.
	Inspect func(path string) ([]diagnostic.Diagnostic, error)
}

type contextKey struct {
	provider *ContextProvider
	owner    string
}

type contextRequest struct {
	key     contextKey
	source  string
	inspect func(string, []string) ([]diagnostic.Diagnostic, error)
}

func (r *Report) inspectContexts(requests []contextRequest, discovered map[string]bool) {
	type resolution struct {
		paths []string
		err   error
	}
	resolved := make(map[contextKey]resolution)
	counted := make(map[string]bool)
	type inspectionKey struct {
		provider *ContextProvider
		path     string
	}
	inspected := make(map[inspectionKey]error)
	for _, request := range requests {
		key := request.key
		result, found := resolved[key]
		if !found {
			result.paths, result.err = key.provider.Resolve(key.owner)
			resolved[key] = result
			if result.err != nil {
				r.record(result.err, "scan.context", key.owner)
			}
		}
		contextOK := result.err == nil
		for _, path := range result.paths {
			r.collectFile(path, nil)
			if path != "" && !discovered[path] && !counted[path] {
				counted[path] = true
				r.ContextFiles = append(r.ContextFiles, path)
				r.Summary.ContextFilesChecked++
			}
			if path != "" && contextOK && key.provider.Inspect != nil {
				inspection := inspectionKey{key.provider, path}
				err, found := inspected[inspection]
				if !found {
					var findings []diagnostic.Diagnostic
					findings, err = key.provider.Inspect(path)
					r.Diagnostics = append(r.Diagnostics, findings...)
					if err != nil {
						r.record(err, "scan.context", path)
					}
					inspected[inspection] = err
				}
				contextOK = err == nil
			}
		}
		if contextOK && request.inspect != nil {
			findings, err := request.inspect(request.source, result.paths)
			r.Diagnostics = append(r.Diagnostics, findings...)
			if err != nil {
				r.record(err, "scan.context", request.source)
			}
		}
	}
}
