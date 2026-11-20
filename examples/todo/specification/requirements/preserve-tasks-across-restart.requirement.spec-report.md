---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart"
---
# Preserve tasks across restart

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry's identity, title and completion state.
Deleted entries must stay absent.

## Applies to

**Targets:** [Return to my list](<../use-cases/return-to-my-list.use-case.spec-report.md>).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Count missing, unexpected or changed entries compared with the last acknowledged list.

### Target

Zero discrepancies across the two views.

## Refines

[Keep recorded work available](<keep-recorded-work-available.requirement.spec-report.md>).

## Relies on

[Same local storage](<../assumptions.md#same-local-storage>).

## Supports

[Less reconstruction](<../jobs/resume-interrupted-work.job.spec-report.md#less-reconstruction>).

## Evidence

The [restart acceptance test](<../../tests/acceptance.py>) compares the acknowledged list
with the entries served after restarting on the same database.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.

## Related specifications

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

**Quality attribute**

- [Durability](<../quality/durability/durability.quality-attribute.spec-report.md>)
