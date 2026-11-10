# dotspec

A framework and CLI for managing software product requirements in `*.spec.md` files.

System specs describe a system's name, purpose, and scope. Subsystems use the
same format, placed alongside the code they describe.
Job specs describe jobs to be done (JTBDs): the accomplishment someone seeks,
their context, the steps involved, and how success is judged. Jobs are associated
with the nearest enclosing system.

## Run

With Go 1.27.3 or later installed:

```sh
go run . system show
go run . system list
go run . job list
go run . job show develop-requirements
```

## Build

```sh
go build -o dotspec .
./dotspec --help
```

## System commands

```sh
dotspec system init                     # Create a system spec here
dotspec system show                     # Describe the nearest enclosing system
dotspec system list                     # Discover systems in this subtree
dotspec system show --path internal/cli
dotspec system list --path internal
```

Use `./dotspec` for a local build, or replace `dotspec` with `go run .` when
running from the repository root.

Every command accepts `--path DIR`, defaulting to the current directory. The
directory must already exist. Use `dotspec system <command> --help` for usage.

- `init` creates `system.spec.md` with editable placeholders for the name,
  purpose, and scope. It never overwrites an existing file.
- `show` reads the selected directory's system spec, or the nearest ancestor's
  if there is no local spec. It displays the name, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the Git repository root
  (including a worktree root), or the filesystem root outside Git.
- `list` displays a tree of systems within the selected directory and its
  descendants, ordered by source path. It does not include ancestor systems
  outside that subtree. An empty result is successful and reports that no
  systems were found.

For both `system list` and `job list`, discovery in Git repositories uses Git's
ignore rules, including nested `.gitignore` files and local exclusions.
Tracked files remain visible even if
an ignore rule matches them. Git must be installed for discovery in a repository.
Outside Git, discovery walks the directory tree without applying ignore rules.
Discovery skips `.git`, nested repositories, and symlinked files and directories;
an explicitly selected symlink directory is resolved before the command runs.

Missing specs for `show`, invalid paths, malformed specs, and command errors
produce a diagnostic on standard error and a nonzero exit status. A malformed
discovered spec makes `list` fail rather than silently omit it. Each list command
reads only its own spec type; `job list` also reads the systems associated with
its jobs.

## System spec format

Each directory can contain one system spec, named `system.spec.md`:

```markdown
---
type: system
---

# Spec parser

## Purpose

Make specifications available to dotspec commands as structured information.

## Scope

Reading and interpreting specification files.
```

The frontmatter currently accepts only `type: system`. The Markdown body must
contain exactly one nonempty document-level H1 naming the system, followed by
nonempty `Purpose` and `Scope` H2 sections. Section names are case-sensitive.
Section bodies can contain Markdown, including lists, code blocks, and deeper
headings. Other sections are allowed but are not displayed by `system show`.
There is no requirements section or requirement inheritance behavior.

The nearest ancestor directory containing a system spec defines the parent:

```text
system.spec.md                  # dotspec
internal/parser/system.spec.md  # Spec parser, a subsystem of dotspec
```

Moving a spec can therefore change its parent. Source paths in command output
are relative to the current directory.

## Job commands

```sh
dotspec job init develop-requirements
dotspec job show develop-requirements
dotspec job list
dotspec job step list develop-requirements
dotspec job outcome list develop-requirements
dotspec job outcome list develop-requirements --step identify-stakeholders
dotspec job map develop-requirements
dotspec job map develop-requirements --outcomes
dotspec job map develop-requirements --format mermaid
dotspec job check develop-requirements
```

The `--path DIR` directory must already exist and defaults to the current
working directory. It can appear before or after the slug. For example:

```sh
dotspec job init understand-tokens --path internal/parser
dotspec job show understand-tokens --path internal/parser
dotspec job list --path internal/parser
```

Use `dotspec job <command> --help` for command-specific options.

- `init <slug>` creates `<slug>.job.spec.md` with an editable template. It never
  overwrites an existing file. Edit the file in your text editor.
