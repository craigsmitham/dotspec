package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMisplacedSpecsFailLocalCommands(t *testing.T) {
	for _, tt := range []struct {
		name, kind, filename, expected string
		commands                       [][]string
	}{
		{"adjacent system", "system", "system.spec.md", "specs/system.spec.md", [][]string{{"init"}, {"show"}, {"list"}}},
		{"system collection", "system", "specs/systems/system.spec.md", "specs/system.spec.md", [][]string{{"init"}, {"show"}, {"list"}}},
		{"nested system", "system", "specs/child/specs/system.spec.md", "specs/system.spec.md", [][]string{{"init"}, {"show"}, {"list"}}},
		{"adjacent job", "job", "work.job.spec.md", "specs/jobs/work.job.spec.md", [][]string{{"init", "work"}, {"show", "work"}, {"list"}, {"map", "work"}, {"check", "work"}}},
		{"ungrouped job", "job", "specs/work.job.spec.md", "specs/jobs/work.job.spec.md", [][]string{{"init", "work"}, {"show", "work"}, {"list"}}},
		{"singular jobs folder", "job", "specs/job/work.job.spec.md", "specs/jobs/work.job.spec.md", [][]string{{"init", "work"}, {"show", "work"}, {"list"}}},
		{"wrong collection", "job", "specs/use-cases/work.job.spec.md", "specs/jobs/work.job.spec.md", [][]string{{"init", "work"}, {"show", "work"}, {"list"}}},
		{"nested jobs", "job", "specs/jobs/group/work.job.spec.md", "specs/jobs/work.job.spec.md", [][]string{{"init", "work"}, {"show", "work"}, {"list"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFixture(t, filepath.Join(dir, tt.filename), "misplaced user content")
			for _, command := range tt.commands {
				args := append([]string{tt.kind}, command...)
				args = append(args, "--path", dir)
				code, out, errOut := invoke(args...)
				if code == 0 || out != "" || !strings.Contains(errOut, "misplaced spec") || !strings.Contains(errOut, tt.expected) {
					t.Errorf("%v: %d %q %q", args, code, out, errOut)
				}
			}
			content, err := os.ReadFile(filepath.Join(dir, tt.filename))
			if err != nil || string(content) != "misplaced user content" {
				t.Fatalf("commands changed the misplaced file: %s, %v", content, err)
			}
			if _, err := os.Stat(filepath.Join(dir, tt.expected)); !os.IsNotExist(err) {
				t.Fatalf("init must not create a second spec: %v", err)
			}
		})
	}
}

func TestDiscoveryRejectsLegacyFilesAlongsideCanonicalSpecs(t *testing.T) {
	for _, repository := range []bool{false, true} {
		for _, kind := range []string{"system", "job"} {
			t.Run(kind+map[bool]string{false: " outside Git", true: " inside Git"}[repository], func(t *testing.T) {
				dir := t.TempDir()
				if repository {
					if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
						t.Fatalf("git init: %v, %s", err, output)
					}
				}
				args := []string{kind, "init"}
				filename := "system.spec.md"
				if kind == "job" {
					args = append(args, "work")
					filename = "work.job.spec.md"
				}
				args = append(args, "--path", dir)
				if code, _, errOut := invoke(args...); code != 0 {
					t.Fatal(errOut)
				}
				writeFixture(t, filepath.Join(dir, filename), "legacy duplicate")
				if code, out, errOut := invoke(kind, "list", "--path", dir); code == 0 || out != "" || !strings.Contains(errOut, "misplaced spec") {
					t.Fatalf("mixed layout: %d %q %q", code, out, errOut)
				}
			})
		}
	}
}

func TestMetadataDirectoriesCannotBecomeProjectScopes(t *testing.T) {
	dir := t.TempDir()
	for _, subdir := range []string{"specs", "specs/jobs", "specs/jobs/extra"} {
		path := filepath.Join(dir, subdir)
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, args := range [][]string{{"system", "init"}, {"system", "show"}, {"system", "list"}, {"job", "init", "work"}, {"job", "show", "work"}, {"job", "list"}} {
			args = append(args, "--path", path)
			if code, out, errOut := invoke(args...); code == 0 || out != "" || !strings.Contains(errOut, "--path selects a project directory") {
				t.Errorf("metadata scope: %v: %d %q %q", args, code, out, errOut)
			}
		}
	}
}

func TestInitDoesNotWriteThroughMetadataSymlinks(t *testing.T) {
	for _, subdir := range []string{"specs", "specs/jobs"} {
		t.Run(subdir, func(t *testing.T) {
			dir, outside := t.TempDir(), t.TempDir()
			link := filepath.Join(dir, subdir)
			if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Fatal(err)
			}
			if code, out, errOut := invoke("job", "init", "work", "--path", dir); code == 0 || out != "" || !strings.Contains(errOut, "symlink") {
				t.Fatalf("symlink: %d %q %q", code, out, errOut)
			}
			entries, err := os.ReadDir(outside)
			if err != nil || len(entries) != 0 {
				t.Fatalf("init wrote outside the scope: %v, %v", entries, err)
			}
		})
	}
}

func TestMisplacedAssociatedSystemDoesNotFallBackToAncestor(t *testing.T) {
	dir := t.TempDir()
	if code, _, errOut := invoke("system", "init", "--path", dir); code != 0 {
		t.Fatal(errOut)
	}
	child := filepath.Join(dir, "child")
	writeFixture(t, filepath.Join(child, "system.spec.md"), "legacy child system")
	if code, _, errOut := invoke("job", "init", "work", "--path", child); code != 0 {
		t.Fatal(errOut)
	}
	for _, args := range [][]string{{"system", "show"}, {"job", "show", "work"}, {"job", "list"}} {
		args = append(args, "--path", child)
		if code, out, errOut := invoke(args...); code == 0 || out != "" || !strings.Contains(errOut, "misplaced spec") {
			t.Errorf("legacy child system: %v: %d %q %q", args, code, out, errOut)
		}
	}
}
