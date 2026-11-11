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
	for _, fixture := range []string{"concept-design", "requirements-traceability", "system-structure", "user-needs", "requirements-review"} {
		if err := os.CopyFS(filepath.Join(dir, "examples", fixture), os.DirFS(filepath.Join("../../examples", fixture))); err != nil {
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
	path := filepath.Join(dir, "docs", "cli.md")
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	changed := append(append([]byte{}, original...), []byte("\nmanual edit\n")...)
	if err := os.WriteFile(path, changed, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err == nil || !strings.Contains(err.Error(), "docs/cli.md") {
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
