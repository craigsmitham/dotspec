# Conventional quality models and system quality attribute adoption

**Implemented technical specification.** A new dotspec project should be able to browse useful quality terms offline and adopt one for its system without creating a quality model or editing configuration. The first product decision is which quality matters to this system. The resulting system quality attribute file gives the team a place to say what that quality means locally.

For example, the Personal task list team could browse **Reliability → Recoverability**, adopt Recoverability, and edit its definition to identify which acknowledged decisions must survive a restart. Alternatively, it could author **Durability** directly, as the example does today. The system quality requirement then states which fields must survive and how the result is assessed. [System quality attribute](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md) · [System quality requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md)

## Domain contract

| Subject | Who maintains it | What it governs |
| --- | --- | --- |
| Quality model and its quality model attributes | dotspec for bundled starters, or the team that authors a reusable model | Vocabulary available for exploration and adoption. |
| Effective quality model | Selected at command time from the convention, project configuration, or explicit argument | Which catalog a browsing or adoption operation uses. |
| System quality attribute | The system’s authors, in one `*.quality-attribute.spec.md` file | Local definition, identity, and grouping within the system. |
| System quality requirement | The requirement’s authors, in one `*.requirement.spec.md` file | The product promise, its applicability and criterion, classified by locally declared attributes. |

A manifest setting is an authoring preference. Product membership, requirement classification, and the system report’s **Quality Attributes** section derive from authored system declarations. A catalog choice alone creates zero members. Editing a model, switching the preference, or upgrading the CLI takes effect on subsequent catalog operations. Revising an adopted definition requires an edit to its system file. [Product specification composition](product-specification-composition.md)

## Project configuration and resolution

A project containing only `name: personal-task-list` uses the bundled `iso-25010-2023` starter. An override takes one scalar in `dotspec.yaml`.

```yaml
name: personal-task-list
quality:
  model: q42
```

For a team vocabulary, replace `q42` with `./specs/quality-models/team-quality.quality-model.md`. Authors can create that file with the existing `quality-model init team-quality --from iso-25010-2023`, then edit, add, or regroup definitions. A one-system extension can go directly in the system’s quality directory. These two workflows cover local tailoring and reusable customization in the first release. [Author a quality model](../../cli-handbook/guides/define-a-quality-model.md)

1. Resolve the command scope using the existing project and path selectors. For a write, identify the target system before looking up its project default. A supplied `--system` path is relative to the effective command path. Otherwise use the enclosing system in the invocation scope, or the sole system discovered in that project if there is no enclosing system. Preserve the caller’s subdirectory for this lookup when the command starts inside a project. If several systems qualify or none exists, report candidate paths where available and how to select a target. Browsing needs no target system.
2. Select the first applicable model source in this order: explicit CLI selection, the target or browsing project’s `quality.model`, the bundled ISO starter. A nested project has its own convention when it omits the setting. A manifest-free invocation in a supported directory uses the convention. A workspace-only invocation uses existing project selection rules.
3. In the manifest, a bare ID names a bundled starter (`iso-25010-2023`, `furps`, or `q42`). A local file reference is a project-relative path ending in `.quality-model.md`, resolved from the declaring manifest. Require it to identify a regular quality model in the same project; apply existing boundary and symlink checks. This gives checked-in settings the same meaning across machines. Explicit CLI arguments retain the existing bundled ID, local slug, and source-path forms; relative CLI paths use the effective command directory.
4. Use strict YAML decoding. An omitted setting or `quality: {}` selects the convention. An explicit empty, null, non-string, unknown setting, or unsupported locator is a configuration error. Separate manifest parsing from target loading: browsing and adoption load the selected model; unrelated operations can read a well-formed manifest without opening the catalog. An invalid selected source produces a located error and zero writes. An explicit valid CLI selection may bypass a missing configured file. Whole-project `check` still validates the configured selection so CI can catch that missing file.
5. Expose the effective source in text and JSON, including the selection basis (`builtin_default`, `project`, or `explicit`), model ID and source locator. An explicit selection of the ISO starter has basis `explicit`. Bundled models resolve from the installed CLI; local models resolve from the checked-out file. Checking a saved product snapshot remains offline and independent of the current manifest.

## CLI contract

Run the examples from the project directory containing the target system.

| Command | Observable behavior |
| --- | --- |
| `dotspec quality-attribute browse` | Show the effective model and its top-level terms, in declaration order, with IDs and short definitions. |
| `dotspec quality-attribute browse reliability` | Show Reliability’s full model definition, parents, and immediate children. A term with several parents has one identity and is reachable under each. |
| `dotspec quality-attribute add recoverability` | Adopt that ID from the effective model. Infer the target when unambiguous. Support `--system`, `--model`, `--dry-run`, and JSON. |
| `dotspec quality-model apply --attributes recoverability,fault-tolerance` | Apply a batch through the same adoption operation. An optional positional model selects the source explicitly. |
| `dotspec quality-model apply` | Browse the effective model. A selection-free invocation, including an existing invocation supplying only `--system`, is read-only. |
| `dotspec quality-model apply --all --system specs/personal-task-list.system.spec.md` | Request every attribute from one model for one system. Here `--all` selects attributes. It never broadens the target to other systems or projects. |

Browse also accepts `--model`. Keep `quality-model list` as the inventory of available models and mark the effective project default. Keep `quality-model show MODEL` as the full-source view. Unknown attribute IDs on `add` or `apply` produce an error naming the model and how to browse it. For an original term, the diagnostic points to the local system quality attribute authoring form. A misspelling must never create a new product concept.

