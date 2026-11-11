---
type: use-case
---
# Restore a document

## Goal

Recover an accidentally deleted document and resume work.

## System boundary

**Subject system:** [Document workspace](../system.spec.md)

The workspace owns retention and restoration; external backups are outside this boundary.

## Supports

[Recover lost work](../jobs/recover.job.spec.md).

## Involves

[DocumentTrash](../system.spec.md#documenttrash).

## Main success scenario

1. The member selects a retained document.
2. The workspace restores it according to [Restore retained content](#restore-content).
   - **Illustrates:** [Restore retained content](#restore-content), [Restore](../concepts/trash.concept.spec.md#restore).
3. The member resumes work.

## Success criterion

The member can resume work on the restored document.

## Obligations

### Restore retained content

**ID:** restore-content

When an authorized workspace member requests restoration before the 30-day recovery period ends, the workspace must restore the retained document content and metadata.

#### Derived from

[Recovery window](../business-rules/retention.business-rule.spec.md) supplies the time boundary.

#### Supports

[Recover lost work](../jobs/recover.job.spec.md).

#### Applies to

**Targets:** [DocumentTrash](../system.spec.md#documenttrash)

#### Verification plan

Exercise restoration with retained content and metadata; compare both to their pre-deletion values. Include requests immediately before the recovery deadline.

## Authority

For this example, the workspace team accepts the named obligation above as the authoritative behavioral expression. The scenario refers to it rather than defining a second rule.

## Primary actor

An authorized workspace member.

## Goal level

User task.

## Sequence

The scenario order is representative; the named obligation governs restoration.

## Extensions

### The recovery deadline has passed

The workspace declines restoration from retained state and explains that the recovery window has ended.

### The member lacks authorization

The workspace declines restoration and leaves the retained document unchanged.

## Validation plan

### Subject

The proposed interaction and alternatives; assess whether they support members in resuming interrupted work.

### Basis

[Recover lost work](../jobs/recover.job.spec.md).

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
