package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/guide"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "check [source-path]", Short: "Check specifications and reference documents in a project subtree",
		Long:    "Check authored specifications and reference documents using the same discovery rules as list. A whole-project check also inspects generated Markdown for projects declared with dotspec.yaml unless report.enabled is false. These automated authoring and selected requirements-quality checks support requirements verification; passing them does not establish full requirements verification, requirements validation, or system verification/validation.\n\n" + exitPolicy,
		Example: "dotspec check\ndotspec check --path internal --json", Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := projectPath(cmd)
			if err != nil {
				return err
			}
			format, err := outputFormat(cmd)
			if err != nil {
				return err
			}
			if len(args) > 0 {
				types := supportedTypes()
				for _, kind := range types {
					if kind.Layout.Match(args[0]) {
						return emitCheckReport(cmd.OutOrStdout(), check.Focused(path, args[0], kind, types...), format)
					}
				}
				return fmt.Errorf("unknown source document type: %s", args[0])
			}
			return emitCheckReport(cmd.OutOrStdout(), inspectProject(path), format)
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
	if report.Scope.Project != nil {
		fmt.Fprintf(&text, "Project: %s\n", report.Scope.Project.Name)
	}
	fmt.Fprintf(&text, "Checking %s (%s; %s discovery)\n", displayPath(report.Scope.Path), strings.Join(report.Scope.Types, ", "), report.Scope.Discovery)
	if report.GeneratedSpec != "" {
		fmt.Fprintf(&text, "Generated specification: %s\n", report.GeneratedSpec)
	}
	for _, path := range report.ContextFiles {
		fmt.Fprintf(&text, "Context: %s\n", displayPath(path))
	}
	for _, d := range report.Diagnostics {
		fmt.Fprintf(&text, "\n%s\n  %s [%s]", diagnosticPath(d.Location), d.Severity, d.Rule)
		if d.Location.Element != "" {
			fmt.Fprintf(&text, " %s", d.Location.Element)
		}
		fmt.Fprintf(&text, "\n%s\n", indent(d.Message))
		if guide.HasRule(d.Rule) {
			fmt.Fprintf(&text, "  Guide: dotspec docs %s\n", d.Rule)
		}
		for _, related := range d.Related {
			fmt.Fprintf(&text, "  Related: %s", diagnosticPath(related))
			if related.Element != "" {
				fmt.Fprintf(&text, " (%s)", related.Element)
			}
			text.WriteByte('\n')
		}
	}
	if report.Complete && report.Summary.FilesDiscovered == 0 && report.Summary.Errors == 0 {
		text.WriteString("\nNo supported source files found.\n")
	}
	fmt.Fprintf(&text, "\nChecked %d of %d discovered source files; additional context files: %d.\nErrors: %d; advisory findings: %d.\n",
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
	return []check.Type{qualitymodel.CheckType(), businessrule.CheckType(), system.CheckType(), system.ComponentCheckType(), system.ProductCheckType(system.FeatureLayout), system.ProductCheckType(system.CapabilityLayout), system.AttributeCheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType()}
}
