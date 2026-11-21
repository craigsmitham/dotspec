---
type: requirement-spec
id: continue-after-invalid-spec
quality_attributes: [diagnosability]
---

# Continue checking after a structural error

## Applies to

Project-wide specification checks.

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Rationale

One malformed specification must not hide problems in other files.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
provides the conceptual context. The
[report contract](../../docs/cli-handbook/reference/inspection.md#check-specifications) describes report fields
and handling within an individual malformed file.
