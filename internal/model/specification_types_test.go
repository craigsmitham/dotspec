package model

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/dotspec-md/dotspec/internal/speckinds"
)

const typedSystem = `---
type: system-spec
---
# Lending desk

## Purpose

Lend equipment to members.

## Scope

Reservations and returns.

## Goals

### Reduce idle equipment

Equipment spends less time unused.

## Stakeholders

### Member

Members borrow equipment.

### Finance office

The finance office cares about replacement cost.

## Environment

A staffed desk with intermittent connectivity.

## Constraints

### Regional retention law

Loan records are kept for seven years.
`

// Every built-in kind needs a root class, a domain mapping and an authored
// suffix that agrees with its layout, or its files cannot enter the model.
func TestBuiltInKindsAreRegisteredWithTheModel(t *testing.T) {
	for _, typ := range speckinds.All() {
		kind := typ.Layout.Kind
		if specfile.DocumentRootClass(kind) == "" {
			t.Errorf("kind %q has no document root class", kind)
		}
		if suffix := specfile.AuthoredSuffix(kind); suffix != typ.Layout.Suffix {
			t.Errorf("kind %q authored suffix %q differs from layout suffix %q", kind, suffix, typ.Layout.Suffix)
		}
		doc := parseMarkdownContent("specs/example"+typ.Layout.Suffix, []byte("---\ntype: "+specfile.AuthoredType(kind)+"\n---\n# Example\n"), true)
		for _, d := range MapDomain(doc, kind).Diagnostics {
			if strings.Contains(String(d, "message"), "no domain mapping") {
				t.Errorf("kind %q has no domain mapping", kind)
			}
		}
	}
}

