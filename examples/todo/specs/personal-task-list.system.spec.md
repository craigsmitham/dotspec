---
type: system-spec
---
# Personal task list

## Description

A personal list of errands the owner can act on, replies they are waiting for and commitments they have settled.

## Purpose

Help one person pick up an errand after an interruption and recognize when to chase an outstanding reply.

## Overview

Each errand sits in one of three views. Open holds work the owner can advance, Waiting holds errands that depend on someone else's reply, and Completed holds commitments the owner has settled. Suppose a guest needs a working sink by Friday. Sending the plumber a photo finishes today's action, so the owner moves the repair to Waiting. They still owe the repair. The [job of seeing personal commitments through](jobs/personal-commitments.job.spec.md) explains what the owner needs to establish when they return.

## Goals

### Resume unfinished work

Reduce the effort of working out what still needs attention.

## Stakeholders

### List owner

The person who keeps the list relies on it to pick up unfinished errands after an interruption. They decide which entries can be discarded.

## Scope

One person records commitments on one computer and reviews them in Open, Waiting and Completed. They choose what reply to watch for and the day to check back. They receive the reply and judge the result outside the list.

## Out of scope

Sending reminders, reading messages, managing delivery deadlines, sharing lists, synchronizing devices and recovering from loss of the database file.

## Environment

The browser and service run on the same computer. The owner can close either between visits and return on a different local calendar day. The [follow-up requirement](requirements/follow-up.requirement.spec.md#conditions) explains what they should see when they open Waiting again.

## Assumptions

### Same local storage

The person returns to the same writable database file.

#### Basis

The owner names the database file when starting the service, or accepts `todo.sqlite` in its working directory. Starting it again the same way reopens the same list. Another path opens a different list, which the restart requirement does not cover.

#### Applicability

Each return to the list on the owner's computer, including after the service restarts.

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

The forms the list owner uses to review and update commitments.

Open offers work the owner can advance. Waiting shows the expected reply and when to chase it. Completed holds settled errands the owner may reopen. After an accepted change the browser returns to the view where the person submitted it.

#### Contracts

[Form template](../app/list.html) supplies the actions and fields.

#### Adopts

##### Redirect after an accepted form submission

[303 See Other](https://www.rfc-editor.org/rfc/rfc9110.html#section-15.4.4).

###### Edition

RFC 9110, June 2022.

###### Applicability

After accepting a change through a task form, Task service redirects the browser to the originating list view.
