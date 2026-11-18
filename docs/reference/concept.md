# Concepts

## Concept commands

```sh
dotspec concept init job-map
dotspec concept show job-map
dotspec concept list
dotspec concept check job-map --json
```

Use concepts to describe a coherent piece of product behavior. Each concept
explains its purpose, abstract state, available actions, and an operational
principle: a characteristic story showing how the behavior fulfills the purpose.
Concept descriptions are optional; a project does not need one for every noun,
specification type, or implementation component.

Every command accepts `--path DIR`, selecting an existing project directory,
and `--json` or `--format json`.

- `init <slug>` creates `specs/concepts/<slug>.concept.spec.md` from an editable
  template. It creates missing metadata folders but never overwrites a file or
  writes through a symlinked metadata directory.
- `show <slug>` selects the concept in that project's `specs/concepts/` folder,
  without searching ancestor projects for the slug. Text includes the complete
  authored Markdown, associated system, and source path. JSON also exposes
  `purpose`, `state`, `actions`, and `operational_principles` as Markdown strings.
- `list` discovers concepts in the selected subtree, grouped by enclosing system.
  Groups are ordered by system name and source path, with unassociated concepts
  last; concepts within each group are ordered by source path. An empty list
  succeeds.
- `check <slug>` checks the selected concept and its enclosing system, reporting
  structural errors and missing-content advice through the shared check report.
  The project-wide `dotspec check` includes concepts as well.

Slugs use lowercase ASCII letters or digits separated by single hyphens.
Different project directories can reuse a slug. The nearest enclosing system
is an association for organization, not a service boundary, inheritance rule,
or statement of exclusive implementation ownership. A concept without an
associated system is valid.

## Concept spec format

Only `type: concept` frontmatter and one nonempty document-level H1 are required.
An optional `id` follows the same document-local identifier rules as jobs and use
cases: begin with an ASCII letter or digit, followed by letters, digits, hyphens,
or underscores. Dots are not allowed. IDs are not resolved across files, and
repository-wide uniqueness is not checked. Unknown frontmatter fields are errors.

```markdown
---
type: concept
id: task
---

# Task

## Purpose

Track a piece of work through completion.

## State

Each task has an identity, text, and an active or completed status.

## Actions

### Create

Given text, create a new active task.

### Complete

For an existing task, set its status to completed.
Completing an already completed task leaves it unchanged.

### Reopen

For an existing task, set its status to active.

## Operational principles

Create a task to remember work. Complete it when the work is finished.
If more work is needed, reopen the same task and complete it again.
```

The four basic section names are case-sensitive H2 headings. Use
`Operational principles` once, with optional named subsections for individual stories. Optional
`Type parameters` and `Introduced types` sections declare named types. Their order
is flexible, and each may appear at most once. Missing or empty sections produce
advice, so unfinished concepts remain readable. Additional sections and nested
headings remain authored Markdown. A nested heading named `State`, for example,
belongs to its containing section and does not replace the concept's H2 `State`.

State describes the information needed to explain behavior; it need not match a
database schema. Actions can use prose, tables, code blocks, or links to governing
artifacts. Action subheadings have ordinary Markdown anchors. Authoring commands
retain their Markdown; model export also creates `ConceptAction` records. Describe
inputs, outputs, conditions, rejection behavior, and effects
where relevant. Distinguish caller assumptions from enforced guards explicitly;
no heading or notation automatically selects contract semantics.

The operational principle explains how the design serves its purpose. Checking
its presence does not establish complete behavior, preserved invariants, or
requirement satisfaction. Checks resolve same-document heading and ID references;
relative file links are validated by check commands, while remote links remain
navigation. A link alone
does not designate an authoritative definition or associate a test result.

Use cases can reference concept actions. Product-level coordination belongs in
explicit `Concept uses` and `Synchronizations` sections of a system specification.
A named use references a reusable concept definition and can bind its type
parameters to primitives or types introduced by another named use. Synchronization
actions identify both the use and the action. `Dependencies` on a use describe its
application-specific role; they do not mean that concepts call one another.

Model export resolves these declarations within the acquired reference scope and
reports mistyped or incorrectly scoped references. Unresolved drafts retain their
source. No field or link implies atomicity, execution, full type-expression checking,
or verified behavior. See [Develop concepts](../develop-concepts.md) for the
complete Markdown conventions and [the generated example](concept-design-example.md)
for the checklist definition, effects table and its use in the task list.

```sh
dotspec model docs --path examples/todo --relative-to docs/reference
```

See [Job Map](../../specs/concepts/job-map.concept.spec.md) and
[Specification Inspection](../../specs/concepts/specification-inspection.concept.spec.md)
for concepts describing dotspec itself. The
[check-before-committing use case](../../specs/use-cases/check-specifications.use-case.spec.md)
shows Specification Inspection in a contributor's workflow.
