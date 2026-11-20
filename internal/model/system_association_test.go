package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/requirement"
)

func associationFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	for _, entry := range []struct{ path, name, attribute string }{
		{"specs/directory-default.system.spec.md", "Directory default", "local"},
		{"remote/specs/explicit-selection.system.spec.md", "Explicit selection", "selected"},
	} {
		conformanceWrite(t, root, entry.path, "---\ntype: system-spec\nid: "+entry.attribute+"-system\n---\n# "+entry.name+"\n## Purpose\nServe users.\n## Scope\nThis service.\n")
		conformanceWrite(t, root, filepath.Join(filepath.Dir(entry.path), "quality/concern/concern.quality-attribute.spec.md"), "---\ntype: quality-attribute-spec\nid: "+entry.attribute+"\n---\n# Concern\n## Definition\nA quality concern.\n")
	}
	body := "## Statement\nRespond promptly.\n\n## Applies to\n**Targets:** [Directory default](../directory-default.system.spec.md)\n\n## Rationale\nKeep users informed.\n\n## Verification plan\nMeasure responses.\n"
	conformanceWrite(t, root, "specs/requirements/explicit.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [selected]\n---\n# Explicit\n\n## Specification system\n[Selected](../../remote/specs/explicit-selection.system.spec.md#selected-system)\n\n"+body)
	conformanceWrite(t, root, "specs/requirements/default.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [local]\n---\n# Default\n\n"+body)
	return root
}

func TestSpecificationSystemSelectionIsConsistent(t *testing.T) {
	root := associationFixture(t)
	for _, slug := range []string{"explicit", "default"} {
		got, err := requirement.Show(root, slug)
		if err != nil {
			t.Fatal(err)
		}
		if got.System == nil || len(got.Classification.Resolved) != 1 || len(got.Classification.Undefined) != 0 {
			t.Fatal(got)
		}
		if (got.System.ID == "selected-system") != (slug == "explicit") {
			t.Fatal(got.System)
		}
	}
	groups, err := requirement.List(root)
	if err != nil || len(groups) != 2 {
		t.Fatal(groups, err)
	}
	report := check.Sources(root, conformanceTypes()...)
	for _, d := range report.Diagnostics {
		if strings.Contains(d.Rule, "undefined-attribute") || d.Severity == "error" {
			t.Fatal(d)
		}
	}
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	r := designRecord(t, s, "Requirement", "Explicit")
	selected := designRecord(t, s, "System", "Explicit selection")
	local := designRecord(t, s, "System", "Directory default")
	if r["specification_system"] != selected["id"] || r["specification_system"] == local["id"] || r["specification_system_origin"] != "explicit" {
		t.Fatal(r)
	}
	if len(Strings(r["quality_attributes"])) != 1 || Strings(r["applies_to"])[0] != local["id"] {
		t.Fatal("quality namespace followed applicability", r)
	}
	index := conformanceIndex(s)
	if index[Strings(r["quality_attributes"])[0]]["owner"] != selected["id"] {
		t.Fatal(r)
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "Requirement", "Explicit")["specification_system"] = local["id"]
	requireFinding(t, forged, "concept.projection")
}

func TestInvalidSpecificationSystemNeverFallsBack(t *testing.T) {
	for _, target := range []string{"../../missing/specs/directory-default.system.spec.md", "https://example.test/directory-default.system.spec.md", "../../remote/specs/explicit-selection.system.spec.md#concern", "../../remote/specs/explicit-selection.system.spec.md) and [Other](../directory-default.system.spec.md"} {
		t.Run(target, func(t *testing.T) {
			root := associationFixture(t)
			traceEdit(t, root, "specs/requirements/explicit.requirement.spec.md", "../../remote/specs/explicit-selection.system.spec.md#selected-system", target)
			if spec, err := requirement.Show(root, "explicit"); err == nil {
				t.Fatal("invalid selection silently accepted", spec)
			}
			result := conformanceExport(t, root)
			if result.ExitCode == 0 {
				t.Fatal("invalid selection accepted")
			}
			r := designRecord(t, result.Snapshot, "Requirement", "Explicit")
			if r["specification_system"] != nil || len(Strings(r["quality_attributes"])) != 0 {
				t.Fatal("fell back to directory context", r)
			}
		})
	}
}

