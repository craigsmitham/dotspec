---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart"
---
# Preserve tasks across restart

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry’s identity, title and state. For a waiting entry, the expected reply and follow-up date must also survive. Deleted entries must stay absent.

## Applies to

**Targets:** [Decide what to do about a waiting errand](<../use-cases/decide-what-to-do-about-a-waiting-errand.use-case.spec-report.md>).

### Conditions

A local file on functioning storage, with each form submission acknowledged
before the process stops.

## Criterion

### Measure

Match entries by identity against the last acknowledged list. Count missing, unexpected or changed entries across all views. Compare the recorded follow-up date. The cue for that date is assessed on the day of return under [Show which replies need follow-up](<show-which-replies-need-follow-up.requirement.spec-report.md>).

### Target

Zero discrepancies across Open, Waiting and Completed.

## Refines

[Keep recorded work available](<keep-recorded-work-available.requirement.spec-report.md>).

## Relies on

[Same local storage](<../assumptions.md#same-local-storage>).

## Supports

[Less reconstruction](<../jobs/see-personal-commitments-through.job.spec-report.md#less-reconstruction>).

## Evidence

The [restart acceptance test](<../../tests/acceptance.py>) compares the acknowledged list
with the entries served after restarting on the same database.

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.

## Related specifications

**Basis for**

- [Acknowledge after commit](<../components/task-service/requirements/acknowledge-after-commit.requirement.spec-report.md>)

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

**Quality attribute**

- [Durability](<../quality/durability/durability.quality-attribute.spec-report.md>)
