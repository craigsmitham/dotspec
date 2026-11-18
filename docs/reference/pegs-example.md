# Project, Environment, Goals, System

Specification project: **todo**. Project membership does not imply a System relationship.

Generated PEGS reading perspective. This view complements business, user and software requirements subsets; it does not replace them or impose a dependency order. It covers only acquired specifications. Referenced project plans and native contracts remain authoritative.

Model `0.19.0-draft`; mapping `dotspec-markdown/14`.

## Project

### [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Project context**

The [walkthrough](../../examples/todo/README.md) identifies the setup, checks and validation exercise.

## Environment and boundaries

Operating conditions and constraints describe what the solution must accommodate. Assumptions are assessed separately for evidence and justified reliance. Explicit exclusions delimit scope. External partners are separate from structural parents.

### BusinessRule: [Completion preserves a commitment](<../../examples/todo/specs/business-rules/completion.business-rule.spec.md#completion-preserves-a-commitment>)

**Statement**

A finished commitment stays on the person's list until they explicitly delete it.

**Applicability**

**Targets:** [Personal task list](../../examples/todo/specs/system.spec.md).

### System: [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Environment**

The person may close the browser or stop the service between work sessions.

**Out of scope**

Sharing lists, synchronizing devices and recovering from loss of the database file.

**Provided interfaces**

#### Task list forms

The browser presents open and completed tasks and accepts changes through HTML forms.

##### Contracts

[Form template](../../examples/todo/app/list.html) supplies the actions and fields.

### Component: [Task service](<../../examples/todo/specs/system.spec.md#task-service>)

## Assumptions: evidence and justified reliance

Assess both empirical support in the stated conditions and whether relying on the premise is appropriate for intended use. An observation supports a bounded claim, not a guarantee that the premise always holds.

### [Same local storage](<../../examples/todo/specs/system.spec.md#same-local-storage>)

**Statement**

The person returns to the same writable database file.

Review required: no Basis, Applicability, Validation plan recorded.

## Goals

### Job: [Resume interrupted work](<../../examples/todo/specs/jobs/resume.job.spec.md#resume-interrupted-work>)

**Statement**

When I put work down, I want to remember what remains so I can pick it up again.

Supports: [Resume unfinished work](<../../examples/todo/specs/system.spec.md#resume-unfinished-work>).

### JobStep: [Recall outstanding work](<../../examples/todo/specs/jobs/resume.job.spec.md#recall-outstanding-work>)

**Description**

Identify the commitments that still need attention.

### JobStep: [Reconsider finished work](<../../examples/todo/specs/jobs/resume.job.spec.md#reconsider-finished-work>)

**Description**

Recognize when a previously finished commitment needs another action.

### DesiredOutcome: [Less reconstruction](<../../examples/todo/specs/jobs/resume.job.spec.md#less-reconstruction>)

**Statement**

Minimize the time spent reconstructing outstanding commitments after an interruption.

### Goal: [Resume unfinished work](<../../examples/todo/specs/system.spec.md#resume-unfinished-work>)

**Statement**

Reduce the effort of working out what still needs attention.

## System

Specifications below describe the solution's subjects, concepts, interactions and obligations. Declaration, structural parentage, applicability and behavioral relationships retain their distinct meanings.

### Concept: [Checklist](<../../examples/todo/specs/concepts/checklist.concept.spec.md#checklist>)

**Purpose**

Keep track of which commitments still need attention.

### Requirement: [Preserve tasks across restart](<../../examples/todo/specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>)

**Statement**

After the service acknowledges a change, stopping and restarting it on the same
database must preserve each surviving entry's identity, title and completion state.
Deleted entries must stay absent.

Supports: [Less reconstruction](<../../examples/todo/specs/jobs/resume.job.spec.md#less-reconstruction>).

Applies to: [Return to my list](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#return-to-my-list>).

### System: [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Purpose**

Help a person pick up unfinished work after an interruption.

**Scope**

One person's list on one computer. The browser submits forms to a local service.

### Component: [Task service](<../../examples/todo/specs/system.spec.md#task-service>)

**Responsibilities**

Accept browser changes and commit them to the local database before confirming them.

**Contracts**

- The browser uses the [HTML forms](../../examples/todo/app/list.html) to submit changes.

### Feature: [Pick up where I left off](<../../examples/todo/specs/system.spec.md#pick-up-where-i-left-off>)

Groups: [Return to my list](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#return-to-my-list>), [Preserve tasks across restart](<../../examples/todo/specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>).

### ConceptUse: [My checklist](<../../examples/todo/specs/system.spec.md#my-checklist>)

Definition: [Checklist](<../../examples/todo/specs/concepts/checklist.concept.spec.md#checklist>).

Supports: [Resume interrupted work](<../../examples/todo/specs/jobs/resume.job.spec.md#resume-interrupted-work>).

### Interface: [Task list forms](<../../examples/todo/specs/system.spec.md#task-list-forms>)

**Description**

The browser presents open and completed tasks and accepts changes through HTML forms.

**Contracts**

[Form template](../../examples/todo/app/list.html) supplies the actions and fields.

### UseCase: [Return to my list](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#return-to-my-list>)

**Goal**

Pick up an outstanding task and bring a finished one back for follow-up.

Supports: [Resume interrupted work](<../../examples/todo/specs/jobs/resume.job.spec.md#resume-interrupted-work>).

Involves: [My checklist](<../../examples/todo/specs/system.spec.md#my-checklist>).

### ObligationClause: [Keep recorded work available](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available>)

**Statement**

The list must make previously recorded entries available when the person returns.

Applies to: [My checklist](<../../examples/todo/specs/system.spec.md#my-checklist>).

## Validation planning

Identify the subject and the needs, intended use or other basis for the assessment. Requirements validation assesses proposed obligations; system validation assesses proposed or realized solutions. A method or supports link alone establishes neither purpose nor a result.

### Concept: [Checklist](<../../examples/todo/specs/concepts/checklist.concept.spec.md#checklist>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### Requirement: [Preserve tasks across restart](<../../examples/todo/specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### System: [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Validation plan**

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

Review separately identified plan details: Basis, Participants, Conditions. The narrative may already address them; no adequacy or completion is inferred.

### UseCase: [Return to my list](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#return-to-my-list>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### ObligationClause: [Keep recorded work available](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

## Open questions

Consolidated authored questions. External tracking references own resolution workflow. Omitted specification fields are not automatically questions; missing metadata below is advisory.

No explicit open questions recorded in this scope.
