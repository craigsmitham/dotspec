package cli

import (
	"io"

	"github.com/craigsmitham/dotspec/internal/guide"
	"github.com/spf13/cobra"
)

func newGuideCommand() *cobra.Command {
	return &cobra.Command{
		Use: "guide [TOPIC|RULE|PATH]", Short: "Read the offline authoring handbook and worked examples",
		Long:    "Open guidance from this version of dotspec. With no argument, list topics and help choose an authoring form. Use examples to browse worked specifications or a diagnostic's suggested rule ID to read its guidance.",
		Example: "dotspec guide\ndotspec guide obligation\ndotspec guide requirement.measurement-basis\ndotspec guide examples",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var topic string
			if len(args) == 1 {
				topic = args[0]
			}
			page, err := guide.Read(topic)
			if err != nil {
				return err
			}
			_, err = io.WriteString(cmd.OutOrStdout(), page)
			return err
		},
	}
}
