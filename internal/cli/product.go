package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func newProductCommand(layout specfile.Layout) *cobra.Command {
	short := "Group stakeholder-visible behavior"
	sections := "## Description\n\nDescribe what a stakeholder can do.\n\n## Groups\n\nLink to the specifications that define this behavior."
	if layout.Kind == "capability" {
		short = "Describe an ability and its providers"
		sections = "## Description\n\nDescribe the ability and what it enables.\n\n## Provided by\n\nLink to the systems or components that provide this ability."
	}
	root := group(layout.Kind, short)
	for _, action := range []struct{ name, short string }{{"init", "Create a standalone specification"}, {"show", "Show the authored specification"}, {"list", "List specifications in the selected subtree"}, {"check", "Check the specification and its references"}} {
		cmd := newSpecAction(layout, action.name, action.short)
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
			if err := validateSelector(layout, action.name, selector); err != nil {
				return err
			}
			switch action.name {
			case "init":
				title := strings.ReplaceAll(selector, "-", " ")
				source := fmt.Sprintf("---\ntype: %s\n---\n\n# %s\n\n%s\n", specfile.AuthoredType(layout.Kind), title, sections)
				path, err := layout.Create(dir, selector, source)
				if err != nil {
					return err
				}
				return writeCreated(out, format, path, "Authoring guide: "+authoringPage(layout.Kind))
			case "check":
				return emitFocusedCheck(out, format, dir, selector, system.ProductCheckType(layout))
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
