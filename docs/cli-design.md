# CLI design

Dotspec serves people reading and editing Markdown, as well as programs invoking
the CLI. These conventions adapt [Command Line Interface Guidelines](https://clig.dev/)
and the lessons from quality.md.

## Commands and input

Use Cobra command definitions as the source of command names, flags, help, and
examples. Keep the tree organized by singular concept (`system`, `job`, `use-case`, `concept`, `requirement`, `quality`) and action.
Each leaf operation validates its arguments before filesystem access. `--path`
selects an existing project directory; flags may precede or follow positional
arguments. `--` ends flag parsing. Commands must work without interactive prompts.

Every new command gets a short description and at least one working example.
Machine-readable formats are explicit, never inferred from pipes. `--json` is the
convenience spelling of `--format json`; using both is a usage error. Mermaid and
other standalone artifacts remain directly redirectable.

## Output and errors

- Primary results go to stdout. Check findings, including incomplete reports,
  are primary results. Invocation errors and ordinary command failures go to stderr.
- Exit `0` means success; advisory findings do not fail a check. Exit `1` means
  specification findings. Exit `2` means invalid invocation or an operational
  failure, including missing inputs, discovery errors, and output failures.
- Format results from typed data. JSON has an explicit schema version and stable
  snake-case keys. Keep arrays empty instead of null. Missing associations may be null.
- Text paths are relative when useful; JSON paths are absolute. Include source
  positions and stable rule IDs when available. Never invent source positions.
- Fang styles terminal help; Lip Gloss styles errors without probing stdout.
  Pipes and `NO_COLOR` use plain Cobra help.
  Neither styling nor terminal detection may change specification data or status.
- A write command reports what it created. `init` must not overwrite an existing
  specification or follow symlinked metadata directories.

## Verification

Test stdout, stderr, and exact exit status together. Test successful empty results,
malformed input, incomplete scans, output failure, flag placement, and current-directory
behavior. Check JSON semantics separately from text snapshots. The generated output
examples and testscript golden files make changes to the human interface reviewable.

Keep parser and checker logic independent of command presentation. Add a dependency
only for an actual boundary or capability; shared document extraction stays in
Goldmark and type-specific meaning stays in the type packages.
