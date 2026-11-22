---
generated_by: dotspec
type: requirement-spec-report
source: "../../../../app/specs/requirements/acknowledge-after-commit.requirement.spec.md#acknowledge-after-commit"
---
# Acknowledge after commit

[Task service](<../task-service.component.spec-report.md>)

## Statement

The task service must commit a change to the local database before acknowledging it to the browser.

## Applies to

**Targets:** [Task service](<../task-service.component.spec-report.md>).

### Conditions

The browser submits an add, edit, state change or deletion. The service acknowledges acceptance by redirecting the browser to the list. A redirect after a failed write, or while a write can still roll back, violates this requirement.

## Derived from

[Preserve tasks across restart](<../../../requirements/preserve-tasks-across-restart.requirement.spec-report.md>). The service’s redirect is the point after which the owner can stop the process and rely on that promise.

## Rationale

A successful response tells the person that it is safe to close the task list.

## Evidence

Inspect the [change handler](<../../../../app/main.go>) to check that database success precedes the redirect. The [restart test](<../../../../tests/acceptance.py>) checks what the browser receives on the next visit after acknowledged changes.

## Related specifications

**Quality attribute**

- [Durability](<../../../quality/durability/durability.quality-attribute.spec-report.md>)
