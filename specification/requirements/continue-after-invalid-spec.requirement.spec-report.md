---
generated_by: dotspec
type: requirement-spec-report
---
# Continue checking after a structural error

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Applies to

Project-wide specification checks.

## System verification plan

Check a project containing two specifications with distinct structural errors.
Verify that the report contains a finding for each specification. Exercise
errors in different supported specification types and positions in scan order.

[View source for “System verification plan”](<../../specs/requirements/continue-after-invalid-spec.requirement.spec.md#system-verification-plan>)

## Rationale

One malformed specification must not hide problems in other files.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
provides the conceptual context. The
[report contract](<../../docs/reference/inspection.md#check-specifications>) describes report fields
and handling within an individual malformed file.

**Quality attribute** · [Diagnosability](<../quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>)

[subject-source]: <../../specs/requirements/continue-after-invalid-spec.requirement.spec.md#continue-checking-after-a-structural-error>
