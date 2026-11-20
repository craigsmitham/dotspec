package model

import (
	"slices"
	"strings"
	"testing"
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

Members borrow equipment. The finance office cares about replacement cost.

## Environment

A staffed desk with intermittent connectivity.

## Constraints

### Regional retention law

Loan records are kept for seven years.
`

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
	if want := []string{"BusinessRule", "Capability", "Component", "Concept", "Feature", "Job", "QualityModel", "Requirement", "System", "SystemQualityAttribute", "UseCase"}; !slices.Equal(roots, want) {
		t.Fatalf("specification classes %v, want %v", roots, want)
	}
	root := t.TempDir()
	conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem+"\n## Quality attributes\n\n")
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
		if r["specification_system_origin"] != "directory" || r["specification_system"] != designRecord(t, s, "System", "Lending desk")["id"] || r["system_context"] != nil {
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
	for field, want := range map[string]string{"stakeholders": "finance office", "environment": "intermittent connectivity", "constraints": "seven years"} {
		if !strings.Contains(d.text(sys[field]), want) {
			t.Errorf("%s narrative lost: %q", field, d.text(sys[field]))
		}
	}
	goal := designRecord(t, s, "Goal", "Reduce idle equipment")
	if !slices.Equal(Strings(sys["goals"]), []string{String(goal, "id")}) || sys["goals_content"] == nil || goal["owner"] != sys["id"] {
		t.Fatal("goals are a section narrative with declared members", sys["goals"])
	}
	for _, e := range s.Elements {
		if String(e, "kind") == "Constraint" || String(e, "kind") == "ExtensionSection" {
			t.Fatal("a constraint stays narrative; an addressable one is a requirement or business rule", e)
		}
	}
	projected := map[string]int{}
	for _, a := range s.Assertions {
		if String(a, "kind") != "Assertion" || !slices.Contains(compiledModel().catalog.Projections, String(a, "projection_path")) {
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

## Success criterion

An obsolete heading is retained, not interpreted.
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
	if len(records(borrow["minimal_guarantees"])) != 1 || len(records(borrow["success_guarantees"])) != 2 || borrow["success_guarantees_content"] == nil || borrow["success_criterion"] != nil {
		t.Fatal("guarantees are sibling condition lists with their section content", borrow)
	}
	if borrow["goal_level"] != "user_task" || borrow["goal_level_content"] == nil || manage["goal_level"] != "summary" {
		t.Fatal("goal level keeps authored content beside its recognized value", borrow["goal_level"])
	}
	if Object(borrow["sequence"])["presence"] != "authored" || borrow["related_jobs"] != nil {
		t.Fatal("sequence is narrative", borrow["sequence"])
	}
	if len(Strings(borrow["extension_sections"])) != 1 {
		t.Fatal("the former Success criterion heading should be retained as unmodelled content", borrow["extension_sections"])
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

// Jobs, business rules and requirements state needs, rules and one obligation.
// They do not contain further obligation clauses.
func TestObligationClausesBelongOnlyToTheirContainers(t *testing.T) {
	clause := "\n## Obligations\n\n### Extra\n\nAn embedded obligation.\n"
	for _, tc := range []struct{ path, source string }{
		{"specs/jobs/borrow.job.spec.md", "---\ntype: job-spec\n---\n# Borrow\n\n## Job\n\nBorrow equipment.\n" + clause},
		{"specs/business-rules/limit.business-rule.spec.md", "---\ntype: business-rule-spec\n---\n# Limit\n\n## Statement\n\nAt most three items.\n" + clause},
		{"specs/requirements/record.requirement.spec.md", "---\ntype: requirement-spec\n---\n# Record\n\n## Statement\n\nRecord each loan.\n" + clause},
	} {
		t.Run(tc.path, func(t *testing.T) {
			root := t.TempDir()
			conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
			conformanceWrite(t, root, tc.path, tc.source)
			result := conformanceExport(t, root)
			found := false
			for _, d := range result.Snapshot.Diagnostics {
				found = found || String(d, "code") == "obligation.container"
			}
			if !found || result.ExitCode == 0 {
				t.Fatal("embedded obligations must be rejected here", result.ExitCode, result.Snapshot.Diagnostics)
			}
			// The rejected section stays visible as retained content, never as unbound source.
			if len(Strings(result.Snapshot.Coverage["unbound_content_nodes"])) != 0 || len(Strings(result.Snapshot.Coverage["extension_sections"])) == 0 {
				t.Fatal("rejected obligations section should be retained as an extension section", result.Snapshot.Coverage)
			}
			checkSourceText(t, result.Snapshot)
			for _, e := range result.Snapshot.Elements {
				if String(e, "kind") == "ObligationClause" {
					t.Fatal("clause created without a container", e)
				}
			}
		})
	}
	t.Run("use cases still contain clauses", func(t *testing.T) {
		root := t.TempDir()
		conformanceWrite(t, root, "specs/lending-desk.system.spec.md", typedSystem)
		conformanceWrite(t, root, "specs/use-cases/borrow.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Borrow\n\n## Goal\n\nBorrow.\n"+clause)
		result := conformanceExport(t, root)
		if result.ExitCode != 0 || len(Strings(designRecord(t, result.Snapshot, "UseCase", "Borrow")["obligations"])) != 1 {
			t.Fatal(result.Snapshot.Diagnostics)
		}
	})
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
	if len(records(step["open_questions"])) != 1 || step["open_questions_content"] == nil || step["extension_sections"] != nil {
		t.Fatal("open question under a job step should be its own field", step)
	}
}
