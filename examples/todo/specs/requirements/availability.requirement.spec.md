---
type: requirement-spec
id: keep-recorded-work-available
verification:
  - name: Review errands after a break and distinguish a failed read
    source: ../../tests/acceptance.py#verify_errands
  - name: Exercise the forms across list views
    source: ../../README.md#run-the-app
---
# Keep recorded work available

## Statement

When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). An intentionally deleted entry must remain absent from every view.

## Supports

[Less time retracing my steps](../jobs/errands.job.spec.md#less-time-retracing-my-steps).

## Criterion

If the books are completed and the parcel is open, returning to Open shows just the parcel. After completing the parcel too, Open reports that it has no tasks. A failed read must report the failure so the owner can distinguish it from an empty list.

## Applies to

**Targets:** [My checklist](../personal-task-list.system.spec.md#my-checklist).

## Derived from

[Keep completed tasks until I remove them](../business-rules/completion.business-rule.spec.md).

## Rationale

The owner can return after closing a tab while the service is still running. A waiting errand that also appears in Open would give them two different answers about whose action comes next.

## Evidence

The [form acceptance test](../../tests/acceptance.py) returns to Open, Waiting and Completed while the service is still running, after entries have been completed, reopened and deleted. Each view shows exactly the entries in that state, once each. After a waiting entry is deleted, a new entry with the same title gets a different identity, and the deleted entry’s old Resume form reports that it is no longer on the list. [Preserve tasks across restart](restart.requirement.spec.md#evidence) covers a return after the service stops.
