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
- `show` reads the nearest enclosing system. Give `show` its source path to select it directly. It displays the name, optional description, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the project boundary, the Git repository root
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

Each project directory can contain one system declaration, directly or in its immediate `specs/` folder:

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
headings. Declare each system quality attribute in its own quality attribute
specification. The optional empty `Quality attributes` heading makes an empty
collection explicit.
There is no requirement inheritance behavior.

### Write the system overview

Use this reading order for the sections that apply. Only Purpose and Scope are required.

1. **Description, Purpose, Goals.** Introduce the system in a sentence. In Purpose, describe the situation that motivates it and the benefit it should provide. Include business background here when it helps explain that benefit. Give a goal its own named entry when readers need to discuss the outcome or trace what contributes to it.
2. **Stakeholders, Scope, Out of scope.** Identify the people and organizations with an interest in the system. Describe the work, users and situations it serves, then name exclusions that clarify its boundary.
3. **Environment, Assumptions, Constraints.** Describe the conditions the system must accommodate, the premises the specification relies on, and any externally imposed limits. For the personal task list, stopping the service between sessions is an operating condition. Returning to the same database file is an assumption behind recovering the list.
4. **System declarations.** Describe the relevant components, interaction partners, interfaces and concept uses. Use the glossary for terms with a project-specific meaning.
5. **Project references.** Link to a plan, responsibility, risk or decision that affects delivery of this system. Say what it affects and link to the record the team maintains. For example, “The release plan names the reviewer for the data migration.” Keep setup and walkthrough instructions in the project README.

See the [personal task list system spec](../../../examples/todo/specs/personal-task-list.system.spec.md) for a worked example. This ordering adapts Meyer's distinction between goals, environment, system and project to a single reader-facing overview. [Meyer, *Handbook of Requirements and Business Analysis*, chapter 3](https://se.inf.ethz.ch/requirements/plan.pdf)

### Include a glossary

Use `## Glossary` with a level-three heading for each term defined in the system
spec. To include terms from another Markdown file, declare its path relative to
the system spec in frontmatter. For example, dotspec uses this declaration.

```yaml
---
type: system-spec
glossary_sources:
  - ../docs/reference/model/glossary.md
---
```

A term needs only a heading and a definition. In an imported vocabulary, use
one H1 for the vocabulary title and one H2 per term. Several systems can
include the same file. A generator can produce it from a model, a catalog, or
another authoritative source.

Use a **topic** when readers need a way to find related entries. A term can
appear under several topics and still have one definition. Use a **context**
when a word has a different meaning for a different group. For example, Billing
can define *Account* as a ledger and Identity can define *Account* as a login.
Set the context once in each vocabulary’s optional frontmatter.

```markdown
---
context: Billing
---
# Billing vocabulary

## Account

A ledger of charges and payments for a customer.

### Topics

- Customers
- Payments

### Relationships

- Receives charges from an [Invoice](#invoice).

### Example

One customer can pay two invoices from the same account.

## Invoice

A request for payment for goods or services.
```

| Optional term subsection | Use |
| --- | --- |
| Alternative names | List other names for the same meaning. |
| Topics | List one plain-text topic name per Markdown list item. |
| Relationships | Write a sentence explaining how this term relates to another. Link the other term by its source heading, including the file path for a different vocabulary. Its entry links back to this passage when both terms appear in the system glossary. |
| Related terms | Point readers to neighboring terms or explain how two contexts use a word. |

A usage example belongs in the term definition. A local heading can introduce it.

Inline terms use the same subsections one heading level deeper. They share
the system’s default context with imported vocabularies that omit `context`.
`dotspec model review` flags preferred names and aliases shared by different
terms within the same system and context. Topic membership has no effect on review.

Dotspec’s `mise run model-generate` takes definitions from class descriptions,
topics from metamodel subsets, and the Feature relationships from its allowed
group members. Run `dotspec spec build` to combine that vocabulary with inline
terms. [Generated vocabulary](../model/glossary.md) · [System glossary](../../../specification/glossary.md)

### Match the filename to the heading

Use the visible H1 as the subject name. Lowercase its letters, keep letters and
digits, and replace each run of spaces or punctuation between them with one
hyphen. For example, `# Spec parser` requires `spec-parser.system.spec.md` and
`type: system-spec`. An optional frontmatter `id` supplies a stable link anchor.
The directory can keep its codebase name. When a subject is renamed, update its
heading, filename and incoming links together. `dotspec check` reports the
expected filename when the two names disagree.

## Component spec format

Place a component declaration in the directory it describes or its immediate `specs/` folder. A component uses the same filename
rule as a system, with the suffix `.component.spec.md`. For example,
`app/specs/task-service.component.spec.md` contains:

```markdown
---
type: component-spec
---

# Task service

## Description

Records task changes and acknowledges each committed update.

## Part of

[Personal task list](../../specs/personal-task-list.system.spec.md)
```

From the To-do project root, the equivalent initializer for an empty component
directory is:

```sh
dotspec component init "Task service" --path app \
  --part-of specs/personal-task-list.system.spec.md
```

A component requires one `Part of` link to a system or another component in
the project. A subsystem links to its parent system with the same heading.
The parent’s `components.md` and `systems.md` indexes list its children from those links.

Use the same collections under either subject's `specs/` directory.
For example, Task service owns
`specs/requirements/acknowledge-after-commit.requirement.spec.md` under its
component directory. The nearest subject supplies the local owner. Following a
component's parent links identifies the default system for quality classification.
A requirement can explicitly select a different specification system and can
name any applicable targets in `Applies to`.

Write each feature and capability in its own typed file with an H1 and a
level-two `Description`. For an inline concept use, put `Description` one level
below its declaration. The report opens with the description and reuses the
first paragraph in the parent’s collection index.

See [structure and feature authoring](../../handbook/relationships.md#describe-systems-components-and-features)
and the [To-do system overview](../../../examples/todo/specification/personal-task-list.system.spec-report.md).
Use `dotspec spec build` in an [opted-in project](../generated-specification.md)
to generate its overview and supporting pages. To inspect system structure and
applicable requirements as a standalone Markdown report, run `dotspec model overview`.

Source paths in command output are relative to the current directory.

## System quality attributes

Place each quality attribute in a named source owned by its system. The filename uses the name from the H1. For example,
`specs/quality/reliability/durability/durability.quality-attribute.spec.md` contains:

```markdown
---
type: quality-attribute-spec
id: durability
subattribute_of: [reliability]
---

# Durability

## Definition

Preservation of acknowledged work between local sessions.
```

The ID defaults to the name slug. Supply an explicit ID to keep references stable
when renaming the attribute. Each ID is unique within the system. The
`subattribute_of` list names every parent. All parents belong to the same system,
and the graph must be acyclic.

A requirement names each applicable quality in `quality_attributes`.
Component-owned requirements resolve these IDs through their governing system.
Subsystems have independent vocabularies. One possible source layout is:


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

The system’s [quality attributes index](../../../examples/todo/specification/quality-attributes.md) links each attribute’s name and short definition to its
individual report. Readers can follow parent and child links or browse all
requirements classified by that attribute. For a complete example, read the
[Durability source](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md)
and its [attribute report](../../../examples/todo/specification/quality/durability/durability.quality-attribute.spec-report.md).
