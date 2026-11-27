package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/testproject"
)

func compositionFixture(t *testing.T) string {
	t.Helper()
	return testproject.Open(t, "concept-design")
}
func designRecord(t *testing.T, s *Snapshot, kind, name string) Record {
	t.Helper()
	d := newDesignResolver(s)
	for _, r := range s.Elements {
		if String(r, "kind") == kind && d.name(r) == name {
			return r
		}
	}
	t.Fatalf("missing %s %s", kind, name)
	return nil
}
func TestConceptCompositionExport(t *testing.T) {
	root := compositionFixture(t)
	result := conformanceExport(t, root)
	s := result.Snapshot
	if result.ExitCode != 0 || len(s.Diagnostics) != 0 {
		t.Fatalf("unexpected findings: %d %v", result.ExitCode, s.Diagnostics)
	}
	checkSourceText(t, s)
	public := designRecord(t, s, "ConceptUse", "PublicLinks")
	internal := designRecord(t, s, "ConceptUse", "InternalLinks")
	if public["definition"] != internal["definition"] || public["id"] == internal["id"] {
		t.Fatal("shared definition conflated named uses")
	}
	shorthand := designRecord(t, s, "Concept", "Shorthand")
	state := newDesignResolver(s).text(shorthand["abstract_state"])
	for _, field := range []string{"Initial state", "Invariants", "Derivation"} {
		if !strings.Contains(state, field) {
			t.Fatal("state passage lost", field)
		}
	}
	expire := designRecord(t, s, "ConceptAction", "Expire")
	if expire["initiation"] == nil || expire["outputs"] == nil {
		t.Fatal("system initiation or action signature lost")
	}
	for _, r := range s.Elements {
		if String(r, "kind") == "ConceptUse" && String(r, "owner") == "" {
			t.Fatal("use has no system owner")
		}
	}
	data, _ := json.Marshal(s)
	second := conformanceExport(t, root)
	again, _ := json.Marshal(second.Snapshot)
	if string(data) != string(again) {
		t.Fatal("composition export is nondeterministic")
	}

}
func TestCompositionAuthoringErrorsRemainFaithful(t *testing.T) {
	cases := []struct{ name, old, replacement, code string }{
		{"wrong definition kind", "[Shorthand](concepts/shorthand.concept.spec.md).", "[System](expiring-links.system.spec.md).", "concept.definition"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := compositionFixture(t)
			path := filepath.Join(root, "specs/expiring-links.system.spec.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), tc.old) {
				t.Fatal("fixture no longer exercises case")
			}
			if err = os.WriteFile(path, []byte(strings.Replace(string(data), tc.old, tc.replacement, 1)), 0644); err != nil {
				t.Fatal(err)
			}
			result := conformanceExport(t, root)
			if result.ExitCode != 1 {
				t.Fatalf("exit %d, want authoring error", result.ExitCode)
			}
			found := false
			for _, d := range result.Snapshot.Diagnostics {
				found = found || String(d, "code") == tc.code && String(d, "severity") == "error"
			}
			if !found {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			checkSourceText(t, result.Snapshot)
		})
	}
}
func TestCompositionCyclesAndPartialDeclarations(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/concepts/simple.concept.spec.md", "---\ntype: concept-spec\n---\n# Simple\n## Purpose\nServe a purpose.\n## Actions\n### Run\nPerform behavior.\n")
	conformanceWrite(t, root, "specs/example.system.spec.md", `---
type: system-spec
---
# Example
## Purpose
P.
## Scope
S.
## Concept uses
### First
#### Definition
[Simple](concepts/simple.concept.spec.md)
#### Dependencies
[Second](#second) is needed in this application.
### Second
#### Definition
[Simple](concepts/simple.concept.spec.md)
#### Dependencies
[First](#first) is needed in this application.
### External draft
#### Definition
[External](https://example.test/concept)
## Synchronizations
### Draft rule
#### When
A condition to investigate.
#### Then
Handling remains undecided.
#### Semantics
A supporting note, not an action.
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	first := designRecord(t, result.Snapshot, "ConceptUse", "First")
	second := designRecord(t, result.Snapshot, "ConceptUse", "Second")
	if Strings(first["dependencies"])[0] != second["id"] || Strings(second["dependencies"])[0] != first["id"] {
		t.Fatal("mutual dependencies lost")
	}
	draft := designRecord(t, result.Snapshot, "ConceptUse", "External draft")
	if draft["definition"] != nil {
		t.Fatal("external definition invented")
	}
	rule := designRecord(t, result.Snapshot, "Synchronization", "Draft rule")
	if !strings.Contains(newDesignResolver(result.Snapshot).text(rule["semantics"]), "supporting note") {
		t.Fatal("supporting notes lost")
	}
}
func TestCompositionRejectsForeignSystemUsesAndLocalNameCollisions(t *testing.T) {
	t.Run("foreign system use", func(t *testing.T) {
		root := compositionFixture(t)
		conformanceWrite(t, root, "other/specs/other.system.spec.md", `---
type: system-spec
---
# Other
## Purpose
P.
## Scope
S.
## Concept uses
### ForeignLinks
#### Definition
[Shorthand](../../specs/concepts/shorthand.concept.spec.md)
#### Dependencies
[PublicLinks](../../specs/expiring-links.system.spec.md#publiclinks)
`)
		result := conformanceExport(t, root)
		if result.ExitCode != 1 {
			t.Fatal("cross-system dependency accepted")
		}
		foreign := designRecord(t, result.Snapshot, "ConceptUse", "ForeignLinks")
		if len(Strings(foreign["dependencies"])) != 0 {
			t.Fatal("foreign dependency materialized")
		}
		found := false
		for _, d := range result.Snapshot.Diagnostics {
			found = found || String(d, "code") == "concept.dependencies" && String(d, "severity") == "error"
		}
		if !found {
			t.Fatal(result.Snapshot.Diagnostics)
		}
	})
	t.Run("duplicate names", func(t *testing.T) {
		root := t.TempDir()
		conformanceWrite(t, root, "specs/concepts/duplicate.concept.spec.md", `---
type: concept-spec
---
# Duplicate
## Actions
### Register
Record a link.
### Register
Record a different link.
`)
		result := conformanceExport(t, root)
		if result.ExitCode != 1 {
			t.Fatal("duplicate declaration name accepted")
		}
	})
}
