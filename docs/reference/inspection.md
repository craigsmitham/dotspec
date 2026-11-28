# Checking and discovering specifications

## Check specifications

From anywhere in your project, run `dotspec check` to locate specification findings across the whole project. Give it a directory or a specification file for a narrower check while correcting a particular file; relative paths resolve against the working directory. For the examples below, substitute an existing subtree for `internal/parser` and an existing job’s file path for the illustrated one.

```sh
dotspec check internal/parser
dotspec check --json
dotspec check specs/jobs/develop-requirements.job.spec.md
```

Use the first finding’s path and location to inspect the source. Correct a structural error and repeat the check. For advice, judge whether the suggested information is known and worth recording. A missing file or an incomplete scan calls for checking the path or restoring access before relying on the report.

`check` discovers the specification and reference types listed in its `scope.types` using the same
Git ignore, repository-boundary, and symlink rules as `list`. Outside Git, ignore files do
not affect discovery. It reports the selected scope, files examined, and whether
the scan completed; the JSON report also records the supported types and
discovery mode. An empty completed scan succeeds and explicitly reports
`No supported source files found.`

The text report lists the findings first. One summary line follows them, naming
the scope, whether the scan completed and the counts:

```text
Checked 16 files in project todo (scan complete): 1 error, 0 advisories; 2 context files.
Generated specification: not checked (source errors prevent building it)

Context files:
  README.md
  dotspec.yaml
```

A whole-project check names its project; a narrower check names its directory.
When the project builds a generated specification, a sentence after the summary
says whether this check covered it. The context files come last. Text paths are
relative to the working directory, so a specification path printed by `list` or
`sources` can be passed to `check` as it stands.

