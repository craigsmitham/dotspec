# Concept design reference

Specification project: **requirements-review**. Project membership does not imply a System relationship.

Generated from authored specifications. Regenerate this view instead of editing it.

Model `0.16.0-draft`; mapping `dotspec-markdown/13`. Relationships reflect explicit declarations. These views do not establish behavioral completeness or implementation guarantees.

## Inspection findings

- **advice** `relationship.motivation`: No explicit motivation or derivation is recorded for this obligation; review its justification.
- **advice** `relationship.measurement_basis`: A measure or measurement plan is recorded, but no separately identified criterion, conditions or applicability is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result.
- **advice** `relationship.verification_plan`: No system verification or measurement plan is recorded for this obligation; review how a realization will be checked.
- **advice** `relationship.motivation`: No explicit motivation or derivation is recorded for this obligation; review its justification.
- **advice** `quality.applies-to`: Identify the subject and applicable conditions, workload, or population.
- **advice** `quality.measurement-basis`: A measure or measurement plan is recorded, but no separately identified criterion, conditions or applicability is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result.
- **advice** `requirement.verification`: Record assessment guidance for checking a realization against the obligation; a plan is not a verification result.

## Concepts

### Trash

[Source](../../examples/requirements-review/specs/concepts/trash.concept.spec.md).

#### Purpose

Make deletion reversible until a retained item expires.

#### Validation plan (needs and intended use)

##### Subject

The proposed reversible-deletion concept; assess the operational principle against intended recovery use.

##### Basis

[Restore a document](../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md).

##### Participants

Workspace members and the recovery-policy stakeholder.

##### Method

Review examples and a prototype with the participants, including counterexamples and unsupported situations.

##### Success criteria

The participants can explain how the proposal addresses recovery needs, and unresolved mismatches are recorded for revision.

##### Conditions

Representative accidental-deletion situations and the stated recovery-window boundary.

##### Evidence to retain

Keep example responses and observations, the reviewed revision, the assessment rationale and any disagreements. Record acceptance separately.

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
