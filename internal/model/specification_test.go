package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInlineSpecificationCommonContent(t *testing.T) {
	for _, tc := range []struct{ kind, section string }{
		{"Component", "Components"}, {"Feature", "Features"}, {"Capability", "Capabilities"},
		{"ConceptUse", "Concept uses"}, {"Synchronization", "Synchronizations"}, {"Interface", "Required interfaces"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			root := t.TempDir()
			source := fmt.Sprintf(`---
type: system
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

#### Validation plan

##### Method

Review with the intended users.
`, tc.section)
			if err := os.MkdirAll(filepath.Join(root, "specs"), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "specs", "system.spec.md"), []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			result := conformanceExport(t, root)
			if result.ExitCode != 0 {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			s := result.Snapshot
			conformancePayloads(t, s)
			r := designRecord(t, s, tc.kind, "Named definition")
			if !designKind(tc.kind, "Specification") || r["authored_id"] != "named-definition" || r["extension_sections"] != nil {
				t.Fatal("inline definition not recognized as a specification", r)
			}
			for _, field := range []string{"preamble", "notes", "evidence", "related_specifications", "validation_plan"} {
				if r[field] == nil {
					t.Fatal("missing common specification content", field)
				}
			}
			evidence := designRecord(t, s, "EvidenceStatement", "Observation")
			if evidence["owner"] != r["id"] {
				t.Fatal("evidence ownership lost", evidence)
			}
		})
	}
}

func TestInterfaceSpecificationRelationships(t *testing.T) {
	root := reviewFixture(t)
	traceEdit(t, root, "specs/system.spec.md", "## Assumptions", `## Features

### Reliable recovery

#### Groups

[Restoration requests](#restoration-requests).

## Assumptions`)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "[DocumentTrash](../system.spec.md#documenttrash)", "[DocumentTrash](../system.spec.md#documenttrash), [Restoration requests](../system.spec.md#restoration-requests)")
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
	req := designRecord(t, s, "GeneralRequirement", "Preserve document metadata on restoration")
	if ids := Strings(req["applies_to"]); len(ids) != 3 || ids[2] != provided["id"] {
		t.Fatal("cross-file interface reference did not resolve", req)
	}
	for kind, object := range map[string]Record{"RequiredInterface": required, "ProvidedInterface": provided} {
		count := 0
		for _, a := range s.Assertions {
			if String(a, "kind") == kind && a["subject"] == sys["id"] && a["object"] == object["id"] {
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
	Object(fSys["provided_interfaces"])["interfaces"] = append(Strings(Object(fSys["provided_interfaces"])["interfaces"]), String(required, "id"))
	requireFinding(t, forged, "interface.direction")
	forged = cloneSnapshot(t, s)
	designRecord(t, forged, "Interface", "Restoration requests")["owner"] = required["id"]
	requireFinding(t, forged, "reference.kind")
}

func TestInterfaceDeclarationIDsAreChecked(t *testing.T) {
	for _, duplicate := range []string{"**ID:** time-readings", "**ID:** restoration-requests\n\n**ID:** another-id"} {
		root := reviewFixture(t)
		traceEdit(t, root, "specs/system.spec.md", "**ID:** restoration-requests", duplicate)
		result := conformanceExport(t, root)
		if result.ExitCode == 0 {
			t.Fatal("duplicate declaration identity accepted")
		}
		conformancePayloads(t, result.Snapshot)
	}
}
