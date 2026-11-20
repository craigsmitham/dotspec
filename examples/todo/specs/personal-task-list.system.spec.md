---
type: system-spec
---
# Personal task list

## Description

A single-user task list that stores open and completed tasks on one computer.

## Purpose

Help a person pick up unfinished work after an interruption.

## Scope

One person's list on one computer. The browser submits forms to a local service.

## Goals

### Resume unfinished work

Reduce the effort of working out what still needs attention.

## Environment

The person may close the browser or stop the service between work sessions.

## Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

## Assumptions

### Same local storage

The person returns to the same writable database file.

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Capabilities

### Local task storage

#### Description

Retains accepted task changes between sessions on the same computer.

#### Provided by

[Task service](../components/task-service/specs/task-service.component.spec.md).

## Features

### Pick up where I left off

#### Description

Reopens the task list with each task in its last saved state.

#### Groups

Follow the interaction in [Return to my list](use-cases/resume-list.use-case.spec.md).
For the storage guarantee and its acceptance test, see
[Preserve tasks across restart](requirements/restart.requirement.spec.md).

## Project context

The [walkthrough](../README.md) identifies the setup, checks and validation exercise.

## Concept uses

### My checklist

#### Description

The checklist used to track one person’s open and completed tasks.

#### Definition

[Checklist](concepts/checklist.concept.spec.md)

#### Supports

[Resume interrupted work](jobs/resume.job.spec.md).

## Provided interfaces

### Task list forms

The browser presents open and completed tasks and accepts changes through HTML forms.

#### Contracts

[Form template](../app/list.html) supplies the actions and fields.

## Quality attributes

The task list is stored locally. Its quality requirements cover what survives between sessions on the same database.

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