func TestExplicitSelectionDoesNotRequireUnselectedDirectoryDefault(t *testing.T) {
	root := associationFixture(t)
	// Select only this subtree. The malformed ancestor default is irrelevant to
	// the explicit association and must not be read solely for quality context.
	conformanceWrite(t, root, "specs/directory-default.system.spec.md", "not a valid system")
	conformanceWrite(t, root, "client/specs/requirements/selected.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [selected]\n---\n# Selected client\n\n## Specification system\n[System](../../../remote/specs/explicit-selection.system.spec.md)\n\n## Statement\nRespond promptly.\n")
	scope := filepath.Join(root, "client")
	got, err := requirement.Show(scope, "selected")
	if err != nil || got.System == nil || got.System.ID != "selected-system" {
		t.Fatal(got, err)
	}
	result := conformanceExport(t, scope)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	r := designRecord(t, result.Snapshot, "Requirement", "Selected client")
	if r["specification_system_origin"] != "explicit" || len(Strings(r["quality_attributes"])) != 1 {
		t.Fatal(r)
	}
	for _, source := range result.Snapshot.Sources {
		if String(source, "locator") == filepath.ToSlash(filepath.Join(root, "specs/directory-default.system.spec.md")) {
			t.Fatal("unselected default was acquired", source)
		}
	}
}

func TestSpecificationSelectionUsesLinkDestinationsNotLabelsOrImages(t *testing.T) {
	root := associationFixture(t)
	traceEdit(t, root, "specs/requirements/explicit.requirement.spec.md", "[Selected](../../remote/specs/explicit-selection.system.spec.md#selected-system)", "[Selected][system?]\n\n![Context illustration](https://example.test/diagram.png)\n\n[system?]: ../../remote/specs/explicit-selection.system.spec.md#selected-system")
	got, err := requirement.Show(root, "explicit")
	if err != nil || got.System.ID != "selected-system" {
		t.Fatal(got, err)
	}
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	r := designRecord(t, result.Snapshot, "Requirement", "Explicit")
	if len(Strings(r["quality_attributes"])) != 1 || len(Strings(r["specification_system_references"])) != 1 {
		t.Fatal(r)
	}
}

func TestOfflineValidationRejectsClassificationWhenExplicitSelectionIsUnresolved(t *testing.T) {
	root := associationFixture(t)
	traceEdit(t, root, "specs/requirements/explicit.requirement.spec.md", "../../remote/specs/explicit-selection.system.spec.md#selected-system", "../../missing/specs/directory-default.system.spec.md")
	result := conformanceExport(t, root)
	s := cloneSnapshot(t, result.Snapshot)
	r := designRecord(t, s, "Requirement", "Explicit")
	selected := designRecord(t, s, "System", "Explicit selection")
	attribute := Strings(selected["quality_attributes"])[0]
	r["quality_attributes"] = []string{attribute}
	index := conformanceIndex(s)
	label := index[Strings(r["quality_attribute_labels"])[0]]
	label["governing_system"] = selected["id"]
	resolution := index[String(label, "resolution")]
	resolution["status"], resolution["target"] = "resolved", attribute
	requireFinding(t, s, "label.context")
}

func TestExplicitSystemPreservesLocalComponentOwnership(t *testing.T) {
	root := associationFixture(t)
	conformanceWrite(t, root, "adapter/specs/intake.component.spec.md", "---\ntype: component-spec\n---\n# Intake\n## Part of\n[Directory default](../../specs/directory-default.system.spec.md)\n")
	conformanceWrite(t, root, "adapter/specs/requirements/dispatch.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [selected]\n---\n# Dispatch\n## Specification system\n[Selected](../../../remote/specs/explicit-selection.system.spec.md)\n## Statement\nDispatch each job once.\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	req := designRecord(t, s, "Requirement", "Dispatch")
	if req["specification_component"] != designRecord(t, s, "Component", "Intake")["id"] || req["specification_system"] != designRecord(t, s, "System", "Explicit selection")["id"] || len(Strings(req["quality_attributes"])) != 1 || len(Strings(req["applies_to"])) != 0 {
		t.Fatal(req)
	}
}
