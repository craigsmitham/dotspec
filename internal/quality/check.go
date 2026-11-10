package quality

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Requirement) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	for _, field := range []struct{ rule, name, value, message string }{
		{"quality.applies-to", "Applies to", spec.AppliesTo, "Identify the subject and applicable conditions, workload, or population."},
		{"quality.criterion", "Criterion", spec.Criterion, "State one coherent, assessable quality obligation."},
		{"quality.measurement", "Measurement", spec.Measurement, "Describe how to assess the criterion and retain the evidence."},
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
	return suggestions
}
