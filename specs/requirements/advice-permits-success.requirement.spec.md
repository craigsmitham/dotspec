---
type: requirement
id: advice-permits-success
---

# Advice permits a successful check

## Applies to

Project-wide and focused specification checks that complete without structural
errors or operational failures.

## Statement

When a completed check reports only advisory findings, dotspec must return
exit status zero.

## Rationale

Unfinished specifications remain useful working material. Suggestions should
not prevent contributors from checking a draft.

## Verification

Check structurally valid drafts that produce advisory findings. Verify exit
status zero and the presence of the advice in both text and JSON reports.
Exercise project-wide checks and focused checks for each applicable type.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
provides the conceptual context. The
[report contract](../../README.md#check-specifications) defines other exit statuses.
