---
generated_by: dotspec
type: job-spec-report
source: "../../specs/jobs/resume.job.spec.md#resume-interrupted-work"
---
# Resume interrupted work

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Job

When an errand needs attention again after a break, I want to recover what I already did and what I still owe, so I can take the next action without repeating a finished step or overlooking a follow-up.

## Performer

The person responsible for seeing their own errands through to completion.

## Situation

Suppose I called the plumber yesterday and moved on to other errands. Today the plumber asks for a photo. I need to recover where we left off before I decide what to do next.

## Supports

[Resume unfinished work](<../goals.md#resume-unfinished-work>).

## Steps

### Recall outstanding work

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/resume.job.spec.md#recall-outstanding-work"} -->

**ID:** recall
**Stage:** locate

Distinguish the errands that still need action from the ones already settled.

**Applies here**

- [Less reconstruction](<#less-reconstruction>)

### Reconsider finished work

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/resume.job.spec.md#reconsider-finished-work"} -->

**ID:** reconsider
**Stage:** modify

Determine what a new request adds to an errand that previously seemed finished.

## Outcomes

### Less reconstruction

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/resume.job.spec.md#less-reconstruction"} -->

**ID:** less-reconstruction

Minimize the time between returning to an interrupted errand and identifying both the last completed action and the next outstanding commitment. In the repair example, I can say that the call is finished and sending the photo is next.

**Supported by**

- [Preserve tasks across restart](<../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

#### Applies to

[Recall outstanding work](<#recall-outstanding-work>).

## Evidence

This example assumes that reconstructing earlier work costs the performer time. The [return scenario](<../use-cases/return-to-my-list.use-case.spec-report.md>) illustrates a finished errand that needs follow-up. In the [participant exercise](<../../README.md#try-the-workflow-with-a-participant>), ask what the person consulted to recover the next action and which details they still had to check.

## Related specifications

**Supported by**

- [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
- [Return to my list](<../use-cases/return-to-my-list.use-case.spec-report.md>)
