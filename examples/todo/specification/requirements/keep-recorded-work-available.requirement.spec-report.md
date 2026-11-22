---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/availability.requirement.spec.md#keep-recorded-work-available"
---
# Keep recorded work available

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Statement

When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). An intentionally deleted entry must remain absent from every view.

## Applies to

**Targets:** [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

## Derived from

[Completion preserves a commitment](<../business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>).

## Rationale

The owner can return after closing a tab while the service is still running. A waiting errand that also appears in Open would give them two different answers about whose action comes next.

## Related specifications

**Refined by**

- [Preserve tasks across restart](<preserve-tasks-across-restart.requirement.spec-report.md>)
