# Use cases

## Use-case commands

```sh
dotspec use-case init check-specifications
dotspec use-case show check-specifications
dotspec use-case list
dotspec use-case check check-specifications --json
dotspec check
```

Each command accepts `--path DIR` (the existing project directory, defaulting to
`.`) and `--json` or `--format json`. Slugs follow the same rules as job slugs.

- `init <slug>` creates `specs/use-cases/<slug>.use-case.spec.md` from an editable
  brief template without overwriting existing content. Use `--template expanded`
  for additional prompts in the same format; neither option sets a maturity state.
- `show <slug-or-source-path>` displays the authored Markdown, associated system and source path.
- `list` discovers use cases in the selected subtree and groups them by enclosing
  system. Groups are ordered by system name and source path, with unassociated
  use cases last; use cases within each group are ordered by source path.
  `--overview` displays actor, goal (falling back to the title), goal level, stated
  boundary, and source. Missing actors, levels, and boundaries are shown as
  unspecified. Multiline prose is collapsed onto one line without interpreting it.
  JSON retains the full grouped records with or without this presentation flag.
- `check <slug>` checks one use case and its enclosing system, reporting structural
  errors and nonfatal missing-content advice. Project-wide `check` includes use
  cases alongside systems and jobs. Ordinary reads do not emit writing advice.

Use cases follow the same discovery, ignore, repository-boundary, and placement
rules as jobs. An enclosing system is used for organization; it does not fill in
or override the use case's stated system boundary. A use case without an enclosing
system is valid and reports `System: None` (`null` in JSON). A malformed enclosing
system is an error.

## Use-case spec format

Only `type: use-case-spec` frontmatter and one nonempty document-level H1 are required.
An optional frontmatter `id` uses the same syntax as job IDs; other frontmatter
keys are rejected. The file slug remains the CLI selector. IDs are local to the
document and are not checked for uniqueness across files.

```markdown
---
type: use-case-spec
id: check-specifications
---

# Check specifications before committing

## Goal

Identify structural problems before committing a change.

## System boundary

The dotspec CLI.

## Primary actor

A contributor preparing a change.

## Goal level

User task (sea level).

## Trigger

The contributor requests a specification check.

## Preconditions

The contributor has selected an existing project directory. Read failures are handled as extensions.

## Success guarantees

The contributor receives a complete report for the selected scope.

## Main success scenario

1. The contributor selects the project subtree to inspect.
2. dotspec inspects its supported specifications.
3. dotspec reports findings and scan completion.

## Extensions

### Unreadable file

#### Condition

A file cannot be read.

#### Handling

Report the failure and mark the scan incomplete.

## Sequence

Representative. The check report contract defines required behavior.

## Supports

[Evaluate requirements](../jobs/develop-requirements.job.spec.md#evaluate-the-requirements).
```

The H2 headings above, plus `Supporting actors`, `Stakeholder interests`, and
`Minimal guarantees`, are conventional sections. They are case-sensitive,
optional, may appear in any order after the title, and may be empty in a draft.
Duplicate conventional sections are errors, even if empty. A conventional section
at document level must use H2. Deeper headings belong to the containing section;
their names do not create additional conventional sections. Code examples and
blockquotes do not define sections. Keep additional prose, examples and local headings within the field they qualify.

JSON exposes conventional section bodies as Markdown strings in `goal`,
`boundary_description`, `primary_actor`, `goal_level`, `trigger`, `preconditions`,
`success_guarantees`, `main_success_scenario`, `extensions`, `sequence`,
`supporting_actors`, `stakeholder_interests`, and `minimal_guarantees`.
Missing bodies are empty strings; `supports` appears only when that section has
content. The object also includes
`name`, `slug`, optional `id`, `markdown`, absolute `path`, `system`, and the title's
source `location`.

A trigger is an initiating event; preconditions describe an already-expected
state. Conditions the system must detect and respond to belong in the interaction
and its extensions rather than being assumed away. Success guarantees state what
holds when the goal is achieved; they can be stated
here or referenced from a governing rule. Minimal guarantees describe what remains
protected even when the goal fails. Supporting actors identify external services
or notification recipients; stakeholder interests can include nonparticipants.
`Sequence` explains whether the narrative order is representative or required;
link to the governing rule when order is required. `Supports` links the jobs,
steps, desired outcomes, or higher-level use cases this use case supports.
These bodies are prose, not
enumerations or executable contracts. Authoring commands keep steps and
extensions as Markdown inside the use case; model export also creates
scenario-step and extension records. Neither makes branch-coverage claims.

`check` suggests a goal, system boundary, primary actor, goal level, and main
success scenario when their conventional sections are missing or empty. Once a
main success scenario is present, it also suggests missing success guarantees,
extensions, and sequence information. An extension can record only a condition
with unresolved handling. These suggestions are possible next writing steps, not
completion gates or a maturity score. The checker does not interpret prose
elsewhere as these fields or judge whether a filled section is adequate. Other
optional sections and IDs produce no missing-content advice.

Same-document links resolve against heading anchors and the optional document
ID. Check commands validate relative file links using the shared
[reference rules](../inspection.md#references); references
such as `Supports` do not establish that an outcome is achieved. Job steps
and use cases can relate many to many. There is no test-result matching or
requirement-satisfaction inference in this version.

The repository's [Check specifications before committing](../../../specs/use-cases/check-specifications.use-case.spec.md)
example connects this interaction to the Develop Requirements job and references
the governing check report contract. The summary use case
[Prepare specifications for review](../../../specs/use-cases/prepare-specifications.use-case.spec.md)
connects it to creating a draft and discovering sources. See
[Develop use cases](../../handbook/use-cases.md) for extension links, behavioral
authority, and delivery slices.
