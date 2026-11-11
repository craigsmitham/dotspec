package businessrule

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// Check suggests missing content without verifying authority or compliance.
func Check(spec *BusinessRule) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	for _, field := range []struct{ rule, name, value, message string }{
		{"business-rule.statement", "Statement", spec.Statement, "State the business constraint independently of enforcement."},
		{"business-rule.applies-to", "Applies to", spec.Applicability, "Describe the business activity, subjects, and applicable conditions."},
		{"business-rule.rationale", "Rationale", spec.Rationale, "Explain why the business rule exists."},
		{"business-rule.authority", "Authority", spec.Authority, "Identify the policy, decision, agreement, or other source establishing the rule."},
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