// A specification is one authored document. Each specification type has one
// concrete class, and a requirement is that type regardless of its classification.
func TestSpecificationTypesCorrespondToDocumentRoots(t *testing.T) {
	catalog := compiledModel().catalog
	roots := []string{}
	for name, class := range catalog.Classes {
		if !class.Abstract && slices.Contains(class.Ancestors, "Specification") {
			roots = append(roots, name)
		}
	}
	slices.Sort(roots)
	var declared []string
	for _, tag := range catalog.Enums["SpecificationType"] {
		kind, ok := strings.CutSuffix(tag, "-spec")
		if !ok || specfile.AuthoredSuffix(kind) != "."+kind+".spec.md" {
			t.Fatalf("unexpected specification form %q", tag)
		}
		declared = append(declared, specfile.DocumentRootClass(kind))
	}
	slices.Sort(declared)
	if !slices.Equal(roots, declared) {
		t.Fatalf("concrete specification classes %v, declared file types map to %v", roots, declared)
	}
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
	conformanceWrite(t, root, "specs/quality/speed/speed.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nid: speed\n---\n# Speed\n\n## Definition\n\nResponse speed.\n")

	conformanceWrite(t, root, "specs/requirements/fast.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [speed]\n---\n# Fast\n\n## Statement\n\nRespond promptly.\n")
	conformanceWrite(t, root, "specs/requirements/plain.requirement.spec.md", "---\ntype: requirement-spec\n---\n# Plain\n\n## Statement\n\nRecord each loan.\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
	if types := Strings(s.Scope["selected_spec_types"]); slices.Contains(types, "quality") || !slices.Contains(types, "requirement-spec") {
		t.Fatalf("selected specification types: %v", types)
	}
	fast, plain := designRecord(t, s, "Requirement", "Fast"), designRecord(t, s, "Requirement", "Plain")
	if len(Strings(fast["quality_attributes"])) != 1 || len(Strings(plain["quality_attributes"])) != 0 {
		t.Fatal("quality is a classification of a requirement, not another type", fast, plain)
	}
	for _, r := range []Record{fast, plain} {
		if r["specification_system_origin"] != "directory" || r["specification_system"] != designRecord(t, s, "System", "Lending desk")["id"] {
			t.Fatal("directory default should be recorded as the specification system", r)
		}
	}
}

func TestSystemNarrativeSectionsAndGoals(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
	d := newDesignResolver(s)
	sys := designRecord(t, s, "System", "Lending desk")
	for field, want := range map[string]string{"environment": "intermittent connectivity", "constraints": "seven years"} {
		if !strings.Contains(d.text(sys[field]), want) {
			t.Errorf("%s narrative lost: %q", field, d.text(sys[field]))
		}
	}
	member := designRecord(t, s, "SystemStakeholderRole", "Member")
	finance := designRecord(t, s, "SystemStakeholderRole", "Finance office")
	if !slices.Equal(Strings(sys["stakeholder_roles"]), []string{String(member, "id"), String(finance, "id")}) || member["owner"] != sys["id"] || finance["owner"] != sys["id"] || !strings.Contains(d.text(finance["description"]), "replacement cost") {
		t.Fatalf("stakeholder roles lost their ownership, order or description: %v %v", member, finance)
	}
	goal := designRecord(t, s, "Goal", "Reduce idle equipment")
	if !slices.Equal(Strings(sys["goals"]), []string{String(goal, "id")}) || goal["owner"] != sys["id"] {
		t.Fatal("goals are a section narrative with declared members", sys["goals"])
	}
	projected := map[string]int{}
	for _, a := range s.Relationships() {
		if !slices.Contains(compiledModel().catalog.Projections, String(a, "projection_path")) {
			t.Fatal("assertions are views of declared projection paths", a)
		}
		projected[String(a, "projection_path")]++
	}
	if projected["System.goals"] != 1 {
		t.Fatal("every declaration membership is projected", projected)
	}
}

func TestUseCaseGuaranteesSupportAndIllustration(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
	conformanceWrite(t, root, "specs/business-rules/limit.business-rule.spec.md", "---\ntype: business-rule-spec\nid: limit\n---\n# Loan limit\n\n## Statement\n\nA member holds at most three items.\n")
	conformanceWrite(t, root, "specs/use-cases/manage.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Manage loans\n\n## Goal\n\nKeep loans in order.\n\n## Goal level\n\nSummary.\n")
	conformanceWrite(t, root, "specs/use-cases/borrow.use-case.spec.md", `---
type: use-case-spec
---
# Borrow an item

## Goal

Take an item home.

## Goal level

User task.

## Supports

[Manage loans](manage.use-case.spec.md).

## Minimal guarantees

- No item leaves without a recorded loan.

## Success guarantees

- The member holds the item.
- The loan is recorded.

## Main success scenario

1. The member requests an item.
2. The desk confirms the loan.
   - **Illustrates:** [Loan limit](../business-rules/limit.business-rule.spec.md).

## Sequence

Representative.

### Example ordering

This sequence is representative.
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
	borrow, manage := designRecord(t, s, "UseCase", "Borrow an item"), designRecord(t, s, "UseCase", "Manage loans")
	rule := designRecord(t, s, "BusinessRule", "Loan limit")
	if !slices.Equal(Strings(borrow["supports"]), []string{String(manage, "id")}) {
		t.Fatal("a use case supports its higher-level use case", borrow["supports"])
	}
	d := newDesignResolver(s)
	if !strings.Contains(d.text(borrow["minimal_guarantees"]), "No item leaves") || !strings.Contains(d.text(borrow["success_guarantees"]), "The loan is recorded") || strings.TrimSpace(d.text(borrow["goal_level"])) != "User task." || borrow["sequence"] == nil {
		t.Fatal(borrow)
	}
	illustrated := false
	for _, e := range s.Elements {
		if String(e, "kind") == "ScenarioStep" && slices.Equal(Strings(e["illustrates"]), []string{String(rule, "id")}) {
			illustrated = true
			if designKind("ScenarioStep", "TraceableElement") || e["supports"] != nil {
				t.Fatal("a scenario step illustrates; it is not itself traceable")
			}
		}
	}
	if !illustrated {
		t.Fatal("a step may illustrate a business rule")
	}
}

func TestBusinessRuleRefinesOnlyBusinessRules(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
	conformanceWrite(t, root, "specs/requirements/record.requirement.spec.md", "---\ntype: requirement-spec\n---\n# Record\n\n## Statement\n\nRecord each loan.\n")
	conformanceWrite(t, root, "specs/business-rules/limit.business-rule.spec.md", "---\ntype: business-rule-spec\n---\n# Limit\n\n## Statement\n\nAt most three items.\n")
	conformanceWrite(t, root, "specs/business-rules/student-limit.business-rule.spec.md", "---\ntype: business-rule-spec\n---\n# Student limit\n\n## Statement\n\nStudents hold at most two items.\n\n## Refines\n\n[Limit](limit.business-rule.spec.md).\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	student, limit := designRecord(t, result.Snapshot, "BusinessRule", "Student limit"), designRecord(t, result.Snapshot, "BusinessRule", "Limit")
	if !slices.Equal(Strings(student["refines"]), []string{String(limit, "id")}) {
		t.Fatal("rule-to-rule refinement lost", student["refines"])
	}
	traceEdit(t, root, "specs/business-rules/student-limit.business-rule.spec.md", "[Limit](limit.business-rule.spec.md)", "[Record](../requirements/record.requirement.spec.md)")
	result = conformanceExport(t, root)
	student = designRecord(t, result.Snapshot, "BusinessRule", "Student limit")
	if result.ExitCode == 0 || len(Strings(student["refines"])) != 0 {
		t.Fatal("a business rule must not refine a software requirement", student["refines"])
	}
}

// Open questions are a field of every domain element, so an explicitly headed
// question is recognized wherever an element's subsections are mapped.
func TestOpenQuestionsBelongToAnyElement(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
	conformanceWrite(t, root, "specs/jobs/borrow.job.spec.md", "---\ntype: job-spec\n---\n# Borrow\n\n## Job\n\nBorrow equipment.\n\n## Steps\n\n### Reserve\n\nChoose an item.\n\n#### Open questions\n\n- How far ahead can members reserve?\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	checkSourceText(t, result.Snapshot)
	step := designRecord(t, result.Snapshot, "JobStep", "Reserve")
	if len(records(step["open_questions"])) != 1 {
		t.Fatal("open question under a job step should be its own field", step)
	}
}

func TestFormerInlineDeclarationsAreLocatedAtTheirHeading(t *testing.T) {
	for _, tc := range []struct{ kind, container string }{
		{"system", ""}, {"component", ""}, {"use-case", ""}, {"concept", ""},
		{"system", "## Provided interfaces\n### Requests\n"}, {"concept", "## Actions\n### Restore\n"},
		{"system", "## Concept uses\n### Archive\n"}, {"component", "## Synchronizations\n### Restore entry\n"},
		{"requirement", ""}, {"business-rule", ""}, {"job", ""},
	} {
		t.Run(tc.kind+"/"+tc.container, func(t *testing.T) {
			heading := "##"
			if tc.container != "" {
				heading = "####"
			}
			prefix := fmt.Sprintf("---\ntype: %s-spec\n---\n# Service\n%s", tc.kind, tc.container)
			path := "service." + tc.kind + ".spec.md"
			doc := ParseContent(path, []byte(prefix+heading+" Obligations\n"+heading+"# Keep entries\nRetain every entry.\n"))
			mapped := MapDomain(doc, tc.kind)
			if len(mapped.Elements) != 0 || len(mapped.Diagnostics) != 1 {
				t.Fatal(mapped)
			}
			finding := mapped.Diagnostics[0]
			if finding["code"] != "structure.inline-requirement" || finding["line"] != strings.Count(prefix, "\n")+1 || !strings.Contains(String(finding, "message"), "*.requirement.spec.md") {
				t.Fatal(finding)
			}
		})
	}
}

func TestScenarioExamplesRetainTheirSource(t *testing.T) {
	source := "---\ntype: use-case-spec\n---\n# Return\n## Main success scenario\n1. The person returns.\n## Success guarantees\nSaved entries are available.\n\n~~~markdown\n## Obligations\n~~~\n"
	doc := ParseContent("return.use-case.spec.md", []byte(source))
	mapped := MapDomain(doc, "use-case")
	if len(mapped.Diagnostics) != 0 || len(mapped.Elements) != 3 || mapped.Elements[1]["kind"] != "Scenario" || mapped.Elements[2]["kind"] != "ScenarioStep" || doc.Source["markdown"] != source {
		t.Fatal(mapped)
	}
}
