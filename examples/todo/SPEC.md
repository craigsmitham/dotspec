---
generated_by: dotspec
type: system-spec-report
documents: "specs/system.spec.md#personal-task-list"
---
# Personal task list

## Purpose

Help a person pick up unfinished work after an interruption.

## Scope

One person's list on one computer. The browser submits forms to a local service.

---

[Requirements](<#requirements>) · [Requirements review](<specification/review.md>)

## Features

- [Pick up where I left off](<specification/features/pick-up-where-i-left-off.md>)

## Goals

### Resume unfinished work

<!-- dotspec-doc {"type":"goal-spec-report","documents":"specs/system.spec.md#resume-unfinished-work"} -->

Reduce the effort of working out what still needs attention.

**Supported by** · [Resume interrupted work](<specification/jobs/resume-interrupted-work.md>)

## Environment

The person may close the browser or stop the service between work sessions.

## Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

## Assumptions

### Same local storage

<!-- dotspec-doc {"type":"assumption-spec-report","documents":"specs/system.spec.md#same-local-storage"} -->

The person returns to the same writable database file.

**Relied on by** · [Preserve tasks across restart](<specification/requirements/preserve-tasks-across-restart.md>)

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Components

- [Task service](<specification/components/task-service.md>)

## Project context

The [walkthrough](<README.md>) identifies the setup, checks and validation exercise.

## Concept uses

- [My checklist](<specification/concept-uses/my-checklist.md>)

## Provided interfaces

- [Task list forms](<specification/interfaces/task-list-forms.md>)

## Quality attributes

| Quality attribute | Description |
| --- | --- |
| [Durability](<specification/quality-attributes.md#durability>) | Preservation of acknowledged work between local sessions. |

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

- [Completion preserves a commitment](<specification/business-rules/completion-preserves-a-commitment.md>)

## Concepts

- [Checklist](<specification/concepts/checklist.md>)

## Jobs

- [Resume interrupted work](<specification/jobs/resume-interrupted-work.md>)

## Quality models

- [Personal data quality](<specification/quality-models/personal-data-quality.md>)

## Requirements

- [Keep recorded work available](<specification/use-cases/return-to-my-list.md#keep-recorded-work-available>)
- [Preserve tasks across restart](<specification/requirements/preserve-tasks-across-restart.md>)

## Use cases

- [Return to my list](<specification/use-cases/return-to-my-list.md>)

## Related specifications

**Applies here** · [Completion preserves a commitment](<specification/business-rules/completion-preserves-a-commitment.md>)

[Edit source](<specs/system.spec.md#personal-task-list>)
