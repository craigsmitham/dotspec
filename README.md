# dotspec

A framework and CLI for managing software product requirements in `*.spec.md` files
and discovering other specification sources through ordinary file references.

System specs describe a system's name, purpose, and scope. Subsystems use the
same format, placed in a `specs/` directory within the part of the project they
describe.
Job specs describe jobs to be done (JTBDs): the accomplishment someone seeks,
their context, the steps involved, and how success is judged. Jobs are associated
with the nearest enclosing system.
Use-case specs describe an actor’s goal and interaction with a stated system
boundary, including success criteria and exceptional conditions.
Concept specs describe a coherent piece of product behavior through its purpose,
abstract state, actions, and operational principle.
Requirement specs state required results, rules, properties, or restrictions directly.
Quality requirements specialize that model with assessable criteria and
measurement methods. Optional quality attributes in a system spec define the
vocabulary used to classify general and quality requirements.

The [CLI reference](docs/cli.md) and [command examples](docs/examples.md) are
generated from the implementation. See [Contributing](CONTRIBUTING.md) for the
local/CI checks and [Releasing](docs/releasing.md) for distribution.

## Design principle: keep meaning in the source; keep connections simple

Connect specification sources through ordinary file references, preserving their
native form. Dotspec provides discovery and reference checking. Consuming projects
determine which sources carry specification weight and ensure those sources are
discoverable; readers and agents interpret their meaning.

This principle guides dotspec's design to reduce duplicate maintenance and keep
sources useful to people and tools. Prefer ordinary files, Markdown links, and
explanatory prose. Add machine-readable structure only when a demonstrated
operation cannot reasonably work without it.

An executable specification can own its required examples: link to it rather
than copying assertions into Markdown. When a test checks a requirement already
stated in prose, retain that definition and explain the test's verification role.
A passing implementation test does not acquire authority merely by being linked.
File discovery makes sources available together; deciding whether they adequately
specify a system remains an authoring and review task.

## Run

With Go 1.27.3 or later installed:

```sh
go run . system show
go run . system list
go run . job list
go run . job show develop-requirements
go run . use-case list
go run . use-case show check-specifications
go run . concept list
go run . concept show job-map
go run . requirement list
go run . requirement show reject-existing-target
go run . sources
go run . check
```

## Build

```sh
go build -o dotspec .
./dotspec --help
```

## Project layout

Each part of the project can have a `specs/` directory. It describes its parent
project directory; folders inside `specs/` organize specifications and do not create
subsystems or new scopes.

