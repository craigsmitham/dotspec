package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	qualityapply "github.com/craigsmitham/dotspec/internal/qualitymodel/apply"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

func newQualityAttributeCommand() *cobra.Command {
	root := group("quality-attribute", "Browse a quality model and adopt system quality attributes")
	browse := &cobra.Command{Use: "browse [attribute-id]", Short: "Browse top-level terms or inspect a term and its neighbors", Example: "dotspec quality-attribute browse reliability", Args: cobra.MaximumNArgs(1)}
	browse.RunE = func(cmd *cobra.Command, args []string) error {
		return qualityOperation(cmd, qualityRequest{focus: firstArg(args)})
	}
	add := &cobra.Command{Use: "add <attribute-id>", Short: "Adopt a quality attribute and its missing parents into one system", Example: "dotspec quality-attribute add recoverability --dry-run", Args: cobra.ExactArgs(1)}
	add.RunE = func(cmd *cobra.Command, args []string) error {
		return qualityOperation(cmd, qualityRequest{attributes: args})
	}
	for _, cmd := range []*cobra.Command{browse, add} {
		cmd.Flags().String("model", "", "Use a bundled model, local slug, or source path for this invocation")
	}
	addPath(browse)
	addOutput(browse, false)
	addAdoptionFlags(add)
	root.AddCommand(browse, add)
	return root
}

func addAdoptionFlags(cmd *cobra.Command) {
	cmd.Flags().String("system", "", "Target system source path, relative to --path")
	cmd.Flags().Bool("dry-run", false, "Preview the decisions and new file contents")
	addPath(cmd)
	addOutput(cmd, false)
}

func firstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

type qualityRequest struct {
	model      string
	focus      string
	attributes []string
	all        bool
}

func qualityOperation(cmd *cobra.Command, req qualityRequest) error {
	dir, err := projectPath(cmd)
	if err != nil {
		return err
	}
	format, err := outputFormat(cmd)
	if err != nil {
		return err
	}
	project, err := specfile.FindProject(dir)
	if err != nil {
		return err
	}
	var target string
	write := req.all || len(req.attributes) > 0
	if write {
		selector, _ := cmd.Flags().GetString("system")
		if selector == "" && cmd.Flags().Changed("system") {
			return fmt.Errorf("--system requires a source path")
		}
		target, err = qualityTarget(dir, project, selector)
		if err != nil {
			return err
		}
		project, err = specfile.FindProject(filepath.Dir(target))
		if err != nil {
			return err
		}
	}
	explicit := req.model
	if cmd.Flags().Changed("model") {
		explicit, _ = cmd.Flags().GetString("model")
		if explicit == "" {
			return fmt.Errorf("--model requires a model ID or path")
		}
	}
	source, selection, err := qualitymodel.Resolve(dir, project, explicit)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if !write {
		return browseQuality(out, format, source, selection, req.focus)
	}
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	result, err := qualityapply.Run(source, target, qualityapply.Options{Attributes: req.attributes, All: req.all, DryRun: dryRun, Revision: citationRevision()})
	if err != nil {
		return err
	}
	if format == "json" {
		return writeJSON(out, struct {
			qualityapply.Result
			Model qualitymodel.Selection `json:"effective_model"`
		}{result, selection})
	}
	var text strings.Builder
	fmt.Fprintln(&text, selectionText(selection))
	for _, entry := range result.Entries {
		label := "Created"
		if dryRun {
			label = "Would create"
		}
		if entry.Decision == "reuse" {
			label = "Already declared"
			if entry.Role == "ancestor" {
				label = "Reused parent"
			}
		}
		fmt.Fprintf(&text, "%s %s — %s\n", label, entry.ID, displayPath(entry.Path))
		if dryRun && entry.Content != "" {
			fmt.Fprintln(&text, entry.Content)
		}
	}
	if !dryRun {
		fmt.Fprintf(&text, "Added %d attributes to %s.\n", len(result.Added), displayPath(target))
	}
	_, err = io.WriteString(out, text.String())
	return err
}

func qualityTarget(dir string, project *specfile.Project, selector string) (string, error) {
	if selector != "" {
		return system.Layout.SourcePath(dir, selector)
	}
	boundary, err := specfile.ContextRoot(dir)
	if err != nil {
		return "", err
	}
	for current := dir; ; current = filepath.Dir(current) {
		file, err := system.Layout.LocalSubject(current)
		if file != "" || err != nil {
			return file, err
		}
		if current == boundary || filepath.Dir(current) == current {
			break
		}
	}
	scope := dir
	if project != nil {
		scope = project.Root
	}
	paths, err := specfile.Discover(scope, system.Layout.Match)
	if err != nil {
		return "", err
	}
	switch len(paths) {
	case 0:
		return "", fmt.Errorf("no system found in %s; create one with dotspec system init or select --system PATH", scope)
	case 1:
		return paths[0], nil
	default:
		return "", fmt.Errorf("select --system PATH from: %s", strings.Join(paths, ", "))
	}
}

func selectionText(s qualitymodel.Selection) string {
	return fmt.Sprintf("Quality model: %s (%s) — %s", s.ID, s.Basis, s.Source)
}

type browseTerm struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

func browseQuality(out io.Writer, format string, model *qualitymodel.QualityModel, selection qualitymodel.Selection, focus string) error {
	view := struct {
		SchemaVersion int                    `json:"schema_version"`
		Model         qualitymodel.Selection `json:"effective_model"`
		Attribute     *browseTerm            `json:"attribute,omitempty"`
		Parents       []browseTerm           `json:"parents"`
		Attributes    []browseTerm           `json:"attributes"`
	}{SchemaVersion: 1, Model: selection, Parents: []browseTerm{}, Attributes: []browseTerm{}}
	var parents []string
	if focus != "" {
		for _, a := range model.QualityAttributes {
			if a.ID == focus {
				view.Attribute = &browseTerm{a.ID, a.Name, a.Definition}
				parents = a.SubattributeOf
				break
			}
		}
		if view.Attribute == nil {
			return qualitymodel.UnknownAttribute(model, focus)
		}
	}
	for _, a := range model.QualityAttributes {
		summary, _, _ := strings.Cut(a.Definition, "\n\n")
		term := browseTerm{a.ID, a.Name, summary}
		if slices.Contains(parents, a.ID) {
			view.Parents = append(view.Parents, term)
		}
		if (focus == "" && len(a.SubattributeOf) == 0) || slices.Contains(a.SubattributeOf, focus) {
			view.Attributes = append(view.Attributes, term)
		}
	}
	if format == "json" {
		return writeJSON(out, view)
	}
	var buffer strings.Builder
	fmt.Fprintln(&buffer, selectionText(selection))
	if a := view.Attribute; a != nil {
		fmt.Fprintf(&buffer, "\n%s (%s)\n%s\n", a.Name, a.ID, a.Definition)
		for _, p := range view.Parents {
			fmt.Fprintf(&buffer, "Parent: %s (%s)\n", p.Name, p.ID)
		}
		if len(view.Attributes) > 0 {
			buffer.WriteString("\nChildren\n")
		}
	}
	for _, a := range view.Attributes {
		fmt.Fprintf(&buffer, "%s (%s) — %s\n", a.Name, a.ID, a.Definition)
	}
	_, err := io.WriteString(out, buffer.String())
	return err
}
