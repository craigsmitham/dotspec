---
type: system-spec
---
# Personal task list

## Description

A single-user task list that stores open and completed tasks on one computer.

## Framing

An errand can need another action after it seemed finished. The [job of seeing personal commitments through](jobs/personal-commitments.job.spec.md) gives the reason to keep both open and completed commitments available when the person returns.

## Purpose

Help a person pick up unfinished work after an interruption.

## Goals

### Resume unfinished work

Reduce the effort of working out what still needs attention.

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Scope

The list tracks one person’s open and completed tasks on one computer.

## Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

## Environment

The person may close the browser or stop the service between work sessions.

## Assumptions

### Same local storage

The person returns to the same writable database file.

## Concept uses

### My checklist

#### Description

The checklist used to track one person’s open and completed tasks.

#### Definition

[Checklist](concepts/checklist.concept.spec.md)

#### Supports

[See personal commitments through](jobs/personal-commitments.job.spec.md).

## Provided interfaces

### Task list forms

The browser presents open and completed tasks and accepts changes through HTML forms.

#### Contracts

[Form template](../app/list.html) supplies the actions and fields.
