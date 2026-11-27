package concept

import (
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Concept) []diagnostic.Diagnostic {
	blank := specfile.NewBlankSections(spec.Path, spec.Location, spec.SectionLocations)
	blank.Suggest("concept.purpose", "Purpose", spec.Purpose, "Describe the benefit this concept provides and why it exists.")
	blank.Suggest("concept.state", "State", spec.State, "Describe or reference the abstract information needed to explain the behavior.")
	blank.Suggest("concept.actions", "Actions", spec.Actions, "Describe or reference the actions, their conditions, and effects.")
	blank.Suggest("concept.operational-principle", "Operational principles", spec.OperationalPrinciples, "Tell a characteristic story showing how the behavior fulfills the purpose.")
	return blank.Advice
}
