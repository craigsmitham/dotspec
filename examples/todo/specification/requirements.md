---
generated_by: dotspec
---
# Requirements

[Personal task list](<personal-task-list.system.spec-report.md>)

- [Keep recorded work available](<use-cases/return-to-my-list.use-case.spec-report.md#keep-recorded-work-available>) — The list must make previously recorded entries available when the person returns.
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry's identity, title and completion state. Deleted entries must stay absent.

## By component

- [Task service](<components/task-service/requirements.md>) — The local service that handles browser requests and saves changes to the task list.
