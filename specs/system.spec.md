---
type: system
---

# dotspec

## Purpose

Help teams manage software product requirements alongside the code
they describe.

## Scope

The system includes the `dotspec` CLI, the `*.spec.md` files it manages, and
discovery of other specification sources through ordinary file references.

## Goals

Reduce the effort of finding structural problems and missing specification
content during authoring and review.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

## Assumptions

The execution conditions required for repeatability comparisons are defined in
[Repeated inspections](quality/repeatable-inspection.quality.spec.md#applies-to).

## Constraints

Repository development and contribution rules are maintained in
[Contributing](../CONTRIBUTING.md). They describe the development process rather
than implying additional product obligations.

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

### Usability

**ID:** usability

How effectively contributors can author, inspect, and maintain specifications.

### Diagnosability

**ID:** diagnosability
**Sub-attribute of:** usability

How effectively findings help contributors understand and locate specification
problems.

### Accessibility

**ID:** accessibility
**Sub-attribute of:** usability

How effectively people with differing abilities and ways of accessing output
can use dotspec.

### Reliability

**ID:** reliability

How consistently dotspec performs its intended operations under stated
conditions.

### Repeatability

**ID:** repeatability
**Sub-attribute of:** reliability

Whether equivalent inspections under unchanged conditions produce equivalent
reports.

### Source preservation

**ID:** source-preservation
**Sub-attribute of:** reliability

Whether inspection preserves the source material it examines.

### Performance

**ID:** performance

The time and computing resources required to perform operations under stated
workloads.

### Responsiveness

**ID:** responsiveness
**Sub-attribute of:** performance

The elapsed time before an operation delivers its result.

No responsiveness target has been adopted. Selecting one requires representative
workloads, defined measurement conditions, and baseline measurements.
