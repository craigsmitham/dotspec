package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/qualitymodel"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

func newQualityModelCommand() *cobra.Command {
	root := group("quality-model", "Define reusable quality attributes and copy them into systems")
	for _, action := range []struct{ name, short string }{
		{"list", "List bundled starters and project quality models"},
		{"show", "Show a bundled starter or a local quality model"},
		{"init", "Create an editable project quality model"},
		{"check", "Check a local quality model's structure and suggest missing content"},
	} {
		cmd := newSpecAction(qualitymodel.Layout, action.name, action.short)
		if action.name == "show" {
			// The show example names a bundled starter rather than a project slug.
			cmd.Example = "dotspec quality-model show furps --path ."
		}
		if action.name == "init" {
			cmd.Flags().String("from", "", "Copy a bundled starter into the new project model")
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
			out := cmd.OutOrStdout()
			if action.name != "list" {
				if err := validateSelector(qualitymodel.Layout, action.name, args[0]); err != nil {
					return err
				}
			}
			if action.name == "init" {
				from, _ := cmd.Flags().GetString("from")
				target, err := qualitymodel.Init(path, args[0], from)
				if err != nil {
					return err
				}
				return writeCreated(out, format, target, "Authoring guide: "+authoringPage("quality-model"))
			}
			if action.name == "list" {
				project, err := specfile.FindProject(path)
				if err != nil {
					return err
				}
				effective, selection, err := qualitymodel.Resolve(path, project, "")
				if err != nil {
					return err
				}
				models, err := qualitymodel.List(path)
				if err != nil {
					return err
				}
				found := false
				for _, m := range models {
					if m.Path == effective.Path {
						found = true
						break
					}
				}
				if !found {
					models = append(models, effective)
				}
				if format == "json" {
					return writeJSON(out, struct {
						Effective     qualitymodel.Selection       `json:"effective_model"`
						SchemaVersion int                          `json:"schema_version"`
						Models        []*qualitymodel.QualityModel `json:"quality_models"`
					}{selection, 3, models})
				}
				if _, err := fmt.Fprintln(out, selectionText(selection)); err != nil {
					return err
				}
				for _, m := range models {
					mark := ""
					if m.Path == effective.Path {
						mark = " (default)"
					}
					if _, err := fmt.Fprintf(out, "%s%s — %s [%s]\n", m.Slug, mark, m.Name, m.Path); err != nil {
						return err
					}
				}
				return nil
			}
			if action.name == "check" {
				builtins, err := qualitymodel.Builtins()
				if err != nil {
					return err
				}
				for _, m := range builtins {
					if m.Slug == args[0] {
						return fmt.Errorf("%s is a bundled starter; check a local model, or copy it with quality-model init <slug> --from %s", m.Slug, m.Slug)
					}
				}
				return emitFocusedCheck(out, format, path, args[0], qualitymodel.CheckType())
			}
			m, err := qualitymodel.Show(path, args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeJSON(out, struct {
					SchemaVersion int                        `json:"schema_version"`
					Model         *qualitymodel.QualityModel `json:"quality_model"`
				}{1, m})
			}
			_, err = fmt.Fprintf(out, "%s\n\nSource: %s\n", m.Markdown, m.Path)
			return err
		}
		addPath(cmd)
		addOutput(cmd, false)
		root.AddCommand(cmd)
	}
	cmd := &cobra.Command{Use: "apply [model]", Short: "Adopt selected quality attributes into one system", Long: "Select IDs with --attributes or request the whole model with --all. Existing system declarations keep their wording and grouping. To explore the model first, use quality-attribute browse.", Example: "dotspec quality-model apply --attributes recoverability --dry-run", Args: cobra.MaximumNArgs(1)}
	cmd.Flags().StringSlice("attributes", nil, "Comma-separated attribute IDs to adopt")
	cmd.Flags().Bool("all", false, "Adopt every attribute from the model into one system")
	cmd.MarkFlagsMutuallyExclusive("all", "attributes")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ids, _ := cmd.Flags().GetStringSlice("attributes")
		if cmd.Flags().Changed("attributes") && (len(ids) == 0 || strings.TrimSpace(strings.Join(ids, "")) == "") {
			return fmt.Errorf("--attributes requires at least one attribute ID")
		}
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(ids) == 0 {
			return fmt.Errorf("select attributes with --attributes ID,ID or --all; use quality-attribute browse to explore the model")
		}
		if len(args) > 0 && args[0] == "" {
			return fmt.Errorf("expected a model ID or path")
		}
		return qualityOperation(cmd, qualityRequest{model: firstArg(args), attributes: ids, all: all})
	}
	addAdoptionFlags(cmd)
	root.AddCommand(cmd)
	return root
}
