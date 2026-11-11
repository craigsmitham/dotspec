---
type: requirement
---
# Preserve document metadata on restoration

## Statement

For the restoration required by Restore retained content, the restored metadata must preserve the original title, author identity, creation time and workspace identity.

## Refines

[Restore retained content](../use-cases/restore.use-case.spec.md#restore-content).

The refinement makes the required metadata explicit without weakening the parent's recovery window or authorization condition.

## Relies on

[Trusted time](../system.spec.md#trusted-time).

## Applies to

**Targets:** [Restore a document](../use-cases/restore.use-case.spec.md), [DocumentTrash](../system.spec.md#documenttrash)

## Rationale

Preserving these metadata fields lets members resume work with the same document identity and attribution. The Refines section identifies the governing obligation.

## Validation plan

### Subject

The proposed obligation to preserve all four metadata fields; assess whether it correctly expresses recovery needs.

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
