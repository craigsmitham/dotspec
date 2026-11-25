package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/craigsmitham/dotspec/internal/model"
	modelschema "github.com/craigsmitham/dotspec/model"
	"github.com/spf13/cobra"
)

func newModelCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "model", Short: "Export and validate complete specification snapshots", Args: cobra.NoArgs,
		Long: "Export source text, named domain fields, located relationships and inspection findings in a versioned specification model.",
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	export := &cobra.Command{Use: "export", Short: "Export a source-backed specification snapshot as JSON", Args: cobra.NoArgs,
		Long:    "Export selected specifications, reusable quality models, enclosing systems, and direct reference context as JSON. Retain readable Markdown and source locations, including diagnostics for invalid documents. Filesystem discovery and authoring findings follow check; external references are not fetched.\n\n" + exitPolicy,
		Example: "dotspec model export\ndotspec model export --path internal",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return fmt.Errorf("validate exported model: %w", err)
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s at %s: %s", len(findings), findings[0].Code, findings[0].Path, findings[0].Message)
			}
			if err := writeJSON(cmd.OutOrStdout(), result.Snapshot); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	addPath(export)
	schema := &cobra.Command{Use: "schema", Short: "Print the embedded snapshot JSON Schema", Args: cobra.NoArgs,
		Example: "dotspec model schema",
		RunE: func(cmd *cobra.Command, _ []string) error {
			workspace, _ := cmd.Flags().GetBool("workspace")
			if workspace {
				schema, err := workspaceExportSchema()
				if err != nil {
					return err
				}
				return writeJSON(cmd.OutOrStdout(), schema)
			}
			_, err := cmd.OutOrStdout().Write(modelschema.SchemaJSON)
			return err
		}}
	schema.Flags().Bool("workspace", false, "print the schema for a complete model export --all envelope")
	validate := &cobra.Command{Use: "validate SNAPSHOT", Short: "Check a snapshot's shape, source ranges, and relationships", Args: cobra.ExactArgs(1),
		Long: "Validate a snapshot file, or use - to read standard input. Snapshot validation checks representation conformance and integrity, not requirements fitness or system behavior; faithfully represented authoring errors do not make a snapshot invalid. No source paths or external targets are fetched.\n\n" + exitPolicy,
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			var data []byte
			if args[0] == "-" {
				data, err = io.ReadAll(cmd.InOrStdin())
			} else {
				data, err = os.ReadFile(args[0])
			}
			if err != nil {
				return err
			}
			findings, err := validateExportJSON(data)
			if err != nil {
				return err
			}
			report := struct {
				Valid    bool            `json:"valid"`
				Findings []model.Finding `json:"findings"`
			}{len(findings) == 0, findings}
			if format == "json" {
				err = writeJSON(cmd.OutOrStdout(), report)
			} else if report.Valid {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "Snapshot is valid for its declared model and mapping profile.")
			} else {
				for _, finding := range findings {
					if _, err = fmt.Fprintf(cmd.OutOrStdout(), "%s [%s] %s\n", finding.Path, finding.Code, finding.Message); err != nil {
						break
					}
				}
			}
			if err != nil {
				return err
			}
			if !report.Valid {
				return &exitError{code: 1}
			}
			return nil
		}}
	addOutput(validate, false)
	docs := &cobra.Command{Use: "docs", Short: "Generate concept references, system compositions, and dependency diagrams as Markdown", Args: cobra.NoArgs,
		Long:    "Render purpose and operational principles before state and actions. Show explicit named uses, synchronization rules and application dependencies. Diagrams contain resolved relationships only; no behavior is executed or inferred.",
		Example: "dotspec model docs --path . --relative-to views",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			base, err := cmd.Flags().GetString("relative-to")
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s", len(findings), findings[0].Message)
			}
			text, err := model.DesignMarkdown(result.Snapshot, base)
			if err != nil {
				return err
			}
			if _, err = io.WriteString(cmd.OutOrStdout(), text); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	trace := &cobra.Command{Use: "trace", Short: "Generate layered requirements traceability and change-review views as Markdown", Args: cobra.NoArgs,
		Long:    "Show motivation, derivation, refinement, assumptions, use-case participation, reverse change-review links and requirement criteria. Resolved explicit relationships only; paths do not establish satisfaction or evidence.",
		Example: "dotspec model trace --path . --relative-to views",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			base, err := cmd.Flags().GetString("relative-to")
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s", len(findings), findings[0].Message)
			}
			text, err := model.TraceMarkdown(result.Snapshot, base)
			if err != nil {
				return err
			}
			if _, err = io.WriteString(cmd.OutOrStdout(), text); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	overview := &cobra.Command{Use: "overview", Short: "Generate system structure, features, and applicable requirements as Markdown", Args: cobra.NoArgs,
		Long:    "Show explicit system parts, non-owning feature groups, applicability, concept uses, native references and requirement criteria. Directory context remains separate; no implementation or requirement inheritance is inferred.",
		Example: "dotspec model overview --path . --relative-to views",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			base, err := cmd.Flags().GetString("relative-to")
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s", len(findings), findings[0].Message)
			}
			text, err := model.OverviewMarkdown(result.Snapshot, base)
			if err != nil {
				return err
			}
			if _, err = io.WriteString(cmd.OutOrStdout(), text); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	needs := &cobra.Command{Use: "needs", Short: "Generate a needs-to-specification map with capability dependencies", Args: cobra.NoArgs,
		Long:    "Show jobs, shared outcomes, features, requirements and optional capabilities. Short labels link to full statements and evidence. Only explicit relationships are projected; contribution does not prove achievement.",
		Example: "dotspec model needs --path . --relative-to views",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			base, err := cmd.Flags().GetString("relative-to")
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s", len(findings), findings[0].Message)
			}
			text, err := model.NeedsMarkdown(result.Snapshot, base)
			if err != nil {
				return err
			}
			if _, err = io.WriteString(cmd.OutOrStdout(), text); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	addPath(docs)
	docs.Flags().String("relative-to", ".", "Directory containing the generated Markdown; source links are relative to it")
	addPath(trace)
	trace.Flags().String("relative-to", ".", "Directory containing the generated Markdown; source links are relative to it")
	addPath(overview)
	overview.Flags().String("relative-to", ".", "Directory containing the generated Markdown; source links are relative to it")
	addPath(needs)
	needs.Flags().String("relative-to", ".", "Directory containing the generated Markdown; source links are relative to it")
	product := &cobra.Command{Use: "product [SYSTEM]", Short: "Show the specification roots, owned elements and references for one system", Args: cobra.MaximumNArgs(1),
		Example: "dotspec model product",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			selector := ""
			if len(args) > 0 {
				selector = args[0]
			}
			view, err := result.Snapshot.Product(selector)
			if err != nil {
				return err
			}
			if root := model.String(model.Object(view.Scope["project"]), "root"); root != "" && view.System != "" {
				evidence, err := model.ReadVerificationEvidence(root)
				if err != nil {
					return err
				}
				if err := view.IncludeEvidence(result.Snapshot, root, evidence); err != nil {
					return err
				}
			}
			if err = writeJSON(cmd.OutOrStdout(), view); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		},
	}
	addPath(product)
	cmd.AddCommand(product)
	cmd.AddCommand(export, schema, validate, docs, trace, overview, needs,
		newModelReadingCommand("review", "Generate advisory requirements review and open-question inventory", model.ReviewMarkdown),
		newModelReadingCommand("pegs", "Generate a Project, Environment, Goals, System reading perspective", model.PEGSMarkdown))
	return cmd
}

func newModelReadingCommand(name, description string, render func(*model.Snapshot, string) (string, error)) *cobra.Command {
	cmd := &cobra.Command{Use: name, Short: description, Args: cobra.NoArgs,
		Long:    description + ". Requirements verification concerns specification quality; requirements and system validation concern needs and intended use. Only acquired content and explicit relationships are shown; missing information produces advisory review prompts, not completeness or endorsement claims.\n\n" + exitPolicy,
		Example: "dotspec model " + name + " --path . --relative-to views",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			base, err := cmd.Flags().GetString("relative-to")
			if err != nil {
				return err
			}
			result, err := model.Export(path, supportedTypes()...)
			if err != nil {
				return err
			}
			findings, err := model.Validate(result.Snapshot)
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("exported model violates %d invariants; first: %s", len(findings), findings[0].Message)
			}
			output, err := render(result.Snapshot, base)
			if err != nil {
				return err
			}
			if _, err = io.WriteString(cmd.OutOrStdout(), output); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	addPath(cmd)
	cmd.Flags().String("relative-to", ".", "Directory containing the generated Markdown; source links are relative to it")
	return cmd
}
