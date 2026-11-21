---
type: requirement-spec
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

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
provides the conceptual context. The
[report contract](../../docs/cli-handbook/reference/inspection.md#check-specifications) defines other exit statuses.
