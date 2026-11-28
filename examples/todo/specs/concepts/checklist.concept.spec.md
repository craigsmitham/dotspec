---
type: concept-spec
specifications:
  - name: Renaming an entry by identity
    defines: concepts
    source: rename.als#rename
    verification:
      - name: Check preservation of other titles
        source: rename.als#OtherTitlesSurvive
---
# Checklist

## Purpose

Decide what I can do next, which replies I need to chase and which commitments I have already settled.

## Operational principles

I promise to get the sink fixed. Once I've asked the plumber for a visit time and chosen a day to chase the reply, I can turn to another errand. When I return, I still find the reply I need and the day I chose. If the plumber asks me for a photo, I resume the same repair entry. After the repair I can complete it. If the sink leaks again, I can reopen that entry without having to remember its title or reconstruct the old wait. The [repair example](repair.md) illustrates that sequence.

## State

Each entry has an identity, a title and one state. Open means I can take the next action. Waiting means I still own the result and have recorded a reply I need. Completed means I consider the commitment settled.

### Initial state

A new list has no entries.

### Invariants

- Each entry appears in exactly one view, named for its state.
- Every Waiting entry carries both an expected reply and a follow-up day. Open and Completed entries carry neither. The day determines when to check back under [When a reply needs chasing](../business-rules/follow-up.business-rule.spec.md).
- Actions select entries by identity. Two errands can have the same title. Acting on one leaves the other alone.

## Actions

### Add

Record a commitment under a title. Ignore surrounding whitespace. Reject an empty or whitespace-only title without creating an entry.

### Rename

Change how the person will recognize the entry on their next visit. The [rename model](#specification.renaming-an-entry-by-identity) defines the effect on titles when several entries have the same name.

#### Inputs

The selected entry and its new title.

#### Guards

- The new title contains at least one non-whitespace character. Ignore surrounding whitespace.

#### Rejections

A blank title leaves the entry as it was.

### Wait for reply

Set an open entry aside until the person wants to check back. Record what reply they expect and the calendar day they choose for follow-up.

#### Guards

- The reply description contains text after trimming surrounding whitespace.
- The follow-up day is a real calendar date. A person who remembers a missed follow-up can enter a day in the past.

#### Rejections

If either detail is missing or invalid, the entry stays Open.

### Update follow-up

Revise a wait after chasing the reply. Start with the saved description and day, change either or both, and keep the commitment in Waiting. The same [guards as Wait for reply](#wait-for-reply) apply. An invalid edit leaves both saved details intact.

### Resume

Take back the next action on a waiting commitment when the person can move the work forward.

### Complete

Mark an open or waiting commitment settled. If the person selected Complete in Open and the entry is now Waiting, reject the old action. They can still choose Complete from Waiting. See the [shared check for actions on an entry](#effects-table) and the [completion rule](../business-rules/completion.business-rule.spec.md).

### Reopen

Return a completed entry to Open when there is more to do.

### Delete

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
