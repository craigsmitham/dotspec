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
	"github.com/spf13/cobra"
)

// Release builds set these through GoReleaser's ldflags.
var version = "dev"
var commit = ""

// Root help states the exit-code policy once. A leaf that can exit 1 says so.
const (
	exitPolicy   = "Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure."
	findingsExit = "Exits 1 when error findings are reported."
)

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
	ctx := context.Background()
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
	// -C changes the process working directory for this invocation only.
	if cwd, err := os.Getwd(); err == nil {
		defer func() { _ = os.Chdir(cwd) }()
	}
	err := misspelledHelp(args)
	if err == nil {
		err = root.ExecuteContext(ctx)
	}
	if out.err != nil {
		err = &exitError{code: 2, err: out.err}
	}
	code := status(err)
	if err != nil {
		var coded *exitError
		if !errors.As(err, &coded) || coded.err != nil {
			message := errorMessage(err)
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

func terminal(w io.Writer) bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return false
	}
	f, ok := w.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(f.Fd())
}

// build is the identity Go records in the binary. A release commit set
// through ldflags takes the place of the recorded VCS settings.
type build struct {
	module   string // main module version, such as v1.2.3, when Go recorded one
	revision string // source commit
	dirty    bool   // the recorded commit had local modifications
}

func readBuild() build {
	b := build{revision: commit}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return b
	}
	if info.Main.Version != "(devel)" {
		b.module = info.Main.Version
	}
	if commit != "" {
		return b
	}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			b.revision = setting.Value
		case "vcs.modified":
			b.dirty = setting.Value == "true"
		}
	}
	return b
}

func buildVersion() string {
	b := readBuild()
	v, revision := version, b.revision
	if v == "dev" && b.module != "" {
		v = b.module
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if revision != "" {
		v += " (" + revision + ")"
	}
	return v
}

// citationRevision names the dotspec source that applied bundled quality
// models cite. A modified working tree has none, so apply cites the embedded
// bytes by digest.
func citationRevision() string {
	b := readBuild()
	switch {
	case b.dirty:
		return ""
	case b.revision != "":
		return b.revision
	case strings.HasPrefix(b.module, "v"):
		return b.module
	}
	return ""
}

// NewRoot is side-effect free: documentation and execution share this tree.
func NewRoot() *cobra.Command {
	root := &cobra.Command{Use: "dotspec", Short: "Specify what software must do and why, alongside its code",
		Long: "Specify what software must do and why, alongside its code: author Markdown specifications, check them and the artifacts that carry their intent, and build one generated specification a reader can navigate.\n\n" +
			"Kinds: run dotspec new --help. Handbook: run dotspec docs.\n\n" + exitPolicy,
		Example: "dotspec init \"Booking desk\"\ndotspec new requirement fast-save\ndotspec build\ndotspec check",
		Version: buildVersion(), SilenceUsage: true, SilenceErrors: true,
		Args: subcommands, RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	root.SetVersionTemplate("dotspec {{.Version}}\n")
	addChdir(root)
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return fmt.Errorf("%w; see %s --help", err, cmd.CommandPath())
	})
	for _, g := range []struct {
		id       string
		commands []*cobra.Command
	}{
		{"Author:", []*cobra.Command{newInitCommand(), newNewCommand()}},
		{"Inspect:", []*cobra.Command{newCheckCommand(), newListCommand(), newSourcesCommand()}},
		{"Publish:", []*cobra.Command{newBuildCommand(), newSnapshotCommand()}},
		{"Help:", []*cobra.Command{newDocsCommand(), newCompletionCommand()}},
	} {
		root.AddGroup(&cobra.Group{ID: g.id, Title: g.id})
		for _, cmd := range g.commands {
			cmd.GroupID = g.id
			root.AddCommand(cmd)
		}
	}
	root.AddCommand(newManCommand())
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()
	help := newHelpCommand()
	help.GroupID = "Help:"
	root.SetHelpCommand(help)
	root.InitDefaultHelpCmd()
	indentExamples(root)
	return root
}

// indentExamples sets examples off from their heading in plain help, as
// Cobra's own commands do. Fang and the generated reference trim the indent.
func indentExamples(cmd *cobra.Command) {
	if cmd.Example != "" {
		cmd.Example = indent(cmd.Example)
	}
	for _, child := range cmd.Commands() {
		indentExamples(child)
	}
}

func group(use, short string) *cobra.Command {
	return &cobra.Command{Use: use, Short: short, Args: subcommands,
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() }}
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
