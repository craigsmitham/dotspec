package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "check", Short: "Check all supported specifications in a project subtree",
		Long:    "Check system, job, use-case, concept, requirement, and quality specifications using the same discovery rules as list.\n\n" + exitPolicy,
		Example: "dotspec check\ndotspec check --path internal --json", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			return emitCheckReport(cmd.OutOrStdout(), check.Project(path, supportedTypes()...), format)
		},
	}
	addPath(cmd)
	addOutput(cmd, false)
	return cmd
}

func emitCheckReport(out io.Writer, report check.Report, format string) error {
	if err := writeCheckReport(out, report, format); err != nil {
		return &exitError{code: 2, err: err}
	}
	if code := report.ExitCode(); code != 0 {
		return &exitError{code: code}
	}
	return nil
}

func writeCheckReport(out io.Writer, report check.Report, format string) error {
	if format == "json" {
		return writeJSON(out, report)
	}
	var text strings.Builder
	fmt.Fprintf(&text, "Checking %s (%s; %s discovery)\n", displayPath(report.Scope.Path), strings.Join(report.Scope.Types, ", "), report.Scope.Discovery)
	for _, path := range report.ContextFiles {
		fmt.Fprintf(&text, "Context: %s\n", displayPath(path))
	}
	for _, d := range report.Diagnostics {
		fmt.Fprintf(&text, "\n%s\n  %s [%s]", diagnosticPath(d.Location), d.Severity, d.Rule)
		if d.Location.Element != "" {
			fmt.Fprintf(&text, " %s", d.Location.Element)
		}
		fmt.Fprintf(&text, "\n%s\n", indent(d.Message))
		for _, related := range d.Related {
			fmt.Fprintf(&text, "  Related: %s", diagnosticPath(related))
			if related.Element != "" {
				fmt.Fprintf(&text, " (%s)", related.Element)
			}
			text.WriteByte('\n')
		}
	}
	if report.Complete && report.Summary.FilesDiscovered == 0 && report.Summary.Errors == 0 {
		text.WriteString("\nNo supported specifications found.\n")
	}
	fmt.Fprintf(&text, "\nChecked %d of %d discovered specifications; additional context files: %d.\nErrors: %d; advisory findings: %d.\n",
		report.Summary.FilesChecked, report.Summary.FilesDiscovered, report.Summary.ContextFilesChecked, report.Summary.Errors, report.Summary.Advice)
	if !report.Complete {
		text.WriteString("Scan incomplete.\n")
	} else {
		text.WriteString("Scan complete.\n")
	}
	_, err := io.WriteString(out, text.String())
	return err
}

func diagnosticPath(location diagnostic.Location) string {
	location.Path = displayPath(location.Path)
	return location.String()
}

func supportedTypes() []check.Type {
	return []check.Type{system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), quality.CheckType(), requirement.CheckType()}
}
