package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// valueName is a string flag value that names itself in help, as in
// --part-of <path>, instead of pflag's generic "string".
type valueName struct {
	value string
	name  string
}

func (p *valueName) String() string     { return p.value }
func (p *valueName) Set(v string) error { p.value = v; return nil }
func (p *valueName) Type() string       { return p.name }

// stringFlag adds a string flag whose help shows its value as name, such as
// <path> or brief|expanded.
func stringFlag(flags *pflag.FlagSet, flag, shorthand, name, value, usage string) {
	flags.VarP(&valueName{value: value, name: name}, flag, shorthand, usage)
}

// flagValue reads a string flag added with stringFlag.
func flagValue(cmd *cobra.Command, flag string) string {
	if f := cmd.Flags().Lookup(flag); f != nil {
		return f.Value.String()
	}
	return ""
}

// usageError states what is wrong with an invocation and where to find its usage.
func usageError(cmd *cobra.Command, problem string) error {
	if cmd.HasParent() {
		problem = strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()+" ") + ": " + problem
	}
	return fmt.Errorf("%s; see %s --help", problem, cmd.CommandPath())
}

// subcommands validates a command group. Cobra suggests corrections only for
// groups without an Args validator, and such a group would accept a misspelled
// subcommand as an operand and print its help with status 0.
func subcommands(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	problem := fmt.Sprintf("unknown command %q for %q", args[0], cmd.CommandPath())
	if cmd.SuggestionsMinimumDistance <= 0 {
		cmd.SuggestionsMinimumDistance = 2
	}
	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 {
		return fmt.Errorf("%s; did you mean %s?", problem, quoteAll(suggestions))
	}
	return fmt.Errorf("%s; see %s --help", problem, cmd.CommandPath())
}

// misspelledHelp rejects a misspelled command followed by --help, which Cobra
// answers with the enclosing group's help and status 0.
func misspelledHelp(args []string) error {
	target, rest, err := NewRoot().Find(args)
	if err != nil || !target.HasSubCommands() {
		return nil
	}
	target.InitDefaultHelpFlag()
	if target.ParseFlags(rest) != nil {
		return nil
	}
	if help, _ := target.Flags().GetBool("help"); !help {
		return nil
	}
	return subcommands(target, target.Flags().Args())
}

func quoteAll(values []string) string {
	quoted := make([]string, len(values))
	for i, value := range values {
		quoted[i] = fmt.Sprintf("%q", value)
	}
	return strings.Join(quoted, " or ")
}

// operands accepts between minimum and maximum positional arguments (a
// negative maximum is unbounded) and names a missing one by its placeholder in
// the command's usage line.
func operands(minimum, maximum int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		switch {
		case len(args) < minimum:
			missing := "operand"
			if names := requiredOperands(cmd.Use); len(args) < len(names) {
				missing = names[len(args)]
			}
			return usageError(cmd, "missing "+missing)
		case maximum >= 0 && len(args) > maximum:
			return usageError(cmd, fmt.Sprintf("unexpected argument %q", args[maximum]))
		}
		return nil
	}
}

// requiredFlags names a missing required flag with its value placeholder, as
// operands names a missing operand. Cobra's own check runs later and gives no
// pointer to the command's usage.
func requiredFlags(cmd *cobra.Command) error {
	var missing []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if required := f.Annotations[cobra.BashCompOneRequiredFlag]; len(required) > 0 && required[0] == "true" && !f.Changed {
			value, _ := pflag.UnquoteUsage(f)
			missing = append(missing, strings.TrimSpace("--"+f.Name+" "+value))
		}
	})
	if len(missing) == 0 {
		return nil
	}
	return usageError(cmd, "missing "+strings.Join(missing, ", "))
}

