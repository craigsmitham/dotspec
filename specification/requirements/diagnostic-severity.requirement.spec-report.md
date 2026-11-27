---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md#diagnostic-severity-is-available-without-styling"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../dotspec.system.spec-report.md#requirements>)

# Diagnostic severity is available without styling

## Statement

Every diagnostic must identify its severity using a text label independently
of terminal styling.

## Applies to

Human-readable reports from project-wide and focused specification checks.

## Criterion

Use the text-label rule in [Statement](<#statement>), with supported severity values
defined by the [check report contract](<../../docs/reference/inspection.md#check-specifications>).

## Rationale

Contributors need to distinguish diagnostic severities when reading plain-text
logs or using tools that do not convey color. This criterion addresses that
specific need; it does not establish overall accessibility.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
describes the report and inspection behavior. The
[check report contract](<../../docs/reference/inspection.md#check-specifications>) defines the
supported severity values and report structure.

**Quality attribute**

- [Accessibility](<../quality/usability/accessibility/accessibility.quality-attribute.spec-report.md>)
- [Diagnosability](<../quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>)

## Companions

**Verified by**

- [quality_contracts_test.go › TestDiagnosticSeverityQualityContract](<../../internal/cli/quality_contracts_test.go#L104>) — CLI quality-contract test for severity without styling

## Verification

Procedure linked.

No matching result recorded.

---

[View source](<../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md#diagnostic-severity-is-available-without-styling>)
