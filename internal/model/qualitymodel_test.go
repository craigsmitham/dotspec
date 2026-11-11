package model

import "testing"

func TestQualityModelOwnershipMappingAndResolution(t *testing.T) {
	root := t.TempDir()
	source := "---\r\ntype: quality-model\r\nid: team\r\n---\r\n# Team\r\n## Description\r\nA reusable model.\r\n## Edition\r\nOne\r\n## Sources\r\n[Original](https://example.test/model).\r\n## Quality attributes\r\nSection narrative.\r\n### Broad\r\n**ID:** broad\r\nA concern.\r\n### Narrow\r\n**ID:** narrow\r\n**Sub-attribute of:** broad, second\r\nA specific concern.\r\n### Second\r\n**ID:** second\r\nAnother grouping.\r\n## Custom\r\nRetain me.\r\n"
	conformanceWrite(t, root, "specs/quality-models/team.quality-model.spec.md", "\ufeff"+source)
	conformanceWrite(t, root, "specs/system.spec.md", "---\ntype: system\n---\n# System\n## Purpose\nP.\n## Scope\nS.\n## Quality attributes\n### Broad\n**ID:** broad\nA distinct local concern.\n")
	conformanceWrite(t, root, "specs/quality/obligation.quality.spec.md", "---\ntype: quality\nattributes: [broad, narrow]\n---\n# Obligation\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatalf("exit %d", result.ExitCode)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	m := elementKind(t, s, "QualityModel")
	if m["owning_system"] != nil {
		t.Fatal("reusable model acquired system ownership")
	}
	for _, field := range []string{"description", "edition", "sources", "quality_attributes_content", "quality_attributes"} {
		if m[field] == nil {
			t.Errorf("missing %s", field)
		}
	}
	if len(Strings(m["extension_sections"])) != 1 {
		t.Fatal("unknown content lost")
	}
	sys := elementKind(t, s, "System")
	for _, label := range s.LabelOccurrences {
		r := snapshotRecord(t, s, String(label, "resolution"))
		if String(label, "role") == "subattribute_of" {
			if String(label, "quality_model_context") != String(m, "id") || label["system_context"] != nil || String(r, "status") != "resolved" {
				t.Fatalf("model parent label: %+v %+v", label, r)
			}
		} else if String(label, "label") == "broad" {
			target := snapshotRecord(t, s, String(r, "target"))
			if String(target, "owner") != String(sys, "id") {
				t.Fatal("requirement classified by reusable definition")
			}
		} else if String(r, "status") != "undefined_label" {
			t.Fatal("model definitions leaked into system classification")
		}
	}
	ids := Strings(m["quality_attributes"])
	rootAttribute := snapshotRecord(t, s, ids[0])
	childAttribute := snapshotRecord(t, s, ids[1])
	if rootAttribute["subattribute_of"] != nil || len(Strings(childAttribute["subattribute_of"])) != 2 || Strings(childAttribute["subattribute_of"])[0] != ids[0] || len(Strings(childAttribute["subattribute_of_labels"])) != 2 {
		t.Fatalf("array grouping: %+v %+v", rootAttribute, childAttribute)
	}
	t.Run("scalar reference rejected by schema", func(t *testing.T) {
		broken := cloneSnapshot(t, s)
		snapshotRecord(t, broken, ids[1])["subattribute_of"] = ids[0]
		report, err := Validate(broken)
		if err != nil {
			t.Fatal(err)
		}
		if len(report) == 0 {
			t.Fatal("scalar accepted for reference array")
		}
	})
	t.Run("no cross-owner parent relation", func(t *testing.T) {
		broken := cloneSnapshot(t, s)
		local := Strings(elementKind(t, broken, "System")["quality_attributes"])[0]
		ids := Strings(elementKind(t, broken, "QualityModel")["quality_attributes"])
		snapshotRecord(t, broken, ids[1])["subattribute_of"] = []string{local}
		requireFinding(t, broken, "quality.scope")
	})
	t.Run("model context cannot classify a requirement", func(t *testing.T) {
		broken := cloneSnapshot(t, s)
		for _, label := range broken.LabelOccurrences {
			if String(label, "role") == "quality_classification" {
				label["quality_model_context"] = m["id"]
			}
		}
		requireFinding(t, broken, "label.context")
	})
}
