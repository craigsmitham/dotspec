---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/discover-specification-sources.use-case.spec.md#discover-specification-sources"
---
# Discover specification sources

[dotspec](<../dotspec.system.spec-report.md>)

## Goal

Find specification files and their directly linked context for review.

## Primary actor

A contributor assembling material for review.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/discover-specification-sources.use-case.spec.md#main-success-scenario"} -->

1. The contributor selects a project subtree and requests its source inventory.
2. Dotspec discovers specifications, their enclosing context, and direct references.
3. Dotspec reports the source inventory and findings for the contributor to inspect.

## Extensions

### Unresolved reference

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/discover-specification-sources.use-case.spec.md#unresolved-reference"} -->

#### Handling

At step 2, a reference target cannot be resolved. Report the finding with the
inventory according to the source-discovery contract. The contributor can repair
the link and request another inventory.

### Operational failure

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/discover-specification-sources.use-case.spec.md#operational-failure"} -->

#### Handling

During discovery or reading, an operational failure prevents completing the
inventory. Report the incomplete operation. End without achieving the goal;
the contributor can repair access and retry.

### Multiple sources have the same slug

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/discover-specification-sources.use-case.spec.md#multiple-sources-have-the-same-slug"} -->

#### Handling

When locating an individual source within a subject, the contributor can give its
typed filename and path. If a slug matches several owned files, dotspec lists
the matching paths so the contributor can select the intended source.

## Success guarantees

The contributor receives a complete source inventory for the selected scope,
including findings about unresolved references. They can follow the reported
paths to consult the original files in their native form.

## System boundary

The dotspec CLI reading the selected project scope and direct reference targets.

## Goal level

User task (sea).

## Sequence

Representative. Discovery order is not a required order of authoring or review.

## Authority

The [source-discovery reference](<../../docs/reference/inspection.md#references>) describes link behavior. The [source layout](<../../docs/reference/layout.md>) defines directory ownership and source selection.
The [source-discovery executable examples](<../../internal/integration/testdata/sources.txtar>)
verify outcomes for their illustrated cases. Contributors use the
inventory to assemble the material for review and decide whether it adequately
specifies the system. The
[contributor guide](<../../CONTRIBUTING.md#specification-tests-and-implementation-tests>) explains how
this project assigns roles to those sources.
