---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/resume-list.use-case.spec.md#return-to-my-list"
---
# Return to my list

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Goal

After a break, recover what I was waiting for and decide what to do about the errand.

## Primary actor

The person who promised to get the errand done.

## Preconditions

- The person has access to the local list where they recorded the errand.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#main-success-scenario"} -->

1. On Monday the person calls about their open errand, “Fix the guest sink”. They record that they are waiting for the plumber's visit time and choose Wednesday to follow up.
   - **Illustrates:** [Wait for reply](<../concepts/checklist.concept.spec-report.md#wait-for-reply>).
2. On Wednesday they return to Waiting. They recognize the repair by its title, see what the plumber owes them and see “Follow up now”.
   - **Illustrates:** [Keep recorded work available](<../requirements/keep-recorded-work-available.requirement.spec-report.md>) and [Show which replies need follow-up](<../requirements/show-which-replies-need-follow-up.requirement.spec-report.md>).
3. They check their messages. The plumber wants a photo before offering a visit. They choose Resume and find the same repair in Open, ready for their next action.
   - **Illustrates:** [Resume](<../concepts/checklist.concept.spec-report.md#resume>).

## Extensions

### The plumber has not answered

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-plumber-has-not-answered"} -->

#### Condition

At step 3 the person finds no reply, even though the chosen day has arrived.

#### Handling

1. The person uses the saved details to chase the plumber. If they still want to follow up today, they leave the entry as it is.
2. Friday passes. They agree a new arrangement with the guest and decide to chase the plumber again on Monday. In Waiting they update the follow-up day on the repair. On Saturday it appears as Upcoming.
   - **Illustrates:** [Update follow-up](<../concepts/checklist.concept.spec-report.md#update-follow-up>).

### The reply settles the repair

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-reply-settles-the-repair"} -->

#### Condition

At step 3 the plumber confirms the sink was repaired during an earlier visit, and the person is satisfied it works.

#### Handling

1. The person chooses Complete in Waiting and finds the repair in Completed.
   - **Illustrates:** [Complete](<../concepts/checklist.concept.spec-report.md#complete>).
2. If the leak recurs the next morning, the person can bring the repair back to Open.
   - **Illustrates:** [Reopen](<../concepts/checklist.concept.spec-report.md#reopen>).

### The person intentionally deleted the errand

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-person-intentionally-deleted-the-errand"} -->

#### Condition

The person still has a Waiting page open, but they have since deleted the repair entry in another window. [Delete](<../concepts/checklist.concept.spec-report.md#delete>).

#### Handling

1. When the reply arrives the person tries Resume on the old page. The list reports that the entry cannot be found.
2. On opening the list again, the person finds no record to resume. If they still owe the repair, they can add a new entry.

## Success guarantees

The person can identify the expected reply and find the commitment in the view that reflects their decision. They can take up the work, continue waiting with a day to check back or record the commitment as settled.

## Minimal guarantees

While the errand remains in Waiting, the person can identify the outstanding reply and the chosen follow-up day. If the entry has been deleted, trying an old form leaves the record absent. [Checklist effects](<../concepts/checklist.concept.spec-report.md#effects-table>).

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](<../personal-task-list.system.spec-report.md>).

Calls and replies happen outside the list. The person interprets what they hear and records what it means for their commitment.

## Supports

[See personal commitments through](<../jobs/see-personal-commitments-through.job.spec-report.md>).

## Involves

[My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

## Sequence

The person can act on a reply as soon as it arrives. Their chosen day governs when to chase an unanswered message. [When a reply needs chasing](<../business-rules/when-a-reply-needs-chasing.business-rule.spec-report.md>).

## Related specifications

**Applies here**

- [Preserve tasks across restart](<../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)
