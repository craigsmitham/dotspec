// Package cli implements dotspec's command-line interface.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/system"
)

const usage = `Usage: dotspec system <command> [--path DIR]

Commands:
  init   Create a system.spec.md template in the selected directory
  show   Describe the nearest enclosing system and identify its parent
  list   Display the hierarchy of systems within the selected subtree

--path defaults to the current directory. Use --help for command help.
`

func Run(args []string, stdout, stderr io.Writer) int {
	if err := run(args, stdout); err != nil {
		fmt.Fprintf(stderr, "dotspec: %v\n", err)
		return 1
	}
	return 0
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 || isHelp(args[0]) {
		_, err := fmt.Fprint(out, usage)
		return err
	}
	if args[0] != "system" {
		return fmt.Errorf("unknown command %q; use dotspec --help", args[0])
	}
	if len(args) == 1 || isHelp(args[1]) {
		_, err := fmt.Fprint(out, usage)
		return err
	}
	command := args[1]
	if command != "init" && command != "show" && command != "list" {
		return fmt.Errorf("unknown system command %q; use dotspec system --help", command)
	}
	flags := flag.NewFlagSet("dotspec system "+command, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	path := flags.String("path", ".", "directory to inspect or initialize")
	if err := flags.Parse(args[2:]); errors.Is(err, flag.ErrHelp) {
		_, err := fmt.Fprintf(out, "Usage: dotspec system %s [--path DIR]\n\n--path defaults to the current directory.\n", command)
		return err
	} else if err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected argument %q; select a directory with --path DIR", flags.Arg(0))
	}
	if strings.TrimSpace(*path) == "" {
		return fmt.Errorf("--path must not be empty")
	}
	switch command {
	case "init":
		target, err := system.Init(*path)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "Created %s\nEdit the system name, Purpose, and Scope.\n", displayPath(target))
		return err
	case "show":
		spec, parent, err := system.Show(*path)
		if err != nil {
			return err
		}
		parentText := "None"
		if parent != nil {
			parentText = fmt.Sprintf("%s (%s)", parent.Name, displayPath(parent.Path))
		}
		_, err = fmt.Fprintf(out, "%s\n\nPurpose\n%s\n\nScope\n%s\n\nParent system: %s\nSource: %s\n",
			spec.Name, indent(spec.Purpose), indent(spec.Scope), parentText, displayPath(spec.Path))
		return err
	case "list":
		nodes, err := system.List(*path)
		if err != nil {
			return err
		}
		if len(nodes) == 0 {
			_, err := fmt.Fprintln(out, "No systems found in the selected subtree.")
			return err
		}
		var output strings.Builder
		output.WriteString("Systems\n")
		writeTree(&output, nodes, "")
		_, err = io.WriteString(out, output.String())
		return err
	}
	return nil
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

func isHelp(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}
