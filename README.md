# dotspec

A framework and CLI for managing software product requirements in `*.spec.md` files.

System specs describe a system's name, purpose, and scope. Subsystems use the
same format, placed in a `specs/` directory within the part of the project they
describe.
Job specs describe jobs to be done (JTBDs): the accomplishment someone seeks,
their context, the steps involved, and how success is judged. Jobs are associated
with the nearest enclosing system.
Use-case specs describe an actor’s goal and interaction with a stated system
boundary, including success criteria and exceptional conditions.

The [CLI reference](docs/cli.md) and [command examples](docs/examples.md) are
generated from the implementation. See [Contributing](CONTRIBUTING.md) for the
local/CI checks and [Releasing](docs/releasing.md) for distribution.

## Run

With Go 1.27.3 or later installed:

```sh
go run . system show
go run . system list
go run . job list
go run . job show develop-requirements
go run . use-case list
go run . use-case show check-specifications
go run . check
```

## Build

```sh
go build -o dotspec .
./dotspec --help
```

## Project layout

Each part of the project can have a `specs/` directory. It describes its parent
project directory; folders inside `specs/` organize concepts and do not create
subsystems or new scopes.

```text
specs/
  system.spec.md
  jobs/
    develop-requirements.job.spec.md
    manage-requirements.job.spec.md
  use-cases/
    check-specifications.use-case.spec.md
internal/
  parser/
    parser.go
    specs/
      system.spec.md
      jobs/
        interpret-specification.job.spec.md
```

The layout is required:

- The system spec is `specs/system.spec.md`, with one system per project directory.
- Jobs are `specs/jobs/<slug>.job.spec.md`, directly inside the plural `jobs/` folder.
- Use cases are `specs/use-cases/<slug>.use-case.spec.md`, directly inside the
  plural `use-cases/` folder.
- CLI namespaces and file type suffixes remain singular: `dotspec job` and
  `.job.spec.md`, or `dotspec use-case` and `.use-case.spec.md`. Collection
  directories are plural. Future types will follow the same convention, for
  example `quality-requirements/`; quality-requirement specs are not implemented yet.
- Create folders only when needed. `init` creates the required spec directories;
  the selected project directory must already exist.

`--path` always selects a project directory, such as `internal/parser`, rather
than `internal/parser/specs` or its `jobs/` folder. Using a metadata directory as
a project scope is an error. A nested `specs/` tree inside another `specs/` tree
is not a separate scope.

Adjacent spec files and singular collection folders are not supported. Move an
existing `system.spec.md` to `specs/system.spec.md`, and move each adjacent
`<slug>.job.spec.md` to `specs/jobs/<slug>.job.spec.md`. Discovery rejects misplaced
files of the requested type with an expected-location diagnostic. Local lookup
and initialization also reject misplaced specs of that type directly in the
selected project directory or anywhere within its `specs/` tree. Having both
old and new locations is an error; there is no location precedence or automatic
migration. Revise any relative file links when moving documents.

## Check specifications

Before committing specification changes, check the selected project subtree:

```sh
dotspec check
dotspec check --path internal/parser
dotspec check --format json
dotspec job check develop-requirements
```

`check` discovers system, job, and use-case specifications using the same Git ignore,
repository-boundary, and symlink rules as `list`. Outside Git, ignore files do
not affect discovery. It reports the selected scope, supported types, discovery
mode, files examined, and whether the scan completed. Other specification types
are not inspected. An empty completed scan succeeds and explicitly reports
`No supported specifications found.`

Structural errors fail the check: misplaced files, invalid metadata, missing
required system sections, duplicate IDs or conventional sections, malformed
step structure, and unresolved same-document references. All supported types
check Markdown heading references; jobs and use cases also resolve their explicit
IDs. File and remote links are not checked.

Writing advice is nonfatal. It uses the job-writing and use-case checks below;
missing maps, outcomes, optional IDs, and stages remain valid. Repeated stage
labels and jobs without an associated system also remain valid. A successful
check makes no claim about completeness, research validation, or whether
desired outcomes are achieved.

