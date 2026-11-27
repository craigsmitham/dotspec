package cli

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/spf13/cobra"
	mdast "github.com/yuin/goldmark/ast"
)

// repositoryRoot is relative to this package's test working directory.
const repositoryRoot = "../.."

// invocationPages are the authored Markdown pages whose dotspec invocations
// must name commands and flags the current tree accepts. Generated reference
// pages are rebuilt from the tree, and plans describe past and future surfaces.
func invocationPages(t *testing.T) []string {
	t.Helper()
	var pages []string
	for _, dir := range []string{"docs/cli-handbook", "docs/reference", "docs/maintainers", "specs", "internal/integration/specs", "skills"} {
		pages = append(pages, markdownFiles(t, dir, func(file string) bool {
			return file != "docs/reference/cli.md" && file != "docs/reference/command-examples.md" &&
				!strings.HasPrefix(file, "docs/maintainers/plans/")
		})...)
	}
	for _, pattern := range []string{"examples/*/README.md", "examples/*/COVERAGE.md", "README.md", "CONTRIBUTING.md", "docs/*.md", "model/*.md"} {
		matches, err := filepath.Glob(filepath.Join(repositoryRoot, pattern))
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range matches {
			rel, err := filepath.Rel(repositoryRoot, match)
			if err != nil {
				t.Fatal(err)
			}
			pages = append(pages, filepath.ToSlash(rel))
		}
	}
	return pages
}

