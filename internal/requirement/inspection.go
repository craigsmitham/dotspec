package requirement

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

func CheckType() check.Type {
	kind := speccollection.CheckType(Layout, readFrom, func(spec *Requirement) *specfile.ReferenceDocument { return spec.References }, Check)
	kind.InspectContext = inspectAttributes
	return kind
}

func inspectAttributes(files *specfile.SourceSet, path string, sources []string) ([]diagnostic.Diagnostic, error) {
	if len(sources) == 0 {
		return nil, nil
	}
	sys, err := system.ReadFrom(files, sources[len(sources)-1])
	if err != nil {
		return nil, err
	}
	spec, err := readFrom(files, path)
	if err != nil {
		return nil, err
	}
	return CheckAttributes(spec, sys), nil
}
