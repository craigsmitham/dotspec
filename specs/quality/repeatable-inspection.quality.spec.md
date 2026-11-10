---
type: quality
id: repeatable-inspection
attributes:
  - repeatability
---

# Repeated inspections produce equivalent reports

## Applies to

Completed project-wide and focused specification checks using the same dotspec
build, arguments, working directory, absolute project path, source tree, Git
state, permissions, and environment. Source files and execution conditions
remain unchanged during and between runs. Completed scans can contain
structural errors or writing advice.

## Criterion

Repeated inspections must produce equivalent structured reports, including the
ordering of array elements.

## Measurement

Run each fixture repeatedly without changing its inputs or execution
conditions. Compare parsed JSON reports for equal values, preserving array
order. Cover empty scans, valid specifications, writing advice, structural
errors, and shared system context. Retain differing reports on test failure.

JSON object key order and insignificant whitespace do not affect equivalence.
This comparison does not equate reports from different paths, builds, or
environments.

The [CLI quality-contract tests](../../internal/cli/quality_contracts_test.go)
check this criterion. The obligation remains defined in this document.

## Rationale

Stable reports make changes in findings easier to identify during review and
automation.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
describes inspection. The [check report contract](../../README.md#check-specifications)
defines report fields and their ordering rules.
