package cli

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Use Cobra's writer-aware generator; Fang's bundled man command writes to the
// process stdout directly, bypassing embedding and output-error handling.
func newManCommand() *cobra.Command {
	return &cobra.Command{Use: "man", Short: "Write the dotspec manual in roff format", Args: cobra.NoArgs,
		Example: "dotspec man",
		RunE: func(cmd *cobra.Command, _ []string) error {
			epoch := time.Unix(0, 0).UTC()
			cmd.Root().DisableAutoGenTag = true
			return doc.GenMan(cmd.Root(), &doc.GenManHeader{Date: &epoch, Title: "DOTSPEC", Section: "1", Source: "dotspec", Manual: "dotspec manual"}, cmd.OutOrStdout())
		},
	}
}
