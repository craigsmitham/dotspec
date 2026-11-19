---
generated_by: dotspec
type: system-spec-report
documents: "specs/system.spec.md#dotspec"
---
# dotspec

## Purpose

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](<specification/jobs/develop-requirements.md>) and
[Manage Requirements](<specification/jobs/manage-requirements.md>) jobs describe the
work dotspec supports.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<specification/use-cases/discover-specification-sources.md>),
[structural checks](<specification/use-cases/check-specifications.md>), and a
[generated requirements specification](<specification/requirements/read-and-verify-the-generated-requirements-specification.md>).

---

[Requirements](<#requirements>) · [Requirements review](<specification/review.md>)

## Goals

- Reduce the effort of finding structural problems and missing specification
  content during authoring and review.
- Reduce the effort of maintaining requirements as the product changes by
  keeping each obligation with its authoritative source.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

## Assumptions

The execution conditions required for repeatability comparisons are defined in
[Repeated inspections](<specification/requirements/repeatable-inspection.md#applies-to>).

## Constraints

### Human-readable specifications

Prefer ordinary files, Markdown links, and explanatory prose. Add
machine-readable structure only when a demonstrated authoring or inspection
operation needs it.

### Source authority

Projects decide which sources define their obligations. Keep policies,
contracts, and executable specifications in their native form and make them
discoverable through ordinary file links from the relevant specifications.

An executable specification can define the required outcomes for its examples.
Link to those examples where they apply. When a test checks an obligation
defined elsewhere, keep the obligation at its source and explain the test's
verification role beside the link. A link alone confers no authority.

The [contributing guide](<CONTRIBUTING.md#specification-tests-and-implementation-tests>)
explains how this repository assigns those roles to its tests.

## Glossary

### Complete scan

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"specs/system.spec.md#complete-scan"} -->

An inspection that finished its selected scope and required context without an
operational failure. A complete scan can contain structural error findings.

### Structural error

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"specs/system.spec.md#structural-error"} -->

A finding that a specification violates a supported structural rule.

### Writing advice

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"specs/system.spec.md#writing-advice"} -->

A suggestion about specification content. Its exit-status consequences are
specified by [Advice permits success](<specification/requirements/advice-permits-success.md>).

## Quality attributes

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.

| Quality attribute | Description |
| --- | --- |
| [Usability](<specification/quality-attributes.md#usability>) | How effectively contributors can author, inspect, and maintain specifications. |
| [Diagnosability](<specification/quality-attributes.md#diagnosability>) | How effectively findings help contributors understand and locate specification problems. |
| [Accessibility](<specification/quality-attributes.md#accessibility>) | How effectively people with differing abilities and ways of accessing output can use dotspec. |
| [Reliability](<specification/quality-attributes.md#reliability>) | How consistently dotspec performs its intended operations under stated conditions. |
| [Repeatability](<specification/quality-attributes.md#repeatability>) | Whether equivalent inspections under unchanged conditions produce equivalent reports. |
| [Source preservation](<specification/quality-attributes.md#source-preservation>) | Whether inspection preserves the source material it examines. |
| [Performance](<specification/quality-attributes.md#performance>) | The time and computing resources required to perform operations under stated workloads. |
| [Responsiveness](<specification/quality-attributes.md#responsiveness>) | The elapsed time before an operation delivers its result. |

## Concepts

- [Command Invocation](<specification/concepts/cli-interaction.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<specification/concepts/job-map.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<specification/concepts/specification-inspection.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

## Jobs

- [Develop Requirements](<specification/jobs/develop-requirements.md>) — Define an agreed set of requirements for a proposed software change, from understanding the intended change to making an agreed baseline available.
- [Manage Requirements](<specification/jobs/manage-requirements.md>) — Keep an agreed requirements set aligned with a proposed product change, from identifying the change's implications to communicating an updated baseline.

## Quality models

- [FURPS starter](<specification/quality-models/furps.md>)
- [ISO/IEC 25010:2023 starter](<specification/quality-models/iso-25010-2023.md>)
- [Q42 starter](<specification/quality-models/q42.md>)

## Requirements

- [Advice permits a successful check](<specification/requirements/advice-permits-success.md>) — Project-wide and focused specification checks that complete without structural errors or operational failures.
- [Continue checking after a structural error](<specification/requirements/continue-after-invalid-spec.md>) — Project-wide specification checks.
- [Diagnostic severity is available without styling](<specification/requirements/diagnostic-severity.md>) — Human-readable reports from project-wide and focused specification checks.
- [Inspection preserves source material](<specification/requirements/inspection-preserves-source.md>)
- [Isolate specification projects in a workspace](<specification/requirements/isolate-specification-projects-in-a-workspace.md>)
- [Preserve source and represented content in model exports](<specification/requirements/model-source-fidelity.md>)
- [Read and verify the generated requirements specification](<specification/requirements/read-and-verify-the-generated-requirements-specification.md>) — Whole-project builds and checks for projects that enable generated specifications.
- [Reject initialization at an existing target](<specification/requirements/reject-existing-target.md>) — Initialization commands for all supported specification types.
- [Repeated inspections produce equivalent reports](<specification/requirements/repeatable-inspection.md>)

## Use cases

- [Check specifications before committing](<specification/use-cases/check-specifications.md>) — Identify structural problems and missing-content advice in the specifications affected by a proposed change before committing it.
- [Create a specification draft](<specification/use-cases/create-specification.md>) — Obtain an editable starting document for a selected kind of specification.
- [Discover specification sources](<specification/use-cases/discover-specification-sources.md>) — Find specification files and their directly linked context for review.
- [Export a specification model for review or integration](<specification/use-cases/export-specification-model.md>) — Obtain a faithful, structured snapshot of a project's specification content for review or consumption by another tool.
- [Prepare specifications for review](<specification/use-cases/prepare-specifications.md>) — Give reviewers a readable account of proposed behavior and its related sources.

[Edit source](<specs/system.spec.md#dotspec>)