Each finding includes a stable rule ID, severity (`error` or `advice`), message,
and source path. Findings include a section or element, line, and column when
available; duplicate definitions include the earlier source location too.
Lines and columns are one-based, columns count Unicode code points, and line
numbers include frontmatter. Layout errors have file locations without invented
line numbers. YAML errors use the reported line with a column at its start.

The checker continues across malformed files. Within a file, it reports the
first structural error and skips dependent checks, including writing advice.
Repairing that error may reveal another on the next run. Findings are sorted
by path, line, column, rule, and message for repeatable output.

Jobs and use cases also need their enclosing system inspected. As with
`job show` and `use-case show`, this direct lookup can read an ancestor system
outside the selected subtree or one excluded
from discovery by an ignore rule. Such files are listed as additional context
and counted separately; shared findings appear once. Context lookup stops at
the repository or worktree root, or the filesystem root outside Git.
Failures to inspect a context directory are reported but do not count that
directory as a context file.

All commands use these exit codes:

| Code | Meaning |
| --- | --- |
| `0` | Command completed successfully; check advice may be present. |
| `1` | Specification findings prevented success. |
| `2` | Invalid invocation, missing input, failed discovery or file access, or output failure. |

Text and JSON reports are written to standard output, including reports with
findings or an incomplete scan. Argument and output errors go to standard error.
Use the exit status as well as the report; a partial report is not a clean check.
For example, a CI step can run `dotspec check --format json > check-report.json`.
The focused `job check <slug>` and `use-case check <slug>` commands use the same
report, diagnostic IDs, and `0`/`1`/`2` exit behavior. Each inspects one selected
specification and its enclosing system, with `scope.discovery: "direct"` and
`scope.types: ["job"]` or `["use-case"]`. A missing specification is an incomplete
check (exit `2`). File counts represent the selected candidate and attempted read,
including a missing file. All check commands send their reports to stdout;
ordinary read-command errors go to stderr. This replaces the earlier focused
check behavior that mixed findings into stderr and used only `0`/`1`.

JSON schema version `1` contains:

| Field | Meaning |
| --- | --- |
| `schema_version` | Report schema version for consumers. |
| `scope` | Absolute resolved project path, supported `types`, and `discovery` mode (`git`, `filesystem`, or `unavailable`). An unresolved scope retains the supplied path; focused checks use `direct` discovery. |
| `complete` | Whether discovery and all attempted file checks completed, independently of structural findings. |
| `context_files` | Additional source paths consulted outside the discovered set. |
| `summary` | `files_discovered`, `files_checked`, `context_files_checked`, `errors`, and `advice`. Checked counts include attempted checks of invalid or unreadable files. |
| `diagnostics` | Findings with `rule`, `severity`, `message`, `location`, and optional `related` locations. Empty findings are an empty array. |

Locations contain `path` and optional `line`, `column`, and `element`. Report
source paths are absolute; text displays them relative to the current directory
when possible. Use rule IDs rather than matching message text in automation.

The current rules are:

| Rule IDs | Meaning |
| --- | --- |
| `layout.system-location`, `layout.job-location`, `layout.job-name`, `layout.use-case-location`, `layout.use-case-name`, `layout.directory` | File placement, filenames, or spec-directory structure. |
| `metadata.frontmatter`, `metadata.invalid`, `metadata.type` | Frontmatter fences, valid YAML and accepted fields, or the declared type. |
| `structure.title`, `structure.section-order`, `structure.required-section`, `structure.duplicate-section`, `structure.section-level`, `structure.step-heading`, `structure.duplicate-field` | Required document structure and unambiguous definitions. |
| `identity.invalid`, `identity.duplicate` | Explicit job/use-case ID syntax and uniqueness within one job document. |
| `reference.invalid`, `reference.unresolved` | Same-document fragment syntax and resolution. |
| `job.statement`, `job.performer`, `job.boundary`, `job.outcome-empty`, `job.outcome-measure` | Advisory job-writing suggestions. |
| `use-case.goal`, `use-case.system-boundary`, `use-case.primary-actor`, `use-case.goal-level`, `use-case.success-criterion`, `use-case.main-success-scenario`, `use-case.sequence` | Advice about missing use-case content or scenario order. |
| `scan.scope`, `scan.discovery`, `scan.read`, `scan.context` | Operational failures that prevent a complete project check. |