```text
specs/
  system.spec.md
  jobs/
    develop-requirements.job.spec.md
    manage-requirements.job.spec.md
  use-cases/
    check-specifications.use-case.spec.md
  concepts/
    job-map.concept.spec.md
    specification-inspection.concept.spec.md
  requirements/
    reject-existing-target.requirement.spec.md
    continue-after-invalid-spec.requirement.spec.md
    advice-permits-success.requirement.spec.md
  quality/
    diagnostic-severity.quality.spec.md
    repeatable-inspection.quality.spec.md
    inspection-preserves-source.quality.spec.md
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
- Concepts are `specs/concepts/<slug>.concept.spec.md`, directly inside the
  plural `concepts/` folder.
- General requirements are `specs/requirements/<slug>.requirement.spec.md`,
  directly inside the plural `requirements/` folder.
- Quality requirements are `specs/quality/<slug>.quality.spec.md`,
  directly inside the `quality/` folder.
- CLI namespaces and file type suffixes remain singular: `dotspec job` and
  `.job.spec.md`, `dotspec use-case` and `.use-case.spec.md`, `dotspec concept`
  and `.concept.spec.md`, `dotspec requirement` and `.requirement.spec.md`, or `dotspec quality` and
  `.quality.spec.md`. Collections use plural names such as `jobs` and `concepts`,
  or the mass noun `quality`.
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

`check` discovers system, job, use-case, concept, requirement, and quality specifications using the same
Git ignore, repository-boundary, and symlink rules as `list`. Outside Git, ignore files do
not affect discovery. It reports the selected scope, supported types, discovery
mode, files examined, and whether the scan completed. Other specification types
are not inspected. An empty completed scan succeeds and explicitly reports
`No supported specifications found.`

Structural errors fail the check: misplaced files, invalid metadata, missing
required system sections, duplicate IDs or conventional sections, malformed
step structure, and unresolved same-document references. All supported types
check Markdown heading references. Jobs, use cases, concepts, general and quality
requirements, and system quality attributes also resolve their explicit document-local IDs.
Quality models check broader-attribute references and reject cycles.
Checks also resolve relative file links and Markdown fragments as described in
[References](#references). Remote URLs remain navigation.

The [advice-permits-success requirement](specs/requirements/advice-permits-success.requirement.spec.md)
governs checks with only advisory findings. Writing advice uses the type-specific content checks
below; missing maps, outcomes, optional IDs, and stages remain valid. Repeated stage
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
The [continued-inspection requirement](specs/requirements/continue-after-invalid-spec.requirement.spec.md)
governs progress through other discovered specifications. Repairing an error
may reveal another on the next run. Findings are sorted
by path, line, column, rule, and message. The
[repeatability requirement](specs/quality/repeatable-inspection.quality.spec.md)
defines equivalence across repeated inspections.
The [diagnostic-severity requirement](specs/quality/diagnostic-severity.quality.spec.md)
governs text presentation of severity, and the
[source-preservation requirement](specs/quality/inspection-preserves-source.quality.spec.md)
governs the effect of inspection on source files.

Jobs, use cases, concepts, general requirements, and quality requirements also need their enclosing
system inspected. General and quality requirements resolve attributes only after that
system passes structural validation; an invalid model does not trigger
cascading undefined-attribute advice. As with
`job show`, `use-case show`, and `concept show`, this direct lookup can read an
ancestor system outside the selected subtree or one excluded from discovery by
an ignore rule. Such files are listed as additional context
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
The focused `job check <slug>`, `use-case check <slug>`,
`concept check <slug>`, `requirement check <slug>`, and `quality check <slug>` commands use the same report, diagnostic IDs, and
`0`/`1`/`2` exit behavior. Each inspects one selected
specification, its enclosing system, and their direct reference targets, with `scope.discovery: "direct"` and
`scope.types: ["job"]`, `["use-case"]`, `["concept"]`, `["requirement"]`, or `["quality"]`. A missing specification
is an incomplete check (exit `2`). File counts represent the selected candidate and attempted read,
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
| `layout.system-location`, `layout.job-location`, `layout.job-name`, `layout.use-case-location`, `layout.use-case-name`, `layout.concept-location`, `layout.concept-name`, `layout.directory` | File placement, filenames, or spec-directory structure. |
| `metadata.frontmatter`, `metadata.invalid`, `metadata.type` | Frontmatter fences, valid YAML and accepted fields, or the declared type. |
| `structure.title`, `structure.section-order`, `structure.required-section`, `structure.duplicate-section`, `structure.section-level`, `structure.step-heading`, `structure.duplicate-field` | Required document structure and unambiguous definitions. |
| `identity.invalid`, `identity.duplicate` | Explicit ID syntax and uniqueness within one document. |
| `reference.invalid`, `reference.unresolved`, `reference.boundary` | Local reference syntax, target resolution, and repository boundaries. |
| `job.statement`, `job.performer`, `job.boundary`, `job.outcome-empty`, `job.outcome-measure` | Advisory job-writing suggestions. |
| `use-case.goal`, `use-case.system-boundary`, `use-case.primary-actor`, `use-case.goal-level`, `use-case.success-criterion`, `use-case.main-success-scenario`, `use-case.sequence` | Advice about missing use-case content or scenario order. |
| `concept.purpose`, `concept.state`, `concept.actions`, `concept.operational-principle` | Advice about missing concept content. |
| `scan.scope`, `scan.discovery`, `scan.read`, `scan.context`, `scan.reference` | Operational failures that prevent a complete check. |

There is one check command rather than a separate lint command. It is read-only:
it does not move files, rewrite prose, assign IDs, or apply fixes.

## Discover specification sources

```sh
dotspec sources
dotspec sources --path internal/integration --json
```

`sources` inventories discovered dotspec files, consulted system context, and
existing files directly referenced by those specs. Ordinary Markdown links work
in every spec type and in any prose section; no special heading or metadata is
required. For example, a module's concept can contain:

```markdown
## Related specifications

