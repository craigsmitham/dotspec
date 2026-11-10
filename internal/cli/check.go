package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "check reported findings"
}

func runCheckArgs(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("dotspec check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("path", ".", "project directory whose subtree is checked")
	format := flags.String("format", "text", "text or json")
	operands, err := parseFlagArgs(flags, args)
	if errors.Is(err, flag.ErrHelp) {
		_, err = fmt.Fprintln(out, "Usage: dotspec check [--path DIR] [--format text|json]\n\nChecks system and job specifications using the same discovery rules as list.\nWriting advice is nonfatal; optional job structure may remain absent.\nExit codes: 0 no structural errors, 1 structural findings, 2 invalid invocation or incomplete scan.")
		if err != nil {
			return &exitError{code: 2, err: err}
		}
		return nil
	}
	if err == nil && len(operands) != 0 {
		err = fmt.Errorf("unexpected argument %q; select a project directory with --path DIR", operands[0])
	}
	if err == nil && strings.TrimSpace(*path) == "" {
		err = fmt.Errorf("--path must not be empty")
	}
	if err == nil && *format != "text" && *format != "json" {
		err = fmt.Errorf("unknown check format %q; use text or json", *format)
	}
	if err != nil {
		return &exitError{code: 2, err: err}
	}
	report := check.Project(*path)
	if err := writeCheckReport(out, report, *format); err != nil {
		return &exitError{code: 2, err: err}
	}
	if code := report.ExitCode(); code != 0 {
		return &exitError{code: code}
	}
	return nil
}

func writeCheckReport(out io.Writer, report check.Report, format string) error {
	if format == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}
	var text strings.Builder
	fmt.Fprintf(&text, "Checking %s (system, job; %s discovery)\n", displayPath(report.Scope.Path), report.Scope.Discovery)
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
	if report.Complete && report.Summary.FilesDiscovered == 0 {
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
