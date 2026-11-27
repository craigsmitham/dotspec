package main

import (
	"bytes"
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
		body := "# Test guide\n\n[Related](list-a-companion.md#choose-a-role) · [Map](../README.md)\n\n`[Example](example.md)`\n"
		if err := os.WriteFile(filepath.Join(dir, handbookDir, slug+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err != nil {
		t.Fatalf("fresh generation drifted: %v", err)
	}
	guide, err := os.ReadFile(filepath.Join(dir, "mintlify", "guides", "describe-a-system.mdx"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"/guides/list-a-companion#choose-a-role", githubSource + "docs/cli-handbook/README.md", "`[Example](example.md)`"} {
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
