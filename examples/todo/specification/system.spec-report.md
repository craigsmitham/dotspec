---
generated_by: dotspec
type: system-spec-report
documents: "../specs/system.spec.md#personal-task-list"
---
# Personal task list

## Purpose

Help a person pick up unfinished work after an interruption.

## Scope

One person's list on one computer. The browser submits forms to a local service.

---

[Requirements](<#requirements>) · [Requirements review](<requirements-review.report.md>)

## Features

- [Pick up where I left off](<features/pick-up-where-i-left-off.feature.spec-report.md>)

## Goals

### Resume unfinished work

<!-- dotspec-doc {"type":"goal-spec-report","documents":"../specs/system.spec.md#resume-unfinished-work"} -->

Reduce the effort of working out what still needs attention.

**Supported by** · [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>)

## Environment

The person may close the browser or stop the service between work sessions.

## Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

## Assumptions

### Same local storage

<!-- dotspec-doc {"type":"assumption-spec-report","documents":"../specs/system.spec.md#same-local-storage"} -->

The person returns to the same writable database file.

**Relied on by** · [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Components

- [Task service](<components/task-service.component.spec-report.md>)

## Project context

The [walkthrough](<../README.md>) identifies the setup, checks and validation exercise.

## Concept uses

- [My checklist](<concept-uses/my-checklist.concept-use.spec-report.md>)

## Provided interfaces

- [Task list forms](<interfaces/task-list-forms.interface.spec-report.md>)

## Quality attributes

| Quality attribute | Description |
| --- | --- |
| [Durability](<quality-attributes.spec-report.md#durability>) | Preservation of acknowledged work between local sessions. |

## Validation plan

### Subject

A person returning to unfinished work using this build.

### Method

Invite someone to record two real errands and mark one finished. After a break,
restart the app with the same database. Ask what still needs attention. Then tell
them that the finished errand needs a follow-up and ask them to resume it.
Let them choose the controls without coaching.

### Success criteria

The participant identifies the outstanding errand and resumes the other without
entering it again. Ask what they expected to find on returning.

### Evidence to retain

With consent, record the build, task outcome, any help given and the participant's
explanation. Review mismatches before deciding whether the workflow meets the need.

## Business rules

- [Completion preserves a commitment](<business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>)

## Concepts

- [Checklist](<concepts/checklist.concept.spec-report.md>)

## Jobs

- [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>)

## Quality models

- [Personal data quality](<quality-models/personal-data-quality.quality-model.spec-report.md>)

## Requirements

- [Keep recorded work available](<use-cases/return-to-my-list.use-case.spec-report.md#keep-recorded-work-available>)
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

## Use cases

- [Return to my list](<use-cases/return-to-my-list.use-case.spec-report.md>)

## Related specifications

**Applies here** · [Completion preserves a commitment](<business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>)

[Edit source](<../specs/system.spec.md#personal-task-list>)
