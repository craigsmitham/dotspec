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

// CheckContext preserves the same nearest-system lookup as Show and List,
// including ignored systems or ancestors outside the selected subtree.
var CheckContext = &check.ContextProvider{Inspect: inspect, Resolve: func(owner string) ([]string, error) {
	spec, err := Enclosing(owner)
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
}}

func readReferences(path string) (*specfile.ReferenceDocument, error) {
	spec, err := Read(path)
	if err != nil {
		return nil, err
	}
	return spec.References, nil
}
