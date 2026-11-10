package quality

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/system"
)

func CheckType() check.Type {
	return check.Type{Layout: Layout, ValidateSlug: ValidateSlug, Context: system.CheckContext, InspectContext: inspectAttributes,
		Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
			spec, err := Read(path)
			if err != nil {
				return nil, err
			}
			return Check(spec), nil
		}}
}

func inspectAttributes(path string, sources []string) ([]diagnostic.Diagnostic, error) {
	if len(sources) == 0 {
		return nil, nil
	}
	sys, err := system.Read(sources[0])
	if err != nil {
		return nil, err
	}
	spec, err := Read(path)
	if err != nil {
		return nil, err
	}
	return CheckAttributes(spec, sys), nil
}
