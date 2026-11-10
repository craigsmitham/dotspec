---
type: requirement
id: reject-existing-target
attributes:
  - source-preservation
---

# Reject initialization at an existing target

## Applies to

Initialization commands for all supported specification types.

## Statement

When the target path already exists, initialization must terminate with an
unsuccessful exit status.

## Rationale

Contributors must not mistake an existing specification for a newly created one.
An existing target requires an explicit editing decision.

## Verification

For each supported specification type, invoke initialization with an existing
target. Verify the unsuccessful exit status. Cover a regular file, directory,
and symbolic link where the platform supports it. Preserve the target fixture
and capture output to diagnose failures.