These scripts specify required outcomes for the illustrated command inputs:

- [Create and inspect specifications](../../testdata/workflow.txtar)
- [Discovery behavior](../../testdata/discovery.txtar)
```

The [CLI Interaction concept](internal/integration/specs/concepts/cli-interaction.concept.spec.md)
is a working example beside its executable acceptance scripts. The
[source-preservation quality requirement](specs/quality/inspection-preserves-source.quality.spec.md)
instead links to tests that check its prose criterion. Dotspec exposes both
links without assigning either file a semantic role.

The command uses the same scan, diagnostics, reference rules, and exit statuses
as `check`, including advisory findings. Results go to stdout even when the scan
fails. Invocation and output errors go to stderr. `--path` selects a directory
subtree, not a declaration that every source applies to one system.

The version-1 JSON result has these fields:

| Field | Meaning |
| --- | --- |
| `schema_version` | `1` |
| `scope` | The same resolved path, supported types, and discovery mode as `check`. |
| `complete` | Whether the scan completed, independently of structural findings. |
| `files` | Existing regular files, sorted by absolute path, each with `path` and `referenced_by`. |
| `diagnostics` | The same findings and location format as `check`. |

Each `referenced_by` entry is the referring spec's location, including its path,
line, column, and authored link destination in `element`. Entries are sorted by
path and source position. Files appear once even if several specs reference them;
multiple authored link locations are retained. Empty collections are `[]`.
Discovered specs and consulted context can have an empty `referenced_by` array.
Malformed but existing spec files remain listed; missing targets appear in
diagnostics rather than as existing files. Do not treat a partial or failed
inventory as proof of coverage.

The inventory follows one reference edge from discovered specs and their system
context. Referenced documents do not start a recursive scan. Directory links do
not include their contents; ignored files can be explicitly referenced within
the existing repository boundaries. URLs and absolute-path links remain external
navigation. Images and links in code samples are excluded. Native files are not
parsed or executed, and native fragments are not resolved by dotspec. Markdown
heading and declared spec-ID fragments retain the existing reference checks.

Projects can compare their expected specification sources with `files[].path`
and use the paths to read or assemble original content. Our
[coverage script](scripts/check-spec-sources/main.go) invokes this command and
requires links to the repository's designated acceptance and quality-verification
tests. It intentionally does not select every test file. See
[Specification tests and implementation tests](CONTRIBUTING.md#specification-tests-and-implementation-tests)
for this repository's policy and the CI task that enforces file coverage.

## Output and automation

The system, job, use-case, concept, requirement, and quality commands (including
`init`), as well as `sources`, accept `--json`
or `--format json`. Use one of these flags, not both. The default is `--format text`; `job map` also
accepts `--format mermaid` and emits standalone diagram source. Redirecting output
does not change its format. JSON, Mermaid, and redirected output contain no ANSI
styling. `NO_COLOR` disables terminal styling, including help and errors.

Every JSON result has `schema_version: 1`. Check results use the report documented
above. Other results contain:

| Command | Result fields |
| --- | --- |
| `system init`, `job init`, `use-case init`, `concept init`, `requirement init`, `quality init` | `path` of the created specification |
| `sources` | `scope`, `complete`, `files`, and `diagnostics` as described above |
| `system show` | `system` and its nullable `parent` |
| `system list` | `systems`, a tree of `system` objects and `children` arrays |
| `job show` | `job`, including its authored Markdown, steps, outcomes, and nullable `system` |
| `job list` | `groups`, each with a nullable `system` and a `jobs` array |
| `use-case show` | `use_case`, including its authored Markdown, conventional section bodies, and nullable `system` |
| `use-case list` | `groups`, each with a nullable `system` and a `use_cases` array |
| `concept show` | `concept`, including its authored Markdown, four conventional section bodies, and nullable `system` |
| `concept list` | `groups`, each with a nullable `system` and a `concepts` array |
| `requirement show` | `requirement`, including its authored Markdown, statement, verification, attributes, and nullable `system` |
| `requirement list` | `groups`, each with a nullable `system` and a `requirements` array containing general and quality requirements |
| `quality show` | `requirement`, including its authored Markdown, criterion, measurement, attributes, and nullable `system` |
| `quality list` | `groups`, each with a nullable `system` and a `requirements` array containing quality requirements |
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
  purpose, and scope. The
  [existing-target requirement](specs/requirements/reject-existing-target.requirement.spec.md)
  governs initialization at an existing path.
- `show` reads the selected project directory's `specs/system.spec.md`, or the
  nearest ancestor's if there is no local spec. It displays the name, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the Git repository root
  (including a worktree root), or the filesystem root outside Git.
- `list` displays a tree of systems within the selected directory and its
  descendants, ordered by source path. It does not include ancestor systems
  outside that subtree. An empty result is successful and reports that no
  systems were found.

For `system list`, `job list`, `use-case list`, `concept list`, and
`quality list`, and `requirement list`, discovery in Git
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
headings. The optional `Quality attributes` section is also displayed by
`system show`
and interpreted as described below. Other sections remain authored Markdown.
There is no requirement inheritance behavior.

The nearest ancestor project directory containing `specs/system.spec.md` defines
the parent system:

```text
specs/system.spec.md                  # dotspec
internal/parser/specs/system.spec.md  # Spec parser, a subsystem of dotspec
```

Moving a spec can therefore change its parent. Source paths in command output
are relative to the current directory.

## System quality attributes

A system may define its local quality vocabulary in a `Quality attributes` H2
section. Omit the section when a vocabulary is not needed. Its H3 headings
introduce attributes; deeper headings belong to their definitions.

```markdown
## Quality attributes

