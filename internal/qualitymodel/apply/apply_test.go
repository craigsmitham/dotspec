package apply

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/system"
)

func target(t *testing.T, raw string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "specs", "system.spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(raw), 0640); err != nil {
		t.Fatal(err)
	}
	return path
}

const plain = "---\ntype: system\n---\n# System\n## Purpose\nP.\n## Scope\nS.\n"

func TestDefaultPreviewApplyAndConflict(t *testing.T) {
	raw := "\ufeff" + strings.ReplaceAll(plain+"## Quality attributes\nLocal narrative.\n### Existing\n**ID:** existing\nKeep this definition.\n## Tail\n<!-- untouched -->\n", "\n", "\r\n")
	path := target(t, raw)
	m, err := qualitymodel.Show(".", qualitymodel.Default)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := Run(m, path, []string{"reliability", "security"}, true)
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if string(before) != raw {
		t.Fatal("preview wrote target")
	}
	result, err := Run(m, path, []string{"reliability", "security"}, false)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != preview.Content || result.Content != preview.Content {
		t.Fatal("preview differs from write")
	}
	inserted := strings.Index(raw, "## Tail")
	if !bytes.HasPrefix(after, []byte(raw[:inserted])) || !bytes.HasSuffix(after, []byte(raw[inserted:])) {
		t.Fatal("existing bytes modified")
	}
	if bytes.Contains(bytes.ReplaceAll(after, []byte("\r\n"), nil), []byte("\n")) {
		t.Fatal("mixed line endings introduced")
	}
	info, _ := os.Stat(path)
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0640 {
		t.Fatal("mode changed")
	}
	parsed, err := system.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.QualityAttributes) != 3 || parsed.QualityAttributes[0].ID != "existing" || parsed.QualityAttributes[1].ID != "reliability" {
		t.Fatalf("wrong attribute order: %+v", parsed)
	}
	if !strings.Contains(parsed.QualityAttributes[1].Definition, "quality model definition") {
		t.Fatal("source link missing")
	}
	if _, err := Run(m, path, []string{"reliability", "maintainability"}, false); err == nil {
		t.Fatal("conflict accepted")
	}
	unchanged, _ := os.ReadFile(path)
	if !bytes.Equal(unchanged, after) {
		t.Fatal("partial write after conflict")
	}
	if _, err := Run(m, path, []string{"unknown"}, false); err == nil {
		t.Fatal("unknown selection accepted")
	}
}

func TestDefaultWholeModelAndAppend(t *testing.T) {
	path := target(t, strings.TrimSuffix(plain, "\n"))
	m, _ := qualitymodel.Show(".", qualitymodel.Default)
	r, err := Run(m, path, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	s, err := system.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Added) != 49 || len(s.QualityAttributes) != 49 {
		t.Fatalf("default population: %+v", r)
	}
}

func TestLocalAncestorClosureAndLinkRelocation(t *testing.T) {
	path := target(t, plain)
	source := filepath.Join(filepath.Dir(path), "quality-models", "local.quality-model.spec.md")
	data := []byte("---\ntype: quality-model\n---\n# Local\n## Sources\n[doc]: ../../guide.md#topic \"Guide\"\n## Quality attributes\n### Narrow\n**ID:** narrow\n**Sub-attribute of:** broad\nSee [details][doc] and [broad](#broad).\n\n`[code](relative.md)`\n\n![image](../../image.png)\n### Broad\n**ID:** broad\nThe parent concern.\n### Other\n**ID:** other\nUnselected.\n")
	m, err := qualitymodel.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m.Path = source
	m.Slug = "local"
	r, err := Run(m, path, []string{"narrow"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(r.Added, ",") != "narrow,broad" {
		t.Fatalf("closure/order: %v", r.Added)
	}
	for _, text := range []string{"[details](<../guide.md#topic> \"Guide\")", "[broad](<quality-models/local.quality-model.spec.md#broad>)", "`[code](relative.md)`", "![image](<../image.png>)"} {
		if !strings.Contains(r.Content, text) {
			t.Errorf("missing %q in %s", text, r.Content)
		}
	}
	if strings.Contains(r.Content, "### Other") {
		t.Fatal("unselected attribute copied")
	}
}

func TestFailureLeavesOriginal(t *testing.T) {
	m, _ := qualitymodel.Show(".", "furps")
	for _, raw := range []string{plain + "```text\nunclosed\n", plain + "## Quality attributes\n```text\nunclosed\n"} {
		path := target(t, raw)
		if _, err := Run(m, path, nil, false); err == nil {
			t.Fatal("unparseable insertion accepted")
		}
		data, _ := os.ReadFile(path)
		if string(data) != raw {
			t.Fatal("failure changed bytes")
		}
	}
	path := target(t, plain)
	link := filepath.Join(t.TempDir(), "system.spec.md")
	if err := os.Symlink(path, link); err == nil {
		if _, err := Run(m, link, nil, false); err == nil {
			t.Fatal("symlink target accepted")
		}
	}
}

func TestISOLeafCopiesItsSingleParent(t *testing.T) {
	m, err := qualitymodel.Show(".", qualitymodel.Default)
	if err != nil {
		t.Fatal(err)
	}
	path := target(t, plain)
	r, err := Run(m, path, []string{"availability"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(r.Added, ",") != "reliability,availability" {
		t.Fatalf("selection: %v", r.Added)
	}
	s, err := system.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(s.FindQualityAttribute("availability").SubattributeOf, ",") != "reliability" {
		t.Fatal("lost parent")
	}
	if s.FindQualityAttribute("recoverability") != nil {
		t.Fatal("unselected sibling copied")
	}
}

func TestQ42SharedQualityCopiesAllDimensions(t *testing.T) {
	m, err := qualitymodel.Show(".", "q42")
	if err != nil {
		t.Fatal(err)
	}
	path := target(t, plain)
	r, err := Run(m, path, []string{"resilience", "resilience"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(r.Added, ",") != "reliable,secure,resilience" {
		t.Fatalf("selection/order: %v", r.Added)
	}
	s, err := system.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	a := s.FindQualityAttribute("resilience")
	if strings.Join(a.SubattributeOf, ",") != "reliable,secure" {
		t.Fatal("lost grouping")
	}
	if !strings.Contains(a.Definition, "creativecommons.org/licenses/by-sa/4.0/") || !strings.Contains(a.Definition, "0185215bfc133bac6a8091029b79b7104e47e351") {
		t.Fatal("lost attribution")
	}
	if got := s.QualityAttributeAncestors([]string{"resilience"}); strings.Join(got, ",") != "reliable,secure" {
		t.Fatalf("navigation: %v", got)
	}
}
