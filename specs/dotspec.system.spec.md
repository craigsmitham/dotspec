---
type: system-spec
glossary_sources:
  - ../docs/reference/model/glossary.md
---

# dotspec

## Description

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

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

dotspec covers the authoring framework, the `*.spec.md` files teams maintain
alongside their code, and the CLI that discovers and checks those files and
builds a requirements specification from them.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.