### Usability

**ID:** usability

How effectively contributors can author, inspect, and maintain specifications.

### Diagnosability

**ID:** diagnosability
**Broader:** usability

How effectively findings help contributors understand and locate problems.

### Accessibility

**ID:** accessibility
**Broader:** usability

How effectively people with differing abilities and ways of accessing output
can use the CLI.
```

Each attribute needs a nonempty, document-unique `ID`. IDs begin with an ASCII
letter or digit and otherwise contain letters, digits, hyphens, or underscores.
They do not encode hierarchy: renaming or regrouping an attribute can preserve
its ID. Same-document links can target these IDs as well as heading anchors.

`Broader` is optional. When present it is a comma-separated list of distinct
attribute IDs in this model; forward references are allowed. Several broader
attributes are allowed. Unknown references, self-references, cycles, duplicate
fields, and duplicate IDs are structural errors. Heading depth does not imply
an attribute relationship. A missing definition produces nonfatal advice.

Only `ID` and `Broader` lines in the attribute's introductory paragraphs are
fields. Examples in code blocks, blockquotes, lists, and deeper subsections remain
prose. Definitions retain their Markdown, including subsections and examples.

`system show` displays the authored model. System JSON adds `quality_model`,
which is `null` when the section is absent. A declared model contains `markdown`
and an `attributes` array, including when that array is empty. Each attribute
exposes `id`, `name`, `definition`, `broader`, and source `location`.

The model belongs to this system. A subsystem does not inherit or merge its
parent's definitions. Attribute grouping imposes no obligations or thresholds.
External models, imports, and automatic edition mappings are not implemented.
The repository's [system spec](specs/system.spec.md#quality-attributes) provides
a local example, including concerns that do not yet have requirements.

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
Relative file links and cross-file fragments are checked by check commands;
see [References](#references). Parsing and ordinary show/list commands do not
follow those links. Remote URLs remain navigation.

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
ID. Check commands validate relative file links using the shared
[reference rules](#references); references
such as `Related jobs` do not establish that an outcome is achieved. Job steps
and use cases can relate many to many. There is no test-result matching or
requirement-satisfaction inference in this version.

The repository's [Check specifications before committing](specs/use-cases/check-specifications.use-case.spec.md)
example connects this interaction to the Develop Requirements job and references
the governing check report contract.

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
  `purpose`, `state`, `actions`, and `operational_principle` as Markdown strings.
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

## Operational principle

Create a task to remember work. Complete it when the work is finished.
If more work is needed, reopen the same task and complete it again.
```

