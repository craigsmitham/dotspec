---
generated_by: dotspec
type: system-spec-report
documents: "../specs/dotspec.system.spec.md#dotspec"
---
# dotspec

## Description

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

## Purpose

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) and
[Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) jobs describe the
work dotspec supports.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<use-cases/discover-specification-sources.use-case.spec-report.md>),
[structural checks](<use-cases/check-specifications.use-case.spec-report.md>), and a
[generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>).

---

[Requirements](<#requirements>) · [Requirements review](<requirements-review.report.md>)

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
[Repeated inspections](<quality/reliability/repeatability/repeatable-inspection.requirement.spec-report.md#applies-to>).

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

The [contributing guide](<../CONTRIBUTING.md#specification-tests-and-implementation-tests>)
explains how this repository assigns those roles to its tests.

## Glossary

### Complete scan

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../specs/dotspec.system.spec.md#complete-scan"} -->

An inspection that finished its selected scope and required context without an
operational failure. A complete scan can contain structural error findings.

### Structural error

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../specs/dotspec.system.spec.md#structural-error"} -->

A finding that a specification violates a supported structural rule.

### Writing advice

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../specs/dotspec.system.spec.md#writing-advice"} -->

A suggestion about specification content. Its exit-status consequences are
specified by [Advice permits success](<requirements/advice-permits-success.requirement.spec-report.md>).

## Quality attributes

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.

| Quality attribute | Description |
| --- | --- |
| [Performance](<quality/performance/performance.quality-attribute.spec-report.md>) | The time and computing resources required to perform operations under stated workloads. |
| [Responsiveness](<quality/performance/responsiveness/responsiveness.quality-attribute.spec-report.md>) | The elapsed time before an operation delivers its result. |
| [Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>) | How consistently dotspec performs its intended operations under stated conditions. |
| [Repeatability](<quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>) | Whether equivalent inspections under unchanged conditions produce equivalent reports. |
| [Source preservation](<quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>) | Whether inspection preserves the source material it examines. |
| [Accessibility](<quality/usability/accessibility/accessibility.quality-attribute.spec-report.md>) | How effectively people with differing abilities and ways of accessing output can use dotspec. |
| [Diagnosability](<quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>) | How effectively findings help contributors understand and locate specification problems. |
| [Usability](<quality/usability/usability.quality-attribute.spec-report.md>) | How effectively contributors can author, inspect, and maintain specifications. |

## Concepts

- [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<concepts/job-map.concept.spec-report.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

## Jobs

- [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) — Define an agreed set of requirements for a proposed software change, from understanding the intended change to making an agreed baseline available.
- [Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) — Keep an agreed requirements set aligned with a proposed product change, from identifying the change's implications to communicating an updated baseline.

## Quality models

- [FURPS starter](<quality-models/furps.quality-model.spec-report.md>) — A dotspec starter covering the five FURPS categories and all 32 subcategory placements (31 distinct subcategories) listed in the cited Rational Unified Process guidance. This is a documented RUP variant, not a claim that FURPS has one universal fixed taxonomy. Explanations are original authoring prompts. Functionality can classify general requirements as well as quality requirements; these categories do not impose targets. The separate FURPS+ constraint categories are outside this starter. Accuracy appears under two categories in the source and is represented once, organized under both reliability and performance.
- [ISO/IEC 25010:2023 starter](<quality-models/iso-25010-2023.quality-model.spec-report.md>) — A dotspec starter for reviewing the nine product-quality characteristics and all 40 sub-characteristics in ISO/IEC 25010:2023. These short explanations are original authoring prompts, not standard definitions or a conformance checklist. Consult the standard for authoritative definitions and guidance.
- [Q42 starter](<quality-models/q42.quality-model.spec-report.md>) — A dotspec starter with all nine Q42 dimensions and 191 distinct qualities from the pinned upstream revision. Dimension tags are represented as Sub-attribute of groupings, preserving all overlaps without implying exclusive decomposition, inherited obligations, or equivalent meanings. The 37 alias records are retained as alternative-name prose on their canonical qualities, not duplicate attributes. Names and groupings are adapted from arc42 Q42; explanatory questions are original dotspec authoring prompts. Upstream definitions, requirements, approaches, and standards are available through the source links.

## Requirements

- [Advice permits a successful check](<requirements/advice-permits-success.requirement.spec-report.md>) — When a completed check reports only advisory findings, dotspec must return exit status zero.
- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>) — When a discovered specification contains a structural error, dotspec must continue inspecting the remaining discovered specifications.
- [Diagnostic severity is available without styling](<quality/usability/diagnosability/diagnostic-severity.requirement.spec-report.md>) — Every diagnostic must identify its severity using a text label independently of terminal styling.
- [Inspection preserves source material](<quality/reliability/source-preservation/inspection-preserves-source.requirement.spec-report.md>) — Inspection must leave source-file contents and paths unchanged, without creating, deleting, or renaming source files.
- [Isolate specification projects in a workspace](<requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>)
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>)
- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>)
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>) — When the target path already exists, initialization must terminate with an unsuccessful exit status.
- [Repeated inspections produce equivalent reports](<quality/reliability/repeatability/repeatable-inspection.requirement.spec-report.md>) — Repeated inspections must produce equivalent structured reports, including the ordering of array elements.

## Use cases

- [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>) — Identify structural problems and missing-content advice in the specifications affected by a proposed change before committing it.
- [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>) — Obtain an editable starting document for a selected kind of specification.
- [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>) — Find specification files and their directly linked context for review.
- [Export a specification model for review or integration](<use-cases/export-specification-model.use-case.spec-report.md>) — Obtain a faithful, structured snapshot of a project's specification content for review or consumption by another tool.
- [Prepare specifications for review](<use-cases/prepare-specifications.use-case.spec-report.md>) — Give reviewers a readable account of proposed behavior and its related sources.

[Edit source](<../specs/dotspec.system.spec.md#dotspec>)
