---
type: use-case-spec
id: discover-specification-sources
companions:
  - target: ../../internal/integration/testdata/sources.txtar
    role: illustrates
    note: Source discovery through ordinary references and companions
  - target: ../../internal/cli/kinds_test.go
    role: illustrates
    note: Lists of every specification kind, grouped by system
---

# Discover specification sources

## Goal

Find specification files and their directly linked context for review.

## System boundary

The dotspec CLI reading the selected project scope and direct reference targets.

## Primary actor

A contributor assembling material for review.

## Goal level

User task (sea).

## Success guarantees

The contributor receives a complete source inventory for the selected scope,
including findings about unresolved references and, for each file a specification
lists as a companion, the declaring specification and the role it states. They can
follow the reported paths to consult the original files in their native form.

## Main success scenario

1. The contributor selects a project subtree and requests its source inventory.
2. Dotspec discovers specifications, their enclosing context, and direct references.
3. Dotspec reports the source inventory and findings for the contributor to inspect.

## Extensions

### Unresolved reference

#### Handling

At step 2, a reference target cannot be resolved. Report the finding with the
inventory according to the source-discovery contract. The contributor can repair
the link and request another inventory.

### Operational failure

#### Handling

During discovery or reading, an operational failure prevents completing the
inventory. Report the incomplete operation. End without achieving the goal;
the contributor can repair access and retry.

### List specifications of one kind by system

#### Handling

At step 1, the contributor needs only the specifications of one kind, such as
the project's requirements. Dotspec lists the files of that kind in the enclosing
project, grouped by the system each belongs to, with each file's path. Systems are
listed by directory and components by their Part of links. The contributor can
pass a listed path to a focused check or open the file, then continue at step 3
with the inventory of that kind.

### Two sources have the same slug

#### Handling

At step 3, two files of one kind share a slug in different directories. Each
listing and finding identifies a file by its source path, so the contributor
selects the intended source by that path.

## Sequence

Representative. Discovery order is not a required order of authoring or review.

## Authority

The [source-discovery reference](../../docs/reference/inspection.md#references) describes link behavior. The [source layout](../../docs/reference/layout.md) defines directory ownership and source selection.
Contributors use the inventory to assemble the material for review and decide
whether it adequately specifies the system. The
[contributor guide](../../CONTRIBUTING.md#specification-tests-and-implementation-tests) explains how
this project assigns roles to those sources.
