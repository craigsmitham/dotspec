package quality

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/obligation"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Requirement) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	for _, field := range []struct{ rule, name, value, message string }{
		{"quality.applies-to", "Applies to", spec.Applicability, "Identify the subject and applicable conditions, workload, or population."},
		{"quality.statement", "Statement", spec.Statement, "State one coherent, assessable quality obligation."},
		{"quality.verification-plan", "Verification plan", spec.VerificationPlan + spec.MeasurementPlan, "Record assessment guidance for checking a realization against the obligation; a plan is not a verification result."},
		{"quality.rationale", "Rationale", spec.Rationale, "Explain why this quality matters and any relevant tradeoffs."},
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
		section := "Measurement plan"
		if strings.TrimSpace(spec.MeasurementPlan) == "" {
			section = "Criterion"
		}
		location := spec.SectionLocations[section]
		location.Path, location.Element = spec.Path, section
		suggestions = append(suggestions, diagnostic.Advice("quality.measurement-basis", obligation.MeasurementMessage(missing), location))
	}
	return suggestions
}
