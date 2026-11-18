# Concept design reference

Specification project: **todo**. Project membership does not imply a System relationship.

Generated from authored specifications. Regenerate this view instead of editing it.

Model `0.19.0-draft`; mapping `dotspec-markdown/14`. Relationships reflect explicit declarations. These views do not establish behavioral completeness or implementation guarantees.

## Concepts

### Checklist

[Source](../specs/concepts/checklist.concept.spec.md).

#### Purpose

Keep track of which commitments still need attention.

#### Operational principles

Record an errand, finish it, then discover a follow-up. Reopening the entry brings
the same errand back into the open list.

#### Introduced types

##### Entry

The identity of one recorded commitment.

#### State

Each entry has an identity, a title and an open or completed state.

##### Initial state

A new list has no entries.

##### Invariants

- Every existing entry belongs to exactly one of the two views.

#### Actions

##### Add

Introduce an open entry with the supplied title.

##### Rename

Change an existing entry's title.

###### Inputs

**Selected entry**

**Type:** [Entry](../specs/concepts/checklist.concept.spec.md#entry)

**New title**

**Primitive:** string

###### Guards

- A title contains at least one non-whitespace character.

###### Rejections

A blank title leaves the entry as it was.

##### Complete

Mark an existing entry finished.

##### Reopen

Return a finished entry to the open list.

##### Delete

Remove the selected entry from the list.

#### Effects table

| Operation | Identity | Title | State |
| --- | --- | --- | --- |
| Add | new | supplied title | open |
| Rename | unchanged | supplied title | unchanged |
| Complete | unchanged | unchanged | completed |
| Reopen | unchanged | unchanged | open |
| Delete | removed | removed | removed |

Blank effects cells are unspecified, not unchanged. This authored matrix is not a proof of completeness.

## System compositions

### Personal task list

[Source](../specs/system.spec.md).

| Use | Definition | Type bindings |
| --- | --- | --- |
| My checklist | [Checklist](../specs/concepts/checklist.concept.spec.md) |  |

#### Coordination

Only resolved action references appear as edges. A rule node does not imply a transaction.

```mermaid
flowchart LR
    u0["My checklist"]
```

#### Application dependencies

An arrow means the source use requires the target use for its role in this system. Cycles are permitted. Bindings and action calls do not imply these edges.

```mermaid
flowchart LR
    u0["My checklist"]
```

#### Concept use declarations

##### My checklist

###### Definition

[Checklist](../specs/concepts/checklist.concept.spec.md)

###### Supports

[Resume interrupted work](../specs/jobs/resume.job.spec.md).
