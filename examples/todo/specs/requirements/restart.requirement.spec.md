---
type: requirement-spec
quality_attributes: [durability]
---
# Preserve tasks across restart

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry's identity, title and completion state.
Deleted entries must stay absent.

## Refines

[Keep recorded work available](../use-cases/resume-list.use-case.spec.md#keep-recorded-work-available).

## Relies on

[Same local storage](../personal-task-list.system.spec.md#same-local-storage).

## Supports

[Less reconstruction](../jobs/resume.job.spec.md#less-reconstruction).

## Applies to

**Targets:** [Return to my list](../use-cases/resume-list.use-case.spec.md).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Count missing, unexpected or changed entries compared with the last acknowledged list.

### Target

Zero discrepancies across the two views.

## Evidence

The [restart acceptance test](../../tests/acceptance.py) compares the acknowledged list
with the entries served after restarting on the same database.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.
