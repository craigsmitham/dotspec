// Package cli implements dotspec's command-line interface.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/spf13/cobra"
)

// Release builds set these through GoReleaser's ldflags.
var version = "dev"
var commit = ""

const exitPolicy = "Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure."

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "specification findings"
}
func (e *exitError) Unwrap() error { return e.err }

// Run owns process output and status; handlers never exit the process.
func Run(args []string, stdout, stderr io.Writer) int {
	return RunContext(context.Background(), args, stdout, stderr)
}

func RunContext(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	out, errOut := &checkedWriter{Writer: stdout}, &checkedWriter{Writer: stderr}
	root := NewRoot()
	root.SetArgs(args)
	root.SetOut(out)
	root.SetErr(errOut)
	// Fang is used only to render requested help. Its error path probes the
	// process stdout for colors, which would contaminate a terminal JSON report.
	if terminal(stdout) {
		root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
			help := NewRoot()
			help.SetOut(out)
			help.SetErr(errOut)
			help.SetArgs(append(strings.Fields(cmd.CommandPath())[1:], "--help"))
			if err := fang.Execute(ctx, help, fang.WithoutManpage(), fang.WithVersion(help.Version),
				fang.WithErrorHandler(func(_ io.Writer, _ fang.Styles, _ error) {})); err != nil && out.err == nil {
				out.err = err
			}
		})
	}
	err := root.ExecuteContext(ctx)
	if out.err != nil {
		err = &exitError{code: 2, err: out.err}
	}
	code := status(err)
	if err != nil {
		var coded *exitError
		if !errors.As(err, &coded) || coded.err != nil {
			message := err.Error()
			if terminal(stderr) {
				message = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)).Render(message)
			}
			_, _ = fmt.Fprintf(errOut, "dotspec: %s\n", message)
		}
	}
	if errOut.err != nil {
		return 2
	}
	return code
}

func status(err error) int {
	if err == nil {
		return 0
	}
	var coded *exitError
	if errors.As(err, &coded) {
		return coded.code
	}
	var finding *diagnostic.Diagnostic
	if errors.As(err, &finding) {
		return 1
	}
	return 2
}

// checkedWriter catches output failures even in upstream help/completion paths
// whose callback interfaces cannot return an error.
type checkedWriter struct {
	io.Writer
	err error
}

func (w *checkedWriter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err == nil && n != len(p) {
		err = io.ErrShortWrite
	}
	if err != nil && w.err == nil {
		w.err = err
	}
	return n, err
}

func (w *checkedWriter) Fd() uintptr {
	if f, ok := w.Writer.(interface{ Fd() uintptr }); ok {
		return f.Fd()
	}
	return ^uintptr(0)
}

func terminal(w io.Writer) bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return false
	}
	f, ok := w.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(f.Fd())
}

func buildVersion() string {
	v, revision := version, commit
	if info, ok := debug.ReadBuildInfo(); ok {
		if v == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
		if revision == "" {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
				}
			}
		}
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if revision != "" {
		v += " (" + revision + ")"
	}
	return v
}

// NewRoot is side-effect free: documentation and execution share this tree.
func NewRoot() *cobra.Command {
	root := &cobra.Command{Use: "dotspec", Short: "Manage software requirements alongside your code",
		Long:    "Manage human-readable software requirements in *.spec.md files.\n\n" + exitPolicy,
		Example: "dotspec system init\ndotspec job init develop-requirements\ndotspec check",
		Version: buildVersion(), SilenceUsage: true, SilenceErrors: true,
		Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	root.SetVersionTemplate("dotspec {{.Version}}\n")
	root.AddCommand(newBusinessRuleCommand(), newSystemCommand(), newJobCommand(), newUseCaseCommand(), newConceptCommand(), newQualityCommand(), newRequirementCommand(), newCheckCommand(), newSourcesCommand(), newModelCommand(), newManCommand(), newCompletionCommand())
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	return root
}

func group(use, short string) *cobra.Command {
	return &cobra.Command{Use: use, Short: short, Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() }}
}

func projectPath(cmd *cobra.Command) (string, error) {
	path, err := cmd.Flags().GetString("path")
	if err == nil && strings.TrimSpace(path) == "" {
		err = fmt.Errorf("--path must not be empty")
	}
	if err != nil {
		return path, err
	}
	// Checks and inventories include scope failures in reports. Other commands report an
	// invalid project selector as an invocation error rather than a spec finding.
	if cmd.Name() != "check" && cmd.Name() != "sources" {
		resolved, scopeErr := specfile.Scope(path)
		if scopeErr != nil {
			return path, &exitError{code: 2, err: scopeErr}
		}
		path = resolved
	}
	return path, nil
}

func addPath(cmd *cobra.Command) {
	cmd.Flags().String("path", ".", "project directory to inspect or initialize")
	_ = cmd.MarkFlagDirname("path")
}

func writeTree(out *strings.Builder, nodes []*system.Node, prefix string) {
	for i, node := range nodes {
		branch, next := "├── ", "│   "
		if i == len(nodes)-1 {
			branch, next = "└── ", "    "
		}
		fmt.Fprintf(out, "%s%s%s  (%s)\n", prefix, branch, node.System.Name, displayPath(node.System.Path))
		writeTree(out, node.Children, prefix+next)
	}
}

func displayPath(path string) string {
	if cwd, err := os.Getwd(); err == nil {
		if resolved, err := filepath.EvalSymlinks(cwd); err == nil {
			cwd = resolved
		}
		if relative, err := filepath.Rel(cwd, path); err == nil {
			return filepath.ToSlash(relative)
		}
	}
	return filepath.ToSlash(path)
}

func indent(value string) string {
	return "  " + strings.ReplaceAll(value, "\n", "\n  ")
}
