# Specification layout

## Project layout

Each part of the project can have a `specs/` directory. It describes its parent
project directory; folders inside `specs/` organize specifications and do not create
subsystems or new scopes.

```text
specs/
  system.spec.md
  jobs/
    develop-requirements.job.spec.md
    manage-requirements.job.spec.md
  use-cases/
    check-specifications.use-case.spec.md
  concepts/
    job-map.concept.spec.md
    specification-inspection.concept.spec.md
  requirements/
    reject-existing-target.requirement.spec.md
    continue-after-invalid-spec.requirement.spec.md
    advice-permits-success.requirement.spec.md
    diagnostic-severity.requirement.spec.md
    repeatable-inspection.requirement.spec.md
    inspection-preserves-source.requirement.spec.md
internal/
  parser/
    parser.go
    specs/
      system.spec.md
      jobs/
        interpret-specification.job.spec.md
```

The layout is required:

- The system spec is `specs/system.spec.md`, with one system per project directory.
- Jobs are `specs/jobs/<slug>.job.spec.md`, directly inside the plural `jobs/` folder.
- Use cases are `specs/use-cases/<slug>.use-case.spec.md`, directly inside the
  plural `use-cases/` folder.
- Concepts are `specs/concepts/<slug>.concept.spec.md`, directly inside the
  plural `concepts/` folder.
- Requirements are `specs/requirements/<slug>.requirement.spec.md`,
  directly inside the plural `requirements/` folder.
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
existing `system.spec.md` to `specs/system.spec.md`, and move each adjacent
`<slug>.job.spec.md` to `specs/jobs/<slug>.job.spec.md`. Discovery rejects misplaced
files of the requested type with an expected-location diagnostic. Local lookup
and initialization also reject misplaced specs of that type directly in the
selected project directory or anywhere within its `specs/` tree. Having both
old and new locations is an error; there is no location precedence or automatic
migration. Revise any relative file links when moving documents.

## Independent projects and monorepos

Register specification-project roots in `dotspec.yaml` to keep unrelated products
and examples separate. Default commands select the containing project; use
`--project NAME` to select another and `--all` for separate workspace-wide results.
Implicit System context stops at project boundaries; explicit cross-project links
remain available. See [workspace configuration and selection](../workspaces.md).

This repository registers dotspec and its example collections separately. Use
`dotspec check --all` to check both the product and examples. The LinkML definition
is `model/metamodel.yaml`; it is not a workspace configuration.

## File formats

[Systems](system.md) · [Jobs](job.md) · [Use cases](use-case.md) ·
[Concepts](concept.md) · [Requirements](requirement.md) ·
[Business rules](../business-rules.md) · [Quality models](../quality-models.md).
