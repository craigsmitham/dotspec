---
type: business-rule-spec
---
# When a reply needs chasing

## Statement

A reply the person is still waiting for is due for follow-up if its current follow-up day is today or has passed.

## Applies to

**Targets:** [My checklist](../personal-task-list.system.spec.md#my-checklist).

The person chooses a day to check back for each reply they are waiting for. Reaching that day leaves the commitment waiting. The person decides whether to chase the reply, take up the work again or mark it settled.

## Authority

A product decision for this worked example. The person keeping the list decides what counts as settling each commitment.

## Rationale

Suppose my guest arrives on Friday. I choose Wednesday to chase the plumber so I still have time to find someone else. A message on Tuesday asking for a photo gives me work to do straight away. Silence on Saturday still leaves me owing my guest a repair or a new arrangement.

## Related specifications

[Show which replies need follow-up](../requirements/follow-up.requirement.spec.md) sets the observable result when the person opens the list.
