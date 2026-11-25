---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/follow-up.requirement.spec.md#show-which-replies-need-follow-up"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Requirements](<../requirements.md>)

# Show which replies need follow-up

## Statement

Each time the owner opens or refreshes Waiting, the list must show each waiting entry’s expected reply, recorded follow-up date and whether follow-up is due under the [follow-up policy](<../business-rules/when-a-reply-needs-chasing.business-rule.spec-report.md>).

## Applies to

**Targets:** [Task list forms](<../interfaces/task-list-forms.interface.spec-report.md>).

### Conditions

Use the local calendar day on the computer serving the list when it prepares the page. A page left open overnight is assessed when it is next loaded.

## Criterion

For an entry marked for Friday, Thursday’s view reads **Upcoming**. Friday’s and Saturday’s views read **Follow up now** if the owner is still waiting. In all three cases the owner can see what reply they are chasing. Moving the entry to Open at midnight would violate the promise even if the follow-up date were still visible there.

## Supports

[Notice a stalled errand](<../jobs/see-personal-commitments-through.job.spec-report.md#notice-a-stalled-errand>).

## Derived from

[When a reply needs chasing](<../business-rules/when-a-reply-needs-chasing.business-rule.spec-report.md>).

## Rationale

Someone who has already called three shops needs to know which call to make again today and what answer to ask for.

## Evidence

The [form acceptance test](<../../tests/acceptance.py>) checks the reply, date and cue for entries dated before, on and after today. On Unix it also opens the same saved entry with the service in two time zones. The reply and date stay the same while its cue changes from **Upcoming** to **Follow up now**.

## Related specifications

**Feature**

- [Pick up where I left off](<../features/pick-up-where-i-left-off.feature.spec-report.md>)

---

[View source](<../../specs/requirements/follow-up.requirement.spec.md#show-which-replies-need-follow-up>)
