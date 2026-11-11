---
type: requirement
id: continue-after-invalid-spec
quality_attributes:
  - diagnosability
---

# Continue checking after a structural error

## Applies to

Project-wide specification checks.

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Rationale

One malformed specification must not hide problems in other files.

## Verification plan

Check a project containing two specifications with distinct structural errors.
Verify that the report contains a finding for each specification. Exercise
errors in different supported specification types and positions in scan order.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
provides the conceptual context. The
[report contract](../../README.md#check-specifications) describes report fields
and handling within an individual malformed file.
