package job

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeJob(t *testing.T, dir, slug, name string) string {
	t.Helper()
	path := filepath.Join(dir, "specs", "jobs", slug+Suffix)
	writeFile(t, path, strings.Replace(validJob, "Develop Requirements", name, 1))
	return path
}

func writeSystem(t *testing.T, dir, name string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "specs", "system.spec.md"), "---\ntype: system\n---\n# "+name+"\n## Purpose\nHelp teams.\n## Scope\nRequirements.\n")
}

func boundLookup(t *testing.T, dir string) {
	t.Helper()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func gitRepo(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, key := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	return dir
}

func summary(groups []*Group) []string {
	var result []string
	for _, group := range groups {
		name := "None"
		if group.System != nil {
			name = group.System.Name
		}
		result = append(result, name)
		for _, job := range group.Jobs {
			result = append(result, "  "+job.Slug+": "+job.Name)
		}
	}
	return result
}

func TestInitAndLocalIdentity(t *testing.T) {
	dir := t.TempDir()
	boundLookup(t, dir)
	path, err := Init(dir, "develop-requirements")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Show(dir, "develop-requirements")
	if err != nil {
		t.Fatal(err)
	}
	if got.Slug != "develop-requirements" || got.System != nil {
		t.Fatalf("unexpected job: %+v", got)
	}
	writeJob(t, dir, "develop-requirements", "A revised title")
	before, _ := os.ReadFile(path)
	if _, err := Init(dir, "develop-requirements"); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("overwrite error: %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatal("existing content changed")
	}
	got, err = Show(dir, "develop-requirements")
	if err != nil || got.Name != "A revised title" || got.Slug != "develop-requirements" {
		t.Fatalf("identity changed: %+v, %v", got, err)
	}
	if _, err := Init(dir, "manage-requirements"); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(dir, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Show(child, "develop-requirements"); err == nil {
		t.Fatal("show must not find an ancestor's job")
	}
	writeJob(t, child, "develop-requirements", "Child job")
	got, err = Show(child, "develop-requirements")
	if err != nil || got.Name != "Child job" {
		t.Fatalf("local job: %+v, %v", got, err)
	}
}

func TestInvalidSlugsAndPaths(t *testing.T) {
	dir := t.TempDir()
	for _, slug := range []string{"", "../escape", "/escape", "a/b", `a\b`, "Uppercase", "two words", "under_score", "-leading", "trailing-", "two--hyphens", "name.job.spec.md", "café"} {
		if _, err := Init(dir, slug); err == nil {
			t.Errorf("init accepted %q", slug)
		}
		if _, err := Show(dir, slug); err == nil || !strings.Contains(err.Error(), "invalid job slug") {
			t.Errorf("show accepted %q: %v", slug, err)
		}
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid init created files: %v, %v", entries, err)
	}
	for _, slug := range []string{"develop-requirements", "job2", "2-steps"} {
		if err := ValidateSlug(slug); err != nil {
			t.Fatal(err)
		}
	}
	file := writeJob(t, dir, "valid", "Valid")
	for _, path := range []string{filepath.Join(dir, "missing"), file} {
		if _, err := Init(path, "test"); err == nil {
			t.Fatal("init accepted an invalid directory")
		}
		if _, err := Show(path, "test"); err == nil {
			t.Fatal("show accepted an invalid directory")
		}
		if _, err := List(path); err == nil {
			t.Fatal("list accepted an invalid directory")
		}
	}
}

func TestAssociationAndSubtreeGrouping(t *testing.T) {
	dir := gitRepo(t)
	writeSystem(t, dir, "dotspec")
	writeJob(t, dir, "develop-requirements", "Develop Requirements")
	writeJob(t, dir, "manage-requirements", "Manage Requirements")
	child := filepath.Join(dir, "parser")
	writeSystem(t, child, "Parser")
	writeJob(t, child, "develop-requirements", "Develop Parser Requirements")
	deep := filepath.Join(child, "tokens")
	writeJob(t, deep, "read-token", "Read a token")
	for _, tt := range []struct{ path, slug, want string }{
		{dir, "manage-requirements", "dotspec"},
		{child, "develop-requirements", "Parser"},
		{deep, "read-token", "Parser"},
	} {
		got, err := Show(tt.path, tt.slug)
		if err != nil || got.System == nil || got.System.Name != tt.want {
			t.Fatalf("association: %+v, %v", got, err)
		}
	}
	groups, err := List(dir)
	want := []string{"Parser", "  develop-requirements: Develop Parser Requirements", "  read-token: Read a token", "dotspec", "  develop-requirements: Develop Requirements", "  manage-requirements: Manage Requirements"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("got %v, %v; want %v", summary(groups), err, want)
	}
	groups, err = List(deep)
	want = []string{"Parser", "  read-token: Read a token"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("subtree: %v, %v", summary(groups), err)
	}
	if filepath.Base(filepath.Dir(filepath.Dir(groups[0].System.Path))) != "parser" {
		t.Fatal("system above subtree was not resolved")
	}
}

func TestUnassociatedJobsAndDistinctSystemsWithSameName(t *testing.T) {
	dir := gitRepo(t)
	writeJob(t, dir, "unassigned", "Unassigned")
	for _, child := range []string{"a", "b"} {
		path := filepath.Join(dir, child)
		writeSystem(t, path, "Same name")
		writeJob(t, path, "work", "Work in "+child)
	}
	groups, err := List(dir)
	want := []string{"Same name", "  work: Work in a", "Same name", "  work: Work in b", "None", "  unassigned: Unassigned"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("groups: %v, %v", summary(groups), err)
	}
	if groups[0].System.Path == groups[1].System.Path {
		t.Fatal("merged distinct systems")
	}
}

func TestAssociationStopsAtRepositoryBoundary(t *testing.T) {
	outer := t.TempDir()
	writeSystem(t, outer, "Outside")
	for _, worktree := range []bool{false, true} {
		dir := filepath.Join(outer, map[bool]string{false: "repository", true: "worktree"}[worktree])
		writeJob(t, dir, "work", "Work")
		if worktree {
			writeFile(t, filepath.Join(dir, ".git"), "gitdir: elsewhere\n")
		} else {
			boundLookup(t, dir)
		}
		got, err := Show(dir, "work")
		if err != nil || got.System != nil {
			t.Fatalf("crossed boundary: %+v, %v", got, err)
		}
	}
}

func TestMalformedJobsAndAssociatedSystems(t *testing.T) {
	dir := gitRepo(t)
	groups, err := List(dir)
	if err != nil || len(groups) != 0 {
		t.Fatalf("empty list: %v, %v", groups, err)
	}
	path := writeJob(t, dir, "work", "Work")
	writeFile(t, path, "malformed")
	if _, err := Show(dir, "work"); err == nil || !strings.Contains(err.Error(), "work.job.spec.md") {
		t.Fatalf("malformed show: %v", err)
	}
	if _, err := List(dir); err == nil || !strings.Contains(err.Error(), "work.job.spec.md") {
		t.Fatalf("malformed list: %v", err)
	}
	writeJob(t, dir, "work", "Work")
	writeFile(t, filepath.Join(dir, "specs", "system.spec.md"), "malformed")
	if _, err := Show(dir, "work"); err == nil || !strings.Contains(err.Error(), "system.spec.md") {
		t.Fatalf("malformed association: %v", err)
	}
	if _, err := List(dir); err == nil || !strings.Contains(err.Error(), "system.spec.md") {
		t.Fatalf("malformed group: %v", err)
	}
	writeSystem(t, dir, "Valid")
	writeFile(t, filepath.Join(dir, "specs", "jobs", "Invalid_slug.job.spec.md"), validJob)
	if _, err := List(dir); err == nil || !strings.Contains(err.Error(), "invalid job slug") {
		t.Fatalf("invalid filename: %v", err)
	}
}

func TestDiscoveryIgnoresAndLinks(t *testing.T) {
	dir := gitRepo(t)
	writeJob(t, dir, "tracked", "Tracked")
	git(t, dir, "add", "specs/jobs/tracked.job.spec.md")
	writeFile(t, filepath.Join(dir, ".gitignore"), "ignored/\ntracked.job.spec.md\nchoices/*\n!choices/keep/\n")
	writeFile(t, filepath.Join(dir, "ignored", "bad.job.spec.md"), "malformed but ignored")
	writeFile(t, filepath.Join(dir, "choices", "omit", "bad.job.spec.md"), "ignored")
	writeJob(t, filepath.Join(dir, "choices", "keep"), "kept", "Kept")
	writeFile(t, filepath.Join(dir, "nested", ".gitignore"), "hidden/\n")
	writeFile(t, filepath.Join(dir, "nested", "hidden", "bad.job.spec.md"), "ignored")
	writeFile(t, filepath.Join(dir, ".git", "bad.job.spec.md"), "not a job")
	outside := t.TempDir()
	outsideFile := writeJob(t, outside, "outside", "Outside")
	if err := os.Symlink(outside, filepath.Join(dir, "linked-directory")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideFile, filepath.Join(dir, "linked-file.job.spec.md")); err != nil {
		t.Fatal(err)
	}
	nestedRepo := filepath.Join(dir, "nested-repo")
	writeJob(t, nestedRepo, "nested", "Nested")
	git(t, nestedRepo, "init", "-q")
	groups, err := List(dir)
	want := []string{"None", "  kept: Kept", "  tracked: Tracked"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("discovery: %v, %v", summary(groups), err)
	}
}

func TestDiscoveryOutsideGit(t *testing.T) {
	dir := t.TempDir()
	writeSystem(t, dir, "Root")
	writeJob(t, dir, "root", "Root job")
	writeJob(t, filepath.Join(dir, "child"), "child", "Child job")
	outside := t.TempDir()
	path := writeJob(t, outside, "outside", "Outside")
	if err := os.Symlink(outside, filepath.Join(dir, "directory-link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(path, filepath.Join(dir, "file-link.job.spec.md")); err != nil {
		t.Fatal(err)
	}
	groups, err := List(dir)
	want := []string{"Root", "  child: Child job", "  root: Root job"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("non-Git discovery: %v, %v", summary(groups), err)
	}
}

func TestRepositoryDirectoryIsNotAJob(t *testing.T) {
	root := gitRepo(t)
	dir := filepath.Join(root, "workspace.job.spec.md")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	groups, err := List(dir)
	if err != nil || len(groups) != 0 {
		t.Fatalf("directory treated as a job: %v, %v", summary(groups), err)
	}
}
