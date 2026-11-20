# Specification layout

## Project layout

Each part of the project can have a `specs/` directory. It describes its parent
project directory; folders inside `specs/` organize specifications and do not create
subsystems or new scopes.

```text
specs/
  dotspec.system.spec.md
  jobs/
    develop-requirements.job.spec.md
    manage-requirements.job.spec.md
  use-cases/
    check-specifications.use-case.spec.md
  concepts/
    job-map.concept.spec.md
    specification-inspection.concept.spec.md
  requirements/
    specification-project-isolation.requirement.spec.md
  quality/
    reliability/
      reliability.quality-attribute.spec.md
      repeatability/
        repeatability.quality-attribute.spec.md
        repeatable-inspection.requirement.spec.md
internal/
  parser/
    parser.go
    specs/
      spec-parser.component.spec.md
      jobs/
        interpret-specification.job.spec.md
```

The layout is required:

- A directory can describe one system or component. Its declaration lives at `specs/<name>.system.spec.md` or `specs/<name>.component.spec.md`. Derive `<name>` from the document heading using the [subject naming rule](formats/system.md#match-the-filename-to-the-heading). Both subjects can own the collections below.
- System quality attributes live in `specs/quality/<attribute>/<attribute>.quality-attribute.spec.md`. Every nested folder identifies a subattribute. A quality requirement that sets an expectation for that attribute can live beside its definition.
- Jobs are `specs/jobs/<slug>.job.spec.md`, directly inside the plural `jobs/` folder.
- Use cases are `specs/use-cases/<slug>.use-case.spec.md`, directly inside the
  plural `use-cases/` folder.
- Concepts are `specs/concepts/<slug>.concept.spec.md`, directly inside the
  plural `concepts/` folder.
- Functional requirements live in the owning system or component’s `specs/requirements/<slug>.requirement.spec.md`. Use `quality_attributes` to link relevant qualities. For example, the To-do restart requirement lives in `requirements/` and links to Durability.
- Business rules are `specs/business-rules/<slug>.business-rule.spec.md`, and
  quality models are `specs/quality-models/<slug>.quality-model.spec.md`,
  directly inside their plural folders.
- CLI namespaces and file type suffixes remain singular: `dotspec job` and
  `.job.spec.md`, `dotspec use-case` and `.use-case.spec.md`, `dotspec concept`
  and `.concept.spec.md`, or `dotspec requirement` and `.requirement.spec.md`.
  Collections use plural names such as `jobs` and `concepts`.
- Create folders only when needed. `init` creates the required spec directories;
  the selected project directory must already exist.

`--path` always selects a project directory, such as `internal/parser`, rather
than `internal/parser/specs` or its `jobs/` folder. Using a metadata directory as
a project scope is an error. A nested `specs/` tree inside another `specs/` tree
is not a separate scope.

Adjacent spec files and singular collection folders are not supported. Move an
existing `<name>.system.spec.md` to `specs/<name>.system.spec.md`, and move each adjacent
`<slug>.job.spec.md` to `specs/jobs/<slug>.job.spec.md`. Discovery rejects misplaced
files of the requested type with an expected-location diagnostic. Local lookup
and initialization also reject misplaced specs of that type directly in the
selected project directory or anywhere within its `specs/` tree. Having both
old and new locations is an error; there is no location precedence or automatic
migration. Revise any relative file links when moving documents.

## Independent projects and monorepos

Put a `dotspec.yaml` containing `name: your-project` in each independent
project's root. Default commands select the nearest enclosing project and skip
nested projects. Use `--project NAME` to select a project by name or `--all` to
inspect every project in the Git worktree. Implicit System lookup stops at the
project root. Explicit links can cross project boundaries.
See [project discovery and selection](workspaces.md).

This repository has a root manifest for dotspec and a second manifest for the
[To-do project](../../examples/todo/README.md). Run `dotspec check --all` to check both.

## File formats

Look up the recognized sections for each of the seven [specification formats](formats/README.md).
