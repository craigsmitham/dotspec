package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

// listResult is the one list envelope every kind shares.
type listResult struct {
	SchemaVersion int         `json:"schema_version"`
	Entries       []listEntry `json:"entries"`
}

func newListCommand() *cobra.Command {
	root := kindGroup("list", "List the specifications of one kind",
		"List the specification files of one kind in the enclosing project, grouped by the system each belongs to. list system prints systems by directory, and list component prints components by their Part of links.")
	for _, k := range kinds() {
		short := k.listShort
		if short == "" {
			short = "List " + k.plural + " by system"
		}
		cmd := &cobra.Command{Use: k.layout.Kind, Short: short, Long: short + ". " + k.short + ".\n\nLearn more: dotspec docs " + k.layout.Kind,
			Example: "dotspec list " + k.layout.Kind, Args: operands(0, 0),
			RunE: func(cmd *cobra.Command, _ []string) error {
				dir, err := projectScope.dir("")
				if err != nil {
					return err
				}
				entries, err := k.members(dir)
				if err != nil {
					return err
				}
				return writeList(cmd.OutOrStdout(), jsonOutput(cmd), k, entries)
			}}
		addJSON(cmd)
		root.AddCommand(cmd)
	}
	return root
}

func writeList(out io.Writer, asJSON bool, k kind, entries []listEntry) error {
	if asJSON {
		return writeJSON(out, listResult{schemaVersion, entries})
	}
	if len(entries) == 0 {
		_, err := fmt.Fprintf(out, "No %s found in the enclosing project.\n", k.plural)
		return err
	}
	if k.tree != nil {
		_, err := io.WriteString(out, k.tree(entries))
		return err
	}
	var text strings.Builder
	text.WriteString(strings.ToUpper(k.plural[:1]) + k.plural[1:] + "\n")
	for i, e := range entries {
		if i == 0 || !sameSystem(entries[i-1].System, e.System) {
			if e.System == nil {
				text.WriteString("\nNo associated system\n")
			} else {
				fmt.Fprintf(&text, "\nSystem: %s (%s)\n", e.System.Name, displayPath(e.System.Path))
			}
		}
		fmt.Fprintf(&text, "  %s  (%s)\n", e.Name, displayPath(e.Path))
	}
	_, err := io.WriteString(out, text.String())
	return err
}

func sameSystem(a, b *systemRef) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Path == b.Path
}
