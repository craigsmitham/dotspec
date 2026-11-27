package apply

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/qualitymodel"
	"github.com/dotspec-md/dotspec/internal/system"
)

func target(t *testing.T, raw string) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "specs", "system.system.spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(raw), 0640); err != nil {
		t.Fatal(err)
	}
	return path
}

const iso = "iso-25010-2023"

// shared groups one term under two roots, so adoption must copy both parents.
const shared = "---\ntype: quality-model\nid: shared\n---\n# Shared qualities\n## Quality attributes\n### Reliable\n**ID:** reliable\nKeeps working.\n### Secure\n**ID:** secure\nResists misuse.\n### Resilience\n**ID:** resilience\n**Sub-attribute of:** reliable, secure\nRecovers from disruption, including attacks.\n### Safe\n**ID:** safe\nAvoids harm.\n"

func sharedModel(t *testing.T) *qualitymodel.QualityModel {
	t.Helper()
	path := filepath.Join(t.TempDir(), "shared.quality-model.md")
	if err := os.WriteFile(path, []byte(shared), 0644); err != nil {
		t.Fatal(err)
	}
	m, err := qualitymodel.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

const plain = "---\ntype: system-spec\n---\n# System\n## Purpose\nP.\n## Scope\nS.\n"

func TestDefaultPreviewAndApply(t *testing.T) {
	raw := plain
	path := target(t, raw)
	m, err := qualitymodel.Show(".", iso)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := Run(m, path, Options{Attributes: []string{"reliability", "security"}, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	result, err := Run(m, path, Options{Attributes: []string{"reliability", "security"}, DryRun: false})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(result.Entries, preview.Entries) {
		t.Fatal("preview differs")
	}
	for _, f := range preview.Entries {
		got, err := os.ReadFile(f.Path)
		if err != nil || string(got) != f.Content {
			t.Fatal(f.Path, err)
		}
	}
	parsed, err := system.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.QualityAttributes) != 2 || parsed.QualityAttributes[0].ID != "reliability" {
		t.Fatalf("wrong attribute order: %+v", parsed)
	}
	if !strings.Contains(parsed.QualityAttributes[0].Definition, "quality model definition") {
		t.Fatal("source link missing")
	}
	if _, err := Run(m, path, Options{Attributes: []string{"reliability", "maintainability"}}); err == nil || !strings.Contains(err.Error(), "reliability already exists") {
		t.Fatal("a named attribute the system declares must be refused", err)
	}
	if next, err := Run(m, path, Options{Attributes: []string{"maintainability"}}); err != nil || strings.Join(next.Added, ",") != "maintainability" {
		t.Fatal("expected one addition", err, next.Added)
	}
	if after, _ := os.ReadFile(path); string(after) != raw {
		t.Fatal("system source changed")
	}
	if _, err := Run(m, path, Options{Attributes: []string{"unknown"}, DryRun: false}); err == nil {
		t.Fatal("unknown selection accepted")
	}
}

func TestEveryTermInOneAdoption(t *testing.T) {
	path := target(t, plain)
	m, _ := qualitymodel.Show(".", iso)
	var ids []string
	for _, a := range m.QualityAttributes {
		ids = append(ids, a.ID)
	}
	r, err := Run(m, path, Options{Attributes: ids, DryRun: false})
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
	source := filepath.Join(filepath.Dir(path), "quality-models", "local.quality-model.md")
	data := []byte("---\ntype: quality-model\n---\n# Local\n## Sources\n[doc]: ../../guide.md#topic \"Guide\"\n## Quality attributes\n### Narrow\n**ID:** narrow\n**Sub-attribute of:** broad\nSee [details][doc] and [broad](#broad).\n\n`[code](relative.md)`\n\n![image](../../image.png)\n### Broad\n**ID:** broad\nThe parent concern.\n### Other\n**ID:** other\nUnselected.\n")
	m, err := qualitymodel.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m.Path = source
	m.Slug = "local"
	r, err := Run(m, path, Options{Attributes: []string{"narrow"}, DryRun: false})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(r.Added, ",") != "narrow,broad" {
		t.Fatalf("closure/order: %v", r.Added)
	}
	for _, text := range []string{"[details](<../../../../guide.md#topic> \"Guide\")", "[broad](<../../../quality-models/local.quality-model.md#broad>)", "`[code](relative.md)`", "![image](<../../../../image.png>)"} {
		if !strings.Contains(r.Entries[0].Content, text) {
			t.Errorf("missing %q in %s", text, r.Entries[0].Content)
		}
	}
	if strings.Contains(r.Entries[0].Content, "### Other") {
		t.Fatal("unselected attribute copied")
	}
}

func TestLinkRelocationKeepsCodeSpanLabels(t *testing.T) {
	path := target(t, plain)
	source := filepath.Join(filepath.Dir(path), "quality-models", "local.quality-model.md")
	data := []byte("---\ntype: quality-model\n---\n# Local\n## Quality attributes\n### Narrow\n**ID:** narrow\nSee [`]`](../../guide.md \"ends in \\\\\").\n")
	m, err := qualitymodel.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m.Path = source
	m.Slug = "local"
	r, err := Run(m, path, Options{Attributes: []string{"narrow"}, DryRun: false})
	if err != nil {
		t.Fatal(err)
	}
	if want := "See [`]`](<../../../guide.md> \"ends in \\\\\")."; !strings.Contains(r.Entries[0].Content, want) {
		t.Fatalf("missing %q in %s", want, r.Entries[0].Content)
	}
}

func TestFailureLeavesOriginal(t *testing.T) {
	m, _ := qualitymodel.Show(".", iso)
	for _, raw := range []string{"invalid specification", strings.ReplaceAll(plain, "# System", "# Wrong name")} {
		path := target(t, raw)
		if _, err := Run(m, path, Options{Attributes: []string{"recoverability"}, DryRun: false}); err == nil {
			t.Fatal("unparseable insertion accepted")
		}
		data, _ := os.ReadFile(path)
		if string(data) != raw {
			t.Fatal("failure changed bytes")
		}
	}
	path := target(t, plain)
	link := filepath.Join(t.TempDir(), "target.system.spec.md")
	if err := os.Symlink(path, link); err == nil {
		if _, err := Run(m, link, Options{Attributes: []string{"recoverability"}, DryRun: false}); err == nil {
			t.Fatal("symlink target accepted")
		}
	}
}

func TestISOLeafCopiesItsSingleParent(t *testing.T) {
	m, err := qualitymodel.Show(".", iso)
	if err != nil {
		t.Fatal(err)
	}
	path := target(t, plain)
	r, err := Run(m, path, Options{Attributes: []string{"availability"}, DryRun: false})
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

func TestSharedQualityCopiesAllDimensions(t *testing.T) {
	m := sharedModel(t)
	path := target(t, plain)
	r, err := Run(m, path, Options{Attributes: []string{"resilience", "resilience"}, DryRun: false})
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
	if s.FindQualityAttribute("safe") != nil {
		t.Fatal("unselected root copied")
	}
	if got := s.QualityAttributeAncestors([]string{"resilience"}); strings.Join(got, ",") != "reliable,secure" {
		t.Fatalf("navigation: %v", got)
	}
}

func TestAdoptUnderTailoredParent(t *testing.T) {
	path := target(t, plain)
	// The system uses its own title, directory, and grouping for the model ID.
	parentDir := filepath.Join(filepath.Dir(path), "quality", "local-resilience")
	parent := filepath.Join(parentDir, "local-resilience.quality-attribute.spec.md")
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatal(err)
	}
	original := "---\ntype: quality-attribute-spec\nid: reliable\n---\n# Local resilience\n## Definition\nKeep acknowledged edits.\n"
	if err := os.WriteFile(parent, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	m := sharedModel(t)
	opts := Options{Attributes: []string{"resilience", "resilience"}, DryRun: true}
	preview, err := Run(m, path, opts)
	if err != nil {
		t.Fatal(err)
	}
	opts.DryRun = false
	actual, err := Run(m, path, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(preview.Entries, actual.Entries) || len(actual.Entries) != 3 {
		t.Fatalf("unexpected decisions %+v", actual.Entries)
	}
	e := actual.Entries[0]
	if e.ID != "reliable" || e.Role != "ancestor" || e.Decision != "reuse" || e.Path != parent || e.Content != "" {
		t.Fatal(e)
	}
	child := actual.Entries[2]
	if child.Path != filepath.Join(parentDir, "resilience", "resilience.quality-attribute.spec.md") || child.Role != "requested" {
		t.Fatal(child)
	}
	s, err := system.Read(path)
	if err != nil || !slices.Equal(s.FindQualityAttribute("resilience").SubattributeOf, []string{"reliable", "secure"}) {
		t.Fatal(err)
	}
	got, err := os.ReadFile(parent)
	if err != nil || string(got) != original {
		t.Fatal("changed existing definition", err)
	}

	// A different catalog groups the same ID under an unrelated root.
	other, err := qualitymodel.Parse([]byte("---\ntype: quality-model\n---\n# Team\n## Quality attributes\n### Other\n**ID:** other\nOther grouping.\n### Reliable\n**ID:** reliable\n**Sub-attribute of:** other\nShared term.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Run(other, path, Options{Attributes: []string{"reliable"}}); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatal("a declared attribute was adopted again", err)
	}
}

func TestBlankAttributes(t *testing.T) {
	path := target(t, plain)
	preview, err := Blank(path, Options{Attributes: []string{"speed", "data-retention"}, DryRun: true})
	if err != nil || len(preview.Entries) != 2 || !strings.Contains(preview.Entries[1].Content, "# Data retention\n") {
		t.Fatal(preview, err)
	}
	if _, err := os.Stat(preview.Entries[0].Path); !os.IsNotExist(err) {
		t.Fatal("preview wrote a file")
	}
	if _, err := Blank(path, Options{Attributes: []string{"speed"}}); err != nil {
		t.Fatal(err)
	}
	parsed, err := system.Read(path)
	if err != nil || len(parsed.QualityAttributes) != 1 || parsed.QualityAttributes[0].ID != "speed" {
		t.Fatal(parsed, err)
	}
	for _, ids := range [][]string{{"speed"}, {"Bad--ID"}, {"twice", "twice"}} {
		if _, err := Blank(path, Options{Attributes: ids}); err == nil {
			t.Errorf("%v: accepted", ids)
		}
	}
}

func TestSiblingAdoptionsSerialize(t *testing.T) {
	path := target(t, plain)
	m, _ := qualitymodel.Show(".", iso)
	ready := make(chan struct{})
	done := make(chan error, 2)
	for _, id := range []string{"recoverability", "fault-tolerance"} {
		go func() { <-ready; _, err := Run(m, path, Options{Attributes: []string{id}}); done <- err }()
	}
	close(ready)
	for range 2 {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	s, err := system.Read(path)
	if err != nil || len(s.QualityAttributes) != 3 {
		t.Fatal(s, err)
	}
}

func TestPublicationRollback(t *testing.T) {
	root := t.TempDir()
	existing := filepath.Join(root, "tailored")
	if err := os.Mkdir(existing, 0755); err != nil {
		t.Fatal(err)
	}
	files := []Entry{{Path: filepath.Join(existing, "child", "first.md"), Content: "first"}, {Path: filepath.Join(root, "second", "second.md"), Content: "second"}}
	failSecond := func(from, to string) error {
		if to == files[1].Path {
			return os.ErrPermission
		}
		return os.Link(from, to)
	}
	if err := publish(root, files, failSecond); err == nil {
		t.Fatal("expected install failure")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 || entries[0].Name() != "tailored" {
		t.Fatal(entries, err)
	}
	children, err := os.ReadDir(existing)
	if err != nil || len(children) != 0 {
		t.Fatal(children, err)
	}
	if err := publish(root, files, os.Link); err != nil {
		t.Fatal(err)
	}
}

func TestRejectPartialSelectionAndDestinationCollision(t *testing.T) {
	path := target(t, plain)
	m, _ := qualitymodel.Show(".", iso)
	q := filepath.Join(filepath.Dir(path), "quality")
	if _, err := Run(m, path, Options{Attributes: []string{"recoverability", "typo"}}); err == nil {
		t.Fatal("accepted unknown attribute")
	}
	if _, err := os.Stat(q); !os.IsNotExist(err) {
		t.Fatal("unexpected write", err)
	}
	// A directory occupies the planned file path. No sibling may be installed.
	collision := filepath.Join(q, "reliability", "recoverability", "recoverability.quality-attribute.spec.md")
	if err := os.MkdirAll(collision, 0755); err != nil {
		t.Fatal(err)
	}
	for _, dry := range []bool{true, false} {
		if _, err := Run(m, path, Options{Attributes: []string{"recoverability", "fault-tolerance"}, DryRun: dry}); err == nil {
			t.Fatal("accepted collision")
		}
		if _, err := os.Stat(filepath.Join(q, "reliability", "reliability.quality-attribute.spec.md")); !os.IsNotExist(err) {
			t.Fatal("partial write", err)
		}
	}
}

func TestBundledCitationUsesReleaseAndSourceHeading(t *testing.T) {
	m, err := qualitymodel.Show(".", iso)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := Run(m, target(t, plain), Options{Attributes: []string{"recoverability"}, DryRun: true, Revision: "fixture-commit"})
	if err != nil {
		t.Fatal(err)
	}
	leaf := plan.Entries[len(plan.Entries)-1]
	if !strings.Contains(leaf.Content, "/blob/fixture-commit/internal/qualitymodel/specs/quality-models/iso-25010-2023.quality-model.md#recoverability") {
		t.Fatal(leaf.Content)
	}
}
