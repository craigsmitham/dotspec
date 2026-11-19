---
generated_by: dotspec
type: requirement-spec-report
documents: "../../specs/requirements/inspection-preserves-source.requirement.spec.md#inspection-preserves-source-material"
---
# Inspection preserves source material

[dotspec](<../system.spec-report.md>)

## Statement

Inspection must leave source-file contents and paths unchanged, without
creating, deleting, or renaming source files.

## Applies to

Project-wide and focused specification checks, including runs that report
structural errors, writing advice, or operational failures. Source material
includes specifications, consulted system context, and other existing project
files. Concurrent edits by other processes and files explicitly written by
the caller's output redirection are outside this criterion.

## Criterion

Apply the preservation rule in [Statement](<#statement>) to the files covered by
[Applies to](<#applies-to>). Filesystem access-time changes are not changes to source
contents or paths.

## Verification plan

Snapshot the relative paths and byte contents of project files and any system
context outside the selected subtree. Capture command output separately from
that snapshot. Compare snapshots before and after each check, covering valid
and malformed specifications, missing input, and unavailable context. Evaluate
differences against [Criterion](<#criterion>).
Retain the differing paths and contents on test failure.

The [CLI quality-contract tests](<../../internal/cli/quality_contracts_test.go>)
check this criterion. The obligation remains defined in this document.

## Rationale

Contributors must be able to inspect unfinished work without inspection
changing the material they are reviewing.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
describes the read-and-report operation. File creation belongs to the separate
`init` commands.

**Quality attribute** · [Source preservation](<../quality-attributes.spec-report.md#source-preservation>)

[Edit source](<../../specs/requirements/inspection-preserves-source.requirement.spec.md#inspection-preserves-source-material>)
