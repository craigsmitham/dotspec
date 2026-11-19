---
generated_by: dotspec
type: component-spec-report
documents: "../../specs/system.spec.md#task-service"
---
# Task service

[Personal task list](<../../SPEC.md>)

## Part of

[Personal task list](<../../SPEC.md>).

## Responsibilities

Accept browser changes and commit them to the local database before confirming them.

## Contracts

- The browser uses the [HTML forms](<../../app/list.html>) to submit changes.

## Implementation

The [Go server](<../../app/main.go>) owns form handling and storage.

[Edit source](<../../specs/system.spec.md#task-service>)
