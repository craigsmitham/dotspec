---
generated_by: dotspec
type: component-spec-report
documents: "../../../components/task-service/specs/task-service.component.spec.md#task-service"
---
# Task service

## Description

The local service that handles browser requests and saves changes to the task list.

[Personal task list](<../../personal-task-list.system.spec-report.md>)

## Part of

[Personal task list](<../../personal-task-list.system.spec-report.md>).

## Responsibilities

Accept browser changes and persist them locally. See [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>) for the confirmation guarantee.

## Contracts

- The browser uses the [HTML forms](<../../../app/list.html>) to submit changes.

## Implementation

The [Go server](<../../../app/main.go>) owns form handling and storage.

## Applicable requirements

- [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>)

## Requirements

- [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>) — The task service must commit a change to the local database before acknowledging it to the browser.

[Edit source](<../../../components/task-service/specs/task-service.component.spec.md#task-service>)
