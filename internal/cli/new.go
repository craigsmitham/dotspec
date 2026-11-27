package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	qualityapply "github.com/craigsmitham/dotspec/internal/qualitymodel/apply"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type createdResult struct {
	SchemaVersion int    `json:"schema_version"`
	Path          string `json:"path"`
}

func newNewCommand() *cobra.Command {
	root := kindGroup("new", "Create a specification draft of one kind",
		"Create a specification draft of one kind below the working directory and report its path with the authoring guide to read next. "+
			"New never overwrites: when the target already exists, it changes nothing and exits 2.")
	for _, k := range kinds() {
		cmd := &cobra.Command{Use: k.layout.Kind + " " + k.operand, Short: "Create a " + k.layout.Kind + " specification",
			Long: newLong(k), Args: operands(1, 1)}
		for _, example := range k.examples {
			cmd.Example = strings.TrimPrefix(cmd.Example+"\ndotspec new "+k.layout.Kind+" "+example, "\n")
		}
		if k.flags != nil {
			k.flags(cmd)
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			dir, err := workingScope.dir("")
			if err != nil {
				return err
			}
			if k.adopt != nil {
				return k.adopt(cmd, dir, args)
			}
			target, err := k.create(cmd, dir, args[0])
			if err != nil {
				return err
			}
			return writeCreated(cmd.OutOrStdout(), jsonOutput(cmd), target, k)
		}
		if k.adopt != nil {
			cmd.Args = operands(1, -1)
		}
		addJSON(cmd)
		root.AddCommand(cmd)
	}
	return root
}

func newLong(k kind) string {
	about := k.about
	if about == "" {
		about = fmt.Sprintf("Writes specs/%s/%s%s below the working directory.", k.layout.Collection, k.operand, k.layout.Suffix)
	}
	return k.short + ". " + about + "\n\nLearn more: dotspec docs " + k.layout.Kind
}

// kindGroup is new or list: a verb whose subcommands are the kinds. Invoked
// without a kind, it names the kinds it accepts.
func kindGroup(verb, short, long string) *cobra.Command {
	var list strings.Builder
	list.WriteString("Kinds:")
	for _, k := range kinds() {
		fmt.Fprintf(&list, "\n  %-19s %s\n  %-19s dotspec docs %s", k.layout.Kind, k.short, "", k.layout.Kind)
	}
	return &cobra.Command{Use: verb + " <kind>", Short: short, Long: long + "\n\n" + list.String(),
		Args: subcommands,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fmt.Errorf("%w\n\n%s", usageError(cmd, "missing <kind>"), list.String())
		}}
}

// writeCreated reports a created specification; guidance and the authoring
// guide follow the created path in text.
func writeCreated(out io.Writer, asJSON bool, target string, k kind) error {
	if asJSON {
		return writeJSON(out, createdResult{schemaVersion, target})
	}
	guidance := ""
	if k.guidance != "" {
		guidance = k.guidance + "\n"
	}
	_, err := fmt.Fprintf(out, "Created %s\n%sAuthoring guide: %s\n", displayPath(target), guidance, authoringPage(k.layout.Kind))
	return err
}

// adoptQualityAttributes adds attributes to the nearest enclosing system,
// blank or copied from a quality model.
func adoptQualityAttributes(cmd *cobra.Command, dir string, ids []string) error {
	from := flagValue(cmd, "from")
	if cmd.Flags().Changed("from") && from == "" {
		return usageError(cmd, "--from requires a model ID or path")
	}
	target, _, err := system.Show(dir)
	if err != nil {
		return err
	}
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	opts := qualityapply.Options{Attributes: ids, DryRun: dryRun, Revision: citationRevision()}
	var result qualityapply.Result
	if from == "" {
		result, err = qualityapply.Blank(target.Path, opts)
	} else {
		var source *qualitymodel.QualityModel
		if source, err = qualitymodel.Show(dir, from); err != nil {
			return err
		}
		result, err = qualityapply.Run(source, target.Path, opts)
	}
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if jsonOutput(cmd) {
		return writeJSON(out, struct {
			SchemaVersion int `json:"schema_version"`
			qualityapply.Result
		}{schemaVersion, result})
	}
	var text strings.Builder
	for _, entry := range result.Entries {
		label := "Created"
		if dryRun {
			label = "Would create"
		}
		if entry.Decision == "reuse" {
			label = "Reused parent"
		}
		fmt.Fprintf(&text, "%s %s — %s\n", label, entry.ID, displayPath(entry.Path))
		if dryRun && entry.Content != "" {
			fmt.Fprintln(&text, entry.Content)
		}
	}
	if !dryRun {
		noun := "attributes"
		if len(result.Added) == 1 {
			noun = "attribute"
		}
		fmt.Fprintf(&text, "Added %d %s to %s.\nAuthoring guide: %s\n", len(result.Added), noun, displayPath(target.Path), authoringPage("quality-attribute"))
	}
	_, err = io.WriteString(out, text.String())
	return err
}
