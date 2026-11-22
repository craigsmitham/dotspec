# CLI design

A contributor running a check in a terminal and a build saving its JSON report both need the findings, even when a specification fails. The report goes to standard output so either caller can keep it. The exit status lets the build decide whether to continue. This applies the distinction between primary output and messages in the *Command Line Interface Guidelines* and quality.md’s CLI guide. For a dotspec check, the findings themselves are the requested output. [Command Line Interface Guidelines, “The basics”](https://clig.dev/#the-basics) · [quality.md, “The basics”](https://github.com/qualitymd/quality.md/blob/6291df0e3173bf281f10a2a446a50eed041b7d21/docs/guides/cli-design.md#the-basics) · [Command Invocation](../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#observe-the-result).

## Commands and input

For a new or changed command, start with the result its caller needs and the specification that governs the operation. Decide which arguments identify its subject and what the caller can learn from an empty result or a failure. Declare the arguments, flags, short description and a worked invocation in the Cobra definition under `internal/cli`. Supply the files that invocation needs in the [command example test](../../internal/cli/contracts_test.go).

Keep the tree organized by singular concept and action, as in `job show` and `requirement check`.
Project-wide `check` and `sources` are direct operations over the selected scope.
Each leaf operation validates its arguments before filesystem access. `--path`
selects an existing project directory; flags may precede or follow positional
arguments. `--` ends flag parsing. Commands must work without interactive prompts.

For commands with a format choice, use `--json` as the convenience spelling of `--format json`. Passing both is a usage error. For a command that already returns an artifact, such as `model export`, write that artifact directly so a caller can redirect it to a file.

## Output and errors

- Primary results go to stdout. Check findings, including incomplete reports,
  are primary results. Invocation errors and ordinary command failures go to stderr.
- A completed check of a draft can suggest adding a performer and exit `0`. Authors can keep working with such a draft. A broken local link warrants exit `1` because the checker found an error in the material. If the checker could not read the file, it exits `2` because part of the requested inspection could not take place. Invalid invocations and output failures also use `2`. [Advice permits a successful check](../../specs/requirements/advice-permits-success.requirement.spec.md) · [Inspection report contract](../reference/inspection.md#check-specifications).
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

From the root of a checkout with the pinned tools installed, run `mise exec -- go test ./internal/cli -run TestDeclaredCommandExamples`. If your example fails, compare its flags with the command definition and its input paths with the test fixture. Fix the invocation or supply the missing fixture, then rerun.

For a change in process behavior, extend the relevant script in `internal/integration/testdata/`. Assert the exact exit status and both output streams. For example, the diagnostics script checks three different outcomes for a draft, a malformed specification and an unavailable source. Run it with `mise exec -- go test ./internal/integration -run TestCLI/diagnostics`. For another script, substitute its filename without the extension. [Diagnostics scenario](../../internal/integration/testdata/diagnostics.txtar).

Test a changed JSON result by parsing its fields. Include empty results, output failure, flag placement or current-directory behavior when they affect the change. Regenerate and review the public reference and command output examples from the checkout with `mise run generate`. [Contributor workflow](../../CONTRIBUTING.md#cli-changes-and-generated-artifacts).

Keep parser and checker logic independent of command presentation. Add a dependency
only for an actual boundary or capability; shared document extraction stays in
Goldmark and type-specific meaning stays in the type packages.
