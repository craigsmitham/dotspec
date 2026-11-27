package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/speckinds"
	"github.com/craigsmitham/dotspec/internal/system"
)

// association resolves a requirement's system and classification the way
// listing and checking do.
func association(t *testing.T, path string) (*system.System, requirement.Classification, error) {
	t.Helper()
	dir, err := specfile.Directory(filepath.Dir(path))
	if err != nil {
		return nil, requirement.Classification{}, err
	}
	path = filepath.Join(dir, filepath.Base(path))
	spec, err := requirement.Read(path)
	if err != nil {
		return nil, requirement.Classification{}, err
	}
	owner, err := system.Associated(requirement.Layout.ContextDir(path), path)
	if err != nil {
		return nil, requirement.Classification{}, err
	}
	return owner, requirement.Classify(spec.QualityAttributes, owner), nil
}

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
	body := "## Statement\nRespond promptly.\n\n## Applies to\n**Targets:** [Directory default](../directory-default.system.spec.md)\n\n## Rationale\nKeep users informed.\n\n## Evidence\nMeasure responses.\n"
	conformanceWrite(t, root, "specs/requirements/explicit.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [selected]\n---\n# Explicit\n\n## Specification system\n[Selected](../../remote/specs/explicit-selection.system.spec.md#selected-system)\n\n"+body)
	conformanceWrite(t, root, "specs/requirements/default.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [local]\n---\n# Default\n\n"+body)
	return root
}

func TestSpecificationSystemSelectionIsConsistent(t *testing.T) {
	root := associationFixture(t)
	for _, slug := range []string{"explicit", "default"} {
		owner, classification, err := association(t, requirement.Layout.Path(root, slug))
		if err != nil {
			t.Fatal(err)
		}
		if owner == nil || len(classification.Resolved) != 1 || len(classification.Undefined) != 0 {
			t.Fatal(owner, classification)
		}
		if (owner.ID == "selected-system") != (slug == "explicit") {
			t.Fatal(owner)
		}
	}
	report := check.Sources(root, speckinds.All()...)
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
	if r["specification_system"] != selected["id"] || r["specification_owner"] != local["id"] || r["specification_system_origin"] != "explicit" {
		t.Fatal(r)
	}
	if len(Strings(r["quality_attributes"])) != 1 || Strings(r["applies_to"])[0] != local["id"] {
		t.Fatal("quality namespace followed applicability", r)
	}
	index := conformanceIndex(s)
	if index[Strings(r["quality_attributes"])[0]]["specification_system"] != selected["id"] {
		t.Fatal(r)
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "Requirement", "Explicit")["specification_system"] = local["id"]
	requireFinding(t, forged, "shape.invalid")
}

func TestInvalidSpecificationSystemNeverFallsBack(t *testing.T) {
	for _, target := range []string{"../../missing/specs/directory-default.system.spec.md", "https://example.test/directory-default.system.spec.md", "../../remote/specs/explicit-selection.system.spec.md#concern", "../../remote/specs/explicit-selection.system.spec.md) and [Other](../directory-default.system.spec.md"} {
		t.Run(target, func(t *testing.T) {
			root := associationFixture(t)
			traceEdit(t, root, "specs/requirements/explicit.requirement.spec.md", "../../remote/specs/explicit-selection.system.spec.md#selected-system", target)
			if owner, _, err := association(t, requirement.Layout.Path(root, "explicit")); err == nil {
				t.Fatal("invalid selection silently accepted", owner)
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
	owner, _, err := association(t, requirement.Layout.Path(scope, "selected"))
	if err != nil || owner == nil || owner.ID != "selected-system" {
		t.Fatal(owner, err)
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
	owner, _, err := association(t, requirement.Layout.Path(root, "explicit"))
	if err != nil || owner == nil || owner.ID != "selected-system" {
		t.Fatal(owner, err)
	}
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	r := designRecord(t, result.Snapshot, "Requirement", "Explicit")
	if len(Strings(r["quality_attributes"])) != 1 || len(newDesignResolver(result.Snapshot).refs(r, "specification_system")) != 1 {
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
	var label Record
	for _, l := range s.LabelOccurrences {
		if l["owner"] == r["id"] {
			label = l
			break
		}
	}
	label["governing_system"] = selected["id"]
	resolution := Object(label["resolution"])
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
	if req["specification_owner"] != designRecord(t, s, "Component", "Intake")["id"] || req["specification_system"] != designRecord(t, s, "System", "Explicit selection")["id"] || len(Strings(req["quality_attributes"])) != 1 || len(Strings(req["applies_to"])) != 0 {
		t.Fatal(req)
	}
}
