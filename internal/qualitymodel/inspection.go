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
			findings = append(findings, diagnostic.Advice("quality.definition", "Define what this quality attribute means.", a.Location))
		}
	}
	return findings
}

func CheckType() check.Type {
	return check.ReadType(Layout, readFrom, func(m *QualityModel) *specfile.ReferenceDocument { return m.References }, Check)
}
