---
generated_by: dotspec
type: system-spec-report
source: "../specs/dotspec.system.spec.md#dotspec"
---
# dotspec

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

| Section | Explore |
| --- | --- |
| [Purpose](<#purpose>) | Help teams develop and maintain software product requirements alongside the code they describe. Teams can describe what people need, what the system promises, and how to assess those promises. |
| [Goals](<#goals>) | [All goals](<goals.md>) |
| [User needs](<#user-needs>) | [Jobs](<jobs.md>) |
| [Scope](<#scope>) | [All assumptions](<assumptions.md>) · [Glossary](<glossary.md>) |
| Use cases | [All use cases](<use-cases.md>) |
| [Concepts](<#concepts>) | [Concept definitions](<concepts.md>) |
| Requirements | [All requirements](<requirements.md>) |
| [Quality attributes](<#quality-attributes>) | [Attribute index](<quality-attributes.md>) · [Quality models](<quality-models.md>) |

## Purpose

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) and
[Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) jobs describe the
work dotspec supports.

## Goals

- [Find gaps during review](<goals.md#find-gaps-during-review>) — Reduce the effort of finding structural problems and missing specification content during authoring and review.

- [Maintain requirements as the product changes](<goals.md#maintain-requirements-as-the-product-changes>) — Reduce the effort of maintaining requirements as the product changes by keeping each obligation with its authoritative source.

## User needs

- [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) — Define an agreed set of requirements for a proposed software change, from understanding the intended change to making an agreed baseline available.
- [Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) — Keep an agreed requirements set aligned with a proposed product change, from identifying the change's implications to communicating an updated baseline.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<use-cases/discover-specification-sources.use-case.spec-report.md>),
[structural checks](<use-cases/check-specifications.use-case.spec-report.md>), and a
[generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).

The specification records stakeholder needs, product obligations, the conditions
in which they apply, and the criteria for satisfaction. Teams maintain test
procedures, validation studies, research records, and issue tracking in their
own tools or documents. Authors can link to them beside the requirements they
inform or assess.

### Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

### Assumptions

- [Repeatable inspection conditions](<assumptions.md#repeatable-inspection-conditions>) — The execution conditions required for repeatability comparisons are defined in Repeated inspections.

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

## Concepts

- [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<concepts/job-map.concept.spec-report.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

## Quality attributes

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.

Requirements classified under Usability: [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).
