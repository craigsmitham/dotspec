---
generated_by: dotspec
type: business-rule-spec-report
source: "../../specs/business-rules/follow-up.business-rule.spec.md#when-a-reply-needs-chasing"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Business rules](<../personal-task-list.system.spec-report.md#business-rules>)

# When a reply needs chasing

## Statement

For a reply the person is still waiting for, the [follow-up decision table](<#specification.follow-up-decision-table>) determines whether to chase it by comparing today with the chosen day.

## Applies to

**Targets:** [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>).

The person chooses a day to check back for each reply they are waiting for. Reaching that day leaves the commitment waiting. The person decides whether to chase the reply, take up the work again or mark it settled.

## Authority

The worked example sets the policy for when follow-up becomes due. The person keeping the list chooses the date for each outstanding reply.

## Rationale

Suppose my guest arrives on Friday. I choose Wednesday to chase the plumber so I still have time to find someone else. A message on Tuesday asking for a photo gives me work to do straight away. Silence on Saturday still leaves me owing my guest a repair or a new arrangement.

## Related specifications

[Show which replies need follow-up](<../requirements/show-which-replies-need-follow-up.requirement.spec-report.md>) sets the observable result when the person opens the list.

**Basis for**

- [Show which replies need follow-up](<../requirements/show-which-replies-need-follow-up.requirement.spec-report.md>)

## Specification sources

<a id="specification.follow-up-decision-table"></a>

### [Follow-up decision table](<../../specs/business-rules/follow-up.csv>)

Defines business rules.

---

[View source](<../../specs/business-rules/follow-up.business-rule.spec.md#when-a-reply-needs-chasing>)
