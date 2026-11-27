package cli

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
)

// schemaVersion versions every command result except the snapshot, which
// carries its own model and representation versions.
const schemaVersion = 1

// addJSON gives a command that returns a result a local --json switch.
func addJSON(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "emit the result as JSON")
}

func jsonOutput(cmd *cobra.Command) bool {
	enabled, _ := cmd.Flags().GetBool("json")
	return enabled
}

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
