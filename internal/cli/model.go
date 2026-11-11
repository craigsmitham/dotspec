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
		Long: "Represent original source bytes, document structure, and named domain content in a versioned specification model. Coverage describes the declared mapping profile, not requirements adequacy or satisfaction.",
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	export := &cobra.Command{Use: "export", Short: "Export a source-backed specification snapshot as JSON", Args: cobra.NoArgs,
		Long:    "Export selected specifications, enclosing systems, and direct reference context as JSON. Preserve invalid documents and original source bytes. Filesystem discovery and authoring findings follow check; external references are not fetched.\n\n" + exitPolicy,
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
			_, err := cmd.OutOrStdout().Write(modelschema.SchemaJSON)
			return err
		}}
	validate := &cobra.Command{Use: "validate SNAPSHOT", Short: "Check a snapshot's shape, references, source fidelity, and coverage", Args: cobra.ExactArgs(1),
		Long: "Validate a snapshot file, or use - to read standard input. Validation checks the model itself; faithfully represented authoring errors do not make a snapshot invalid. No source paths or external targets are fetched. External payload locators can be checked structurally; their byte availability is not verified offline.\n\n" + exitPolicy,
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
			findings, err := model.ValidateJSON(data)
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
	cmd.AddCommand(export, schema, validate)
	return cmd
}
