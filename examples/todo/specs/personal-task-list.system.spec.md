---
type: system-spec
---
# Personal task list

## Description

A personal checklist for errands I want to remember and work I have already finished.

## Purpose

Give one person a place to put a task down now and find out what is left when they return.

## Overview

At breakfast I remember to return the library books and mail a parcel. I write both down. That afternoon I return the books and cross them off. The next morning the parcel is still in Open. If I find another library book in my bag, I can reopen the same task. [Take care of an errand I have to leave for later](jobs/errands.job.spec.md).

Sometimes the next action belongs to someone else. The [follow-up example](features/follow-up.feature.spec.md) shows how I keep track of a reply.

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

I can cross off one errand, leave another open and put a task back if there is more to do. Two errands with the same title can be at different stages.

#### Definition

[Checklist](concepts/checklist.concept.spec.md)

#### Supports

[Take care of an errand I have to leave for later](jobs/errands.job.spec.md) and [See personal commitments through](jobs/personal-commitments.job.spec.md).

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
