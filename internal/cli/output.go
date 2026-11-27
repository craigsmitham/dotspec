package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// outputFormatsAnnotation records, on each command, the --format values addOutput allows.
const outputFormatsAnnotation = "dotspec.output-formats"

func addOutput(cmd *cobra.Command, mermaid bool) {
	formats := []string{"text", "json"}
	if mermaid {
		formats = append(formats, "mermaid")
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[outputFormatsAnnotation] = strings.Join(formats, ",")
	cmd.Flags().String("format", "text", "output format: "+strings.Join(formats, " or "))
	cmd.Flags().Bool("json", false, "emit structured JSON (equivalent to --format json)")
	cmd.MarkFlagsMutuallyExclusive("json", "format")
	_ = cmd.RegisterFlagCompletionFunc("format", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return formats, cobra.ShellCompDirectiveNoFileComp
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
	if slices.Contains(strings.Split(cmd.Annotations[outputFormatsAnnotation], ","), format) {
		return format, nil
	}
	return "", fmt.Errorf("unknown output format %q; see %s --help", format, cmd.CommandPath())
}

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
