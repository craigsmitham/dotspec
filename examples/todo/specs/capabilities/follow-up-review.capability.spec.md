---
type: capability-spec
---
# Follow-up review

## Description

When the owner reviews waiting errands, they can see which replies they planned to chase by today and what to ask for.

## Depends on

[Remember task progress between visits](local-task-storage.capability.spec.md) keeps the owner's note about the expected reply and their chosen follow-up day available between visits.

## Supports

[Notice a stalled errand](../jobs/personal-commitments.job.spec.md#notice-a-stalled-errand).

## Enablers

The owner records what answer they need and chooses when to check back. When they review the list, they consult their messages to decide whether to chase the reply.

## Provided by

[Task service](../../app/specs/task-service.component.spec.md).