There is one check command rather than a separate lint command. It is read-only:
it does not move files, rewrite prose, assign IDs, or apply fixes.

## Output and automation

All system, job, and use-case commands, including `init`, accept `--json` or `--format json`.
Use one of these flags, not both. The default is `--format text`; `job map` also
accepts `--format mermaid` and emits standalone diagram source. Redirecting output
does not change its format. JSON, Mermaid, and redirected output contain no ANSI
styling. `NO_COLOR` disables terminal styling, including help and errors.

Every JSON result has `schema_version: 1`. Check results use the report documented
above. Other results contain:

| Command | Result fields |
| --- | --- |
| `system init`, `job init`, `use-case init` | `path` of the created specification |
| `system show` | `system` and its nullable `parent` |
| `system list` | `systems`, a tree of `system` objects and `children` arrays |
| `job show` | `job`, including its authored Markdown, steps, outcomes, and nullable `system` |
| `job list` | `groups`, each with a nullable `system` and a `jobs` array |
| `use-case show` | `use_case`, including its authored Markdown, conventional section bodies, and nullable `system` |
| `use-case list` | `groups`, each with a nullable `system` and a `use_cases` array |
| `job step list` | `name`, `path`, and `steps` in authored order |
| `job outcome list` | `name`, `path`, `overall` outcomes, and selected `steps` with their outcomes |
| `job map` | `name`, `path`, `overall`, and `steps`; outcomes are included only with `--outcomes` |

Object keys use snake case. Empty collections are arrays; absent system associations
are `null`. Paths are absolute in JSON and relative where possible in text. Source
locations use one-based lines and Unicode-code-point columns. Step descriptions and
Markdown remain authored content. Overall outcomes are never inherited by steps.
`--step ID` leaves `overall` empty and selects only that step.

Use `dotspec --version` for build metadata, `dotspec completion bash|zsh|fish|powershell`
for shell completion scripts, and `dotspec man` for the root manual in roff format.
See the generated [CLI reference](docs/cli.md) for flags and examples.

## System commands

```sh
dotspec system init                     # Create specs/system.spec.md here
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

- `init` creates `specs/system.spec.md` with editable placeholders for the name,
  purpose, and scope. It never overwrites an existing file.
- `show` reads the selected project directory's `specs/system.spec.md`, or the
  nearest ancestor's if there is no local spec. It displays the name, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the Git repository root
  (including a worktree root), or the filesystem root outside Git.
- `list` displays a tree of systems within the selected directory and its
  descendants, ordered by source path. It does not include ancestor systems
  outside that subtree. An empty result is successful and reports that no
  systems were found.

For `system list`, `job list`, and `use-case list`, discovery in Git repositories uses Git's
ignore rules, including nested `.gitignore` files and local exclusions.
Tracked files remain visible even if
an ignore rule matches them. Git must be installed for discovery in a repository.
Outside Git, discovery walks the directory tree without applying ignore rules.
Discovery skips `.git`, nested repositories, and symlinked files and directories;
an explicitly selected symlink directory is resolved before the command runs.

Missing specs for `show`, invalid paths, malformed specs, and command errors
produce a diagnostic on standard error and a nonzero exit status. A malformed
discovered spec makes `list` fail rather than silently omit it. Each list command
reads only its own spec type; `job list` and `use-case list` also read the systems
associated with their specifications.

## System spec format

Each project directory can contain one system spec at `specs/system.spec.md`:

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

The nearest ancestor project directory containing `specs/system.spec.md` defines
the parent system:

```text
specs/system.spec.md                  # dotspec
internal/parser/specs/system.spec.md  # Spec parser, a subsystem of dotspec
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
- `outcome list <slug>` displays success criteria under `Overall job outcomes`
  and `Outcomes for step: ...`. `--step ID` selects just that step's outcomes,
  excluding overall job outcomes. An unknown step is an error; an existing step
  without outcomes is a successful empty result.
- `map <slug>` displays the job and its documented steps in authored order.
  `--format mermaid` emits standalone Mermaid flowchart source for a renderer
  or a fenced `mermaid` block. `--outcomes` includes the whole job's outcomes and
  attaches each step's outcomes to that step, with scope labels in text and Mermaid.
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

