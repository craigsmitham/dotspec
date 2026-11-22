---
generated_by: dotspec
---
# Requirements

[Personal task list](<personal-task-list.system.spec-report.md>)

- [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>) — When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). An intentionally deleted entry must remain absent from every view.
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry’s identity, title and state. For a waiting entry, the expected reply and follow-up date must also survive. Deleted entries must stay absent.
  - Refines [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>)
- [Show which replies need follow-up](<requirements/show-which-replies-need-follow-up.requirement.spec-report.md>) — Each time the owner opens or refreshes Waiting, the list must show each waiting entry’s expected reply, recorded follow-up date and whether follow-up is due under the follow-up policy.

## By component

- [Task service](<components/task-service/requirements.md>) — The local service that decides whether a submitted action still applies to the selected entry and records accepted changes.
