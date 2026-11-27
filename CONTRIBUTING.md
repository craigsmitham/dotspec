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

`mise tasks` lists every task with its description, which [mise.toml](mise.toml)
owns. The `depends` of its `check` task are the steps of the full gate. Pass CLI
arguments after `--`, as in `mise run run -- list job`. [Releasing](docs/maintainers/releasing.md)
covers `snapshot`.

`check` is read-only. It fails on generated-file drift and asks you to regenerate;
it does not stage changes or install Git hooks. Run focused tests while iterating,
then the full gate before submitting a change.

## Documentation

This repository is the source of truth for the public docs at `dotspec.md`. Mintlify publishes the [`mintlify/` directory](mintlify/) from `main`. Whether you use a local editor or Mintlify’s editor, land the change in this repository to publish it.

| What are you changing? | Where to edit |
| --- | --- |
| Website landing page, quickstart, or navigation | Edit the MDX page under [`mintlify/`](mintlify/) or the navigation in [`docs.json`](mintlify/docs.json). |
| How to write a specification | Edit the [CLI handbook](docs/cli-handbook/README.md). The CLI ships these guides offline. Run `mise run generate` to refresh the selected guides in `mintlify/guides/`. The [web handbook](mintlify/handbook.mdx) links to both native pages and additional topics in the repository. |
| Commands, flags, or command output examples | Edit the Cobra command or example and run `mise run generate` to refresh the [CLI reference](docs/reference/cli.md) and [command examples](docs/reference/command-examples.md), along with the selected pages in `mintlify/cli/`. |
| A worked example | Edit the [To-do project](examples/todo/README.md), rebuild its report if the specifications changed, and run `mise run example-check`. Link readers to the relevant source, report section, or test. |

The [website generator](scripts/generate/web.go) selects the published guides and commands. To publish another, add it there and to the [site navigation](mintlify/docs.json), then run `mise run generate`. CI checks the generated files against their sources with `mise run generated-check`.

For website changes, run the same checks as the web docs workflow from the repository root.

```sh
cd mintlify
npm exec --yes --package=mint@4.3.97 -- mint validate
npm exec --yes --package=mint@4.3.97 -- mint broken-links --check-anchors
```

