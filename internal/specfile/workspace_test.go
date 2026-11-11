package specfile

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func workspaceFixture(t *testing.T) string {
	t.Helper()
	root, err := Directory(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"specs", "apps/demo/specs", "apps/other/specs", "apps/demo/internal"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"specs/system.spec.md", "apps/demo/specs/system.spec.md", "apps/other/specs/system.spec.md"} {
		if err := os.WriteFile(filepath.Join(root, path), []byte("example"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, WorkspaceFilename), []byte("projects:\n  tool: .\n  demo: apps/demo\n  other: apps/other\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return root
}
func TestWorkspacePartitionsDiscovery(t *testing.T) {
	for _, git := range []bool{false, true} {
		t.Run(map[bool]string{false: "filesystem", true: "git"}[git], func(t *testing.T) {
			root := workspaceFixture(t)
			if git {
				if output, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
					t.Fatal(err, string(output))
				}
			}
			for _, dir := range []string{".", "apps/demo", "apps/other"} {
				paths, err := Discover(filepath.Join(root, dir), func(name string) bool { return name == "system.spec.md" })
				want := []string{filepath.Join(root, dir, "specs/system.spec.md")}
				if err != nil || !reflect.DeepEqual(paths, want) {
					t.Fatal(paths, want, err)
				}
			}
			boundary, err := ContextRoot(filepath.Join(root, "apps/demo/internal"))
			if err != nil || boundary != filepath.Join(root, "apps/demo") {
				t.Fatal(boundary, err)
			}
		})
	}
}
func TestWorkspaceRejectsInvalidDeclarations(t *testing.T) {
	for _, body := range []string{
		"projects: {}", "projects:\n  one: .\n  two: .", "projects:\n  one: ../outside", "projects:\n  one: missing", "projects:\n  one: specs", "projects:\n  one: .\nunknown: true", "projects:\n  one: .\n  one: apps/demo", "projects:\n  one: .\n---\nprojects:\n  two: .",
	} {
		t.Run(body, func(t *testing.T) {
			root := workspaceFixture(t)
			if err := os.WriteFile(filepath.Join(root, WorkspaceFilename), []byte(body), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := FindWorkspace(root); err == nil {
				t.Fatal("invalid manifest accepted")
			}
		})
	}
}
func TestWorkspaceWithoutRootProjectRequiresSelection(t *testing.T) {
	root := workspaceFixture(t)
	if err := os.WriteFile(filepath.Join(root, WorkspaceFilename), []byte("projects:\n  demo: apps/demo\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ProjectContext(root); err == nil {
		t.Fatal("workspace became an implicit root project")
	}
	if _, _, err := ProjectContext(filepath.Join(root, "apps/demo")); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceRejectsSymlinksAndNestedRepositories(t *testing.T) {
	t.Run("symlink", func(t *testing.T) {
		root := workspaceFixture(t)
		if err := os.Symlink(filepath.Join(root, "apps/demo"), filepath.Join(root, "alias")); err != nil {
			t.Skip(err)
		}
		if err := os.WriteFile(filepath.Join(root, WorkspaceFilename), []byte("projects:\n  alias: alias\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := FindWorkspace(root); err == nil {
			t.Fatal("symlink root accepted")
		}
	})
	t.Run("nested repository", func(t *testing.T) {
		root := workspaceFixture(t)
		if err := os.Mkdir(filepath.Join(root, "apps/demo/.git"), 0755); err != nil {
			t.Fatal(err)
		}
		if _, err := FindWorkspace(root); err == nil {
			t.Fatal("nested repository accepted")
		}
	})
}
