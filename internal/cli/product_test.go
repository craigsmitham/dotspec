package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/model"
)

func TestProductChangesAndVerificationObservations(t *testing.T) {
	projectFixture(t)
	const req = "specs/requirements/recover.requirement.spec.md"
	const protocol = "protocol.md"
	writeFixture(t, protocol, "# Protocol\n\n## Restart\nReopen the service and read the saved value.\n")
	writeFixture(t, "tests/restart.sh", "#!/bin/sh\nexit 0\n")
	writeFixture(t, "specs/quality/continuity/continuity.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\n---\n# Continuity\n## Definition\nPreserve an acknowledged change.\n")
	read := func(path string) string {
		t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	run := func(args ...string) string {
		t.Helper()
		code, out, errs := invoke(args...)
		if code != 0 {
			t.Fatalf("%v: %d %s %s", args, code, out, errs)
		}
		return out
	}
	product := func(args ...string) model.ProductSpecification {
		t.Helper()
		var view model.ProductSpecification
		if err := json.Unmarshal([]byte(run(append([]string{"model", "product", "tool"}, args...)...)), &view); err != nil {
			t.Fatal(err)
		}
		return view
	}
	run("spec", "build")
	if text := read("specification/tool.system.spec-report.md"); !strings.Contains(text, "## Quality Attributes") || !strings.Contains(text, "Preserve an acknowledged change.") {
		t.Fatal(text)
	}
	const body = "---\ntype: requirement-spec\nquality_attributes: [continuity]\ncompanions:\n  - target: ../../tests/restart.sh\n    role: verifies\n  - target: ../../protocol.md#restart\n    role: verifies\n    note: Manual restart protocol\n---\n# Recover\n\n## Statement\nRestore acknowledged work.\n\n## Criterion\nThe saved value survives restart.\n\n## Evidence\n[Citation](../../protocol.md#restart).\n"
	writeFixture(t, req, body)
	before := product()
	if !before.Complete || len(before.Members) != 3 || len(before.Adoptions) != 0 {
		t.Fatal(before)
	}
	if partial := product("--path", "specs/requirements"); partial.Complete || len(partial.Members) != 2 {
		t.Fatal(partial)
	}
	writeFixture(t, req, body+"\n## Adopts\n\n### Recovery protocol\n\n[Restart](../../protocol.md#restart).\n\n#### Edition\n\nevolving\n\n#### Applicability\n\nAfter acknowledged writes.\n")
	adopted := product()
	if len(adopted.Members) != 3 || len(adopted.Adoptions) != 1 || adopted.Adoptions[0]["provision"] != "../../protocol.md#restart" {
		t.Fatal(adopted.Adoptions)
	}
	run("spec", "build")
	const report = "specification/requirements/recover.requirement.spec-report.md"
	if text := read(report); !strings.Contains(text, "Procedure linked.") || !strings.Contains(text, "No matching result recorded.") || !strings.Contains(text, "After acknowledged writes.") {
		t.Fatal(text)
	}
	if !strings.Contains(read("specification/quality/continuity/continuity.quality-attribute.spec-report.md"), "recover.requirement.spec-report.md") {
		t.Fatal("classified requirement missing from attribute")
	}
	ctx := model.VerificationContext{ProductRevision: "r2", Conditions: "local storage", SpecificationRevision: "v2"}
	evidence := model.VerificationEvidence{Selection: ctx}
	for i, outcome := range []string{"fail", "pass", "fail", "pass", "inconclusive", "fail"} {
		result := model.VerificationResult{VerificationContext: ctx, Subject: req, Procedure: "tests/restart.sh", Outcome: outcome, Observation: "https://example.test/observations/" + outcome}
		switch i {
		case 1:
			result.Procedure = protocol + "#restart"
		case 2:
			result.ProductRevision = "r1"
		case 3:
			result.Conditions = "other"
		case 4:
			result.SpecificationRevision = "v1"
		case 5:
			result.ProductRevision = ""
		}
		evidence.Results = append(evidence.Results, result)
	}
	data, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, "verification-results.json", string(data))
	run("spec", "build")
	text := read(report)
	if strings.Count(text, "Matching observation ·") != 2 || strings.Count(text, "Other recorded observation ·") != 4 || strings.Contains(text, "No matching result recorded.") {
		t.Fatal(text)
	}
	observed := product()
	if len(observed.Verification) != 1 || len(observed.Verification[0].Observations) != 6 || observed.VerificationSelection == nil {
		t.Fatal(observed.Verification)
	}
	if !reflect.DeepEqual(adopted.Members, observed.Members) {
		t.Fatal("recording runs changed product membership")
	}
	run("check")
	writeFixture(t, "verification-results.json", `{"selection":{},"results":[]}`)
	if code, _, _ := invoke("spec", "build"); code == 0 || read(report) != text {
		t.Fatal("invalid evidence replaced the last valid report")
	}
}

