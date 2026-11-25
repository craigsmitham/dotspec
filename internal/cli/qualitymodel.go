package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
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
		{"check", "Check quality-model structure and suggest missing content"},
	} {
		cmd := &cobra.Command{Use: action.name + " <slug>", Short: action.short, Args: cobra.ExactArgs(1), Example: "dotspec quality-model " + action.name + " furps --path ."}
		if action.name == "init" {
			cmd.Example = "dotspec quality-model init organization-quality --path ."
			cmd.Flags().String("from", "", "Copy a bundled starter into the new project model")
		}
		if action.name == "list" {
			cmd.Use = "list"
			cmd.Args = cobra.NoArgs
			cmd.Example = "dotspec quality-model list --path ."
		}
		if action.name == "check" {
			cmd.Long = action.short + ".\n\n" + exitPolicy
		}
		if action.name == "show" || action.name == "check" {
			cmd.Use = action.name + " <slug-or-source-path>"
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
			if action.name == "init" {
				from, _ := cmd.Flags().GetString("from")
				target, err := qualitymodel.Init(path, args[0], from)
				if err != nil {
					return err
				}
				if format == "json" {
					return writeJSON(out, createdResult{1, target})
				}
				_, err = fmt.Fprintf(out, "Created %s\nAuthoring guide: %s\n", displayPath(target), authoringPage("quality-model"))
				return err
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
						Default       string                       `json:"default"`
						Models        []*qualitymodel.QualityModel `json:"quality_models"`
					}{selection, 2, effective.Slug, models})
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
			// Local checks retain standard report semantics, including unreadable files.
			if action.name == "check" {
				builtins, err := qualitymodel.Builtins()
				if err != nil {
					return err
				}
				for _, m := range builtins {
					if m.Slug == args[0] {
						findings := qualitymodel.Check(m)
						report := check.Report{SchemaVersion: 1, Scope: check.Scope{Path: m.Path, Types: []string{"quality-model"}, Discovery: "bundled"}, Complete: true, ContextFiles: []string{}, Summary: check.Summary{FilesDiscovered: 1, FilesChecked: 1, Advice: len(findings)}, Diagnostics: findings}
						return emitCheckReport(out, report, format)
					}
				}
				return emitCheckReport(out, check.Focused(path, args[0], qualitymodel.CheckType(), supportedTypes()...), format)
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
	cmd := &cobra.Command{Use: "apply [model]", Short: "Adopt selected quality attributes into one system", Long: "Browse the effective model, select IDs with --attributes or request the whole model with --all. Existing system declarations keep their wording and grouping.", Example: "dotspec quality-model apply --attributes recoverability --dry-run", Args: cobra.MaximumNArgs(1)}
	cmd.Flags().StringSlice("attributes", nil, "Comma-separated attribute IDs to adopt")
	cmd.Flags().Bool("all", false, "Adopt every attribute from the model into one system")
	cmd.MarkFlagsMutuallyExclusive("all", "attributes")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ids, _ := cmd.Flags().GetStringSlice("attributes")
		if cmd.Flags().Changed("attributes") && (len(ids) == 0 || strings.TrimSpace(strings.Join(ids, "")) == "") {
			return fmt.Errorf("--attributes requires at least one attribute ID")
		}
		all, _ := cmd.Flags().GetBool("all")
		if len(args) > 0 && args[0] == "" {
			return fmt.Errorf("expected a model ID or path")
		}
		return qualityOperation(cmd, qualityRequest{model: firstArg(args), attributes: ids, all: all})
	}
	addAdoptionFlags(cmd)
	root.AddCommand(cmd)
	return root
}
