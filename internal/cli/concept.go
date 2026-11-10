package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/spf13/cobra"
)

type conceptResult struct {
	SchemaVersion int              `json:"schema_version"`
	Concept       *concept.Concept `json:"concept"`
}

type conceptsResult struct {
	SchemaVersion int              `json:"schema_version"`
	Groups        []*concept.Group `json:"groups"`
}

func newConceptCommand() *cobra.Command {
	root := group("concept", "Describe purposeful product behavior through state and actions")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a concept specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover concepts, grouped by associated system"},
		{"check", "Check structure and suggest missing concept content"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short,
			Example: "dotspec concept " + action.name + " job-map --path .", Args: cobra.ExactArgs(1)}
		if action.name == "list" {
			cmd.Use, cmd.Example, cmd.Args = "list", "dotspec concept list --path .", cobra.NoArgs
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
				return listConcepts(path, format, cmd.OutOrStdout())
			}
			if err := concept.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runConcept(action.name, path, args[0], format, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runConcept(command, path, slug, format string, out io.Writer) error {
	if command == "init" {
		target, err := concept.Init(path, slug)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nDescribe the purpose, state, actions, and operational principle as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.Focused(path, slug, concept.CheckType(), supportedTypes()...), format)
	}
	spec, err := concept.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, conceptResult{1, spec})
	}
	_, err = fmt.Fprintf(out, "%s\n\nSystem: %s\nSource: %s\n", spec.Markdown, systemLabel(spec.System), displayPath(spec.Path))
	return err
}

func listConcepts(path, format string, out io.Writer) error {
	groups, err := concept.List(path)
	if err != nil {
		return err
	}
	result := conceptsResult{1, nonnil(groups)}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No concepts found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Concepts\n")
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.Concepts {
			fmt.Fprintf(&output, "  %s  (%s)\n", spec.Name, displayPath(spec.Path))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}
