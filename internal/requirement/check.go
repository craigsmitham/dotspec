package requirement

import (
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Requirement) []diagnostic.Diagnostic {
	// Quality-specific wording follows the classification, not a separate type.
	statement, rationale := "State one required result, rule, property, or restriction.", "Explain why this requirement is needed and reference its source or governing goal."
	if len(spec.QualityAttributes) > 0 {
		statement, rationale = "State one coherent, assessable quality requirement.", "Explain why this quality matters and any relevant tradeoffs."
	}
	blank := specfile.NewBlankSections(spec.Path, spec.Location, spec.SectionLocations)
	blank.Suggest("requirement.applies-to", "Applies to", spec.Applicability, "Identify the subject and applicable conditions, workload, or population.")
	blank.Suggest("requirement.statement", "Statement", spec.Statement, statement)
	blank.Suggest("requirement.rationale", "Rationale", spec.Rationale, rationale)
	suggestions := blank.Advice
	if missing := spec.MeasurementGuidance.Missing(); len(missing) > 0 {
		location := spec.SectionLocations["Criterion"]
		location.Path, location.Element = spec.Path, "Criterion"
		suggestions = append(suggestions, diagnostic.Advice("requirement.measurement-basis", MeasurementMessage(missing), location))
	}
	return suggestions
}
