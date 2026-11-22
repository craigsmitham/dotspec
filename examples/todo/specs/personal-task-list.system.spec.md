---
type: system-spec
---
# Personal task list

## Description

A personal list of errands the owner can act on, replies they are waiting for and commitments they have settled.

## Framing

Suppose a guest needs a working sink by Friday. Sending the plumber a photo finishes today's action. The owner still owes the repair. The [job of seeing personal commitments through](jobs/personal-commitments.job.spec.md) explains what the owner needs to establish when they return.

## Purpose

Help one person pick up an errand after an interruption and recognize when to chase an outstanding reply.

## Goals

### Resume unfinished work

Reduce the effort of working out what still needs attention.

## Stakeholders

The person who keeps the list also decides which entries can be discarded.

## Scope

One person records commitments on one computer and reviews them in Open, Waiting and Completed. They choose what reply to watch for and the day to check back. They receive the reply and judge the result outside the list.

## Out of scope

Sending reminders, reading messages, managing delivery deadlines, sharing lists, synchronizing devices and recovering from loss of the database file.

## Environment

The browser and service run on the same computer. The person can close either between visits. The service uses that computer’s local calendar when the person loads Waiting.

## Assumptions

### Same local storage

The person returns to the same writable database file.

## Concept uses

### My checklist

#### Description

Each entry follows one commitment through the owner's next action, an outstanding reply or settlement. Two errands with the same title can be at different stages.

#### Definition

[Checklist](concepts/checklist.concept.spec.md)

#### Supports

[See personal commitments through](jobs/personal-commitments.job.spec.md).

## Provided interfaces

### Task list forms

Open offers work the owner can advance. Waiting shows the expected reply and when to chase it. Completed holds settled errands the owner may reopen. After an accepted change the browser returns to the view where the person submitted it.

#### Contracts

[Form template](../app/list.html) supplies the actions and fields.
