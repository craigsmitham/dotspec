package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func invoke(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestSystemWorkflow(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	code, out, errOut := invoke("system", "init")
	if code != 0 || !strings.Contains(out, "Created system.spec.md") || errOut != "" {
		t.Fatalf("init: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("system", "show")
	for _, want := range []string{"New system", "Purpose", "Scope", "Parent system: None", "Source: system.spec.md"} {
		if code != 0 || !strings.Contains(out, want) || errOut != "" {
			t.Fatalf("show missing %q: %d %q %q", want, code, out, errOut)
		}
	}
	child := filepath.Join(dir, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := invoke("system", "init", "--path", child); code != 0 {
		t.Fatal(errOut)
	}
	code, out, errOut = invoke("system", "list")
	want := "Systems\n└── New system  (system.spec.md)\n    └── New system  (child/system.spec.md)\n"
	if code != 0 || out != want || errOut != "" {
		t.Fatalf("list: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("system", "show", "--path", child)
	if code != 0 || !strings.Contains(out, "Parent system: New system (system.spec.md)") || errOut != "" {
		t.Fatalf("child show: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("system", "init")
	if code == 0 || out != "" || !strings.Contains(errOut, "refusing to overwrite") {
		t.Fatalf("overwrite: %d %q %q", code, out, errOut)
	}
}

func TestHelpAndUsageErrors(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"system"}, {"system", "--help"}, {"system", "init", "--help"}, {"system", "show", "--help"}, {"system", "list", "--help"}} {
		code, out, errOut := invoke(args...)
		if code != 0 || !strings.Contains(out, "Usage:") || errOut != "" {
			t.Errorf("help %v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{{"bogus"}, {"system", "bogus"}, {"system", "show", "unexpected"}, {"system", "list", "--unknown"}, {"system", "show", "--path"}, {"system", "show", "--path", ""}} {
		code, out, errOut := invoke(args...)
		if code == 0 || out != "" || !strings.HasPrefix(errOut, "dotspec:") {
			t.Errorf("usage error %v: %d %q %q", args, code, out, errOut)
		}
	}
}

func TestMissingAndMalformedSystems(t *testing.T) {
	dir := t.TempDir()
	// A .git marker bounds ancestor lookup without requiring Git for show.
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := invoke("system", "show", "--path", dir)
	if code == 0 || out != "" || !strings.Contains(errOut, "no system spec") {
		t.Fatalf("missing: %d %q %q", code, out, errOut)
	}
	if err := os.WriteFile(filepath.Join(dir, "system.spec.md"), []byte("broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("system", "show", "--path", dir)
	if code == 0 || out != "" || !strings.Contains(errOut, "system.spec.md") {
		t.Fatalf("malformed: %d %q %q", code, out, errOut)
	}
}

func TestEmptyList(t *testing.T) {
	code, out, errOut := invoke("system", "list", "--path", t.TempDir())
	if code != 0 || !strings.Contains(out, "No systems found") || errOut != "" {
		t.Fatalf("empty list: %d %q %q", code, out, errOut)
	}
}
