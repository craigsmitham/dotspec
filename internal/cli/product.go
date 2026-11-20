package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func newProductCommand(layout specfile.Layout) *cobra.Command {
	short, example := "Group stakeholder-visible behavior", "resume-work"
	sections := "## Description\n\nDescribe what a stakeholder can do.\n\n## Groups\n\nLink to the specifications that define this behavior."
	if layout.Kind == "capability" {
		short, example = "Describe an ability and its providers", "local-storage"
		sections = "## Description\n\nDescribe the ability and what it enables.\n\n## Provided by\n\nLink to the systems or components that provide this ability."
	}
	root := group(layout.Kind, short)
	for _, action := range []struct{ name, short string }{{"init", "Create a standalone specification"}, {"show", "Show the authored specification"}, {"list", "List specifications in the selected subtree"}, {"check", "Check the specification and its references"}} {
		cmd := &cobra.Command{Use: action.name + " <slug-or-source-path>", Short: action.short, Args: cobra.ExactArgs(1), Example: "dotspec " + layout.Kind + " " + action.name + " " + example}
		if action.name == "list" {
			cmd.Use = "list"
			cmd.Args = cobra.NoArgs
			cmd.Example = "dotspec " + layout.Kind + " list --path ."
		}
		if action.name == "init" {
			cmd.Use = "init <slug>"
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			dir, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			if action.name == "list" {
				paths, err := specfile.Discover(dir, layout.Match)
				if err != nil {
					return err
				}
				specs := []*system.ProductSpec{}
				for _, path := range paths {
					item, err := system.ReadProductSpec(layout, path)
					if err != nil {
						return err
					}
					specs = append(specs, item)
				}
				if format == "json" {
					return writeJSON(out, map[string]any{"schema_version": 1, layout.Collection: specs})
				}
				if len(specs) == 0 {
					_, err = fmt.Fprintf(out, "No %s found in the selected subtree.\n", layout.Collection)
					return err
				}
				for _, item := range specs {
					if _, err := fmt.Fprintf(out, "%s  (%s)\n", item.Name, displayPath(item.Path)); err != nil {
						return err
					}
				}
				return nil
			}
			selector := args[0]
			if action.name == "init" || !layout.Match(selector) {
				if err := system.ValidateProductSlug(layout.Kind, selector); err != nil {
					return err
				}
			}
			switch action.name {
			case "init":
				title := strings.ReplaceAll(selector, "-", " ")
				source := fmt.Sprintf("---\ntype: %s\n---\n\n# %s\n\n%s\n", specfile.AuthoredType(layout.Kind), title, sections)
				path, err := layout.Create(dir, selector, source)
				if err != nil {
					return err
				}
				if format == "json" {
					return writeJSON(out, createdResult{1, path})
				}
				_, err = fmt.Fprintf(out, "Created %s\n", displayPath(path))
				return err
			case "check":
				return emitCheckReport(out, check.Focused(dir, selector, system.ProductCheckType(layout), supportedTypes()...), format)
			default:
				path, err := layout.LocalPath(dir, selector)
				if err != nil {
					return err
				}
				item, err := system.ReadProductSpec(layout, path)
				if err != nil {
					return err
				}
				if format == "json" {
					return writeJSON(out, map[string]any{"schema_version": 1, layout.Kind: item})
				}
				_, err = fmt.Fprintf(out, "%s\n\nSource: %s\n", item.Markdown, displayPath(path))
				return err
			}
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	return root
}
