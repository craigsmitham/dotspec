---
type: use-case-spec
id: build-generated-specification
verification:
  - name: Rebuilding the To-do generated specification and following its links
    source: ../../scripts/check-example/check.py
---

# Build the generated specification

## Goal

Give readers one navigable Markdown account of a project's specifications that
matches the authored sources.

## System boundary

The dotspec CLI reading every specification in a project declared with
`dotspec.yaml` and writing its generated reports beneath `specification/` and updating the selected agent instruction files.

## Primary actor

A [contributor](../dotspec.system.spec.md#contributor) preparing specifications
for readers.

## Goal level

User task (sea level).

## Trigger

The contributor has changed a specification, or a check reports that the
generated specification is missing or stale.

## Preconditions

The contributor works in a project declared with `dotspec.yaml`.

## Minimal guarantees

Nothing is written while a source has error findings, and no authored file is
replaced.

## Success guarantees

The generated specification presents every specification in the project as the
[generated specification requirement](../requirements/generated-specification.requirement.spec.md)
describes, and a whole-project check confirms that it is current.

## Main success scenario

1. The contributor requests a build of the whole project.
2. Dotspec inspects every specification in the project, as a
   [check](check-specifications.use-case.spec.md) would.
3. Dotspec writes the reports and updates the agent instruction links, removes
   obsolete generated pages, and reports what changed.
4. The contributor reads the result from the entry point named in the output and commits it with its
   sources.

## Extensions

### No project declared

#### Handling

At step 1, no `dotspec.yaml` encloses the working directory. Dotspec writes
nothing and directs the contributor to
[start a specification project](start-specification-project.use-case.spec.md).

### Generated specification turned off

#### Handling

At step 1, the manifest sets `build.enabled` to `false`. Dotspec writes nothing
and names the setting.

### Part of a project selected

#### Handling

At step 1, the contributor selects a directory below the project root. Dotspec
writes nothing and names the project, because a build covers the whole project.

### Source has error findings

#### Handling

At step 2, a source has error findings. Dotspec writes nothing, states that
the specification was not written and how many source errors prevent it, and reports the
findings as a check would. The contributor repairs the sources and repeats
step 1.

### Authored file at a generated destination

#### Handling

At step 3, an authored file occupies a planned destination. Dotspec reports the
collision and leaves the authored file in place. The contributor moves it,
updates links to it and repeats step 1.

### Nothing changed

#### Handling

At step 3, the sources are unchanged since the last build. Dotspec leaves the
generated files as they are and reports that nothing was written or removed.

## Sequence

Representative. Building and checking can repeat in any useful order.

## Authority

The [build command](../../docs/reference/cli.md#dotspec-build) and the linked
generated specification requirement define command behavior.

[A build without a project, a first build and the check that confirms it](../../internal/integration/testdata/init.txtar).
