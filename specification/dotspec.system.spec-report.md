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
| [Use cases](<#use-cases>) | [All use cases](<use-cases.md>) |
| [Concepts](<#concepts>) | [Concept definitions](<concepts.md>) |
| [Requirements](<#requirements>) | [All requirements](<requirements.md>) |
| [Quality attributes](<#quality-attributes>) | [All quality attributes](<quality-attributes.md>) · [Quality models](<quality-models.md>) |

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

## Use cases

- [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>) — Identify structural problems and missing-content advice in the specifications affected by a proposed change before committing it.
- [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>) — Obtain an editable starting document for a selected kind of specification.
- [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>) — Find specification files and their directly linked context for review.
- [Export a specification model for review or integration](<use-cases/export-specification-model.use-case.spec-report.md>) — Obtain a faithful, structured snapshot of a project's specification content for review or consumption by another tool.
- [Prepare specifications for review](<use-cases/prepare-specifications.use-case.spec-report.md>) — Give reviewers a readable account of proposed behavior and its related sources.

## Concepts

### Concept definitions

- [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<concepts/job-map.concept.spec-report.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

## Requirements

- [Advice permits a successful check](<requirements/advice-permits-success.requirement.spec-report.md>) — When a completed check reports only advisory findings, dotspec must return exit status zero.
- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>) — When a discovered specification contains a structural error, dotspec must continue inspecting the remaining discovered specifications.
- [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>) — Every diagnostic must identify its severity using a text label independently of terminal styling.
- [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>) — Inspection must leave source-file contents and paths unchanged, without creating, deleting, or renaming source files.
- [Isolate specification projects in a workspace](<requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>) — Each dotspec.yaml declares the name of the specification project rooted in its directory. A spec belongs to the nearest enclosing project. Discovery must exclude nested projects, and implicit System lookup must stop at the project root. Explicit references may acquire context across project boundaries.
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>) — Export must preserve the meaning and order of authored Markdown, including conditions, criteria, reference destinations, and the locations of unsupported or malformed content. Source locations refer to the normalized Markdown included in the snapshot. The original files remain authoritative for exact bytes. A native contract, executable specification, or other non-Markdown target is identified by its path and read status.
- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>) — For an opted-in project with one root system, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and supporting references, and locate the editable source.
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>) — When the target path already exists, initialization must terminate with an unsuccessful exit status.
- [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>) — Repeated inspections must produce equivalent structured reports, including the ordering of array elements.

## Quality attributes

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.

Requirements classified under Usability: [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).
