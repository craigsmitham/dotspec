---
generated_by: dotspec
---
[Personal task list](<personal-task-list.system.spec-report.md>)

# Requirements

- [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>) — When the person returns to the same local database, each saved entry must appear once in the view for its recorded state (Open, Waiting or Completed). …
- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>) — After the service acknowledges a change, stopping and restarting it on the same database must preserve each surviving entry’s identity, title and state. …
  - Refines [Keep recorded work available](<requirements/keep-recorded-work-available.requirement.spec-report.md>)
- [Show which replies need follow-up](<requirements/show-which-replies-need-follow-up.requirement.spec-report.md>) — Each time the owner opens or refreshes Waiting, the list must show each waiting entry’s expected reply, recorded follow-up date and whether follow-up is due under the follow-up policy.

## By component

- [Task service requirements](<components/task-service/requirements.md>)
