---
type: requirement-spec
quality_attributes: [durability]
verification:
  - name: Restart acceptance test
    source: ../../tests/acceptance.py#verify
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

**Targets:** [My checklist](../personal-task-list.system.spec.md#my-checklist).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Use the acknowledged changes to establish which entries should survive and in which view. Match each identity across all three views after restart. Count a discrepancy for each identity that is missing, appears unexpectedly or more than once, or has a different title, state or saved waiting detail. The cue for the saved date is assessed on the day of return under [Show which replies need follow-up](follow-up.requirement.spec.md).

### Target

Zero discrepancies across Open, Waiting and Completed.

## Evidence

The [restart acceptance test](../../tests/acceptance.py) records entries in Open, Waiting and Completed through the forms, kills the service process without a clean shutdown and starts it again on the same database. It compares each view’s identities and titles with the acknowledged list. A missing, duplicated or unexpected entry fails the run, as would the return of either entry deleted before the stop. For both waiting entries it also compares the expected reply and follow-up date. The test does not stop the service during a change or simulate a power loss.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.
