package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type systemResult struct {
	SchemaVersion   int            `json:"schema_version"`
	System          *system.System `json:"system"`
	EnclosingSystem *system.System `json:"enclosing_system"`
}
type systemsResult struct {
	SchemaVersion int            `json:"schema_version"`
	Systems       []*system.Node `json:"systems"`
}
type createdResult struct {
	SchemaVersion int    `json:"schema_version"`
	Path          string `json:"path"`
}

func newSystemCommand() *cobra.Command {
	root := group("system", "Describe systems and discover their directory context")
	for _, action := range []struct{ name, short string }{
		{"init", "Create a named system specification"},
		{"show", "Describe the nearest system and its enclosing system document"},
		{"list", "Display systems in the selected subtree"},
	} {
		cmd := &cobra.Command{Use: action.name, Short: action.short, Args: cobra.NoArgs,
			Example: "dotspec system " + action.name + " --path .",
			RunE: func(cmd *cobra.Command, args []string) error {
				path, err := projectPath(cmd)
				if err != nil {
					return err
				}
				format, err := outputFormat(cmd)
				if err != nil {
					return err
				}
				return runSystem(cmd.Name(), path, format, cmd.OutOrStdout(), args...)
			},
		}
		if action.name == "show" {
			cmd.Use = "show [source-path]"
			cmd.Args = cobra.MaximumNArgs(1)
		}
		if action.name == "init" {
			cmd.Use = "init [name]"
			cmd.Args = cobra.MaximumNArgs(1)
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runSystem(command, path, format string, out io.Writer, args ...string) error {
	switch command {
	case "init":
		name := "New system"
		if len(args) > 0 {
			name = args[0]
		}
		target, err := system.InitNamed(path, name)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit Description, Purpose, and Scope. Keep the filename in sync if you rename the system.\nAuthoring guide: dotspec guide system\n", displayPath(target))
		return err
	case "show":
		if len(args) > 0 {
			file, err := system.Layout.SourcePath(path, args[0])
			if err != nil {
				return err
			}
			path = specfile.DeclarationDir(file)
		}
		spec, parent, err := system.Show(path)
		if err != nil {
			return err
		}
		result := systemResult{1, spec, parent}
		if format == "json" {
			return writeJSON(out, result)
		}
		parentText := "None"
		if result.EnclosingSystem != nil {
			parentText = fmt.Sprintf("%s (%s)", result.EnclosingSystem.Name, displayPath(result.EnclosingSystem.Path))
		}
		modelText := ""
		for _, field := range []struct{ name, value string }{{"Goals", spec.Goals}, {"Environment", spec.Environment}, {"Stakeholders", spec.Stakeholders}, {"Assumptions", spec.Assumptions}, {"Constraints", spec.Constraints}, {"Glossary", spec.Glossary}, {"Part of", spec.PartOf}, {"Components", spec.Components}, {"Concept uses", spec.ConceptUses}, {"Synchronizations", spec.Synchronizations}} {
			if field.value != "" {
				modelText += "\n\n" + field.name + "\n" + indent(field.value)
			}
		}
		for _, a := range spec.QualityAttributes {
			if system.AttributeLayout.Match(a.Location.Path) {
				if a.ID == spec.QualityAttributes[0].ID {
					modelText += "\n\nQuality attributes"
				}
				modelText += "\n  " + a.Name + " (" + a.ID + ") — " + displayPath(a.Location.Path)
			}
		}

		description := ""
		if spec.Description != "" {
			description = "\n\nDescription\n" + indent(spec.Description)
		}
		if spec.Framing != "" {
			description += "\n\nFraming\n" + indent(spec.Framing)
		}
		_, err = fmt.Fprintf(out, "%s%s\n\nPurpose\n%s\n\nScope\n%s%s\n\nEnclosing system document: %s\nSource: %s\n",
			result.System.Name, description, indent(result.System.Purpose), indent(result.System.Scope), modelText, parentText, displayPath(result.System.Path))
		return err
	case "list":
		nodes, err := system.List(path)
		if err != nil {
			return err
		}
		if nodes == nil {
			nodes = []*system.Node{}
		}
		result := systemsResult{1, nodes}
		if format == "json" {
			return writeJSON(out, result)
		}
		if len(result.Systems) == 0 {
			_, err = fmt.Fprintln(out, "No systems found in the selected subtree.")
			return err
		}
		var output strings.Builder
		output.WriteString("Systems by directory\n")
		writeTree(&output, result.Systems, "")
		_, err = io.WriteString(out, output.String())
		return err
	}
	return fmt.Errorf("unknown system command %q", command)
}

func newComponentCommand() *cobra.Command {
	root := group("component", "Describe a component in its own specification")
	var parent string
	cmd := &cobra.Command{Use: "init <name>", Short: "Create a named component specification", Args: cobra.ExactArgs(1),
		Example: "dotspec component init Intake --path internal --part-of specs/new-system.system.spec.md",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			target, err := system.InitComponent(path, args[0], parent)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeJSON(cmd.OutOrStdout(), createdResult{1, target})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", displayPath(target))
			return err
		}}
	cmd.Flags().StringVar(&parent, "part-of", "", "Parent system or component source file, relative to the working directory")
	_ = cmd.MarkFlagRequired("part-of")
	addPath(cmd)
	addOutput(cmd, false)
	root.AddCommand(cmd)
	return root
}