For batch adoption, reject `--all` combined with `--attributes`, and reject an explicitly empty attribute list. Both write commands accept `--dry-run`; print the same ordered plan and proposed file content the write would use. With no selection, `apply --dry-run` returns the read-only browse view. Return successful no-op and preview results with status 0. Invalid invocations and failed writes use the project’s existing CLI error contract. Whole-project checks retain the distinction between an invalid declaration and incomplete inspection. [CLI design](../cli-design.md)

## Adoption and repeatability

Adding Fault tolerance after Recoverability reuses the local Reliability definition and installs one new file under its directory. [Adoption writer](../../../internal/qualitymodel/apply/apply.go)

1. Under the target-system lock, read its current declarations and index them by local ID and actual file path. Resolve **all explicitly requested IDs** in the chosen model before planning changes. Deduplicate repeated selections.
2. Treat an existing local ID as the system’s authoritative meaning. Mark it **already declared** for a directly requested ID or **reused parent** for an implicit ancestor, and show its local path. Stop traversing the model’s parents at that declaration. A team may already have tailored its wording, title, or grouping. The same ID from a different model also resolves to that local declaration; the result must expose the reuse and must never claim a new adoption from that model. Distinct local meanings require distinct IDs.
3. For each missing ID, include missing ancestors up to those boundaries. Preserve every declared parent ID. Use the first parent’s actual local directory for the child’s physical location, including when the parent’s title and ID differ. Traverse all parent branches, emit each new declaration once, preserve source order for presentation, and validate the combined local graph and destination paths before writing.
4. Copy the selected definition into the ordinary system quality attribute form, relocate relative links, and cite the exact model attribute. For bundled citations, use a shipped source revision or release permalink when available; include a reproducible bundled-model identifier and revision in development builds. For local citations, use a relative URI plus the attribute anchor. The user can edit the local definition immediately.
5. Stage all new files. Install only planned additions, including new child directories beneath existing parents, with exclusive creation, existing project boundary checks, and rollback of additions made by the failed operation. Revalidate under the lock on execution, including after an earlier dry run. Parallel invocations of dotspec targeting the same system serialize. A retry after success returns the existing declarations and creates zero files.

Text output should let the author see what was created, what was reused, and which definition to edit. Version the adoption JSON result when adding the effective model and plan entries. Each entry carries the attribute ID, target path, requested or ancestor role, and create or reuse decision. Include proposed content for creations. Keep command results deterministic for identical inputs and output arrays empty when there are no entries.

## Integration and delivery

1. **Resolve and inspect the default.** Extend the [project manifest parser](../../../internal/specfile/workspace.go) with the optional quality configuration. Put catalog selection and model loading in [the quality model package](../../../internal/qualitymodel/store.go), leaving manifest decoding independent of domain parsing. Use that resolver for browsing, default markers, adoption, and configured-model checks. In `check --all`, resolve each manifest independently and attach findings to the owning project.
2. **Ship incremental adoption.** Refactor the existing writer into a shared plan and publish path, then expose the attribute-first command and the revised [batch command](../../../internal/cli/qualitymodel.go). Make bulk selection explicit in the caller and the writer API. The writer requires an explicit whole-model selection. Update callers, fixtures, and shell examples together.
3. **Publish the authoring workflow.** Lead the handbook with browsing, adopting, and tailoring one system quality attribute. Put model authoring on the reusable-vocabulary path. Document the manifest field, JSON contracts, and migration. Existing bulk scripts must add `--all`; existing source files require no migration. Keep per-system model overrides, model inheritance, and synchronized updates for a later design if projects demonstrate a need.

## Acceptance cases

| Fixture or action | Evidence required |
| --- | --- |
| New project, no quality configuration | Offline browse identifies the bundled ISO starter. File inventory and product membership are identical before and after browsing. |
| Bundled override, local file override, nested project, and explicit CLI override | Text and JSON identify the correct source and basis from the project root and a subdirectory. A one-off selection leaves the manifest untouched. |
| Invalid YAML value, unknown bundled ID, missing file, malformed model, or path into another project | The check identifies the manifest field or source passage, respects exit semantics, and writes zero files. An unrelated project’s default still resolves independently. |
| Two systems plus a component subdirectory | An enclosing system is selected where determinable. An ambiguous project-root write identifies the candidates. Explicit selection places files only in the chosen system. |
| Apply Recoverability, tailor Reliability, then add Fault tolerance | The first write creates Reliability and Recoverability. The second creates only Fault tolerance beneath the existing parent. Its report links to the local Reliability declaration. The tailored parent is byte-identical. |
| Multi-parent selection, divergent local grouping, duplicate selection, repeat run, and ID collision across models | Every new file has the intended parent IDs. Existing local definitions retain their identities, paths, wording, and grouping. Reuse is visible in both output formats. |
| Unknown ID, destination collision, simultaneous sibling additions, and injected publication failure | Unknown selections and collisions produce zero additions. Concurrent successful additions both survive. A failed invocation restores the prior authored tree; a retry produces the expected complete selection. |
| Switch the project default after classifying a system quality requirement | The catalog shown by browse changes. Local attribute files, classified requirement endpoints, product membership, and the system summary’s declared qualities stay the same. An unadopted catalog term is still unresolved as a requirement classification. |
| Empty selection, dry run, and explicit bulk adoption | Empty selection browses. Preview enumerates the planned new files and reuses. A write from the same state matches the preview; bulk adoption requires `--all`. |

Verify through subprocess CLI fixtures, saved snapshot validation, generated system and attribute reports, and a temporary copy of the To-do project. Refresh the generated command reference and examples with `mise run generate`, then run `mise run check`.