The repository's [Develop Requirements](specs/jobs/develop-requirements.job.spec.md) and
[Manage Requirements](specs/jobs/manage-requirements.job.spec.md) examples each contain a
bounded job statement, performer, proposed steps, and outcomes. Both are marked
as hypotheses pending research. They cover collaboration and judgment outside
the CLI as well as work dotspec may support.

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
  template without overwriting existing content.
- `show <slug>` selects that file in the chosen project directory, displays all
  authored Markdown, and reports its enclosing system and source path. It does
  not search ancestor directories for the same slug.
- `list` discovers use cases in the selected subtree and groups them by enclosing
  system. Groups are ordered by system name and source path, with unassociated
  use cases last; use cases within each group are ordered by source path.
- `check <slug>` checks one use case and its enclosing system, reporting structural
  errors and nonfatal missing-content advice. Project-wide `check` includes use
  cases alongside systems and jobs. Ordinary reads do not emit writing advice.

Use cases follow the same discovery, ignore, repository-boundary, and placement
rules as jobs. An enclosing system is used for organization; it does not fill in
or override the use case's stated system boundary. A use case without an enclosing
system is valid and reports `System: None` (`null` in JSON). A malformed enclosing
system is an error.

## Use-case spec format

Only `type: use-case` frontmatter and one nonempty document-level H1 are required.
An optional frontmatter `id` uses the same syntax as job IDs; other frontmatter
keys are rejected. The file slug remains the CLI selector. IDs are local to the
document and are not checked for uniqueness across files.

```markdown
---
type: use-case
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

The selected project directory exists and is readable.

## Success criterion

The contributor receives a complete report for the selected scope.

## Main success scenario

1. The contributor selects the project subtree to inspect.
2. Dotspec inspects its supported specifications.
3. Dotspec reports findings and scan completion.

## Extensions

- A file cannot be read: report the failure and mark the scan incomplete.

## Sequence

Representative. The check report contract defines required behavior.

## Related jobs

[Evaluate requirements](../jobs/develop-requirements.job.spec.md#evaluate-the-requirements).
```

These eleven H2 headings are conventional sections. They are case-sensitive,
optional, may appear in any order after the title, and may be empty in a draft.
Duplicate conventional sections are errors, even if empty. A conventional section
at document level must use H2. Deeper headings belong to the containing section;
their names do not create additional conventional sections. Code examples and
blockquotes do not define sections. Additional prose and sections are retained.

JSON exposes conventional section bodies as Markdown strings in `goal`,
`system_boundary`, `primary_actor`, `goal_level`, `trigger`, `preconditions`,
`success_criterion`, `main_success_scenario`, `extensions`, `sequence`, and
`related_jobs`. Missing bodies are empty strings. The object also includes
`name`, `slug`, optional `id`, `markdown`, absolute `path`, `system`, and the title's
source `location`.

A trigger is an initiating event; preconditions describe an already-expected
state. Success criteria can be stated here or referenced from a governing rule.
`Sequence` explains whether the narrative order is representative or required;
link to the governing rule when order is required. These bodies are prose, not
enumerations or executable contracts. Steps and extensions remain Markdown
inside the use case, without separate records, IDs, or branch-coverage claims.

`check` suggests a goal, system boundary, primary actor, goal level, success
criterion, and main success scenario when their conventional sections are missing
or empty. It also suggests `Sequence` when a main success scenario is present but
that section is missing or empty. It does not interpret prose elsewhere as these
fields or judge whether a filled section is adequate. Missing triggers,
preconditions, extensions, job links, and IDs produce no advice.

Same-document links resolve against heading anchors and the optional document
ID. File and remote links remain navigation and are not validated; references
such as `Related jobs` do not establish that an outcome is achieved. Job steps
and use cases can relate many to many. There is no test-result matching or
requirement-satisfaction inference in this version.

The repository's [Check specifications before committing](specs/use-cases/check-specifications.use-case.spec.md)
example connects this interaction to the Develop Requirements job and references
the governing check report contract.

## Development

```sh
go test ./...
go vet ./...
```

Tests that exercise Git ignore rules and repository boundaries require Git.
