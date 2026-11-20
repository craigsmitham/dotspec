---
type: system-spec
---

# dotspec

## Description

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

## Purpose

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](jobs/develop-requirements.job.spec.md) and
[Manage Requirements](jobs/manage-requirements.job.spec.md) jobs describe the
work dotspec supports.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](use-cases/discover-specification-sources.use-case.spec.md),
[structural checks](use-cases/check-specifications.use-case.spec.md), and a
[generated requirements specification](requirements/generated-specification.requirement.spec.md).

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
[Repeated inspections](quality/reliability/repeatability/repeatable-inspection.requirement.spec.md#applies-to).

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

The [contributing guide](../CONTRIBUTING.md#specification-tests-and-implementation-tests)
explains how this repository assigns those roles to its tests.

## Glossary

### Complete scan

An inspection that finished its selected scope and required context without an
operational failure. A complete scan can contain structural error findings.

### Structural error

A finding that a specification violates a supported structural rule.

### Writing advice

A suggestion about specification content. Its exit-status consequences are
specified by [Advice permits success](requirements/advice-permits-success.requirement.spec.md).

## Quality attributes

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.
