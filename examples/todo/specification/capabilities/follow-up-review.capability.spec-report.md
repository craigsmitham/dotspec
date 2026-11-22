---
generated_by: dotspec
type: capability-spec-report
source: "../../specs/capabilities/follow-up-review.capability.spec.md#follow-up-review"
---
# Follow-up review

## Description

Identifies the outstanding replies the owner planned to chase by today. A person returning after the weekend can see both which errand needs attention and what to ask for.

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Depends on

[Local task storage](<local-task-storage.capability.spec-report.md>) supplies the recorded reply and day. The owner can return after the chosen day has passed.

## Supports

[Notice a stalled errand](<../jobs/see-personal-commitments-through.job.spec-report.md#notice-a-stalled-errand>).

## Provided by

[Task service](<../components/task-service/task-service.component.spec-report.md>).