Structural errors fail the check: conflicting subject declarations, invalid metadata, missing
required system sections, duplicate IDs or conventional sections, malformed
step structure, and unresolved same-document references. All supported types
check Markdown heading references. Jobs, use cases, concepts, requirements,
and system quality attributes also resolve their explicit document-local IDs.
Quality attribute checks validate parent references and reject cycles.
Checks also resolve relative file links and Markdown fragments as described in
[References](#references). Remote URLs remain navigation.

The [advice-permits-success requirement](../../specs/requirements/advice-permits-success.requirement.spec.md)
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
The [continued-inspection requirement](../../specs/requirements/continue-after-invalid-spec.requirement.spec.md)
governs progress through other discovered specifications. Repairing an error
may reveal another on the next run. Findings are sorted
by path, line, column, rule, and message. The
[repeatability requirement](../../specs/quality/reliability/repeatability/repeatable-inspection.requirement.spec.md)
defines equivalence across repeated inspections.
The [diagnostic-severity requirement](../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md)
governs text presentation of severity, and the
[source-preservation requirement](../../specs/quality/reliability/source-preservation/inspection-preserves-source.requirement.spec.md)
governs the effect of inspection on source files.

Kinds that read system context, such as requirements, also need their enclosing
system inspected. Requirements resolve attributes only after that
system passes structural validation; an invalid system does not trigger
cascading undefined-attribute advice. This direct lookup can read an
ancestor system outside the selected subtree or one excluded from discovery by
an ignore rule. Such files are listed as additional context
and counted separately; shared findings appear once. The target of an unresolved
reference is reported as a `reference.unresolved` finding, not listed as context
in text; JSON `context_files` still includes it. Context lookup stops at
the nearest project root, or at the repository/worktree root (filesystem root
outside Git) when there is no enclosing project manifest.
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
An error names the command and the problem, such as `new requirement: missing <slug>`,
and then the help to consult; a misspelled command suggests the nearest one.
Filesystem errors name the path, relative to the working directory where possible,
and state the problem in plain words. The help of each command that can exit with
status `1` says so.
Use the exit status as well as the report; a partial report is not a clean check.
For example, a CI step can run `dotspec check --json > check-report.json`.
A focused `check <file>`, given one specification's source path, uses the same report, diagnostic IDs, and
`0`/`1`/`2` exit behavior. It inspects that
specification, its system context where that kind has one, and their direct reference targets, with `scope.discovery: "direct"` and
`scope.types` set to the selected kind, such as `["requirement"]`. A missing specification
is an incomplete check (exit `2`). File counts represent the selected candidate and attempted read,
including a missing file. All check commands send their reports to stdout;
ordinary read-command errors go to stderr.

JSON schema version `1` contains:

| Field | Meaning |
| --- | --- |
| `schema_version` | `1`, always the first key. |
| `generated_spec` | Present on a whole-project check of a declared project: whether the generated specification was checked. |
| `scope` | Absolute resolved project path, supported `types`, and `discovery` mode (`git`, `filesystem`, or `unavailable`). An unresolved scope retains the supplied path; focused checks use `direct` discovery. |
| `complete` | Whether discovery and all attempted file checks completed, independently of structural findings. |
| `context_files` | Additional source paths consulted outside the discovered set. |
| `summary` | `files_discovered`, `files_checked`, `context_files_checked`, `errors`, and `advice`. Checked counts include attempted checks of invalid or unreadable files. |
| `diagnostics` | Findings with `rule`, `severity`, `message`, `location`, and optional `related` locations. Empty findings are an empty array. |

Locations contain `path` and optional `line`, `column`, and `element`. Report
source paths are absolute; text displays them relative to the current directory
when possible. Every JSON result uses this location shape. Use rule IDs rather than matching message text in automation.

Examples of diagnostic rules include the following.

| Rule IDs | Meaning |
| --- | --- |
| `layout.subject-ambiguous`, `layout.subject-name`, `layout.job-name`, `layout.use-case-name`, `layout.concept-name` | Conflicting subject declarations or invalid filenames. |
| `metadata.frontmatter`, `metadata.invalid`, `metadata.type` | Frontmatter fences, valid YAML and accepted fields, or the declared type. |
| `structure.title`, `structure.section-order`, `structure.required-section`, `structure.duplicate-section`, `structure.section-level`, `structure.duplicate-field` | Required document structure and unambiguous definitions. |
| `identity.invalid`, `identity.duplicate` | Explicit ID syntax and uniqueness within one document. |
| `reference.invalid`, `reference.unresolved`, `reference.boundary` | Local reference syntax, target resolution, and repository boundaries. |
| `catalog.invalid`, `catalog.duplicate`, `catalog.unlisted`, `catalog.unchecked` | Declaration shape, name collisions, candidate files and uninspected locators. |
| `job.statement`, `job.performer`, `job.outcome-empty` | Advice about a missing job statement, performer, or outcome statement. |
| `use-case.goal`, `use-case.system-boundary`, `use-case.primary-actor`, `use-case.goal-level`, `use-case.success-guarantees`, `use-case.main-success-scenario`, `use-case.extensions`, `use-case.sequence` | Advice about missing use-case content or scenario order. |
| `concept.purpose`, `concept.state`, `concept.actions`, `concept.operational-principle` | Advice about missing concept content. |
| `assumption.basis`, `assumption.applicability`, `interface.contract` | Advice about a system assumption with no Basis or Applicability, or a system or component interface that links to no native contract. |
| `scan.scope`, `scan.discovery`, `scan.read`, `scan.context`, `scan.reference` | Operational failures that prevent a complete check. |

There is one check command rather than a separate lint command. It is read-only:
it does not move files, rewrite prose, assign IDs, or apply fixes.

## Discover specification sources

Run from anywhere in your project. For a narrower inventory, substitute a directory containing your specifications for `internal/integration`.

```sh
dotspec sources
dotspec sources internal/integration --json
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

A declaration reports its name, the referring specification's location and its authored source. A governing source also reports its category. The [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md#evidence) declares a verification procedure and cites its Python file in Evidence. The source inventory shows both occurrences beside the same file.

The command uses the same scan, diagnostics, reference rules, and exit statuses
as `check`, including advisory findings. Results go to stdout even when the scan
fails. Invocation and output errors go to stderr. A directory operand selects a subtree. A regular-file operand shows its incoming references across the enclosing project and includes its canonical path as `selected_file` in JSON. For a shared file outside a project the lookup uses the project of the invoking directory.

The version-1 JSON result has these fields:

| Field | Meaning |
| --- | --- |
| `schema_version` | `1`, always the first key. |
| `scope` | The same resolved path, supported types, and discovery mode as `check`. |
| `complete` | Whether the scan completed, independently of structural findings. |
| `files` | Existing regular files, sorted by absolute path, each with `path`, `referenced_by` and `declarations`. |
| `diagnostics` | The same findings and location format as `check`. |

Each `referenced_by` occurrence records path, line, column and element. A frontmatter occurrence points to the source scalar. Each `declarations` entry records the name, source and location. A governing declaration also has `defines`; a nested verification procedure identifies its `verified_specification`. Ordinary contextual links appear only in `referenced_by`. Existing regular files are listed once with incoming occurrences in source order. External URLs remain in the snapshot and generated catalog.

The text inventory groups files under the project that contains them, naming each
project once, and ends with a summary line such as
`Listed 21 source files in project todo (scan complete).`
Malformed but existing spec files remain listed; missing targets appear in
diagnostics rather than as existing files. Do not treat a partial or failed
inventory as proof of coverage.

The inventory resolves links and named sources by the rules in
[References](#references). A directory link does not include the directory's
contents, and native files are not executed or parsed beyond finding the
definition a native fragment names.

To check coverage, compare the source files your project expects with `files[].path` in a complete inventory. For an absent specification, verify the selected directory, typed filename suffix, project boundary and Git ignore rules. For a missing contract or test, inspect the relative link in the referring specification and rerun `dotspec sources` after correcting it. Use `referenced_by` to find the specification that links to a test or contract, and `declarations` to see which specifications declare it and what it governs or assesses.

## Output and automation

To save a snapshot for other tools and validate it, see [Export and validate a snapshot](specification-model.md#export-and-validate-a-snapshot). Offline, run `dotspec snapshot schema` to print the snapshot format; the [mapping reference](../../model/markdown-mapping.md) explains how authored Markdown becomes its fields.

A command that returns a result accepts a local `--json` flag, listed with its other
flags in the [CLI reference](cli.md); without it the result is text. `docs`, `help`,
`completion` and `man` have no result to encode, and `build` writes the generated
specification instead; none of them accepts `--json`. The snapshot artifacts, `snapshot export`
and `snapshot schema`, always write JSON. The snapshot is the integration contract for tools
that consume the whole specification; `--json` on other commands is a per-command
convenience for scripts. Redirecting output does not change its format. JSON and redirected output contain no ANSI
styling. `NO_COLOR` disables terminal styling, including help and errors.

Every command result begins with `"schema_version": 1`. The snapshot is the exception: `snapshot export` carries the `model_version` and `representation_version` of the [specification model](specification-model.md) instead. Check results use the report documented above.

| Command | Result fields |
| --- | --- |
| `init` | `project` name, absolute `manifest` and `system` paths, and the nullable `enclosing_project` (`name` and `root`) whose scans now skip the new project |
| `new <kind>` | `path` of the created specification |
| `new quality-attribute` | the `target` system, each `entries` decision (`id`, `role` of `requested` or `ancestor`, `decision` of `create` or `reuse`, `path` and, for a new file, its `content`), the `added` IDs and `dry_run` |
| `list <kind>` | `entries`, each with its `kind`, `slug`, `name`, absolute `path` and nullable `system` (`name` and `path`): a specification's associated system, a component's governing system or a system's enclosing system. A component entry also carries the `part_of` path of its parent |
| `sources` | `scope`, `complete`, `files`, and `diagnostics` as described above |
| `snapshot validate` | `valid`, and the `findings` that make a snapshot invalid, each with its `path`, `code` and `message` |

Object keys use snake case. Empty collections are arrays; absent system associations
are `null`. Paths are absolute in JSON and relative where possible in text. Every
`location` carries the `path` of its source
file; a bundled quality model, which has no file, is named `builtin:<slug>` instead. Source
locations use one-based lines and Unicode-code-point columns. Step descriptions and
Markdown remain authored content. In a snapshot, each job step carries only the outcomes it declares. To see which steps an outcome applies to, open the job’s generated report, where each step lists its outcomes under **Applies here**. [Outcome ownership and applicability](../cli-handbook/guides/write-a-job.md#record-steps-and-shared-outcomes).

Use `dotspec --version` for build metadata and `dotspec completion bash|zsh|fish|powershell`
for shell completion scripts. Each release archive includes the manual page
`man/man1/dotspec.1`.
See the generated [CLI reference](cli.md) for flags and examples.

## References

A `structure.adoption` finding points to an entry under **Adopts**. Identify a single provision and give both its edition and local applicability. [Source mapping](../../model/markdown-mapping.md#adopting-a-provision)

Every check command validates Markdown links and declared sources
in successfully parsed selected specifications and their associated system context. Same-document fragments
are checked during parsing as well. Links inside code samples and images are
not interpreted as specification references.

Relative paths resolve from the referring file. A link can address a whole file
or directory, or a Markdown heading or declared ID in another specification,
for example `../jobs/develop-requirements.job.spec.md#OUT-009`.
IDs remain document-local; the file path disambiguates identical IDs in different
files. The target parser supplies recognized IDs, including job steps and
outcomes and system quality attributes. This does not introduce IDs for arbitrary
paragraphs.

Direct targets can be outside the selected subtree or ignored by discovery.
Targets are counted once as additional context unless already discovered.
Only the outgoing links of selected specs and their system context are checked;
a referenced document does not start a recursive scan. A supported target spec
is structurally parsed, without collecting its writing advice. Invalid targets
produce their structural diagnostic without cascading missing-fragment errors.

Ordinary Markdown targets provide heading anchors without needing a dotspec
wrapper. Supported native fragments resolve by declaration name for both named sources and contextual links. Other native fragments retain their format's semantics. [Name a definition](../cli-handbook/guides/catalog-a-specification.md#name-a-definition).
URLs and absolute paths remain external navigation. No network requests are made.

Relative links cannot cross repository/worktree boundaries, enter nested
repositories or Git metadata, or traverse symbolic links. Outside Git, lookup is
bounded by the filesystem root. Missing files or Markdown targets produce
`reference.unresolved`; malformed relative references produce `reference.invalid`;
prohibited paths produce `reference.boundary`. I/O failures produce
`scan.reference` and leave the scan incomplete.

References establish navigation and explicit connections. They do not infer
normative authority, automatic inheritance, consistency, or requirement
satisfaction.

For explicit boundary, job-support and applicability links, see
[the Markdown relationship contract](../../model/markdown-mapping.md#relationship-contract).
The generated [relationship reference](model/relationships.md)
documents the metamodel's edge categories, cardinalities, scope and projections.

## Check the generated requirements specification

A whole-project check includes the generated Markdown by default for projects
declared with `dotspec.yaml`. Set `build.enabled` to `false` for source-only checks. Commit the output of `dotspec build` with its
sources. Check all configured projects in CI with `dotspec check --all`.
See [Generate a requirements specification](generated-specification.md) for
document types, workspace boundaries, reader navigation and recovery from drift.

When the generated spec participates, the JSON check report includes
`generated_spec`. Its value is `checked` after the files were inspected,
`blocked-by-source` when source findings prevent constructing the expected
spec, `out-of-scope` for a narrower selection, or `incomplete` when an
operational error interrupted document inspection. The text report states the
same in a sentence after its summary line: `Generated specification: checked`,
`not checked (source errors prevent building it)`, `not checked (subtree check)`,
`not checked (focused check)` or `check incomplete (operational error)`. Read the
findings and exit status to determine whether a checked spec passed.

## Adopt quality attributes

`new quality-attribute <id>...` adopts each named attribute into the nearest system that encloses the working directory. Without `--from`, each attribute starts with a blank definition in `specs/quality/<id>/`. `--from <model>` copies each term and its missing parents from the bundled `iso-25010-2023` model, a local quality model slug, or a model source path; there is no default model. Read the bundled model's terms offline with `dotspec docs quality-models/iso-25010-2023`.

Adoption follows the [no-overwrite rule](../../specs/requirements/reject-existing-target.requirement.spec.md): when the system already declares a named attribute, it exits 2 and writes nothing. A parent the system already declares is reused, keeping its wording and grouping. `new quality-attribute --json` reports `schema_version`, `target`, `dry_run`, the IDs in `added`, and ordered `entries`. Each entry supplies its `id`, `path`, `role` (requested or ancestor), and `decision` (create or reuse). Creations include their complete proposed `content`. `--dry-run` describes the same decisions and preserves the filesystem. `list quality-model` lists the project's own models; the bundled model is read through `dotspec docs`. [Define a System quality attribute](../cli-handbook/guides/define-a-system-quality-attribute.md)
