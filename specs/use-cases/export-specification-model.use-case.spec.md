---
type: use-case-spec
id: export-specification-model
companions:
  - target: ../../internal/integration/testdata/model.txtar
    role: illustrates
    note: Export and offline validation
---

# Export a specification model for review or integration

## Goal

Obtain a faithful, structured snapshot of a project's specification content for
review or consumption by another tool.

## System boundary

The dotspec CLI reading the selected project subtree, enclosing system context
and direct reference targets, then writing the snapshot and validating saved
snapshots against the embedded model contract.

## Primary actor

A contributor preparing specification material for review or tool integration.

## Goal level

User task (sea level).

## Trigger

The contributor requests a specification-model export.

## Preconditions

The contributor has selected an existing project directory. Reading files and
repository discovery may still fail; those outcomes are handled below.

## Supporting actors

Git supplies repository discovery information when the directory belongs to a
Git repository. A downstream tool may consume the JSON snapshot and its schema.

## Stakeholder interests

Contributors need their authored sources and qualifications preserved. Reviewers
need explicit domain roles, provenance, and visible unresolved or uninterpreted
content. Integrators need a checkable representation and an honest account of
which scope was inspected.

## Minimal guarantees

Export preserves source identity, normalized Markdown, and content order according to
the [model source-fidelity requirement](../requirements/model-source-fidelity.requirement.spec.md).
Failed reads and incomplete operations remain visible. References remain authored
occurrences unless resolution or an explicit mapping establishes more.

## Success guarantees

The contributor receives a complete snapshot for the selected scope, with normalized
Markdown and native reference identities, ordered document content, named domain mappings,
located relationship occurrences, resolutions and diagnostics. Structural findings
about authored specifications can be part of this successful information-gathering
outcome. The contributor can validate the saved representation without rereading
its original source paths.

The [model contract](../../docs/reference/specification-model.md) defines command output,
validation and completion distinctions. A conforming snapshot does not establish
that its specifications are accepted, adequate or satisfied.

## Main success scenario

1. The contributor selects a project subtree and requests a model export.
2. Dotspec discovers specifications and their required enclosing context and
   records normalized Markdown and reference identities.
3. Dotspec represents ordered document content, maps explicit domain structures,
   and inspects direct references within the declared boundary.
4. Dotspec writes the snapshot with provenance, diagnostics and operation
   completion information.
5. The contributor saves or passes the snapshot to another tool and validates its
   representation when needed.

## Extensions

### Malformed specification

#### Handling

During steps 2–3, a specification fails syntax or authoring validation. Retain the
available text and locations, report the failure, omit an invalid domain object,
and continue with other files. Resume step 4 with the retained findings. A
complete snapshot containing authoring errors still achieves the review goal.

### Unknown section or unfinished content

#### Handling

During step 3, a section has no named domain mapping, or the author leaves content
empty or handling explicitly unresolved. Preserve its content and appropriate
interpretation state, expose unknown sections and mapping gaps, and continue to
step 4. Do not invent missing behavior or claim stakeholder adequacy.

### Unresolved or uninspected reference

#### Handling

During step 3, a target cannot be resolved or lies outside the traversal boundary.
Retain the occurrence and its resolution state. External resources are not
fetched, and links within direct targets do not trigger a recursive scan.
Continue to step 4 so the contributor can inspect the result.

### Unreadable input or failed discovery

#### Handling

During steps 2–3, an operational failure prevents completion. Retain acquired
material, identify the failed source or operation, and report an incomplete
snapshot when one is available. End without achieving the complete-snapshot
goal; the contributor can repair access and retry.

### Empty scope

#### Handling

At step 2, the selected scope contains no supported specifications. Report a
completed empty snapshot at step 4. This is an alternative successful outcome.

### Invalid saved snapshot

#### Handling

At step 5, validation finds a malformed representation, invalid source location,
invalid reference or invalid source range. Report the contract findings
without changing the snapshot or its source files. The contributor repairs or
regenerates the representation before relying on it.

## Sequence

Representative. The source-fidelity requirement and model contract govern
preservation, traversal, status and validation behavior. Numbered scenario order
does not impose an additional authoring workflow.

## Related specifications

[Specification Inspection](../concepts/specification-inspection.concept.spec.md)
provides the existing inspection context. The source, document-content and domain
layers are explained in [Specification model](../../docs/reference/specification-model.md).

## Authority

The linked source-fidelity requirement specifies source preservation and location fidelity. The [Markdown mapping profile](../../model/markdown-mapping.md)
defines named mappings.
