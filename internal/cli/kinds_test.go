package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every kind in the table is created with new and found by list through one
// convention: new refuses an existing target and names a missing operand, and
// list returns the shared envelope with each entry's system.
func TestKindTableContracts(t *testing.T) {
	operands := map[string][]string{
		"system":            {"Billing", "-C", "billing", "--part-of", "../specs/new-system.system.spec.md"},
		"component":         {"Intake", "-C", "internal", "--part-of", "../specs/new-system.system.spec.md"},
		"quality-attribute": {"durability"},
	}
	for _, k := range kinds() {
		t.Run(k.layout.Kind, func(t *testing.T) {
			prepareExample(t, "dotspec new")
			if err := os.Mkdir("billing", 0o755); err != nil {
				t.Fatal(err)
			}
			args, ok := operands[k.layout.Kind]
			if !ok {
				args = []string{"work"}
			}
			if k.layout.Kind != "system" {
				code, out, errOut := invoke("list", k.layout.Kind)
				if code != 0 || out != "No "+k.plural+" found in the enclosing project.\n" || errOut != "" {
					t.Fatalf("empty list: %d %q %q", code, out, errOut)
				}
			}
			code, out, errOut := invoke(append([]string{"new", k.layout.Kind}, args...)...)
			if code != 0 || errOut != "" || !strings.HasPrefix(out, "Created ") || !strings.Contains(out, "Authoring guide: dotspec docs guides/") {
				t.Fatalf("new: %d %q %q", code, out, errOut)
			}
			if code, out, errOut := invoke(append([]string{"new", k.layout.Kind}, args...)...); code != 2 || out != "" || !strings.Contains(errOut, "refusing to overwrite") {
				t.Fatalf("existing target: %d %q %q", code, out, errOut)
			}
			missing := strings.TrimSuffix(k.operand, "...")
			if code, out, errOut := invoke("new", k.layout.Kind); code != 2 || out != "" || !strings.Contains(errOut, "missing "+missing) {
				t.Fatalf("missing operand: %d %q %q", code, out, errOut)
			}
			code, out, errOut = invoke("list", k.layout.Kind, "--json")
			if code != 0 || errOut != "" {
				t.Fatalf("list: %d %q %q", code, out, errOut)
			}
			assertResultContract(t, []string{"list", k.layout.Kind}, out)
			var result listResult
			if err := json.Unmarshal([]byte(out), &result); err != nil {
				t.Fatal(err)
			}
			slug := strings.ToLower(args[0])
			var created *listEntry
			for i, e := range result.Entries {
				if e.Slug == slug {
					created = &result.Entries[i]
				}
			}
			if created == nil || created.Kind != k.layout.Kind || !filepath.IsAbs(created.Path) || created.System == nil || created.System.Name != "New system" {
				t.Fatalf("entry for %s: %+v", slug, result.Entries)
			}
			code, out, errOut = invoke("list", k.layout.Kind)
			if code != 0 || errOut != "" || !strings.Contains(out, created.Name+"  (") {
				t.Fatalf("list text: %d %q %q", code, out, errOut)
			}
			if err := os.WriteFile(created.Path, []byte("broken"), 0o644); err != nil {
				t.Fatal(err)
			}
			if code, out, errOut := invoke("list", k.layout.Kind); code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
				t.Fatalf("malformed: %d %q %q", code, out, errOut)
			}
		})
	}
}

func TestStructuralLists(t *testing.T) {
	prepareExample(t, "dotspec new")
	for _, args := range [][]string{
		{"-C", "internal", "new", "component", "Intake", "--part-of", "../specs/new-system.system.spec.md"},
		{"-C", "internal/parser", "new", "component", "Parser", "--part-of", "../specs/intake.component.spec.md"},
		{"-C", "billing", "new", "system", "Billing", "--part-of", "../specs/new-system.system.spec.md"},
	} {
		if err := os.MkdirAll(args[1], 0o755); err != nil {
			t.Fatal(err)
		}
		if code, out, errOut := invoke(args...); code != 0 {
			t.Fatalf("%v: %d %q %q", args, code, out, errOut)
		}
	}
	data, err := os.ReadFile("billing/specs/billing.system.spec.md")
	if err != nil || !strings.Contains(string(data), "## Part of\n\n[Parent specification](<../../specs/new-system.system.spec.md>)") {
		t.Fatalf("nested system: %s %v", data, err)
	}
	for _, tc := range []struct{ kind, want string }{
		{"system", "Systems by directory\n└── New system  (specs/new-system.system.spec.md)\n    └── Billing  (billing/specs/billing.system.spec.md)\n"},
		{"component", "Components by Part of\n\nNew system  (specs/new-system.system.spec.md)\n└── Intake  (internal/specs/intake.component.spec.md)\n    └── Parser  (internal/parser/specs/parser.component.spec.md)\n"},
	} {
		if code, out, errOut := invoke("list", tc.kind); code != 0 || out != tc.want || errOut != "" {
			t.Errorf("list %s: %d %q %q", tc.kind, code, out, errOut)
		}
	}
	if code, out, errOut := invoke("new", "system", "Orphan"); code != 2 || out != "" || !strings.Contains(errOut, "missing --part-of <path>") {
		t.Fatalf("root system from new: %d %q %q", code, out, errOut)
	}
	if code, _, errOut := invoke("-C", "billing", "new", "system", "Inner", "--part-of", "../internal/specs/intake.component.spec.md"); code != 2 || !strings.Contains(errOut, "expected a system declaration") {
		t.Fatalf("system part of component: %d %q", code, errOut)
	}
}

