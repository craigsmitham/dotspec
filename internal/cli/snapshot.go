package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/dotspec-md/dotspec/internal/model"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	modelschema "github.com/dotspec-md/dotspec/model"
	"github.com/spf13/cobra"
)

func newSnapshotCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "snapshot", Short: "Export and validate specification snapshots", Args: subcommands,
		Long: "Export a project's source text, named domain fields, located relationships and inspection findings as one versioned JSON snapshot, print the snapshot's JSON Schema, or validate a saved snapshot.\n\nLearn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md",
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	export := &cobra.Command{Use: "export [<path>]", Short: "Export a source-backed specification snapshot as JSON", Args: operands(0, 1),
		Long:    "Export selected specifications, reusable quality models, enclosing systems, and direct reference context as JSON. Retain readable Markdown and source locations, including diagnostics for invalid documents. With no <path>, export covers the enclosing project; a directory selects that subtree. Filesystem discovery and authoring findings follow check; external references are not fetched. " + findingsExit + "\n\nLearn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md",
		Example: "dotspec snapshot export\ndotspec snapshot export internal",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := projectScope.dir(optionalPath(args))
			if err != nil {
				return err
			}
			result, err := exportValidated(path)
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), result.Snapshot); err != nil {
				return err
			}
			if result.ExitCode != 0 {
				return &exitError{code: result.ExitCode}
			}
			return nil
		}}
	schema := &cobra.Command{Use: "schema", Short: "Print the embedded snapshot JSON Schema", Args: operands(0, 0),
		Example: "dotspec snapshot schema",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write(modelschema.SchemaJSON)
			return err
		}}
	validate := &cobra.Command{Use: "validate <file>|-", Short: "Check a snapshot's shape, source ranges, and relationships", Args: operands(1, 1),
		Long: "Validate a saved snapshot file, or - to read the snapshot from standard input. Source paths and external targets are not fetched.\n\nA valid snapshot conforms to its declared model and mapping profile, and its source ranges and relationships are consistent; authoring errors it records faithfully do not make it invalid. Exits 1 when the snapshot is invalid.\n\nLearn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md",
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte
			var err error
			if args[0] == "-" {
				data, err = io.ReadAll(cmd.InOrStdin())
			} else {
				data, err = os.ReadFile(args[0])
			}
			if err != nil {
				return err
			}
			findings, err := model.ValidateJSON(data)
			if err != nil {
				return err
			}
			report := struct {
				SchemaVersion int             `json:"schema_version"`
				Valid         bool            `json:"valid"`
				Findings      []model.Finding `json:"findings"`
			}{schemaVersion, len(findings) == 0, findings}
			if jsonOutput(cmd) {
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
	addJSON(validate)
	cmd.AddCommand(export, schema, validate)
	return cmd
}

// exportValidated exports the project at path and rejects a snapshot that
// violates the model's invariants, so commands never emit a malformed model.
func exportValidated(path string) (model.ExportResult, error) {
	result, err := model.Export(path, speckinds.All()...)
	if err != nil {
		return result, err
	}
	findings, err := model.Validate(result.Snapshot)
	if err != nil {
		return result, fmt.Errorf("validate exported model: %w", err)
	}
	if len(findings) > 0 {
		return result, fmt.Errorf("exported model violates %d invariants; first: %s at %s: %s", len(findings), findings[0].Code, findings[0].Path, findings[0].Message)
	}
	return result, nil
}
