# Contributing

When changing a domain concept, an authoring convention, or a quality check,
start with [Maintaining the handbook](docs/maintainers/handbook.md). It points
to the owning text, the sources behind each choice, and examples for review.

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
| `mise run generate` | Regenerate the guide map, command reference, output gallery, and the dotspec and To-do requirements specifications |
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

## Documentation

Use the [documentation map](docs/README.md) to find the source for a CLI guide topic. The [handbook maintenance guide](docs/maintainers/handbook.md#changing-a-reference) identifies the inputs and regeneration commands for published tables, starters and examples.

## Architecture

- `main.go` owns process exit. `internal/cli` defines Cobra commands and handles
  presentation, JSON encoding, Fang terminal help, Lip Gloss errors, and status mapping.
- `internal/check` scans explicitly supplied type integrations and produces typed
  reports and the source inventory. Both operations share discovery, context, and
  reference resolution. It owns report accounting and context deduplication, without importing
  the concrete specification types. Checks do not print or exit.
- `internal/qualitymodel`, `internal/businessrule`, `internal/system`, `internal/job`, `internal/usecase`, `internal/concept`
  and `internal/requirement` own type-specific parsing and validation,
  metadata schemas, layouts, templates, and repository operations. Their `CheckType`
  functions supply inspection behavior; the CLI registers the supported types.
  Parsing functions accept bytes without I/O.
- `internal/specfile` owns common frontmatter, Goldmark outlines, source locations,
  document-local identifiers, reference extraction, and shared filesystem/discovery
  behavior. All specification types use the same Markdown extraction path.
- `internal/speccollection` implements owner-scoped lookup and subtree listing for business rules, concepts, jobs, requirements and use cases. Each family supplies its reader, slug rules, group result and any enrichment, such as classifying a requirement against its associated system. Reusable quality models have their own selection policy, including bundled starters.
- `internal/diagnostic` owns the finding representation and source locations. The checker or parser that detects a condition supplies its rule ID and message.
- `model/metamodel.yaml` owns the experimental specification metamodel. Generated
  JSON Schema and the inherited range catalog are embedded by the `model` Go
  package. [Model maintenance](model/README.md) describes the pinned toolchain.
- `internal/model` records normalized Markdown and ordered document content, maps
  named domain fields, resolves bounded references, and validates snapshots.
  Schema-backed tagged records preserve subtype fields. The
  [Markdown mapping profile](model/markdown-mapping.md) defines deterministic
  extraction; generated catalog ranges own relationship endpoint restrictions.

Keep CLI dependencies out of the parser and checker. Preserve authored Markdown
and locations; semantic validation must not silently rewrite files. Test a new
specification construct with valid, invalid, and source-location examples. Preserve
unknown prose and distinguish draft omissions from structural errors.

### Where to change a behavior

| Change | Owner and consumers | Independent check |
| --- | --- | --- |
| Reorder a system report | `systemOutline` in `internal/specdoc/overview_layout.go` sets reading order and navigation groups. The planner uses that outline for passage placement, anchors and navigation. Components have their own outline. | The To-do reader journey in `scripts/check-example/check.py` spells out the expected group labels, heading order and link destinations. |
| Add an allowed endpoint to a relationship | Edit the slot’s reference ranges in `model/metamodel.yaml`, keep any `trace_targets` annotation in agreement, and regenerate the catalog. The exporter uses those ranges when it resolves an authored link. The Markdown mapping owns the heading that supplies the link. Report wording and the placement of each link direction live in `internal/specdoc/relationship_layout.go`. | Exercise both an accepted and a wrong-kind target. The offline validator checks typed occurrences against catalog ranges and the claims in the snapshot. |
| Change a diagnostic | Change the rule where the condition is detected. `internal/check` aggregates findings and computes inspection status. `internal/cli` formats the result and maps command failures to process status. Markdown model commands share their export, validation and output path in `newModelReadingCommand`. | The process tests check output streams and exit status. Update the diagnostic reference if the public rule changes. |
| Change system grouping in a family list | `internal/speccollection` groups by system source path, sorts equally named systems by path, and puts unassociated members last. Each family constructs its own JSON result. | The job and use-case store fixtures include distinct systems with the same name, unassociated members and nested scopes. |

For a generated report, follow the data from `internal/cli/spec.go`. The CLI supplies the type integrations to `check.InspectSources`. The inventory carries the observed source bytes into `model.ExportInspection`, which maps authored content and resolves references. `specdoc.Generate` plans pages from the snapshot. `specdoc.Write` reconciles generated files, and `specdoc.Check` inspects the delivered Markdown against source and navigation contracts. Offline snapshot validation can also run directly on exported JSON.

### Adding a specification type

Ask two questions before adding a specification type. **What product intent does this subject own? Why must an author maintain it independently?** An independently authored reference, such as a quality model, has its own document type and filename suffix. Reserve `*.spec.md` and `SpecificationType` for roots that meet both tests. [What belongs in a specification?](docs/cli-handbook/explanations/what-belongs-in-a-specification.md)

For example, a new concept action already has an editable home and an address within its concept. A new specification root needs a concrete specification class and a Markdown mapping that explains its declaration. Add its frontmatter value to `SpecificationType` in the metamodel. [Specification authority](docs/maintainers/plans/specification-authority.md#what-each-source-owns) · [Model maintenance](model/README.md#regeneration).

Define a `specfile.Layout` in the owning type package with its specification type
(the `Kind` field) and its filename or suffix. For a collection such as jobs, provide the default collection directory. For a named directory owner, see the system and component declarations in `internal/specfile/subject.go`. Use that layout for discovery,
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
sources. Business rules, jobs, use cases, concepts, and requirements share
`system.CheckContext`; other types need not
depend on system specifications. Providers report attempted source paths even
when reads fail, and the checker distinguishes structural findings from incomplete
scans. Sharing a provider instance resolves it once per owning directory. Its optional
`Inspect` function contributes findings from context files. A type can supply
`References` to expose validated identities and authored links, and
`InspectContext` to check relationships against those resolved sources after
context validation succeeds; requirements use this for attribute labels.

The checker follows one reference edge from selected specs and their system
context. Pass the supported target types to `check.Focused` so target parsers
can expose type-specific IDs. Targets contribute structural diagnostics and
context accounting, but do not recursively expand the scan.

Connect the kind to its class in `internal/specfile/profile.go` so its headings can be checked against the model. Register the checker in `supportedTypes()` in `internal/cli/check.go` and pass the same integration to any focused check command. Map the parsed subject and its fields in `internal/model/mapping.go`. Give its report an address in the document plan and add its concrete class to the specification-type tests. [Document planning](internal/specdoc/plan.go) · [Specification-type checks](internal/model/specification_types_test.go).
`internal/check/testdata/quality` is a minimal integration fixture demonstrating
another type with its own metadata and no dependency on existing domain packages;
its deliberately different metadata is test-only. It is not a supported
specification type.

From the repository root, run `mise exec -- go test ./internal/check -run IndependentType` to exercise the fixture’s discovery, direct inspection and advice. For a public type, follow a real authored example from creation and inspection through model export and generated report. Check that its source yields one primary record and that the report leads back to it. Add its authoring topic and diagnostic destinations, exercise the public commands in the example gate, regenerate, and run the full check. [Model export criteria](specs/requirements/model-source-fidelity.requirement.spec.md#criterion) · [Handbook publication](docs/maintainers/handbook.md#distribute-the-same-guides).

## CLI changes and generated artifacts

Follow [CLI design](docs/maintainers/cli-design.md). Declare commands, flags, and examples in
the Cobra tree. `scripts/generate` reads that tree without executing it for the
reference; its output gallery runs real commands in a temporary project. Templates
come from the type packages. Edit those sources, then run `mise run generate`.

## Specification tests and implementation tests

This repository keeps each product decision in the specification for its subject. When adding or changing a CLI example, identify the governing requirement, concept action, use case, or native interface contract. If the example introduces a product requirement, author its requirement file in the same change.

- The scripts in `internal/integration/testdata/*.txtar` own their example inputs, commands and expected observations. They verify the public behavior identified alongside their links in Command Invocation.
- `internal/cli/quality_contracts_test.go` checks the conditions and criteria in the quality requirements that cite it.
- Implementation tests can exercise private parser decisions and regressions. Changing the implementation can change these checks while preserving the specified behavior.

Review changes to example assertions against their governing source. Update that source when the intended behavior changes.

The [CLI Interaction concept](internal/integration/specs/concepts/cli-interaction.concept.spec.md)
links to the acceptance scripts. Quality requirements explain their verification links.
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

The command-example test executes examples declared on the Cobra tree. When the README quick start changes, try it in a new directory with the rebuilt binary on PATH and follow the link in the resulting `SPEC.md`.

The model generation pipeline also validates subset classifications, explicit class
memberships, and the focused documentation view. Edit `model/views/core-domain.yaml`
for presentation choices; keep types and relationship semantics in `model/metamodel.yaml`.

See [Releasing](docs/maintainers/releasing.md) for snapshot and publication procedures.

Bundled quality models are authored Markdown under
`internal/qualitymodel/specs/quality-models/`. The CLI embeds those files; the
standard generator produces their reference pages. Add or edit definitions in
the source files, retain attribution, and run `mise run generate`. The `internal/qualitymodel/apply` package relocates Markdown references while copying definitions into system attribute files. Try a selective copy in a disposable project and inspect its source attribution. [Apply selected attributes](docs/cli-handbook/guides/define-a-system-quality-attribute.md).

The Q42 starter is adapted from arc42's CC BY-SA 4.0 catalog at the revision named
in its Edition section. Preserve the per-attribute source and license notices so
selective copies retain attribution. Its generated reference has the same license.
Refresh names and dimension assignments against a pinned upstream revision; use
original prompts rather than importing third-party quotations from upstream pages.

## Verify the runnable example

Run `mise run example-check` from the repository root. The same read-only task
runs in CI through `mise run check`. It builds the CLI, checks the public
capability inventory, tests the app across restart and rejects deliberately
broken temporary copies. See the [To-do walkthrough](examples/todo/README.md).

To enable the tracked pre-push hook in this clone, run
`git config --local core.hooksPath .githooks`. This selects the repository's
hook directory for the clone. If you already use another hooks directory, add
`mise run spec-check` and `mise run example-check` to its pre-push hook. Git must be able to find
mise on PATH. To undo the local setting, restore its previous value or run
`git config --local --unset core.hooksPath` if it was previously unset.
