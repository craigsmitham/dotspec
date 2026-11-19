# Jobs

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

The `--path DIR` project directory must already exist and defaults to the current
working directory. It can appear before or after the slug. For example:

```sh
dotspec job init understand-tokens --path internal/parser
dotspec job show understand-tokens --path internal/parser
dotspec job list --path internal/parser
```

Use `dotspec job <command> --help` for command-specific options.

- `init <slug>` creates `specs/jobs/<slug>.job.spec.md` with an editable template.
  It never overwrites an existing file. Edit the file in your text editor.
- `show <slug>` displays all authored Markdown, including additional sections,
  plus the associated system and source path. It selects the job in the chosen
  project directory's `specs/jobs/` folder without searching other scopes for the
  same slug.
- `list` discovers jobs within the selected subtree, grouped by associated
  system. Groups are ordered by system name and source path; jobs within each
  group are ordered by source path. Jobs without a system appear last under
  `System: None`. An empty list is successful.
- `step list <slug>` displays steps in document order with their IDs, optional
  stages, and descriptions.
- `outcome list <slug>` displays success criteria under `Outcomes declared by job`
  and `Outcomes declared by step: ...`. `--step ID` selects outcomes declared by that step. Use `model needs` to inspect
  explicit applicability across declaration owners. An unknown step is an error; an existing step
  without outcomes is a successful empty result.
- `map <slug>` displays the job and its documented steps in authored order.
  `--format mermaid` emits standalone Mermaid flowchart source for a renderer
  or a fenced `mermaid` block. `--outcomes` includes the whole job's outcomes and
  attaches declared outcomes to their owner, with declaration labels in text and Mermaid.
  This view describes declaration structure; `model needs` shows explicit applicability.
  Without steps, the command displays the job, any requested overall outcomes,
  and `No steps defined.`. Arrows show authored step order; they do not require
  that execution sequence. No format assesses the map's completeness.
  Loops and branches are not modeled in this first version.
- `check <slug>` reports structural errors and advisory writing suggestions on
  stdout. Use `--json` for the same structured report as the project-wide check.
  Suggestions do not produce a failure exit status. Ordinary reading and mapping
  commands do not emit writing advice.

A file slug contains lowercase ASCII letters or digits, separated by single
hyphens: for example, `develop-requirements`. Slugs cannot contain paths, spaces,
or file extensions. Different directories can reuse a slug. Changing a job's
name does not change its file slug.

Each job belongs to the project directory containing its `specs/` tree. It is
associated with that directory's `specs/system.spec.md`, or the nearest ancestor
project directory's system, even if that system is outside the subtree selected for
`job list`. Lookup stops at the Git repository or worktree root, or the filesystem
root outside Git. A job without an enclosing system is valid and displays
`System: None`. A malformed enclosing system is an error. System association does
not mean that the software performs every step of the customer's job.

## Job spec format

Jobs use one Markdown convention with no profiles or method selection. Only
`type: job-spec` frontmatter and one nonempty document-level H1 are required. An
optional `id` in frontmatter supplies an explicit job identity; the file slug
remains the CLI selector. Other frontmatter keys are currently rejected.

For example, a draft with a partial map:

```markdown
---
type: job-spec
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

`Intended progress` describes the narrative result the person seeks. Use that
heading consistently: the old singular `Desired outcome` heading is now ordinary
extension content and does not populate intended progress.

`Outcomes` contains desired results or criteria for judging success. Each outcome
has one declaration owner: the job or a step. Declaration ownership does not
establish applicability or inheritance. A shared outcome can explicitly apply to
several jobs or steps; declare it once and link its subjects with `Applies to`.
Step success does not establish overall job success.

Each direct list item defines an outcome; use a bold ID followed by a colon, as in
`- **OUT-001:** Minimize ...`, when it needs a reference. IDs are optional for list
items. Continuation paragraphs and nested support bullets remain prose.

For a concise map label and explicit applicability, use a named outcome under
`Outcomes`, with an introductory `**ID:**` and full statement:

```markdown
### Timely awareness

**ID:** timely-awareness

Minimize the time before recognizing that the journey plan needs attention.

#### Applies to

[Monitor the journey](#monitor-journey)
[Adapt travel plans](#adapt-plans)
```

Under a step's `#### Outcomes`, use H5 outcome names and H6 `Applies to` headings.
A heading without an explicit ID remains supporting extension content; it does
not silently become another outcome. The statement is separate from the short
heading label. `Supporting content`, `Supports`, and `Derived from` sections can
retain qualifications and provenance. Ordinary links elsewhere remain citations.

Empty statements receive advisory guidance. ODI direction/measurement wording
advice is enabled only by `job check --preset odi`; qualitative emotional and
social outcomes are accepted by the default check.

A job can have overall outcomes without any steps or IDs, for example:

```markdown
---
type: job-spec
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
outcomes declared by that step. Headings in code fences and blockquotes do
not define steps or outcomes. Progress cases and evidence retain their full
Markdown and named model fields; they are not interpreted as job-map steps.

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
Relative file links and cross-file fragments are checked by check commands;
see [References](../inspection.md#references). Parsing and ordinary show/list commands do not
follow those links. Remote URLs remain navigation.

### Checks and incomplete work

Malformed metadata, duplicate IDs or conventional sections, invalid step heading
structure, duplicate step metadata, and unresolved local references fail with a
nonzero exit status. Missing sections and unfinished content remain readable.

The default writing check suggests a missing job statement, performer, or empty
outcome statement. Missing IDs, maps, stages and outcome collections are valid.
Optional `--preset odi`, `--preset progress`, and `--preset job-story` guidance
adds method-specific suggestions. ODI prompts for a functional boundary, an
optional map, and measurable wording. Progress prompts for cases, circumstances,
participants, struggling moments, alternatives and trade-offs. Job-story guidance
prompts for situation, motivation and intended progress. All are advisory;
none establishes completeness, correctness or research validation.

Job-map and outcome-list JSON use `schema_version: 2` and `declared_by_job`
for outcomes declared by the job. The name describes ownership, not inferred scope.

The same presets work with `job init`. `general` is the default, and the selection
is authoring guidance rather than a classification stored on the job:

```sh
dotspec job init reach-work --preset progress
dotspec job check reach-work --preset odi
dotspec model needs
```

See [Jobs and user needs](../../handbook/jobs.md) for the method distinctions,
progress-case fields, optional capabilities, and research evidence boundaries.
The [generated job section](../../../examples/todo/specification/jobs/resume-interrupted-work.md) shows how the restart requirement supports the
outcome of spending less time reconstructing outstanding work.

The repository's [Develop Requirements](../../../specs/jobs/develop-requirements.job.spec.md) and
[Manage Requirements](../../../specs/jobs/manage-requirements.job.spec.md) examples each contain a
bounded job statement, performer, proposed steps, and outcomes. Both are marked
as hypotheses pending research. They cover collaboration and judgment outside
the CLI as well as work dotspec may support.
