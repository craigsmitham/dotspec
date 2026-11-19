package system

import (
	"errors"
	"os"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func CheckType() check.Type {
	return check.Type{Layout: Layout, References: readReferences, Inspect: inspect}
}

func inspect(path string) ([]diagnostic.Diagnostic, error) {
	spec, err := Read(path)
	if err != nil {
		return nil, err
	}
	return Check(spec), nil
}

// CheckContext preserves the same explicit selection or directory default as
// Show and List, including ignored systems outside the selected subtree.
var CheckContext = &check.ContextProvider{Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
	if ComponentLayout.Match(path) {
		_, err := ReadComponent(path)
		return nil, err
	}
	return inspect(path)
},
	Resolve: SubjectChain,
	ResolveSource: func(owner, source string) ([]string, error) {
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}
		_, explicit, err := specfile.SpecificationSystem(data)
		if err != nil {
			return nil, diagnostic.WithPath(err, source)
		}
		if explicit {
			paths, err := contextPaths(Associated(owner, source))
			if err != nil {
				return paths, err
			}
			// An explicit selection can override an invalid directory system. Only a local
			// component contributes additional ownership context.
			boundary, err := specfile.ContextRoot(owner)
			if err != nil {
				return nil, err
			}
			local, lookupErr := nearestSubject(owner, boundary)
			if lookupErr == nil && local != "" && ComponentLayout.Match(local) {
				paths = append([]string{local}, paths...)
			}
			return paths, nil
		}
		return SubjectChain(owner)
	},
}

func contextPaths(spec *System, err error) ([]string, error) {
	if err != nil {
		path := ""
		var finding *diagnostic.Diagnostic
		var pathError *os.PathError
		if errors.As(err, &finding) {
			path = finding.Location.Path
		} else if errors.As(err, &pathError) {
			path = pathError.Path
		}
		if path != "" && Layout.Match(path) {
			return []string{path}, err
		}
		return nil, err
	}
	if spec == nil {
		return nil, nil
	}
	return []string{spec.Path}, nil
}

func readReferences(path string) (*specfile.ReferenceDocument, error) {
	spec, err := Read(path)
	if err != nil {
		return nil, err
	}
	return spec.References, nil
}