// requiredOperands lists the placeholders written outside brackets, such as
// <slug> in "init <slug> [<path>]".
func requiredOperands(use string) []string {
	_, operands, _ := strings.Cut(use, " ")
	var names []string
	depth := 0
	for _, word := range strings.Fields(operands) {
		if depth == 0 && strings.HasPrefix(word, "<") {
			name, _, _ := strings.Cut(strings.TrimSuffix(word, "..."), "|")
			names = append(names, name)
		}
		depth += strings.Count(word, "[") - strings.Count(word, "]")
	}
	return names
}

// errorMessage renders a failure for standard error. A filesystem failure
// names its path and the problem in plain words, and paths near the working
// directory print relative to it.
func errorMessage(err error) string {
	message := err.Error()
	var pathError *fs.PathError
	if errors.As(err, &pathError) {
		message = strings.Replace(message, pathError.Error(), pathError.Path+": "+diagnostic.FileProblem(pathError), 1)
	}
	return relativePaths(message)
}

// relativePaths rewrites absolute paths within or up to two levels above the
// working directory. Other absolute paths, such as one the caller typed, stay
// as written.
func relativePaths(message string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return message
	}
	bases := []string{cwd}
	if resolved, err := filepath.EvalSymlinks(cwd); err == nil && resolved != cwd {
		bases = append(bases, resolved)
	}
	type prefix struct{ dir, relative string }
	var prefixes []prefix
	for _, base := range bases {
		dir, relative := base, "."
		for range 3 {
			if filepath.Dir(dir) == dir {
				break
			}
			prefixes = append(prefixes, prefix{dir, relative})
			dir, relative = filepath.Dir(dir), filepath.Join(relative, "..")
		}
	}
	// Replace longer prefixes first so a parent never claims a nested path.
	sort.SliceStable(prefixes, func(i, j int) bool { return len(prefixes[i].dir) > len(prefixes[j].dir) })
	for _, p := range prefixes {
		message = replacePath(message, p.dir, p.relative)
	}
	return message
}

// replacePath replaces dir where it appears as a whole path or as the start of
// one, never in the middle of a longer name.
func replacePath(message, dir, relative string) string {
	var out strings.Builder
	for {
		i := strings.Index(message, dir)
		if i < 0 {
			out.WriteString(message)
			return out.String()
		}
		before, after := message[:i], message[i+len(dir):]
		out.WriteString(before)
		switch {
		case before != "" && pathByte(before[len(before)-1]):
			out.WriteString(dir)
		case strings.HasPrefix(after, string(filepath.Separator)):
			if relative == "." {
				after = after[1:]
			} else {
				out.WriteString(relative)
			}
		case pathEnd(after):
			out.WriteString(relative)
			if strings.HasPrefix(after, ".") {
				// A closing full stop must not read as part of "." or "..".
				out.WriteString(string(filepath.Separator))
			}
		default:
			out.WriteString(dir)
		}
		message = after
	}
}

// pathEnd reports whether a path ends before rest, which may start with a
// closing full stop.
func pathEnd(rest string) bool {
	rest = strings.TrimPrefix(rest, ".")
	return rest == "" || !pathByte(rest[0])
}

func pathByte(c byte) bool {
	return c == '/' || c == '\\' || c == '.' || c == '-' || c == '_' ||
		c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// newHelpCommand replaces Cobra's help command, which prints the root help
// with status 0 for a misspelled topic.
func newHelpCommand() *cobra.Command {
	return &cobra.Command{Use: "help [command]", Short: "Help about any command",
		Long: "Show the help for any command, such as dotspec help check.",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, prefix string) ([]string, cobra.ShellCompDirective) {
			target, _, err := cmd.Root().Find(args)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var names []string
			for _, child := range target.Commands() {
				if child.IsAvailableCommand() && strings.HasPrefix(child.Name(), prefix) {
					names = append(names, child.Name()+"\t"+child.Short)
				}
			}
			return names, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			target, rest, err := cmd.Root().Find(args)
			if err != nil {
				return err
			}
			if len(rest) > 0 {
				return subcommands(target, rest)
			}
			target.InitDefaultHelpFlag()
			target.InitDefaultVersionFlag()
			return target.Help()
		},
	}
}
