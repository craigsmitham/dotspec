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
