package specdoc

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const (
	systemPath = "specification/task-list.system.spec-report.md"
	detailLink = "requirements/restart.requirement.spec-report.md"
	detailPath = "specification/" + detailLink
)

func testPlan() *Plan {
	entry := &Page{
		Title: "Task list", RequiredLinks: []string{systemPath},
		Markdown: "---\ngenerated_by: dotspec\n---\n# Task list\n\n[Read the system specification](" + systemPath + ")\n",
	}
	root := &Page{
		Document: Document{Type: "system-spec-report", Documents: "../specs/task-list.system.spec.md#task-list"}, Title: "Task list",
		RequiredLinks: []string{detailLink},
		Markdown:      "---\ntype: system-spec-report\ndocuments: ../specs/task-list.system.spec.md#task-list\ngenerated_by: dotspec\n---\n# Task list\n\n[Resume tasks](" + detailLink + ")\n",
	}
	detail := &Page{
		Document: Document{Type: "requirement-spec-report", Documents: "../../specs/restart.spec.md#restart"}, Title: "Restart",
		RequiredContent: []string{"**Preserve tasks** on the same database."},
		Markdown:        "---\ntype: requirement-spec-report\ndocuments: ../../specs/restart.spec.md#restart\ngenerated_by: dotspec\n---\n# Restart\n\n**Preserve tasks** on the same database.\n",
	}
	return &Plan{Pages: map[string]*Page{"SPEC.md": entry, systemPath: root, detailPath: detail}}
}

func TestCheckRenderedSourceAndRoutes(t *testing.T) {
	dir := t.TempDir()
	p := testPlan()
	if _, err := Write(dir, p); err != nil {
		t.Fatal(err)
	}
	unchanged, err := Write(dir, p)
	if err != nil || len(unchanged.Written)+len(unchanged.Removed) != 0 {
		t.Fatalf("second build: %+v, %v", unchanged, err)
	}
	if ds, err := Check(dir, p); err != nil || len(ds) > 0 {
		t.Fatalf("clean check: %+v, %v", ds, err)
	}
	// Simulate a buggy renderer. Both broken files exactly match its generated
	// bytes; explicit content and navigation rules still detect both losses.
	p.Pages[systemPath].Markdown = strings.ReplaceAll(p.Pages[systemPath].Markdown, "[Resume tasks]("+detailLink+")", "Resume tasks")
	p.Pages[detailPath].Markdown = strings.ReplaceAll(p.Pages[detailPath].Markdown, " on the same database", "")
	for rel, pg := range p.Pages {
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(pg.Markdown), 0644); err != nil {
			t.Fatal(err)
		}
	}
	ds, err := Check(dir, p)
	if err != nil {
		t.Fatal(err)
	}
	rules := map[string]bool{}
	for _, d := range ds {
		rules[d.Rule] = true
	}
	if !rules["spec.navigation"] || !rules["spec.content"] || !rules["spec.unreachable"] || rules["spec.stale"] {
		t.Fatalf("independent verification: %+v", ds)
	}
}

func TestWritePreservesAuthoredFiles(t *testing.T) {
	dir := t.TempDir()
	p := testPlan()
	root := filepath.Join(dir, "SPEC.md")
	if err := os.WriteFile(root, []byte("# Authored\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Write(dir, p)
	var finding *diagnostic.Diagnostic
	if !errors.As(err, &finding) || finding.Rule != "spec.collision" {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "specification")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("partial build: %v", err)
	}
	data, err := os.ReadFile(root)
	if err != nil || string(data) != "# Authored\n" {
		t.Fatal(err, string(data))
	}
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	if _, err := Write(dir, p); err != nil {
		t.Fatal(err)
	}
	notes := filepath.Join(dir, "specification", "requirements", "notes.md")
	if err := os.WriteFile(notes, []byte("# Authored notes\n"), 0644); err != nil {
		t.Fatal(err)
	}
	delete(p.Pages, detailPath)
	p.Pages[systemPath].RequiredLinks = nil
	p.Pages[systemPath].Markdown = strings.Split(p.Pages[systemPath].Markdown, "[")[0]
	changes, err := Write(dir, p)
	if err != nil || len(changes.Removed) != 1 || changes.Removed[0] != detailPath {
		t.Fatal(changes, err)
	}
	if _, err := os.Stat(notes); err != nil {
		t.Fatal(err)
	}
}

func TestWriteRejectsSymlinkedDestination(t *testing.T) {
	dir, other := t.TempDir(), t.TempDir()
	if err := os.Symlink(other, filepath.Join(dir, "specification")); err != nil {
		t.Skip(err)
	}
	if _, err := Write(dir, testPlan()); err == nil {
		t.Fatal("wrote through a linked directory")
	}
	for _, name := range []string{filepath.Join(dir, "SPEC.md"), filepath.Join(other, filepath.Base(systemPath))} {
		if _, err := os.Stat(name); !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(name, err)
		}
	}
}

