---
generated_by: dotspec
type: requirement-spec-report
documents: "../../../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md#diagnostic-severity-is-available-without-styling"
---
# Diagnostic severity is available without styling

[dotspec](<../../../dotspec.system.spec-report.md>)

## Statement

Every diagnostic must identify its severity using a text label independently
of terminal styling.

## Applies to

Human-readable reports from project-wide and focused specification checks.

## Criterion

Use the text-label rule in [Statement](<#statement>), with supported severity values
defined by the [check report contract](<../../../../docs/reference/inspection.md#check-specifications>).

## Verification plan

Render reports containing each supported diagnostic severity. Capture output
with terminal-style environment settings and in redirected-output conditions,
including with color disabled. Compare each diagnostic with
[Criterion](<#criterion>). Retain captured output on test failure.

The [CLI quality-contract tests](<../../../../internal/cli/quality_contracts_test.go>)
check this criterion. The obligation remains defined in this document.

## Rationale

Contributors need to distinguish diagnostic severities when reading plain-text
logs or using tools that do not convey color. This criterion addresses that
specific need; it does not establish overall accessibility.

## Related specifications

[Specification Inspection](<../../../concepts/specification-inspection.concept.spec-report.md>)
describes the report and inspection behavior. The
[check report contract](<../../../../docs/reference/inspection.md#check-specifications>) defines the
supported severity values and report structure.

**Quality attribute**

- [Accessibility](<../accessibility/accessibility.quality-attribute.spec-report.md>)
- [Diagnosability](<diagnosability.quality-attribute.spec-report.md>)

[Edit source](<../../../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md#diagnostic-severity-is-available-without-styling>)
