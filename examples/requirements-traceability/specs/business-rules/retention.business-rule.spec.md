---
type: business-rule
---
# Recovery window

## Statement

Workspace members may recover accidentally deleted documents during the 30 days following deletion.

## Supports

[Reduce lost work](../system.spec.md#reduce-lost-work).

## Rationale

The organization commits to a bounded recovery period. This business rule does not itself specify the storage implementation.

## Applies to

**Targets:** [Document workspace](../system.spec.md)

Applies to workspace members and documents deleted within that workspace.

## Authority

The workspace team's recovery policy establishes this example rule.
