# Contributing

Install [mise](https://mise.jdx.dev), then run these commands from the repository root:

```sh
mise trust
mise install
mise run model-setup
mise run check
```

`mise.toml` pins Go, golangci-lint, GoReleaser, and the Python version used for
model generation. `model-setup` explicitly installs the locked LinkML toolchain
in the ignored `.venv-model/` environment. The CLI itself runs entirely in Go.
Runtime and test dependencies
are pinned in `go.mod` and `go.sum`. Git is required for repository discovery and
integration tests. The same check task runs on Linux, macOS, and Windows in CI.
Git attributes keep text files at LF line endings on every platform so formatting
and generated-file checks compare the same bytes.

| Task | Purpose |
| --- | --- |
| `mise run build` | Build `dist/dotspec` for the current platform |
| `mise run run -- job list` | Run the CLI from source |
| `mise run fmt` | Format Go source |
| `mise run test` | Run unit tests and process-level testscript scenarios |
| `mise run check` | Check formatting, module tidiness, lint, tests, generated artifacts, specs, specification-source coverage, and release configuration |
| `mise run spec-sources-check` | Check that designated specification-test files appear in the public source inventory |
| `mise run generate` | Regenerate command reference, output examples, and starter specs |
| `mise run model-setup` | Provision the isolated model toolchain from its dependency lock |
| `mise run model-generate` | Regenerate JSON Schema, reference catalog, and model documentation |
| `mise run model-check` | Verify generated model projections without writing |
| `mise run model-test` | Test model and catalog generation |
| `mise run fuzz` | Fuzz each Markdown parser for ten seconds; also runs in CI |
| `mise run snapshot` | Build the complete release matrix locally without publishing |

`check` is read-only. It fails on generated-file drift and asks you to regenerate;
it does not stage changes or install Git hooks. Run focused tests while iterating,
then the full gate before submitting a change. The standard golangci-lint set
includes errcheck, govet, ineffassign, staticcheck, and unused.

## Architecture

- `main.go` owns process exit. `internal/cli` defines Cobra commands and handles
  presentation, JSON encoding, Fang terminal help, Lip Gloss errors, and status mapping.
- `internal/check` scans explicitly supplied type integrations and produces typed
  reports and the source inventory. Both operations share discovery, context, and
  reference resolution. It owns report accounting and context deduplication, without importing
  the concrete specification types. Checks do not print or exit.
- `internal/businessrule`, `internal/system`, `internal/job`, `internal/usecase`, `internal/concept`,
  `internal/quality` and `internal/requirement` own type-specific parsing and validation,
  metadata schemas, layouts, templates, and repository operations. Their `CheckType`
  functions supply inspection behavior; the CLI registers the supported types.
  Parsing functions accept bytes without I/O.
- `internal/obligation` shares general and quality requirement content parsing
  and quality-attribute classification.
- `internal/specfile` owns common frontmatter, Goldmark outlines, source locations,
  document-local identifiers, reference extraction, and shared filesystem/discovery
  behavior. All specification types use the same Markdown extraction path.
- `internal/diagnostic` owns stable findings and source locations.
- `model/dotspec.yaml` owns the experimental specification metamodel. Generated
  JSON Schema and the inherited range catalog are embedded by the `model` Go
  package. [Model maintenance](model/README.md) describes the pinned toolchain.
- `internal/model` captures original bytes and ordered document content, maps
  named domain fields, resolves bounded references, and validates snapshots.
  Schema-backed tagged records preserve subtype fields. The
  [Markdown mapping profile](model/markdown-mapping.md) defines deterministic
  extraction; generated catalog ranges own relationship endpoint restrictions.

Keep CLI dependencies out of the parser and checker. Preserve authored Markdown
and locations; semantic validation must not silently rewrite files. Test a new
specification construct with valid, invalid, and source-location examples. Preserve
unknown prose and distinguish draft omissions from structural errors.

### Adding a specification type

Define a `specfile.Layout` in the owning type package with its kind, collection,
and either a singleton filename or a slug suffix. Use that layout for discovery,
placement validation, creation, and direct lookup so those operations agree.
Keep slug rules and filename diagnostics in the type package.

Define a metadata struct embedding `specfile.Metadata` with `yaml:",inline"` and
pass a fresh pointer to `specfile.ParseOutline` (or `Parse` for required-section
extraction). The shared decoder rejects fields outside that schema. Interpret
sections in the type package; register any explicit IDs with `specfile.Identifiers`
before validating local references. Preserve authored Markdown and source locations.

Expose a `CheckType` integration supplying the layout, slug validation, and an
inspector that reads the file and returns writing advice only after structural
validation succeeds. Attach a context provider only when the type needs other
sources. Jobs, use cases, concepts, general requirements, and quality requirements share
`system.CheckContext`; other types need not
depend on system specifications. Providers report attempted source paths even
when reads fail, and the checker distinguishes structural findings from incomplete
scans. Sharing a provider instance resolves it once per owning directory. Its optional
`Inspect` function contributes findings from context files. A type can supply
`References` to expose validated identities and authored links, and
`InspectContext` to check relationships against those resolved sources after
context validation succeeds; general and quality requirements use this for attribute labels.

The checker follows one reference edge from selected specs and their system
context. Pass the supported target types to `check.Focused` so target parsers
can expose type-specific IDs. Targets contribute structural diagnostics and
context accounting, but do not recursively expand the scan.

Register the type in the CLI's project check and pass the same integration to any
focused check command. Preserve the existing command and report contracts.
`internal/check/testdata/quality` is a minimal integration fixture demonstrating
another type with its own metadata and no dependency on existing domain packages;
its deliberately different metadata is test-only. The supported quality
format belongs to `internal/quality`.

## CLI changes and generated artifacts

Follow [CLI design](docs/cli-design.md). Declare commands, flags, and examples in
the Cobra tree. `scripts/generate` reads that tree without executing it for the
reference; its output gallery runs real commands in a temporary project. Templates
come from the type packages. Edit those sources, then run `mise run generate`.

## Specification tests and implementation tests

This repository adopts the [design principle](README.md#design-principle-keep-meaning-in-the-source-keep-connections-simple)
by linking to native sources and enforcing file coverage with a small consumer
script. Test roles are an authoring decision, not a classification inferred by
dotspec or a consequence of using a particular test framework.

- `internal/integration/testdata/*.txtar` are designated executable specification
  examples. Their assertions define required CLI outcomes for the illustrated
  inputs and conditions. Sequences are representative unless a required ordering
  is explicitly stated. New cases belong here only when reviewers intend to
  preserve their behavior across changes in implementation.
- `internal/cli/quality_contracts_test.go` verifies the quality obligations
  stated in the linked Markdown quality specs. Those prose criteria retain
  authority; a test passing is evidence for its exercised cases.
- Other Go tests and the integration harness remain development artifacts by
  default. They can test private structure, parser mechanics, and regressions
  without making every expectation a product obligation. A unit test may be
  deliberately adopted as a specification test; an integration test may instead
  be implementation-specific. Choose by intent and scope, not test level.

When a new test only freezes incidental behavior, keep it outside the designated
collection. When a test file mixes normative examples and implementation checks,
separate those concerns before designating the file. Review changes to designated
assertions as specification changes, including fixtures that determine expected
results. Existing general rules and executable examples can constrain the same
behavior; resolve disagreements explicitly rather than treating a green run as
precedence.

The [CLI Interaction concept](internal/integration/specs/concepts/cli-interaction.concept.spec.md)
links to the acceptance scripts. Quality specs explain their verification links.
Use ordinary file links from the relevant specs, with prose stating the role.
An entire directory link does not provide file coverage.

`mise run spec-sources-check` runs `scripts/check-spec-sources`, a repository-owned
consumer of `dotspec sources --json`. It selects the two sets above, requires a
successful complete inventory, and fails with the missing paths. The same task
runs in `mise run check` locally and in CI. The script checks file coverage, not
assertion adequacy or whether a link's prose correctly explains its role. Adapt
its selection rules in other projects; dotspec has no built-in test-file policy.

Tests in `internal/integration/testdata/*.txtar` describe real subprocess workflows
and their files together. `status CODE COMMAND ...` asserts the exact process exit
code; `stdout`, `stderr`, and `cmp` verify results. Add scenarios there for changes
to path handling, pipelines, command composition, or error behavior. Keep parser
edge cases in package-local Go tests and seed relevant cases into the fuzz targets.

Tests also execute the examples declared on the command tree and the quick-start
commands in the README, so documentation syntax cannot silently drift.

The model generation pipeline also validates subset classifications, explicit class
memberships, and the focused documentation view. Edit `model/views/core-domain.yaml`
for presentation choices; keep types and relationship semantics in `model/dotspec.yaml`.

See [Releasing](docs/releasing.md) for snapshot and publication procedures.
