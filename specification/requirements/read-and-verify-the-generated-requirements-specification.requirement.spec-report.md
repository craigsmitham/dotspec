---
generated_by: dotspec
type: requirement-spec-report
documents: "../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification"
---
# Read and verify the generated requirements specification

[dotspec](<../dotspec.system.spec-report.md>)

## Statement

For an opted-in project with one root system, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and verification guidance, and locate the editable source.

## Applies to

Whole-project builds and checks for projects that enable generated specifications.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. From that report, reach the same promise through its feature and through the system’s Requirements section. Both routes explain why the promise is listed and lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

## Verification plan

Run the [To-do reader journey and regression probes](<../../scripts/check-example/check.py>). Inspect the artifact and diagnostics using the [generated specification guide](<../../docs/reference/generated-specification.md>).

## Rationale

Reviewers can inspect product intent in Git and follow a promise to its context and evidence without learning the repository layout.

## Related specifications

**Quality attribute**

- [Repeatability](<../quality-attributes.spec-report.md#repeatability>)
- [Usability](<../quality-attributes.spec-report.md#usability>)

[Edit source](<../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification>)
