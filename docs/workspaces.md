# Separate specification projects in one repository

Use a workspace when a repository contains independent products, a reference app,
or example collections that should not enter the product's own specification
scan. A **specification project** is a tooling partition. A **System** is a domain
subject described by specifications. Project membership creates no `part_of`,
implementation, ownership or requirement-inheritance relationship between Systems.

## Declare the boundaries

Create `dotspec.yaml` at the workspace root:

```yaml
projects:
  product: .
  reference-app: examples/reference-app
  billing: apps/billing
```

The directories must already exist. Names start with a letter or digit and contain
letters, digits, dots, underscores or hyphens. Paths are relative to the manifest;
`.` registers the root project. Roots must be distinct real directories inside the
workspace, outside `specs/` metadata trees and nested Git repositories. Unknown
configuration fields, duplicate names/roots, symlinks and missing roots are errors.

The most specific containing root determines membership. Thus `product` excludes
both registered child projects. Ordinary subdirectories within a project remain
available for colocated specs and related subsystem definitions. Each project may
contain several System specifications. It need not have a System yet.

A workspace-only repository can omit the `.` project. Commands at that root then
require `--project NAME` or `--all`; dotspec does not invent a root product.
Without a manifest, the existing repository/filesystem discovery behavior remains.

The workspace manifest is separate from dotspec's canonical LinkML
[metamodel](../model/metamodel.yaml).

## Select what to inspect

```sh
dotspec check                           # current registered project
dotspec check --project reference-app
dotspec check --project product --path internal/parser
dotspec model overview --project billing
dotspec check --all                     # every registered project separately
dotspec model export --all > workspace.json
dotspec model validate workspace.json
dotspec model schema --workspace > workspace.schema.json
```

Dotspec finds the nearest ancestor `dotspec.yaml`, stopping at the Git/worktree
boundary (or filesystem root outside Git). With no explicit selector, the
containing project is selected; running from one of its subdirectories still
selects that project root. Use `--path` to narrow the scan. On its own, `--path` is
relative to the working directory. With `--project`, it is relative to that
project's root and cannot select another project. Absolute paths within that
project also work.

All path-based commands accept `--project`, including initialization and focused
checks. `--all` is available on collection checks, source inventories, lists and
model exports/reading views. It cannot be combined with `--project`, `--path`, or
an individual specification slug. Initialization and other mutations operate on
one selected project.

Workspace text/Markdown output separates projects with headings. JSON output has
`schema_version: 1`, `workspace`, and a sorted `projects` array. Each entry contains
`name`, `root`, `exit_code`, and its command `result`; an operational failure may
instead supply `error` without a result. All projects are attempted. The aggregate
exit code is the highest result: 0 success, 1 specification findings, 2 invocation
or operational failure.

Workspace model exports contain independent snapshots rather than one combined
System. `model validate` checks each snapshot and its envelope identity offline;
a missing snapshot is incomplete and fails validation. The workspace schema
covers complete export envelopes, including snapshots that faithfully record
specification findings.

## Keep context explicit

Implicit System lookup stops at the selected project boundary. A project with no
System does not borrow its enclosing directory's System or quality attributes.
Directory-derived `enclosing_system` also stops there. Inside a project, normal
nearest-System lookup and explicit structural relationships continue to work.

Explicit references can cross project boundaries under the existing repository,
symlink and direct-reference rules. For example, a requirement can select another
project's System explicitly, or an Interface can cite a shared native contract.
These references acquire context, not another recursive project scan.

Snapshot `scope.project` identifies the selected project. Acquired source records
carry their containing `project` when registered, and the consulted manifest is
retained as `workspace_configuration` context. Source inventories also identify
registered project membership. Generated reading views name the selected project
and identify acquired domain context from other projects. These are recorded
acquisition facts, not claims that the workspace or requirements are complete.

## This repository

The root project is `dotspec`. Existing example collections are registered as
separate projects, such as `system-structure` and `requirements-review`. Default
commands inspect dotspec; `mise run spec-check` and the full check gate inspect
all projects. Add a future runnable reference app by creating its directory and
specs and registering a new project root. No fake subsystem or Git submodule is
needed for isolation.
