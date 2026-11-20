---
generated_by: dotspec
type: component-spec-report
source: "../../../app/specs/task-service.component.spec.md#task-service"
---
# Task service

## Description

The local service that handles browser requests and saves changes to the task list.

[Personal task list](<../../personal-task-list.system.spec-report.md>)

- [Capabilities](<capabilities.md>)
- [Requirements](<requirements.md>)

## Part of

[Personal task list](<../../personal-task-list.system.spec-report.md>).

## Responsibilities

Accept browser changes and persist them locally. See [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>) for the confirmation guarantee.

## Contracts

- The browser uses the [HTML forms](<../../../app/list.html>) to submit changes.

## Implementation

The [Go server](<../../../app/main.go>) owns form handling and storage.
