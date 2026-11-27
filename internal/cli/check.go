package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/guide"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "check [<path>]", Short: "Check specifications and reference documents in a project subtree",
		Long:    "Check authored specifications and the reference documents they cite, using the same discovery rules as list. With no <path>, check covers the enclosing project; a directory narrows it to that subtree, and a specification file checks that file and the context it requests. A whole-project check in a project declared with dotspec.yaml also checks that the generated specification is current, unless build.enabled is false. With --all, check covers every project discovered in the workspace separately and exits with the highest project exit code.\n\nA passing check means no checked file has an error finding under dotspec's automated authoring rules. " + findingsExit + "\n\nLearn more: dotspec docs docs/reference/inspection.md#check-specifications",
		Example: "dotspec check\ndotspec check internal --json", Args: operands(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			asJSON := jsonOutput(cmd)
			if all, _ := cmd.Flags().GetBool("all"); all {
				if len(args) > 0 {
					return usageError(cmd, "--all checks every project; omit <path>")
				}
				return checkWorkspace(cmd, asJSON)
			}
			target := optionalPath(args)
			if kind := check.MatchType(target, speckinds.All()); target != "" && kind != nil {
				// A specification file is checked from its own directory, so it
				// belongs to its enclosing project wherever the check starts.
				dir, err := workingScope.dir(filepath.Dir(target))
				if err != nil {
					return err
				}
				return emitCheckReport(cmd.OutOrStdout(), check.Focused(dir, filepath.Base(target), *kind, speckinds.All()...), asJSON)
			}
			if info, err := os.Stat(target); target != "" && err == nil && !info.IsDir() {
				return fmt.Errorf("unknown source document type: %s", target)
			}
			dir, err := projectScope.dir(target)
			if err != nil {
				return err
			}
			return emitCheckReport(cmd.OutOrStdout(), inspectProject(dir), asJSON)
		},
	}
	cmd.Flags().Bool("all", false, "check every discovered project separately")
	addJSON(cmd)
	return cmd
}

func emitCheckReport(out io.Writer, report check.Report, asJSON bool) error {
	if err := writeCheckReport(out, report, asJSON); err != nil {
		return &exitError{code: 2, err: err}
	}
	if code := report.ExitCode(); code != 0 {
		return &exitError{code: code}
	}
	return nil
}

func writeCheckReport(out io.Writer, report check.Report, asJSON bool) error {
	if asJSON {
		return writeJSON(out, report)
	}
	var text strings.Builder
	for _, d := range report.Diagnostics {
		fmt.Fprintf(&text, "%s\n  %s [%s]", diagnosticPath(d.Location), d.Severity, d.Rule)
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
		text.WriteByte('\n')
	}
	if report.Complete && report.Summary.FilesDiscovered == 0 && report.Summary.Errors == 0 {
		text.WriteString("No supported source files found.\n")
	}
	completeness := "scan complete"
	if !report.Complete {
		completeness = "scan incomplete"
	}
	context := report.ExistingContextFiles()
	fmt.Fprintf(&text, "Checked %s in %s (%s): %s, %s; %s.\n",
		count(report.Summary.FilesChecked, "file", "files"), scopeName(report.Scope), completeness,
		count(report.Summary.Errors, "error", "errors"), count(report.Summary.Advice, "advisory", "advisories"),
		count(len(context), "context file", "context files"))
	if state := generatedSpecState(report); state != "" {
		fmt.Fprintf(&text, "Generated specification: %s\n", state)
	}
	if len(context) > 0 {
		text.WriteString("\nContext files:\n")
		for _, path := range context {
			fmt.Fprintf(&text, "  %s\n", displayPath(path))
		}
	}
	_, err := io.WriteString(out, text.String())
	return err
}

// scopeName names a report's scope for text output: the project when the
// scope is the whole project, otherwise the directory.
func scopeName(scope check.Scope) string {
	if scope.Project != nil && scope.Path == scope.Project.Root && scope.Discovery != "direct" {
		return "project " + scope.Project.Name
	}
	return displayPath(scope.Path)
}

// generatedSpecState states in words whether a check covered the generated
// specification, or "" when the project does not build one.
func generatedSpecState(report check.Report) string {
	switch report.GeneratedSpec {
	case "checked":
		return "checked"
	case "out-of-scope":
		if report.Scope.Discovery == "direct" {
			return "not checked (focused check)"
		}
		return "not checked (subtree check)"
	case "blocked-by-source":
		return "not checked (source errors prevent building it)"
	case "incomplete":
		return "check incomplete (operational error)"
	}
	return ""
}

func count(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

func diagnosticPath(location diagnostic.Location) string {
	location.Path = displayPath(location.Path)
	return location.String()
}
