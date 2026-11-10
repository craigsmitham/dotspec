# Contributing

Install [mise](https://mise.jdx.dev), then run these commands from the repository root:

```sh
mise trust
mise install
mise run check
```

`mise.toml` pins Go, golangci-lint, and GoReleaser. Runtime and test dependencies
are pinned in `go.mod` and `go.sum`. Git is required for repository discovery and
integration tests. The same check task runs on Linux, macOS, and Windows in CI.

| Task | Purpose |
| --- | --- |
| `mise run build` | Build `dist/dotspec` for the current platform |
| `mise run run -- job list` | Run the CLI from source |
| `mise run fmt` | Format Go source |
| `mise run test` | Run unit tests and process-level testscript scenarios |
| `mise run check` | Check formatting, module tidiness, lint, tests, generated artifacts, specs, and release configuration |
| `mise run generate` | Regenerate command reference, output examples, and starter specs |
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
  reports. It owns report accounting and context deduplication, without importing
  the concrete specification types. Checks do not print or exit.
- `internal/system`, `internal/job`, `internal/usecase`, `internal/concept`, and
  `internal/quality` own type-specific parsing and validation,
  metadata schemas, layouts, templates, and repository operations. Their `CheckType`
  functions supply inspection behavior; the CLI registers the supported types.
  Parsing functions accept bytes without I/O.
- `internal/specfile` owns common frontmatter, Goldmark outlines, source locations,
  document-local identifiers, local references, and shared filesystem/discovery
  behavior. All specification types use the same Markdown extraction path.
- `internal/diagnostic` owns stable findings and source locations.

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
sources. Jobs, use cases, concepts, and quality requirements share
`system.CheckContext`; other types need not
depend on system specifications. Providers report attempted source paths even
when reads fail, and the checker distinguishes structural findings from incomplete
scans. Sharing a provider instance resolves it once per owning directory. Its optional
`Inspect` function contributes findings from context files. A type can supply
`InspectContext` to check relationships against those resolved sources after
context validation succeeds; quality requirements use this for attribute labels.

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

Tests in `internal/integration/testdata/*.txtar` describe real subprocess workflows
and their files together. `status CODE COMMAND ...` asserts the exact process exit
code; `stdout`, `stderr`, and `cmp` verify results. Add scenarios there for changes
to path handling, pipelines, command composition, or error behavior. Keep parser
edge cases in package-local Go tests and seed relevant cases into the fuzz targets.

Tests also execute the examples declared on the command tree and the quick-start
commands in the README, so documentation syntax cannot silently drift.

See [Releasing](docs/releasing.md) for snapshot and publication procedures.
