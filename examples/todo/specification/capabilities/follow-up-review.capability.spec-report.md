---
generated_by: dotspec
type: capability-spec-report
source: "../../specs/capabilities/follow-up-review.capability.spec.md#follow-up-review"
---
# Follow-up review

## Description

When the owner reviews waiting errands, they can see which replies they planned to chase by today and what to ask for.

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Depends on

[Local task storage](<local-task-storage.capability.spec-report.md>) keeps the owner's note about the expected reply and their chosen follow-up day available between visits.

## Supports

[Notice a stalled errand](<../jobs/see-personal-commitments-through.job.spec-report.md#notice-a-stalled-errand>).

## Provided by

[Task service](<../components/task-service/task-service.component.spec-report.md>).
