package requirement

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Requirement) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	// Quality-specific wording follows the classification, not a separate type.
	statement, rationale := "State one required result, rule, property, or restriction.", "Explain why this requirement is needed and reference its source or governing goal."
	if len(spec.QualityAttributes) > 0 {
		statement, rationale = "State one coherent, assessable quality requirement.", "Explain why this quality matters and any relevant tradeoffs."
	}
	for _, field := range []struct{ rule, name, value, message string }{
		{"requirement.applies-to", "Applies to", spec.Applicability, "Identify the subject and applicable conditions, workload, or population."},
		{"requirement.statement", "Statement", spec.Statement, statement},
		{"requirement.rationale", "Rationale", spec.Rationale, rationale},
	} {
		if strings.TrimSpace(field.value) != "" {
			continue
		}
		location, found := spec.SectionLocations[field.name]
		if !found {
			location = spec.Location
		}
		location.Path, location.Element = spec.Path, field.name
		suggestions = append(suggestions, diagnostic.Advice(field.rule, field.message, location))
	}
	if missing := spec.MeasurementGuidance.Missing(); len(missing) > 0 {
		location := spec.SectionLocations["Criterion"]
		location.Path, location.Element = spec.Path, "Criterion"
		suggestions = append(suggestions, diagnostic.Advice("requirement.measurement-basis", MeasurementMessage(missing), location))
	}
	return suggestions
}
