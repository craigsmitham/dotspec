package cli

import (
	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func newBusinessRuleCommand() *cobra.Command {
	return specKind[businessrule.BusinessRule, *businessrule.Group]{
		short:  "Describe business constraints independently of their enforcement",
		plural: "business rules", checkShort: "Check structure and suggest missing business-rule content", key: "business_rule",
		guidance: "Describe the business constraint, applicability, rationale, authority, and exceptions as they become known.\nAuthoring guide: " + authoringPage("business-rule"),
		layout:   businessrule.Layout, checkType: businessrule.CheckType, init: businessrule.Init, show: businessrule.Show, list: businessrule.List,
		source: func(spec *businessrule.BusinessRule) specSource {
			return specSource{spec.Name, spec.Markdown, spec.Path, spec.System}
		},
		members: func(group *businessrule.Group) (*system.System, []*businessrule.BusinessRule) {
			return group.System, group.BusinessRules
		},
	}.command()
}
