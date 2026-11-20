package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Elements declared by name inside a system specification are addressable and
// traceable, but only an authored document is a specification.
func TestDeclaredElementsAreNotSpecifications(t *testing.T) {
	for _, tc := range []struct{ kind, section string }{
		{"ConceptUse", "Concept uses"}, {"Synchronization", "Synchronizations"}, {"Interface", "Required interfaces"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			root := t.TempDir()
			source := fmt.Sprintf(`---
type: system-spec
---
# Example

## Purpose

Describe a bounded example.

## Scope

A single specification declaration.

## %s

### Named definition

**ID:** named-definition

An authored definition.

#### Notes

Keep this note on the definition.

#### Evidence

##### Observation

An observation awaiting assessment.

#### Related specifications

[This definition](#named-definition).

#### System validation plan

##### Method

Review with the intended users.
`, tc.section)
			if err := os.MkdirAll(filepath.Join(root, "specs"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "specs", "example.system.spec.md"), []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			result := conformanceExport(t, root)
			if result.ExitCode != 0 {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			s := result.Snapshot
			conformancePayloads(t, s)
			r := designRecord(t, s, tc.kind, "Named definition")
			if designKind(tc.kind, "Specification") || !designKind(tc.kind, "TraceableElement") || r["authored_id"] != "named-definition" {
				t.Fatal("declared element must be traceable and addressable without being a specification", r)
			}
			for _, field := range []string{"notes", "assessment_plans"} {
				if r[field] == nil {
					t.Fatal("missing element content", field)
				}
			}
			for _, field := range []string{"preamble", "evidence", "related_specifications", "slug"} {
				if r[field] != nil {
					t.Fatal("document-level content mapped onto a declared element", field)
				}
			}
			// Sections that belong to specifications stay with the element as retained, unmodelled content.
			if len(Strings(r["extension_sections"])) != 2 {
				t.Fatal("Evidence and Related specifications should be retained as the element's extension sections", r["extension_sections"])
			}
		})
	}
}

func TestInterfaceSpecificationRelationships(t *testing.T) {
	root := reviewFixture(t)
	conformanceWrite(t, root, "specs/recovery.feature.spec.md", `---
type: feature-spec
---
# Reliable recovery

## Groups

[Restoration requests](document-workspace.system.spec.md#restoration-requests).
`)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "[DocumentTrash](../document-workspace.system.spec.md#documenttrash)", "[DocumentTrash](../document-workspace.system.spec.md#documenttrash), [Restoration requests](../document-workspace.system.spec.md#restoration-requests)")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	sys := designRecord(t, s, "System", "Document workspace")
	provided := designRecord(t, s, "Interface", "Restoration requests")
	required := designRecord(t, s, "Interface", "Time readings")
	clause := designRecord(t, s, "ObligationClause", "Preserve request identity")
	feature := designRecord(t, s, "Feature", "Reliable recovery")
	if provided["owner"] != sys["id"] || required["owner"] != sys["id"] || clause["owner"] != provided["id"] {
		t.Fatal("interface or clause ownership lost")
	}
	if ids := Strings(feature["groups"]); len(ids) != 1 || ids[0] != provided["id"] {
		t.Fatal("feature cannot group an interface", feature)
	}
	if ids := Strings(clause["applies_to"]); len(ids) != 1 || ids[0] != provided["id"] {
		t.Fatal("interface ID did not resolve as a subject", clause)
	}
	req := designRecord(t, s, "Requirement", "Preserve document metadata on restoration")
	if ids := Strings(req["applies_to"]); len(ids) != 3 || ids[2] != provided["id"] {
		t.Fatal("cross-file interface reference did not resolve", req)
	}
	for kind, object := range map[string]Record{"System.required_interfaces": required, "System.provided_interfaces": provided} {
		count := 0
		for _, a := range s.Assertions {
			if String(a, "projection_path") == kind && a["subject"] == sys["id"] && a["object"] == object["id"] {
				count++
			}
		}
		if count != 1 {
			t.Fatal("interface direction missing from graph", kind, count)
		}
	}
	text, err := OverviewMarkdown(s, root)
	if err != nil || !strings.Contains(text, "## Interfaces") || !strings.Contains(text, "Interface: Restoration requests") {
		t.Fatal("interface absent from overview", err, text)
	}
	forged := cloneSnapshot(t, s)
	fSys := designRecord(t, forged, "System", "Document workspace")
	fSys["provided_interfaces"] = append(Strings(fSys["provided_interfaces"]), String(required, "id"))
	requireFinding(t, forged, "interface.direction")
	forged = cloneSnapshot(t, s)
	designRecord(t, forged, "Interface", "Restoration requests")["owner"] = required["id"]
	requireFinding(t, forged, "reference.kind")
}

func TestInterfaceDeclarationIDsAreChecked(t *testing.T) {
	for _, duplicate := range []string{"**ID:** time-readings", "**ID:** restoration-requests\n\n**ID:** another-id"} {
		root := reviewFixture(t)
		traceEdit(t, root, "specs/document-workspace.system.spec.md", "**ID:** restoration-requests", duplicate)
		result := conformanceExport(t, root)
		if result.ExitCode == 0 {
			t.Fatal("duplicate declaration identity accepted")
		}
		conformancePayloads(t, result.Snapshot)
	}
}
