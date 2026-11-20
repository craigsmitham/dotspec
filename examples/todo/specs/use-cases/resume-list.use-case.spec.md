---
type: use-case-spec
---
# Return to my list

## Goal

Pick up an outstanding task and bring a finished one back for follow-up.

## Primary actor

The person who keeps the list.

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](../personal-task-list.system.spec.md).

## Supports

[Resume interrupted work](../jobs/resume.job.spec.md).

## Involves

[My checklist](../personal-task-list.system.spec.md#my-checklist).

## Preconditions

- The person has access to their local database.

## Minimal guarantees

- A rejected rename preserves the existing title.

## Main success scenario

1. The person records an errand and later marks it complete.
   - **Illustrates:** [Complete](../concepts/checklist.concept.spec.md#complete).
2. On returning, the person finds the errand among the completed tasks.
   - **Illustrates:** [Keep recorded work available](#keep-recorded-work-available).
3. The person discovers a follow-up and reopens the errand.
   - **Illustrates:** [Reopen](../concepts/checklist.concept.spec.md#reopen).

## Success guarantees

- The person can continue with the original entry.

## Sequence

The example follows one errand across two work sessions.

## Extensions

### The person intentionally deleted the errand

#### Condition

The person previously used [Delete](../concepts/checklist.concept.spec.md#delete).

#### Handling

1. The person looks through the two views. Neither contains the deleted entry.
2. If the commitment still matters, the person records a new entry.

## Obligations

### Keep recorded work available

**ID:** keep-recorded-work-available

The list must make previously recorded entries available when the person returns.

#### Derived from

[Completion preserves a commitment](../business-rules/completion.business-rule.spec.md).

#### Applies to

**Targets:** [My checklist](../personal-task-list.system.spec.md#my-checklist).

#### System verification plan

Record a task, complete it and return to both views. Reopen it and compare its
identity with the original. For an interruption between sessions, follow the
[restart verification plan](../requirements/restart.requirement.spec.md#system-verification-plan).
