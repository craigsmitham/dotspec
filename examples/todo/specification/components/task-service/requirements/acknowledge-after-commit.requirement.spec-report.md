---
generated_by: dotspec
type: requirement-spec-report
---
# Acknowledge after commit

[Task service](<../task-service.component.spec-report.md>) · [View source][subject-source]

## Statement

The task service must commit a change to the local database before acknowledging it to the browser.

## Applies to

**Targets:** [Task service](<../task-service.component.spec-report.md>).

## System verification plan

Use the assertions in the [restart test](<../../../../tests/acceptance.py>) to check that an acknowledged change survives a fresh process.

[View source for “System verification plan”](<../../../../app/specs/requirements/acknowledge-after-commit.requirement.spec.md#system-verification-plan>)

## Rationale

A successful response tells the person that it is safe to close the task list.

## Related specifications

**Quality attribute** · [Durability](<../../../quality/durability/durability.quality-attribute.spec-report.md>)

[subject-source]: <../../../../app/specs/requirements/acknowledge-after-commit.requirement.spec.md#acknowledge-after-commit>
