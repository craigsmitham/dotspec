package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSpec(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, "specs", Filename)
	writeFile(t, path, strings.Replace(validSpec, "# Parser", "# "+name, 1))
	return path
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func isolatedGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Fatal("Git is required for repository discovery tests")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, name := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR"} {
		t.Setenv(name, "")
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInitNeverOverwrites(t *testing.T) {
	dir := t.TempDir()
	path, err := Init(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err != nil {
		t.Fatalf("template must be a valid system: %v", err)
	}
	writeFile(t, path, "existing user content")
	if _, err := Init(dir); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("expected overwrite error, got %v", err)
	}
	contents, _ := os.ReadFile(path)
	if string(contents) != "existing user content" {
		t.Fatal("existing content changed")
	}
}

func TestShowResolvesNearestSystemAndParent(t *testing.T) {
	isolatedGit(t)
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	rootPath := writeSpec(t, dir, "Root")
	parser := filepath.Join(dir, "internal", "parser")
	parserPath := writeSpec(t, parser, "Parser")
	tokens := filepath.Join(parser, "tokens")
	if err := os.MkdirAll(tokens, 0o755); err != nil {
		t.Fatal(err)
	}
	spec, parent, err := Show(tokens)
	if err != nil || spec.Name != "Parser" || parent.Name != "Root" {
		t.Fatalf("show: spec=%+v parent=%+v err=%v", spec, parent, err)
	}
	// macOS temporary directories can have symlinked ancestors.
	wantParser, _ := filepath.EvalSymlinks(parserPath)
	wantRoot, _ := filepath.EvalSymlinks(rootPath)
	if spec.Path != wantParser || parent.Path != wantRoot {
		t.Fatalf("wrong sources: %s, %s", spec.Path, parent.Path)
	}
	spec, parent, err = Show(dir)
	if err != nil || spec.Name != "Root" || parent != nil {
		t.Fatalf("root show: spec=%+v parent=%+v err=%v", spec, parent, err)
	}
}

func TestShowHonorsRepositoryAndWorktreeBoundaries(t *testing.T) {
	isolatedGit(t)
	outer := t.TempDir()
	writeSpec(t, outer, "Outside")
	repo := filepath.Join(outer, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "init", "-q")
	git(t, repo, "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "-c", "commit.gpgsign=false", "commit", "--allow-empty", "-qm", "Initial")
	worktree := filepath.Join(outer, "worktree")
	git(t, repo, "worktree", "add", "--detach", worktree)
	for _, dir := range []string{repo, worktree} {
		if _, _, err := Show(dir); err == nil || !strings.Contains(err.Error(), "no system spec") {
			t.Fatalf("lookup must stop at %s: %v", dir, err)
		}
		writeSpec(t, dir, "Inside")
		if spec, parent, err := Show(dir); err != nil || spec.Name != "Inside" || parent != nil {
			t.Fatalf("boundary show: spec=%+v parent=%+v err=%v", spec, parent, err)
		}
	}
}

func TestShowReportsMalformedAncestor(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "specs", Filename), "not a spec")
	child := filepath.Join(dir, "child")
	writeSpec(t, child, "Child")
	if _, _, err := Show(child); err == nil || !strings.Contains(err.Error(), Filename) {
		t.Fatalf("expected source-specific parent error, got %v", err)
	}
}

func treeNames(nodes []*Node) []string {
	var names []string
	var visit func([]*Node, string)
	visit = func(nodes []*Node, prefix string) {
		for _, node := range nodes {
			names = append(names, prefix+node.System.Name)
			visit(node.Children, prefix+"  ")
		}
	}
	visit(nodes, "")
	return names
}

func TestListHierarchyAndSubtree(t *testing.T) {
	dir := t.TempDir()
	writeSpec(t, dir, "Root")
	writeSpec(t, filepath.Join(dir, "a"), "A")
	writeSpec(t, filepath.Join(dir, "a", "gap", "child"), "Child")
	writeSpec(t, filepath.Join(dir, "b"), "B")
	// Unrelated types must not be interpreted as systems.
	writeFile(t, filepath.Join(dir, "work.job.spec.md"), "not a system")
	for _, tt := range []struct {
		path string
		want []string
	}{
		{dir, []string{"Root", "  A", "    Child", "  B"}},
		{filepath.Join(dir, "a"), []string{"A", "  Child"}},
		{filepath.Join(dir, "a", "gap"), []string{"Child"}},
	} {
		nodes, err := List(tt.path)
		if err != nil || !reflect.DeepEqual(treeNames(nodes), tt.want) {
			t.Fatalf("list %s: got %v, %v; want %v", tt.path, treeNames(nodes), err, tt.want)
		}
	}
}

