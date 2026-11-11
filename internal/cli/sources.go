package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/spf13/cobra"
)

func newSourcesCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "sources", Short: "List specification files, system context, and directly referenced files",
		Long:    "Inventory sources using the same scan and reference checks as check. File references do not assign authority or classify their contents.\n\n" + exitPolicy,
		Example: "dotspec sources\ndotspec sources --path internal --json", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			r := check.Sources(path, supportedTypes()...)
			if err := writeSources(cmd.OutOrStdout(), r, format); err != nil {
				return &exitError{code: 2, err: err}
			}
			if code := r.ExitCode(); code != 0 {
				return &exitError{code: code}
			}
			return nil
		}}
	addPath(cmd)
	addOutput(cmd, false)
	return cmd
}

func writeSources(out io.Writer, r check.SourceReport, format string) error {
	if format == "json" {
		return writeJSON(out, r)
	}
	var text strings.Builder
	if r.Scope.Project != nil {
		fmt.Fprintf(&text, "Project: %s\n", r.Scope.Project.Name)
	}
	fmt.Fprintf(&text, "Sources in %s (%s discovery)\n", displayPath(r.Scope.Path), r.Scope.Discovery)
	for _, file := range r.Files {
		fmt.Fprintf(&text, "\n%s\n", displayPath(file.Path))
		if file.Project != nil {
			fmt.Fprintf(&text, "  Project: %s\n", file.Project.Name)
		}
		for _, location := range file.ReferencedBy {
			fmt.Fprintf(&text, "  Referenced by %s", diagnosticPath(location))
			if location.Element != "" {
				fmt.Fprintf(&text, " (%s)", location.Element)
			}
			text.WriteByte('\n')
		}
	}
	if len(r.Files) == 0 {
		text.WriteString("\nNo source files found.\n")
	}
	for _, d := range r.Diagnostics {
		fmt.Fprintf(&text, "\n%s\n  %s [%s] %s\n", diagnosticPath(d.Location), d.Severity, d.Rule, d.Message)
		for _, related := range d.Related {
			fmt.Fprintf(&text, "  Related: %s\n", diagnosticPath(related))
		}
	}
	fmt.Fprintf(&text, "\nSource files: %d.\n", len(r.Files))
	if r.Complete {
		text.WriteString("Scan complete.\n")
	} else {
		text.WriteString("Scan incomplete.\n")
	}
	_, err := io.WriteString(out, text.String())
	return err
}
