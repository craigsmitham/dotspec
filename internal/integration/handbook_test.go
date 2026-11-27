package integration_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/docs"
	"github.com/dotspec-md/dotspec/examples"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/yuin/goldmark/ast"
)

// Read the author's links independently of the offline and repository footers.
// Otherwise removing a bundled file could silently remove it from test coverage.
func TestHandbookPublication(t *testing.T) {
	repository := os.DirFS("../..")
	open := offlineReader(t)
	pages := []string{docs.Root + "/README.md"}
	err := fs.WalkDir(repository, docs.Root, func(file string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && path.Dir(file) != docs.Root && path.Ext(file) == ".md" {
			pages = append(pages, file)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	// Include published walkthroughs alongside handbook pages. The specimens
	// and generated specification they cite are source artifacts.
	err = fs.WalkDir(examples.Files, ".", func(file string, entry fs.DirEntry, err error) error {
		if entry != nil && entry.IsDir() && entry.Name() == "specification" {
			return fs.SkipDir
		}
		if err == nil && !entry.IsDir() && path.Ext(file) == ".md" &&
			path.Base(file) != "SPEC.md" && !strings.HasSuffix(file, ".spec.md") {
			pages = append(pages, "examples/"+file)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range pages {
		t.Run(file, func(t *testing.T) {
			data, err := fs.ReadFile(repository, file)
			if err != nil {
				t.Fatal(err)
			}
			if err := open(file); err != nil {
				t.Fatalf("%s: %v", file, err)
			}
			if err := checkHandbookLinks(repository, file, data, open); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestMain installs the dotspec entry point on PATH. Start every invocation
// in the same empty directory and check each distinct destination once.
func offlineReader(t *testing.T) func(string) error {
	t.Helper()
	dir := t.TempDir()
	results := make(map[string]error)
	return func(address string) error {
		if err, ok := results[address]; ok {
			return err
		}
		cmd := exec.Command("dotspec", "docs", address)
		cmd.Dir = dir
		cmd.Stdout = io.Discard
		var stderr strings.Builder
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			err = fmt.Errorf("dotspec docs %s: %w: %s", address, err, strings.TrimSpace(stderr.String()))
		}
		results[address] = err
		return err
	}
}

func checkHandbookLinks(repository fs.FS, source string, data []byte, open func(string) error) error {
	md := specfile.NewMarkdown(data)
	var findings []error
	err := ast.Walk(md.Tree, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		link, ok := node.(*ast.Link)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		err := checkHandbookLink(repository, source, link, md.Body, open)
		if err != nil {
			findings = append(findings, fmt.Errorf("%s: link %q: %w", source, link.Destination, err))
		}
		return ast.WalkContinue, nil
	})
	return errors.Join(append(findings, err)...)
}

func checkHandbookLink(repository fs.FS, source string, link *ast.Link, body []byte, open func(string) error) error {
	u, err := url.Parse(string(link.Destination))
	if err != nil {
		return err
	}
	if u.Scheme == "file" || path.IsAbs(u.Path) && u.Host == "" {
		return errors.New("use a relative path for local reading")
	}
	if u.IsAbs() || u.Host != "" {
		return nil
	}
	file := source
	if u.Path != "" {
		file = path.Join(path.Dir(source), u.Path)
	}
	target := file
	if u.Fragment != "" {
		target += "#" + u.Fragment
	}
	// The exception is visible to the reader and belongs to this one citation.
	last, ok := link.LastChild().(*ast.Text)
	if ok && strings.HasSuffix(string(last.Value(body)), " (dotspec checkout)") {
		data, err := fs.ReadFile(repository, file)
		if err != nil {
			return err
		}
		if u.Fragment != "" && !specfile.NewMarkdown(data).AnchorReferences().Anchors[u.Fragment] {
			return fmt.Errorf("checkout target %s has no heading %q", file, u.Fragment)
		}
		return nil
	}
	// The full address can only succeed by exact lookup, even when a
	// missing page’s basename happens to be a topic or search term.
	return open(target)
}

func TestHandbookPublicationRegressions(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "study.md"), []byte("# Research note\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	open := offlineReader(t)
	const source = "walkthrough.md"
	for _, tc := range []struct{ link, wantError string }{
		{"[required study](study.md)", "cannot open"},
		{"[removed heading](docs/cli-handbook/guides/write-a-job.md#former-job-heading)", "unknown section"},
		{"[study protocol (dotspec checkout)](study.md#research-note)", ""},
		{"[study protocol (dotspec checkout)](missing.md)", "missing.md"},
		{"[study protocol (dotspec checkout)](study.md#former-heading)", "no heading"},
		{"`[a specimen path](missing.md)`", ""},
	} {
		err := checkHandbookLinks(os.DirFS(root), source, []byte(tc.link), open)
		if tc.wantError == "" {
			if err != nil {
				t.Errorf("%s: %v", tc.link, err)
			}
		} else if err == nil || !strings.Contains(err.Error(), source) || !strings.Contains(err.Error(), tc.wantError) {
			t.Errorf("%s: got %v, want source and %q", tc.link, err, tc.wantError)
		}
	}
}
