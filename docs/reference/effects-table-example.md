# Concept design reference

Generated from authored specifications. Regenerate this view instead of editing it.

Model `0.12.0-draft`; mapping `dotspec-markdown/11`. Relationships reflect explicit declarations. These views do not establish behavioral completeness or implementation guarantees.

## Inspection findings

- **advice** `relationship.verification_plan`: No verification or measurement plan is recorded for this obligation; review how it will be checked.
- **advice** `requirement.verification`: Describe how to assess the statement and retain the evidence.

## Concepts

### Trash

[Source](../../examples/requirements-review/specs/concepts/trash.concept.spec.md).

#### Purpose

Make deletion reversible until a retained item expires.

#### Operational principles

If an item is deleted and then restored before it expires, it becomes available again with its retained content.

#### State

Items are either available, retained after deletion, or expired.

#### Actions

##### Delete

Move an available item into retained state with its deletion time.

##### Restore

Move a retained, unexpired item into available state.

##### Expire

Remove an expired item from retained state.

#### Effects table

| Operation | Available(item) | Retained(item) | Deletion time(item) |
| --- | --- | --- | --- |
| Delete | false | true | current time |
| Restore | true | false | unchanged |
| Expire | unchanged | false | |

Blank cells are unspecified; explicit unchanged cells assert preservation.

Blank effects cells are unspecified, not unchanged. This authored matrix is not a proof of completeness.

#### Implementation suggestions (nonbinding)

An index could accelerate expiration queries without changing the specified effects.

## System compositions

### Document workspace

[Source](../../examples/requirements-review/specs/system.spec.md).

| Use | Definition | Type bindings |
| --- | --- | --- |
| DocumentTrash | [Trash](../../examples/requirements-review/specs/concepts/trash.concept.spec.md) |  |

#### Coordination

Only resolved action references appear as edges. A rule node does not imply a transaction.

```mermaid
flowchart LR
    u0["DocumentTrash"]
```

#### Application dependencies

An arrow means the source use requires the target use for its role in this system. Cycles are permitted. Bindings and action calls do not imply these edges.

```mermaid
flowchart LR
    u0["DocumentTrash"]
```

#### Concept use declarations

##### DocumentTrash

###### Definition

[Trash](../../examples/requirements-review/specs/concepts/trash.concept.spec.md)

###### Supports

[Recover lost work](../../examples/requirements-review/specs/jobs/recover.job.spec.md).

###### Relies on

[Trusted time](../../examples/requirements-review/specs/system.spec.md#trusted-time) bounds retention intervals.
