package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/cli"
	"github.com/craigsmitham/dotspec/internal/concept"
)

func put(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func cliSources(root string) ([]byte, error) {
	var out, stderr bytes.Buffer
	code := cli.Run([]string{"sources", "--path", root, "--json"}, &out, &stderr)
	if code != 0 {
		return nil, fmt.Errorf("sources exit %d: %s %s", code, &out, &stderr)
	}
	return out.Bytes(), nil
}

func setup(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	put(t, dir, "internal/integration/testdata/workflow.txtar", "# An acceptance example.\n")
	put(t, dir, "internal/cli/quality_contracts_test.go", "package cli\n")
	return dir
}

func TestCoverageUsesOnlyDesignatedFilesAndPublicInventory(t *testing.T) {
	dir := setup(t)
	spec := concept.Template + "\n[Example](../../internal/integration/testdata/workflow.txtar)\n[Verification](../../internal/cli/quality_contracts_test.go)\n"
	put(t, dir, "specs/concepts/cli.concept.spec.md", spec)
	if count, err := run(dir, cliSources); err != nil || count != 2 {
		t.Fatalf("coverage: %d %v", count, err)
	}
	put(t, dir, "internal/cli/private_test.go", "package cli\n// Implementation detail.\n")
	if count, err := run(dir, cliSources); err != nil || count != 2 {
		t.Fatalf("implementation test opted in: %d %v", count, err)
	}
	put(t, dir, "internal/integration/testdata/new.txtar", "# New acceptance example.\n")
	if _, err := run(dir, cliSources); err == nil || !strings.Contains(err.Error(), "testdata/new.txtar") {
		t.Fatalf("missed unreferenced source: %v", err)
	}
	spec += "[New example](../../internal/integration/testdata/new.txtar)\n"
	put(t, dir, "specs/concepts/cli.concept.spec.md", spec)
	if count, err := run(dir, cliSources); err != nil || count != 3 {
		t.Fatalf("reference did not repair coverage: %d %v", count, err)
	}
	if data, err := os.ReadFile(filepath.Join(dir, "specs/concepts/cli.concept.spec.md")); err != nil || string(data) != spec {
		t.Fatal("coverage rewrote source", err)
	}
	put(t, dir, "specs/concepts/cli.concept.spec.md", spec+"\n[Missing](missing.tla)\n")
	if _, err := run(dir, cliSources); err == nil || !strings.Contains(err.Error(), "sources exit 1") {
		t.Fatalf("accepted failed inventory: %v", err)
	}
}

func TestCoverageRejectsUnusableInventory(t *testing.T) {
	dir := setup(t)
	for _, data := range []string{
		`not json`,
		`{"schema_version":1,"complete":false}`,
		`{"schema_version":2,"complete":true}`,
		`{"schema_version":1,"complete":true,"scope":{"path":"elsewhere"}}`,
	} {
		if _, err := run(dir, func(string) ([]byte, error) { return []byte(data), nil }); err == nil {
			t.Fatalf("accepted %s", data)
		}
	}
	if _, err := run(dir, func(string) ([]byte, error) { return nil, fmt.Errorf("inventory unavailable") }); err == nil {
		t.Fatal("ignored command failure")
	}
}
