---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/reject-existing-target.requirement.spec.md#reject-initialization-at-an-existing-target"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../dotspec.system.spec-report.md#requirements>)

# Reject initialization at an existing target

## Statement

When the target path already exists, initialization must terminate with an
unsuccessful exit status.

## Applies to

Project initialization with `dotspec init`, and creation with `dotspec new <kind>` for every supported specification kind, including adoption of a named quality attribute that the system already declares.

## Rationale

Contributors must not mistake an existing specification for a newly created one.
An existing target requires an explicit editing decision.

## Related specifications

**Quality attribute**

- [Source preservation](<../quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>)

---

[View source](<../../specs/requirements/reject-existing-target.requirement.spec.md#reject-initialization-at-an-existing-target>)
