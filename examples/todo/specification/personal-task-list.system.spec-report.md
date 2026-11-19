---
generated_by: dotspec
type: system-spec-report
documents: "../specs/personal-task-list.system.spec.md#personal-task-list"
---
# Personal task list

## Description

A single-user task list that stores open and completed tasks on one computer.

## Purpose

Help a person pick up unfinished work after an interruption.

## Scope

One person's list on one computer. The browser submits forms to a local service.

---

[Requirements](<#requirements>) · [Requirements review](<requirements-review.report.md>)

## Features

- [Pick up where I left off](<features/pick-up-where-i-left-off.feature.spec-report.md>) — Reopens the task list with each task in its last saved state.

## Goals

### Resume unfinished work

<!-- dotspec-doc {"type":"goal-spec-report","documents":"../specs/personal-task-list.system.spec.md#resume-unfinished-work"} -->

Reduce the effort of working out what still needs attention.

**Supported by** · [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>)

## Environment

The person may close the browser or stop the service between work sessions.

## Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

## Assumptions

### Same local storage

<!-- dotspec-doc {"type":"assumption-spec-report","documents":"../specs/personal-task-list.system.spec.md#same-local-storage"} -->

The person returns to the same writable database file.

**Relied on by** · [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Capabilities

- [Local task storage](<capabilities/local-task-storage.capability.spec-report.md>) — Retains accepted task changes between sessions on the same computer.

## Project context

The [walkthrough](<../README.md>) identifies the setup, checks and validation exercise.

## Concept uses

- [My checklist](<concept-uses/my-checklist.concept-use.spec-report.md>) — The checklist used to track one person’s open and completed tasks.

## Provided interfaces

- [Task list forms](<interfaces/task-list-forms.interface.spec-report.md>) — The browser presents open and completed tasks and accepts changes through HTML forms.

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

- [Completion preserves a commitment](<business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>) — A finished commitment stays on the person's list until they explicitly delete it.

## Components

- [Task service](<components/task-service/task-service.component.spec-report.md>) — The local service that handles browser requests and saves changes to the task list.

## Concepts

- [Checklist](<concepts/checklist.concept.spec-report.md>) — Keep track of which commitments still need attention.

## Jobs

- [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>) — When I put work down, I want to remember what remains so I can pick it up again.

## Quality models

- [Personal data quality](<quality-models/personal-data-quality.quality-model.spec-report.md>) — Select concerns for tools in which a person keeps work between sessions.

## Requirements

- [Keep recorded work available](<use-cases/return-to-my-list.use-case.spec-report.md#keep-recorded-work-available>)
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry's identity, title and completion state. Deleted entries must stay absent.

## Use cases

- [Return to my list](<use-cases/return-to-my-list.use-case.spec-report.md>) — Pick up an outstanding task and bring a finished one back for follow-up.

## Related specifications

**Applies here** · [Completion preserves a commitment](<business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>)

[Edit source](<../specs/personal-task-list.system.spec.md#personal-task-list>)
