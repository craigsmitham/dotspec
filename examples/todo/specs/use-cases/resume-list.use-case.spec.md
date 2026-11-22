---
type: use-case-spec
---
# Decide what to do about a waiting errand

## Goal

Decide whether to take up a waiting errand, keep waiting or mark it settled, and see that decision reflected in my list.

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

- The person has access to their local list. On Monday they recorded “Fix the guest sink” as waiting for the plumber’s visit time and chose Wednesday to follow up. [Wait for reply](../concepts/checklist.concept.spec.md#wait-for-reply).

## Main success scenario

1. On Wednesday the person checks which waiting errands need attention.
2. The list shows “Fix the guest sink” with the saved note “Waiting for the plumber’s visit time”. The Wednesday follow-up reads “Follow up now”.
   - **Illustrates:** [Keep recorded work available](../requirements/availability.requirement.spec.md) and [Show which replies need follow-up](../requirements/follow-up.requirement.spec.md).
3. The person checks their messages. The plumber needs a photo before offering a visit, so the person takes back the next action on the repair.
4. The list shows the same repair in Open. The person knows it is theirs to act on.
   - **Illustrates:** [Resume](../concepts/checklist.concept.spec.md#resume).

## Success guarantees

The list reflects the person’s decision about the same errand. They can tell whether it is theirs to act on, still waiting for a reply with a chosen follow-up day, or settled.

## Minimal guarantees

While the errand remains in Waiting, the person can identify the outstanding reply and the chosen follow-up day. If the entry has been deleted, trying an old form leaves the record absent. [Checklist effects](../concepts/checklist.concept.spec.md#effects-table).

## Sequence

The person can act on a reply as soon as it arrives. Their chosen day governs when to chase an unanswered message. [When a reply needs chasing](../business-rules/follow-up.business-rule.spec.md).

## Extensions

### The plumber has not answered

#### Condition

At step 3 the person finds no reply, even though the chosen day has arrived.

#### Handling

1. The person uses the saved details to chase the plumber.
2. They choose when to check again. Keeping today’s date leaves the repair in Waiting with “Follow up now”. Choosing Friday keeps it in Waiting with the new date and “Upcoming”.
   - **Illustrates:** [Update follow-up](../concepts/checklist.concept.spec.md#update-follow-up).

### The reply settles the repair

#### Condition

At step 3 the plumber confirms the sink was repaired during an earlier visit, and the person is satisfied it works.

#### Handling

1. The person chooses Complete in Waiting and finds the repair in Completed.
   - **Illustrates:** [Complete](../concepts/checklist.concept.spec.md#complete).

### The person intentionally deleted the errand

#### Condition

At step 3 the person tries to resume the repair from their Waiting page after deleting the entry in another window. [Delete](../concepts/checklist.concept.spec.md#delete).

#### Handling

1. The list reports that it cannot find the entry.
2. The person opens the list again and sees that the deleted repair is still absent.
