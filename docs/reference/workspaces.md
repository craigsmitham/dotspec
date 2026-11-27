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

To start a project, run `dotspec init "Billing" --name billing` in its directory. Init writes the manifest and the project's root system specification, and refuses to replace either file. Without `--name`, the project name is the system name's slug. When the directory lies inside another project, init names that project, whose scans now skip the new directory.

The name identifies the project in commands and reports. It starts with a letter
or digit and can contain letters, digits, dots, underscores and hyphens. Moving a
directory carries its project name and boundary with it. Use one YAML document
with a required `name` field and an optional `build.enabled` setting. The generated specification is enabled by default. Any other field is a configuration error. Workspace-wide selection requires unique names.

A spec belongs to its nearest enclosing project. A parent scan skips nested
project directories. Ordinary subdirectories can contain colocated specs, and
a project can have zero, one or several System specifications.

For a workspace containing only child projects, leave the repository root
without a manifest. At that root, select a project with `-C <dir>` or check every
project with `dotspec check --all`. Repositories with no project manifests still
support ordinary path-based commands.

## Select what to inspect

Run the following commands in the example workspace root. Substitute your project directories for `examples/reference-app` and `apps/billing`. For the narrower check, use an existing directory inside the selected project.

```sh
dotspec check                                   # current project
dotspec -C examples/reference-app check
dotspec check internal/parser                   # one subtree of the current project
dotspec -C apps/billing snapshot export
dotspec check --all                             # every discovered project separately
```

`-C <dir>`, or `--chdir <dir>`, runs dotspec as if it had started in that
directory, like `git -C`. Every relative path, whether an operand such as the
`<path>` of `check` or a flag value such as `--part-of`, then resolves against
it. dotspec finds the nearest ancestor `dotspec.yaml`, stopping at the Git/worktree
boundary (or filesystem root outside Git). Reads such as `check`, `sources`,
lists, `snapshot export` and `build` cover the enclosing project from its root,
wherever in the project you run them. `check`, `sources` and `snapshot export`
take an optional `<path>`, a directory that narrows the scan to one subtree;
`check` also accepts a single specification file, which it checks within that
file's own project. `init` and `new` write below the working directory, and
`new quality-attribute` adopts into the nearest system that encloses it.

For `check --all`, dotspec searches the entire Git worktree from any
subdirectory. Tracked and non-ignored manifests participate. Outside Git, the
search starts at the outermost enclosing project or, if there is none, at the
working directory. Discovery skips symlinks and nested repositories. A project manifest or repository root creates a boundary even in a directory named `specs/`.
A local manifest still bounds direct lookup and parent scans even when Git ignores
it. Track it to include that project in the workspace check.

`--all` is available only on `check`, the one workspace-wide inspection, and
cannot be combined with a `<path>`. Every other command, including
initialization and other mutations, operates on one project.

Workspace text output separates projects with `# Project: <name>` headings.
JSON output has `schema_version: 1`, the searched directory in `workspace`, and
a name-sorted `projects` array. Each entry contains `name`, `root`, `exit_code`,
and that project's check report in `result`. Every project is checked. The
aggregate exit code is the highest project result: 0 success, 1 specification
findings, 2 invocation or operational failure.

Check that the workspace report names each project you expected and that findings appear under their owning project. If a project is missing, check its manifest name and whether Git ignores its directory. When two manifests use the same name, give the projects distinct names before running the workspace check again. If a project’s generated specification is missing or out of date, run `dotspec -C apps/billing build` with that project’s directory, review the result, and repeat the check. [Generate a requirements specification](generated-specification.md).

To save a project's snapshot, export it from within the project, or name the
project's directory with `-C`. Each export covers one project with its own
acquisition scope and identities.

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
