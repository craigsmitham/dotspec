# Specification layout

## Project layout

A directory can declare one system or component. Place the declaration directly
in that directory or directly in its optional `specs/` folder. For example,
`app/task-service.component.spec.md` and `app/specs/task-service.component.spec.md`
both claim `app/`. Creating both is a conflict. The filename follows the
[subject naming rule](formats/system.md#match-the-filename-to-the-heading).

Each other directory can claim its own subject. For example,
`app/specs/storage/storage.component.spec.md` claims `app/specs/storage/`.
A component declares its architectural parent with its `Part of` link. An independent project or repository rooted at a directory named `specs/` owns its own declarations.

The nearest enclosing declaration owns a source. Lookup stops at the project
boundary. A requirement with no enclosing declaration still appears in discovery.
Moving a requirement within its owner's directory tree preserves its ownership
until it crosses into another subject. Update relative Markdown links when moving files.

```text
specs/
  personal-task-list.system.spec.md
app/
  main.go
  specs/
    task-service.component.spec.md
    requirements/
      acknowledge-after-commit.requirement.spec.md
  handlers/
    retry.requirement.spec.md
```

Both requirements here belong to Task service. A requirement's `Applies to`
section identifies what it governs. Its `quality_attributes` field names its
qualities. Attribute definitions declare their parents with `subattribute_of`.

## Creating and finding sources

`init` uses predictable destinations beneath the selected directory. Systems
and components go in `specs/`. Collection defaults are `specs/jobs/`,
`specs/use-cases/`, `specs/concepts/`, `specs/requirements/`,
`specs/features/`, `specs/capabilities/`, `specs/business-rules/` and `specs/quality-models/`. You can move the files
next to code or group them by feature. Keep their typed filename suffixes.

For `show` and focused `check`, a slug searches the selected owner's files
and skips child subjects. With no enclosing subject, it searches the selected
subtree up to child subjects. If several files match, the diagnostic lists each
path. A source path selects one file directly and resolves relative to the
selected directory (the working directory by default for a single source). Use `--path` to select a directory or `--project` to start at
a named project root.

```sh
dotspec requirement show retry --path app
dotspec requirement show handlers/retry.requirement.spec.md --path app
dotspec check app/handlers/retry.requirement.spec.md
```

## Independent projects and monorepos

Put a `dotspec.yaml` containing `name: your-project` in each independent
project's root. Default commands select the nearest enclosing project and skip
nested projects. Use `--project NAME` to select a project by name or `--all` to
inspect every project in the Git worktree. Implicit System lookup stops at the
project root. Explicit links can cross project boundaries.
See [project discovery and selection](workspaces.md).

This repository has a root manifest for dotspec and a second manifest for the
[To-do project](../../../examples/todo/README.md). Run `dotspec check --all` to check both.

## File formats

Look up the recognized sections for each of the [specification formats](formats/README.md).
