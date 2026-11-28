---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/save-errand.use-case.spec.md#save-an-errand-for-later"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Use cases](<../personal-task-list.system.spec-report.md#use-cases>)

# Save an errand for later

## Goal

Put an errand down so I can get back to what I was doing.

## Primary actor

[List owner](<../personal-task-list.system.spec-report.md#list-owner>).

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/save-errand.use-case.spec.md#main-success-scenario"} -->

1. At breakfast I add “Mail the parcel” with just its title.
2. The list shows the new task in Open.
   - **Illustrates:** [Add](<../concepts/checklist.concept.spec-report.md#add>).

## Extensions

### I submit an empty title

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/save-errand.use-case.spec.md#i-submit-an-empty-title"} -->

#### Condition

At step 1 the title contains only spaces.

#### Handling

1. The list asks for a task name and leaves my tasks as they were. I can enter the title and try again.
   - **Illustrates:** [Add](<../concepts/checklist.concept.spec-report.md#add>).

## Success guarantees

I can find the errand when I next review my list.

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](<../personal-task-list.system.spec-report.md>).

## Supports

[Take care of an errand I have to leave for later](<../jobs/take-care-of-an-errand-i-have-to-leave-for-later.job.spec-report.md>).

## Sequence

The person can add errands in whichever order they remember them.

## Related specifications

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

---

[View source](<../../specs/use-cases/save-errand.use-case.spec.md#save-an-errand-for-later>)
