---
generated_by: dotspec
type: requirement-spec-report
---
# Preserve tasks across restart

[Personal task list](<../personal-task-list.system.spec-report.md>) · [View source][subject-source]

## Statement

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry's identity, title and completion state.
Deleted entries must stay absent.

[View source for “Statement”](<../../specs/requirements/restart.requirement.spec.md#statement>)

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

[View source for “Criterion”](<../../specs/requirements/restart.requirement.spec.md#criterion>)

## System verification plan

### Subject

The running [task list service](<../../app/main.go>) before and after restart on the same database.

### Basis

The [restart statement](<#statement>) and its [criterion](<#criterion>).

### Method

The [restart acceptance test](<../../tests/acceptance.py>) submits the public forms,
stops the process, starts a new process on the same file and compares both views.

[View source for “System verification plan”](<../../specs/requirements/restart.requirement.spec.md#system-verification-plan>)

## Refines

[Keep recorded work available](<../use-cases/return-to-my-list.use-case.spec-report.md#keep-recorded-work-available>).

## Relies on

[Same local storage](<../personal-task-list.system.spec-report.md#same-local-storage>).

## Supports

[Less reconstruction](<../jobs/resume-interrupted-work.job.spec-report.md#less-reconstruction>).

## Requirements verification plan

### Subject

The [restart statement](<#statement>).

### Basis

Clarity, consistency with [same local storage](<../personal-task-list.system.spec-report.md#same-local-storage>),
and whether each promised result can be checked.

### Method

Ask a reviewer to identify when preservation starts, which entries must survive,
and what would count as a discrepancy. Record any ambiguity in the wording.

[View source for “Requirements verification plan”](<../../specs/requirements/restart.requirement.spec.md#requirements-verification-plan>)

## Rationale

A list that loses track of earlier decisions forces the person to reconstruct them.

## Related specifications

**Feature** · [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

**Quality attribute** · [Durability](<../quality/durability/durability.quality-attribute.spec-report.md>)

[subject-source]: <../../specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>
