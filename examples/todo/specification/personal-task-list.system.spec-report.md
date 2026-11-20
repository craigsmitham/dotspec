---
generated_by: dotspec
type: system-spec-report
source: "../specs/personal-task-list.system.spec.md#personal-task-list"
---
# Personal task list

A single-user task list that stores open and completed tasks on one computer.

[Requirements review](<requirements-review.report.md>)

| Section | Explore |
| --- | --- |
| [Purpose](<#purpose>) | Help a person pick up unfinished work after an interruption. |
| [Goals](<#goals>) |  |
| [Business rules](<#business-rules>) | [Business rules](<business-rules.md>) |
| [User needs](<#user-needs>) | [Jobs](<jobs.md>) |
| [Scope](<#scope>) | Out of scope · Environment · Assumptions · Project context |
| [System topology](<#system-topology>) | [Components](<components.md>) · [Interfaces](<interfaces.md>) |
| [Requirements](<#requirements>) | [Features](<features.md>) · [Capabilities](<capabilities.md>) · [Use cases](<use-cases.md>) · [Concept definitions](<concepts.md>) · [Concept uses](<concept-uses.md>) · [All requirements](<requirements.md>) |
| [Quality attributes](<#quality-attributes>) | [Attribute index](<quality-attributes.md>) · [Quality models](<quality-models.md>) |
| [Validation](<#validation>) | System validation plan |

## Purpose

Help a person pick up unfinished work after an interruption.

## Goals

### Resume unfinished work

<!-- dotspec-doc {"type":"goal-spec-report","documents":"../specs/personal-task-list.system.spec.md#resume-unfinished-work"} -->

Reduce the effort of working out what still needs attention.

**Supported by** · [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>)

## Business rules

[Business rules](<business-rules.md>) — Policies that govern the system.

## User needs

### Stakeholders

The person who keeps the list also decides which entries can be discarded.

[Jobs](<jobs.md>) — The work people need to accomplish and the outcomes they seek.

- [Resume interrupted work](<jobs/resume-interrupted-work.job.spec-report.md>) — When I put work down, I want to remember what remains so I can pick it up again.

## Scope

One person's list on one computer. The browser submits forms to a local service.

### Out of scope

Sharing lists, synchronizing devices and recovering from loss of the database file.

### Environment

The person may close the browser or stop the service between work sessions.

### Assumptions

#### Same local storage

<!-- dotspec-doc {"type":"assumption-spec-report","documents":"../specs/personal-task-list.system.spec.md#same-local-storage"} -->

The person returns to the same writable database file.

**Relied on by** · [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

### Project context

The [walkthrough](<../README.md>) identifies the setup, checks and validation exercise.

## System topology

```mermaid
flowchart LR
    n0["Personal task list"]
    n1["Task service"]
    n0 -->|contains| n1
    n2["Task list forms"]
    n0 -->|provides interface| n2
```

- **Component** · [Task service](<components/task-service/task-service.component.spec-report.md>) — The local service that handles browser requests and saves changes to the task list.
  - Provides [Local task storage](<capabilities/local-task-storage.capability.spec-report.md>).
- **Provided interface** · [Task list forms](<interfaces/task-list-forms.interface.spec-report.md>) — The browser presents open and completed tasks and accepts changes through HTML forms.

[Components](<components.md>) — Parts of this system and their responsibilities.

[Interfaces](<interfaces.md>) — Required and provided interfaces and their contracts.

## Requirements

### Behavior

[Features](<features.md>) — Product behavior and the specifications each feature groups.

- [Pick up where I left off](<features/pick-up-where-i-left-off.feature.spec-report.md>) — Reopens the task list with each task in its last saved state.

[Capabilities](<capabilities.md>) — What the system can provide and who provides it.

- [Local task storage](<capabilities/local-task-storage.capability.spec-report.md>) — Retains accepted task changes between sessions on the same computer.

[Use cases](<use-cases.md>) — Actor goals, interactions, and outcomes.

- [Return to my list](<use-cases/return-to-my-list.use-case.spec-report.md>) — Pick up an outstanding task and bring a finished one back for follow-up.

### Concepts

[Concept definitions](<concepts.md>) — Definitions of purposeful behavior, state, and actions.

- [Checklist](<concepts/checklist.concept.spec-report.md>) — Keep track of which commitments still need attention.

[Concept uses](<concept-uses.md>) — How this system uses each concept.

- [My checklist](<concept-uses/my-checklist.concept-use.spec-report.md>) — The checklist used to track one person’s open and completed tasks.

### Obligations

[Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) refines [Keep recorded work available](<use-cases/return-to-my-list.use-case.spec-report.md#keep-recorded-work-available>).

[All requirements](<requirements.md>) — Standalone requirements and embedded obligations, including routes by component.

## Quality attributes

The task list is stored locally. Its quality requirements cover what survives between sessions on the same database.

**[Durability](<quality/durability/durability.quality-attribute.spec-report.md>)** — Preservation of acknowledged work between local sessions.

[Attribute index](<quality-attributes.md>)

[Quality models](<quality-models.md>) — Reference classifications and their definitions.

## Validation

### System validation plan

#### Subject

A person returning to unfinished work using this build.

#### Method

Invite someone to record two real errands and mark one finished. After a break,
restart the app with the same database. Ask what still needs attention. Then tell
them that the finished errand needs a follow-up and ask them to resume it.
Let them choose the controls without coaching.

#### Success criteria

The participant identifies the outstanding errand and resumes the other without
entering it again. Ask what they expected to find on returning.

#### Evidence to retain

With consent, record the build, task outcome, any help given and the participant's
explanation. Review mismatches before deciding whether the workflow meets the need.

[Requirements review](<requirements-review.report.md>)
