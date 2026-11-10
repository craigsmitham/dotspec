package cli

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func sourceSnapshot(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[path] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// Exercises the criteria in repeatable-inspection and
// inspection-preserves-source through public CLI entry points.
func TestInspectionQualityContracts(t *testing.T) {
	for _, scenario := range []struct {
		name  string
		args  []string
		setup func(*testing.T)
		code  int
	}{
		{"empty", []string{"check"}, func(t *testing.T) { prepareExample(t, "dotspec") }, 0},
		{"source inventory", []string{"sources"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"valid", []string{"check"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"advice", []string{"check"}, func(t *testing.T) {
			prepareExample(t, "read")
			writeFixture(t, "specs/quality/draft.quality.spec.md", "---\ntype: quality\n---\n# Draft\n")
		}, 0},
		{"malformed", []string{"check"}, func(t *testing.T) {
			prepareExample(t, "read")
			writeFixture(t, "specs/quality/broken.quality.spec.md", "broken")
		}, 1},
		{"shared context outside subtree", []string{"check", "--path", "internal"}, func(t *testing.T) {
			prepareExample(t, "read")
			writeFixture(t, "internal/specs/quality/a.quality.spec.md", "---\ntype: quality\n---\n# Draft A\n")
			writeFixture(t, "internal/specs/quality/b.quality.spec.md", "---\ntype: quality\n---\n# Draft B\n")
		}, 0},
		{"focused job", []string{"job", "check", "develop-requirements"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"focused use case", []string{"use-case", "check", "check-specifications"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"focused concept", []string{"concept", "check", "job-map"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"focused general requirement", []string{"requirement", "check", "reject-existing-target"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"focused quality requirement", []string{"quality", "check", "diagnostic-severity"}, func(t *testing.T) { prepareExample(t, "read") }, 0},
		{"missing input", []string{"quality", "check", "missing"}, func(t *testing.T) { prepareExample(t, "read") }, 2},
		{"unavailable context", []string{"quality", "check", "diagnostic-severity"}, func(t *testing.T) {
			prepareExample(t, "read")
			if err := os.Remove("specs/system.spec.md"); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir("specs/system.spec.md", 0o755); err != nil {
				t.Fatal(err)
			}
		}, 2},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.setup(t)
			writeFixture(t, "notes.txt", "Unrelated source material.\n")
			before := sourceSnapshot(t)
			var previous any
			for i := 0; i < 2; i++ {
				code, out, errOut := invoke(append(scenario.args, "--json")...)
				if code != scenario.code || errOut != "" {
					t.Fatalf("check: %d %s %s", code, out, errOut)
				}
				var report any
				if err := json.Unmarshal([]byte(out), &report); err != nil {
					t.Fatal(err)
				}
				if i != 0 && !reflect.DeepEqual(previous, report) {
					t.Fatalf("reports differ:\n%v\n%v", previous, report)
				}
				previous = report
				if after := sourceSnapshot(t); !reflect.DeepEqual(before, after) {
					t.Fatalf("inspection changed sources:\nbefore: %v\nafter: %v", before, after)
				}
			}
		})
	}
}

// Textual severity must survive changes in styling and redirection settings.
func TestDiagnosticSeverityQualityContract(t *testing.T) {
	prepareExample(t, "read")
	writeFixture(t, "specs/quality/draft.quality.spec.md", "---\ntype: quality\n---\n# Draft\n")
	writeFixture(t, "specs/quality/broken.quality.spec.md", "broken")
	for _, env := range []struct{ term, color, noColor string }{{"dumb", "", "1"}, {"xterm-256color", "1", ""}} {
		t.Setenv("TERM", env.term)
		t.Setenv("CLICOLOR_FORCE", env.color)
		t.Setenv("NO_COLOR", env.noColor)
		code, out, errOut := invoke("check")
		if code != 1 || errOut != "" {
			t.Fatalf("check: %d %s %s", code, out, errOut)
		}
		for _, label := range []string{"error [metadata.frontmatter]", "advice [quality.criterion]"} {
			if !strings.Contains(out, label) {
				t.Fatalf("missing %q in %s", label, out)
			}
		}
	}
}