The four conventional section names are case-sensitive H2 headings. Their order
is flexible, and each may appear at most once. Missing or empty sections produce
advice, so unfinished concepts remain readable. Additional sections and nested
headings remain authored Markdown. A nested heading named `State`, for example,
belongs to its containing section and does not replace the concept's H2 `State`.

State describes the information needed to explain behavior; it need not match a
database schema. Actions can use prose, tables, code blocks, or links to governing
artifacts. Action subheadings have ordinary Markdown anchors, but are not separate
records. Describe inputs, outputs, conditions, rejection behavior, and effects
where relevant. Distinguish caller assumptions from enforced guards explicitly;
no heading or notation automatically selects contract semantics.

The operational principle explains how the design serves its purpose. Checking
its presence does not establish complete behavior, preserved invariants, or
requirement satisfaction. Checks resolve same-document heading and ID references;
relative file links are validated by check commands, while remote links remain
navigation. A link alone
does not designate an authoritative definition or associate a test result.

Use cases can reference concept actions. Product-level rules composing several
concepts belong in a shared composition description, linked from the participating
views. Generic bindings such as `TaskLists = Collection<TaskId>` can be described
in an additional section with a link to the generic definition. Typed bindings,
synchronization records, cross-file identity, and executable authority are not
part of this format.

See [Job Map](specs/concepts/job-map.concept.spec.md) and
[Specification Inspection](specs/concepts/specification-inspection.concept.spec.md)
for concepts describing dotspec itself. The
[check-before-committing use case](specs/use-cases/check-specifications.use-case.spec.md)
shows Specification Inspection in a contributor's workflow.

## Development

```sh
go test ./...
go vet ./...
```

Tests that exercise Git ignore rules and repository boundaries require Git.

## Quality commands

```sh
dotspec quality init diagnostic-severity
dotspec quality show diagnostic-severity
dotspec quality list
dotspec quality list --attribute usability --json
dotspec quality check diagnostic-severity --json
```

Every command accepts `--path DIR`, selecting an existing project directory,
and `--json` or `--format json`.

- `init <slug>` creates `specs/quality/<slug>.quality.spec.md`
  without overwriting an existing file or following symlinked metadata directories.
- `show <slug>` reads the requirement in that project, displays its authored
  Markdown and classifications, and reports its associated system and source.
  Slug lookup does not search ancestor projects.
- `list` discovers requirements in the selected subtree, grouped by associated
  system. `--attribute ID` includes exact labels and attributes narrower than
  that ID in each requirement's own model. Each matching requirement appears
  once. Without a model, only exact labels match. Same-named attributes in
  different systems remain separate definitions.
- `check <slug>` checks the requirement and its system context with the shared
  report format and exit behavior. Project-wide `dotspec check` includes this type.

## Quality format

A quality spec states one assessable quality requirement. The shorter `quality`
name is used for the type, commands, folder, and file suffix.

```markdown
---
type: quality
id: diagnostic-severity
attributes:
  - diagnosability
  - accessibility
---

# Diagnostic severity is available without styling

## Applies to

Human-readable reports from project-wide and focused specification checks.

## Criterion

Every diagnostic must identify its severity using a text label independently
of terminal styling.

## Measurement

Capture output containing each severity, including with color disabled.
Compare the labels with the reported severities and retain failures.

## Rationale

Contributors need to distinguish severities without relying on color.
```

Only `type: quality` and one nonempty document H1 are required.
`id` is optional and follows the document-local ID rules. `attributes` is an
optional list of unique, nonempty string IDs; an empty list is allowed. A scalar,
null, non-string element, invalid ID, or duplicate label is a structural error.
Unknown metadata fields are rejected.

