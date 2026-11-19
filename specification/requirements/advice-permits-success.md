---
type: specification-section
generated_by: dotspec
subject_kind: requirement
---
# Advice permits a successful check

[dotspec](<../../SPEC.md>) · [Contents](<../index.md>)

## Statement

When a completed check reports only advisory findings, dotspec must return
exit status zero.

## Applies to

Project-wide and focused specification checks that complete without structural
errors or operational failures.

## Rationale

Unfinished specifications remain useful working material. Suggestions should
not prevent contributors from checking a draft.

## Verification plan

Check structurally valid drafts that produce advisory findings. Verify exit
status zero and the presence of the advice in both text and JSON reports.
Exercise project-wide checks and focused checks for each applicable type.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.md>)
provides the conceptual context. The
[report contract](<../../docs/reference/inspection.md#check-specifications>) defines other exit statuses.

[Edit source](<../../specs/requirements/advice-permits-success.requirement.spec.md#advice-permits-a-successful-check>)