func markdownFiles(t *testing.T, dir string, keep func(string) bool) []string {
	t.Helper()
	var files []string
	err := fs.WalkDir(os.DirFS(repositoryRoot), dir, func(file string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && path.Ext(file) == ".md" && keep(file) {
			files = append(files, file)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

type invocation struct {
	source string // repository-relative file and line
	text   string
}

func TestDocumentedInvocationsResolve(t *testing.T) {
	var found []invocation
	for _, page := range invocationPages(t) {
		data, err := os.ReadFile(filepath.Join(repositoryRoot, page))
		if err != nil {
			t.Fatal(err)
		}
		found = append(found, markdownInvocations(page, data)...)
	}
	goFound, err := goInvocations(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	found = append(found, goFound...)
	if len(found) < 50 {
		t.Fatalf("found only %d invocations; the extraction no longer reads the documentation", len(found))
	}
	for _, inv := range found {
		if err := resolveInvocation(inv.text); err != nil {
			t.Errorf("%s: `%s`: %v", inv.source, inv.text, err)
		}
	}
}

func TestCommandReferenceFragmentsResolve(t *testing.T) {
	anchors := map[string]bool{}
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		if cmd.Hidden {
			return
		}
		// The generated reference has one heading per command path.
		anchors[strings.ReplaceAll(cmd.CommandPath(), " ", "-")] = true
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
	pages := append(invocationPages(t), markdownFiles(t, "docs/maintainers/plans", func(string) bool { return true })...)
	checked := 0
	for _, page := range pages {
		data, err := os.ReadFile(filepath.Join(repositoryRoot, page))
		if err != nil {
			t.Fatal(err)
		}
		for _, link := range commandReferenceLinks(page, data) {
			checked++
			if !anchors[link.text] {
				t.Errorf("%s: link to cli.md#%s names no command in the reference", link.source, link.text)
			}
		}
	}
	if checked == 0 {
		t.Fatal("found no command reference fragments; the extraction no longer reads the documentation")
	}
}

func TestInvocationResolution(t *testing.T) {
	for _, tc := range []struct {
		text, wantError string
	}{
		{"dotspec check", ""},
		{"dotspec", ""},
		{"dotspec --version", ""},
		{"dotspec check --all", ""},
		{"dotspec check --json > check-report.json", ""},
		{"dotspec docs <topic>", ""},
		{"dotspec docs success criteria", ""},
		{"dotspec completion bash|zsh|fish|powershell", ""},
		{`dotspec new component "Task service" --part-of specs/a.system.spec.md`, ""},
		{"dotspec check # comment --bogus", ""},
		{"dotspec %s", ""},
		{"dotspec list <kind> --json", ""},
		{"dotspec new <kind> --help", ""},
		{"dotspec snapshot <command> --json", ""},
		{"dotspec -C examples/todo check", ""},
		{"dotspec check --path .", "unknown flag"},
		{"dotspec snapshot <command> --bogus", "no subcommand"},
		{"dotspec list requirment", "unknown command"},
		{"dotspec lst requirement", "unknown command"},
		{"dotspec snapshot exportx", "unknown command"},
		{"dotspec check --bogus", "unknown flag"},
		{"dotspec completion bash|zhs", "unknown command"},
		{"dotspec check a b", "unexpected argument"},
		{"dotspec list job --json -C", "needs an argument"},
	} {
		err := resolveInvocation(tc.text)
		if tc.wantError == "" {
			if err != nil {
				t.Errorf("%s: %v", tc.text, err)
			}
		} else if err == nil || !strings.Contains(err.Error(), tc.wantError) {
			t.Errorf("%s: got %v, want %q", tc.text, err, tc.wantError)
		}
	}
}

func TestInvocationExtraction(t *testing.T) {
	page := "Run `dotspec check` or `dotspec.yaml`.\n\n```sh\n# setup\ndotspec docs job\n$ dotspec sources --json\n```\n\n```text\ndotspec ignored output\n```\n\nOr `mise exec -- go run . check --all`.\n"
	var got []string
	for _, inv := range markdownInvocations("page.md", []byte(page)) {
		got = append(got, inv.source+" "+inv.text)
	}
	want := []string{"page.md:1 dotspec check", "page.md:5 dotspec docs job", "page.md:6 dotspec sources --json", "page.md:13 dotspec check --all"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("markdown: got %q, want %q", got, want)
	}
	got = nil
	for _, value := range []string{
		"no system; start a project with dotspec init <system-name>",
		"obsolete page; run dotspec build in this project",
		"Title\n  dotspec docs %s\n",
		"unknown ID (dotspec docs quality-attribute)",
		"dotspec sources failed: %w",
		"the dotspec repository root",
		"dotspec.yaml",
	} {
		got = append(got, stringInvocations(value)...)
	}
	want = []string{"dotspec init <system-name>", "dotspec build", "dotspec docs %s", "dotspec docs quality-attribute"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("go: got %q, want %q", got, want)
	}
}

// checkoutCommands run this checkout's CLI in place of an installed dotspec.
var checkoutCommands = []string{"mise exec -- go run . ", "go run . "}

// commandLine returns value as a dotspec invocation when it runs the CLI. A
// code span that runs this checkout's sources names a dotspec command, but a
// shell block may first change to another module, such as the example app.
func commandLine(value string, checkout bool) (string, bool) {
	for _, prefix := range checkoutCommands {
		if rest, ok := strings.CutPrefix(value, prefix); ok && checkout {
			value = "dotspec " + rest
			break
		}
	}
	return value, value == "dotspec" || strings.HasPrefix(value, "dotspec ")
}

// markdownInvocations reads code spans that start with the dotspec command or
// run it from this checkout's sources, and command lines in shell code blocks.
func markdownInvocations(page string, data []byte) []invocation {
	md := specfile.NewMarkdown(data)
	offset := len(md.Source) - len(md.Body)
	line := func(at int) string {
		return fmt.Sprintf("%s:%d", page, bytes.Count(md.Source[:offset+at], []byte("\n"))+1)
	}
	var found []invocation
	_ = mdast.Walk(md.Tree, func(node mdast.Node, entering bool) (mdast.WalkStatus, error) {
		if !entering {
			return mdast.WalkContinue, nil
		}
		switch n := node.(type) {
		case *mdast.CodeSpan:
			var text strings.Builder
			start := -1
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				if segment, ok := child.(*mdast.Text); ok {
					if start < 0 {
						start = segment.Segment.Start
					}
					text.Write(segment.Segment.Value(md.Body))
				}
			}
			if value, ok := commandLine(strings.TrimSpace(text.String()), true); ok {
				found = append(found, invocation{line(max(start, 0)), value})
			}
			return mdast.WalkSkipChildren, nil
		case *mdast.FencedCodeBlock:
			switch string(n.Language(md.Body)) {
			case "sh", "shell", "bash", "console", "zsh":
			default:
				return mdast.WalkSkipChildren, nil
			}
			lines := n.Lines()
			for i := range lines.Len() {
				segment := lines.At(i)
				value := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(segment.Value(md.Body))), "$ "))
				if value, ok := commandLine(value, false); ok {
					found = append(found, invocation{line(segment.Start), value})
				}
			}
			return mdast.WalkSkipChildren, nil
		}
		return mdast.WalkContinue, nil
	})
	return found
}

func commandReferenceLinks(page string, data []byte) []invocation {
	md := specfile.NewMarkdown(data)
	var found []invocation
	_ = mdast.Walk(md.Tree, func(node mdast.Node, entering bool) (mdast.WalkStatus, error) {
		link, ok := node.(*mdast.Link)
		if !entering || !ok {
			return mdast.WalkContinue, nil
		}
		u, err := url.Parse(string(link.Destination))
		if err != nil || u.IsAbs() || !strings.HasPrefix(u.Fragment, "dotspec") {
			return mdast.WalkContinue, nil
		}
		if path.Join(path.Dir(page), u.Path) == "docs/reference/cli.md" {
			found = append(found, invocation{page, u.Fragment})
		}
		return mdast.WalkContinue, nil
	})
	return found
}

// goInvocations reads dotspec commands that non-test Go source tells a reader
// to run, such as diagnostic hints and help text.
func goInvocations(root string) ([]invocation, error) {
	var found []invocation
	err := filepath.WalkDir(root, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := entry.Name()
		if entry.IsDir() {
			if file != root && (strings.HasPrefix(name, ".") || name == "dist" || name == "agent_extensions" || name == "testdata") {
				return fs.SkipDir
			}
			return nil
		}
		if filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		parsed, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		if ast.IsGenerated(parsed) {
			return nil
		}
		rel, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			lit, ok := node.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(lit.Value)
			if err != nil {
				return true
			}
			for _, text := range stringInvocations(value) {
				found = append(found, invocation{fmt.Sprintf("%s:%d", filepath.ToSlash(rel), fset.Position(lit.Pos()).Line), text})
			}
			return true
		})
		return nil
	})
	return found, err
}

