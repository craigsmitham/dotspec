---
generated_by: dotspec
---
# Assumptions

[Personal task list](<personal-task-list.system.spec-report.md>)

## Same local storage

<!-- dotspec-doc {"type":"assumption-spec-report","documents":"../specs/personal-task-list.system.spec.md#same-local-storage"} -->

The person returns to the same writable database file.

**Relied on by**

- [Preserve tasks across restart](<requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

### Basis

The owner names the database file when starting the service, or accepts `todo.sqlite` in its working directory. Starting it again the same way reopens the same list. Another path opens a different list, which the restart requirement does not cover.

### Applicability

Each return to the list on the owner's computer, including after the service restarts.
