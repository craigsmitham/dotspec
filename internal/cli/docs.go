package cli

import (
	"io"
	"strings"

	"github.com/dotspec-md/dotspec/internal/guide"
	"github.com/spf13/cobra"
)

func newDocsCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "docs [<topic> | <page>[#<section>] | <rule> | <words>...]",
		Short:   "Browse the offline authoring handbook",
		Long:    "Browse topics, open a page or section, or search this version's handbook and bundled examples. Topic names and exact addresses open directly. Enter words to find matching pages.",
		Example: "dotspec docs\ndotspec docs job\ndotspec docs success criteria\ndotspec docs guides/write-a-job\ndotspec docs job.performer",
		ValidArgsFunction: func(_ *cobra.Command, args []string, prefix string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return guide.Complete(prefix), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			page, err := guide.Read(strings.Join(args, " "))
			if err != nil {
				return err
			}
			_, err = io.WriteString(cmd.OutOrStdout(), page)
			return err
		},
	}
}

func authoringPage(topic string) string {
	return "dotspec docs " + guide.StartingPage(topic)
}
