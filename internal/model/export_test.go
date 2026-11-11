package model

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func exportTypes() []check.Type {
	return []check.Type{system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType(), quality.CheckType()}
}

func TestExportEmptyScope(t *testing.T) {
	result, err := Export(t.TempDir(), exportTypes()...)
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("empty export: %v %+v", err, result)
	}
	if len(result.Snapshot.Sources) != 0 || String(result.Snapshot.Coverage, "status") != "complete_for_profile" {
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
		"specs/concepts/media.concept.spec.md": []byte("---\ntype: concept\n---\n# Media\n## Purpose\nInspect sources.\n## State\n![Image](../../assets/pixel.bin)\n[Directory](../../assets/)\n[Native selector](../../contract.txt#rule)\n## Actions\nRead.\n## Operational principle\nFollow links.\n"),
		"assets/pixel.bin":                     {0, 1, 0xff, 0x0d, 0x0a},
		"contract.txt":                         []byte("native selector semantics\n"),
	}
	for path, data := range files {
		if err := os.WriteFile(filepath.Join(dir, path), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Export(dir, exportTypes()...)
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
	for _, resolution := range result.Snapshot.Resolutions {
		native = native || String(resolution, "rule") == "native-fragment-semantics/1" && String(resolution, "status") == "not_attempted"
	}
	if !directory || !image || !native {
		t.Fatalf("directory=%v image=%v native-selector-preserved=%v", directory, image, native)
	}
	first, _ := json.Marshal(result.Snapshot)
	again, err := Export(dir, exportTypes()...)
	if err != nil {
		t.Fatal(err)
	}
	second, _ := json.Marshal(again.Snapshot)
	if !bytes.Equal(first, second) {
		t.Fatal("unchanged input did not produce an identical snapshot")
	}
}

func TestExportDoesNotSkipInaccessibleParentSystem(t *testing.T) {
	dir := t.TempDir()
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, output)
	}
	body := []byte("---\ntype: system\n---\n# System\n## Purpose\nRead.\n## Scope\nLocal.\n")
	for _, path := range []string{"specs/system.spec.md", "middle/child/specs/system.spec.md"} {
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
	if err := os.Symlink(filepath.Join(dir, "specs/system.spec.md"), filepath.Join(dir, "middle/specs/system.spec.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	result, err := Export(dir, exportTypes()...)
	if err != nil {
		t.Fatal(err)
	}
	for _, element := range result.Snapshot.Elements {
		if String(element, "kind") == "System" && element["parent_system"] != nil {
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
		"specs/system.spec.md":                              "---\ntype: system\n---\n# Root\n## Purpose\nRead.\n## Scope\nLocal.\n## Quality attributes\n### Useful\n**ID:** useful\nA concern.\n",
		"consumer/specs/concepts/entry.concept.spec.md":     "---\ntype: concept\n---\n# Entry\n## Purpose\nRead.\n## State\n[Target](../../../middle/child/specs/quality/target.quality.spec.md)\n## Actions\nRead.\n## Operational principle\nRead a target.\n",
		"middle/child/specs/quality/target.quality.spec.md": "---\ntype: quality\nattributes: [useful]\n---\n# Target\n## Criterion\nA criterion.\n",
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
	if err := os.Symlink(filepath.Join(dir, "specs/system.spec.md"), filepath.Join(dir, "middle/specs/system.spec.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	result, err := Export(filepath.Join(dir, "consumer"), exportTypes()...)
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("export: %v %d", err, result.ExitCode)
	}
	if len(result.Snapshot.LabelOccurrences) != 1 {
		t.Fatalf("expected one retained label, got %d", len(result.Snapshot.LabelOccurrences))
	}
	label := result.Snapshot.LabelOccurrences[0]
	for _, resolution := range result.Snapshot.Resolutions {
		if String(resolution, "id") == String(label, "resolution") && String(resolution, "status") != "not_attempted" {
			t.Fatalf("classification skipped nearer inaccessible vocabulary: %+v", resolution)
		}
	}
	findings, err := Validate(result.Snapshot)
	if err != nil || len(findings) > 0 {
		t.Fatalf("snapshot: %v %+v", err, findings)
	}
}
