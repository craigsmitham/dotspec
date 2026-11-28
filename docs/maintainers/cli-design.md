# CLI design

A contributor running a check in a terminal and a build saving its JSON report both need the findings, even when a specification fails. The report goes to standard output so either caller can keep it. The exit status lets the build decide whether to continue. This applies the distinction between primary output and messages in the *Command Line Interface Guidelines* and quality.md’s CLI guide. For a dotspec check, the findings themselves are the requested output. [Command Line Interface Guidelines, “The basics”](https://clig.dev/#the-basics) · [quality.md, “The basics”](https://github.com/qualitymd/quality.md/blob/6291df0e3173bf281f10a2a446a50eed041b7d21/docs/guides/cli-design.md#the-basics) · [Command Invocation](../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#observe-the-result).

## Commands and input

For a new or changed command, start with the result its caller needs and the specification that governs the operation. Decide which arguments identify its subject and what the caller can learn from an empty result or a failure. Declare the arguments, flags, short description and a worked invocation in the Cobra definition under `internal/cli`. Supply the files that invocation needs in the [command example test](../../internal/cli/contracts_test.go). In usage lines, write `<name>` for a required operand, `[name]` or `[<name>]` for an optional one and `...` for repetition, all in lowercase. Write flag help as a lowercase phrase with no closing period. Add a flag that takes a value with `stringFlag`, which names the value in help, as in `--part-of <path>` or `--template brief|expanded`, instead of pflag's generic `string`; read it with `flagValue`. The same test file checks all three conventions. Validate a command group's arguments with `subcommands` and a leaf's with `operands` in `internal/cli/usage.go`: a misspelled subcommand then suggests the nearest one and exits `2`, even when `--help` follows it, and a missing operand or required flag is named by its placeholder, as in `new requirement: missing <slug>` or `new component: missing --part-of <path>`. Write examples without indentation; the tree indents them for plain help.

Write a Long description in one shape: what the command does in domain terms; for `check`, `sources` and snapshot validation, one sentence on what the result establishes; and a closing `Learn more: dotspec docs <address>` line whose address opens a handbook topic or page. State the exit-code policy only in root help. On a leaf that can exit `1`, say so in one sentence.

Keep the grammar verb-first. The verbs are the same for every kind, so `new` and
`list` take the kind as a subcommand, as in `dotspec new requirement fast-save` and
`dotspec list use-case`, the way `kubectl create <type>` does. Both are generated
from the kind table in `internal/cli/kinds.go`, keyed by `specfile.Layout`: each
entry supplies the kind's domain description, plural, `new` operand, kind-specific
flags, creation and listing. Adding a kind adds one table entry and no command
wiring. Leaf Shorts come from the table ("Create a requirement specification",
"List requirements by system"), and `new --help`, `list --help` and the error for
a missing kind list every kind with its description and `dotspec docs <kind>`
route. A kind-specific flag, such as `--part-of`, `--template`, `--attribute` or
`--from`, lives on that kind's subcommand, never on a generic verb. Named subjects
(`system`, `component`) take a `<name>` that becomes the heading, collection kinds
take a `<slug>`, and quality attributes take catalog `<id>...`; Part of is always
an explicit, required `--part-of`. Every list returns one envelope: `schema_version`
first, then `entries` with `kind`, `slug`, `name`, absolute `path` and nullable
`system`. `check`, `sources` and `snapshot export` act on files rather than kinds,
so they take a path.
Root help groups the commands as Author (`init`, `new`), Inspect (`check`, `list`,
`sources`), Publish (`build`, `snapshot`) and Help (`docs`, `help`, `completion`).
A command ships only when a use case, requirement or concept in `specs/` governs it.
Project-wide `check` and `sources` are direct operations over the enclosing project.
The project lifecycle sits at the root too: `init` starts a project in the working
directory, and `build` writes its generated specification. Call that artifact the
*generated specification* in help, messages and docs; *report* names check reports
and the per-subject pages inside it.
Each leaf operation validates its arguments before filesystem access. The global
`-C, --chdir <dir>` is the only location flag: it runs dotspec as if started in
`<dir>`, and every relative path, positional or flag, resolves against the
effective working directory. Reads default to the enclosing project; `check`,
`sources` and `snapshot export` narrow it with one optional `[<path>]` operand,
never a flag. `--all` means only "every discovered project" and exists only on
`check`. Declare where a command starts as a typed `scope` when the command is
built (`internal/cli/workspace.go`); never choose it by matching the command's
name. Flags may precede or follow positional arguments. `--` ends flag parsing.
Commands must work without interactive prompts.

`--json` is the only output switch. Give each command that returns a result a local `--json` flag with `addJSON`; never make it persistent, and leave it off `docs`, `help`, `completion` and `man`, which have no result to encode, and off `build`, which writes files. The [contract test](../../internal/cli/contracts_test.go) enforces both. For a command that already returns an artifact, such as `snapshot export`, write that artifact directly so a caller can redirect it to a file.

## Output and errors

- Primary results go to stdout. Check findings, including incomplete reports,
  are primary results. Invocation errors and ordinary command failures go to stderr.
- Write an error as `dotspec: <what>; <fix>`. A usage error ends with the help to
  consult, as in `see dotspec check --help`. A filesystem error names the path and
  the problem in plain words, never the failed system call; `Run` prints paths near
  the working directory relative to it.
- A completed check of a draft can suggest adding a performer and exit `0`. Authors can keep working with such a draft. A broken local link warrants exit `1` because the checker found an error in the material. If the checker could not read the file, it exits `2` because part of the requested inspection could not take place. Invalid invocations and output failures also use `2`. [Advice permits a successful check](../../specs/requirements/advice-permits-success.requirement.spec.md) · [Inspection report contract](../reference/inspection.md#check-specifications).
- Format results from typed data. Every JSON command result starts with
  `"schema_version": 1`; the snapshot carries the model's `model_version` and
  `representation_version` instead. Bump a result family's version only after a release
  changes it. Keys are stable snake case. Keep arrays empty instead of null. Missing
  associations may be null. Every `location` has the same shape and carries `path`.
- Text paths are relative to the working directory, so any specification path
  printed by `list` or `sources` can be pasted into `check`; JSON paths are
  absolute. Include source positions and stable rule IDs when available. Never
  invent source positions.
- Text output is prose for a reader, not a dump of the JSON. Print no enum values
  such as `blocked-by-source`; state them in words. Omit an incidental field that
  has no value rather than printing `None`; where the missing value is the subject,
  say so in a sentence, as in "No associated system".
- A text report leads with what the reader acts on. `check` prints its findings,
  then one summary line with the scope, completeness and counts, then the context
  files. A blocked `build` leads with `Specification not written: N source errors`.
  `sources` names each project once, above the files it contains, and prints a
  named source or procedure with its declared purpose rather than as a bare
  reference.
- Fang styles terminal help; Lip Gloss styles errors without probing stdout.
  Pipes and `NO_COLOR` use plain Cobra help.
  Neither styling nor terminal detection may change specification data or status.
- A write command reports what it created. `init` and `new` must not overwrite an
  existing manifest, specification or adopted quality attribute, or follow symlinked
  metadata directories.

## Verification

From the root of a checkout with the pinned tools installed, run `mise exec -- go test ./internal/cli -run TestDeclaredCommandExamples`. If your example fails, compare its flags with the command definition and its input paths with the test fixture. Fix the invocation or supply the missing fixture, then rerun.

When you rename or remove a command or flag, run `mise exec -- go test ./internal/cli -run 'TestDocumentedInvocationsResolve|TestCommandReferenceFragmentsResolve'`. The [invocation reference test](../../internal/cli/references_test.go) resolves every `dotspec` invocation, or `go run .` from the checkout, written in code spans and shell blocks of the handbook, reference, maintainer docs (except plans), model and skill pages, specifications, example walkthroughs, README and Contributing. It also resolves every command named in a Go string after a cue such as "run" or "use"; write `dotspec` before a command named in help or a diagnostic so the test can find it. It also checks each link to a `cli.md#dotspec-*` heading. Update each reported reference to the current command.

For a change in process behavior, extend the relevant script in `internal/integration/testdata/`. Assert the exact exit status and both output streams. For example, the diagnostics script checks three different outcomes for a draft, a malformed specification and an unavailable source. Run it with `mise exec -- go test ./internal/integration -run TestCLI/diagnostics`. For another script, substitute its filename without the extension. [Diagnostics scenario](../../internal/integration/testdata/diagnostics.txtar).

Test a changed JSON result by parsing its fields. Include empty results, output failure, flag placement or current-directory behavior when they affect the change. Regenerate and review the public reference and command output examples from the checkout with `mise run generate`. [Contributor workflow](../../CONTRIBUTING.md#cli-changes-and-generated-artifacts).

Keep parser and checker logic independent of command presentation. Add a dependency
only for an actual boundary or capability; shared document extraction stays in
Goldmark and type-specific meaning stays in the type packages.
