---
generated_by: dotspec
type: concept-spec-report
source: "../../specs/concepts/checklist.concept.spec.md#checklist"
---
# Checklist

[Personal task list](<../personal-task-list.system.spec-report.md>)

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

## Actions

### Add

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#add"} -->

Introduce an open entry with the supplied title.

### Rename

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#rename"} -->

Change an existing entry's title.

#### Inputs

The selected entry and its new title.

#### Guards

- A title contains at least one non-whitespace character.

#### Rejections

A blank title leaves the entry as it was.

### Complete

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#complete"} -->

Mark an existing entry finished.

### Reopen

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#reopen"} -->

Return a finished entry to the open list.

### Delete

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#delete"} -->

Remove the selected entry from the list.

## Effects table

| Operation | Identity | Title | State |
| --- | --- | --- | --- |
| Add | new | supplied title | open |
| Rename | unchanged | supplied title | unchanged |
| Complete | unchanged | unchanged | completed |
| Reopen | unchanged | unchanged | open |
| Delete | removed | removed | removed |

## Related specifications

**Used by** · [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
