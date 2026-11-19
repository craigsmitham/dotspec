---
type: specification-section
generated_by: dotspec
subject_kind: requirement
---
# Repeated inspections produce equivalent reports

[dotspec](<../../SPEC.md>) · [Contents](<../index.md>)

## Statement

Repeated inspections must produce equivalent structured reports, including the
ordering of array elements.

## Applies to

Completed project-wide and focused specification checks using the same dotspec
build, arguments, working directory, absolute project path, source tree, Git
state, permissions, and environment. Source files and execution conditions
remain unchanged during and between runs. Completed scans can contain
structural errors or writing advice.

## Criterion

For the equivalence required by [Statement](<#statement>), parsed JSON values must
be equal, preserving array order. JSON object key order and insignificant whitespace
do not affect equivalence. The comparison applies only within the operating context
in [Applies to](<#applies-to>).

## Verification plan

Run each fixture repeatedly without changing its inputs or execution
conditions. Compare reports using [Criterion](<#criterion>).
Cover empty scans, valid specifications, writing advice, structural
errors, and shared system context. Retain differing reports on test failure.

The [CLI quality-contract tests](<../../internal/cli/quality_contracts_test.go>)
check this criterion. The obligation remains defined in this document.

## Rationale

Stable reports make changes in findings easier to identify during review and
automation.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.md>)
describes inspection. The [check report contract](<../../docs/reference/inspection.md#check-specifications>)
defines report fields and their ordering rules.

## Related specifications

- Quality attribute — [Repeatability](<../../SPEC.md#repeatability>)

[Edit source](<../../specs/requirements/repeatable-inspection.requirement.spec.md#repeated-inspections-produce-equivalent-reports>)
