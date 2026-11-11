package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/model"
)

func TestModelExportValidationAndSourcePreservation(t *testing.T) {
	t.Chdir(t.TempDir())
	system := "---\ntype: system\n---\n# System\n## Purpose\nRead specs.\n## Scope\nA CLI.\n"
	useCase := "\xef\xbb\xbf---\r\ntype: use-case\r\nid: check\r\n---\r\n# Check a specification\r\n\r\n## Goal\r\nRead [missing material](../../missing.md).  \r\n\r\n## Main success scenario\r\n1. Ask for a report.\r\n2. Receive findings.\r\n\r\n## Extensions\r\n### Invalid source\r\nHandling unresolved.\r\n\r\n<!-- retain this comment -->\r\n"
	writeFixture(t, "specs/system.spec.md", system)
	writeFixture(t, "specs/use-cases/check.use-case.spec.md", useCase)
	code, output, errOut := invoke("model", "export")
	if code != 1 || errOut != "" {
		t.Fatalf("export with source finding: %d %s", code, errOut)
	}
	var snapshot model.Snapshot
	if err := json.Unmarshal([]byte(output), &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Elements) < 5 || snapshot.ModelVersion != model.ModelVersion {
		t.Fatalf("missing typed domain records: %d", len(snapshot.Elements))
	}
	for path, expected := range map[string]string{"specs/system.spec.md": system, "specs/use-cases/check.use-case.spec.md": useCase} {
		actual, err := os.ReadFile(path)
		if err != nil || string(actual) != expected {
			t.Fatalf("source modified: %s: %v", path, err)
		}
	}
	if err := os.WriteFile("snapshot.json", []byte(output), 0o600); err != nil {
		t.Fatal(err)
	}
	code, validation, errOut := invoke("model", "validate", "snapshot.json", "--json")
	if code != 0 || errOut != "" || !strings.Contains(validation, `"valid": true`) {
		t.Fatalf("faithful source-error snapshot rejected: %d %s %s", code, validation, errOut)
	}
	cmd := NewRoot()
	cmd.SetArgs([]string{"model", "validate", "-", "--json"})
	cmd.SetIn(strings.NewReader(output))
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := cmd.Execute(); err != nil || !strings.Contains(stdout.String(), `"valid": true`) {
		t.Fatalf("stdin validation: %v %s %s", err, &stdout, &stderr)
	}
}

func TestModelSchemaAndValidationFailures(t *testing.T) {
	t.Chdir(t.TempDir())
	code, output, errOut := invoke("model", "schema")
	if code != 0 || errOut != "" || !json.Valid([]byte(output)) || !strings.Contains(output, `"ScenarioStep"`) {
		t.Fatalf("schema: %d %s", code, errOut)
	}
	for _, content := range []string{`{}`, `{"model_version":`, `{"model_version":"a","model_version":"b"}`} {
		path := filepath.Join(t.TempDir(), "snapshot.json")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		code, output, errOut := invoke("model", "validate", path, "--json")
		if code != 1 || errOut != "" || !strings.Contains(output, `"valid": false`) {
			t.Fatalf("invalid snapshot %q: %d %s %s", content, code, output, errOut)
		}
	}
	for _, args := range [][]string{{"model", "validate", "absent.json"}, {"model", "export", "--path", "absent"}, {"model", "validate"}, {"model", "schema", "unexpected"}} {
		code, output, errOut := invoke(args...)
		if code != 2 || output != "" || errOut == "" {
			t.Fatalf("invocation/read failure %v: %d %s %s", args, code, output, errOut)
		}
	}
}

type modelBrokenWriter struct{}

func (modelBrokenWriter) Write([]byte) (int, error) { return 0, errors.New("broken output") }

