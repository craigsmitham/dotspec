---
generated_by: dotspec
type: system-spec-report
source: "../specs/personal-task-list.system.spec.md#personal-task-list"
---
# Personal task list

A personal list of errands the owner can act on, replies they are waiting for and commitments they have settled.

<details>
<summary>On this page</summary>

**Intent**\
[Framing](<#framing>) · [Purpose](<#purpose>) · [Goals](<#goals>) · [Stakeholders](<#stakeholders>) · [Jobs](<#jobs>)

**Context**\
[Scope](<#scope>) · [Environment](<#environment>) · [Assumptions](<#assumptions>) · [Business rules](<#business-rules>)

**Behavior**\
[Features](<#features>) · [Capabilities](<#capabilities>) · [Use cases](<#use-cases>)

[Quality Attributes](<#quality-attributes>)

[Requirements](<#requirements>)

**Design**\
[Concepts](<#concepts>) · [Concept uses](<#concept-uses>) · [Architecture](<#architecture>)

</details>

## Framing

Suppose a guest needs a working sink by Friday. Sending the plumber a photo finishes today's action. The owner still owes the repair. The [job of seeing personal commitments through](<jobs/see-personal-commitments-through.job.spec-report.md>) explains what the owner needs to establish when they return.

## Purpose

Help one person pick up an errand after an interruption and recognize when to chase an outstanding reply.

## Goals

- [Resume unfinished work](<goals.md#resume-unfinished-work>) — Reduce the effort of working out what still needs attention.

## Stakeholders

### List owner

<!-- dotspec-doc {"type":"system-stakeholder-role-spec-report","documents":"../specs/personal-task-list.system.spec.md#list-owner"} -->

The person who keeps the list relies on it to pick up unfinished errands after an interruption. They decide which entries can be discarded.

## Jobs

- [See personal commitments through](<jobs/see-personal-commitments-through.job.spec-report.md>) — See personal commitments through from accepting responsibility for an errand to establishing that the promised result has been delivered or someone else has agreed to take responsibility for it.

## Scope

One person records commitments on one computer and reviews them in Open, Waiting and Completed. They choose what reply to watch for and the day to check back. They receive the reply and judge the result outside the list.

### Out of scope

Sending reminders, reading messages, managing delivery deadlines, sharing lists, synchronizing devices and recovering from loss of the database file.

## Environment

The browser and service run on the same computer. The owner can close either between visits and return on a different local calendar day. The [follow-up requirement](<requirements/show-which-replies-need-follow-up.requirement.spec-report.md#conditions>) explains what they should see when they open Waiting again.

## Assumptions

- [Same local storage](<assumptions.md#same-local-storage>) — The person returns to the same writable database file.

## Business rules

- [Completion preserves a commitment](<business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>) — A settled commitment must stay on the person's list until they delete it.
- [When a reply needs chasing](<business-rules/when-a-reply-needs-chasing.business-rule.spec-report.md>) — For a reply the person is still waiting for, follow-up is upcoming before the currently chosen day and due from that day onward.

## Features

- [Pick up where I left off](<features/pick-up-where-i-left-off.feature.spec-report.md>) — Find an errand after a break, see whether a reply needs chasing and take up the work again when an answer arrives.

## Capabilities

- [Follow-up review](<capabilities/follow-up-review.capability.spec-report.md>) — When the owner reviews waiting errands, they can see which replies they planned to chase by today and what to ask for.
- [Local task storage](<capabilities/local-task-storage.capability.spec-report.md>) — Keeps a commitment and the owner's last recorded decision together between visits. …

## Use cases

- [Decide what to do about a waiting errand](<use-cases/decide-what-to-do-about-a-waiting-errand.use-case.spec-report.md>) — Decide whether to take up a waiting errand, keep waiting or mark it settled, and see that decision reflected in my list.

## Quality Attributes

- [Durability](<quality/durability/durability.quality-attribute.spec-report.md>) — Preservation of the owner’s acknowledged decisions between local sessions, including the reply they expect for a waiting errand and the day they chose to follow up. …

## Requirements

- [Acknowledge after commit](<components/task-service/requirements/acknowledge-after-commit.requirement.spec-report.md>) — The task service must record each change in a single commit to the local database and acknowledge it to the browser only after that commit succeeds.
- [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>) — When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). …
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry’s identity, title and state. …
- [Show which replies need follow-up](<requirements/show-which-replies-need-follow-up.requirement.spec-report.md>) — Each time the owner opens or refreshes Waiting, the list must show each waiting entry’s expected reply, recorded follow-up date and whether follow-up is due under the follow-up policy.

## Concepts

- [Checklist](<concepts/checklist.concept.spec-report.md>) — Decide what I can do next, which replies I need to chase and which commitments I have already settled.

## Concept uses

- [My checklist](<concept-uses/my-checklist.concept-use.spec-report.md>) — Each entry follows one commitment through the owner's next action, an outstanding reply or settlement. …

## Architecture

```mermaid
flowchart LR
    subgraph owner["Personal task list"]
        n1["Task service"]
        n2(["Task list forms"])
    end
```

### Components

- [Task service](<components/task-service/task-service.component.spec-report.md>) — Checks whether a submitted action still applies to the selected entry and records accepted changes.

### Interfaces

- [Task list forms](<interfaces/task-list-forms.interface.spec-report.md>) — The forms the list owner uses to review and update commitments.

---

[View source](<../specs/personal-task-list.system.spec.md#personal-task-list>) · [Quality models](<quality-models.md>)
