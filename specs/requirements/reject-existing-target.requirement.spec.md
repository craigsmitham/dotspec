---
type: requirement-spec
id: reject-existing-target
quality_attributes: [source-preservation]
---

# Reject initialization at an existing target

## Applies to

Project initialization with `dotspec init`, and creation with `dotspec new <kind>` for every supported specification kind, including adoption of a named quality attribute that the system already declares.

## Statement

When the target path already exists, initialization must terminate with an
unsuccessful exit status.

## Rationale

Contributors must not mistake an existing specification for a newly created one.
An existing target requires an explicit editing decision.
