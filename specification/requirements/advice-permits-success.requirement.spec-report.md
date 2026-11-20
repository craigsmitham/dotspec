---
generated_by: dotspec
type: requirement-spec-report
---
# Advice permits a successful check

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Statement

When a completed check reports only advisory findings, dotspec must return
exit status zero.

## Applies to

Project-wide and focused specification checks that complete without structural
errors or operational failures.

## System verification plan

Check structurally valid drafts that produce advisory findings. Verify exit
status zero and the presence of the advice in both text and JSON reports.
Exercise project-wide checks and focused checks for each applicable type.

[View source for “System verification plan”](<../../specs/requirements/advice-permits-success.requirement.spec.md#system-verification-plan>)

## Rationale

Unfinished specifications remain useful working material. Suggestions should
not prevent contributors from checking a draft.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
provides the conceptual context. The
[report contract](<../../docs/reference/inspection.md#check-specifications>) defines other exit statuses.

[subject-source]: <../../specs/requirements/advice-permits-success.requirement.spec.md#advice-permits-a-successful-check>
