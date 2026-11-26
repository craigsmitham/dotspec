---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Requirements](<../personal-task-list.system.spec-report.md#requirements>)

# Preserve tasks across restart

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry’s identity, title and state. For a waiting entry, the expected reply and follow-up date must also survive. Deleted entries must stay absent.

## Applies to

**Targets:** [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Use the acknowledged changes to establish which entries should survive and in which view. Match each identity across all three views after restart. Count a discrepancy for each identity that is missing, appears unexpectedly or more than once, or has a different title, state or saved waiting detail. The cue for the saved date is assessed on the day of return under [Show which replies need follow-up](<show-which-replies-need-follow-up.requirement.spec-report.md>).

### Target

Zero discrepancies across Open, Waiting and Completed.

## Refines

[Keep recorded work available](<keep-recorded-work-available.requirement.spec-report.md>).

## Relies on

[Same local storage](<../assumptions.md#same-local-storage>).

## Supports

[Less reconstruction](<../jobs/see-personal-commitments-through.job.spec-report.md#less-reconstruction>).

## Evidence

The [restart acceptance test](<../../tests/acceptance.py>) records entries in Open, Waiting and Completed through the forms, kills the service process without a clean shutdown and starts it again on the same database. It compares each view’s identities and titles with the acknowledged list. A missing, duplicated or unexpected entry fails the run, as would the return of either entry deleted before the stop. For both waiting entries it also compares the expected reply and follow-up date. The test does not stop the service during a change or simulate a power loss.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.

## Companions

**Verified by**

- [acceptance.py › verify](<../../tests/acceptance.py#L149>) — Restart acceptance test; kills the service and compares each view's identities after it starts again

## Verification

Procedure linked.

No matching result recorded.

## Related specifications

**Basis for**

- [Acknowledge after commit](<../components/task-service/requirements/acknowledge-after-commit.requirement.spec-report.md>)

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

**Quality attribute**

- [Durability](<../quality/durability/durability.quality-attribute.spec-report.md>)

---

[View source](<../../specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>)
