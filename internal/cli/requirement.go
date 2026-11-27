package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func labelList(labels []string) string {
	if len(labels) == 0 {
		return "None"
	}
	return strings.Join(labels, ", ")
}

type requirementListResult struct {
	SchemaVersion int                  `json:"schema_version"`
	Groups        []*requirement.Group `json:"groups"`
	Attribute     string               `json:"attribute,omitempty"`
}

func newRequirementCommand() *cobra.Command {
	root := group("requirement", "Specify required results, rules, properties, and restrictions")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a requirement specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover requirements, grouped by associated system"},
		{"check", "Check requirement specifications and attribute references"},
	} {
		cmd := newSpecAction(requirement.Layout, action.name, action.short)
		if action.name == "list" {
			cmd.Flags().String("attribute", "", "Include requirements with this attribute or a sub-attribute in their system's model")
		}
		if action.name == "init" {
			cmd.Flags().String("attribute", "", "Write this governing system attribute ID in quality_attributes.")
			cmd.Flags().String("preset", "general", "Authoring template: general, or quality to prompt for an explicit criterion")
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			preset := "general"
			if action.name == "init" {
				preset, _ = cmd.Flags().GetString("preset")
				if err := requirement.ValidatePreset(preset); err != nil {
					return err
				}
			}
			if action.name == "list" {
				attribute, err := cmd.Flags().GetString("attribute")
				if err != nil {
					return err
				}
				if cmd.Flags().Changed("attribute") {
					if attribute == "" {
						return fmt.Errorf("attribute selector must not be empty")
					}
					var ids specfile.Identifiers
					if err := ids.Register(attribute, "selector", diagnostic.Location{}); err != nil {
						return fmt.Errorf("invalid attribute selector: %s", err)
					}
				}
				return listRequirement(path, format, attribute, cmd.OutOrStdout())
			}
			if err := validateSelector(requirement.Layout, action.name, args[0]); err != nil {
				return err
			}
			attribute := ""
			if action.name == "init" {
				attribute, _ = cmd.Flags().GetString("attribute")
				if cmd.Flags().Changed("attribute") && attribute == "" {
					return fmt.Errorf("attribute selector must not be empty")
				}
			}
			return runRequirement(action.name, path, args[0], format, preset, attribute, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runRequirement(command, path, slug, format, preset, attribute string, out io.Writer) error {
	if command == "init" {
		target, err := requirement.InitWithAttribute(path, slug, preset, attribute)
		if err != nil {
			return err
		}
		return writeCreated(out, format, target, "Describe the scope, statement, any criterion, rationale, and verification as they become known.\nAuthoring guide: "+authoringPage("requirement"))
	}
	if command == "check" {
		return emitFocusedCheck(out, format, path, slug, requirement.CheckType())
	}
	spec, err := requirement.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, readResult{"requirement", spec})
	}
	attributeSource := "None (labels are unbound)"
	if spec.Classification.SystemPath != "" {
		attributeSource = displayPath(spec.Classification.SystemPath)
	}
	details := fmt.Sprintf("Quality attributes: %s\nAncestor quality attributes: %s\nUndefined quality attributes: %s\nAttribute declarations: %s\n", labelList(spec.QualityAttributes), labelList(spec.Classification.Ancestors), labelList(spec.Classification.Undefined), attributeSource)
	return writeShown(out, spec.Markdown, details, spec.System, spec.Path)
}

func listRequirement(path, format, attribute string, out io.Writer) error {
	groups, err := requirement.List(path)
	if err != nil {
		return err
	}
	result := requirementListResult{readSchemaVersion, nonnil(requirement.Filter(groups, attribute)), attribute}
	if format == "json" {
		return writeJSON(out, result)
	}
	heading := "Requirements"
	if attribute != "" {
		heading += fmt.Sprintf("\nAttribute: %s (including sub-attributes)", attribute)
	}
	owner := func(group *requirement.Group) *system.System { return group.System }
	return writeGroupedList(out, "requirements", heading, result.Groups, owner, func(output *strings.Builder, group *requirement.Group) {
		for _, spec := range group.Requirements {
			writeListEntry(output, spec.Name, spec.Path)
			fmt.Fprintf(output, "    Quality attributes: %s\n", labelList(spec.QualityAttributes))
		}
	})
}
