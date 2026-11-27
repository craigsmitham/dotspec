# dotspec CLI design review: recommendations for the first release

## Status

Implemented on the `cli-v1-surface` branch, following the sequencing in section 9. The owner adopted the recommendation for every open decision in section 8: verb-first `new <kind>` and `list <kind>` generated from one kind table; `show` deferred, with no System line added to `check`; adoption as `new quality-attribute <id>... [--from <model>] [--dry-run]`, with only ISO 25010 bundled and published as `dotspec docs quality-models/iso-25010-2023`; `-C, --chdir` in place of `--path` and `--project`, with `--all` only on `check`; `model` renamed `snapshot`; `report.enabled` renamed `build.enabled`; and `relationship.motivation` deferred. The rest of this document is the design record as reviewed.

Deliberate deviations from the text below:

- **Learn more addresses.** No `check` or `inspection` topics were added to `topics.json`. The `check`, `sources` and `build` help points at exact `docs/reference/inspection.md#...` addresses instead, because the topic catalog is tested to match the authored document types.
- **`build` has no `--json`.** It writes the generated specification rather than returning a result, so the `--json` placement test exempts it.
- **Moved advisories.** They are `assumption.basis`, `assumption.applicability` and `interface.contract`. A focused check reports them only for the files it checks, not for a system it reads as context. The views' "no named interfaces recorded" prompt and `model review`'s shared-term-name prompt were removed with the views rather than moved.
- **`relationship.motivation`.** It was produced by the snapshot export as well as `model trace`, so it was removed from the export too.
- **Job checks.** Removing presets also removed the four advisories only presets produced: `job.boundary`, `job.outcome-measure`, `job.story-context` and `job.map`. The generated job page says "No steps recorded." when a job declares no steps.
- **Product composition.** `Snapshot.Product` lost the evidence and verification fields that only `model product` used.
- **Focused checks.** `check <file>` checks a specification file within its own project, even when the working directory is in another project.
- **Help details.**
  - The generated CLI reference lists the global `-C` once, under the root.
  - `list system` and `list component` keep their own Shorts ("List systems by directory", "List components by Part of").
  - The root example runs `build` before `check`, so the check passes.
  - A value flag names its value in help, as in `--part-of <path>`. A missing required flag is reported like a missing operand. A misspelled command followed by `--help` exits 2 with a suggestion.
- **Shell completion** shows the generated leaf Shorts ("Create a job specification"), not the domain Shorts section 4 lists, because Cobra completes from `Short` and section 7 asks for generated Shorts.
- **`new system --part-of`** must name a system, matching the metamodel. A component parent is rejected.
- **Check text** names a whole-project scope as `project <name>`. Types and discovery mode appear only in JSON.
- **Draft headings.** `new <kind> <slug>` titles a collection draft from its slug (`fast-save` becomes "Fast save") instead of "New requirement". Each draft therefore gets a distinct heading and its own generated page name.
- **Manifest errors.** An unknown `dotspec.yaml` key is reported by name, with the keys the manifest accepts, instead of as a YAML decoding error that names a Go type.
- **Bundled quality model page.** `quality-models/<slug>` opens by address and search but is not listed in `topics.json`, because knowledge packaging reads catalog pages from repository paths. `CONTRIBUTING.md` bundles a further model only when it is MIT-compatible and a worked example or governing specification uses it.

## 1. Bottom line

The CLI's runtime conventions are sound. Exit codes, output streams, the no-overwrite rule, offline docs and JSON hygiene all work. The command tree is the problem. It has 19 root entries (not counting Cobra's automatic `help`), about 55 leaf commands, and 19 distinct flags besides `-h` and `-v`. Most of that size comes from three sources:

- **Per-kind verb groups.** The same `init/list/show/check` set is repeated on 9 kind groups, each built a little differently.
- **A second reporting pipeline.** Five Markdown views and `model product` under `model` predate `spec build` and now duplicate it.
- **Parallel paths.** Two quality-adoption commands, job presets, job map views, and `--format` alongside `--json`.

The scope flags are also incoherent:

- `--path` says its default is `.`, but leaving it out behaves differently from passing `--path .`.
- Relative paths resolve against three different bases.
- Scope defaults are chosen by matching command-name strings (`workspace.go:93-118`).

Core recommendation: adopt one verb-first grammar in which the kind is a generated subcommand of the verb, as `kubectl create <type>`, `hugo new content` and `dotnet new <template>` do. Cut every command that no use case, requirement or concept in `specs/` governs. That leaves:

- 9 root entries plus `help`: `init`, `new`, `list`, `check`, `sources`, `build`, `snapshot`, `docs` and `completion`.
- 7 command flags plus a local `--json`.
- One global `-C, --chdir <dir>`.

Adding a kind then costs one table entry, not four hand-wired commands.

The cut removes roughly 2,000 lines of Go source, about 1,400 lines of tests and about 2,100 lines of bundled vocabulary. It removes no domain meaning. The generated specification already renders everything the removed views showed, apart from three advisories. Two of those move into `check`. The third, `relationship.motivation`, needs your decision (Open decision 7) because it would add 10 advisories to dotspec's own check. No release tag exists, so every break is free now and costs something later.

## 2. What's working (keep)