func TestListGitIgnoresAndSymlinks(t *testing.T) {
	isolatedGit(t)
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	writeSpec(t, dir, "Root")
	writeSpec(t, filepath.Join(dir, "tracked"), "Tracked")
	writeSpec(t, filepath.Join(dir, "deleted"), "Deleted")
	git(t, dir, "add", ".")
	if err := os.Remove(filepath.Join(dir, "deleted", "specs", Filename)); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, ".gitignore"), "ignored/\ntracked/\nchoices/*\n!choices/keep/\n")
	writeFile(t, filepath.Join(dir, "ignored", "specs", Filename), "malformed but ignored")
	writeSpec(t, filepath.Join(dir, "choices", "keep"), "Kept")
	writeFile(t, filepath.Join(dir, "choices", "omit", "specs", Filename), "ignored")
	writeFile(t, filepath.Join(dir, "nested", ".gitignore"), "hidden/\n")
	writeFile(t, filepath.Join(dir, "nested", "hidden", "specs", Filename), "ignored")
	writeSpec(t, filepath.Join(dir, "nested", "visible"), "Visible")
	writeFile(t, filepath.Join(dir, ".git", "info", "exclude"), "excluded/\n")
	writeFile(t, filepath.Join(dir, "excluded", "specs", Filename), "ignored")
	writeFile(t, filepath.Join(dir, ".git", "metadata", "specs", Filename), "not a spec")
	outside := t.TempDir()
	writeSpec(t, outside, "Outside")
	if err := os.Symlink(outside, filepath.Join(dir, "linked-directory")); err != nil {
		t.Fatal(err)
	}
	linkedFile := filepath.Join(dir, "linked-file")
	if err := os.MkdirAll(filepath.Join(linkedFile, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "specs", Filename), filepath.Join(linkedFile, "specs", Filename)); err != nil {
		t.Fatal(err)
	}
	// A tracked directory replaced with a symlink must not expose outside files.
	writeSpec(t, filepath.Join(dir, "replaced"), "Replaced")
	git(t, dir, "add", "replaced")
	if err := os.RemoveAll(filepath.Join(dir, "replaced")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "replaced")); err != nil {
		t.Fatal(err)
	}
	nodes, err := List(dir)
	want := []string{"Root", "  Kept", "  Visible", "  Tracked"}
	if err != nil || !reflect.DeepEqual(treeNames(nodes), want) {
		t.Fatalf("got %v, %v; want %v", treeNames(nodes), err, want)
	}
	nodes, err = List(filepath.Join(dir, "nested"))
	if err != nil || !reflect.DeepEqual(treeNames(nodes), []string{"Visible"}) {
		t.Fatalf("scoped Git discovery: %v, %v", treeNames(nodes), err)
	}
}

func TestListOutsideGitSkipsLinksAndNestedRepositories(t *testing.T) {
	dir := t.TempDir()
	writeSpec(t, dir, "Root")
	outside := t.TempDir()
	writeSpec(t, outside, "Outside")
	if err := os.Symlink(outside, filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "nested", ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSpec(t, filepath.Join(dir, "nested"), "Nested repository")
	nodes, err := List(dir)
	if err != nil || !reflect.DeepEqual(treeNames(nodes), []string{"Root"}) {
		t.Fatalf("got %v, %v", treeNames(nodes), err)
	}
}

func TestListEmptyAndMalformed(t *testing.T) {
	dir := t.TempDir()
	nodes, err := List(dir)
	if err != nil || len(nodes) != 0 {
		t.Fatalf("empty list: %v, %v", nodes, err)
	}
	path := filepath.Join(dir, "child", "specs", Filename)
	writeFile(t, path, "bad spec")
	if _, err := List(dir); err == nil || !strings.Contains(err.Error(), Filename) {
		t.Fatalf("expected source-specific parse error, got %v", err)
	}
}

func TestCommandsRejectInvalidPaths(t *testing.T) {
	dir := t.TempDir()
	file := writeSpec(t, dir, "Root")
	for _, path := range []string{filepath.Join(dir, "missing"), file} {
		if _, err := Init(path); err == nil {
			t.Fatal("init accepted invalid directory")
		}
		if _, _, err := Show(path); err == nil {
			t.Fatal("show accepted invalid directory")
		}
		if _, err := List(path); err == nil {
			t.Fatal("list accepted invalid directory")
		}
	}
}
