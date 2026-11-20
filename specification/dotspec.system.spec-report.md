---
generated_by: dotspec
type: system-spec-report
---
# dotspec

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

[Requirements review](<requirements-review.report.md>) · [View source][subject-source]

| Section | Explore |
| --- | --- |
| [Goals](<#goals>) | Purpose · Objectives |
| [User needs](<#user-needs>) | [Jobs](<jobs.md>) |
| [Scope](<#scope>) | [Glossary](<glossary.md>) |
| [Requirements](<#requirements>) | [Use cases](<use-cases.md>) · [Concept definitions](<concepts.md>) · [All requirements](<requirements.md>) |
| [Quality attributes](<#quality-attributes>) | [Attribute index](<quality-attributes.md>) · [Quality models](<quality-models.md>) |

## Goals

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) and
[Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) jobs describe the
work dotspec supports.

[View source for “Purpose”](<../specs/dotspec.system.spec.md#purpose>)

### Objectives

- Reduce the effort of finding structural problems and missing specification
  content during authoring and review.
- Reduce the effort of maintaining requirements as the product changes by
  keeping each obligation with its authoritative source.

## User needs

[Jobs](<jobs.md>) — The work people need to accomplish and the outcomes they seek.

- [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) — Define an agreed set of requirements for a proposed software change, from understanding the intended change to making an agreed baseline available.
- [Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) — Keep an agreed requirements set aligned with a proposed product change, from identifying the change's implications to communicating an updated baseline.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<use-cases/discover-specification-sources.use-case.spec-report.md>),
[structural checks](<use-cases/check-specifications.use-case.spec-report.md>), and a
[generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).

[View source for “Scope”](<../specs/dotspec.system.spec.md#scope>)

### Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

### Assumptions

The execution conditions required for repeatability comparisons are defined in
[Repeated inspections](<requirements/repeatable-inspection.requirement.spec-report.md#applies-to>).

### Constraints

#### Human-readable specifications

Prefer ordinary files, Markdown links, and explanatory prose. Add
machine-readable structure only when a demonstrated authoring or inspection
operation needs it.

#### Source authority

Projects decide which sources define their obligations. Keep policies,
contracts, and executable specifications in their native form and make them
discoverable through ordinary file links from the relevant specifications.

An executable specification can define the required outcomes for its examples.
Link to those examples where they apply. When a test checks an obligation
defined elsewhere, keep the obligation at its source and explain the test's
verification role beside the link. A link alone confers no authority.

The [contributing guide](<../CONTRIBUTING.md#specification-tests-and-implementation-tests>)
explains how this repository assigns those roles to its tests.

[View source for “Source authority”](<../specs/dotspec.system.spec.md#source-authority>)

[Glossary](<glossary.md>) — Definitions of terms used in this system.

## Requirements

### Behavior

[Use cases](<use-cases.md>) — Actor goals, interactions, and outcomes.

### Concepts

[Concept definitions](<concepts.md>) — Definitions of purposeful behavior, state, and actions.

- [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<concepts/job-map.concept.spec-report.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

### Obligations

[All requirements](<requirements.md>) — Standalone requirements and embedded obligations, including routes by component.

## Quality attributes

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.

[View source for “Quality attributes”](<../specs/dotspec.system.spec.md#quality-attributes>)

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

- [Responsiveness](<quality/performance/responsiveness/responsiveness.quality-attribute.spec-report.md>) — The elapsed time before an operation delivers its result.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

- [Repeatability](<quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>) — Whether equivalent inspections under unchanged conditions produce equivalent reports.
- [Source preservation](<quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>) — Whether inspection preserves the source material it examines.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.

- [Accessibility](<quality/usability/accessibility/accessibility.quality-attribute.spec-report.md>) — How effectively people with differing abilities and ways of accessing output can use dotspec.
- [Diagnosability](<quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>) — How effectively findings help contributors understand and locate specification problems.

[Attribute index](<quality-attributes.md>)

[Quality models](<quality-models.md>) — Reference classifications and their definitions.

[subject-source]: <../specs/dotspec.system.spec.md#dotspec>
