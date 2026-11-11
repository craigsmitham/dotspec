---
type: concept
---
# Trash

## Purpose

Make deletion reversible until a retained item expires.

## Operational principles

If an item is deleted and then restored before it expires, it becomes available again with its retained content.

## State

Items are either available, retained after deletion, or expired.

## Actions

### Delete

Move an available item into retained state with its deletion time.

### Restore

Move a retained, unexpired item into available state.

### Expire

Remove an expired item from retained state.

## Effects table

| Operation | Available(item) | Retained(item) | Deletion time(item) |
| --- | --- | --- | --- |
| Delete | false | true | current time |
| Restore | true | false | unchanged |
| Expire | unchanged | false | |

Blank cells are unspecified; explicit unchanged cells assert preservation.

## Implementation hints

An index could accelerate expiration queries without changing the specified effects.

## Validation plan

### Subject

The proposed reversible-deletion concept; assess the operational principle against intended recovery use.

### Basis

[Restore a document](../use-cases/restore.use-case.spec.md).

### Participants

Workspace members and the recovery-policy stakeholder.

### Method

Review examples and a prototype with the participants, including counterexamples and unsupported situations.

### Success criteria

The participants can explain how the proposal addresses recovery needs, and unresolved mismatches are recorded for revision.

### Conditions

Representative accidental-deletion situations and the stated recovery-window boundary.

### Evidence to retain

Keep example responses and observations, the reviewed revision, the assessment rationale and any disagreements. Record acceptance separately.
