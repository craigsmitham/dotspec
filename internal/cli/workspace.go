package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

type projectSelectionKey struct{}
type workspaceResult struct {
	SchemaVersion int             `json:"schema_version"`
	Workspace     string          `json:"workspace"`
	Projects      []projectResult `json:"projects"`
}
type projectResult struct {
	Name     string          `json:"name"`
	Root     string          `json:"root"`
	ExitCode int             `json:"exit_code"`
	Result   json.RawMessage `json:"result,omitempty"`
	Error    string          `json:"error,omitempty"`
}

func selectProject(cmd *cobra.Command, path string) (string, error) {
	if selected, ok := cmd.Context().Value(projectSelectionKey{}).(string); ok {
		return selected, nil
	}
	name, _ := cmd.Flags().GetString("project")
	if cmd.Flags().Changed("project") && name == "" {
		return "", fmt.Errorf("--project must not be empty")
	}
	if name != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		w, err := specfile.FindWorkspace(cwd)
		if err != nil {
			return "", err
		}
		if w == nil {
			return "", fmt.Errorf("--project requires a dotspec.yaml workspace")
		}
		for _, p := range w.Projects {
			if p.Name == name {
				selected := p.Root
				if cmd.Flags().Changed("path") {
					selected = filepath.Join(p.Root, path)
					if filepath.IsAbs(path) {
						selected = path
					}
				}
				selected, err = specfile.Directory(selected)
				if err != nil {
					return "", err
				}
				owner, ok := w.Containing(selected)
				if !ok || owner.Name != name {
					return "", fmt.Errorf("--path must stay within project %s", name)
				}
				return selected, nil
			}
		}
		return "", fmt.Errorf("unknown project %q in %s", name, w.Path)
	}
	// Preserve existing diagnostic reports for invalid path selectors.
	dir, err := specfile.Directory(path)
	if err != nil {
		return path, nil
	}
	w, p, err := specfile.ProjectContext(dir)
	if err != nil {
		return "", err
	}
	if w != nil && !cmd.Flags().Changed("path") {
		return p.Root, nil
	}
	return path, nil
}

func addProjectSelection(cmd *cobra.Command) {
	cmd.Flags().String("project", "", "select a project registered in dotspec.yaml; --path is relative to its root")
	switch cmd.Name() {
	case "check", "sources", "list", "export", "docs", "trace", "overview", "needs", "review", "pegs":
	default:
		return
	}
	if len(strings.Fields(cmd.Use)) != 1 {
		return
	}
	cmd.Flags().Bool("all", false, "inspect all registered workspace projects separately")
	cmd.MarkFlagsMutuallyExclusive("project", "all")
	cmd.MarkFlagsMutuallyExclusive("path", "all")
	run := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all {
			return run(cmd, args)
		}
		if len(args) > 0 {
			return fmt.Errorf("--all cannot select an individual specification")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		w, err := specfile.FindWorkspace(cwd)
		if err != nil {
			return err
		}
		if w == nil {
			return fmt.Errorf("--all requires a dotspec.yaml workspace")
		}
		format := "text"
		if cmd.Flags().Lookup("format") != nil {
			format, err = outputFormat(cmd)
			if err != nil {
				return err
			}
		}
		if cmd.Name() == "export" {
			format = "json"
		}
		result := workspaceResult{SchemaVersion: 1, Workspace: filepath.ToSlash(w.Path), Projects: []projectResult{}}
		out, ctx := cmd.OutOrStdout(), cmd.Context()
		defer func() { cmd.SetOut(out); cmd.SetContext(ctx) }()
		code := 0
		for _, p := range w.Projects {
			var buffer bytes.Buffer
			cmd.SetOut(&buffer)
			cmd.SetContext(context.WithValue(ctx, projectSelectionKey{}, p.Root))
			runErr := run(cmd, args)
			entry := projectResult{Name: p.Name, Root: filepath.ToSlash(p.Root), ExitCode: status(runErr)}
			if entry.ExitCode > code {
				code = entry.ExitCode
			}
			if runErr != nil {
				var coded *exitError
				if e, ok := runErr.(*exitError); ok {
					coded = e
				}
				if coded == nil || coded.err != nil {
					entry.Error = runErr.Error()
				}
			}
			if format == "json" {
				if buffer.Len() > 0 {
					if !json.Valid(buffer.Bytes()) {
						return fmt.Errorf("project %s produced invalid JSON", p.Name)
					}
					entry.Result = append(json.RawMessage(nil), buffer.Bytes()...)
				}
			} else {
				if _, err := fmt.Fprintf(out, "# Project: %s\n\n%s\n", p.Name, buffer.String()); err != nil {
					return err
				}
				if entry.Error != "" {
					if _, err := fmt.Fprintf(out, "Error: %s\n\n", entry.Error); err != nil {
						return err
					}
				}
			}
			result.Projects = append(result.Projects, entry)
		}
		if format == "json" {
			if err := writeJSON(out, result); err != nil {
				return err
			}
		}
		if code != 0 {
			return &exitError{code: code}
		}
		return nil
	}
}
