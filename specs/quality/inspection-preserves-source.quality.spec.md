---
type: quality
id: inspection-preserves-source
attributes:
  - source-preservation
---

# Inspection preserves source material

## Applies to

Project-wide and focused specification checks, including runs that report
structural errors, writing advice, or operational failures. Source material
includes specifications, consulted system context, and other existing project
files. Concurrent edits by other processes and files explicitly written by
the caller's output redirection are outside this criterion.

## Criterion

Inspection must leave source-file contents and paths unchanged, without
creating, deleting, or renaming source files.

## Measurement

Snapshot the relative paths and byte contents of project files and any system
context outside the selected subtree. Capture command output separately from
that snapshot. Compare snapshots before and after each check, covering valid
and malformed specifications, missing input, and unavailable context.
Retain the differing paths and contents on test failure.

Filesystem access-time changes are not changes to source contents or paths.

## Rationale

Contributors must be able to inspect unfinished work without inspection
changing the material they are reviewing.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
describes the read-and-report operation. File creation belongs to the separate
`init` commands.