func TestKindOptions(t *testing.T) {
	prepareExample(t, "dotspec new")
	if code, out, errOut := invoke("new", "use-case", "sketch"); code != 0 || errOut != "" {
		t.Fatalf("brief: %d %s %s", code, out, errOut)
	}
	data, err := os.ReadFile("specs/use-cases/sketch.use-case.spec.md")
	if err != nil || strings.Contains(string(data), "## Preconditions") || !strings.Contains(string(data), "## System boundary") {
		t.Fatalf("brief: %s %v", data, err)
	}
	if code, _, _ := invoke("new", "use-case", "long", "--template", "expanded"); code != 0 {
		t.Fatal("expanded template")
	}
	if data, err := os.ReadFile("specs/use-cases/long.use-case.spec.md"); err != nil || !strings.Contains(string(data), "## Preconditions") {
		t.Fatalf("expanded: %s %v", data, err)
	}
	for _, args := range [][]string{
		{"new", "use-case", "unused", "--template", "unknown"},
		{"new", "requirement", "unused", "--attribute", ""},
		{"new", "requirement", "unused", "--attribute", "absent"},
		{"new", "quality-model", "unused", "--from", ""},
		{"new", "quality-attribute", "unused", "--from", ""},
		{"new", "requirement", "Upper"}, {"new", "requirement", "../escape"},
		{"new", "requirement", "one", "two"}, {"new", "requirement", "unused", "--json=text"},
	} {
		if code, out, errOut := invoke(args...); code != 2 || out != "" || errOut == "" {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, path := range []string{"specs/use-cases/unused.use-case.spec.md", "specs/requirements/unused.requirement.spec.md", "specs/quality-models/unused.quality-model.md"} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("invalid invocation created %s", path)
		}
	}
}

// Adoption copies a term and its missing parents, reuses parents the system
// already declares and refuses an attribute it already declares.
func TestQualityAttributeAdoption(t *testing.T) {
	prepareExample(t, "dotspec new")
	code, out, errOut := invoke("new", "quality-attribute", "recoverability", "--from", "iso-25010-2023", "--dry-run")
	if code != 0 || errOut != "" || !strings.Contains(out, "Would create reliability") || !strings.Contains(out, "Would create recoverability") {
		t.Fatalf("preview: %d %q %q", code, out, errOut)
	}
	if _, err := os.Stat("specs/quality"); !os.IsNotExist(err) {
		t.Fatal("preview wrote files")
	}
	if code, out, errOut := invoke("new", "quality-attribute", "recoverability", "--from", "iso-25010-2023"); code != 0 || !strings.Contains(out, "Added 2 attributes") {
		t.Fatalf("adopt: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("new", "quality-attribute", "fault-tolerance", "--from", "iso-25010-2023", "--json")
	var result struct {
		Added   []string `json:"added"`
		Entries []struct{ ID, Decision string }
	}
	if err := json.Unmarshal([]byte(out), &result); code != 0 || err != nil || strings.Join(result.Added, ",") != "fault-tolerance" || result.Entries[0].Decision != "reuse" {
		t.Fatalf("reuse parent: %d %q %q %v", code, out, errOut, err)
	}
	if code, out, errOut := invoke("new", "quality-attribute", "reliability", "--from", "iso-25010-2023"); code != 2 || out != "" || !strings.Contains(errOut, "already exists") {
		t.Fatalf("existing attribute: %d %q %q", code, out, errOut)
	}
	if code, out, errOut := invoke("new", "quality-attribute", "speed", "latency"); code != 0 || !strings.Contains(out, "Added 2 attributes") {
		t.Fatalf("blank: %d %q %q", code, out, errOut)
	}
	if data, err := os.ReadFile("specs/quality/latency/latency.quality-attribute.spec.md"); err != nil || !strings.Contains(string(data), "# Latency\n\n## Definition") {
		t.Fatalf("blank attribute: %s %v", data, err)
	}
}
