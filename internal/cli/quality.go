package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

type qualityResult struct {
	SchemaVersion int                  `json:"schema_version"`
	Quality       *quality.Requirement `json:"quality"`
}

type qualityListResult struct {
	SchemaVersion int              `json:"schema_version"`
	Groups        []*quality.Group `json:"groups"`
	Attribute     string           `json:"attribute,omitempty"`
}

func newQualityCommand() *cobra.Command {
	root := group("quality", "Specify assessable quality obligations and their attributes")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a quality specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover quality requirements, grouped by associated system"},
		{"check", "Check quality specifications and attribute references"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short,
			Example: "dotspec quality " + action.name + " diagnostic-severity --path .", Args: cobra.ExactArgs(1)}
		if action.name == "list" {
			cmd.Use, cmd.Example, cmd.Args = "list", "dotspec quality list --path .", cobra.NoArgs
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
				return listQuality(path, format, attribute, cmd.OutOrStdout())
			}
			if err := quality.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runQuality(action.name, path, args[0], format, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runQuality(command, path, slug, format string, out io.Writer) error {
	if command == "init" {
		target, err := quality.Init(path, slug)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nDescribe the scope, criterion, measurement, and rationale as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.Focused(path, slug, quality.CheckType(), supportedTypes()...), format)
	}
	spec, err := quality.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, qualityResult{1, spec})
	}
	attributeSource := "None (labels are unbound)"
	if spec.Classification.SystemPath != "" {
		attributeSource = displayPath(spec.Classification.SystemPath)
	}
	_, err = fmt.Fprintf(out, "%s\n\nQuality attributes: %s\nAncestor quality attributes: %s\nUndefined quality attributes: %s\nAttribute declarations: %s\nSystem: %s\nSource: %s\n", spec.Markdown, labelList(spec.QualityAttributes), labelList(spec.Classification.Ancestors), labelList(spec.Classification.Undefined), attributeSource, systemLabel(spec.System), displayPath(spec.Path))
	return err
}

func listQuality(path, format, attribute string, out io.Writer) error {
	groups, err := quality.List(path)
	if err != nil {
		return err
	}
	result := qualityListResult{1, nonnil(quality.Filter(groups, attribute)), attribute}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No quality requirements found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Quality requirements\n")
	if attribute != "" {
		fmt.Fprintf(&output, "Attribute: %s (including sub-attributes)\n", attribute)
	}
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.Requirements {
			fmt.Fprintf(&output, "  %s  (%s)\n    Quality attributes: %s\n", spec.Name, displayPath(spec.Path), labelList(spec.QualityAttributes))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}

func labelList(labels []string) string {
	if len(labels) == 0 {
		return "None"
	}
	return strings.Join(labels, ", ")
}
