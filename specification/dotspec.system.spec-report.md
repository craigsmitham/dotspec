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
| User needs | [Jobs](<jobs.md>) |
| [Scope](<#scope>) | [All assumptions](<assumptions.md>) · [Glossary](<glossary.md>) |
| Use cases | [All use cases](<use-cases.md>) |
| Concepts | [Concept definitions](<concepts.md>) |
| Requirements | [All requirements](<requirements.md>) |
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

- [Maintain requirements as the product changes](<goals.md#maintain-requirements-as-the-product-changes>) — Reduce the effort of maintaining requirements as the product changes by keeping each requirement with its authoritative source.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<use-cases/discover-specification-sources.use-case.spec-report.md>),
[structural checks](<use-cases/check-specifications.use-case.spec-report.md>), and a
[generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).

The specification records stakeholder needs, product requirements, the conditions
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

Each authored specification declares one primary subject of its stated type. That file owns the subject and its named parts. A use case owns the interaction and its guarantees. A concept owns its actions and effects. Each product requirement has one requirement specification. Other specifications cite that file or declare a relationship to it.

External policies and native contracts govern their own provisions. Cite the relevant provision where it applies. Executable examples record inputs and assertions. Link them to the specifications they exercise so reviewers can compare expected behavior with observed results.

The [contributing guide](<../CONTRIBUTING.md#specification-tests-and-implementation-tests>)
explains how this repository assigns those roles to its tests.

## Quality attributes

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.
