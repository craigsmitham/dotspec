---
type: system-spec
glossary_sources:
  - ../docs/reference/model/glossary.md
---

# dotspec

## Description

dotspec is a requirements framework and CLI for teams that record product decisions in Markdown, formal models, executable examples, diagrams, and native contracts. Its generated report brings those decisions together as one product specification.

## Purpose

Help teams elicit, analyze, specify, validate, verify, and maintain requirements alongside their code. The goal is to make it inexpensive to confirm that a promise meets a stakeholder’s need, verify that an implementation keeps it, and revise it when the need changes.

## Overview

A product’s intent layer records what the software must do and why. It consists of the decisions that would have to survive if the implementation were discarded and rebuilt.

A specification records product intent for one independently maintained subject. Each `*.spec.md` file declares one root subject and owns its parts. A use case owns its steps and extensions. A requirement that several use cases rely on has its own declaration.

The specifications in product scope and the relationships among them compose the product specification that dotspec manages. Authors can develop it incrementally and keep each declaration alongside the product or component it describes. The generated specification lets people and agents read those decisions together.

Authors designate one governing source for each decision. A Markdown declaration can name a schema that defines an interface, a state diagram that defines permitted transitions, or an executable example that defines behavior. An author can also cite a diagram to explain a transition or a procedure to assess a promise. The declaration tells readers which role the artifact plays.

### Inspiration

Chad Fowler’s writing on regenerative software inspires dotspec’s practice of
recording which behaviors a component’s replacement must preserve. A team can
return to those requirements when it rewrites the component.
[The Phoenix Architecture](https://chadfowler.com/regenerative-software/)

Joel Spolsky’s case for functional specifications motivates dotspec’s use of
ordinary prose. Describing a feature gives teammates and customers something
concrete to question while the design is still cheap to change.
[Painless Functional Specifications – Part 1: Why Bother?](https://www.joelonsoftware.com/2000/10/02/painless-functional-specifications-part-1-why-bother/)

Hillel Wayne explains how several implementations can satisfy the same
specification. This gives dotspec authors a way to judge when a requirement
needs more detail. If a reasonable implementation could satisfy the wording
and still disappoint its users, refine the requirement. When the remaining
choices all meet the need, leave them to the implementer.
[A sufficiently comprehensive spec is not (necessarily) code](https://buttondown.com/hillelwayne/archive/a-sufficiently-comprehensive-spec-is-not/)

## Stakeholders

### Contributor

A contributor authors or changes specifications and needs findings that identify the passage to revise.

### Reviewer

A reviewer examines a proposed change and needs to know which sources were inspected and where the evidence for a requirement comes from.

## Scope

dotspec covers the authoring framework, the Markdown declarations teams maintain alongside their code, and the CLI that catalogs their governing sources, checks the declarations and references, and builds the navigable specification. Teams keep native specifications in their own formats and use their own tools to assess them. dotspec records the procedures and observations they supply.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.
