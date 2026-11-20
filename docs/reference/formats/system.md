# Systems

## System commands

```sh
dotspec system init "Spec parser"       # Create specs/spec-parser.system.spec.md
dotspec system show                     # Describe the nearest enclosing system
dotspec system list                     # Discover systems in this subtree
dotspec system show --path internal/cli
dotspec system list --path internal
```

Use `./dotspec` for a local build, or replace `dotspec` with `go run .` when
running from the repository root.

Every system command accepts `--path DIR`, defaulting to the current project directory.
The project directory must already exist. Use `dotspec system <command> --help`
for usage.

- `init` creates `specs/<name>.system.spec.md` using the supplied display name and editable
  sections for description, purpose, and scope. With no argument it uses “New system”. The
  [existing-target requirement](../../../specs/requirements/reject-existing-target.requirement.spec.md)
  governs initialization at an existing path.
- `show` reads the selected project directory's `specs/<name>.system.spec.md`, or the
  nearest ancestor's if there is no local spec. It displays the name, optional description, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the Git repository root
  (including a worktree root), or the filesystem root outside Git.
- `list` displays a tree of systems within the selected directory and its
  descendants, ordered by source path. It does not include ancestor systems
  outside that subtree. An empty result is successful and reports that no
  systems were found.

For `system list`, `job list`, `use-case list`, `concept list`, and
`requirement list`, discovery in Git
repositories uses Git's ignore rules, including nested `.gitignore` files and
local exclusions.
Tracked files remain visible even if
an ignore rule matches them. Git must be installed for discovery in a repository.
Outside Git, discovery walks the directory tree without applying ignore rules.
Discovery skips `.git`, nested repositories, and symlinked files and directories;
an explicitly selected symlink directory is resolved before the command runs.

Missing specs for `show`, invalid paths, malformed specs, and command errors
produce a diagnostic on standard error and a nonzero exit status. A malformed
discovered spec makes `list` fail rather than silently omit it. Each list command
reads only its own spec type; all non-system list commands also
read the systems associated with their specifications.

## System spec format

Each project directory can contain one system spec at `specs/<name>.system.spec.md`:

```markdown
---
type: system-spec
---

# Spec parser

## Description

A Markdown reader that extracts structured requirements from specification files.

## Purpose

Make specifications available to dotspec commands as structured information.

## Scope

Reading and interpreting specification files.
```

The frontmatter accepts `type: system-spec` and an optional `id` for a stable authored identity. The Markdown body must
contain exactly one nonempty document-level H1 naming the system, followed by
nonempty `Purpose` and `Scope` H2 sections. Section names are case-sensitive.
Section bodies can contain Markdown, including lists, code blocks, and deeper
headings. The optional `Quality attributes` section is also displayed by
`system show`
and interpreted as described below. Other sections remain authored Markdown.
There is no requirement inheritance behavior.

Use an optional `## Description` to introduce the system in one short paragraph.
Write it so that someone seeing only the name and description can tell what the
system is. Use `Purpose` for the reason it exists and `Scope` for its boundaries.

### Match the filename to the heading

Use the visible H1 as the subject name. Lowercase its letters, keep letters and
digits, and replace each run of spaces or punctuation between them with one
hyphen. For example, `# Spec parser` requires `spec-parser.system.spec.md` and
`type: system-spec`. An optional frontmatter `id` supplies a stable link anchor.
The directory can keep its codebase name. When a subject is renamed, update its
heading, filename and incoming links together. `dotspec check` reports the
expected filename when the two names disagree.

## Component spec format

Give each component its own code directory with a `specs/` folder. Put one
named declaration directly under that folder. A component uses the same filename
rule as a system, with the suffix `.component.spec.md`. For example,
`components/task-service/specs/task-service.component.spec.md` contains:

```markdown
---
type: component-spec
---

# Task service

## Description

Records task changes and acknowledges each committed update.

## Part of

[Personal task list](../../../specs/personal-task-list.system.spec.md)
```

From the To-do project root, the equivalent initializer for an empty component
directory is:

```sh
dotspec component init "Task service" --path components/task-service \
  --part-of specs/personal-task-list.system.spec.md
```

A component requires one `Part of` link to a system or another component in
the project. A subsystem links to its parent system with the same heading.
The parent report lists its children from those links.

Use the same collections under either subject's `specs/` directory.
For example, Task service owns
`specs/requirements/acknowledge-after-commit.requirement.spec.md` under its
component directory. The nearest subject supplies the local owner. Following a
component's parent links identifies the default system for quality classification.
A requirement can explicitly select a different specification system and can
name any applicable targets in `Applies to`.

For a feature, capability, or concept use, put `Description` one heading level
below its declaration. The report opens with the description and reuses the
first paragraph in its parent's directory.

See [structure and feature authoring](../../handbook/relationships.md#describe-systems-components-and-features)
and the [To-do system overview](../../../examples/todo/specification/personal-task-list.system.spec-report.md).
Use `dotspec spec build` in an [opted-in project](../generated-specification.md)
to generate its overview and supporting pages. To inspect system structure and
applicable obligations as a standalone Markdown report, run `dotspec model overview`.

Source paths in command output are relative to the current directory.

## System quality attributes

Give each quality attribute a directory under the system’s `specs/quality/`.
The directory and filename use the name from the H1. For example,
`specs/quality/reliability/durability/durability.quality-attribute.spec.md` contains:

```markdown
---
type: quality-attribute-spec
id: durability
---

# Durability

## Definition

Preservation of acknowledged work between local sessions.
```

The ID defaults to the name slug. Supply an explicit ID to keep references stable
when renaming the attribute. Each ID is unique within the system. Nesting the
Durability folder under Reliability makes Reliability a parent. For additional
parents, add a frontmatter list such as `subattribute_of: [recoverability]`. A
flat folder can use that list for every parent. All parents belong to the same
system, and the graph must be acyclic.

Keep functional requirements under the owner’s `specs/requirements/`. They can name relevant attributes in `quality_attributes`. A quality obligation beside an attribute source takes its primary classification from that folder.
Use `quality_attributes` in its frontmatter for any additional IDs. Requirements
in a component's own `specs/requirements/` can use the same IDs through their
governing system. Subsystems have independent vocabularies.

```text
specs/
  quality/
    reliability/
      reliability.quality-attribute.spec.md
      durability/
        durability.quality-attribute.spec.md
        data-loss-tolerance.requirement.spec.md
  requirements/
    restart.requirement.spec.md
  quality-models/
    personal-data.quality-model.spec.md
```

The system overview links each attribute’s name and short definition to its
individual report. Readers can follow parent and child links or browse all
obligations classified by that attribute. For a complete example, read the
[Durability source](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md)
and its [attribute report](../../../examples/todo/specification/quality/durability/durability.quality-attribute.spec-report.md).