Use the [documentation map](docs/README.md) to find other topics. The [handbook maintenance guide](docs/maintainers/handbook.md#changing-a-reference) lists the inputs and regeneration commands for generated references and examples.

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
  It decodes `dotspec.yaml` strictly, without parsing any domain content.
- `internal/speccollection` implements owner-scoped lookup, subtree listing and the check type for business rules, concepts, jobs, requirements and use cases. Each family supplies its layout, reader, writing advice, group result and any enrichment, such as classifying a requirement against its associated system. Reusable quality models live in `internal/qualitymodel`, which also embeds the bundled ISO/IEC 25010 model.
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
| Reorder a system report | Move the section in `systemOutline` in `internal/specdoc/overview_layout.go`. Move Scope with its Out of scope subsection. The planner uses the outline to place passages, omit empty sections, allocate anchors and build navigation. Components have their own outline. | The To-do reader journey in `scripts/check-example/check.py` names the expected heading order, nesting and destinations. Change those expectations when the intended reading order changes. |
| Accept a section heading | `specfile.SectionSlot` in `internal/specfile/profile.go` owns the vocabulary. A heading introduces the slot it names in sentence case, so `Out of scope` introduces `out_of_scope`, when its class declares that slot and `sectionSlots` lists it as authored in Markdown. `irregularHeadings` names the titles that depart from the convention and `headingOwners` restricts titles to the classes that author them. Export maps a rich-text section through the same rule; give a structured section, such as a collection or a link-bearing field, its own arm in the kind's mapper in `internal/model`. | `internal/specfile/testdata/section-headings.tsv` pins every accepted class and heading. `TestCheckAndExportAcceptTheSameHeadings` in `internal/model` fails when check accepts a heading that export does not map. |
| Add an allowed relationship endpoint | Edit the slot’s `range` or `any_of` in `model/metamodel.yaml` and run `mise run model-generate`. Export and offline validation consume the generated ranges. If a core diagram shows that relationship, include the new endpoint in its node selection. | Use a named accepted target and a named wrong-kind target in an authored example. For scenario steps, see `TestIllustrationIsExplicitAndRejectsWrongEndpoint` and the To-do example’s check that a step illustrates the availability requirement. |
| Revise the invalid-slug message | `specfile.Layout.ValidateSlug` owns the collection slug grammar and repair text. `Layout.Create` and `Layout.LocalPath` apply it when a file is created or looked up by slug. Each type supplies its optional example through `SlugExample`. `specfile.ReadMember` reports a malformed collection filename as `layout.<kind>-name`. Named systems and components derive filenames from their headings through `NameSlug`. | `diagnostics.txtar` checks the literal text, stream and process status for invalid `new job` and `new feature` slugs and a malformed job filename. |
| Change a finding’s severity or rule ID | Change the producer that detects the condition. `internal/check` aggregates located findings. `internal/cli` formats reports and maps failures to process status. Update the diagnostic reference when the public rule changes. | The process tests exercise complete scans with findings, incomplete scans and invalid invocations. |
| Change system grouping in a family list | `internal/speccollection` groups by system source path, sorts equally named systems by path, and puts unassociated members last. `speccollection.Members` gives `dotspec list <kind>` that order for every collection kind, and the CLI writes one list envelope for all of them. | The job store fixtures in `internal/job/store_test.go` include distinct systems with the same name, unassociated members and nested scopes. `internal/speccollection/collection_test.go` covers listing from a subtree. |

### Follow an authored link into a report

The **Illustrates** field in the To-do use case links a scenario step to the availability requirement. To follow that link through the CLI, start with the [authored step](examples/todo/specs/use-cases/resume-list.use-case.spec.md#main-success-scenario).

1. `check.InspectSources` discovers the selected layouts, acquires the Markdown in a `specfile.SourceSet`, and runs the use-case inspector. The shared Markdown parser preserves source positions. The use-case parser interprets the scenario structure. The checker follows direct references within the repository and records the target file in its inventory.
2. `model.ExportInspection` takes that inventory, including the type registrations used during discovery. The mapper associates the link occurrence and its source span with the scenario step. Resolution identifies the target, and the generated slot contract determines whether the step may illustrate it. A saved snapshot carries the authored evidence for offline validation.
3. `specdoc.Generate` plans pages and fragment destinations from the snapshot. `internal/specdoc/relationship_layout.go` supplies reader-facing labels and placement for each direction. `specdoc.Write` reconciles generated files, and `specdoc.Check` inspects the delivered Markdown, including its links and attribution.
4. The To-do example reads exported JSON and follows links in the built reports. Keep the expected target, heading order and source links in the example expressed in the project's own terms. Its acceptance procedure exercises the running app across a restart. [Example verification](scripts/check-example/check.py) · [App acceptance procedure](examples/todo/tests/acceptance.py).

When a context provider reads Markdown during inspection, pass the inspection’s source set through to the type reader. This keeps context checks and export on the same observed file revision. A new inspection acquires subsequent edits. Check records each source’s context resolution in `Inventory.Contexts`, and export maps those results rather than calling the provider again, so both apply one selection rule.

### Adding a specification type

Ask two questions before adding a specification type. **What product intent does this subject own? Why must an author maintain it independently?** An independently authored reference, such as a quality model, has its own document type and filename suffix. Reserve `*.spec.md` and `SpecificationType` for roots that meet both tests. [What belongs in a specification?](docs/cli-handbook/explanations/what-belongs-in-a-specification.md)

For example, a new concept action already has an editable home and an address within its concept. A new specification root needs a concrete specification class and a Markdown mapping that explains its declaration. Add its frontmatter value to `SpecificationType` in the metamodel. [Specification authority](docs/maintainers/plans/specification-authority.md#what-each-source-owns) · [Model maintenance](model/README.md#regeneration).

Define a `specfile.Layout` in the owning type package with its specification type
(the `Kind` field) and its suffix. For a collection such as jobs, provide the default collection directory. For a named directory owner, see the system and component declarations in `internal/specfile/subject.go`. Use that layout for discovery,
placement validation, creation, and direct lookup so those operations agree.
For a collection with user-supplied slugs, `Layout.Create` and `Layout.LocalPath` enforce the slug grammar from `Layout.ValidateSlug`. `Layout.Create` takes the working directory, checks the slug before it touches the filesystem, then resolves the directory, so a type's `Init` only chooses its template. Set `SlugExample` when the repair text should illustrate a valid name. Read each file through `specfile.ReadMember`, which reports a malformed filename as `layout.<kind>-name` at the file's path, then set the slug and path on the result.

Define a metadata struct embedding `specfile.Metadata` with `yaml:",inline"` and
pass a fresh pointer to `(*specfile.Markdown).TypedOutline`. The shared decoder rejects fields outside that schema. Interpret
sections in the type package; register any explicit IDs with `specfile.Identifiers`
before validating local references with `(*specfile.Outline).ValidateReferences`. Preserve authored Markdown and source locations.

Expose a `CheckType` integration supplying the layout and an
inspector that reads the file and returns writing advice only after structural
validation succeeds. `check.ReadType` builds the inspector and `References` from
the type's reader. Attach a context provider only when the type needs other
sources. Types built with `speccollection.CheckType`, and the product and
component types, attach `system.CheckContext`; other types need not
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

Connect the kind to its class in `internal/specfile/profile.go` so its headings can be checked against the model. Register the checker in `All()` in `internal/speckinds/speckinds.go`; `check` uses it for project and focused checks. Add one entry to the kind table in `internal/cli/kinds.go` with the kind's domain description, `new` operand and flags, creator and lister; `dotspec new <kind>` and `dotspec list <kind>` are generated from it, and `TestKindTableContracts` covers the new entry. Map the parsed subject and its fields in `internal/model/mapping.go`; a guard test fails when a registered kind lacks a root class, a domain mapping or an authored suffix that matches its layout. Give its report a family directory in `families` in the document plan; the report's filename follows from the authored suffix and its type from the model's report class. Add its concrete class to the specification-type tests. [Document planning](internal/specdoc/plan.go) · [Specification-type checks](internal/model/specification_types_test.go).
`internal/check/testdata/quality` is a minimal integration fixture demonstrating
another type with its own metadata and no dependency on existing domain packages;
its deliberately different metadata is test-only. It is not a supported
specification type.

From the repository root, run `mise exec -- go test ./internal/check -run IndependentType` to exercise the fixture’s discovery, direct inspection and advice. For a public type, follow a real authored example from creation and inspection through snapshot export and the generated specification. Check that its source yields one primary record and that its report page leads back to it. Add its authoring topic and diagnostic destinations, exercise the public commands in the example gate, regenerate, and run the full check. [Snapshot export criteria](specs/requirements/model-source-fidelity.requirement.spec.md#criterion) · [Handbook publication](docs/maintainers/handbook.md#distribute-the-same-guides).

## CLI changes and generated artifacts

Follow [CLI design](docs/maintainers/cli-design.md). Declare commands, flags, and examples in
the Cobra tree. `scripts/generate` reads that tree without executing it for the
reference; its command output examples run real commands in a temporary project. Templates
come from the type packages. Edit those sources, then run `mise run generate`.

## Specification tests and implementation tests

This repository keeps each product decision in the specification for its subject. When adding or changing a CLI example, identify the governing requirement, concept action, use case, or native interface contract. If the example introduces a product requirement, author its requirement file in the same change.

- The scripts in `internal/integration/testdata/*.txtar` own their example inputs, commands and expected observations. They verify the public behavior of the specifications that list them as companions, or of the sources paired with them in Command Invocation.
- `internal/cli/quality_contracts_test.go` checks the conditions and criteria in the quality requirements that list it as a companion.
- Implementation tests can exercise private parser decisions and regressions. Changing the implementation can change these checks while preserving the specified behavior.

Review changes to example assertions against their governing source. Update that source when the intended behavior changes.

A specification lists its designated tests as companions in its frontmatter. Use
`role: verifies` for a test that assesses the specification, or `role: illustrates`
when an executable example is itself part of the specification. When the file holds several
tests, name the one that applies in the fragment, such as `#TestInspectionQualityContracts`;
a txtar script can be listed whole. Add a short note saying what it covers. Keep any account of what a
test exercises under **Evidence**. For example, the
[repeatable-inspection requirement](specs/quality/reliability/repeatability/repeatable-inspection.requirement.spec.md)
lists its quality-contract test. Scripts that no specification lists yet are linked from the
[Command Invocation concept](internal/integration/specs/concepts/cli-interaction.concept.spec.md),
beside the sources they exercise. An entire directory link does not provide file coverage.

`mise run spec-sources-check` runs `scripts/check-spec-sources`, a repository-owned
consumer of `dotspec sources --json`. The inventory includes companion targets as
well as links in the body, so the script sees tests listed either way. It selects
the two sets above, requires a successful complete inventory, and fails with the
missing paths. The same task runs in `mise run check` locally and in CI. The script
checks file coverage, not assertion adequacy or whether a declared role or a link's
prose correctly describes the test. Adapt its selection rules in other projects;
dotspec has no built-in test-file policy.

Tests in `internal/integration/testdata/*.txtar` describe real subprocess workflows
and their files together. `status CODE COMMAND ...` asserts the exact process exit
code; `stdout`, `stderr`, and `cmp` verify results. Add scenarios there for changes
to path handling, pipelines, command composition, or error behavior. Keep parser
edge cases in package-local Go tests and seed relevant cases into the fuzz targets.

The command-example test executes examples declared on the Cobra tree. When the README quick start changes, try it in a new directory with the rebuilt binary on PATH and follow the link in the resulting `SPEC.md`.

The bundled ISO/IEC 25010:2023 model is authored Markdown under
`internal/qualitymodel/specs/quality-models/`. The CLI embeds that file and
publishes it offline as `dotspec docs quality-models/iso-25010-2023`. Edit
definitions in the source file, retain attribution, and run `mise run generate`. The `internal/qualitymodel/apply` package relocates Markdown references while copying definitions into system attribute files. Try a selective copy in a disposable project and inspect its source attribution. [Apply selected attributes](docs/cli-handbook/guides/define-a-system-quality-attribute.md).

Keep the bundled model's prompts original rather than importing quotations from
the standard. Bundle another model only when it is MIT-compatible and a worked
example or governing specification uses it.

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
