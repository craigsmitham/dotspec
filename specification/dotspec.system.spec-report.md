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
- [Scope and environment](<#scope-and-environment>)
- [Quality attributes](<#quality-attributes>)

</details>

<details>
<summary>Browse specifications</summary>

- [Concept definitions](<concepts.md>)
- [Glossary](<glossary.md>)
- [Jobs](<jobs.md>)
- [Quality attributes](<quality-attributes.md>)
- [Quality models](<quality-models.md>)
- [Requirements](<requirements.md>)
- [Use cases](<use-cases.md>)

</details>

## Framing

The idea beheind dotspec is that having a distinctly seperated functional specification for a system and its parts is beneficial for both humans and agents

* Provides shared model and language for describing the intent and shape of a system
* Easy to infer (agents) and comprhened (humans)

dotspec consists of *.spec.md files that are co-located with the system or part/component of the system they specify. These files comprise the foundational "specification" layer for the system. Specs for a system can take many forms (formal methods, executable specifications, etc.), and dotspec is not intended to be the sole source of truth for all specifications, just those it chooses to specify. It is intended to support "progressive specifcitaion" with "just enough" specs to the extent they are helpful.

dotspec provides a rich but simple model for specifying requirements. Not all need to be used, but they can all be used and composed together if desired.

- identify/trace to other sources of truth for specs
- provide affordances for ensuring verification coverage/traceability
- non goal: perform verification of requirements or implementation

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

## Scope and environment

dotspec covers the authoring framework, the `*.spec.md` files teams maintain
alongside their code, and the CLI that discovers and checks those files and
builds a requirements specification from them.

### Operating environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

### Constraints

## Quality attributes

**[Performance](<quality/performance/performance.quality-attribute.spec-report.md>)** — The time and computing resources required to perform operations under stated workloads.

**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)** — How consistently dotspec performs its intended operations under stated conditions.

**[Usability](<quality/usability/usability.quality-attribute.spec-report.md>)** — How effectively contributors can author, inspect, and maintain specifications.

---

[View source](<../specs/dotspec.system.spec.md#dotspec>)

<a id="empty.goals"></a>

<a id="empty.assumptions"></a>
