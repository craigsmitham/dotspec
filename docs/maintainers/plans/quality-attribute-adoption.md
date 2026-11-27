# Quality attribute adoption

**Implemented technical specification.** A dotspec project can read a quality model's terms offline and adopt one for its system. The first product decision is which quality matters to this system. The resulting system quality attribute file gives the team a place to say what that quality means locally.

For example, the Personal task list team could read **Reliability → Recoverability** in the bundled ISO/IEC 25010:2023 model, adopt Recoverability, and edit its definition to identify which acknowledged decisions must survive a restart. Alternatively, it could author **Durability** directly, as the example does today. The system quality requirement then states which fields must survive and how the result is assessed. [System quality attribute](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md) · [System quality requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md)

## Domain contract

| Subject | Who maintains it | What it governs |
| --- | --- | --- |
| Quality model and its quality model attributes | dotspec for the bundled ISO/IEC 25010:2023 model, or the team that authors a reusable model | Vocabulary available for reading and adoption. |
| System quality attribute | The system’s authors, in one `*.quality-attribute.spec.md` file | Local definition, identity, and grouping within the system. |
| System quality requirement | The requirement’s authors, in one `*.requirement.spec.md` file | The product promise, its applicability and criterion, classified by locally declared attributes. |

Product membership, requirement classification, and the system report’s **Quality Attributes** section derive from authored system declarations. A quality model alone creates zero members. Editing a model or upgrading the CLI takes effect on subsequent adoptions. Revising an adopted definition requires an edit to its system file. [Product specification composition](product-specification-composition.md)

No project or system has a default model, and the manifest has no quality setting. Each adoption names its model. Per-system model defaults, model inheritance, and synchronized updates of adopted definitions are left for a later design if projects demonstrate a need.

## Model and target selection

1. `dotspec new quality-attribute <id>...` adopts each named ID. `--from <model>` names the bundled `iso-25010-2023` model, a local quality model slug, or a source path; without it, each ID starts as a blank declaration with a definition prompt. A slug resolves among the models owned by the system or component that contains the working directory. Relative paths use the working directory, which `-C` selects. Authors read the bundled model with `dotspec docs quality-models/iso-25010-2023` and a local model in its source file. [Author a quality model](../../cli-handbook/guides/define-a-quality-model.md)
2. The target is the nearest system that encloses the working directory. To adopt into another system, run the command from its directory or select it with `-C`. When no system encloses the working directory, report that none was found and how to start one.
3. Bundled models resolve from the installed CLI; local models resolve from the checked-out file. An invalid model source produces a located error and zero writes. Checking a saved product snapshot remains offline.

## CLI contract

Run the command from a directory inside the target system. A term with several parents has one identity, and adoption copies every parent branch. [Adopt quality attributes](../../reference/inspection.md#adopt-quality-attributes) · [`new quality-attribute`](../../reference/cli.md#dotspec-new-quality-attribute) · [`new quality-model`](../../reference/cli.md#dotspec-new-quality-model)

An unknown attribute ID produces an error naming the model and where to read its terms. For an original term, the diagnostic points to the local system quality attribute authoring form. A misspelling must never create a new product concept.

`--dry-run` prints the same ordered plan and proposed file content the write would use. Return successful previews with status 0. Invalid invocations and failed writes use the project’s existing CLI error contract. [CLI design](../cli-design.md)

## Adoption and repeatability

Adding Fault tolerance after Recoverability reuses the local Reliability definition and installs one new file under its directory. [Adoption writer](../../../internal/qualitymodel/apply/apply.go)

1. Under the target-system lock, read its current declarations and index them by local ID and actual file path. Resolve **all explicitly requested IDs** in the chosen model before planning changes. Deduplicate repeated selections.
2. Treat an existing local ID as the system’s authoritative meaning. Refuse a directly requested ID that the system already declares, under the [no-overwrite rule](../../../specs/requirements/reject-existing-target.requirement.spec.md), and write nothing. Mark an implicit ancestor that already exists as a **reused parent** and show its local path. Stop traversing the model’s parents at that declaration. A team may already have tailored its wording, title, or grouping. The same ID from a different model also resolves to that local declaration; the result must expose the reuse and must never claim a new adoption from that model. Distinct local meanings require distinct IDs.
3. For each missing ID, include missing ancestors up to those boundaries. Preserve every declared parent ID. Use the first parent’s actual local directory for the child’s physical location, including when the parent’s title and ID differ. Traverse all parent branches, emit each new declaration once, preserve source order for presentation, and validate the combined local graph and destination paths before writing.
4. Copy the selected definition into the ordinary system quality attribute form, relocate relative links, and cite the exact model attribute. For bundled citations, use a shipped source revision or release permalink when available; include a reproducible bundled-model identifier and revision in development builds. For local citations, use a relative URI plus the attribute anchor. The user can edit the local definition immediately.
5. Stage all new files. Install only planned additions, including new child directories beneath existing parents, with exclusive creation, existing project boundary checks, and rollback of additions made by the failed operation. Revalidate under the lock on execution, including after an earlier dry run. Parallel invocations of dotspec targeting the same system serialize. A retry after success returns the existing declarations and creates zero files.

Text output should let the author see what was created, what was reused, and which definition to edit. Keep command results deterministic for identical inputs and output arrays empty when there are no entries. [Adopt quality attributes](../../reference/inspection.md#adopt-quality-attributes)

## Acceptance cases

| Fixture or action | Evidence required |
| --- | --- |
| New project | The bundled model is readable offline through `dotspec docs`. File inventory and the manifest are unchanged. |
| Missing or empty `--model`, unknown model, malformed local model | The command fails with exit status 2 or 1, identifies the source where one exists, and writes zero files. |
| Two systems plus a component subdirectory | An enclosing system is selected where determinable. An ambiguous project-root write identifies the candidates. Explicit selection places files only in the chosen system. |
| Adopt Recoverability, tailor Reliability, then add Fault tolerance | The first write creates Reliability and Recoverability. The second creates only Fault tolerance beneath the existing parent. Its report links to the local Reliability declaration. The tailored parent is byte-identical. |
| Multi-parent selection, divergent local grouping, duplicate selection, repeat run, and ID collision across models | Every new file has the intended parent IDs. Existing local definitions retain their identities, paths, wording, and grouping. Parent reuse is visible in both output formats, and a repeat run for a declared ID exits 2. |
| Unknown ID, destination collision, simultaneous sibling additions, and injected publication failure | Unknown selections and collisions produce zero additions. Concurrent successful additions both survive. A failed invocation restores the prior authored tree; a retry produces the expected complete selection. |
| Dry run | Preview enumerates the planned new files and reuses. A write from the same state matches the preview. |

Verify through subprocess CLI fixtures, saved snapshot validation, generated system and attribute reports, and a temporary copy of the To-do project. Refresh the generated command reference and examples with `mise run generate`, then run `mise run check`.
