---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/resume-list.use-case.spec.md#decide-what-to-do-about-a-waiting-errand"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Use cases](<../use-cases.md>)

# Decide what to do about a waiting errand

## Goal

Decide whether to take up a waiting errand, keep waiting or mark it settled, and see that decision reflected in my list.

## Primary actor

[The person who promised to get the errand done](<../personal-task-list.system.spec-report.md#list-owner>).

## Preconditions

- The person can open their local list and has a waiting errand with a recorded reply to watch for and a chosen follow-up day. [Wait for reply](<../concepts/checklist.concept.spec-report.md#wait-for-reply>).

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#main-success-scenario"} -->

[View section source](<../../specs/use-cases/resume-list.use-case.spec.md#main-success-scenario>)

1. On Wednesday the person opens Waiting to see whether it is time to chase the plumber about the sink.
2. The list shows “Fix the guest sink” with the saved note “Waiting for the plumber’s visit time”. The Wednesday follow-up reads “Follow up now”.
   - **Illustrates:** [Keep recorded work available](<../requirements/keep-recorded-work-available.requirement.spec-report.md>) and [Show which replies need follow-up](<../requirements/show-which-replies-need-follow-up.requirement.spec-report.md>).
3. The person checks their messages. The plumber needs a photo before offering a visit, so the person takes back the next action on the repair.
4. The repair disappears from Waiting. The person opens the Open list and finds the same repair there, ready for them to act on.
   - **Illustrates:** [Resume](<../concepts/checklist.concept.spec-report.md#resume>).

## Extensions

[The plumber has not answered](<#the-plumber-has-not-answered>) · [The reply settles the repair](<#the-reply-settles-the-repair>) · [The person intentionally deleted the errand](<#the-person-intentionally-deleted-the-errand>) · [The list cannot record the decision](<#the-list-cannot-record-the-decision>)

### The plumber has not answered

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-plumber-has-not-answered"} -->

[View section source](<../../specs/use-cases/resume-list.use-case.spec.md#the-plumber-has-not-answered>)

#### Condition

At step 3 the person finds no reply, even though the chosen day has arrived.

#### Handling

1. The person uses the saved details to chase the plumber.
2. They choose when to check again. Keeping today’s date leaves the repair in Waiting with “Follow up now”. Choosing Friday keeps it in Waiting with the new date and “Upcoming”.
   - **Illustrates:** [Update follow-up](<../concepts/checklist.concept.spec-report.md#update-follow-up>).

### The reply settles the repair

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-reply-settles-the-repair"} -->

[View section source](<../../specs/use-cases/resume-list.use-case.spec.md#the-reply-settles-the-repair>)

#### Condition

At step 3 the plumber confirms the sink was repaired during an earlier visit, and the person is satisfied it works.

#### Handling

1. The person chooses Complete in Waiting. The repair disappears from that view. They open Completed and find the settled repair there.
   - **Illustrates:** [Complete](<../concepts/checklist.concept.spec-report.md#complete>).

### The person intentionally deleted the errand

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-person-intentionally-deleted-the-errand"} -->

[View section source](<../../specs/use-cases/resume-list.use-case.spec.md#the-person-intentionally-deleted-the-errand>)

#### Condition

At step 3 the person tries to resume the repair from their Waiting page after deleting the entry in another window. [Delete](<../concepts/checklist.concept.spec-report.md#delete>).

#### Handling

1. The list reports that it cannot find the entry.
2. The person opens the list again and sees that the deleted repair is still absent.

### The list cannot record the decision

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/resume-list.use-case.spec.md#the-list-cannot-record-the-decision"} -->

[View section source](<../../specs/use-cases/resume-list.use-case.spec.md#the-list-cannot-record-the-decision>)

#### Condition

At step 3 the person chooses Resume or Complete, but the list cannot save the change.

#### Handling

1. The list shows “Could not access the task list.” instead of returning to Waiting.
   - **Illustrates:** [Acknowledge after commit](<../components/task-service/requirements/acknowledge-after-commit.requirement.spec-report.md>).
2. The person opens Waiting again and finds the repair with the same expected reply and follow-up day. The decision was not recorded, and the errand is still waiting.
   - **Illustrates:** [Keep recorded work available](<../requirements/keep-recorded-work-available.requirement.spec-report.md>).

## Success guarantees

The list reflects the person’s decision about the same errand. They can tell whether it is theirs to act on, still waiting for a reply with a chosen follow-up day, or settled.

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

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

---

[View source](<../../specs/use-cases/resume-list.use-case.spec.md#decide-what-to-do-about-a-waiting-errand>)
