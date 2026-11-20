---
generated_by: dotspec
type: requirement-spec-report
---
# Reject initialization at an existing target

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Statement

When the target path already exists, initialization must terminate with an
unsuccessful exit status.

## Applies to

Initialization commands for all supported specification types.

## System verification plan

For each supported specification type, invoke initialization with an existing
target. Verify the unsuccessful exit status. Cover a regular file, directory,
and symbolic link where the platform supports it. Preserve the target fixture
and capture output to diagnose failures.

[View source for “System verification plan”](<../../specs/requirements/reject-existing-target.requirement.spec.md#system-verification-plan>)

## Rationale

Contributors must not mistake an existing specification for a newly created one.
An existing target requires an explicit editing decision.

## Related specifications

**Quality attribute** · [Source preservation](<../quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>)

[subject-source]: <../../specs/requirements/reject-existing-target.requirement.spec.md#reject-initialization-at-an-existing-target>
