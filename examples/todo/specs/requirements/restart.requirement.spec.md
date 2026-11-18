---
type: requirement
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

[Same local storage](../system.spec.md#same-local-storage).

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

## Verification plan

### Method

The [restart acceptance test](../../tests/acceptance.py) submits the public forms,
stops the process, starts a new process on the same file and compares both views.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.
