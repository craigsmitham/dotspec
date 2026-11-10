package concept

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// Check suggests missing content. It does not judge the design's coherence,
// the adequacy of its behavior, or whether its purpose has been achieved.
func Check(spec *Concept) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	for _, field := range []struct{ rule, name, value, message string }{
		{"concept.purpose", "Purpose", spec.Purpose, "Describe the benefit this concept provides and why it exists."},
		{"concept.state", "State", spec.State, "Describe or reference the abstract information needed to explain the behavior."},
		{"concept.actions", "Actions", spec.Actions, "Describe or reference the actions, their conditions, and effects."},
		{"concept.operational-principle", "Operational principle", spec.OperationalPrinciple, "Tell a characteristic story showing how the behavior fulfills the purpose."},
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
