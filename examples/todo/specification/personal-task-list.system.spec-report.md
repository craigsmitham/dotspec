---
generated_by: dotspec
type: system-spec-report
source: "../specs/personal-task-list.system.spec.md#personal-task-list"
---
# Personal task list

A personal list of errands the owner can act on, replies they are waiting for and commitments they have settled.

<details>
<summary>On this page</summary>

- [Framing](<#framing>)
- [Purpose](<#purpose>)
- [Goals](<#goals>)
- [Stakeholders](<#stakeholders>)
- [Scope](<#scope>)
- [Environment](<#environment>)
- [Assumptions](<#assumptions>)
- [Architecture](<#architecture>)
- [Capabilities](<#capabilities>)
- [Features](<#features>)
- [Concepts](<#concepts>)
- [Quality Attributes](<#quality-attributes>)

</details>

<details>
<summary>Browse specifications</summary>

- [Business rules](<business-rules.md>)
- [Jobs](<jobs.md>)
- [Requirements](<requirements.md>)
- [Use cases](<use-cases.md>)

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

[View section source](<../specs/personal-task-list.system.spec.md#list-owner>)

## Scope

One person records commitments on one computer and reviews them in Open, Waiting and Completed. They choose what reply to watch for and the day to check back. They receive the reply and judge the result outside the list.

### Out of scope

Sending reminders, reading messages, managing delivery deadlines, sharing lists, synchronizing devices and recovering from loss of the database file.

## Environment

The browser and service run on the same computer. The owner can close either between visits and return on a different local calendar day. The [follow-up requirement](<requirements/show-which-replies-need-follow-up.requirement.spec-report.md#conditions>) explains what they should see when they open Waiting again.

## Assumptions

- [Same local storage](<assumptions.md#same-local-storage>) — The person returns to the same writable database file.

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

## Capabilities

| Capability | What it enables | Provided by |
| --- | --- | --- |
| [Follow-up review](<capabilities/follow-up-review.capability.spec-report.md>) | When the owner reviews waiting errands, they can see which replies they planned to chase by today and what to ask for. | [Task service](<components/task-service/task-service.component.spec-report.md>) |
| [Local task storage](<capabilities/local-task-storage.capability.spec-report.md>) | Keeps a commitment and the owner's last recorded decision together between visits. … | [Task service](<components/task-service/task-service.component.spec-report.md>) |

## Features

- [Pick up where I left off](<features/pick-up-where-i-left-off.feature.spec-report.md>) — Find an errand after a break, see whether a reply needs chasing and take up the work again when an answer arrives.

## Concepts

<a id="concept.checklist"></a>

### [Checklist](<concepts/checklist.concept.spec-report.md>)

Decide what I can do next, which replies I need to chase and which commitments I have already settled.

**Use in this system**

- [My checklist](<concept-uses/my-checklist.concept-use.spec-report.md>) — Each entry follows one commitment through the owner's next action, an outstanding reply or settlement. Two errands with the same title can be at different stages.

## Quality Attributes

**[Durability](<quality/durability/durability.quality-attribute.spec-report.md>)** — Preservation of the owner’s acknowledged decisions between local sessions, including the reply they expect for a waiting errand and the day they chose to follow up. …

---

[View source](<../specs/personal-task-list.system.spec.md#personal-task-list>) · [Quality models](<quality-models.md>)
