package cli

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

type initResult struct {
	SchemaVersion    int               `json:"schema_version"`
	Project          string            `json:"project"`
	Manifest         string            `json:"manifest"`
	System           string            `json:"system"`
	EnclosingProject *specfile.Project `json:"enclosing_project"`
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "init <system-name>", Short: "Start a specification project",
		Long: "Start a specification project in the working directory. Init writes dotspec.yaml, which declares the project, and the root system specification, whose heading is <system-name> and whose filename is derived from it. The project name defaults to the system name's slug. " +
			"Init refuses to overwrite either file and exits 2 instead. When the directory lies inside another project, init says so: that project's discovery now stops at the new manifest.\n\nLearn more: dotspec docs system",
		Example: "dotspec init Bookings", Args: operands(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := flagValue(cmd, "name")
			if project == "" {
				project = specfile.NameSlug(args[0])
				if project != "" && !specfile.ValidProjectName(project) {
					return fmt.Errorf("invalid project name %q derived from the system name; choose one with --name", project)
				}
			}
			// Read the enclosing manifest first, so a failure to read it writes nothing.
			enclosing, err := specfile.EnclosingProject(".")
			if err != nil {
				return err
			}
			manifest, root, err := system.InitProject(".", args[0], project)
			if err != nil {
				return err
			}
			if jsonOutput(cmd) {
				return writeJSON(cmd.OutOrStdout(), initResult{schemaVersion, project, manifest, root, enclosing})
			}
			var out strings.Builder
			fmt.Fprintf(&out, "Created %s for project %s\nCreated %s\n", displayPath(manifest), project, displayPath(root))
			if enclosing != nil {
				fmt.Fprintf(&out, "This project is nested in project %s (%s). Checks and builds of %s no longer include this directory.\n", enclosing.Name, displayPath(enclosing.Manifest()), enclosing.Name)
			}
			fmt.Fprintf(&out, "Edit Description, Purpose, and Scope. Keep the filename in sync if you rename the system.\nAuthoring guide: %s\n", authoringPage("system"))
			_, err = fmt.Fprint(cmd.OutOrStdout(), out.String())
			return err
		},
	}
	stringFlag(cmd.Flags(), "name", "", "<project>", "", "project name written to dotspec.yaml; defaults to the system name's slug")
	addJSON(cmd)
	return cmd
}
