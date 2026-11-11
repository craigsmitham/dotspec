package cli

import (
	"fmt"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	qualityapply "github.com/craigsmitham/dotspec/internal/qualitymodel/apply"
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
				_, err = fmt.Fprintf(out, "Created %s\n", displayPath(target))
				return err
			}
			if action.name == "list" {
				models, err := qualitymodel.List(path)
				if err != nil {
					return err
				}
				if format == "json" {
					return writeJSON(out, struct {
						SchemaVersion int                          `json:"schema_version"`
						Default       string                       `json:"default"`
						Models        []*qualitymodel.QualityModel `json:"quality_models"`
					}{1, qualitymodel.Default, models})
				}
				for _, m := range models {
					mark := ""
					if m.Builtin && m.Slug == qualitymodel.Default {
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
	cmd := &cobra.Command{Use: "apply [model]", Short: "Copy selected model attributes into a system without overwriting declarations", Long: "Copy attributes from a bundled starter or local quality model. The default model is iso-25010-2023. With no --attributes selection, copy all attributes. Include ancestors automatically; refuse ID conflicts before writing. --dry-run prints the exact resulting document.\n\nLocal declarations remain independent of their source model.", Example: "dotspec quality-model apply --system specs/system.spec.md --dry-run", Args: cobra.MaximumNArgs(1)}
	cmd.Flags().String("system", "", "Target system.spec.md path, relative to --path")
	cmd.Flags().StringSlice("attributes", nil, "Comma-separated attribute IDs to copy, including ancestors")
	cmd.Flags().Bool("dry-run", false, "Preview the complete resulting document without writing")
	_ = cmd.MarkFlagRequired("system")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		path, err := projectPath(cmd)
		if err != nil {
			return err
		}
		format, err := outputFormat(cmd)
		if err != nil {
			return err
		}
		slug := qualitymodel.Default
		if len(args) > 0 {
			slug = args[0]
		}
		m, err := qualitymodel.Show(path, slug)
		if err != nil {
			return err
		}
		target, _ := cmd.Flags().GetString("system")
		if !filepath.IsAbs(target) {
			target = filepath.Join(path, target)
		}
		selected, _ := cmd.Flags().GetStringSlice("attributes")
		if cmd.Flags().Changed("attributes") && len(selected) == 0 {
			return fmt.Errorf("--attributes requires at least one attribute ID")
		}
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		result, err := qualityapply.Run(m, target, selected, dryRun)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeJSON(cmd.OutOrStdout(), result)
		}
		if dryRun {
			_, err = fmt.Fprint(cmd.OutOrStdout(), result.Content)
		} else {
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Added %d attributes to %s from %s.\n", len(result.Added), displayPath(result.Target), m.Slug)
		}
		return err
	}
	addPath(cmd)
	addOutput(cmd, false)
	root.AddCommand(cmd)
	return root
}