func TestBuildsSerialize(t *testing.T) {
	dir := t.TempDir()
	p := testPlan()
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			if _, err := Write(dir, p); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if ds, err := Check(dir, p); err != nil || len(ds) != 0 {
		t.Fatal(ds, err)
	}
}

func TestCheckScopedLinksAndIndexMembership(t *testing.T) {
	source := []byte("---\ntype: use-case-spec-report\ndocuments: source.md#return\ngenerated_by: dotspec\n---\n# Return\n\n## First promise\n\n## Second promise\n\n[Edit source](../test.py)\n")
	p := parseGenerated(source)
	page := Page{Title: "Return", Document: Document{Type: "use-case-spec-report", Documents: "source.md#return"}, RequiredSectionLinks: map[string][]string{"first-promise": {"../test.py"}}}
	var rules []string
	report := func(rule, _ string) { rules = append(rules, rule) }
	checkPage(&page, p, report)
	if !slices.Equal(rules, []string{"spec.navigation"}) {
		t.Fatal("link outside governing clause", rules)
	}
	page.RequiredSectionLinks = map[string][]string{"second-promise": {"../test.py"}}
	rules = nil
	checkPage(&page, p, report)
	if len(rules) != 0 {
		t.Fatal(rules)
	}
	p = parseGenerated([]byte("---\ntype: system-spec-report\ndocuments: source.md#tasks\ngenerated_by: dotspec\n---\n# Tasks\n\n## Requirements\n\n- [Restart](restart.md)\n- [Unexpected](other.md)\n"))
	page = Page{Document: Document{Type: "system-spec-report", Documents: "source.md#tasks"}, Title: "Tasks", DirectoryMembers: map[string][]string{"requirements": {"restart.md"}}}
	checkPage(&page, p, report)
	if !slices.Equal(rules, []string{"spec.membership"}) {
		t.Fatal("misplaced index member", rules)
	}
}

func TestCheckRestoresOrderAndHandlesCheckoutLineEndings(t *testing.T) {
	dir := t.TempDir()
	p := testPlan()
	if _, err := Write(dir, p); err != nil {
		t.Fatal(err)
	}
	for rel, page := range p.Pages {
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(strings.ReplaceAll(page.Markdown, "\n", "\r\n")), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if ds, err := Check(dir, p); err != nil || len(ds) != 0 {
		t.Fatal(ds, err)
	}
	page := p.Pages[detailPath]
	page.Markdown += "\n1. Remove the task.\n2. Start the service.\n"
	page.RequiredContent = append(page.RequiredContent, "Start the service.", "Remove the task.")
	_, err := Write(dir, p)
	var d *diagnostic.Diagnostic
	if !errors.As(err, &d) || d.Rule != "spec.content" {
		t.Fatalf("reversed source order: %v", err)
	}
}

func TestWriteRespectsNestedProjectBoundaries(t *testing.T) {
	for _, rel := range []string{"specification", "specification/requirements"} {
		t.Run(filepath.Base(rel), func(t *testing.T) {
			dir := t.TempDir()
			nested := filepath.Join(dir, rel)
			if err := os.MkdirAll(nested, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(nested, "dotspec.yaml"), []byte("name: child\n"), 0644); err != nil {
				t.Fatal(err)
			}
			stale := filepath.Join(nested, "owned-by-child.md")
			if err := os.WriteFile(stale, []byte("---\ngenerated_by: dotspec\n---\n# Child\n"), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := Write(dir, testPlan())
			var d *diagnostic.Diagnostic
			if !errors.As(err, &d) || d.Rule != "spec.collision" {
				t.Fatalf("nested project: %v", err)
			}
			if _, err := os.Stat(stale); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(dir, "SPEC.md")); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("partial build: %v", err)
			}
		})
	}
}

func TestDuplicateReportClaimsIdentifyBothFiles(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	plan := testPlan()
	if _, err := Write(dir, plan); err != nil {
		t.Fatal(err)
	}
	extra := filepath.Join(dir, "specification", "extra.requirement.spec-report.md")
	copied := strings.ReplaceAll(plan.Pages[detailPath].Markdown, "../../specs/", "../specs/")
	if err := os.WriteFile(extra, []byte(copied), 0644); err != nil {
		t.Fatal(err)
	}
	ds, err := Check(dir, plan)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range ds {
		if d.Rule == "spec.duplicate-report" && len(d.Related) == 1 {
			paths := []string{d.Location.Path, d.Related[0].Path}
			if slices.Contains(paths, extra) && slices.Contains(paths, filepath.Join(dir, detailPath)) {
				return
			}
		}
	}
	t.Fatalf("missing duplicate claim: %+v", ds)
}