func TestReferenceOnlyCatalog(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	writeFixture(t, "dotspec.yaml", "name: vocabulary\n")
	writeFixture(t, "specs/quality-models/team.quality-model.md", "---\ntype: quality-model\n---\n# Team\n## Quality attributes\n### Continuity\n**ID:** continuity\nKeep the record.\n")
	for _, args := range [][]string{{"spec", "build"}, {"check"}, {"model", "product"}} {
		code, out, errs := invoke(args...)
		if code != 0 {
			t.Fatal(args, out, errs)
		}
		if args[0] == "model" {
			var view model.ProductSpecification
			if err := json.Unmarshal([]byte(out), &view); err != nil || view.System != "" || len(view.Members) != 0 || len(view.ReferenceCatalog) != 1 {
				t.Fatal(out)
			}
		}
	}
	page, err := os.ReadFile(filepath.Join(root, "specification", "quality-models.md"))
	if err != nil || !strings.Contains(string(page), "team.quality-model-report.md") {
		t.Fatal(err, string(page))
	}
}

func TestProductWithStructuralParentInAnotherProject(t *testing.T) {
	projectFixture(t)
	writeFixture(t, "specs/tool.system.spec.md", "---\ntype: system-spec\n---\n# Tool\n## Purpose\nAssist Demo.\n## Scope\nA local tool.\n## Part of\n[Demo](../apps/demo/specs/demo.system.spec.md).\n")
	for _, args := range [][]string{{"spec", "build"}, {"check"}, {"model", "product"}} {
		code, out, errs := invoke(args...)
		if code != 0 {
			t.Fatal(args, out, errs)
		}
		if args[0] == "model" {
			var v model.ProductSpecification
			if err := json.Unmarshal([]byte(out), &v); err != nil || !v.Complete || len(v.Members) != 1 || v.Members[0].Name != "Tool" {
				t.Fatal(out)
			}
		}
	}
	report, err := os.ReadFile("specification/tool.system.spec-report.md")
	if err != nil || !strings.Contains(string(report), "### Part of") || !strings.Contains(string(report), "demo.system.spec.md") {
		t.Fatal(err, string(report))
	}
}

func TestLocallyAuthoredRequirementForExternalSystem(t *testing.T) {
	projectFixture(t)
	writeFixture(t, "specs/requirements/remote.requirement.spec.md", "---\ntype: requirement-spec\n---\n# Remote\n## Specification system\n[Demo](../../apps/demo/specs/demo.system.spec.md)\n## Statement\nNotify an operator.\n")
	for _, args := range [][]string{{"spec", "build"}, {"check"}, {"model", "product", "tool"}} {
		code, out, errs := invoke(args...)
		if code != 0 {
			t.Fatal(args, out, errs)
		}
		if args[0] == "model" {
			var v model.ProductSpecification
			if err := json.Unmarshal([]byte(out), &v); err != nil || !v.Complete || len(v.Members) != 1 {
				t.Fatal(out)
			}
		}
	}
	index, err := os.ReadFile("specification/requirements-for-other-systems.md")
	if err != nil || !strings.Contains(string(index), "remote.requirement.spec-report.md") {
		t.Fatal(err, string(index))
	}
}
