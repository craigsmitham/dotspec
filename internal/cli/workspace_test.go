package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func projectFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"specs", "apps/demo/specs/jobs", "apps/demo/internal"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{
		"dotspec.yaml":                          "name: tool\n",
		"apps/demo/dotspec.yaml":                "name: demo\n",
		"specs/tool.system.spec.md":             "---\ntype: system-spec\n---\n# Tool\n\n## Purpose\nTool work.\n\n## Scope\nTool only.\n",
		"apps/demo/specs/demo.system.spec.md":   "---\ntype: system-spec\n---\n# Demo\n\n## Purpose\nDemo work.\n\n## Scope\nDemo only.\n",
		"apps/demo/specs/jobs/work.job.spec.md": "---\ntype: job-spec\n---\n# Work\n\n## Job\nWhen I need help, I want to work, so I make progress.\n",
	}
	for path, body := range files {
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
	return root
}
func TestChdirSelectsProjectAndSeparateWorkspaceResults(t *testing.T) {
	root := projectFixture(t)
	for _, args := range [][]string{{"list", "system", "--json"}, {"-C", "apps/demo", "list", "system", "--json"}} {
		code, out, errOut := invoke(args...)
		if code != 0 {
			t.Fatal(code, errOut)
		}
		want, absent := "Tool", "Demo"
		if args[0] == "-C" {
			want, absent = absent, want
		}
		if !strings.Contains(out, want) || strings.Contains(out, absent) {
			t.Fatal(out)
		}
	}
	// -C lasts for one invocation.
	if cwd, err := os.Getwd(); err != nil || mustEvalSymlinks(t, cwd) != mustEvalSymlinks(t, root) {
		t.Fatal(cwd, err)
	}
	for _, dir := range []string{".", "apps/demo"} {
		if code, out, errOut := invoke("--chdir", dir, "build"); code != 0 {
			t.Fatal(code, out, errOut)
		}
	}
	code, out, errOut := invoke("check", "--all", "--json")
	var result workspaceResult
	if code != 0 || json.Unmarshal([]byte(out), &result) != nil || len(result.Projects) != 2 {
		t.Fatal(code, out, errOut)
	}
	if result.Projects[0].Name != "demo" || result.Projects[1].Name != "tool" {
		t.Fatal(result)
	}
	code, out, errOut = invoke("-C", "apps/demo", "snapshot", "export")
	if code != 0 || strings.Contains(out, `"enclosing_system"`) || !strings.Contains(out, `"project_configuration"`) {
		t.Fatal(code, errOut, out)
	}
	// A missing System never falls through to the tool's System.
	if err := os.Remove(filepath.Join(root, "apps/demo/specs/demo.system.spec.md")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("-C", "apps/demo", "list", "job", "--json")
	if code != 0 || strings.Contains(out, `"name": "Tool"`) || !strings.Contains(out, `"system": null`) {
		t.Fatal(code, out, errOut)
	}
	code, _, errOut = invoke("-C", "missing", "check")
	if code != 2 || !strings.Contains(errOut, "--chdir missing: directory does not exist") {
		t.Fatal(code, errOut)
	}
	code, _, errOut = invoke("check", "--all", "apps/demo")
	if code != 2 || !strings.Contains(errOut, "omit <path>") {
		t.Fatal(code, errOut)
	}
}

func mustEvalSymlinks(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func TestExplicitCrossProjectContextAndAllFailures(t *testing.T) {
	root := projectFixture(t)
	path := filepath.Join(root, "apps/demo/specs/jobs/work.job.spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte("\n## Specification system\n\n[Tool](../../../../specs/tool.system.spec.md).\n")...)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := invoke("-C", "apps/demo", "snapshot", "export")
	if code != 0 || !strings.Contains(out, `"role": "specification_system"`) || !strings.Contains(out, `"name": "tool"`) {
		t.Fatal(code, errOut, out)
	}
	if err := os.WriteFile(filepath.Join(root, "specs/tool.system.spec.md"), []byte("not a system"), 0600); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("check", "--all", "--json")
	var result workspaceResult
	if code != 1 || json.Unmarshal([]byte(out), &result) != nil || len(result.Projects) != 2 {
		t.Fatal(code, out, errOut)
	}
	code, out, _ = invoke("check", "--all")
	if code != 1 || !strings.Contains(out, "# Project: demo") || !strings.Contains(out, "# Project: tool") {
		t.Fatal(code, out)
	}
}

func TestWorkspaceRootWithoutProjectAndPathNarrowing(t *testing.T) {
	root := projectFixture(t)
	if output, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatal(err, string(output))
	}
	code, out, errOut := invoke("-C", "apps/demo", "check", "internal", "--json")
	if code != 0 || !strings.Contains(out, `"files_discovered": 0`) {
		t.Fatal(code, out, errOut)
	}
	// Paths resolve against the working directory that -C selects.
	code, out, errOut = invoke("-C", "apps/demo", "check", "specs/jobs/work.job.spec.md", "--json")
	if code != 0 || !strings.Contains(out, `"files_checked": 1`) {
		t.Fatal(code, out, errOut)
	}
	// A named specification file is checked within its own project, not the
	// project enclosing the working directory.
	code, out, errOut = invoke("check", "apps/demo/specs/jobs/work.job.spec.md", "--json")
	if code != 0 || !strings.Contains(out, `"files_checked": 1`) || !strings.Contains(out, `"name": "demo"`) {
		t.Fatal(code, out, errOut)
	}
	t.Run("discover project from a child directory", func(t *testing.T) {
		t.Chdir(filepath.Join(root, "apps/demo/internal"))
		code, out, errOut := invoke("list", "system", "--json")
		if code != 0 || !strings.Contains(out, "Demo") {
			t.Fatal(code, out, errOut)
		}
	})
	if err := os.Remove(filepath.Join(root, "dotspec.yaml")); err != nil {
		t.Fatal(err)
	}
	code, _, errOut = invoke("check")
	if code != 2 || !strings.Contains(errOut, "outside any specification project") {
		t.Fatal(code, errOut)
	}
	if code, out, errOut := invoke("-C", "apps/demo", "build"); code != 0 {
		t.Fatal(code, out, errOut)
	}
	code, out, errOut = invoke("check", "apps/demo", "--json")
	if code != 0 || !strings.Contains(out, `"name": "demo"`) || !strings.Contains(out, `"generated_spec": "checked"`) {
		t.Fatal(code, out, errOut)
	}
	code, _, errOut = invoke("check", "--all")
	if code != 0 {
		t.Fatal(code, errOut)
	}
	if err := os.Rename(filepath.Join(root, "apps/demo"), filepath.Join(root, "apps/moved")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("-C", "apps/moved", "sources", "--json")
	if code != 0 || !strings.Contains(filepath.ToSlash(out), "moved") || strings.Contains(out, "apps/demo") {
		t.Fatal(code, out, errOut)
	}
}

// Adoption writes to the nearest enclosing system of the working directory's
// project, never to a nested project's system.
func TestWorkspaceMutationCannotTargetAnotherProject(t *testing.T) {
	root := projectFixture(t)
	target := filepath.Join(root, "apps/demo/specs/demo.system.spec.md")
	before, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	code, out, errOut := invoke("new", "quality-attribute", "recoverability", "--from", "iso-25010-2023")
	if code != 0 || !strings.Contains(out, "specs/tool.system.spec.md") {
		t.Fatal(code, out, errOut)
	}
	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "apps/demo/specs/quality")); string(before) != string(after) || !os.IsNotExist(err) {
		t.Fatal("adoption changed another project")
	}
}
