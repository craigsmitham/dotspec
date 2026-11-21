---
type: job-spec
---
# Resume interrupted work

## Job

When an errand needs attention again after a break, I want to recover what I already did and what I still owe, so I can take the next action without repeating a finished step or overlooking a follow-up.

## Performer

The person responsible for seeing their own errands through to completion.

## Situation

Suppose I called the plumber yesterday and moved on to other errands. Today the plumber asks for a photo. I need to recover where we left off before I decide what to do next.

## Supports

[Resume unfinished work](../personal-task-list.system.spec.md#resume-unfinished-work).

## Steps

### Recall outstanding work

**ID:** recall
**Stage:** locate

Distinguish the errands that still need action from the ones already settled.

### Reconsider finished work

**ID:** reconsider
**Stage:** modify

Determine what a new request adds to an errand that previously seemed finished.

## Outcomes

### Less reconstruction

**ID:** less-reconstruction

Minimize the time between returning to an interrupted errand and identifying both the last completed action and the next outstanding commitment. In the repair example, I can say that the call is finished and sending the photo is next.

#### Applies to

[Recall outstanding work](#recall-outstanding-work).

## Evidence

This example assumes that reconstructing earlier work costs the performer time. The [return scenario](../use-cases/resume-list.use-case.spec.md) illustrates a finished errand that needs follow-up. In the [participant exercise](../../README.md#try-the-workflow-with-a-participant), ask what the person consulted to recover the next action and which details they still had to check.
