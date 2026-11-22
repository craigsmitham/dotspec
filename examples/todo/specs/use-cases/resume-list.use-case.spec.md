---
type: use-case-spec
---
# Return to my list

## Goal

After a break, recover what I was waiting for and decide what to do about the errand.

## Primary actor

The person who promised to get the errand done.

## Goal level

User task.

## System boundary

**Subject system:** [Personal task list](../personal-task-list.system.spec.md).

Calls and replies happen outside the list. The person interprets what they hear and records what it means for their commitment.

## Supports

[See personal commitments through](../jobs/personal-commitments.job.spec.md).

## Involves

[My checklist](../personal-task-list.system.spec.md#my-checklist).

## Preconditions

- The person has access to the local list where they recorded the errand.

## Main success scenario

1. On Monday the person calls about their open errand, “Fix the guest sink”. They record that they are waiting for the plumber's visit time and choose Wednesday to follow up.
   - **Illustrates:** [Wait for reply](../concepts/checklist.concept.spec.md#wait-for-reply).
2. On Wednesday they return to Waiting. They recognize the repair by its title, see what the plumber owes them and see “Follow up now”.
   - **Illustrates:** [Keep recorded work available](../requirements/availability.requirement.spec.md) and [Show which replies need follow-up](../requirements/follow-up.requirement.spec.md).
3. They check their messages. The plumber wants a photo before offering a visit. They choose Resume and find the same repair in Open, ready for their next action.
   - **Illustrates:** [Resume](../concepts/checklist.concept.spec.md#resume).

## Success guarantees

The person can identify the expected reply and find the commitment in the view that reflects their decision. They can take up the work, continue waiting with a day to check back or record the commitment as settled.

## Minimal guarantees

While the errand remains in Waiting, the person can identify the outstanding reply and the chosen follow-up day. If the entry has been deleted, trying an old form leaves the record absent. [Checklist effects](../concepts/checklist.concept.spec.md#effects-table).

## Sequence

The person can act on a reply as soon as it arrives. Their chosen day governs when to chase an unanswered message. [When a reply needs chasing](../business-rules/follow-up.business-rule.spec.md).

## Extensions

### The plumber has not answered

#### Condition

At step 3 the person finds no reply, even though the chosen day has arrived.

#### Handling

1. The person uses the saved details to chase the plumber. If they still want to follow up today, they leave the entry as it is.
2. Friday passes. They agree a new arrangement with the guest and decide to chase the plumber again on Monday. In Waiting they update the follow-up day on the repair. On Saturday it appears as Upcoming.
   - **Illustrates:** [Update follow-up](../concepts/checklist.concept.spec.md#update-follow-up).

### The reply settles the repair

#### Condition

At step 3 the plumber confirms the sink was repaired during an earlier visit, and the person is satisfied it works.

#### Handling

1. The person chooses Complete in Waiting and finds the repair in Completed.
   - **Illustrates:** [Complete](../concepts/checklist.concept.spec.md#complete).
2. If the leak recurs the next morning, the person can bring the repair back to Open.
   - **Illustrates:** [Reopen](../concepts/checklist.concept.spec.md#reopen).

### The person intentionally deleted the errand

#### Condition

The person still has a Waiting page open, but they have since deleted the repair entry in another window. [Delete](../concepts/checklist.concept.spec.md#delete).

#### Handling

1. When the reply arrives the person tries Resume on the old page. The list reports that the entry cannot be found.
2. On opening the list again, the person finds no record to resume. If they still owe the repair, they can add a new entry.
