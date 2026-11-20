---
type: requirement-spec
quality_attributes: [durability]
---
# Acknowledge after commit

## Statement

The task service must commit a change to the local database before acknowledging it to the browser.

## Applies to

**Targets:** [Task service](../task-service.component.spec.md).

## Rationale

A successful response tells the person that it is safe to close the task list.

## Verification plan

Use the assertions in the [restart test](../../../tests/acceptance.py) to check that an acknowledged change survives a fresh process.
