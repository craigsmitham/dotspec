package model

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dotspec-md/dotspec/internal/speckinds"
)

func TestExportEmptyScope(t *testing.T) {
	result, err := Export(t.TempDir(), speckinds.All()...)
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("empty export: %v %+v", err, result)
	}
	if len(result.Snapshot.Sources) != 0 {
		t.Fatalf("unexpected empty scope: %+v", result.Snapshot)
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) != 0 {
		t.Fatalf("empty snapshot: %v %+v", err, findings)
	}
}

func TestExportDirectoryImageAndNativeFragments(t *testing.T) {
	dir := t.TempDir()
	for _, directory := range []string{"specs/concepts", "assets"} {
		if err := os.MkdirAll(filepath.Join(dir, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string][]byte{
		"specs/concepts/media.concept.spec.md": []byte("---\ntype: concept-spec\n---\n# Media\n## Purpose\nInspect sources.\n## State\n![Image](../../assets/pixel.bin)\n[Directory](../../assets/)\n[Native selector](../../contract.txt#rule)\n## Actions\n### Read\n## Operational principles\nFollow links.\n"),
		"assets/pixel.bin":                     {0, 1, 0xff, 0x0d, 0x0a},
		"contract.txt":                         []byte("native selector semantics\n"),
	}
	for path, data := range files {
		if err := os.WriteFile(filepath.Join(dir, path), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Export(dir, speckinds.All()...)
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("export: %v %d", err, result.ExitCode)
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) != 0 {
		t.Fatalf("snapshot: %v %+v", err, findings)
	}
	directory, image, native := false, false, false
	for _, source := range result.Snapshot.Sources {
		directory = directory || String(source, "kind") == "SourceDirectory"
		image = image || filepath.Base(String(source, "locator")) == "pixel.bin"
	}
	for _, ref := range append(result.Snapshot.References, result.Snapshot.LabelOccurrences...) {
		resolution := Object(ref["resolution"])
		native = native || String(resolution, "rule") == "native-fragment-semantics/1" && String(resolution, "status") == "not_attempted"
	}
	if !directory || !image || !native {
		t.Fatalf("directory=%v image=%v native-selector-preserved=%v", directory, image, native)
	}
	first, _ := json.Marshal(result.Snapshot)
	again, err := Export(dir, speckinds.All()...)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := json.Marshal(again.Snapshot)
	if !bytes.Equal(first, second) {
		t.Fatal("unchanged input did not produce an identical snapshot")
	}
}

func TestExportRecordsNativeDefinitionLine(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "specs/concepts"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"specs/concepts/queue.concept.spec.md": "---\ntype: concept-spec\nspecifications:\n  - source: ../../queue.qnt#inv\n    defines: requirements\n    name: Queue invariant\nverification:\n  - source: ../../queue.qnt#missing\n    name: Queue check\n---\n# Queue\n## Purpose\nOrder work.\n## State\nItems.\n## Actions\n### Push\n## Operational principles\nFirst in, first out.\n",
		"queue.qnt":                            "module queue {\n  val inv = true\n}\n",
	}
	for path, data := range files {
		if err := os.WriteFile(filepath.Join(dir, path), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Export(dir, speckinds.All()...)
	if err != nil {
		t.Fatal(err)
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) != 0 {
		t.Fatalf("snapshot: %v %+v", err, findings)
	}
	lines := map[string]any{}
	for _, ref := range result.Snapshot.References {
		if resolution := Object(ref["resolution"]); String(resolution, "rule") == "native-definition/1" {
			lines[String(ref, "authored_target")+" "+String(resolution, "status")] = resolution["line"]
		}
	}
	if len(lines) != 2 || lines["../../queue.qnt#inv resolved"] != 2 || lines["../../queue.qnt#missing unresolved"] != nil {
		t.Fatalf("native definition resolutions: %v", lines)
	}
}

func TestExportResolvesContextForSourceCheckRejects(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"specs/app.system.spec.md":                    "---\ntype: system-spec\n---\n# App\n## Purpose\nRead.\n## Scope\nLocal.\n",
		"specs/requirements/save.requirement.spec.md": "---\ntype: requirement-spec\n---\n# Save\n## Criterion\nA criterion.\n## Bogus heading\nText.\n",
	}
	for path, data := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Export(dir, speckinds.All()...)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, source := range result.Snapshot.Sources {
		ids[filepath.Base(String(source, "locator"))] = String(source, "id")
	}
	for _, context := range result.Snapshot.ContextSources {
		if String(context, "role") == "system_context" && String(context, "source") == ids["app.system.spec.md"] {
			for _, requester := range Strings(context["requested_by"]) {
				if requester == ids["save.requirement.spec.md"] {
					return
				}
			}
		}
	}
	t.Fatalf("rejected requirement did not request its system context: %+v", result.Snapshot.ContextSources)
}

func TestExportDoesNotSkipInaccessibleParentSystem(t *testing.T) {
	dir := t.TempDir()
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, output)
	}
	body := []byte("---\ntype: system-spec\n---\n# System\n## Purpose\nRead.\n## Scope\nLocal.\n")
	for _, path := range []string{"specs/system.system.spec.md", "middle/child/specs/system.system.spec.md"} {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, body, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "middle/specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "specs/system.system.spec.md"), filepath.Join(dir, "middle/specs/system.system.spec.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	result, err := Export(dir, speckinds.All()...)
	if err != nil {
		t.Fatal(err)
	}
	for _, element := range result.Snapshot.Elements {
		if String(element, "kind") == "System" && element["enclosing_system"] != nil {
			t.Fatalf("parent relation crossed a nearer symlink boundary: %+v", element)
		}
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) > 0 {
		t.Fatalf("snapshot: %v %+v", err, findings)
	}
}

func TestExportReferencedClassificationDoesNotSkipNearerBoundary(t *testing.T) {
	dir := t.TempDir()
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, output)
	}
	files := map[string]string{
		"specs/root.system.spec.md":                             "---\ntype: system-spec\n---\n# Root\n## Purpose\nRead.\n## Scope\nLocal.\n",
		"specs/quality/useful/useful.quality-attribute.spec.md": "---\ntype: quality-attribute-spec\nid: useful\n---\n# Useful\n\n## Definition\n\nA concern.\n",

		"consumer/specs/concepts/entry.concept.spec.md":              "---\ntype: concept-spec\n---\n# Entry\n## Purpose\nRead.\n## State\n[Target](../../../middle/child/specs/requirements/target.requirement.spec.md)\n## Actions\n### Read\n## Operational principles\nRead a target.\n",
		"middle/child/specs/requirements/target.requirement.spec.md": "---\ntype: requirement-spec\nquality_attributes: [useful]\n---\n# Target\n## Criterion\nA criterion.\n",
	}
	for path, data := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "middle/specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "specs/root.system.spec.md"), filepath.Join(dir, "middle/specs/root.system.spec.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	result, err := Export(filepath.Join(dir, "consumer"), speckinds.All()...)
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("export: %v %d", err, result.ExitCode)
	}
	if len(result.Snapshot.LabelOccurrences) != 1 {
		t.Fatalf("expected one retained label, got %d", len(result.Snapshot.LabelOccurrences))
	}
	label := result.Snapshot.LabelOccurrences[0]
	for _, ref := range append(result.Snapshot.References, result.Snapshot.LabelOccurrences...) {
		resolution := Object(ref["resolution"])
		if String(ref, "id") == String(label, "id") && String(resolution, "status") != "not_attempted" {
			t.Fatalf("classification skipped nearer inaccessible attribute declarations: %+v", resolution)
		}
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) > 0 {
		t.Fatalf("snapshot: %v %+v", err, findings)
	}
}
