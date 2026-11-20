---
type: component-spec
---
# Task service

## Description

The local service that handles browser requests and saves changes to the task list.

## Part of

[Personal task list](../../specs/personal-task-list.system.spec.md).

## Responsibilities

Accept browser changes and persist them locally. See [Acknowledge after commit](requirements/acknowledge-after-commit.requirement.spec.md) for the confirmation guarantee.

## Contracts

- The browser uses the [HTML forms](../list.html) to submit changes.

## Implementation

The [Go server](../main.go) owns form handling and storage.
