---
type: requirement-spec
id: keep-recorded-work-available
verification:
  - name: Exercise the forms across list views
    source: ../../README.md#run-the-app
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

## Evidence

The [form acceptance test](../../tests/acceptance.py) returns to Open, Waiting and Completed while the service is still running, after entries have been completed, reopened and deleted. Each view shows exactly the entries in that state, once each. After a waiting entry is deleted, a new entry with the same title gets a different identity, and the deleted entry’s old Resume form reports that it is no longer on the list. [Preserve tasks across restart](restart.requirement.spec.md#evidence) covers a return after the service stops.
