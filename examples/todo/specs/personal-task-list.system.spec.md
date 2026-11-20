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
