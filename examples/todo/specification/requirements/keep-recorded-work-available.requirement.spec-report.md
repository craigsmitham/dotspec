---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/availability.requirement.spec.md#keep-recorded-work-available"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Requirements](<../personal-task-list.system.spec-report.md#requirements>)

# Keep recorded work available

## Statement

When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). An intentionally deleted entry must remain absent from every view.

## Applies to

**Targets:** [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

## Criterion

If the books are completed and the parcel is open, returning to Open shows just the parcel. After completing the parcel too, Open reports that it has no tasks. A failed read must report the failure so the owner can distinguish it from an empty list.

## Supports

[Less time retracing my steps](<../jobs/take-care-of-an-errand-i-have-to-leave-for-later.job.spec-report.md#less-time-retracing-my-steps>).

## Derived from

[Keep completed tasks until I remove them](<../business-rules/keep-completed-tasks-until-i-remove-them.business-rule.spec-report.md>).

## Rationale

The owner can return after closing a tab while the service is still running. A waiting errand that also appears in Open would give them two different answers about whose action comes next.

## Evidence

The [form acceptance test](<../../tests/acceptance.py>) returns to Open, Waiting and Completed while the service is still running, after entries have been completed, reopened and deleted. Each view shows exactly the entries in that state, once each. After a waiting entry is deleted, a new entry with the same title gets a different identity, and the deleted entry’s old Resume form reports that it is no longer on the list. [Preserve tasks across restart](<preserve-tasks-across-restart.requirement.spec-report.md#evidence>) covers a return after the service stops.

## Verification

- [Review errands after a break and distinguish a failed read](<../../tests/acceptance.py#L148>)

- [Exercise the forms across list views](<../../README.md#run-the-app>)

No matching result recorded.

## Related specifications

**Feature**

- [Follow up when someone owes me a reply](<../features/follow-up-when-someone-owes-me-a-reply.feature.spec-report.md>)
- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

**Refined by**

- [Preserve tasks across restart](<preserve-tasks-across-restart.requirement.spec-report.md>)

---

[View source](<../../specs/requirements/availability.requirement.spec.md#keep-recorded-work-available>)
