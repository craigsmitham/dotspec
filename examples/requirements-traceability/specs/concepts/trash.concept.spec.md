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
