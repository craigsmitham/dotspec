package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

type requirementResult struct {
	SchemaVersion int                      `json:"schema_version"`
	Requirement   *requirement.Requirement `json:"requirement"`
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
		{"list", "Discover general and quality requirements, grouped by associated system"},
		{"check", "Check requirement specifications and attribute references"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short,
			Example: "dotspec requirement " + action.name + " reject-existing-target --path .", Args: cobra.ExactArgs(1)}
		if action.name == "list" {
			cmd.Use, cmd.Example, cmd.Args = "list", "dotspec requirement list --path .", cobra.NoArgs
			cmd.Flags().String("attribute", "", "Include requirements with this attribute or a sub-attribute in their system's model")
		}
		if action.name == "check" {
			cmd.Long = action.short + ".\n\n" + exitPolicy
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
			if err := requirement.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runRequirement(action.name, path, args[0], format, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runRequirement(command, path, slug, format string, out io.Writer) error {
	if command == "init" {
		target, err := requirement.Init(path, slug)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nDescribe the scope, statement, rationale, and verification as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.Focused(path, slug, requirement.CheckType(), supportedTypes()...), format)
	}
	spec, err := requirement.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, requirementResult{1, spec})
	}
	attributeSource := "None (labels are unbound)"
	if spec.Classification.SystemPath != "" {
		attributeSource = displayPath(spec.Classification.SystemPath)
	}
	_, err = fmt.Fprintf(out, "%s\n\nAttributes: %s\nAncestor attributes: %s\nUndefined attributes: %s\nAttribute declarations: %s\nSystem: %s\nSource: %s\n", spec.Markdown, labelList(spec.Attributes), labelList(spec.Classification.Ancestors), labelList(spec.Classification.Undefined), attributeSource, systemLabel(spec.System), displayPath(spec.Path))
	return err
}

func listRequirement(path, format, attribute string, out io.Writer) error {
	groups, err := requirement.List(path)
	if err != nil {
		return err
	}
	result := requirementListResult{1, nonnil(requirement.Filter(groups, attribute)), attribute}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No requirements found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Requirements\n")
	if attribute != "" {
		fmt.Fprintf(&output, "Attribute: %s (including sub-attributes)\n", attribute)
	}
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.Requirements {
			fmt.Fprintf(&output, "  %s [%s]  (%s)\n    Attributes: %s\n", spec.Name, spec.Type, displayPath(spec.Path), labelList(spec.Attributes))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}
