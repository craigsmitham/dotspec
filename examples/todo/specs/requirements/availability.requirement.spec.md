---
type: requirement-spec
id: keep-recorded-work-available
---
# Keep recorded work available

## Statement

When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). An intentionally deleted entry must remain absent from every view.

## Applies to

**Targets:** [My checklist](../personal-task-list.system.spec.md#my-checklist).

## Derived from

[Completion preserves a commitment](../business-rules/completion.business-rule.spec.md).

## Rationale

The owner can return after closing a tab while the service is still running. A waiting errand that also appears in Open would give them two different answers about whose action comes next.
