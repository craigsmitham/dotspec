---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/advice-permits-success.requirement.spec.md#advice-permits-a-successful-check"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../requirements.md>)

# Advice permits a successful check

[View source](<../../specs/requirements/advice-permits-success.requirement.spec.md#advice-permits-a-successful-check>)

## Statement

When a completed check reports only advisory findings, dotspec must return
exit status zero.

## Applies to

Project-wide and focused specification checks that complete without structural
errors or operational failures.

## Rationale

Unfinished specifications remain useful working material. Suggestions should
not prevent contributors from checking a draft.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
provides the conceptual context. The
[report contract](<../../docs/reference/inspection.md#check-specifications>) defines other exit statuses.
