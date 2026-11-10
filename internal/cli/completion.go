package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Resolve the writer at execution time. Cobra's default completion command
// captures it during construction, before Run attaches the caller's streams.
func newCompletionCommand() *cobra.Command {
	root := group("completion", "Generate a shell completion script")
	root.Long = "Generate a completion script, then load it using your shell's completion setup."
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		cmd := &cobra.Command{Use: shell, Short: "Generate the " + shell + " completion script",
			Example: "dotspec completion " + shell, Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				noDescriptions, err := cmd.Flags().GetBool("no-descriptions")
				if err != nil {
					return err
				}
				out, tree := cmd.OutOrStdout(), cmd.Root()
				switch cmd.Name() {
				case "bash":
					return tree.GenBashCompletionV2(out, !noDescriptions)
				case "zsh":
					if noDescriptions {
						return tree.GenZshCompletionNoDesc(out)
					}
					return tree.GenZshCompletion(out)
				case "fish":
					return tree.GenFishCompletion(out, !noDescriptions)
				case "powershell":
					if noDescriptions {
						return tree.GenPowerShellCompletion(out)
					}
					return tree.GenPowerShellCompletionWithDesc(out)
				}
				return fmt.Errorf("unsupported shell %q", cmd.Name())
			},
		}
		cmd.Flags().Bool("no-descriptions", false, "omit completion descriptions")
		root.AddCommand(cmd)
	}
	return root
}
