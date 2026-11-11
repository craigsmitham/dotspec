package cli

import (
	"encoding/json"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"os"
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
		"dotspec.yaml":                          "projects:\n  tool: .\n  demo: apps/demo\n",
		"specs/system.spec.md":                  "---\ntype: system\n---\n# Tool\n\n## Purpose\nTool work.\n\n## Scope\nTool only.\n",
		"apps/demo/specs/system.spec.md":        "---\ntype: system\n---\n# Demo\n\n## Purpose\nDemo work.\n\n## Scope\nDemo only.\n",
		"apps/demo/specs/jobs/work.job.spec.md": "---\ntype: job\n---\n# Work\n\n## Statement\nWhen I need help, I want to work, so I make progress.\n",
	}
	for path, body := range files {
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
	return root
}
func TestProjectSelectorsAndSeparateWorkspaceResults(t *testing.T) {
	root := projectFixture(t)
	for _, args := range [][]string{{"system", "list", "--json"}, {"system", "list", "--project", "demo", "--json"}} {
		code, out, errOut := invoke(args...)
		if code != 0 {
			t.Fatal(code, errOut)
		}
		want, absent := "Tool", "Demo"
		if len(args) > 3 {
			want, absent = absent, want
		}
		if !strings.Contains(out, want) || strings.Contains(out, absent) {
			t.Fatal(out)
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
	code, out, errOut = invoke("system", "show", "--project", "demo", "--json")
	if code != 0 || !strings.Contains(out, `"enclosing_system": null`) {
		t.Fatal(out, errOut)
	}
	code, out, errOut = invoke("model", "export", "--project", "demo")
	if code != 0 || strings.Contains(out, `"enclosing_system"`) || !strings.Contains(out, `"workspace_configuration"`) {
		t.Fatal(code, errOut, out)
	}
	// A missing System never falls through to the tool's System.
	if err := os.Remove(filepath.Join(root, "apps/demo/specs/system.spec.md")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("job", "show", "work", "--project", "demo", "--json")
	if code != 0 || strings.Contains(out, `"name": "Tool"`) {
		t.Fatal(code, out, errOut)
	}
	for _, args := range [][]string{{"check", "--project", "missing"}, {"check", "--project", "tool", "--path", "apps/demo"}, {"check", "--all", "--project", "demo"}} {
		code, _, _ := invoke(args...)
		if code != 2 {
			t.Fatal(args, code)
		}
	}
}
func TestExplicitCrossProjectContextAndAllFailures(t *testing.T) {
	root := projectFixture(t)
	path := filepath.Join(root, "apps/demo/specs/jobs/work.job.spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte("\n## Specification system\n\n[Tool](../../../../specs/system.spec.md).\n")...)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := invoke("model", "export", "--project", "demo")
	if code != 0 || !strings.Contains(out, `"specification_system_origin": "explicit"`) || !strings.Contains(out, `"name": "tool"`) {
		t.Fatal(code, errOut, out)
	}
	if err := os.WriteFile(filepath.Join(root, "specs/system.spec.md"), []byte("not a system"), 0600); err != nil {
		t.Fatal(err)
	}
	code, out, errOut = invoke("check", "--all", "--json")
	var result workspaceResult
	if code != 1 || json.Unmarshal([]byte(out), &result) != nil || len(result.Projects) != 2 {
		t.Fatal(code, out, errOut)
	}
}

func TestWorkspaceExportsValidateOffline(t *testing.T) {
	root := projectFixture(t)
	code, out, errOut := invoke("model", "export", "--all")
	if code != 0 {
		t.Fatal(code, errOut)
	}
	schema, err := workspaceExportSchema()
	if err != nil {
		t.Fatal(err)
	}
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	var schemaValue any
	if err := json.Unmarshal(schemaBytes, &schemaValue); err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("https://example.test/workspace.json", schemaValue); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile("https://example.test/workspace.json")
	if err != nil {
		t.Fatal(err)
	}
	var value any
	if err := json.Unmarshal([]byte(out), &value); err != nil {
		t.Fatal(err)
	}
	if err := compiled.Validate(value); err != nil {
		t.Fatal(err)
	}
	findings, err := validateExportJSON([]byte(out))
	if err != nil || len(findings) != 0 {
		t.Fatal(findings, err)
	}
	// Acquisition must have recorded the manifest, project provenance and separate snapshots.
	var envelope workspaceResult
	if err := json.Unmarshal([]byte(out), &envelope); err != nil {
		t.Fatal(err)
	}
	if len(envelope.Projects) != 2 {
		t.Fatal(envelope)
	}
	path := filepath.Join(root, "export.json")
	if err := os.WriteFile(path, []byte(out), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "dotspec.yaml")); err != nil {
		t.Fatal(err)
	}
	code, _, errOut = invoke("model", "validate", path)
	if code != 0 {
		t.Fatal(code, errOut)
	}
	envelope.Projects[0].Name = "forged"
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	findings, err = validateExportJSON(data)
	if err != nil || len(findings) == 0 || findings[0].Code != "workspace.scope" {
		t.Fatal(findings, err)
	}
}

func TestWorkspaceRootWithoutProjectAndPathNarrowing(t *testing.T) {
	root := projectFixture(t)
	code, out, errOut := invoke("check", "--project", "demo", "--path", "internal", "--json")
	if code != 0 || !strings.Contains(out, `"files_discovered": 0`) {
		t.Fatal(code, out, errOut)
	}
	t.Chdir(filepath.Join(root, "apps/demo/internal"))
	code, out, errOut = invoke("system", "list", "--json")
	if code != 0 || !strings.Contains(out, "Demo") {
		t.Fatal(code, out, errOut)
	}
	t.Chdir(root)
	if err := os.WriteFile(filepath.Join(root, "dotspec.yaml"), []byte("projects:\n  demo: apps/demo\n"), 0600); err != nil {
		t.Fatal(err)
	}
	code, _, _ = invoke("check")
	if code != 2 {
		t.Fatal(code)
	}
	code, _, errOut = invoke("check", "--all")
	if code != 0 {
		t.Fatal(code, errOut)
	}
}

func TestWorkspaceMutationCannotTargetAnotherProject(t *testing.T) {
	root := projectFixture(t)
	target := filepath.Join(root, "apps/demo/specs/system.spec.md")
	before, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	code, _, _ := invoke("quality-model", "apply", "--project", "tool", "--system", "apps/demo/specs/system.spec.md")
	if code != 2 {
		t.Fatal(code)
	}
	after, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("cross-project target changed")
	}
}