// A Go string names a command to run after a cue such as "run" or "use",
// inside backticks or parentheses, after a label colon, or on an indented line.
// Prose such as "the dotspec repository" is not an invocation.
var goInvocationStart = regexp.MustCompile("(?m)(?:\\b(?:[Rr]un|[Uu]se|[Ss]ee|with) |: |[`(]|^[ \\t]+)dotspec ")

// Prose that follows an invocation in a sentence ends it.
var goInvocationStops = map[string]bool{"to": true, "or": true, "and": true, "in": true, "then": true, "for": true, "from": true}

func stringInvocations(value string) []string {
	var found []string
	for _, match := range goInvocationStart.FindAllStringIndex(value, -1) {
		rest := value[match[1]:]
		if end := strings.IndexAny(rest, "\n`"); end >= 0 {
			rest = rest[:end]
		}
		words := []string{"dotspec"}
		for _, word := range strings.Fields(rest) {
			trimmed := strings.TrimRight(word, ".,;:)")
			if goInvocationStops[trimmed] {
				break
			}
			if trimmed != "" {
				words = append(words, trimmed)
			}
			if trimmed != word {
				break
			}
		}
		found = append(found, strings.Join(words, " "))
	}
	return found
}

// resolveInvocation resolves a written invocation against the command tree:
// the command path, each flag and its value, and any positional arguments
// written. A reference that omits required arguments still names a command.
func resolveInvocation(text string) error {
	words, err := shellWords(text)
	if err != nil {
		return err
	}
	if len(words) == 0 || words[0] != "dotspec" {
		return errors.New("not a dotspec invocation")
	}
	// Expand alternatives written as a|b into one invocation per choice.
	for i, word := range words {
		if strings.Contains(word, "|") && !strings.HasPrefix(word, "-") {
			var errs []error
			for _, choice := range strings.Split(word, "|") {
				expanded := append(append(append([]string{}, words[:i]...), choice), words[i+1:]...)
				errs = append(errs, resolveArgs(expanded[1:]))
			}
			return errors.Join(errs...)
		}
	}
	return resolveArgs(words[1:])
}

