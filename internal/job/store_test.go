package job

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/speccollection"
	"github.com/dotspec-md/dotspec/internal/specfile"
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
	path := filepath.Join(dir, "specs", "jobs", slug+Layout.Suffix)
	writeFile(t, path, strings.Replace(validJob, "Develop Requirements", name, 1))
	return path
}

func writeSystem(t *testing.T, dir, name string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "specs", specfile.NameSlug(name)+".system.spec.md"), "---\ntype: system-spec\n---\n# "+name+"\n## Purpose\nHelp teams.\n## Scope\nRequirements.\n")
}

func boundLookup(t *testing.T, dir string) {
	t.Helper()
	git(t, dir, "init", "-q")
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

// members lists jobs below dir as dotspec list job does.
func members(dir string) ([]*speccollection.Member, error) {
	return speccollection.Members(dir, Layout, "jobs", func(path string) (string, error) {
		job, err := Read(path)
		if err != nil {
			return "", err
		}
		return job.Name, nil
	})
}

// member is the listed job whose filename has slug.
func member(t *testing.T, dir, slug string) *speccollection.Member {
	t.Helper()
	listed, err := members(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range listed {
		if filepath.Base(m.Path) == slug+Layout.Suffix {
			return m
		}
	}
	t.Fatalf("job %s not listed below %s", slug, dir)
	return nil
}

func summary(listed []*speccollection.Member) []string {
	var result []string
	for i, m := range listed {
		if i == 0 || (m.System == nil) != (listed[i-1].System == nil) || (m.System != nil && m.System.Path != listed[i-1].System.Path) {
			name := "None"
			if m.System != nil {
				name = m.System.Name
			}
			result = append(result, name)
		}
		result = append(result, "  "+strings.TrimSuffix(filepath.Base(m.Path), Layout.Suffix)+": "+m.Name)
	}
	return result
}

func TestInitAndIdentity(t *testing.T) {
	dir := t.TempDir()
	boundLookup(t, dir)
	path, err := Init(dir, "develop-requirements")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Read(path)
	if err != nil || got.Slug != "develop-requirements" || member(t, dir, "develop-requirements").System != nil {
		t.Fatalf("unexpected job: %+v, %v", got, err)
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
	got, err = Read(path)
	if err != nil || got.Name != "A revised title" || got.Slug != "develop-requirements" {
		t.Fatalf("identity changed: %+v, %v", got, err)
	}
}

func TestInvalidSlugsAndPaths(t *testing.T) {
	dir := t.TempDir()
	for _, slug := range []string{"", "../escape", "/escape", "a/b", `a\b`, "Uppercase", "two words", "under_score", "-leading", "trailing-", "two--hyphens", "name.job.spec.md", "café"} {
		if _, err := Init(dir, slug); err == nil {
			t.Errorf("init accepted %q", slug)
		}
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid init created files: %v, %v", entries, err)
	}
	for _, slug := range []string{"develop-requirements", "job2", "2-steps"} {
		if err := Layout.ValidateSlug(slug); err != nil {
			t.Fatal(err)
		}
	}
	file := writeJob(t, dir, "valid", "Valid")
	for _, path := range []string{filepath.Join(dir, "missing"), file} {
		if _, err := Init(path, "test"); err == nil {
			t.Fatal("init accepted an invalid directory")
		}
		if _, err := members(path); err == nil {
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
		if got := member(t, tt.path, tt.slug); got.System == nil || got.System.Name != tt.want {
			t.Fatalf("association: %+v", got)
		}
	}
	groups, err := members(dir)
	want := []string{"Parser", "  develop-requirements: Develop Parser Requirements", "  read-token: Read a token", "dotspec", "  develop-requirements: Develop Requirements", "  manage-requirements: Manage Requirements"}
	if err != nil || !reflect.DeepEqual(summary(groups), want) {
		t.Fatalf("got %v, %v; want %v", summary(groups), err, want)
	}
	groups, err = members(deep)
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
	groups, err := members(dir)
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
			git(t, dir, "init", "-q", "--separate-git-dir="+filepath.Join(outer, "worktree-git"))
		} else {
			boundLookup(t, dir)
		}
		if got := member(t, dir, "work"); got.System != nil {
			t.Fatalf("crossed boundary: %+v", got)
		}
	}
}

func TestMalformedJobsAndAssociatedSystems(t *testing.T) {
	dir := gitRepo(t)
	groups, err := members(dir)
	if err != nil || len(groups) != 0 {
		t.Fatalf("empty list: %v, %v", groups, err)
	}
	path := writeJob(t, dir, "work", "Work")
	writeFile(t, path, "malformed")
	if _, err := members(dir); err == nil || !strings.Contains(err.Error(), "work.job.spec.md") {
		t.Fatalf("malformed list: %v", err)
	}
	writeJob(t, dir, "work", "Work")
	writeFile(t, filepath.Join(dir, "specs", "new-system.system.spec.md"), "malformed")
	if _, err := members(dir); err == nil || !strings.Contains(err.Error(), "new-system.system.spec.md") {
		t.Fatalf("malformed group: %v", err)
	}
	writeSystem(t, dir, "Valid")
	writeFile(t, filepath.Join(dir, "specs", "jobs", "Invalid_slug.job.spec.md"), validJob)
	if _, err := members(dir); err == nil || !strings.Contains(err.Error(), "invalid job slug") {
		t.Fatalf("invalid filename: %v", err)
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
	groups, err := members(dir)
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
	groups, err := members(dir)
	if err != nil || len(groups) != 0 {
		t.Fatalf("directory treated as a job: %v, %v", summary(groups), err)
	}
}
