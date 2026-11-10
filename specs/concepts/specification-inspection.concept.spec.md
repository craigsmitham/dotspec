---
type: concept
id: specification-inspection
---

# Specification Inspection

## Purpose

Make structural problems and missing-content advice visible within an explicit
inspection scope, so contributors can decide what to revise.

## State

Each inspection produces a report containing the selected scope, supported types,
discovery mode, examined-file counts, additional context sources, diagnostics,
and whether the scan completed. Diagnostics distinguish structural errors from
writing advice and locate the relevant source where possible.

This is the state of one run. A later run produces a new report from its inputs;
there is no stored finding lifecycle or acknowledgement state. The
[check report contract](../../README.md#check-specifications) defines report
fields, diagnostic identifiers, and exit statuses.

Quality obligations are defined by
[textual diagnostic severity](../quality/diagnostic-severity.quality.spec.md),
[repeatable reports](../quality/repeatable-inspection.quality.spec.md),
and [source preservation](../quality/inspection-preserves-source.quality.spec.md).

## Actions

### Inspect a project

Select an existing project directory and inspect supported specifications in its
subtree, consulting enclosing system specifications when needed. Report findings
and scan completion. An empty completed scan explicitly reports that no
supported specifications were found.

### Inspect one specification

Select a supported specification by type, local slug, and project directory.
Inspect it and its required system context using the same diagnostic rules as
the project check. Other well-placed sibling specifications are outside this
focused inspection.

### Read findings and inspect again

Use source locations to find problems, edit the source, and request another
inspection. Within a malformed file, only the first structural error is
reported and dependent checks are skipped; other discovered files continue to
be inspected. Fixing one error can therefore expose another on the next run.
Advice does not fail a check. An access or discovery failure leaves the scan
incomplete, separately from structural findings in readable files.

## Operational principle

Before committing a change, inspect the project and read the findings. Correct
an invalid specification and inspect again to see the current findings. A
complete report shows that the declared scan finished, even if structural
errors remain. An incomplete report exposes an operational gap that must be
resolved before treating the scan as complete. Neither result establishes
requirements completeness, stakeholder agreement, or satisfaction.

The [check-before-committing use case](../use-cases/check-specifications.use-case.spec.md)
places this behavior in a contributor's workflow.
