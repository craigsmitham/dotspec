---
generated_by: dotspec
type: requirement-spec-report
documents: "../../../../components/task-service/specs/requirements/acknowledge-after-commit.requirement.spec.md#acknowledge-after-commit"
---
# Acknowledge after commit

[Task service](<../task-service.component.spec-report.md>)

## Statement

The task service must commit a change to the local database before acknowledging it to the browser.

## Applies to

**Targets:** [Task service](<../task-service.component.spec-report.md>).

## Verification plan

Use the assertions in the [restart test](<../../../../tests/acceptance.py>) to check that an acknowledged change survives a fresh process.

## Rationale

A successful response tells the person that it is safe to close the task list.

## Related specifications

**Quality attribute** · [Durability](<../../../quality/durability/durability.quality-attribute.spec-report.md>)

[Edit source](<../../../../components/task-service/specs/requirements/acknowledge-after-commit.requirement.spec.md#acknowledge-after-commit>)
