package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
)

type useCaseResult struct {
	SchemaVersion int              `json:"schema_version"`
	UseCase       *usecase.UseCase `json:"use_case"`
}

type useCasesResult struct {
	SchemaVersion int              `json:"schema_version"`
	Groups        []*usecase.Group `json:"groups"`
}

func newUseCaseCommand() *cobra.Command {
	root := group("use-case", "Describe actor goals and interactions with a system")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a use-case specification without overwriting an existing file"},
		{"show", "Show authored Markdown and the associated system"},
		{"list", "Discover use cases, grouped by associated system"},
		{"check", "Check structure and suggest missing use-case content"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short,
			Example: "dotspec use-case " + action.name + " check-specifications --path .", Args: cobra.ExactArgs(1)}
		if action.name == "list" {
			cmd.Use, cmd.Example, cmd.Args = "list", "dotspec use-case list --path .", cobra.NoArgs
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
				return listUseCases(path, format, cmd.OutOrStdout())
			}
			if err := usecase.ValidateSlug(args[0]); err != nil {
				return err
			}
			return runUseCase(action.name, path, args[0], format, cmd.OutOrStdout())
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runUseCase(command, path, slug, format string, out io.Writer) error {
	if command == "init" {
		target, err := usecase.Init(path, slug)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.\n", displayPath(target))
		return err
	}
	if command == "check" {
		return emitCheckReport(out, check.FocusedUseCase(path, slug), format)
	}
	spec, err := usecase.Show(path, slug)
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, useCaseResult{1, spec})
	}
	_, err = fmt.Fprintf(out, "%s\n\nSystem: %s\nSource: %s\n", spec.Markdown, systemLabel(spec.System), displayPath(spec.Path))
	return err
}

func listUseCases(path, format string, out io.Writer) error {
	groups, err := usecase.List(path)
	if err != nil {
		return err
	}
	result := useCasesResult{1, nonnil(groups)}
	if format == "json" {
		return writeJSON(out, result)
	}
	if len(result.Groups) == 0 {
		_, err = fmt.Fprintln(out, "No use cases found in the selected subtree.")
		return err
	}
	var output strings.Builder
	output.WriteString("Use cases\n")
	for _, group := range result.Groups {
		fmt.Fprintf(&output, "\nSystem: %s\n", systemLabel(group.System))
		for _, spec := range group.UseCases {
			fmt.Fprintf(&output, "  %s  (%s)\n", spec.Name, displayPath(spec.Path))
		}
	}
	_, err = io.WriteString(out, output.String())
	return err
}