func TestModelOutputFailure(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, args := range [][]string{{"model", "export"}, {"model", "schema"}} {
		var stderr bytes.Buffer
		if code := Run(args, modelBrokenWriter{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "broken output") {
			t.Fatalf("%v: %d %s", args, code, &stderr)
		}
	}
}

func TestConceptDesignCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/concept-design")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	code, output, errOut := invoke("model", "docs", "--relative-to", "docs")
	if code != 0 || errOut != "" {
		t.Fatalf("docs %d: %s", code, errOut)
	}
	for _, part := range []string{"## Concepts", "#### Coordination", "#### Application dependencies", "Resource = PublicLinks.Code", "../specs/concepts/shorthand.concept.spec.md"} {
		if !strings.Contains(output, part) {
			t.Fatalf("missing %q", part)
		}
	}
	code, output, errOut = invoke("system", "show", "--json")
	if code != 0 || errOut != "" || !strings.Contains(output, `"concept_uses":`) || !strings.Contains(output, `"synchronizations":`) {
		t.Fatalf("system composition omitted: %d %s %s", code, output, errOut)
	}
	path := filepath.Join(root, "specs/system.spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "[Register](concepts/shorthand.concept.spec.md#register)", "[Expire](concepts/expiry.concept.spec.md#expire)", 1))
	if err = os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	code, output, errOut = invoke("model", "docs")
	if code != 1 || errOut != "" || !strings.Contains(output, "concept.action") || !strings.Contains(output, "#### Synchronization rules") {
		t.Fatalf("source findings not retained: %d %s %s", code, output, errOut)
	}
	code, output, errOut = invoke("model", "export")
	if code != 1 || errOut != "" {
		t.Fatalf("export: %d %s", code, errOut)
	}
	if err = os.WriteFile("snapshot.json", []byte(output), 0600); err != nil {
		t.Fatal(err)
	}
	code, output, errOut = invoke("model", "validate", "snapshot.json", "--json")
	if code != 0 || errOut != "" || !strings.Contains(output, `"valid": true`) {
		t.Fatalf("faithful composition findings rejected: %d %s %s", code, output, errOut)
	}
}

func TestRequirementsTraceCommand(t *testing.T) {
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/requirements-traceability")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	code, output, errOut := invoke("model", "trace", "--relative-to", "docs")
	if code != 0 || errOut != "" || !strings.Contains(output, "## Change review") || !strings.Contains(output, "../specs/use-cases/restore.use-case.spec.md#restore-retained-content") {
		t.Fatalf("trace: %d %s %s", code, output, errOut)
	}
	var stderr bytes.Buffer
	if code := Run([]string{"model", "trace"}, modelBrokenWriter{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "broken output") {
		t.Fatalf("trace output failure: %d %s", code, &stderr)
	}
	path := filepath.Join(root, "specs/requirements/restoration.requirement.spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "[Trusted time](../system.spec.md#trusted-time)", "[Goal](../system.spec.md#reduce-lost-work)", 1))
	if err = os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	code, output, errOut = invoke("model", "trace")
	if code != 1 || errOut != "" || !strings.Contains(output, "relationship.relies_on") {
		t.Fatalf("trace findings: %d %s %s", code, output, errOut)
	}
}

func TestSystemOverviewCommand(t *testing.T) {
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/system-structure")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	code, output, errOut := invoke("model", "overview", "--relative-to", "docs")
	if code != 0 || errOut != "" || !strings.Contains(output, "# System specification overview") || !strings.Contains(output, "../specs/system.spec.md#checkout") || !strings.Contains(output, "## Features") {
		t.Fatalf("overview: %d %s %s", code, output, errOut)
	}
	var stderr bytes.Buffer
	if code := Run([]string{"model", "overview"}, modelBrokenWriter{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "broken output") {
		t.Fatal(code, &stderr)
	}
}

func TestModelReviewAndPEGSCommandsKeepAdvisoryDraftsValid(t *testing.T) {
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/requirements-review")); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	for _, tc := range []struct{ command, want string }{{"review", "No system verification or measurement guidance recorded"}, {"pegs", "## Project"}} {
		code, output, errOut := invoke("model", tc.command, "--relative-to", "docs")
		if code != 0 || errOut != "" || !strings.Contains(output, tc.want) || !strings.Contains(output, "../specs/system.spec.md") {
			t.Fatalf("%s: %d %s %s", tc.command, code, errOut, output)
		}
	}
}