func resolveCommand(args []string) error {
	cmd, rest, err := NewRoot().Find(args)
	if err != nil {
		return err
	}
	cmd.InitDefaultHelpFlag()
	var written []string
	for _, arg := range rest {
		// An optional placeholder such as [<path>] stands for nothing written.
		if !strings.HasPrefix(arg, "[") {
			written = append(written, arg)
		}
	}
	if err := cmd.ParseFlags(written); err != nil {
		// A misspelled subcommand explains the unknown flag that follows it.
		for _, arg := range written {
			if cmd.HasSubCommands() && !strings.HasPrefix(arg, "-") {
				if argErr := cmd.ValidateArgs([]string{arg}); argErr != nil {
					return argErr
				}
				break
			}
		}
		return err
	}
	positional := cmd.Flags().Args()
	if len(positional) == 0 {
		return nil
	}
	return cmd.ValidateArgs(positional)
}

// resolveArgs resolves args, letting a placeholder where a subcommand may
// stand, as in "dotspec new <kind> --help", be any subcommand that accepts
// the rest of the invocation.
func resolveArgs(args []string) error {
	cmd, rest, err := NewRoot().Find(args)
	if err != nil {
		return err
	}
	if cmd.HasSubCommands() {
		for i, arg := range rest {
			if strings.HasPrefix(arg, "-") {
				continue
			}
			if !placeholder(arg) {
				break
			}
			path := strings.Fields(cmd.CommandPath())[1:]
			var first error
			for _, child := range cmd.Commands() {
				if !child.IsAvailableCommand() {
					continue
				}
				expanded := append(append(append([]string{}, path...), rest[:i]...), child.Name())
				expanded = append(expanded, rest[i+1:]...)
				if err := resolveCommand(expanded); err == nil {
					return nil
				} else if first == nil {
					first = err
				}
			}
			if first == nil {
				return nil
			}
			return fmt.Errorf("no subcommand of %q accepts %q: %w", cmd.CommandPath(), strings.Join(rest, " "), first)
		}
	}
	return resolveCommand(args)
}

func placeholder(word string) bool {
	return strings.HasPrefix(word, "<") || strings.HasPrefix(word, "%") || strings.HasPrefix(word, "$")
}

// shellWords splits a command line as a POSIX shell would for these simple
// cases, ending at a comment or a redirection, pipe or command separator.
func shellWords(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	inWord, quote := false, rune(0)
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				word.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			if inWord {
				words = append(words, word.String())
				word.Reset()
				inWord = false
			}
		case !inWord && (r == '#' || r == '>' || r == ';' || r == '&' || r == '|' || r == '<' && isRedirect(string(runes[i:]))):
			return words, nil
		default:
			word.WriteRune(r)
			inWord = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote in %q", line)
	}
	if inWord {
		words = append(words, word.String())
	}
	return words, nil
}

// isRedirect distinguishes "< file" from a placeholder such as <topic>.
func isRedirect(rest string) bool {
	end := strings.IndexRune(rest, '>')
	return end < 0 || strings.ContainsAny(rest[:end], " \t")
}
