---
type: requirement-spec
---
# Show which replies need follow-up

## Statement

Each time the owner opens or refreshes Waiting, the list must show each waiting entry’s expected reply, recorded follow-up date and whether follow-up is due under the [follow-up policy](../business-rules/follow-up.business-rule.spec.md).

## Supports

[Notice a stalled errand](../jobs/personal-commitments.job.spec.md#notice-a-stalled-errand).

## Applies to

**Targets:** [Task list forms](../personal-task-list.system.spec.md#task-list-forms).

### Conditions

Use the local calendar day on the computer serving the list when it prepares the page. A page left open overnight is assessed when it is next loaded.

## Criterion

For an entry marked for Friday, Thursday’s view reads **Upcoming**. Friday’s and Saturday’s views read **Follow up now** if the owner is still waiting. In all three cases the owner can see what reply they are chasing. Moving the entry to Open at midnight would violate the promise even if the follow-up date were still visible there.

## Derived from

[When a reply needs chasing](../business-rules/follow-up.business-rule.spec.md).

## Rationale

Someone who has already called three shops needs to know which call to make again today and what answer to ask for.

## Evidence

The [form acceptance test](../../tests/acceptance.py) checks the reply, date and cue for entries dated before, on and after today. The [time-zone cases](../quality/durability/durability.quality-attribute.spec.md#specification.saved-follow-up-dates-across-time-zones) cover returning to a saved list on a different local day.
