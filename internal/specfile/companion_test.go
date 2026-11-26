package specfile

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestCompanionFrontmatter(t *testing.T) {
	spec := func(frontmatter string) []byte {
		return []byte("---\ntype: requirement-spec\n" + frontmatter + "---\n# Keep\n## Statement\nKeep it.\n")
	}
	outline, err := ParseOutline(spec("companions:\n  - target: ../models/tasks.qnt#exactlyOneView\n    role: formalizes\n    note: Task lifecycle\n  - {target: https://example.com/doc, role: explains}\n"), "requirement", &Metadata{})
	if err != nil {
		t.Fatal(err)
	}
	want := []Companion{
		{Target: "../models/tasks.qnt#exactlyOneView", Role: "formalizes", Note: "Task lifecycle", Location: diagnostic.Location{Line: 4, Column: 13, Element: CompanionsField}},
		{Target: "https://example.com/doc", Role: "explains", Location: diagnostic.Location{Line: 7, Column: 14, Element: CompanionsField}},
	}
	if !reflect.DeepEqual(outline.Companions, want) {
		t.Fatalf("%#v", outline.Companions)
	}
	links := outline.References().Links
	if len(links) < 2 || links[0].Destination != want[0].Target || links[0].Location != want[0].Location || links[1].Destination != want[1].Target {
		t.Fatalf("companion targets must be ordinary references: %#v", links)
	}

	for _, c := range []struct{ frontmatter, rule, message string }{
		{"companions: models/tasks.qnt\n", "companion.invalid", "YAML sequence"},
		{"companions:\n  - models/tasks.qnt\n", "companion.invalid", "mapping"},
		{"companions:\n  - role: formalizes\n", "companion.invalid", "target"},
		{"companions:\n  - target: x.qnt\n", "companion.role", "Give each companion a role"},
		{"companions:\n  - target: x.qnt\n    role: proves\n", "companion.role", `"proves"`},
		{"companions:\n  - target: x.qnt\n    role: formalizes\n    kind: quint\n", "companion.invalid", `"kind"`},
		{"companions:\n  - target: '#statement'\n    role: explains\n", "companion.invalid", "another file"},
		{"companions:\n  - target: x.qnt\n    role: formalizes\n    note: ''\n", "companion.invalid", "empty note"},
		{"companions:\n  - target: x.qnt\n    role: formalizes\n    note: |\n      two\n      lines\n", "companion.invalid", "single-line"},
		{"companions:\n  - target: [x.qnt]\n    role: formalizes\n", "companion.invalid", "single-line text"},
		{"companions:\n  - target: x.qnt\n    role: verifies\n  - target: x.qnt\n    role: verifies\n", "companion.duplicate", "already lists"},
	} {
		_, err := ParseOutline(spec(c.frontmatter), "requirement", &Metadata{})
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != c.rule || !strings.Contains(d.Message, c.message) || d.Location.Element != CompanionsField {
			t.Errorf("%q: %v", c.frontmatter, err)
		}
	}
	// The same target can serve two roles, such as a model that also checks itself.
	if _, err := ParseOutline(spec("companions:\n  - target: x.qnt\n    role: formalizes\n  - target: x.qnt\n    role: verifies\n"), "requirement", &Metadata{}); err != nil {
		t.Fatal(err)
	}
	model := []byte("---\ntype: quality-model\ncompanions:\n  - target: x.md\n    role: explains\n---\n# Model\n")
	if _, err := ParseOutline(model, "quality-model", &Metadata{}); err == nil || !strings.Contains(err.Error(), "reference material") {
		t.Fatal(err)
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
	if line := NativeDefinitionLine("x.QNT", []byte("module x {\r\n  val inv = true\r\n}"), "inv"); line != 2 {
		t.Fatal(line)
	}
}
