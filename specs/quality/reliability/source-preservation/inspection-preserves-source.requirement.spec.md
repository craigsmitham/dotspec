---
type: requirement-spec
quality_attributes: [source-preservation]
id: inspection-preserves-source
---

# Inspection preserves source material

## Applies to

Project-wide and focused specification checks, including runs that report
structural errors, writing advice, or operational failures. Source material
includes specifications, consulted system context, and other existing project
files. Concurrent edits by other processes and files explicitly written by
the caller's output redirection are outside this criterion.

## Statement

Inspection must leave source-file contents and paths unchanged, without
creating, deleting, or renaming source files.

## Criterion

Apply the preservation rule in [Statement](#statement) to the files covered by
[Applies to](#applies-to). Filesystem access-time changes are not changes to source
contents or paths.

## Evidence

Checks for this requirement include [CLI quality-contract tests](../../../../internal/cli/quality_contracts_test.go).

## Rationale

Contributors must be able to inspect unfinished work without inspection
changing the material they are reviewing.

## Related specifications

[Specification Inspection](../../../concepts/specification-inspection.concept.spec.md)
describes the read-and-report operation. File creation belongs to the separate
`init` commands.
