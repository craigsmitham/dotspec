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

## Rationale

A successful response tells the person that it is safe to close the task list.

## Evidence

The [restart test](<../../../../tests/acceptance.py>) exercises persistence of acknowledged changes.

## Related specifications

**Quality attribute** · [Durability](<../../../quality/durability/durability.quality-attribute.spec-report.md>)
