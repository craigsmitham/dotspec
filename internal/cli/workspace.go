package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

// scope is where a command starts when no positional path narrows it. Each
// command declares its scope when it is built.
type scope int

const (
	// projectScope reads the enclosing project from its manifest root.
	projectScope scope = iota
	// workingScope creates files and looks up the nearest subject from the
	// working directory, which may lie inside a component.
	workingScope
)

// dir resolves the directory a command reads or writes. A path narrows the
// scope and is relative to the working directory. Inspections report a path
// that cannot be read, so an unreadable path is returned unresolved.
func (s scope) dir(path string) (string, error) {
	narrowed := path != ""
	if !narrowed {
		path = "."
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		if narrowed {
			return path, nil
		}
		return "", err
	}
	p, err := specfile.FindProject(dir)
	if err != nil {
		return "", err
	}
	if p == nil {
		found, err := specfile.HasProjects(dir)
		if err != nil {
			return "", err
		}
		if found {
			return "", fmt.Errorf("%s is outside any specification project; run dotspec -C <project-dir>, or dotspec check --all", displayPath(dir))
		}
		return dir, nil
	}
	if narrowed || s == workingScope {
		return dir, nil
	}
	return p.Root, nil
}

// optionalPath returns a command's [<path>] operand, or "" when it is omitted.
func optionalPath(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

// addChdir gives the root the global -C, --chdir. Every relative path,
// positional or flag, then resolves against that directory.
func addChdir(root *cobra.Command) {
	stringFlag(root.PersistentFlags(), "chdir", "C", "<dir>", "", "run as if dotspec was started in <dir>")
	_ = root.MarkPersistentFlagDirname("chdir")
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if err := requiredFlags(cmd); err != nil {
			return err
		}
		dir := flagValue(cmd, "chdir")
		if !cmd.Flags().Changed("chdir") {
			return nil
		}
		if dir == "" {
			return usageError(cmd, "--chdir requires a directory")
		}
		info, err := os.Stat(dir)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return fmt.Errorf("--chdir %s: directory does not exist", dir)
		case err != nil:
			return fmt.Errorf("--chdir %s: %w", dir, err)
		case !info.IsDir():
			return fmt.Errorf("--chdir %s: not a directory", dir)
		}
		return os.Chdir(dir)
	}
}

// workspaceResult is check --all's result: each discovered project's report,
// checked separately.
type workspaceResult struct {
	SchemaVersion int             `json:"schema_version"`
	Workspace     string          `json:"workspace"`
	Projects      []projectResult `json:"projects"`
}
type projectResult struct {
	Name     string          `json:"name"`
	Root     string          `json:"root"`
	ExitCode int             `json:"exit_code"`
	Result   json.RawMessage `json:"result"`
}

// checkWorkspace checks every project discovered from the working directory
// and exits with the highest project exit code.
func checkWorkspace(cmd *cobra.Command, asJSON bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	w, err := specfile.FindWorkspace(cwd)
	if err != nil {
		return err
	}
	if len(w.Projects) == 0 {
		return fmt.Errorf("no specification projects found in %s", displayPath(w.Root))
	}
	out := cmd.OutOrStdout()
	result := workspaceResult{SchemaVersion: schemaVersion, Workspace: filepath.ToSlash(w.Root), Projects: []projectResult{}}
	code := 0
	for _, p := range w.Projects {
		report := inspectProject(p.Root)
		var buffer bytes.Buffer
		if err := writeCheckReport(&buffer, report, asJSON); err != nil {
			return err
		}
		entry := projectResult{Name: p.Name, Root: filepath.ToSlash(p.Root), ExitCode: report.ExitCode()}
		code = max(code, entry.ExitCode)
		if asJSON {
			entry.Result = append(json.RawMessage(nil), buffer.Bytes()...)
			result.Projects = append(result.Projects, entry)
			continue
		}
		if _, err := fmt.Fprintf(out, "# Project: %s\n\n%s\n", p.Name, buffer.String()); err != nil {
			return err
		}
	}
	if asJSON {
		if err := writeJSON(out, result); err != nil {
			return err
		}
	}
	if code != 0 {
		return &exitError{code: code}
	}
	return nil
}
