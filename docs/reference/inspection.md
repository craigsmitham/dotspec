# Checking and discovering specifications

## Check specifications

Before committing specification changes, check the selected project subtree:

```sh
dotspec check
dotspec check --path internal/parser
dotspec check --format json
dotspec job check develop-requirements
```

`check` discovers system, job, use-case, concept, requirement, business-rule, and
quality-model specifications using the same
Git ignore, repository-boundary, and symlink rules as `list`. Outside Git, ignore files do
not affect discovery. It reports the selected scope, supported types, discovery
mode, files examined, and whether the scan completed. An empty completed scan
succeeds and explicitly reports `No supported specifications found.`

Structural errors fail the check: misplaced files, invalid metadata, missing
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
[repeatability requirement](../../specs/requirements/repeatable-inspection.requirement.spec.md)
defines equivalence across repeated inspections.
The [diagnostic-severity requirement](../../specs/requirements/diagnostic-severity.requirement.spec.md)
governs text presentation of severity, and the
[source-preservation requirement](../../specs/requirements/inspection-preserves-source.requirement.spec.md)
governs the effect of inspection on source files.

Jobs, use cases, concepts, requirements, and business rules also need their enclosing
system inspected. Requirements resolve attributes only after that
system passes structural validation; an invalid system does not trigger
cascading undefined-attribute advice. As with
`job show`, `use-case show`, and `concept show`, this direct lookup can read an
ancestor system outside the selected subtree or one excluded from discovery by
an ignore rule. Such files are listed as additional context
and counted separately; shared findings appear once. Context lookup stops at
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
Use the exit status as well as the report; a partial report is not a clean check.
For example, a CI step can run `dotspec check --format json > check-report.json`.
The focused `job check <slug>`, `use-case check <slug>`,
`concept check <slug>`, `requirement check <slug>`, and `business-rule check <slug>` commands use the same report, diagnostic IDs, and
`0`/`1`/`2` exit behavior. Each inspects one selected
specification, its enclosing system, and their direct reference targets, with `scope.discovery: "direct"` and
`scope.types: ["job"]`, `["use-case"]`, `["concept"]`, `["requirement"]`, or `["business-rule"]`. A missing specification
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
| `use-case.goal`, `use-case.system-boundary`, `use-case.primary-actor`, `use-case.goal-level`, `use-case.success-guarantees`, `use-case.main-success-scenario`, `use-case.extensions`, `use-case.sequence` | Advice about missing use-case content or scenario order. |
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

The [CLI Interaction concept](../../internal/integration/specs/concepts/cli-interaction.concept.spec.md)
is a working example beside its executable acceptance scripts. The
[source-preservation quality requirement](../../specs/requirements/inspection-preserves-source.requirement.spec.md)
instead links to tests that check its prose criterion. dotspec exposes both
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
[coverage script](../../scripts/check-spec-sources/main.go) invokes this command and
requires links to the repository's designated acceptance and quality-verification
tests. It intentionally does not select every test file. See
[Specification tests and implementation tests](../../CONTRIBUTING.md#specification-tests-and-implementation-tests)
for this repository's policy and the CI task that enforces file coverage.

## Output and automation

The experimental [specification model](specification-model.md) exposes
source bytes, ordered document content, and named domain records together:

```sh
dotspec model export --path . > specification.json
dotspec model validate specification.json --json
dotspec model schema > specification.schema.json
```

`model export` emits JSON directly. Invalid authored documents retain their source
and diagnostics. Its exit status follows the authoring/operational findings;
`model validate` instead checks whether the snapshot faithfully satisfies its
schema, source, reference, ownership, and coverage invariants. A valid snapshot
can describe invalid source documents.

The model distinguishes empty fields, absent fields, undecomposed narrative,
unknown sections, and failed extraction. Unknown content remains available with
an explicit mapping limitation. `complete_for_profile` reports coverage under
the declared mapping profile; it does not establish requirements completeness,
agreement, or satisfaction. See the [canonical model](../../model/README.md) and
[mapping reference](../../model/markdown-mapping.md).

The system, job, use-case, concept, requirement, business-rule, and quality-model
commands (including `init`), as well as `sources`, accept `--json`
or `--format json`. Use one of these flags, not both. The default is `--format text`; `job map` also
accepts `--format mermaid` and emits standalone diagram source. Redirecting output
does not change its format. JSON, Mermaid, and redirected output contain no ANSI
styling. `NO_COLOR` disables terminal styling, including help and errors.

Every JSON result has `schema_version: 1`. Check results use the report documented
above. Other results contain:

| Command | Result fields |
| --- | --- |
| `system init`, `job init`, `use-case init`, `concept init`, `requirement init` | `path` of the created specification |
| `sources` | `scope`, `complete`, `files`, and `diagnostics` as described above |
| `system show` | `system` and its nullable `enclosing_system` (directory context) |
| `system list` | `systems`, a directory tree of `system` objects and `children` arrays |
| `job show` | `job`, including its authored Markdown, steps, outcomes, and nullable `system` |
| `job list` | `groups`, each with a nullable `system` and a `jobs` array |
| `use-case show` | `use_case`, including its authored Markdown, conventional section bodies, and nullable `system` |
| `use-case list` | `groups`, each with a nullable `system` and a `use_cases` array |
| `concept show` | `concept`, including its authored Markdown, conventional section bodies, optional type declarations, and nullable `system` |
| `concept list` | `groups`, each with a nullable `system` and a `concepts` array |
| `requirement show` | `requirement`, including its authored Markdown, statement, criterion, verification plan, attributes, and nullable `system` |
| `requirement list` | `groups`, each with a nullable `system` and a `requirements` array; `attribute` when `--attribute` selects one |
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
See the generated [CLI reference](cli.md) for flags and examples.

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
outcomes, system quality attributes, and obligation clauses given an `**ID:**`
under an `Obligations` section. This does not introduce IDs for arbitrary
paragraphs.

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

For explicit boundary, job-support and applicability links, see
[Declare relationships between specifications](../handbook/relationships.md).
The generated [relationship reference](model/relationships.md)
documents the metamodel's edge categories, cardinalities, scope and projections.