- `show <slug>` displays all authored Markdown, including additional sections,
  plus the associated system and source path. It selects the job in the chosen
  directory without searching other directories for the same slug.
- `list` discovers jobs within the selected subtree, grouped by associated
  system. Groups are ordered by system name and source path; jobs within each
  group are ordered by source path. Jobs without a system appear last under
  `System: None`. An empty list is successful.
- `step list <slug>` displays steps in document order with their IDs, optional
  stages, and descriptions.
- `outcome list <slug>` displays success criteria under `Overall job outcomes`
  and `Outcomes for step: ...`. `--step ID` selects just that step's outcomes,
  excluding overall job outcomes. An unknown step is an error; an existing step
  without outcomes is a successful empty result.
- `map <slug>` displays the job and its documented steps in authored order.
  `--format mermaid` emits standalone Mermaid flowchart source for a renderer
  or a fenced `mermaid` block. `--outcomes` includes the whole job's outcomes and
  attaches each step's outcomes to that step, with scope labels in both formats.
  Without steps, the command displays the job, any requested overall outcomes,
  and `No steps defined.`. Arrows show authored step order; they do not require
  that execution sequence. Neither format assesses the map's completeness.
  Loops and branches are not modeled in this first version.
- `check <slug>` reports structural errors and offers advisory writing suggestions.
  Suggestions do not produce a failure exit status. Ordinary reading and mapping
  commands do not emit writing advice.

A file slug contains lowercase ASCII letters or digits, separated by single
hyphens: for example, `develop-requirements`. Slugs cannot contain paths, spaces,
or file extensions. Different directories can reuse a slug. Changing a job's
name does not change its file slug.

Each job is associated with the nearest `system.spec.md` in its directory or an
ancestor directory, even if that system is outside the subtree selected for
`job list`. Lookup stops at the Git repository or worktree root, or the filesystem
root outside Git. A job without an enclosing system is valid and displays
`System: None`. A malformed enclosing system is an error. System association does
not mean that the software performs every step of the customer's job.

## Job spec format

Jobs use one Markdown convention with no profiles or method selection. Only
`type: job` frontmatter and one nonempty document-level H1 are required. An
optional `id` in frontmatter supplies an explicit job identity; the file slug
remains the CLI selector. Other frontmatter keys are currently rejected.

For example, a draft with a partial map:

```markdown
---
type: job
id: develop-requirements
---

# Develop Requirements

## Job

Define an agreed set of requirements for a proposed software change.

## Performer

The person responsible for defining the change, working with stakeholders.

## Situation

I’m defining a new product or changing an existing capability.

## Motivation

I want to turn stakeholder needs and product intent into clear requirements.

## Intended progress

The team can agree on what to build and how to determine whether it meets
those needs.

## Outcomes

- **OUT-001:** Minimize the likelihood that an important stakeholder need is
  discovered only after implementation begins.

## Steps

### Identify affected stakeholders

**ID:** identify-stakeholders
**Stage:** locate

Identify people whose needs or responsibilities may be affected.

#### Outcomes

- **OUT-002:** Minimize the time needed to identify affected stakeholders.

### Locate existing commitments

**Stage:** locate

Find the commitments that may constrain the proposed change.

## Evidence

This partial map and its outcomes are hypotheses awaiting research.
```

The conventional H2 sections are `Job`, `Performer`, `Situation`, `Motivation`,
`Intended progress`, `Outcomes`, `Steps`, `Progress cases`, and `Evidence`. All are
optional, may be empty in a draft, and can be developed incrementally. Their
names are case-sensitive, and duplicate conventional sections are errors.
Additional sections and ordinary Markdown are retained. Existing three-section
job stories remain valid.

`Intended progress` is the job story's narrative result. The previous heading,
`Desired outcome`, remains a supported alias. Using both headings, or repeating
either one, is a duplicate-section error even if a section is empty. `show`
preserves the authored heading and Markdown; existing documents need no migration.

