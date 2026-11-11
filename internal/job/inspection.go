package job

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

func CheckType() check.Type { return CheckTypeWithPreset("general") }

func CheckTypeWithPreset(preset string) check.Type {
	return check.Type{Layout: Layout, References: readReferences, ValidateSlug: ValidateSlug, Context: system.CheckContext,
		Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
			spec, err := Read(path)
			if err != nil {
				return nil, err
			}
			return CheckWithPreset(spec, preset), nil
		}}
}

func readReferences(path string) (*specfile.ReferenceDocument, error) {
	spec, err := Read(path)
	if err != nil {
		return nil, err
	}
	return spec.References, nil
}
