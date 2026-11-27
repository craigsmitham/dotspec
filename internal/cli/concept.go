package cli

import (
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func newConceptCommand() *cobra.Command {
	return specKind[concept.Concept, *concept.Group]{
		short:  "Describe purposeful product behavior through state and actions",
		plural: "concepts", checkShort: "Check structure and suggest missing concept content", key: "concept",
		guidance: "Describe the purpose, state, actions, and operational principle as they become known.\nAuthoring guide: " + authoringPage("concept"),
		layout:   concept.Layout, checkType: concept.CheckType, init: concept.Init, show: concept.Show, list: concept.List,
		source: func(spec *concept.Concept) specSource {
			return specSource{spec.Name, spec.Markdown, spec.Path, spec.System}
		},
		members: func(group *concept.Group) (*system.System, []*concept.Concept) { return group.System, group.Concepts },
	}.command()
}
