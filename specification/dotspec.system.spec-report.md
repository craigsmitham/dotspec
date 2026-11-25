---
generated_by: dotspec
type: system-spec-report
source: "../specs/dotspec.system.spec.md#dotspec"
---
# dotspec

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

<details>
<summary>On this page</summary>

- [Framing](<#framing>)
- [Purpose](<#purpose>)
- [Stakeholders](<#stakeholders>)
- [Scope](<#scope>)
- [Environment](<#environment>)
- [Concepts](<#concepts>)
- [Quality Attributes](<#quality-attributes>)

</details>

<details>
<summary>Browse specifications</summary>

- [Glossary](<glossary.md>)
- [Jobs](<jobs.md>)
- [Requirements](<requirements.md>)
- [Use cases](<use-cases.md>)

</details>

## Framing

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

## Purpose

Help teams develop and maintain software requirements alongside the
code they describe.

## Stakeholders

### Contributor

<!-- dotspec-doc {"type":"system-stakeholder-role-spec-report","documents":"../specs/dotspec.system.spec.md#contributor"} -->

A contributor authors or changes specifications and needs findings that identify the passage to revise.

[View section source](<../specs/dotspec.system.spec.md#contributor>)

### Reviewer

<!-- dotspec-doc {"type":"system-stakeholder-role-spec-report","documents":"../specs/dotspec.system.spec.md#reviewer"} -->

A reviewer examines a proposed change and needs to know which sources were inspected and where the evidence for a requirement comes from.

[View section source](<../specs/dotspec.system.spec.md#reviewer>)

## Scope

dotspec covers the authoring framework, the `*.spec.md` files teams maintain
alongside their code, and the CLI that discovers and checks those files and
builds a requirements specification from them.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

## Concepts

### [Command Invocation](<concepts/cli-interaction.concept.spec-report.md>)

Let a person or program submit an operation and determine its outcome through a consistent process interface.

### [Job Map](<concepts/job-map.concept.spec-report.md>)

Make an accomplishment's steps and success criteria inspectable, so people can discuss what the work involves and how success is judged.

### [Specification Inspection](<concepts/specification-inspection.concept.spec-report.md>)

Make structural problems and missing-content advice visible within an explicit inspection scope, so contributors can decide what to revise.

## Quality Attributes

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

- [Responsiveness](<quality/performance/responsiveness/responsiveness.quality-attribute.spec-report.md>) — The elapsed time before an operation delivers its result.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

- [Repeatability](<quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>) — Whether equivalent inspections under unchanged conditions produce equivalent reports.
- [Source preservation](<quality/reliability/source-preservation/source-preservation.quality-attribute.spec-report.md>) — Whether inspection preserves the source material it examines.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.

- [Accessibility](<quality/usability/accessibility/accessibility.quality-attribute.spec-report.md>) — How effectively people with differing abilities and ways of accessing output can use dotspec.
- [Diagnosability](<quality/usability/diagnosability/diagnosability.quality-attribute.spec-report.md>) — How effectively findings help contributors understand and locate specification problems.

<a id="empty.goals"></a>

<a id="empty.assumptions"></a>

<a id="empty.constraints"></a>

---

[View source](<../specs/dotspec.system.spec.md#dotspec>) · [Quality models](<quality-models.md>)
