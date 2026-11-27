---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/build-generated-specification.use-case.spec.md#build-the-generated-specification"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Use cases](<../dotspec.system.spec-report.md#use-cases>)

# Build the generated specification

## Goal

Give readers one navigable Markdown account of a project's specifications that
matches the authored sources.

## Primary actor

A [contributor](<../dotspec.system.spec-report.md#contributor>) preparing specifications
for readers.

## Trigger

The contributor has changed a specification, or a check reports that the
generated specification is missing or stale.

## Preconditions

The contributor works in a project declared with `dotspec.yaml`.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#main-success-scenario"} -->

1. The contributor requests a build of the whole project.
2. Dotspec inspects every specification in the project, as a
   [check](<check-specifications.use-case.spec-report.md>) would.
3. Dotspec writes `SPEC.md` and the pages beneath `specification/`, removes
   obsolete generated pages, and reports what changed.
4. The contributor reads the result from `SPEC.md` and commits it with its
   sources.

## Extensions

[No project declared](<#no-project-declared>) · [Generated specification turned off](<#generated-specification-turned-off>) · [Part of a project selected](<#part-of-a-project-selected>) · [Source has error findings](<#source-has-error-findings>) · [Authored file at a generated destination](<#authored-file-at-a-generated-destination>) · [Nothing changed](<#nothing-changed>)

### No project declared

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#no-project-declared"} -->

#### Handling

At step 1, no `dotspec.yaml` encloses the working directory. Dotspec writes
nothing and directs the contributor to
[start a specification project](<start-specification-project.use-case.spec-report.md>).

### Generated specification turned off

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#generated-specification-turned-off"} -->

#### Handling

At step 1, the manifest sets `build.enabled` to `false`. Dotspec writes nothing
and names the setting.

### Part of a project selected

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#part-of-a-project-selected"} -->

#### Handling

At step 1, the contributor selects a directory below the project root. Dotspec
writes nothing and names the project, because a build covers the whole project.

### Source has error findings

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#source-has-error-findings"} -->

#### Handling

At step 2, a source has error findings. Dotspec writes nothing, states that
`SPEC.md` was not written and how many source errors prevent it, and reports the
findings as a check would. The contributor repairs the sources and repeats
step 1.

### Authored file at a generated destination

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#authored-file-at-a-generated-destination"} -->

#### Handling

At step 3, an authored file occupies a planned destination. Dotspec reports the
collision and leaves the authored file in place. The contributor moves it,
updates links to it and repeats step 1.

### Nothing changed

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/build-generated-specification.use-case.spec.md#nothing-changed"} -->

#### Handling

At step 3, the sources are unchanged since the last build. Dotspec leaves the
generated files as they are and reports that nothing was written or removed.

## Success guarantees

The generated specification presents every specification in the project as the
[generated specification requirement](<../requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>)
describes, and a whole-project check confirms that it is current.

## Minimal guarantees

Nothing is written while a source has error findings, and no authored file is
replaced.

## System boundary

The dotspec CLI reading every specification in a project declared with
`dotspec.yaml` and writing that project's generated specification: `SPEC.md`
and the pages beneath `specification/`.

## Goal level

User task (sea level).

## Sequence

Representative. Building and checking can repeat in any useful order.

## Authority

The [build command](<../../docs/reference/cli.md#dotspec-build>) and the linked
generated specification requirement define command behavior.

## Companions

**Illustrated by**

- [init.txtar](<../../internal/integration/testdata/init.txtar>) — A build without a project, a first build and the check that confirms it
- [check.py](<../../scripts/check-example/check.py>) — Rebuilding the To-do generated specification and following its links

---

[View source](<../../specs/use-cases/build-generated-specification.use-case.spec.md#build-the-generated-specification>)
