---
generated_by: dotspec
type: system-spec-report
source: "../specs/dotspec.system.spec.md#dotspec"
---
# dotspec

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

<details>
<summary>On this page</summary>

**Intent**\
[Purpose](<#purpose>) · [Overview](<#overview>) · [Stakeholders](<#stakeholders>) · [Jobs](<#jobs>)

**Context**\
[Scope](<#scope>) · [Environment](<#environment>)

**Behavior**\
[Use cases](<#use-cases>)

[Quality Attributes](<#quality-attributes>)

[Requirements](<#requirements>)

**Design**\
[Concepts](<#concepts>)

</details>

<details>
<summary>Browse specifications</summary>

- [Companions](<companions.md>)
- [Glossary](<glossary.md>)

</details>

## Purpose

Help teams develop and maintain software requirements alongside the
code they describe.

## Overview

A specification records product intent for one independently maintained subject. Each `*.spec.md` file declares one root subject and owns its parts. A use case owns its steps and extensions. A requirement that several use cases rely on has its own declaration.

The specifications in product scope and the relationships among them compose the product specification that dotspec manages. Authors can develop it incrementally and keep each declaration alongside the product or component it describes. Generated reports let people and agents read those decisions together.

A specification can cite an external contract for the provision it adopts, a reusable quality model for a term, or an executable test for evidence. Each retains its own authority. dotspec helps readers trace those sources and inspect verification coverage. Assessing the running implementation takes additional evidence.

### Inspiration

Chad Fowler’s writing on regenerative software inspires dotspec’s practice of
recording which behaviors a component’s replacement must preserve. A team can
return to those requirements when it rewrites the component.
[The Phoenix Architecture](<https://chadfowler.com/regenerative-software/>)

Joel Spolsky’s case for functional specifications motivates dotspec’s use of
ordinary prose. Describing a feature gives teammates and customers something
concrete to question while the design is still cheap to change.
[Painless Functional Specifications – Part 1: Why Bother?](<https://www.joelonsoftware.com/2000/10/02/painless-functional-specifications-part-1-why-bother/>)

Hillel Wayne explains how several implementations can satisfy the same
specification. This gives dotspec authors a way to judge when a requirement
needs more detail. If a reasonable implementation could satisfy the wording
and still disappoint its users, refine the requirement. When the remaining
choices all meet the need, leave them to the implementer.
[A sufficiently comprehensive spec is not (necessarily) code](<https://buttondown.com/hillelwayne/archive/a-sufficiently-comprehensive-spec-is-not/>)

## Stakeholders

### Contributor

<!-- dotspec-doc {"type":"system-stakeholder-role-spec-report","documents":"../specs/dotspec.system.spec.md#contributor"} -->

A contributor authors or changes specifications and needs findings that identify the passage to revise.

### Reviewer

<!-- dotspec-doc {"type":"system-stakeholder-role-spec-report","documents":"../specs/dotspec.system.spec.md#reviewer"} -->

A reviewer examines a proposed change and needs to know which sources were inspected and where the evidence for a requirement comes from.

## Jobs

- [Develop Requirements](<jobs/develop-requirements.job.spec-report.md>) — Define an agreed set of requirements for a proposed software change, from understanding the intended change to making an agreed baseline available.
- [Manage Requirements](<jobs/manage-requirements.job.spec-report.md>) — Keep an agreed requirements set aligned with a proposed product change, from identifying the change's implications to communicating an updated baseline.

## Scope

dotspec covers the authoring framework, the `*.spec.md` files teams maintain
alongside their code, and the CLI that discovers and checks those files and
builds a requirements specification from them.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

## Use cases

- [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>) — Identify structural problems and missing-content advice in the specifications affected by a proposed change before committing it.
- [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>) — Obtain an editable starting document for a selected kind of specification.
- [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>) — Find specification files and their directly linked context for review.
- [Export a specification model for review or integration](<use-cases/export-specification-model.use-case.spec-report.md>) — Obtain a faithful, structured snapshot of a project's specification content for review or consumption by another tool.
- [Prepare specifications for review](<use-cases/prepare-specifications.use-case.spec-report.md>) — Give reviewers a readable account of proposed behavior and its related sources.

## Quality Attributes

- [Accessibility](<quality/usability/accessibility/accessibility.quality-attribute.spec-report.md>) — How effectively people with differing abilities and ways of accessing output can use dotspec.
- [Diagnosability](<quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>) — How effectively findings help contributors understand and locate specification problems.
- [Performance](<quality/performance/performance.quality-attribute.spec-report.md>) — The time and computing resources required to perform operations under stated workloads.
- [Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>) — How consistently dotspec performs its intended operations under stated conditions.
- [Repeatability](<quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>) — Whether equivalent inspections under unchanged conditions produce equivalent reports.
- [Responsiveness](<quality/performance/responsiveness/responsiveness.quality-attribute.spec-report.md>) — The elapsed time before an operation delivers its result.
- [Source preservation](<quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>) — Whether inspection preserves the source material it examines.
- [Usability](<quality/usability/usability.quality-attribute.spec-report.md>) — How effectively contributors can author, inspect, and maintain specifications.

## Requirements

- [Advice permits a successful check](<requirements/advice-permits-success.requirement.spec-report.md>) — When a completed check reports only advisory findings, dotspec must return exit status zero.
- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>) — When a discovered specification contains a structural error, dotspec must continue inspecting the remaining discovered specifications.
- [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>) — Every diagnostic must identify its severity using a text label independently of terminal styling.
- [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>) — A contributor must be able to reach the relevant versioned guidance from a specification type, a question in their own words, or a mapped diagnostic, then read the selected page, section, or worked example offline. …
- [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>) — Inspection must leave source-file contents and paths unchanged, without creating, deleting, or renaming source files.
- [Isolate specification projects in a workspace](<requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>) — Each dotspec.yaml declares the name of the specification project rooted in its directory. …
- [List the companions that carry a specification's intent](<requirements/list-the-companions-that-carry-a-specification-s-intent.requirement.spec-report.md>) — An author must be able to list, in a specification's frontmatter, each artifact that carries part of its intent in another form, such as a formal model, an executable example, a test or an explanatory document, …
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>) — Export must preserve the meaning and order of authored Markdown, including conditions, criteria, reference destinations, and the locations of unsupported or malformed content. …
- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>) — For a project with one root system and reports enabled, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. …
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>) — When the target path already exists, initialization must terminate with an unsuccessful exit status.
- [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>) — Repeated inspections must produce equivalent structured reports, including the ordering of array elements.

## Concepts

- [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>) — Let a person or program submit an operation and determine its outcome through a consistent process interface.
- [Job Map](<concepts/job-map.concept.spec-report.md>) — Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.
- [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>) — Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

<a id="empty.goals"></a>

<a id="empty.assumptions"></a>

<a id="empty.constraints"></a>

---

[View source](<../specs/dotspec.system.spec.md#dotspec>) · [Quality models](<quality-models.md>)
