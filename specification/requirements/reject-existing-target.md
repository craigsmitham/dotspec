---
type: specification-section
generated_by: dotspec
subject_kind: requirement
---
# Reject initialization at an existing target

[dotspec](<../../SPEC.md>) · [Contents](<../index.md>)

## Statement

When the target path already exists, initialization must terminate with an
unsuccessful exit status.

## Applies to

Initialization commands for all supported specification types.

## Rationale

Contributors must not mistake an existing specification for a newly created one.
An existing target requires an explicit editing decision.

## Verification plan

For each supported specification type, invoke initialization with an existing
target. Verify the unsuccessful exit status. Cover a regular file, directory,
and symbolic link where the platform supports it. Preserve the target fixture
and capture output to diagnose failures.

## Related specifications

- Quality attribute — [Source preservation](<../../SPEC.md#source-preservation>)

[Edit source](<../../specs/requirements/reject-existing-target.requirement.spec.md#reject-initialization-at-an-existing-target>)
