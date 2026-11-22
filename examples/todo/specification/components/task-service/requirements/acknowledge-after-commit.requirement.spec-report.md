---
generated_by: dotspec
type: requirement-spec-report
source: "../../../../app/specs/requirements/acknowledge-after-commit.requirement.spec.md#acknowledge-after-commit"
---
# Acknowledge after commit

[Task service](<../task-service.component.spec-report.md>)

## Statement

The task service must record each change in a single commit to the local database and acknowledge it to the browser only after that commit succeeds.

## Applies to

**Targets:** [Task service](<../task-service.component.spec-report.md>).

### Conditions

The browser submits an add, edit, state change or deletion. The service acknowledges acceptance by redirecting the browser to the list.

## Criterion

A redirect after a failed write, or while the write can still roll back, violates this requirement. So does an error response after which any part of the change is saved, such as a failed Complete that leaves a waiting errand in Open without its expected reply. An error response that leaves every entry as it was, including a waiting errand’s reply and follow-up day, meets it.

## Derived from

[Preserve tasks across restart](<../../../requirements/preserve-tasks-across-restart.requirement.spec-report.md>). The service’s redirect is the point after which the owner can stop the process and rely on that promise.

## Rationale

A successful response tells the person that it is safe to close the task list.

## Evidence

The [form acceptance test](<../../../../tests/acceptance.py>) submits Complete for a waiting errand and makes the database reject the write that would mark it completed. A second case submits Resume and makes the database reject removal of the waiting details. Each submission returns the error “Could not access the task list.” without a redirect. Waiting still shows the same entry, expected reply and follow-up day. The restart scenario in that test checks the saved entries after acknowledged changes. Inspect the [change handler](<../../../../app/main.go>) to verify that every path to the success redirect passes through a successful commit.

## Related specifications

**Quality attribute**

- [Durability](<../../../quality/durability/durability.quality-attribute.spec-report.md>)
