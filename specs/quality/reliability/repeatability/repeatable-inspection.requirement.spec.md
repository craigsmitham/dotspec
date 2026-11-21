---
type: requirement-spec
quality_attributes: [repeatability]
id: repeatable-inspection
---

# Repeated inspections produce equivalent reports

## Applies to

Completed project-wide and focused specification checks using the same dotspec
build, arguments, working directory, absolute project path, source tree, Git
state, permissions, and environment. Source files and execution conditions
remain unchanged during and between runs. Completed scans can contain
structural errors or writing advice.

## Statement

Repeated inspections must produce equivalent structured reports, including the
ordering of array elements.

## Criterion

For the equivalence required by [Statement](#statement), parsed JSON values must
be equal, preserving array order. JSON object key order and insignificant whitespace
do not affect equivalence. The comparison applies only within the operating context
in [Applies to](#applies-to).

## Evidence

Checks for this requirement include [CLI quality-contract tests](../../../../internal/cli/quality_contracts_test.go).

## Rationale

Stable reports make changes in findings easier to identify during review and
automation.

## Related specifications

[Specification Inspection](../../../concepts/specification-inspection.concept.spec.md)
describes inspection. The [check report contract](../../../../docs/cli-handbook/reference/inspection.md#check-specifications)
defines report fields and their ordering rules.