`Outcomes` contains success criteria. A job's H2 `Outcomes` section describes
success for the overall job; each step's H4 `Outcomes` section describes success
at that step. These are separate collections: overall outcomes are not inherited
by steps, and step success does not establish overall job success. Either
collection can be absent or empty, meaning no outcomes are recorded at that scope.
Narrative progress does not create an outcome entry.

Each direct list item defines an outcome; use a bold ID followed by a colon, as in
`- **OUT-001:** Minimize ...`, when the outcome needs to be referenced.
IDs are optional. Continuation lines and paragraphs are supported.
Unidentified and unfinished outcome statements remain readable; `check` offers
wording advice for empty statements or statements without recognizable direction
and measurement terms. Nested bullets are preserved in `show` as additional prose
and do not define separate outcomes.

A job can have overall outcomes without any steps or IDs, for example:

```markdown
---
type: job
---

# Assess a proposed change

## Job

Assess a proposed product change from receiving the proposal to deciding
whether to accept it.

## Performer

A reviewer responsible for the affected product commitments.

## Outcomes

- Minimize the likelihood of accepting a change that violates an applicable commitment.

## Evidence

This job and outcome are hypotheses awaiting research. No job map is recorded.
```

Under `Steps`, each H3 introduces a step. Its `**ID:**` and optional `**Stage:**`
fields belong in the introductory prose. A step's H4 `Outcomes` section contains
only that step's success criteria. Headings in code fences and blockquotes do
not define steps or outcomes. Progress cases and evidence are currently retained
as Markdown; their subsections are not interpreted as job-map structure.

Step order in the document defines the displayed map order, not a required
execution sequence. Mark partial maps in the introductory prose under `Steps`
or in `Evidence`; these notes remain authored Markdown available through `show`.
IDs remain independent of names and positions. Explicit IDs use ASCII letters
or digits followed by letters, digits, hyphens, or underscores, and are
case-sensitive. Job, step, and outcome
IDs share a namespace inside one job document. Duplicate explicit IDs are errors;
IDs may be reused in other job documents. Add an ID when a step or outcome needs
stable references, or when selecting a step with `--step`. Missing IDs are valid;
if a selection fails and the job contains unidentified steps, the diagnostic
explains how to assign one.

`Stage` is an optional label. The eight ODI stages are useful conventions:
`define`, `locate`, `prepare`, `confirm`, `execute`, `monitor`, `modify`, and
`conclude`. Custom stage labels are accepted without configuration. Stages can
repeat or be omitted; they are elicitation aids rather than required slots.
The partial example above includes two `locate` steps, one with no recorded
outcomes or ID.

Local Markdown links such as `[the criterion](#OUT-002)` are checked against
explicit IDs and document heading anchors. An unresolved local fragment is an
error. Explicit ID links are semantic references understood by dotspec; an
ordinary Markdown renderer may require an HTML anchor to make those IDs clickable.
External URLs and links to other files are retained without validation in this
version. Cross-file identity resolution is not implemented yet.

### Checks and incomplete work

Malformed metadata, duplicate IDs or conventional sections, invalid step heading
structure, duplicate step metadata, and unresolved local references fail with a
nonzero exit status. Missing sections and unfinished content remain readable.

Writing checks are simple, explained heuristics: they suggest a job statement
and performer when absent; suggest clarifying the execution boundary when a
statement has no marker such as “for,” “when,” or “during”; and look for an
improvement direction and a measurement word in authored outcome statements.
Missing IDs, maps, stages, and outcome collections do not produce completion
suggestions. These checks are advisory, not proof of clarity,
completeness, measurability, or research validation. They do not assess whether a
custom stage is correct or make missing evidence an error.

The repository's [Develop Requirements](develop-requirements.job.spec.md) and
[Manage Requirements](manage-requirements.job.spec.md) examples each contain a
bounded job statement, performer, proposed steps, and outcomes. Both are marked
as hypotheses pending research. They cover collaboration and judgment outside
the CLI as well as work dotspec may support.

## Development

```sh
go test ./...
go vet ./...
```

Tests that exercise Git ignore rules and repository boundaries require Git.
