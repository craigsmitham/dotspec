package specfile

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
)

func TestCatalogFrontmatter(t *testing.T) {
	for _, c := range []struct{ field, rule string }{
		{"specifications: never", "catalog.invalid"},
		{"verification: [broken]", "catalog.invalid"},
		{"specifications: [{name: x, source: x}]", "catalog.invalid"},
		{"specifications: [{name: x, source: x, defines: gherkin}]", "catalog.invalid"},
		{"verification: [{name: x, source: '#statement'}]", "catalog.invalid"},
		{"verification: [{name: x, source: []}]", "catalog.invalid"},
		{"verification: [{name: x, source: a, role: verifies}]", "catalog.invalid"},
		{"verification: [{name: x, source: a}, {name: X, source: b}]", "catalog.duplicate"},
		{"companions: []", "metadata.invalid"},
	} {
		_, err := NewMarkdown([]byte("---\ntype: requirement-spec\n"+c.field+"\n---\n# Keep\n")).TypedOutline("requirement", &Metadata{})
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != c.rule || d.Location.Line != 3 {
			t.Errorf("%q: %v", c.field, err)
		}
	}
}

func TestNativeDefinitions(t *testing.T) {
	for _, c := range []struct {
		file, source string
		want         []NativeDefinition
	}{
		{"tasks.qnt", "module tasks {\n  const N: int\n  var views: str -> str\n  pure def size(s) = 0\n  val exactlyOneView = true\n  action step = all {}\n  temporal deletedStaysAbsent = true\n  run resume = init\n}\n",
			[]NativeDefinition{{"tasks", 1}, {"N", 2}, {"views", 3}, {"size", 4}, {"exactlyOneView", 5}, {"step", 6}, {"deletedStaysAbsent", 7}, {"resume", 8}}},
		{"balls.als", "module balls\nabstract sig Node, Edge {}\nvar sig Ball {}\nfact connected {}\npred Ball.move[n: Node] {}\nassert NoCollision {}\ncheck NoCollision for 5\n",
			[]NativeDefinition{{"balls", 1}, {"Node", 2}, {"Edge", 2}, {"Ball", 3}, {"connected", 4}, {"move", 5}, {"NoCollision", 6}}},
		{"Commit.tla", "---- MODULE Commit ----\nCONSTANTS RM, Max\nVARIABLES state\nInit == state = 0\nNext(x) == TRUE\n  Nested == FALSE\nTHEOREM Safe == Spec => []TypeOK\n====\n",
			[]NativeDefinition{{"Commit", 1}, {"RM", 2}, {"Max", 2}, {"state", 3}, {"Init", 4}, {"Next", 5}, {"Safe", 7}}},
		{"views.feature", "Feature: Task views\n  Scenario: Resume a waiting entry\n  Scenario Outline: Resume a waiting entry\n  Rule: One view\n",
			[]NativeDefinition{{"task-views", 1}, {"resume-a-waiting-entry", 2}, {"resume-a-waiting-entry-1", 3}, {"one-view", 4}}},
		{"store_test.go", "package store\nfunc TestRestart(t *testing.T) {}\nfunc (s *Store) Save() {}\ntype Store struct{}\n",
			[]NativeDefinition{{"TestRestart", 2}, {"Save", 3}, {"Store", 4}}},
		{"acceptance.py", "class Page:\n    def entries(self):\n        pass\nasync def verify(binary):\n    pass\n",
			[]NativeDefinition{{"Page", 1}, {"entries", 2}, {"verify", 4}}},
		{"workflow.txtar", "exec dotspec\n-- input.md --\n# Input\n--  b.txt  --\n",
			[]NativeDefinition{{"input.md", 2}, {"b.txt", 4}}},
		{"queue.tla", "VARIABLES\n  pc,\n  stack\nCONSTANTS q, sent,\n          head\nInit == TRUE\nApply(F(_), x) == F(x)\nLOCAL Helper == 1\n",
			[]NativeDefinition{{"pc", 2}, {"stack", 3}, {"q", 4}, {"sent", 4}, {"head", 5}, {"Init", 6}, {"Apply", 7}, {"Helper", 8}}},
		{"grouped.go", "package x\ntype (\n\tFoo struct{}\n\tBar = int\n)\nfunc After() {}\n",
			[]NativeDefinition{{"Foo", 3}, {"Bar", 4}, {"After", 6}}},
		{"birthday.als", "module examples/toys/birthday\nsig Name {}\n",
			[]NativeDefinition{{"birthday", 1}, {"Name", 2}}},
	} {
		if got := NativeDefinitions(c.file, []byte(c.source)); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %v", c.file, got)
		}
	}
	if NativeDefinitionFormat("diagram.puml") || NativeDefinitions("diagram.puml", []byte("@startuml")) != nil {
		t.Fatal("unsupported formats keep their own fragment semantics")
	}
	if line := nativeDefinitionLine("x.QNT", []byte("module x {\r\n  val inv = true\r\n}"), "inv"); line != 2 {
		t.Fatal(line)
	}
	file := filepath.Join(t.TempDir(), "x.qnt")
	if err := os.WriteFile(file, []byte("module x {\n  val inv = true\n}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if line, err := NativeDefinitionAt(file, "inv"); err != nil || line != 2 {
		t.Fatal(line, err)
	}
	if line, err := NativeDefinitionAt(file, "missing"); err != nil || line != 0 {
		t.Fatal(line, err)
	}
	if _, err := NativeDefinitionAt(filepath.Join(filepath.Dir(file), "absent.qnt"), "inv"); err == nil {
		t.Fatal("missing file read")
	}
}
