package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type systemResult struct {
	SchemaVersion int            `json:"schema_version"`
	System        *system.System `json:"system"`
	Parent        *system.System `json:"parent"`
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
	root := group("system", "Describe systems and discover their hierarchy")
	for _, action := range []struct{ name, short string }{
		{"init", "Create specs/system.spec.md without overwriting an existing file"},
		{"show", "Describe the nearest enclosing system and its parent"},
		{"list", "Display systems in the selected subtree"},
	} {
		cmd := &cobra.Command{Use: action.name, Short: action.short, Args: cobra.NoArgs,
			Example: "dotspec system " + action.name + " --path .",
			RunE: func(cmd *cobra.Command, _ []string) error {
				path, err := projectPath(cmd)
				if err != nil {
					return err
				}
				format, err := outputFormat(cmd)
				if err != nil {
					return err
				}
				return runSystem(cmd.Name(), path, format, cmd.OutOrStdout())
			},
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}

func runSystem(command, path, format string, out io.Writer) error {
	switch command {
	case "init":
		target, err := system.Init(path)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(out, createdResult{1, target})
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit the system name, Purpose, and Scope.\n", displayPath(target))
		return err
	case "show":
		spec, parent, err := system.Show(path)
		if err != nil {
			return err
		}
		result := systemResult{1, spec, parent}
		if format == "json" {
			return writeJSON(out, result)
		}
		parentText := "None"
		if result.Parent != nil {
			parentText = fmt.Sprintf("%s (%s)", result.Parent.Name, displayPath(result.Parent.Path))
		}
		modelText := ""
		if spec.QualityModel != nil {
			modelText = "\n\nQuality attributes\n" + indent(spec.QualityModel.Markdown)
		}
		_, err = fmt.Fprintf(out, "%s\n\nPurpose\n%s\n\nScope\n%s%s\n\nParent system: %s\nSource: %s\n",
			result.System.Name, indent(result.System.Purpose), indent(result.System.Scope), modelText, parentText, displayPath(result.System.Path))
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
		output.WriteString("Systems\n")
		writeTree(&output, result.Systems, "")
		_, err = io.WriteString(out, output.String())
		return err
	}
	return fmt.Errorf("unknown system command %q", command)
}
