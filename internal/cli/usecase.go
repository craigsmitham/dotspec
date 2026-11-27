package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
)

func newUseCaseCommand() *cobra.Command {
	template, overview := "brief", false
	return specKind[usecase.UseCase, *usecase.Group]{
		short:  "Describe actor goals and interactions with a system",
		plural: "use cases", checkShort: "Check structure and suggest missing use-case content", key: "use_case",
		guidance: "Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.\nAuthoring guide: " + authoringPage("use-case"),
		layout:   usecase.Layout, checkType: usecase.CheckType, show: usecase.Show, list: usecase.List,
		init: func(path, slug string) (string, error) { return usecase.InitTemplate(path, slug, template) },
		source: func(spec *usecase.UseCase) specSource {
			return specSource{spec.Name, spec.Markdown, spec.Path, spec.System}
		},
		members: func(group *usecase.Group) (*system.System, []*usecase.UseCase) { return group.System, group.UseCases },
		flags: func(action string, cmd *cobra.Command) {
			switch action {
			case "init":
				cmd.Flags().StringVar(&template, "template", "brief", "Writing prompts: brief or expanded (same document format)")
				cmd.Example += "\ndotspec use-case init review-change --template expanded --path ."
			case "list":
				cmd.Flags().BoolVar(&overview, "overview", false, "Show actor, goal, level, and stated boundary in text; JSON retains the full records")
				cmd.Example += "\ndotspec use-case list --overview --path ."
			}
		},
		writeEntry: func(output *strings.Builder, spec *usecase.UseCase) {
			if !overview {
				writeListEntry(output, spec.Name, spec.Path)
				return
			}
			goal := spec.Goal
			if strings.TrimSpace(goal) == "" {
				goal = spec.Name
			}
			fmt.Fprintf(output, "\n  %s (%s)\n    Actor: %s\n    Goal: %s\n    Level: %s\n    Stated boundary: %s\n", spec.Name, displayPath(spec.Path), overviewCell(spec.PrimaryActor), overviewCell(goal), overviewCell(spec.GoalLevel), overviewCell(spec.BoundaryDescription))
		},
	}.command()
}

// Collapse layout whitespace only; never infer an actor or boundary from ownership.
func overviewCell(value string) string {
	if strings.TrimSpace(value) == "" {
		return "(unspecified)"
	}
	return strings.Join(strings.Fields(value), " ")
}
