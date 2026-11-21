package specfile

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
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
	for _, path := range []string{"specs/new-system.system.spec.md", "apps/demo/specs/new-system.system.spec.md", "apps/other/specs/new-system.system.spec.md"} {
		if err := os.WriteFile(filepath.Join(root, path), []byte("example"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for dir, name := range map[string]string{".": "tool", "apps/demo": "demo", "apps/other": "other"} {
		writeManifest(t, filepath.Join(root, dir), "name: "+name)
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
				paths, err := Discover(filepath.Join(root, dir), func(name string) bool { return name == "new-system.system.spec.md" })
				want := []string{filepath.Join(root, dir, "specs/new-system.system.spec.md")}
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

func writeManifest(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ProjectFilename), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestProjectManifestValidation(t *testing.T) {
	for _, body := range []string{"", "{}", "name: ''", "name: invalid/name", "projects: {tool: .}", "name: demo\nunknown: true", "name: demo\ngenerated_spec: false", "name: demo\nreport: true", "name: demo\nreport:\n  enabled: perhaps", "name: demo\nreport:\n  unknown: false", "name: demo\nname: other", "name: demo\n---\nname: other"} {
		t.Run(body, func(t *testing.T) {
			root := workspaceFixture(t)
			child := filepath.Join(root, "apps/demo")
			writeManifest(t, child, body)
			if _, err := FindProject(child); err == nil || !strings.Contains(err.Error(), filepath.Join(child, ProjectFilename)) {
				t.Fatalf("expected error identifying invalid manifest, got %v", err)
			}
			if _, err := FindWorkspace(root); err == nil {
				t.Fatal("invalid manifest accepted")
			}
			paths, err := Discover(root, func(name string) bool { return name == "new-system.system.spec.md" })
			if err != nil || !reflect.DeepEqual(paths, []string{filepath.Join(root, "specs/new-system.system.spec.md")}) {
				t.Fatal(paths, err)
			}
		})
	}
}

func TestReportConfig(t *testing.T) {
	for _, tc := range []struct {
		section string
		enabled bool
	}{
		{"", true},
		{"report: {}", true},
		{"report:\n  enabled: true", true},
		{"report:\n  enabled: false", false},
	} {
		root := t.TempDir()
		writeManifest(t, root, "name: demo\n"+tc.section)
		project, err := readProject(root)
		if err != nil {
			t.Fatal(err)
		}
		if project.Report.Enabled != tc.enabled {
			t.Errorf("report.enabled for %q = %t, want %t", tc.section, project.Report.Enabled, tc.enabled)
		}
	}
}

func TestWorkspaceDiscoveryBoundaries(t *testing.T) {
	for _, git := range []bool{false, true} {
		t.Run(map[bool]string{false: "filesystem", true: "git"}[git], func(t *testing.T) {
			root := workspaceFixture(t)
			if git {
				if output, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
					t.Fatal(err, string(output))
				}
				if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("/ignored/\n"), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(filepath.Join(root, "ignored"), 0755); err != nil {
					t.Fatal(err)
				}
				writeManifest(t, filepath.Join(root, "ignored"), "name: ignored")
			}
			other := filepath.Join(root, "apps/other")
			if err := os.Mkdir(filepath.Join(other, ".git"), 0755); err != nil {
				t.Fatal(err)
			}
			writeManifest(t, filepath.Join(root, "specs"), "name: metadata")
			if err := os.Symlink(filepath.Join(root, "apps/demo"), filepath.Join(root, "alias")); err != nil {
				t.Skip(err)
			}
			w, err := FindWorkspace(filepath.Join(root, "apps/demo/internal"))
			if err != nil {
				t.Fatal(err)
			}
			want := []Project{
				{Name: "demo", Root: filepath.Join(root, "apps/demo"), Report: ReportConfig{Enabled: true}},
				{Name: "metadata", Root: filepath.Join(root, "specs"), Report: ReportConfig{Enabled: true}},
				{Name: "tool", Root: root, Report: ReportConfig{Enabled: true}},
			}
			if w.Root != root || !reflect.DeepEqual(w.Projects, want) {
				t.Fatal(w)
			}
			p, err := FindProject(filepath.Join(root, "specs"))
			if err != nil || p == nil || p.Root != filepath.Join(root, "specs") {
				t.Fatal(p, err)
			}
			writeManifest(t, filepath.Join(root, "apps/demo"), "name: tool")
			if _, err := FindWorkspace(root); err == nil || !strings.Contains(err.Error(), "duplicate project name") {
				t.Fatal(err)
			}
		})
	}
}
