---
generated_by: dotspec
type: concept-spec-report
source: "../../specs/concepts/checklist.concept.spec.md#checklist"
---
# Checklist

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Purpose

Decide what I can do next, which replies I need to chase and which commitments I have already settled.

## Operational principles

I promise to get the sink fixed. Once I've asked the plumber for a visit time, I can turn to another errand and still find the reply I need when I return. If the plumber asks me for a photo, I resume the same repair entry. After the repair I can complete it. If the sink leaks again, I can reopen that entry without having to remember its title or reconstruct the old wait.

## State

Each entry has an identity, a title and one state. Open means I can take the next action. Waiting means I still own the result and have recorded a reply I need. Completed means I consider the commitment settled.

### Initial state

A new list has no entries.

### Invariants

- Each entry appears in exactly one view, named for its state.
- Only Waiting entries carry a reply description and follow-up day. The day determines when to check back under [When a reply needs chasing](<../business-rules/when-a-reply-needs-chasing.business-rule.spec-report.md>).
- Actions select entries by identity. Two errands can have the same title. Acting on one leaves the other alone.

## Actions

### Add

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#add"} -->

Record a commitment under a title. Ignore surrounding whitespace. Reject an empty or whitespace-only title without creating an entry.

### Rename

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#rename"} -->

Change how the person will recognize the entry on their next visit.

#### Inputs

The selected entry and its new title.

#### Guards

- The new title contains at least one non-whitespace character. Ignore surrounding whitespace.

#### Rejections

A blank title leaves the entry as it was.

### Wait for reply

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#wait-for-reply"} -->

Set an open entry aside until the person wants to check back. Record what reply they expect and the calendar day they choose for follow-up.

#### Guards

- The reply description contains text after trimming surrounding whitespace.
- The follow-up day is a real calendar date. A person who remembers a missed follow-up can enter a day in the past.

#### Rejections

If either detail is missing or invalid, the entry stays Open.

### Update follow-up

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#update-follow-up"} -->

Revise a wait after chasing the reply. Start with the saved description and day, change either or both, and keep the commitment in Waiting. The same [guards as Wait for reply](<#wait-for-reply>) apply. An invalid edit leaves both saved details intact.

### Resume

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#resume"} -->

Take back the next action on a waiting commitment when the person can move the work forward.

### Complete

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#complete"} -->

Mark an open or waiting commitment settled. The [completion rule](<../business-rules/completion-preserves-a-commitment.business-rule.spec-report.md>) governs retention.

### Reopen

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#reopen"} -->

Return a completed entry to Open when there is more to do.

### Delete

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/checklist.concept.spec.md#delete"} -->

Discard the selected entry from any view. Adding another errand with the same title creates a new identity.

## Effects table

Each row describes one accepted action. For an existing entry, the view where the person submits the action must match the entry's current state. A mismatch or a deleted entry causes the action to be rejected without changing the list.

| Operation | Identity | Title | State | Reply and follow-up day |
| --- | --- | --- | --- | --- |
| Add | new | supplied title | open | absent |
| Rename | unchanged | supplied title | unchanged | unchanged |
| Wait for reply | unchanged | unchanged | waiting | supplied details |
| Update follow-up | unchanged | unchanged | waiting | revised details |
| Resume | unchanged | unchanged | open | cleared |
| Complete | unchanged | unchanged | completed | cleared |
| Reopen | unchanged | unchanged | open | absent |
| Delete | removed | removed | removed | removed |

## Related specifications

**Used by**

- [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
