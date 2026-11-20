package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/model"
)

func TestInitDoesNotWriteThroughMetadataSymlinks(t *testing.T) {
	for _, subdir := range []string{"specs", "specs/jobs", "specs/concepts", "specs/requirements"} {
		t.Run(subdir, func(t *testing.T) {
			dir, outside := t.TempDir(), t.TempDir()
			link := filepath.Join(dir, subdir)
			if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Fatal(err)
			}
			kind := "job"
			if subdir == "specs/requirements" {
				kind = "requirement"
			}
			if subdir == "specs/concepts" {
				kind = "concept"
			}
			if code, out, errOut := invoke(kind, "init", "work", "--path", dir); code == 0 || out != "" || !strings.Contains(errOut, "symlink") {
				t.Fatalf("symlink: %d %q %q", code, out, errOut)
			}
			entries, err := os.ReadDir(outside)
			if err != nil || len(entries) != 0 {
				t.Fatalf("init wrote outside the scope: %v, %v", entries, err)
			}
		})
	}
}

func TestMoveRequirementAcrossSourceLayouts(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: movement\ngenerated_spec: true\n")
	writeFixture(t, "product.system.spec.md", "---\ntype: system-spec\n---\n# Product\n## Purpose\nStore work.\n## Scope\nThe product.\n")
	writeFixture(t, "app/specs/service.component.spec.md", "---\ntype: component-spec\n---\n# Service\n## Part of\n[Product](../../product.system.spec.md)\n")
	writeFixture(t, "app/specs/storage/storage.component.spec.md", "---\ntype: component-spec\n---\n# Storage\n## Part of\n[Service](../service.component.spec.md)\n")
	writeFixture(t, "qualities/reliable.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\n---\n# Reliable\n")
	writeFixture(t, "qualities/durable.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nsubattribute_of: [reliable]\n---\n# Durable\n")
	const requirement = "---\ntype: requirement-spec\nquality_attributes: [durable]\n---\n# Save\n## Statement\nPersist each accepted edit.\n## Applies to\n**Targets:** [Product](%sproduct.system.spec.md).\n"
	run := func(args ...string) string {
		t.Helper()
		code, out, stderr := invoke(args...)
		if code != 0 {
			t.Fatalf("%v: exit %d\n%s\n%s", args, code, out, stderr)
		}
		return out
	}
	var previous string
	for _, stage := range []struct{ file, up, owner, report string }{
		{"app/specs/requirements/save.requirement.spec.md", "../../../", "Service", "components/service"},
		{"app/handlers/save.requirement.spec.md", "../../", "Service", "components/service"},
		{"app/specs/storage/save.requirement.spec.md", "../../../", "Storage", "components/service/components/storage"},
	} {
		if previous != "" {
			if err := os.Remove(previous); err != nil {
				t.Fatal(err)
			}
		}
		writeFixture(t, stage.file, fmt.Sprintf(requirement, stage.up))
		previous = stage.file
		run("check", stage.file)
		out := run("requirement", "show", "save", "--path", filepath.Dir(stage.file))
		if !strings.Contains(out, "Ancestor quality attributes: reliable") {
			t.Fatal(out)
		}
		var snapshot model.Snapshot
		if err := json.Unmarshal([]byte(run("model", "export")), &snapshot); err != nil {
			t.Fatal(err)
		}
		if findings, err := model.Validate(&snapshot); err != nil || len(findings) > 0 {
			t.Fatal(findings, err)
		}
		records := map[string]model.Record{}
		var req model.Record
		for _, r := range snapshot.Elements {
			records[model.String(r, "id")] = r
			if model.String(r, "kind") == "Requirement" {
				req = r
			}
		}
		name := func(r model.Record) string {
			spans, _ := model.Object(r["name"])["spans"].([]any)
			if len(spans) == 0 {
				return ""
			}
			span := model.Object(spans[0])
			for _, src := range snapshot.Sources {
				if src["id"] == span["source"] {
					return strings.TrimSpace(strings.TrimLeft(model.String(src, "markdown")[model.Number(span["start"]):model.Number(span["end"])], "#"))
				}
			}
			return ""
		}
		var targets, attributes []string
		for _, edge := range snapshot.Relationships() {
			if edge["subject"] != req["id"] {
				continue
			}
			switch edge["projection_path"] {
			case "Requirement.applies_to":
				targets = append(targets, model.String(edge, "object"))
			case "Requirement.quality_attributes":
				attributes = append(attributes, model.String(edge, "object"))
			}
		}
		if name(records[model.String(req, "specification_owner")]) != stage.owner || len(targets) != 1 || name(records[targets[0]]) != "Product" || len(attributes) != 1 || name(records[attributes[0]]) != "Durable" {
			t.Fatalf("unexpected ownership, targets or attributes in %v", req)
		}
		run("spec", "build")
		report := filepath.Join("specification", stage.report, "requirements/save.requirement.spec-report.md")
		content, err := os.ReadFile(report)
		if err != nil || !strings.Contains(string(content), "Persist each accepted edit.") {
			t.Fatal(report, err)
		}
		run("check")
	}
}

func TestOwnerLookupAmbiguityAndExactPaths(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: local\n")
	writeFixture(t, "app/app.system.spec.md", "---\ntype: system-spec\n---\n# App\n## Purpose\nSupport work.\n## Scope\nThe app.\n")
	const a = "app/handlers/retry.requirement.spec.md"
	const b = "app/specs/retry.requirement.spec.md"
	for _, file := range []string{a, b} {
		writeFixture(t, file, "---\ntype: requirement-spec\n---\n# Retry\n")
	}
	for _, action := range []string{"show", "check"} {
		code, out, errOut := invoke("requirement", action, "retry", "--path", "app")
		if code == 0 || !strings.Contains(out+errOut, a) || !strings.Contains(out+errOut, b) {
			t.Fatal(code, out, errOut)
		}
		for _, file := range []string{a, b} {
			if code, out, errOut := invoke("requirement", action, file); code != 0 {
				t.Fatal(code, out, errOut)
			}
		}
	}
	writeFixture(t, "app/specs/wrong.component.spec.md", "---\ntype: component-spec\n---\n# Wrong\n")
	if code, out, errOut := invoke("check"); code != 1 || !strings.Contains(out, "layout.subject-ambiguous") {
		t.Fatal(code, out, errOut)
	}
}

func TestSpecsDirectoryAtProjectBoundary(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: outer\n")
	writeFixture(t, "specs/dotspec.yaml", "name: inner\n")
	writeFixture(t, "specs/tool.system.spec.md", "---\ntype: system-spec\n---\n# Tool\n## Purpose\nHelp.\n## Scope\nInner project.\n")
	if code, out, errOut := invoke("system", "list"); code != 0 || !strings.Contains(out, "No systems") {
		t.Fatal(code, out, errOut)
	}
	if code, out, errOut := invoke("system", "show", "--path", "specs"); code != 0 || !strings.Contains(out, "Inner project.") {
		t.Fatal(code, out, errOut)
	}
	if code, out, errOut := invoke("system", "init", "Outer"); code == 0 || !strings.Contains(errOut, "another project") {
		t.Fatal(code, out, errOut)
	}
}
