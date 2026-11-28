package main

import (
	"bytes"
	"github.com/dotspec-md/dotspec/internal/cli"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratedArtifactsAndReadOnlyDriftCheck(t *testing.T) {
	dir := t.TempDir()
	if err := os.CopyFS(filepath.Join(dir, "examples"), os.DirFS("../../examples")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, handbookDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docs", "cli-handbook", "README.md"), []byte("# Documentation map\n\n<!-- Generated topic navigation -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, slug := range publishedGuides {
		body := "# Test guide\n\n[Related](catalog-a-specification.md#choose-the-source) · [Map](../README.md)\n\n`[Example](example.md)`\n"
		if err := os.WriteFile(filepath.Join(dir, handbookDir, slug+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.CopyFS(filepath.Join(dir, "mintlify"), os.DirFS("../../mintlify")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err != nil {
		t.Fatalf("fresh generation drifted: %v", err)
	}
	root := cli.NewRoot()
	extra := &cobra.Command{Use: "temporary-command"}
	root.AddCommand(extra)
	assertRoute := func(name string, want bool) {
		t.Helper()
		files, err := webPages(root)
		if err != nil {
			t.Fatal(err)
		}
		page, listed := false, false
		for _, f := range files {
			if f.path == "mintlify/cli/"+name+".mdx" {
				page = true
			}
			if f.path == "mintlify/docs.json" {
				listed = strings.Contains(f.content, `"cli/`+name+`"`)
			}
		}
		if page != want || listed != want {
			t.Fatalf("%s page=%v listed=%v want=%v", name, page, listed, want)
		}
	}
	assertRoute("temporary-command", true)
	extra.Use = "renamed-command"
	assertRoute("temporary-command", false)
	assertRoute("renamed-command", true)
	extra.Hidden = true
	assertRoute("renamed-command", false)
	extra.Hidden = false
	root.RemoveCommand(extra)
	assertRoute("renamed-command", false)
	orphan := "mintlify/cli/retired.mdx"
	if err := os.WriteFile(orphan, []byte("{/* Generated from Cobra command tree by mise run generate. */}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err == nil || !strings.Contains(err.Error(), "retired.mdx") {
		t.Fatalf("missing stale page diagnostic: %v", err)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatal(err)
	}
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("obsolete page survived regeneration")
	}
	guide, err := os.ReadFile(filepath.Join(dir, "mintlify", "guides", "describe-a-system.mdx"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"/guides/catalog-a-specification#choose-the-source", githubSource + "docs/cli-handbook/README.md", "`[Example](example.md)`"} {
		if !strings.Contains(string(guide), expected) {
			t.Fatalf("web guide is missing %q", expected)
		}
	}
	path := filepath.Join(dir, "docs", "reference", "cli.md")
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	changed := append(append([]byte{}, original...), []byte("\nmanual edit\n")...)
	if err := os.WriteFile(path, changed, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err == nil || !strings.Contains(err.Error(), "docs/reference/cli.md") {
		t.Fatalf("missing drift diagnostic: %v", err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, changed) {
		t.Fatal("check mode rewrote the working tree")
	}
}

func TestGalleryNormalizesWindowsPathsWithoutChangingAuthoredEscapes(t *testing.T) {
	input := `{"path":"C:\\work\\project\\specs\\system.spec.md","scope":"C:\\work\\project","markdown":"Literal \\text\n[link](a\\b)","unrelated":"C:\\work\\project-other\\file"}`
	want := `{"path":"$PROJECT/specs/system.spec.md","scope":"$PROJECT","markdown":"Literal \\text\n[link](a\\b)","unrelated":"C:\\work\\project-other\\file"}`
	if got := galleryOutput(input, `C:\work\project`); got != want {
		t.Fatalf("normalized output:\n%s\nwant:\n%s", got, want)
	}
	if got := galleryOutput(`{"path":"/tmp/project/specs/system.spec.md"}`, "/tmp/project"); got != `{"path":"$PROJECT/specs/system.spec.md"}` {
		t.Fatal(got)
	}
}
