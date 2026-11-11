---
type: quality
id: diagnostic-severity
quality_attributes:
  - diagnosability
  - accessibility
---

# Diagnostic severity is available without styling

## Applies to

Human-readable reports from project-wide and focused specification checks.

## Statement

Every diagnostic must identify its severity using a text label independently
of terminal styling.

## Criterion

Use the text-label rule in [Statement](#statement), with supported severity values
defined by the [check report contract](../../README.md#check-specifications).

## Measurement plan

Render reports containing each supported diagnostic severity. Capture output
with terminal-style environment settings and in redirected-output conditions,
including with color disabled. Compare each diagnostic with
[Criterion](#criterion). Retain captured output on test failure.

The [CLI quality-contract tests](../../internal/cli/quality_contracts_test.go)
check this criterion. The obligation remains defined in this document.

## Rationale

Contributors need to distinguish diagnostic severities when reading plain-text
logs or using tools that do not convey color. This criterion addresses that
specific need; it does not establish overall accessibility.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
describes the report and inspection behavior. The
[check report contract](../../README.md#check-specifications) defines the
supported severity values and report structure.
