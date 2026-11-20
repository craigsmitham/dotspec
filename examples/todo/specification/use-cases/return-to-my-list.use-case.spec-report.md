---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/resume-list.use-case.spec.md#return-to-my-list"
---
# Return to my list

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Goal

Pick up an outstanding task and bring a finished one back for follow-up.

## Primary actor

The person who keeps the list.

## Preconditions

- The person has access to their local database.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#main-success-scenario"} -->

1. The person records an errand and later marks it complete.
   - **Illustrates:** [Complete](<../concepts/checklist.concept.spec-report.md#complete>).
2. On returning, the person finds the errand among the completed tasks.
   - **Illustrates:** [Keep recorded work available](<../requirements/keep-recorded-work-available.requirement.spec-report.md>).
3. The person discovers a follow-up and reopens the errand.
   - **Illustrates:** [Reopen](<../concepts/checklist.concept.spec-report.md#reopen>).

## Extensions

### The person intentionally deleted the errand

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-person-intentionally-deleted-the-errand"} -->

#### Condition

The person previously used [Delete](<../concepts/checklist.concept.spec-report.md#delete>).

#### Handling

1. The person looks through the two views. Neither contains the deleted entry.
2. If the commitment still matters, the person records a new entry.

## Success guarantees

- The person can continue with the original entry.

## Minimal guarantees

- A rejected rename preserves the existing title.

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](<../personal-task-list.system.spec-report.md>).

## Supports

[Resume interrupted work](<../jobs/resume-interrupted-work.job.spec-report.md>).

## Involves

[My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

## Sequence

The example follows one errand across two work sessions.

## Related specifications

**Applies here**

- [Preserve tasks across restart](<../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)
