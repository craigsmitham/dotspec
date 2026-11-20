---
generated_by: dotspec
---
# Requirements

[Personal task list](<personal-task-list.system.spec-report.md>)

- [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>) — When the person returns to the same local database, each saved entry must be available in its open or completed view. An intentionally deleted entry must remain absent from both views.
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry's identity, title and completion state. Deleted entries must stay absent.
  - Refines [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>)

## By component

- [Task service](<components/task-service/requirements.md>) — The local service that handles browser requests and saves changes to the task list.
