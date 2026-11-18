package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/guide"
	"github.com/spf13/cobra"
)

func newGuideCommand() *cobra.Command {
	var list bool
	cmd := &cobra.Command{
		Use: "guide [TOPIC|RULE|PATH]", Short: "Read the offline authoring handbook and worked examples",
		Long:    "Open guidance from this version of dotspec. With no argument, list topics and help choose an authoring form. Use --list to see topic descriptions and source paths. Use examples to browse worked specifications or a diagnostic's suggested rule ID to read its guidance.",
		Example: "dotspec guide\ndotspec guide --list\ndotspec guide obligation\ndotspec guide requirement.measurement-basis\ndotspec guide examples",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if list {
				if len(args) != 0 {
					return fmt.Errorf("--list takes no guide topic")
				}
				var out strings.Builder
				for _, route := range guide.Topics() {
					fmt.Fprintf(&out, "%s  %s\n", route.Name, route.Description)
					if route.Page != "" {
						fmt.Fprintf(&out, "  %s\n", route.Page)
					}
				}
				_, err := io.WriteString(cmd.OutOrStdout(), out.String())
				return err
			}
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
	cmd.Flags().BoolVar(&list, "list", false, "List guide topics with their descriptions and source paths")
	return cmd
}
