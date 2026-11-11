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

## Verification

Compare all four metadata fields before deletion and after restoration for each supported document type. Record the implementation revision and test results outside this example.

## Rationale

Preserving these metadata fields lets members resume work with the same document identity and attribution. The Refines section identifies the governing obligation.
