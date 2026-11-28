---
type: use-case-spec
---
# See what is left to do

## Goal

Work out which of my saved errands I still need to do.

## Primary actor

[List owner](../personal-task-list.system.spec.md#list-owner).

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](../personal-task-list.system.spec.md).

## Supports

[Take care of an errand I have to leave for later](../jobs/errands.job.spec.md).

## Main success scenario

1. Before going out again I open my list. Yesterday I checked off the library books.
2. I see “Mail the parcel” in Open. I find the books in Completed.
   - **Illustrates:** [Keep recorded work available](../requirements/availability.requirement.spec.md).

## Success guarantees

I can tell what remains, including when I have finished everything.

## Extensions

### I find another book in my bag

#### Condition

At step 2 I discover the library errand is unfinished.

#### Handling

1. I reopen the books from Completed and find the same task back in Open alongside the parcel.
   - **Illustrates:** [Reopen](../concepts/checklist.concept.spec.md#reopen).

### The list cannot be read

#### Condition

At step 1 the service cannot read my tasks.

#### Handling

1. The service reports that it could not access the task list. An empty-list message would tell me I could leave everything at home.

## Sequence

I can consult Open or Completed first, depending on whether I want to find unfinished work or check something I crossed off.
