package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/spf13/cobra"
)

type businessruleResult struct {
	SchemaVersion int                        `json:"schema_version"`
	BusinessRule  *businessrule.BusinessRule `json:"business_rule"`
}

type businessRulesResult struct {
	SchemaVersion int                   `json:"schema_version"`
	Groups        []*businessrule.Group `json:"groups"`
}

func newBusinessRuleCommand() *cobra.Command {
	root := group("business-rule", "Describe business constraints independently of their enforcement")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a business-rule specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover business rules, grouped by associated system"},
		{"check", "Check structure and suggest missing business-rule content"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short,
			Example: "dotspec business-rule " + action.name + " purchase-approval --path .", Args: cobra.ExactArgs(1)}
		if action.name == "list" {
			cmd.Use, cmd.Example, cmd.Args = "list", "dotspec business-rule list --path .", cobra.NoArgs
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
				return listBusinessRules(path, format, cmd.OutOrStdout())
			}
			if err := businessrule.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runBusinessRule(action.name, path, args[0], format, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runBusinessRule(command, path, slug, format string, out io.Writer) error {
	if command == "init" {
		target, err := businessrule.Init(path, slug)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nDescribe the business constraint, applicability, rationale, authority, and exceptions as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.Focused(path, slug, businessrule.CheckType(), supportedTypes()...), format)
	}
	spec, err := businessrule.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, businessruleResult{1, spec})
	}
	_, err = fmt.Fprintf(out, "%s\n\nSystem: %s\nSource: %s\n", spec.Markdown, systemLabel(spec.System), displayPath(spec.Path))
	return err
}

func listBusinessRules(path, format string, out io.Writer) error {
	groups, err := businessrule.List(path)
	if err != nil {
		return err
	}
	result := businessRulesResult{1, nonnil(groups)}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No business rules found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Business rules\n")
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.BusinessRules {
			fmt.Fprintf(&output, "  %s  (%s)\n", spec.Name, displayPath(spec.Path))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}
