package apply

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/system"
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

const plain = "---\ntype: system-spec\n---\n# System\n## Purpose\nP.\n## Scope\nS.\n"

func TestDefaultPreviewApplyAndConflict(t *testing.T) {
	raw := "\ufeff" + strings.ReplaceAll(plain+"## Quality attributes\n## Environment\n<!-- untouched -->\n", "\n", "\r\n")
	path := target(t, raw)
	m, err := qualitymodel.Show(".", qualitymodel.Default)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := Run(m, path, Options{Attributes: []string{"reliability", "security"}, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if string(before) != raw {
		t.Fatal("preview wrote target")
	}
	result, err := Run(m, path, Options{Attributes: []string{"reliability", "security"}, DryRun: false})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != raw || !slices.Equal(result.Entries, preview.Entries) {
		t.Fatal("system source changed or preview differs")
	}
	for _, f := range preview.Entries {
		got, err := os.ReadFile(f.Path)
		if err != nil || string(got) != f.Content {
			t.Fatal(f.Path, err)
		}
	}
	info, _ := os.Stat(path)
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0640 {
		t.Fatal("mode changed")
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
	if next, err := Run(m, path, Options{Attributes: []string{"reliability", "maintainability"}}); err != nil || strings.Join(next.Added, ",") != "maintainability" {
		t.Fatal("expected reuse and one addition", err, next.Added)
	}
	unchanged, _ := os.ReadFile(path)
	if !bytes.Equal(unchanged, after) {
		t.Fatal("partial write after conflict")
	}
	if _, err := Run(m, path, Options{Attributes: []string{"unknown"}, DryRun: false}); err == nil {
		t.Fatal("unknown selection accepted")
	}
}

func TestDefaultWholeModelAndAppend(t *testing.T) {
	path := target(t, strings.TrimSuffix(plain, "\n"))
	m, _ := qualitymodel.Show(".", qualitymodel.Default)
	r, err := Run(m, path, Options{All: true, DryRun: false})
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

func TestFailureLeavesOriginal(t *testing.T) {
	m, _ := qualitymodel.Show(".", "furps")
	for _, raw := range []string{"invalid specification", strings.ReplaceAll(plain, "# System", "# Wrong name")} {
		path := target(t, raw)
		if _, err := Run(m, path, Options{All: true, DryRun: false}); err == nil {
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
		if _, err := Run(m, link, Options{All: true, DryRun: false}); err == nil {
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

func TestQ42SharedQualityCopiesAllDimensions(t *testing.T) {
	m, err := qualitymodel.Show(".", "q42")
	if err != nil {
		t.Fatal(err)
	}
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
	if !strings.Contains(a.Definition, "creativecommons.org/licenses/by-sa/4.0/") || !strings.Contains(a.Definition, "0185215bfc133bac6a8091029b79b7104e47e351") {
		t.Fatal("lost attribution")
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
	m, _ := qualitymodel.Show(".", "q42")
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
	repeat, err := Run(other, path, Options{Attributes: []string{"reliable"}})
	if err != nil || len(repeat.Added) != 0 || len(repeat.Entries) != 1 || repeat.Entries[0].Path != parent || repeat.Entries[0].Role != "requested" {
		t.Fatal(repeat, err)
	}
}

func TestSiblingAdoptionsSerialize(t *testing.T) {
	path := target(t, plain)
	m, _ := qualitymodel.Show(".", qualitymodel.Default)
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
	files := []File{{filepath.Join(existing, "child", "first.md"), "first"}, {filepath.Join(root, "second", "second.md"), "second"}}
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
	m, _ := qualitymodel.Show(".", qualitymodel.Default)
	q := filepath.Join(filepath.Dir(path), "quality")
	for _, opts := range []Options{{}, {Attributes: []string{"recoverability", "typo"}}, {All: true, Attributes: []string{"availability"}}} {
		if _, err := Run(m, path, opts); err == nil {
			t.Fatal("accepted", opts)
		}
		if _, err := os.Stat(q); !os.IsNotExist(err) {
			t.Fatal("unexpected write", err)
		}
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
	m, err := qualitymodel.Show(".", "furps")
	if err != nil {
		t.Fatal(err)
	}
	old := releaseCommit
	releaseCommit = "fixture-commit"
	t.Cleanup(func() { releaseCommit = old })
	plan, err := Run(m, target(t, plain), Options{Attributes: []string{"failure-frequency-and-severity"}, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	leaf := plan.Entries[len(plan.Entries)-1]
	if !strings.Contains(leaf.Content, "/blob/fixture-commit/internal/qualitymodel/specs/quality-models/furps.quality-model.md#frequency-and-severity-of-failure") {
		t.Fatal(leaf.Content)
	}
}
