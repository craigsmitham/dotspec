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
	Layout       specfile.Layout
	ValidateSlug func(string) error
	Inspect      func(path string) ([]diagnostic.Diagnostic, error)
	Context      *ContextProvider
}

// ContextProvider inspects additional sources needed by a type. Resolve returns
// their absolute paths, including attempted reads that failed, alongside any
// error. Paths identify files, not directories where lookup failed.
// Share a provider between types to resolve it once per owner in a project scan.
type ContextProvider struct {
	Resolve func(owner string) ([]string, error)
}

type contextRequest struct {
	provider *ContextProvider
	owner    string
}

func (r *Report) inspectContexts(requests []contextRequest, discovered map[string]bool) {
	resolved := make(map[contextRequest]bool)
	counted := make(map[string]bool)
	for _, request := range requests {
		if resolved[request] {
			continue
		}
		resolved[request] = true
		paths, err := request.provider.Resolve(request.owner)
		if err != nil {
			r.record(err, "scan.context", request.owner)
		}
		for _, path := range paths {
			if path != "" && !discovered[path] && !counted[path] {
				counted[path] = true
				r.ContextFiles = append(r.ContextFiles, path)
				r.Summary.ContextFilesChecked++
			}
		}
	}
}
