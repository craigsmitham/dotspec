---
type: use-case-spec
---
# Save an errand for later

## Goal

Put an errand down so I can get back to what I was doing.

## Primary actor

[List owner](../personal-task-list.system.spec.md#list-owner).

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](../personal-task-list.system.spec.md).

## Supports

[Take care of an errand I have to leave for later](../jobs/errands.job.spec.md).

## Main success scenario

1. At breakfast I add “Mail the parcel” with just its title.
2. The list shows the new task in Open.
   - **Illustrates:** [Add](../concepts/checklist.concept.spec.md#add).

## Success guarantees

I can find the errand when I next review my list.

## Extensions

### I submit an empty title

#### Condition

At step 1 the title contains only spaces.

#### Handling

1. The list asks for a task name and leaves my tasks as they were. I can enter the title and try again.
   - **Illustrates:** [Add](../concepts/checklist.concept.spec.md#add).

## Sequence

The person can add errands in whichever order they remember them.
