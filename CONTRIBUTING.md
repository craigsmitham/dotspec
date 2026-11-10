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
- `internal/check` produces typed reports. Checks do not print or exit.
- `internal/system` and `internal/job` own type-specific parsing and validation,
  templates, and repository operations. Parsing functions accept bytes without I/O.
- `internal/specfile` owns common frontmatter, Goldmark outlines, source locations,
  local references, and shared filesystem/discovery behavior. Both specification
  types use the same Markdown extraction path.
- `internal/diagnostic` owns stable findings and source locations.

Keep CLI dependencies out of the parser and checker. Preserve authored Markdown
and locations; semantic validation must not silently rewrite files. Test a new
specification construct with valid, invalid, and source-location examples. Preserve
unknown prose and distinguish draft omissions from structural errors.

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
