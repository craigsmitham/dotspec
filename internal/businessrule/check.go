package businessrule

import (
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Check suggests missing content without verifying authority or compliance.
func Check(spec *BusinessRule) []diagnostic.Diagnostic {
	blank := specfile.NewBlankSections(spec.Path, spec.Location, spec.SectionLocations)
	blank.Suggest("business-rule.statement", "Statement", spec.Statement, "State the business constraint independently of enforcement.")
	blank.Suggest("business-rule.applies-to", "Applies to", spec.Applicability, "Describe the business activity, subjects, and applicable conditions.")
	blank.Suggest("business-rule.rationale", "Rationale", spec.Rationale, "Explain why the business rule exists.")
	blank.Suggest("business-rule.authority", "Authority", spec.Authority, "Identify the policy, decision, agreement, or other source establishing the rule.")
	return blank.Advice
}