- **Exit codes.** 0 means success (advice doesn't fail), 1 means error findings, 2 means bad invocation or an operational failure. The policy is stated in root help (`cli.go:27`) and matches Vale's convention.
- **Streams.** Results go to stdout, including check reports with findings and blocked builds. Errors go to stderr with a `dotspec:` prefix.
- **Safe and scriptable.** Commands never prompt. `init` refuses to overwrite or to follow a symlink (exit 2). Flags can go anywhere, `--` works and `NO_COLOR` is respected.
- **JSON.**
  - Output is parseable, and empty arrays come back as `[]` rather than null.
  - `--dry-run` appears only on commands that write.
  - JSON paths are absolute (`cli-design.md:24`), and `scripts/check-spec-sources` relies on that.
- **Offline guidance.** `dotspec docs` addresses topics, `page#section`, rule IDs and words. Every finding prints `Guide: dotspec docs <rule>`. This is the explain-a-code pattern, grounded in the offline-authoring-guidance requirement.
- **`check` covers the generated spec.** A whole-project check also verifies that the generated specification is up to date ("Generated specification: checked"). A focused check reports it as `out-of-scope`.
- **Project isolation.** A nested `dotspec.yaml` bounds discovery and needs no flags.
- **Generated reference.** `docs/reference/cli.md` is generated from `NewRoot()`, and `TestDeclaredCommandExamples` executes every example. Keep both; they make the restructure cheap.
- **Domain-grounded Shorts.** For example, "Specify required results, rules, properties, and restrictions". Reuse them as the kind table's descriptions.
- **A kind registry (`speckinds`) to build on.** It currently drives only check and export.

## 3. Key findings

### A. The grammar repeats one lifecycle per kind, and every copy has drifted (high)

**Problem.** The kinds share one lifecycle, but each kind group is wired by hand, so the copies have drifted apart.

**Evidence.**
- `cli.go:198` hand-lists 11 constructors using five construction patterns:
  - Only concept, business-rule and use-case go through `specKind` (`speckind.go:120-170`).
  - Requirement, job and quality-model use bespoke `newSpecAction` loops.
  - Feature and capability use `product.go`.
  - System and component are hand-built.
- `component` has only `init`.
- `system` has no `check`.
- `system show app/specs/task-service.component.spec.md` fails with "expected a system source path".
- Feature and capability drift from the other kinds:
  - Their Shorts are generic ("Create a standalone specification", `product.go:20`).
  - Their lists are flat, with no system grouping.
  - They return `schema_version` 1 where the other kinds return 2.
- `<kind> check <slug>` is byte-identical to `check <file>`: `diff` is empty for both `requirement check restart` and `job check personal-commitments`.
- The two check forms take different arguments. `dotspec check restart` fails with "unknown source document type: restart" (exit 2), and `dotspec check specs` fails the same way.

**Recommendation.**
- Make the verbs top-level. Generate one kind subcommand under `new` and under `list` from a CLI-side kind table in `internal/cli`, keyed by `specfile.Layout`. Keep presentation out of `internal/check`, as `cli-design.md` requires.
- Make `check [<path>]` the only check command. It takes a file or a directory, like Vale, `redocly lint` and eslint.
- Delete every `<kind> check`.
- Rewrite `docs/maintainers/cli-design.md:7-9` ("organized by singular concept and action") to describe the new grammar.

### B. `model` holds a second reporting pipeline, and its name is ambiguous (high)

**Problem.** The `model` group mixes the snapshot contract with five Markdown views and a JSON product view. All of them duplicate the generated specification.

**Evidence.**
- `model.go:141-157` registers trace, overview, needs, review and pegs. Their renderers are `internal/model/{needs,overview,review,trace,project,design}_docs.go`, 1,215 lines in all. Only the CLI calls them; `spec build` uses `internal/specdoc`.
- They were added on 2026-11-11, before `specdoc` (2026-11-19).
- `review` and `pegs` emit the same "Environment and boundaries" blocks.
- The overview repeats the system report's Architecture, Features and Requirements sections. Trace's change-review links repeat the requirement report's "Refined by".
- No spec governs these views.
- `model product` admits it skips the export invariant guard (`model.go:106-108`). Its only consumers are `check.py:101` and `:532` and `product_test.go`. `Snapshot.Product` itself stays, because `specdoc/plan.go:308` uses it.
- "Model" means three things on the surface:
  - the `model` group;
  - `--model`, the quality catalog;
  - `quality-model`.

  Next to `check`, `model validate` reads like "validate my specs".

**Unique content that must move or be decided first.** Three advisories appear only in these views:
- `relationship.motivation` (`traceability.go:107`): a requirement with no supports, derived_from or refines. `check` doesn't emit it. `model trace` flags 10 of dotspec's own 11 requirements, while `check` reports 0 advisories today. It is a review prompt about justification, not an authoring defect. Whether it belongs in `check` is an owner decision (Open decision 7).
- Assumption with no Basis or Applicability (`review_docs.go:259-278`). I confirmed that `check` reports 0 advisories where `model pegs` prints "Review required". `describe-a-system.md:27` relies on this.
- Interface with no native contract reference (`review_docs.go:253`).

`requirement.measurement-basis` is already a check rule, so it needs no move.

**Recommendation.**
- Move the assumption and interface-contract advisories into `check` with stable rule IDs. Handle `relationship.motivation` according to Open decision 7.
- Delete the five views, `--relative-to` and `model product`.
- Keep `Snapshot.Product` internal. Move the `check.py` membership assertions into Go tests or onto the generated system report.
- Rename the remaining group to `snapshot export|schema|validate`. `Snapshot` is the metamodel class (`metamodel.yaml:310`), and the use case already calls the artifact "a faithful, structured snapshot".

### C. Scope and path rules are incoherent (high)

**Problem.** Scope depends on the command's name, relative paths have three bases, and project selection is wired into commands that never use it.

**Evidence.**

From `examples/todo/app`:

| Command | Result |
|---|---|
| `requirement list` | 4 items |
| `requirement list --path .` | 1 item |
| `check` | "Checking .." |
| `check --path .` | "Checking ." |
| `requirement show restart` | exit 2, although `list` printed `restart` |

The causes:
- `selectProject` switches on `cmd.Name()` (`workspace.go:96-103`). It gates `--all` on `cmd.Use == "check [source-path]"` (`:116`) and forces JSON when the name is `export` (`:153`).
- Relative paths have three bases:
  - `--part-of` resolves against the working directory (`system.go:166`).
  - `--system` resolves against `--path` (`qualityattribute.go:38`).
  - Positional paths also resolve against `--path`. `show examples/todo/app/... --path examples/todo` looks for `.../examples/todo/examples/todo/...`.
- `--project` is on 51 commands and `--all` on 18. Their only real uses are `check --all` (`mise.toml:61`) and `spec build --project` (`mise.toml:37`).
- `--all` means "all projects" in one place (`workspace.go:123`) and "all attributes" in another (`qualitymodel.go:127`).
- `model schema --workspace` and `workspace_model.go` (92 lines) exist only for the `model export --all` envelope.

**Recommendation.**
- Add one global `-C, --chdir <dir>` ("run as if started in DIR"), following `git -C`, `make -C` and `terraform -chdir`. pflag needs a long name for every flag, and rule 9 bans shorthand-only forms.
- Resolve every relative path, positional or flag, against the effective working directory.
- Default every read to the enclosing project.
- Narrow `check`, `sources` and `snapshot export` with one optional positional `[<path>]`.
  - `snapshot export [<path>]` keeps the subtree export that `export-specification-model.use-case.spec.md` step 1 requires ("selects a project subtree").
  - Use a single path rather than a variadic list. The report has one `scope.path`, `scripts/check-spec-sources/main.go:90` asserts that it equals the root, and repeatable-inspection compares these reports.
- Keep `--all` on `check` only, because the isolation requirement's "highest exit code" rule matters for CI. Document its envelope as check's workspace result.
- Remove `--path`, `--project`, the export envelope and `schema --workspace`.
- Amend `specification-project-isolation.requirement.spec.md` in two places:
  - Replace project selection by name with `-C`.
  - Narrow "inspect all discovered projects" and the Applies to ("source inventories, model exports and reading views") to `check --all`.
- Declare scope as typed data when each command is built, and delete the name switches.
- Deferring `show` (Open decision 2) removes slug lookup from v1, and with it the list/show scope mismatch.

### D. Duplicate and speculative authoring paths (high)

**Quality adoption.**

Problem:
- `quality-attribute add recoverability --dry-run --json` and `quality-model apply --attributes recoverability --dry-run --json` give identical output; both call `qualityOperation`.
- The two commands use mirror-image conventions: a positional ID with `--model`, against a positional model with `--attributes`.
- `quality-attribute browse` overlaps `quality-model show`.
- The `quality-attribute` noun acts on a catalog, while `*.quality-attribute.spec.md` files are the system's adopted attributes.

Recommendation:
- Adoption becomes creation of the kind: `new quality-attribute <id>... [--from <model>]`. This parallels the existing `quality-model init --from`. `--from` takes a required value, so there is no optional-value flag and no implicit default.
- It follows the no-overwrite requirement. Today, re-adding an existing attribute exits 0 with `"decision": "reuse"` (`quality-convention.txtar:49`). Under `new`, it exits 2 when a *named* attribute already exists, and existing parents may still be reused. Amend `reject-existing-target.requirement.spec.md`'s Applies to so that it covers `init` and `new`.
- `list` stays reserved for project files. `list quality-attribute` lists the adopted attributes, like every other kind.
- To read a catalog, publish the bundled ISO 25010 model as an offline docs page, for example `dotspec docs quality-models/iso-25010-2023`. Today `dotspec docs iso 25010` finds no bundled model. The alternative is to defer catalog browsing. Either way, `qualitymodel/store.go:164` stops pointing at `quality-attribute browse`.

**Presets and templates.**

Problem:
- Three names do the same job: `job --preset general|odi|job-story`, `requirement --preset general|quality` and `use-case --template brief|expanded`.
- `job check X --preset odi` reports advisories that `dotspec check` never reports. `inspection.go:10` hard-codes `general`, and the preset isn't recorded in the file.
- `job-story` differs from `general` only in prompt wording.
- The general requirement template invites a Criterion in prose: "Add a Criterion section if…" (`requirement/spec.go:29`). It has no Criterion heading, which the quality preset adds and `check.py:474` asserts. The preset therefore saves one heading. That isn't enough to justify a second flag and template.

Recommendation:
- Drop `--preset` everywhere, along with `presets.go`, the preset branches in `job/check.go:31-62` and `requirement.QualityTemplate`.
- Move the ODI and job-story prompts into `write-a-job.md` as prose.
- Keep `--template` only for use-case brief and expanded, which the create-specification extension "More detailed use-case prompts" governs.

**Job views.**

Problem:
- `job map --outcomes --json` and `job outcome list --json` give identical output.
- `job outcome` is the only three-level path in the tree, and Mermaid is the only reason `--format` exists.
- `job show` and the generated job page already present steps and outcomes in authored order.

Recommendation:
- Drop `job map`, `job outcome list`, `--outcomes` and `--step`.
- Revise `specs/concepts/job-map.concept.spec.md`. It has three actions, and a static page can meet only two:
  - Keep "Describe the job and its steps", met by authoring.
  - Keep "View the map", met by the generated job page, which should say "No steps recorded" explicitly.
  - Delete "Inspect a step's outcomes". It selects a step by ID and treats an unknown step as an error, which a static page can't do.
  - Remove "and selection" from State, the needs-map reference, and both `cli.md#dotspec-job-map` and `#dotspec-job-outcome-list` anchors.
- Rewrite `write-a-job.md:86-92`, which teaches `--step`.

**One-off list flags.**

Problem:
- `use-case list --overview` prints links relative to the source file, which don't resolve from the terminal (`../personal-task-list.system.spec.md#list-owner`).
- `requirement list --attribute` is the only list filter. The quality-attribute report already lists the classified requirements.

Recommendation: drop both.

**Bundled starters.**

Problem:
- FURPS is 269 lines and Q42 1,875, against 354 for the default ISO 25010 model. No worked example uses either.
- Q42 is CC BY-SA content inside an MIT binary.

Recommendation:
- Defer both.
- Replace Q42's multi-parent test coverage (`quality-convention.txtar:15-17`) with a local fixture.
- Delete the `quality: {model: q42}` build-stability test (`check.py:541-563`) together with the `quality.model` key.

### E. The lifecycle has a missing first step and an awkward last step (medium)

**Problem.** No command creates the project manifest, and the build command sits under an ambiguous one-child group.

**Evidence.**
- No command writes `dotspec.yaml`. `README.md:12` says to create it by hand.
- In an empty directory, `system init`, `requirement init` and `check` all succeed. Then `spec build` fails with "requires a project declared in dotspec.yaml" (exit 2), with no hint.
- `spec` is a group with one child, and every authored file is already a "spec".
- The artifact goes by four names:
  - "report" (`report.enabled`);
  - "generated specification" (the requirement title, the `generated_spec` JSON key and the check status line);
  - "SPEC.md";
  - "spec build".

**Recommendation.**
- Add `dotspec init <system-name> [--name <project>]`. It writes `dotspec.yaml` and the root system spec, refuses to overwrite, and says when the new project is nested in another.
- `init` then means only "initialize a project", as in git, cargo and npm, and `new` creates items, as in `hugo new`.
- Make `dotspec build` top-level, as `cargo build` and `mkdocs build` are. Update the `spec.stale` and `spec.obsolete` hints (`internal/specdoc/files.go:171,174`), which `check` emits.
- Use one user-facing noun for the artifact, "generated specification". "Report" already names check reports and per-subject `*.spec-report.md` pages.
- Make build's missing-manifest error name `dotspec init`.

### F. The machine contracts differ from command to command (medium)

**Problem.** JSON versions and shapes vary by command, and two flags select the same output.

**Evidence.**
- `schema_version` varies by command:
  - 1 on check, sources, feature, capability and init;
  - 2 on requirement, job and use-case;
  - 3 on `system show` and `quality-model list`.

  `model product` has none.
- `schema_version` isn't the first key on feature, check or quality-model list, which contradicts `speckind.go:21`.
- Job `location` has no `path`, while requirement `location` does.
- `--json` and `--format` both appear on 44 commands and are mutually exclusive. The conflict produces raw Cobra text: "if any flags in the group [json format] are set none of the others can be".

**Recommendation.**
- Reset every command result to `"schema_version": 1`, emitted as the first key. The exception is the snapshot itself (`snapshot export`), which stays versioned by `metamodel.yaml` through `model_version` (0.36.0-draft) and `representation_version` (0.4.0).
- Use one list envelope from the kind table, with entries carrying `kind`, `slug`, `name`, `path` and `system`. `path` is absolute, like every JSON path, because `scripts/check-spec-sources` depends on it.
- Keep `--json` only, as gh does. Remove `--format`, the mutual-exclusion group and the txtar cases that test rejecting `--format mermaid`.
- Make `--json` a local flag on commands that return a result, not a persistent one, so that `docs` and `completion` don't accept a flag that does nothing. Enforce this in `contracts_test.go`.
- Document `snapshot export` as the integration contract, and `--json` on other commands as a per-command convenience.

### G. Help and errors (medium)

**Problem.** Help text leads with qualifications instead of usage, root help gives no structure, and the first errors a new user hits read as library output.

**Evidence.**
- Long help opens with disclaimers:
  - `check.go:19` says "...does not establish full requirements verification...".
  - A shared `readingGuidance` paragraph is appended to two views.
  - `job_views.go:128` prints "Completeness is not assessed." in its output.
- Root help has no structure:
  - It lists commands alphabetically, so `man` and `completion` sit among the domain nouns.
  - Short and Long disagree (`cli.go:190-191`).
  - Examples print flush-left in plain output.
- Typo suggestions are disabled. `dotspec requirment list` prints only `unknown command "requirment"`, because `Args: cobra.NoArgs` (`cli.go:194`, `:218`) skips Cobra's suggestion step.
- A missing operand prints `accepts 1 arg(s), received 0`.
- Filesystem errors leak `lstat /nonexistent: no such file or directory` and `find requirement /abs/path: file does not exist`.
- The exit-code line appears on `review` and `pegs` but not on `needs` or `spec build`, which also exit 1.

**Recommendation.** Apply the conventions in sections 6 and 7. Three cautions:
- Use a validator that suggests corrections. Don't just remove `Args`: on a group, that would let `requirement lst` exit 0 (`cli_test.go:69`).
- `docs <no-match>` stays at exit 2, as the offline-authoring-guidance requirement and `docs.txtar:76` require.
- Keep Cobra's `help [command]` and list it in the Help group. It is conventional (`git help`, `gh help`), and `docs.txtar:5,7` use it. Under verb-first, `dotspec help job` no longer exists. Amend the offline-authoring-guidance criterion "Type help opens a route to the topic" (line 24) so that it names `dotspec new <kind> --help` and the kind list in `new --help`. Update `docs.txtar:7` to match.

### H. Nothing ties commands to the project's own spec (medium)

**Problem.** The project's spec never describes building the generated specification or starting a project, and nothing checks that the commands the docs mention still exist.

**Evidence.**
- `prepare-specifications.use-case.spec.md` runs create, discover, check, make available. It never builds the generated specification.
- No use case covers creating a project.
- `list`, `show`, the quality commands and the model views have no governing spec. That is how the speculative surface accumulated.
- `handbook_test.go` checks links but not `dotspec ...` invocations, so stale `model needs` or `job map` instructions would go unnoticed.
- Command names also appear outside Markdown:
  - Go diagnostics: `internal/specdoc/files.go:171,174`, `internal/system/store.go:45` ("use dotspec system init") and `internal/qualitymodel/store.go:164` ("use dotspec quality-attribute browse").
  - Anchors into the generated `cli.md`: in `job-map.concept.spec.md`, in `cli-interaction.concept.spec.md` (`#dotspec-quality-model-apply`) and in `plans/conventional-quality-model.md:40`.
- `scripts/generate/main.go:192-214` and `scripts/check-example/check.py` hardcode today's command tree. They are covered in section 9.

**Recommendation.**
- Adopt a v1 rule: a command ships only if a use case, requirement or concept in `specs/` governs it.
- Add four use cases or extensions:
  - "Start a specification project" (`init`).
  - "Build the generated specification" (also step 4 of prepare-specifications).
  - "List specifications of one kind by system" (an extension of discover-specification-sources).
  - "Adopt a quality attribute for a system".
- Make the spec amendments the cut requires:
  - The isolation requirement (`-C`, and `--all` only on `check`).
  - `reject-existing-target` (Applies to covers `init` and `new`).
  - The offline-authoring-guidance "Type help" criterion.
  - The discover-specification-sources extension "Multiple sources have the same slug": delete it, or rewrite it for path-based lookup.
  - The specification-inspection concept ("Inspect one specification", by path).
  - The job-map concept (finding D).
- Add a test that resolves each referenced command with `root.Find` plus flag lookup. It reads three kinds of reference:
  - Backticked `dotspec <words> [--flags]` invocations in `docs/cli-handbook`, `docs/reference` (except the generated files), `docs/maintainers`, `specs/` and `examples/*/{README,COVERAGE}.md`.
  - `dotspec ...` invocations inside Go string literals.
  - Every `cli.md#dotspec-*` fragment.

## 4. Proposed v1 command tree

```
dotspec [-C, --chdir <dir>] <command> [args] [flags]

Author
  dotspec init <system-name> [--name <project>]     writes dotspec.yaml + root <slug>.system.spec.md
  dotspec new system <name> --part-of <path>        nested system; heading is the identity, filename derived
  dotspec new component <name> --part-of <path>     Part of is an explicit assertion, never inferred
  dotspec new feature <slug>
  dotspec new capability <slug>
  dotspec new job <slug>
  dotspec new use-case <slug> [--template brief|expanded]
  dotspec new concept <slug>
  dotspec new business-rule <slug>
  dotspec new requirement <slug> [--attribute <id>]
  dotspec new quality-attribute <id>... [--from <model>] [--dry-run]
                                                    blank without --from; with --from copies terms and
                                                    reuses existing parents; exits 2 if a named attribute exists
  dotspec new quality-model <slug> [--from <model>]

Inspect
  dotspec check [<path>] [--all]                    file or directory; default: enclosing project;
                                                    when checking the whole project, also checks that
                                                    the generated specification is current
  dotspec list <kind>                               every kind above, project files only; grouped by system;
                                                    list system | list component print the structural tree
  dotspec sources [<path>]                          specifications, companions, context and referenced files

Publish
  dotspec build                                     SPEC.md + specification/ for the enclosing project
  dotspec snapshot export [<path>]                  default: enclosing project; <path> selects a subtree
  dotspec snapshot schema
  dotspec snapshot validate <file>|-

Help
  dotspec docs [<topic> | <page>[#<section>] | <rule> | <words>...]
  dotspec help [<command>]
  dotspec completion <shell>
  (dotspec man: hidden; dotspec.1 generated into release archives)

Global: -C, --chdir <dir>, -h/--help, -v/--version
--json: local flag on every command with a result (not docs, help, completion);
        snapshot export/schema always write JSON
```

That gives:

- 9 root entries plus `help`, down from 19 plus `help`.
- 8 verbs.
- About 34 leaves, down from about 55. That counts 11 `new`, 11 `list`, the 4 completion shells, `init`, `check`, `sources`, `build`, `docs` and 3 `snapshot` leaves.
- 7 command flags plus `--json` and `-C, --chdir`, down from 19 (excluding `-h` and `-v`). The 7 are `--name`, `--part-of`, `--template`, `--attribute`, `--from`, `--dry-run` and `--all`.

**Why this grammar.** The verbs are the same for every kind. kubectl faces the same situation and uses `kubectl get <type>`. gh is different: there, `pr merge` and `issue close` really do differ by noun. clig.dev accepts either order but asks for one. Kind-specific flags, such as `--part-of`, `--template` and `--from`, live on the generated kind subcommand, as in `kubectl create deployment --image`. So `dotspec new component --help` shows only what applies, and no generic verb collects kind-specific options.

The domain vocabulary stays on the first screens:
- Root help says "Kinds: run `dotspec new --help`".
- Four places list each kind with its domain Short and its `dotspec docs <kind>` route: `new --help`, `list --help`, the error for `new` with no kind, and shell completion.
- Every invocation still names the domain term: `dotspec new requirement fast-save`, `dotspec list use-case`.

`check`, `sources` and `snapshot export` act on files, not kinds, so they take a path. Adding a kind adds one table entry and no new code paths.

## 5. Removed or deferred

| Current surface | Disposition | Reason | Footprint (approx.) |
|---|---|---|---|
| `<kind> init` (10 kinds) | folded into `new <kind>` | `init` is reserved for the project | wiring replaced by the kind table |
| `system init` | folded into `init` (root system) and `new system --part-of` (nested) | no command wrote `dotspec.yaml`; without `--part-of`, a second root system would fail the "one root system" rule at build | — |
| `system init` default "New system" | drop | silent placeholder; the name is now required | — |
| `component init --part-of` (required) | becomes `new component --part-of` (still required) | Part of is an explicit domain assertion, not inferred from file location | — |
| `<kind> check` (8 groups) | folded into `check [<path>]` | byte-identical output | 8 leaves |
| `<kind> list` (9 groups + `system list`) | folded into `list <kind>` | one renderer and one envelope | — |
| `<kind> show` (8 groups + `system show`) | defer | no governing spec; prints authored Markdown plus a System line; removing it removes slug lookup (amend the discover "same slug" extension) | 9 leaves |
| `job map`, `--outcomes` | drop | duplicates `job show` and the generated job page | `job_views.go` 163 + tests 216 |
| `job outcome list`, `--step` | drop | only three-level path; same JSON as `map --outcomes`; job-map concept's step-selection action removed | (in the above) |
| `job init/check --preset` | drop | check results depend on an unrecorded flag | `presets.go` 72 + tests 64; `check.go` branches |
| `requirement init --preset quality` | drop | adds only a Criterion heading that the general template already invites in prose | `QualityTemplate` |
| `use-case list --overview` | drop | text-only; prints broken relative links | `usecase.go:30-54` |
| `requirement list --attribute` | drop | only list filter; the attribute report lists classified requirements | — |
| `model trace/overview/needs/review/pegs` | drop (two advisories moved to `check`; `relationship.motivation` per Open decision 7) | second reporting pipeline; generated spec covers it | 1,215 src + ~490 test |
| `--relative-to` | drop | exists only because views wrote linked Markdown to stdout | — |
| `model product` | drop (`Snapshot.Product` stays internal) | test oracle exposed as public surface; weaker validation | `model.go:99-139`; `product_test.go` 176 |
| `model export/schema/validate` | renamed `snapshot ...`; export keeps a `[<path>]` subtree | removes the "model" ambiguity and the `validate`/`check` confusion | — |
| `model schema --workspace`, export `--all` envelope | drop | no consumer outside tests | `workspace_model.go` 92 |
| `--all` on lists, sources, export, views, `quality-model list` | drop (kept on `check`) | unused; conflicting meaning; isolation requirement narrowed | part of `workspace.go` |
| `quality-model apply --all` | drop | reuses `--all` with another meaning; whole-model adoption ungoverned | — |
| `--project` (51 commands) | folded into `-C, --chdir` | same effect via directory; requirement amended | — |
| `--path` (51 commands) | folded into `-C, --chdir` and positional `[<path>]` | four meanings; false default | `cmd.Name()` switch |
| `--format` (44 commands) | drop | duplicates `--json`; existed for Mermaid | `output.go` annotation + exclusion group |
| `quality-attribute add` | folded into `new quality-attribute --from` (no-overwrite applies to named attributes) | adoption creates the kind | `qualityattribute.go` 233 → table entry + hook |
| `quality-attribute browse` | replaced by the bundled model as an offline docs page, or deferred | `list` stays reserved for project files | (in the above) |
| `quality-model apply`, `--attributes` | drop | duplicate of `add` | `qualitymodel.go:125-144` |
| `--model` | renamed `--from` | one name for "copy from this catalog" | — |
| `--system` (add/apply) | drop | system context comes from the nearest enclosing system, as for every collection kind; use `-C` | — |
| `quality-model init/list/show/check` | `new` / `list` / deferred / folded into `check` | uniform kinds | — |
| `quality.model` manifest key, `Resolve`, `effective_model` | defer | no consumer once `--from` is explicit; `check.py` q42 build-stability test deleted with it | `check.go:39`; `qualitymodel.go:57` |
| bundled FURPS and Q42 | defer | unused; Q42 is share-alike licensed | 269 + 1,875 content lines |
| `spec build`, `spec` group | renamed `build` | one-child group | — |
| `man` | hidden; generated at release | packaging aid | — |
| Cobra `help` | keep, listed under Help | conventional; used by `docs.txtar` | — |
| per-command exit-policy boilerplate | replaced by one sentence where exit 1 is possible | — | — |

## 6. Flag and naming conventions

1. Each flag name has one meaning across the CLI. `--all` means only "every discovered project" and exists only on `check`.
2. `-C, --chdir <dir>` is the only location flag. Every relative path, positional or flag, resolves against the effective working directory.
3. Reads default to the enclosing project. Narrow them with one positional `[<path>]` (`check`, `sources`, `snapshot export`), never with a flag.
4. `--json` is the only output switch. It is a local flag on commands that return a result, and `contracts_test.go` enforces that. Artifact commands (`snapshot export`, `snapshot schema`) always write JSON.
5. Kind-specific flags live on the kind subcommand (`new component`, `new use-case`), never on a generic verb.
6. `--from <model>` always means "copy from this quality model". It appears only on `new quality-attribute` and `new quality-model`. `--template <name>` always means "starting prompts", and appears only where a kind has real alternatives.
7. No flag changes check findings. Findings depend on file content and project configuration only.
8. Offer `--dry-run` only where what gets written depends on project state the caller can't predict. In v1 that is only adoption (`new quality-attribute`).
9. No optional-value flags, no abbreviations, no shorthand-only flags and no hidden aliases. v1 has no `--verbose` or `--quiet`.
10. Arguments follow the kind:
    - Named subjects (`system`, `component`) take `<name>`, which is the heading. The filename is derived from it, and `layout.subject-name` enforces the match.
    - Collection kinds take `<slug>`.
    - Quality attributes take `<id>`, the catalog term ID that also becomes the attribute's slug.
    - Structural relationships (`--part-of`) are always explicit and required.
11. Placeholders use `<x>` for required, `[x]` for optional and `...` for variadic, always lowercase (no `SNAPSHOT` or `[SYSTEM]`).
12. Flag help starts lowercase with no trailing period. A contract test in `contracts_test.go` enforces rules 11 and 12.

## 7. Help and output conventions

1. **Root help.**
   - Group commands with Cobra `AddGroup` as Author, Inspect, Publish and Help, as gh and Terraform do. Fang renders the groups.
   - Short and Long state one sentence covering the whole intent layer.
   - The example is a workflow: `init "Booking desk"`, `new requirement fast-save`, `check`, `build`.
   - Root help names the handbook (`dotspec docs`) and the kind list (`dotspec new --help`), as the offline-authoring-guidance requirement asks.
2. **Long help.** Use one shape:
   - what the command does, in domain terms;
   - only for `check`, `sources` and `snapshot validate`, one precise sentence on what the result establishes;
   - `Learn more: dotspec docs <address>`, with an address that resolves (add `check` and `inspection` topics to `topics.json`).

   Everything else moves to the use-case specs and `inspection.md`.
3. **Shorts.**
   - Generate leaf Shorts from the kind table with one verb per action: "Create a `<kind>` specification", "List `<kind plural>` by system".
   - Never use "Describe" to mean display.
   - State the no-overwrite guarantee once, in the shared `new` Long.
4. **Exit codes.** State the 0/1/2 policy once, in root help and `inspection.md`. On leaves, add "Exits 1 when error findings are reported" only to `check`, `sources`, `build` and `snapshot export/validate`.
5. **Examples.** Indent them two spaces, include no redundant flags (no `--path .`), and keep them executed by `TestDeclaredCommandExamples`.
6. **Errors.**
   - Use the form `dotspec: <what>; <fix>`, ending with `see dotspec <cmd> --help`.
   - Add typo suggestions through a shared validator.
   - Name the missing operand, for example `new requirement: missing <slug>`.
   - Wrap filesystem errors, for example `--chdir /nonexistent: directory does not exist`.
   - Never show Go internals or absolute paths.
7. **JSON.**
   - `schema_version: 1` is the first key of every command result. The snapshot keeps `model_version` and `representation_version`.
   - Lists use one envelope. Arrays are never null.
   - JSON paths stay absolute.
   - `check --all --json` emits the documented workspace result.
   - Bump versions per result family only after release.
8. **Text paths.** This rule applies to text output only.
   - Print paths relative to the working directory, so any path printed by `list` or `sources` can be pasted into `check`.
   - Print no enum strings (`builtin_default`, `blocked-by-source`).
   - Print no `None` for incidental fields. Where a field is the subject, say so in a sentence ("No enclosing system").
9. **`check` text.**
   - Findings come first, then one summary line that keeps the scope and completeness the check-specifications use case requires: `Checked 16 files in examples/todo (scan complete): 1 error, 0 advisories; 6 context files.`
   - Context files are listed after the summary, not before the findings.
   - The generated-spec state is a sentence, for example "Generated specification: not checked (focused check)".
   - Unresolved targets are not listed as context.
10. **`build` failure.** A blocked build leads with `SPEC.md not written: 1 source error`.
11. **`sources`.** Print the project once per group, not on every line.

## 8. Open decisions for the owner

1. **Grammar.** The choice is verb-first with generated kind subcommands (recommended), or noun-first (`dotspec requirement new`) generated from the same table.
   - Noun-first matches the current `cli-design.md:7`, keeps `dotspec help <kind>` working, and puts kinds on the root screen. The cost is 11 root groups with identical verb pairs and about the same leaf count.
   - Both remove the drift.
   - Verb-first also cuts the root from 19 entries to 9, and `new --help` still teaches the vocabulary.
2. **`show`.** Recommendation: defer.
   - No spec governs it. The generated spec is the reading surface, and `snapshot export` is the parsed form.
   - Deferring it removes slug resolution and its ambiguity rules from v1. That means deleting or rewriting the discover-specification-sources extension "Multiple sources have the same slug".
   - Adding `show <path>` later is an additive change.
   - The cost: `system show` is the only command that answers "which system am I in?". `check` has no System line in text or JSON; its `scope` carries only project, path, types and discovery. Either accept the loss, since `list system` shows "Systems by directory", or add a System line to `check` deliberately.
3. **Quality adoption in v1.** Recommendation: `new quality-attribute <id>... --from <model>`, keeping `internal/qualitymodel/apply` (429 lines). It exits 2 when a named attribute exists. Publish the bundled ISO 25010 model as an offline docs page for browsing.
   - The alternative is to defer adoption entirely. Authors would write attribute files by hand, and `apply`, citations and the dirty-build plumbing would be deleted.
   - Either way, drop the `quality.model` default and bundle only ISO 25010 for now.
4. **Name-based project selection.** Recommendation: drop `--project` in favour of `-C, --chdir`.
   - Amend `specification-project-isolation.requirement.spec.md` for both name selection and the workspace-wide clause, which becomes `check --all` only.
   - Keep `check --all`, and change `mise.toml` to `build` plus `build -C examples/todo`.
5. **`model` → `snapshot`.** Recommendation: rename, and retitle the use case "Export a specification snapshot". Keeping `model` is defensible ("specification model"), but it collides with `quality-model` and blurs `validate` against `check`.
6. **Manifest key.** Recommendation: rename `report.enabled` to `build.enabled`, so the key matches the command and doesn't reuse "report".
7. **`relationship.motivation`.** It is a review prompt about justification, and moving it into `check` adds 10 advisories to dotspec's own check run. There are three options:
   - Move it into `check` as an advisory, and update txtar and the command-example advisory counts.
   - Defer it, and let it disappear with `model trace`.
   - Fix dotspec's own requirements first, then move it.

   Recommendation: defer, because no spec requires the prompt. Revisit when a requirement on justification exists.

## 9. Suggested sequencing

Each step is a separate change. Each one updates its own specs, handbook, txtar, generated docs and `scripts/check-example/check.py`, and passes `mise run check`.

`check.py` depends on the tree well beyond lines 101 and 532:
- The coverage gate at `:620-630` hardcodes the groups `{TYPES, quality-attribute, component, spec, model, job outcome, completion}`.
- `:150` calls the views.
- `:420-433` call `<kind> list/show/check`, `--attribute`, `--project` and `--format mermaid`.
- `:458-474` use presets, `--template` and `show`.
- `:482-515` use `apply`, `browse`, `add --system` and `--from furps`.
- `:541-563` check `builtin_default` and the q42 test.

`inventory.py:99` also runs `model export --path`.

`scripts/generate/main.go:192-214` hardcodes the command-examples list, so the first removal would fail `mise run generated-check`.

1. **Safety net.** Add the invocation test from Finding H, covering Markdown, Go string literals and `cli.md#dotspec-*` fragments. Add the placeholder and flag-help contract test.
2. **Help hygiene.** Add the suggesting Args validator and named-operand errors, wrap filesystem errors, indent examples, hide `man` and generate the man page in a goreleaser hook.
3. **Advisories.** Move the assumption Basis/Applicability prompt and the interface-contract prompt into `check` with rule IDs. Handle `relationship.motivation` according to Open decision 7.
4. **Model views.**
   - Remove the five views and `--relative-to`.
   - Repoint the 10+ handbook guides, `docs.txtar:94-104`, `COVERAGE.md:53`, the todo README and `maintainers/handbook.md:82` to report sections.
   - Amend the isolation requirement's Applies to.
   - Update `scripts/generate/main.go` and `check.py:150`.
5. **`model product`.** Remove it and move the `check.py:101` and `:532` assertions into Go tests or onto the generated report.
6. **Job extras and presets.**
   - Remove `job map`, `job outcome`, `--outcomes`, `--step` and `--preset` (job and requirement).
   - Revise the job-map concept: drop the step-selection action, "selection" in State, the needs-map reference and both anchors.
   - Rewrite `write-a-job.md:86-92`.
   - Update `scripts/generate/main.go` and `check.py:420-474`.
7. **Output flags.** Remove `--format`. Reset command `schema_version` to 1 with the version first (the snapshot excepted), and align `location` shapes. Update `scripts/generate/main.go`.
8. **Quality.**
   - Remove `quality-model apply` and its `--all`.
   - Defer FURPS and Q42, add a local multi-parent fixture, and remove the `quality.model` key along with the `check.py` q42 test.
   - Publish the ISO 25010 model as a docs page and update `CONTRIBUTING.md`.
9. **Lifecycle.**
   - Add `dotspec init` and rename `spec build` to `build`, including the `spec.stale` and `spec.obsolete` hints in `specdoc/files.go`.
   - Settle the "generated specification" naming and the missing-manifest hint.
   - Add the "Start a specification project" and "Build the generated specification" use cases.
10. **Scope and grammar, in one change.** Doing them separately would rewrite the roughly 80 handbook invocations twice.
    - **Scope:**
      - Add `-C, --chdir`, and `[<path>]` on `check`, `sources` and `snapshot export`.
      - Make every path relative to the working directory.
      - Remove `--path`, `--project`, the `--all` copies, the export envelope and `schema --workspace`.
      - Replace the `cmd.Name()` switches with typed scope.
    - **Grammar:**
      - Add the CLI kind table and generate `new <kind>` and `list <kind>`. `--part-of` stays required for component and nested system, and `new quality-attribute` refuses existing named attributes.
      - Delete the per-kind groups, `product.go`, the `qualityattribute.go` wiring and the per-kind tests, and add one parameterized contract test.
      - Rename `model` to `snapshot`, add root groups, and list `help` under Help.
    - **Rewrites:**
      - `mise.toml`, the `check.py` coverage gate and command calls, `inventory.py:99`, and `scripts/generate/main.go`.
      - The Go diagnostics in `system/store.go:45` and `qualitymodel/store.go:164`, and the `cli.md#dotspec-*` anchors.
      - The handbook invocations, `cli-design.md` and `inspection.md:83`.
    - **Spec amendments:**
      - The isolation requirement (workspace-wide clause narrowed to `check --all`).
      - `reject-existing-target` (Applies to covers `init` and `new`).
      - The offline-authoring-guidance "Type help" criterion, plus `docs.txtar:7`.
      - The discover "same slug" extension and the specification-inspection concept.
      - New list and adoption use cases.
11. **Text output.** Findings-first check, the blocked-build message, `sources` grouping, no `None`, relative text paths, and quality-model names without enum strings.
12. **Later, additive.** Companion roles in `sources` text and JSON.

## Appendix: dropped findings and recommendation parts

No finding was refuted outright. These claims or proposals were dropped or replaced:

- **Fold `component init` into `system init --part-of`.** Dropped. Component is a first-class domain kind with its own suffix, check type, pages and guide.
- **Default `--part-of` to the nearest enclosing system or component.** Dropped. Part of is a domain assertion. The isolation requirement keeps system relationships "explicit domain assertions, independent of tooling partitions", and a file's location establishes only a reference.
- **Bare slugs resolved project-wide or across kinds, and a `kind/slug` ref grammar.** Dropped for v1. Deferring `show` and making `check` path-based removes the need, and the nearest-subject lookup was deliberate.
- **A new `quality` group with `adopt --every`.** Dropped. It adds a non-kind root group and an unusual flag name.
- **Keep the `quality-attribute browse|add` group.** Not chosen. It keeps a noun that names catalog operations rather than the system's attribute files.
- **`list quality-attribute --from <model>` for catalog browsing.** Dropped. It would give `--from` a second meaning and make one `list` read a catalog instead of project files. The catalog is read through `dotspec docs` instead.
- **`quality-model show [model] [attribute-id]`.** Dropped. Two optional positionals are ambiguous.
- **`--from` with an optional value, or an implicit effective model.** Dropped. Optional-value flags are ambiguous, and every project has a default model, so "blank versus adopt" would have no stable meaning.
- **`--attribute` selects the quality template.** Factually wrong (`store.go:18-24`). Replaced by a single template, with `--attribute` as classification only.
- **Keep `job map` with `--step`, or add `--format mermaid` or `--mermaid` to `show`.** Dropped. It brings back kind-specific output on a generic verb. Mermaid can go into the generated job page later.
- **Remove `--json` from list and show (one JSON contract).** Dropped. It contradicts clig.dev and `cli-design.md:15`, and `--json` has consumers.
- **A global `--json`.** Dropped. It would be accepted and ignored on `docs`, `help` and `completion`.
- **Variadic `check [<path>...]` and `sources [<path>...]`.** Dropped for v1. The report has a single `scope.path`, and its consumers assume one.
- **Make `model product` text by default, or keep it as `snapshot product`.** Dropped. It would need a new renderer or keep a test oracle as public surface.
- **"Report" as the single noun for the generated artifact.** Replaced by "generated specification", because "report" already names check reports and per-subject pages.
- **`list --references` as an alternative to `sources`.** Dropped. It is a second spelling and blurs specifications with the file inventory.
- **`build --all`.** Dropped in favour of two `-C` invocations in `mise.toml`.
- **`dotspec docs <no-match>` exits 0.** Dropped. The offline-authoring-guidance requirement and `docs.txtar:76` require a failed lookup.
- **Remove `Args: cobra.NoArgs` from groups.** Dropped. It would let group typos exit 0. A suggesting validator replaces it.
- **Group Shorts as noun phrases.** Dropped. Cobra and clig Shorts are imperative, and the current domain wording is already intent-revealing.
- **Strip the job-map Mermaid disclaimer.** Moot, because `job map` is removed.
- **Move check context lines behind `--verbose`.** Dropped. The check-specifications use case requires the report to identify context, so context is listed after the summary instead, with no new flag.
- **Give quality-attribute its own `init/list/show/check`.** Superseded by the uniform `new` and `list`. `check` is project-wide, and `show` is deferred.
- **Add product membership to `system show --json`.** Dropped. No reader task needs it yet.
- **Overstated evidence, corrected above:**
  - `measurement_basis` is unique to the views. It isn't; `check` already emits it as `requirement.measurement-basis`.
  - ODI presets change check verdicts. They are advisory only.
  - The workspace envelope is export-only. It is generic, and `check --all --json` uses it.
  - "Users can't see adopted attributes". `system show` and `requirement list` show them.
  - The general requirement template already prompts for a Criterion. It only invites one in prose and has no Criterion heading.
  - `check`'s System line can replace `system show`. `check` has no System line.
  - `check` always checks that the generated specification is current. It does so only for a whole-project check.
  - The job-map concept has two actions. It has three, and one of them needs step selection.