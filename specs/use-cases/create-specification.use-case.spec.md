---
type: use-case
id: create-specification
---

# Create a specification draft

## Goal

Obtain an editable starting document for a selected kind of specification.

## System boundary

The dotspec CLI and its creation of a specification file in the selected project.

## Primary actor

A contributor beginning a behavioral account.

## Goal level

User task (sea).

## Success criterion

The contributor receives the path to a newly created draft of the selected kind.

## Minimum guarantees

An existing target is preserved according to the
[reject-existing-target requirement](../requirements/reject-existing-target.requirement.spec.md).

## Main success scenario

1. The contributor selects a project directory, specification kind, and local name.
2. Dotspec creates the draft and reports its path.
3. The contributor opens the document and replaces prompts with known behavior.

## Extensions

### More detailed use-case prompts

At step 1, the contributor selects the expanded use-case template. Dotspec creates
additional context prompts in the same document format, then continues at step 3.

### Target already exists

At step 2, the target exists. Dotspec refuses to overwrite it. This creation attempt
ends without a new draft; the contributor can edit the existing source instead.

## Sequence

Representative. Editing follows creation in this example; existing sources can
be edited without invoking init.

## Authority and evidence

The [init commands](../../docs/cli.md) and linked requirement define command behavior.
The [use-case command examples](../../internal/integration/testdata/use-cases.txtar)
define the illustrated template-selection and non-overwrite outcomes.
