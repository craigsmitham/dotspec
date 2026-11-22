# Separate specification projects in one repository

The To-do example lives inside the dotspec checkout so a contributor can run it alongside the CLI. If directory nesting alone put it in the product specification, a reader could mistake keeping a personal task list for a responsibility of dotspec. Each has its own project manifest, so checking the CLI from its project root leaves the example to its own check.

A **specification project** determines which sources a scan discovers together. A **System** is a subject the authors describe. They can declare that one system forms part of another or link to a contract in a neighboring project. Those decisions require evidence about the products, beyond where their files live. [Project isolation requirement](../../specs/requirements/specification-project-isolation.requirement.spec.md) · [Describe a system and its boundary](../cli-handbook/guides/describe-a-system.md).

## Declare the boundaries

Place a `dotspec.yaml` in each project's root directory. For a product with a
reference app and a separate billing project, the layout is:

```text
dotspec.yaml                         # name: product
examples/reference-app/dotspec.yaml  # name: reference-app
apps/billing/dotspec.yaml            # name: billing
```

For example, the entire billing manifest is:

```yaml
name: billing
```

The name identifies the project in commands and reports. It starts with a letter
or digit and can contain letters, digits, dots, underscores and hyphens. Moving a
directory carries its project name and boundary with it. Use one YAML document
with a required `name` field and optional `report` section. Its boolean `enabled` field defaults to `true`. Workspace-wide selection requires unique names.

A spec belongs to its nearest enclosing project. A parent scan skips nested
project directories. Ordinary subdirectories can contain colocated specs, and
a project can have zero, one or several System specifications.

For a workspace containing only child projects, leave the repository root
without a manifest. At that root, use `--project NAME` or `--all`. Repositories
with no project manifests still support ordinary path-based commands.

## Select what to inspect

```sh
dotspec check                           # current project
dotspec check --project reference-app
dotspec check --project product --path internal/parser
dotspec model overview --project billing
dotspec check --all                     # every discovered project separately
dotspec model export --all > workspace.json
dotspec model validate workspace.json
dotspec model schema --workspace > workspace.schema.json
```

dotspec finds the nearest ancestor `dotspec.yaml`, stopping at the Git/worktree
boundary (or filesystem root outside Git). With no explicit selector, project inventories start at the containing project root. Initialization, show, and focused checks start at the working directory. Use `--path` to narrow the scan. On its own, `--path` is
relative to the working directory. With `--project`, it is relative to that
project's root and cannot select another project. Absolute paths within that
project also work.

For `--project` and `--all`, dotspec searches the entire Git worktree from any
subdirectory. Tracked and non-ignored manifests participate. Outside Git, the
search starts at the outermost enclosing project or, if there is none, at the
working directory. Discovery skips symlinks and nested repositories. A project manifest or repository root creates a boundary even in a directory named `specs/`.
A local manifest still bounds direct lookup and parent scans even when Git ignores
it. Track it to include that project in worktree-wide commands.

All path-based commands accept `--project`, including initialization and focused
checks. `--all` is available on collection checks, source inventories, lists and
model exports/reading views. It cannot be combined with `--project`, `--path`, or
an individual specification slug. Initialization and other mutations operate on
one selected project.

Workspace text/Markdown output separates projects with headings. JSON output has
`schema_version: 1`, the searched directory in `workspace`, and a name-sorted
`projects` array. Each entry contains
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

Snapshot `scope.project` records the selected name and root. Acquired files carry
their containing project when one exists. The selected project's `dotspec.yaml`
is retained as `project_configuration` context and appears in the source inventory.
Reading views name the selected project and identify acquired domain context
from other projects.

## This repository

In the dotspec checkout, the root manifest declares `dotspec`. The To-do example has
its own manifest. Both projects use the default generated reports. From the root,
`dotspec check` checks the product and its generated spec. Use
`dotspec check --all` to include To-do. To add a reference app, place a
`dotspec.yaml` beside its code and specs.

For generated Markdown and whole-project checks, see [Generate a requirements specification](generated-specification.md).
