---
type: use-case-spec
id: adopt-quality-attribute
companions:
  - target: ../../internal/integration/testdata/quality-convention.txtar
    role: illustrates
    note: Adoption into the nearest system, refusal of a declared attribute and parent reuse
  - target: ../../internal/integration/testdata/quality-models.txtar
    role: illustrates
    note: Preview, bundled and local models, and a term with two parents
---

# Adopt a quality attribute for a system

## Goal

Give a system a declaration for a quality that matters to it, so that its
requirements can be classified by that quality and its authors can say what the
quality means for this system.

## System boundary

The dotspec CLI reading the nearest system that encloses the working directory
and, when one is named, a quality model, then writing quality attribute
declarations into that system.

## Primary actor

A [contributor](../dotspec.system.spec.md#contributor) specifying the qualities
of one system.

## Goal level

User task (sea level).

## Preconditions

The working directory lies inside the system that will own the attributes.

## Minimal guarantees

Existing declarations are preserved according to the
[reject-existing-target requirement](../requirements/reject-existing-target.requirement.spec.md).
A failed or previewed adoption writes nothing.

## Success guarantees

The system declares each requested attribute in its own file, and the
contributor knows each file's path and how to read the guidance for defining the
attribute. An attribute copied from a quality model cites the term it came from,
and its parents are declared by the same system.

## Main success scenario

1. The contributor names one or more attribute IDs and, optionally, the quality
   model to copy them from, and requests their adoption.
2. Dotspec finds the nearest enclosing system and plans one declaration per
   requested ID. With a model, each declaration copies the model term's
   definition and cites it, and the plan adds the term's missing parents.
3. Dotspec writes the declarations and reports each path with a direct command
   for the [authoring guidance](../requirements/offline-authoring-guidance.requirement.spec.md).
4. The contributor replaces the definition with what the quality means for this
   system, then classifies requirements by its ID.

## Extensions

### Preview the adoption

#### Handling

At step 1, the contributor asks for a preview. At step 3, dotspec reports each
planned decision and the content of each new file, and writes nothing. The
contributor can repeat step 1 without the preview.

### Attribute already declared

#### Handling

At step 2, the system already declares a requested ID. Dotspec refuses to
overwrite it, writes nothing and ends unsuccessfully. The contributor edits the
existing declaration instead.

### Parent already declared

#### Handling

At step 2, the system already declares a parent of a requested term. Dotspec
reuses that declaration, with its wording and grouping, and plans no copy of
it. Continue at step 3.

### Unknown term or unavailable model

#### Handling

At step 2, the named model cannot be read or has no term with a requested ID.
Dotspec writes nothing and reports the model it could not read, or the unknown
ID with where to read the model's terms. The contributor corrects the ID or
model and repeats step 1, or omits the model to start from a blank definition.

### No enclosing system

#### Handling

At step 2, no system encloses the working directory. Dotspec writes nothing and
reports how to start one. The contributor selects a directory inside a system
and repeats step 1.

## Sequence

Representative. A contributor can also write a quality attribute declaration
by hand in the system's quality directory.

## Authority

The [`new quality-attribute` command](../../docs/reference/cli.md#dotspec-new-quality-attribute),
the [adoption reference](../../docs/reference/inspection.md#adopt-quality-attributes)
and the linked requirements define command behavior.
