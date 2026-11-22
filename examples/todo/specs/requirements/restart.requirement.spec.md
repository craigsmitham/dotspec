---
type: requirement-spec
quality_attributes: [durability]
---
# Preserve tasks across restart

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry’s identity, title and state. For a waiting entry, the expected reply and follow-up date must also survive. Deleted entries must stay absent.

## Refines

[Keep recorded work available](availability.requirement.spec.md).

## Relies on

[Same local storage](../personal-task-list.system.spec.md#same-local-storage).

## Supports

[Less reconstruction](../jobs/personal-commitments.job.spec.md#less-reconstruction).

## Applies to

**Targets:** [Return to my list](../use-cases/resume-list.use-case.spec.md).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Match entries by identity against the last acknowledged list. Count missing, unexpected or changed entries across all views. Compare the recorded follow-up date. The cue for that date is assessed on the day of return under [Show which replies need follow-up](follow-up.requirement.spec.md).

### Target

Zero discrepancies across Open, Waiting and Completed.

## Evidence

The [restart acceptance test](../../tests/acceptance.py) compares the acknowledged list
with the entries served after restarting on the same database.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.
