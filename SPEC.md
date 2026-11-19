---
type: requirements-specification
generated_by: dotspec
subject_kind: system
---
# dotspec

## Purpose

Help teams develop and maintain software product requirements alongside the
code they describe. Teams can describe what people need, what the system
promises, and how to assess those promises.

The [Develop Requirements](<specification/jobs/develop-requirements.md>) and
[Manage Requirements](<specification/jobs/manage-requirements.md>) jobs describe the
work dotspec supports.

## Scope

The system includes the authoring framework, the `dotspec` CLI, and the
`*.spec.md` files it manages. The CLI supports
[source discovery](<specification/use-cases/discover-specification-sources.md>),
[structural checks](<specification/use-cases/check-specifications.md>), and a
[generated requirements specification](<specification/requirements/read-and-verify-the-generated-requirements-specification.md>).

[Contents](<specification/index.md>) · [Review](<specification/review.md>)

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
[Repeated inspections](<specification/requirements/repeatable-inspection.md#applies-to>).

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

The [contributing guide](<CONTRIBUTING.md#specification-tests-and-implementation-tests>)
explains how this repository assigns those roles to its tests.

## Glossary

### Complete scan

An inspection that finished its selected scope and required context without an
operational failure. A complete scan can contain structural error findings.

### Structural error

A finding that a specification violates a supported structural rule.

### Writing advice

A suggestion about specification content. Its exit-status consequences are
specified by [Advice permits success](<specification/requirements/advice-permits-success.md>).

## Quality attributes

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.

### Usability

**ID:** usability

How effectively contributors can author, inspect, and maintain specifications.

**Classifies** · [Read and verify the generated requirements specification](<specification/requirements/read-and-verify-the-generated-requirements-specification.md>)

### Diagnosability

**ID:** diagnosability
**Sub-attribute of:** usability

How effectively findings help contributors understand and locate specification
problems.

**Classifies** · [Continue checking after a structural error](<specification/requirements/continue-after-invalid-spec.md>)

**Classifies** · [Diagnostic severity is available without styling](<specification/requirements/diagnostic-severity.md>)

### Accessibility

**ID:** accessibility
**Sub-attribute of:** usability

How effectively people with differing abilities and ways of accessing output
can use dotspec.

**Classifies** · [Diagnostic severity is available without styling](<specification/requirements/diagnostic-severity.md>)

### Reliability

**ID:** reliability

How consistently dotspec performs its intended operations under stated
conditions.

### Repeatability

**ID:** repeatability
**Sub-attribute of:** reliability

Whether equivalent inspections under unchanged conditions produce equivalent
reports.

**Classifies** · [Read and verify the generated requirements specification](<specification/requirements/read-and-verify-the-generated-requirements-specification.md>)

**Classifies** · [Repeated inspections produce equivalent reports](<specification/requirements/repeatable-inspection.md>)

### Source preservation

**ID:** source-preservation
**Sub-attribute of:** reliability

Whether inspection preserves the source material it examines.

**Classifies** · [Inspection preserves source material](<specification/requirements/inspection-preserves-source.md>)

**Classifies** · [Preserve source and represented content in model exports](<specification/requirements/model-source-fidelity.md>)

**Classifies** · [Reject initialization at an existing target](<specification/requirements/reject-existing-target.md>)

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

## Explore this system

- [Concepts](<specification/concepts/index.md>)
- [Jobs](<specification/jobs/index.md>)
- [Quality models](<specification/quality-models/index.md>)
- [Requirements](<specification/requirements/index.md>)
- [Use cases](<specification/use-cases/index.md>)
- [Full specification contents](<specification/index.md>)
- [Requirements review](<specification/review.md>)

[Edit source](<specs/system.spec.md#dotspec>)
