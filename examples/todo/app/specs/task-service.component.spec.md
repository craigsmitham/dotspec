---
type: component-spec
---
# Task service

## Description

The local service that decides whether a submitted action still applies to the selected entry and records accepted changes.

## Part of

[Personal task list](../../specs/personal-task-list.system.spec.md).

## Responsibilities

A browser tab can still show Open after the owner has moved an errand to Waiting in another tab. The service checks the entry’s current state when applying the form and reports a conflict if the view has changed. It also supplies the current local day when it prepares the Waiting view. [Acknowledge after commit](requirements/acknowledge-after-commit.requirement.spec.md) governs when it may return success.

## Contracts

- The browser uses the [HTML forms](../list.html) to submit changes.

## Implementation

The [Go server](../main.go) owns form handling and storage.