`Applies to`, `Criterion`, `Measurement`, and `Rationale` are optional H2 sections.
Missing or empty sections produce advice; repeated conventional sections are
errors. Deeper headings belong to their enclosing sections. Additional sections,
links, and other authored Markdown are preserved. Check commands validate
relative file links using the shared [reference rules](#references).

The nearest enclosing system supplies the model for label resolution. This is
independent of the subject and conditions stated in `Applies to`. If no model
exists, labels remain unbound without advice. With a model, each undefined label
produces `quality.undefined-attribute` advice at its source location.
An explicitly empty model consequently reports every used label as undefined.

JSON retains `attributes` exactly as authored. The separate `classification`
object exposes `model_path`, `resolved` attribute definitions, `broader` ancestor
IDs in model declaration order, and `undefined` IDs in authored order. With no
model, `model_path` is empty and all three arrays are empty. Broader groups are
computed navigation, not additional authored labels or inherited requirements.
The show envelope uses `quality`; list groups use `requirements`.
Both envelopes include `schema_version: 1`. Empty collections are arrays.

Use multiple attributes when they classify the same coherent obligation. Split
obligations that need independent management, even if they share an attribute.
Choose the most specific applicable labels; broad attributes can also be used.
Grouping beneath an attribute does not establish that the attribute is satisfied.

State the relevant subject, conditions, measure, target, and assessment method.
For quantitative criteria include units, sampling, aggregation, and observation
windows where applicable. If a target is unknown, explicitly record the open
decision and the evidence needed to choose it. The checker does not evaluate
prose, invent thresholds, assess test results, or decide requirement satisfaction.

The repository examples are
[diagnostic severity](specs/quality/diagnostic-severity.quality.spec.md),
[repeatable inspection](specs/quality/repeatable-inspection.quality.spec.md),
and [source preservation](specs/quality/inspection-preserves-source.quality.spec.md).
They own their criteria; the report-format reference and inspection concept link
to them. Responsiveness has a model definition but no adopted threshold.

Additional rules for these formats:

| Rule ID | Meaning |
| --- | --- |
| `layout.quality-location`, `layout.quality-name` | Required layout and slug |
| `structure.attribute-heading` | Quality attributes must use nonempty H3 headings |
| `identity.required` | Attribute definition needs an explicit ID |
| `quality.broader` | Broader list is empty or contains an undefined reference |
| `quality.cycle` | Broader graph contains a cycle |
| `quality.definition` | Attribute definition is missing (advice) |
| `quality.attributes` | Attribute list has an invalid shape or element |
| `quality.undefined-attribute` | Label is absent from the selected model (advice) |
| `quality.applies-to`, `quality.criterion`, `quality.measurement`, `quality.rationale` | Missing section content (advice) |

Existing identity, duplicate-field, section, and reference rules also apply.


## Requirements

Use a general requirement to state an obligation without first writing a concept
or use case. A requirement can express a response, postcondition, invariant,
relationship, or imposed restriction. Its type does not implicitly mean
“functional”; functional/non-functional classification is not required.

```sh
dotspec requirement init reject-existing-target
dotspec requirement show reject-existing-target --json
dotspec requirement list
dotspec requirement list --attribute reliability --json
dotspec requirement check reject-existing-target --json
```

`init`, `show`, and focused `check` address general requirements at
`specs/requirements/<slug>.requirement.spec.md`. They follow the same path,
slug, output, and initialization conventions as other types.

`requirement list` is the common inventory of general and quality requirements,
grouped by their associated system. Every entry includes its authored `type`,
source path, and Markdown. In this inventory, quality `Criterion` and
`Measurement` become `statement` and `verification`; the source file remains
unchanged. `quality list` remains the specialized inventory of quality specs.
Attribute filtering applies equally to both forms and includes narrower
attributes in the selected system's model.

```markdown
---
type: requirement
id: continue-after-invalid-spec
attributes: [diagnosability]
---

# Continue checking after a structural error

## Applies to

Project-wide specification checks.

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Rationale

One malformed specification must not hide problems in other files.

## Verification

Check two files with distinct structural errors and verify that both appear
in the report.
```

Only the type and one nonempty H1 are required. Optional `id` and `attributes`
follow the same rules as quality specs. `Applies to`, `Statement`, `Rationale`,
and `Verification` are optional H2 sections: omissions produce advice, duplicate
sections or invalid levels produce structural errors. Nested headings, contract
notation, tables, and additional prose are retained. No sentence pattern or
particular verification technique is required.

The show envelope uses `requirement`; the inventory uses `groups` containing
`requirements`. Both include `schema_version: 1`. Each requirement exposes its
`type`, `statement`, `verification`, `applies_to`, `rationale`, `attributes`,
`classification`, `markdown`, source identity, and associated system.

Keep one authoritative statement of each obligation. Concepts organize related
rules and quality requirements; use cases illustrate interactions and link to
the rules that govern them. Preserve required logical dependencies without
turning a representative sequence into an obligation. Distinguish caller
assumptions from conditions the system must enforce or reject.

Examples are [existing-target rejection](specs/requirements/reject-existing-target.requirement.spec.md),
[continued inspection](specs/requirements/continue-after-invalid-spec.requirement.spec.md),
and [successful advisory checks](specs/requirements/advice-permits-success.requirement.spec.md).
General requirements and quality requirements share parsing, identity,
classification, and reference behavior. Neither form establishes fulfillment
merely by being present or passing a structural check.

Rules specific to general requirements are `layout.requirement-location`,
`layout.requirement-name`, `requirement.attributes`,
`requirement.undefined-attribute`, `requirement.applies-to`,
`requirement.statement`, `requirement.rationale`, and `requirement.verification`.
The shared structural and identity rules also apply.

## System context

System specs may include `Goals`, `Environment`, `Assumptions`, `Constraints`,
and `Glossary` H2 sections. These are authored Markdown displayed by
`system show` and exposed in JSON when nonempty. Duplicate sections are errors;
missing or empty context sections produce no advice. `Purpose` and `Scope`
remain required.

Use goals for intended organizational or product results, environment for
relevant surrounding conditions, and assumptions for conditions relied on.
Record imposed restrictions under constraints, referencing independently
managed obligations or existing project policies when appropriate. Keep each
definition in one place. A glossary can use H3 headings for individual meanings;
these are heading targets, not a new term record format. Link to a concept when
it already defines the term.

## References

Every check command validates Markdown links in successfully parsed selected
specifications and their associated system context. Same-document fragments
are checked during parsing as well. Links inside code samples and images are
not interpreted as specification references.

Relative paths resolve from the referring file. A link can address a whole file
or directory, or a Markdown heading or declared ID in another specification,
for example `../jobs/develop-requirements.job.spec.md#OUT-009`.
IDs remain document-local; the file path disambiguates identical IDs in different
files. The target parser supplies recognized IDs, including job steps and
outcomes and system quality attributes. This does not introduce IDs for arbitrary
paragraphs or embedded requirement clauses.

Direct targets can be outside the selected subtree or ignored by discovery.
Targets are counted once as additional context unless already discovered.
Only the outgoing links of selected specs and their system context are checked;
a referenced document does not start a recursive scan. A supported target spec
is structurally parsed, without collecting its writing advice. Invalid targets
produce their structural diagnostic without cascading missing-fragment errors.

Ordinary Markdown targets provide heading anchors without needing a dotspec
wrapper. Other native file formats are checked for existence; their fragments
retain the owning format's semantics. URLs and absolute paths remain external
navigation. No network requests are made.

Relative links cannot cross repository/worktree boundaries, enter nested
repositories or Git metadata, or traverse symbolic links. Outside Git, lookup is
bounded by the filesystem root. Missing files or Markdown targets produce
`reference.unresolved`; malformed relative references produce `reference.invalid`;
prohibited paths produce `reference.boundary`. I/O failures produce
`scan.reference` and leave the scan incomplete.

References establish navigation and explicit connections. They do not infer
normative authority, automatic inheritance, consistency, or requirement
satisfaction.
