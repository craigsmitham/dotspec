---
type: concept-spec
---
# Checklist

## Purpose

Keep track of which commitments still need attention.

## Operational principles

Record an errand, finish it, then discover a follow-up. Reopening the entry brings
the same errand back into the open list.

## State

Each entry has an identity, a title and an open or completed state.

### Initial state

A new list has no entries.

### Invariants

- Every existing entry belongs to exactly one of the two views.

## Introduced types

### Entry

The identity of one recorded commitment.

## Actions

### Add

Introduce an open entry with the supplied title.

### Rename

Change an existing entry's title.

#### Inputs

##### Selected entry

**Type:** [Entry](#entry)

##### New title

**Primitive:** string

#### Guards

- A title contains at least one non-whitespace character.

#### Rejections

A blank title leaves the entry as it was.

### Complete

Mark an existing entry finished.

### Reopen

Return a finished entry to the open list.

### Delete

Remove the selected entry from the list.

## Effects table

| Operation | Identity | Title | State |
| --- | --- | --- | --- |
| Add | new | supplied title | open |
| Rename | unchanged | supplied title | unchanged |
| Complete | unchanged | unchanged | completed |
| Reopen | unchanged | unchanged | open |
| Delete | removed | removed | removed |
