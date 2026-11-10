package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func addOutput(cmd *cobra.Command, mermaid bool) {
	formats := "text or json"
	if mermaid {
		formats += " or mermaid"
	}
	cmd.Flags().String("format", "text", "output format: "+formats)
	cmd.Flags().Bool("json", false, "emit structured JSON (equivalent to --format json)")
	cmd.MarkFlagsMutuallyExclusive("json", "format")
	_ = cmd.RegisterFlagCompletionFunc("format", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		values := []string{"text", "json"}
		if mermaid {
			values = append(values, "mermaid")
		}
		return values, cobra.ShellCompDirectiveNoFileComp
	})
}

func outputFormat(cmd *cobra.Command) (string, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return "", err
	}
	jsonMode, err := cmd.Flags().GetBool("json")
	if err != nil {
		return "", err
	}
	if jsonMode {
		return "json", nil
	}
	if format == "text" || format == "json" || (cmd.Name() == "map" && format == "mermaid") {
		return format, nil
	}
	return "", fmt.Errorf("unknown output format %q; see %s --help", format, cmd.CommandPath())
}

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
