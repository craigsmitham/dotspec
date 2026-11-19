---
type: use-case-spec
id: discover-specification-sources
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
including findings about unresolved references. They can follow the reported
paths to consult the original files in their native form.

## Main success scenario

1. The contributor selects a project subtree and requests its source inventory.
2. Dotspec discovers specifications, their enclosing context, and direct references.
3. Dotspec reports the source inventory and findings for the contributor to inspect.

## Extensions

### Unresolved reference

At step 2, a reference target cannot be resolved. Report the finding with the
inventory according to the source-discovery contract. The contributor can repair
the link and request another inventory.

### Operational failure

During discovery or reading, an operational failure prevents completing the
inventory. Report the incomplete operation. End without achieving the goal;
the contributor can repair access and retry.

## Sequence

Representative. Discovery order is not a required order of authoring or review.

## Authority and evidence

The [source-discovery reference](../../docs/reference/inspection.md#references) describes link behavior.
The [source-discovery executable examples](../../internal/integration/testdata/sources.txtar)
define required outcomes for their illustrated cases. Contributors use the
inventory to assemble the material for review and decide whether it adequately
specifies the system. The
[source authority principle](../system.spec.md#source-authority) explains how
projects assign roles to those sources.
