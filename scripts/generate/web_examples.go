package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotspec-md/dotspec/internal/cli"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

const restartReport = "examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md"
const assumptionReport = "examples/todo/specification/assumptions.md"

var reportRoutes = map[string]string{restartReport: "examples/restart-requirement", assumptionReport: "examples/storage-assumption"}

func publishReport(file string) (generatedFile, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return generatedFile{}, err
	}
	_, body, ok := strings.Cut(string(specfile.Normalize(data)), "\n---\n")
	if !ok {
		return generatedFile{}, fmt.Errorf("missing report frontmatter: %s", file)
	}
	_, body, ok = strings.Cut(body, "\n# ")
	if !ok {
		return generatedFile{}, fmt.Errorf("missing report title: %s", file)
	}
	title, body, _ := strings.Cut(body, "\n")
	var lines []string
	for _, line := range strings.Split(body, "\n") {
		if !strings.HasPrefix(line, "<!-- dotspec-doc ") {
			lines = append(lines, line)
		}
	}
	content, err := webGuide(file, []byte(strings.Join(lines, "\n")))
	if err != nil {
		return generatedFile{}, err
	}
	if file == restartReport {
		content += "\n\n[Run the restart acceptance test](/examples/restart#run-the-two-checks).\n"
	}
	return webPage(reportRoutes[file], title, file, content), nil
}

// Capture the CLI's real streams and exit codes in a fresh project.
func checkWalkthrough() (content string, err error) {
	dir, err := os.MkdirTemp("", "dotspec-check-example-")
	if err != nil {
		return "", err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	run := func(code int, args ...string) (string, error) {
		var out, diagnostics bytes.Buffer
		got := cli.Run(append([]string{"-C", dir}, args...), &out, &diagnostics)
		if got != code {
			return "", fmt.Errorf("%v exited %d: %s%s", args, got, out.String(), diagnostics.String())
		}
		return strings.TrimSpace(strings.ReplaceAll(out.String(), "\r\n", "\n")), nil
	}
	if _, err = run(0, "init", "Booking desk"); err != nil {
		return "", err
	}
	if _, err = run(0, "build"); err != nil {
		return "", err
	}
	file := filepath.Join(dir, "specs", "booking-desk.system.spec.md")
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	if err = os.WriteFile(file, bytes.ReplaceAll(data, []byte("Describe the system's boundary and what it includes."), []byte("The desk takes reservations for one hotel.")), 0644); err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("\n## Check an edited project\n\nFrom the project root, check after editing the system scope. The previously built report still contains the earlier wording.\n")
	for _, step := range []struct {
		name string
		code int
	}{{"check", 1}, {"build", 0}, {"check", 0}} {
		out, err := run(step.code, step.name)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "\n```sh\ndotspec %s\n```\n\n```text nocopy\n%s\n```\n", step.name, out)
	}
	b.WriteString("\n[Repair a stale specification](/troubleshooting/spec-stale).\n\n## Use the result in CI\n\nRun this as the final command in a shell step. Redirection saves the report and preserves the process exit status.\n\n```sh\ndotspec check --json > check-report.json\n```\n\nExit 0 means the check completed successfully and may include advice. Exit 1 means specification findings prevented success. Exit 2 covers invalid invocations, missing input, discovery, file access, or output failures. Reports go to stdout, including incomplete scans. Invocation and output errors go to stderr. Inspect the scan status before relying on a partial report. A passing check assesses the authored files and report. Exercise the application separately to verify its behavior. [Inspection contract](" + githubSource + "docs/reference/inspection.md) · [Restart example](/examples/restart).\n\nRead the scope and findings guide offline with `dotspec docs docs/reference/inspection.md#check-specifications`.\n")
	return b.String(), nil
}
