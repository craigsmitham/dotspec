package qualitymodel

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func Check(spec *QualityModel) []diagnostic.Diagnostic {
	findings := []diagnostic.Diagnostic{}
	loc := spec.Location
	loc.Path = spec.Path
	if spec.Description == "" {
		findings = append(findings, diagnostic.Advice("quality-model.description", "Describe the purpose and scope of this model.", loc))
	}
	if len(spec.QualityAttributes) == 0 {
		findings = append(findings, diagnostic.Advice("quality-model.attributes", "Declare quality attributes when the model is ready to populate a system.", loc))
	}
	for _, a := range spec.QualityAttributes {
		if a.Definition == "" {
			loc := a.Location
			loc.Path = spec.Path
			findings = append(findings, diagnostic.Advice("quality.definition", "Define what this quality attribute means.", loc))
		}
	}
	return findings
}

func CheckType() check.Type {
	return check.Type{Layout: Layout, ValidateSlug: ValidateSlug,
		Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
			m, err := Read(path)
			if err != nil {
				return nil, err
			}
			return Check(m), nil
		},
		References: func(path string) (*specfile.ReferenceDocument, error) {
			m, err := Read(path)
			if err != nil {
				return nil, err
			}
			return m.References, nil
		},
	}
}
