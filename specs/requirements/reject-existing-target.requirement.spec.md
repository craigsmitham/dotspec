---
type: requirement-spec
id: reject-existing-target
quality_attributes: [source-preservation]
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
