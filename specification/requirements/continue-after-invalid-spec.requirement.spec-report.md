---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/continue-after-invalid-spec.requirement.spec.md#continue-checking-after-a-structural-error"
---
# Continue checking after a structural error

[dotspec](<../dotspec.system.spec-report.md>)

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Applies to

Project-wide specification checks.

## Rationale

One malformed specification must not hide problems in other files.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
provides the conceptual context. The
[report contract](<../../docs/cli-handbook/reference/inspection.md#check-specifications>) describes report fields
and handling within an individual malformed file.

**Quality attribute**

- [Diagnosability](<../quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>)
