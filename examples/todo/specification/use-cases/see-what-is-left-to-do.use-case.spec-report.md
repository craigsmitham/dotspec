---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/see-remaining.use-case.spec.md#see-what-is-left-to-do"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Use cases](<../personal-task-list.system.spec-report.md#use-cases>)

# See what is left to do

## Goal

Work out which of my saved errands I still need to do.

## Primary actor

[List owner](<../personal-task-list.system.spec-report.md#list-owner>).

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/see-remaining.use-case.spec.md#main-success-scenario"} -->

1. Before going out again I open my list. Yesterday I checked off the library books.
2. I see “Mail the parcel” in Open. I find the books in Completed.
   - **Illustrates:** [Keep recorded work available](<../requirements/keep-recorded-work-available.requirement.spec-report.md>).

## Extensions

[I find another book in my bag](<#i-find-another-book-in-my-bag>) · [The list cannot be read](<#the-list-cannot-be-read>)

### I find another book in my bag

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/see-remaining.use-case.spec.md#i-find-another-book-in-my-bag"} -->

#### Condition

At step 2 I discover the library errand is unfinished.

#### Handling

1. I reopen the books from Completed and find the same task back in Open alongside the parcel.
   - **Illustrates:** [Reopen](<../concepts/checklist.concept.spec-report.md#reopen>).

### The list cannot be read

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/see-remaining.use-case.spec.md#the-list-cannot-be-read"} -->

#### Condition

At step 1 the service cannot read my tasks.

#### Handling

1. The service reports that it could not access the task list. An empty-list message would tell me I could leave everything at home.

## Success guarantees

I can tell what remains, including when I have finished everything.

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](<../personal-task-list.system.spec-report.md>).

## Supports

[Take care of an errand I have to leave for later](<../jobs/take-care-of-an-errand-i-have-to-leave-for-later.job.spec-report.md>).

## Sequence

I can consult Open or Completed first, depending on whether I want to find unfinished work or check something I crossed off.

## Related specifications

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

---

[View source](<../../specs/use-cases/see-remaining.use-case.spec.md#see-what-is-left-to-do>)
