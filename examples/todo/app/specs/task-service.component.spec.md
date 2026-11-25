---
type: component-spec
---
# Task service

## Description

Checks whether a submitted action still applies to the selected entry and records accepted changes.

## Part of

[Personal task list](../../specs/personal-task-list.system.spec.md).

## Responsibilities

A browser tab can still show Open after the owner has moved an errand to Waiting in another tab. When the service receives the old form, it compares the submitting view with the entry’s current state. A mismatch leaves the entry as it was and the browser receives a conflict. [Checklist action effects](../../specs/concepts/checklist.concept.spec.md#effects-table).

The service supplies the current local day when it prepares the Waiting view. [Acknowledge after commit](requirements/acknowledge-after-commit.requirement.spec.md) governs when it may report a successful change.

## Contracts

- The browser submits changes through [Task list forms](../../specs/personal-task-list.system.spec.md#task-list-forms), whose [form template](../list.html) supplies the actions and fields.

## Implementation

The [Go server](../main.go) owns form handling and storage.
