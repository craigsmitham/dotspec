---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/quality/reliability/repeatability/repeatable-inspection.requirement.spec.md#repeated-inspections-produce-equivalent-reports"
---
# Repeated inspections produce equivalent reports

[dotspec](<../dotspec.system.spec-report.md>)

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

## Evidence

Checks for this requirement include [CLI quality-contract tests](<../../internal/cli/quality_contracts_test.go>).

## Rationale

Stable reports make changes in findings easier to identify during review and
automation.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
describes inspection. The [check report contract](<../../docs/reference/inspection.md#check-specifications>)
defines report fields and their ordering rules.

**Quality attribute** · [Repeatability](<../quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>)
