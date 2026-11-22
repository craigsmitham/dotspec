# Capability coverage

Use this inventory when reviewing what the To-do example demonstrates or adding a new CLI path. In the dotspec checkout, complete the [pinned tool setup](README.md#run-the-app). From the repository root, run `mise run example-check`. For the command-by-command and field-by-field results, run `mise exec -- python -I scripts/check-example/check.py --inventory`. Both commands build their own CLI. Find the public command in the output and follow its evidence to the walkthrough, group help or subprocess test. For a field marked `GAP`, find its review assignment below.

## CLI workflows

| Surface | Example and check |
| --- | --- |
| Project discovery, selection and scope | Select To-do by path, then by project name for the attribute filter. The gate checks the scope in the JSON inventory. The [workspace tests](../../internal/cli/workspace_test.go) exercise separate multi-project results with `--all`. The [discovery fixture](../../internal/integration/testdata/discovery.txtar) distinguishes tracked, ignored and nested repository files. |
| Generated specification | From [the overview](specification/personal-task-list.system.spec-report.md), follow a feature and the system’s [requirements index](specification/requirements.md) to the restart requirement, its edit link and verification method. The gate rebuilds from sources and deliberately breaks a condition and a navigation link in a temporary project. |
| Specification families | List and show each collection. Create a named `component` source and verify the Task service report links to its local acknowledgement requirement. Run each available focused check, then the shared project check. Show the restart requirement through the durability attribute filter. Follow the Durability report to both the colocated restart requirement and the component-owned acknowledgement requirement. |
| Authoring and diagnostics | In a temporary To-do copy, initialize a system and each collection. Exercise general, ODI and job-story jobs, brief and expanded use cases, and general and quality requirements. Preview and apply the local quality model, check duplicate rejection leaves the target intact and break a link to obtain a structural finding. The [quality-model fixture](../../internal/integration/testdata/quality-models.txtar) also exercises bundled starters and ancestor copying. |
| Job navigation | Map the personal commitments job in text and Mermaid. List its steps in authored order and find the named outcome in the job's outcome list. |
| Source inventory and model interchange | Find the acceptance file and its referring requirement, export a snapshot to a temporary file and validate it. Inspect the typed edges. Compare the scope with the [independent snapshot exercise](../../internal/integration/testdata/model.txtar). |
| Guides and diagnostics | Follow each topic to its bundled pages and each diagnostic to its assigned passage. Open the worked job from the Jobs topic. The [guide scenario](../../internal/integration/testdata/docs.txtar) starts in a consumer directory. |
| Help, completion and manual | Visit the command groups through public help. Run the existing [subprocess scenario](../../internal/integration/testdata/help.txtar) for version, plain output, four shell completions and the roff manual. |

For the model views, run `mise exec -- go run . model VIEW --path examples/todo` from the repository root, substituting one of the view names below.

| View | What to look for |
| --- | --- |
| `overview` | The system provides Task list forms. |
| `pegs` | The follow-up policy says who still owns a result when the chosen day passes. |
| `needs` | Less reconstruction appears as a desired outcome. |
| `docs` | The checklist effects table preserves identity and title on Reopen. |
| `trace` | The restart requirement refines the standalone availability requirement. |
| `review` | The return scenario illustrates Resume after a request for a photo, and Complete when a reply settles the commitment. |

## Handbook and authoring exercises

With dotspec on your PATH, use `dotspec docs TOPIC` for each topic below. To use this checkout’s CLI from the repository root, substitute `mise exec -- go run . docs TOPIC`. The questions identify the decision a reader should be able to review.

Each case changes one fact from the walkthrough. The answers come from the authored specifications.

| Kind | Case | Answer and source |
| --- | --- | --- |
| `job` | A neighbor agrees to collect a parcel. Have I finished the job of getting it to my guest? | Only if the guest has received it or somebody has accepted responsibility for the remaining result. I might establish that by phone even if I never open the list. [Job endpoint](specs/jobs/personal-commitments.job.spec.md#settle-the-commitment) |
| `system` | The courier emails a delivery receipt. Does the entry move? | The owner reads the message elsewhere and decides what to record. The list has no inbox. [Scope](specs/personal-task-list.system.spec.md#scope) |
| `feature` | I closed my laptop while waiting on three shops. Which behavior helps me pick one to chase? | I can see which replies I’m still waiting for and which are due for follow-up. [Feature scope](specs/features/pick-up-where-i-left-off.feature.spec.md#groups) |
| `capability` | The app remembers “Collect parcel” but loses the expected dispatch confirmation. Can it still provide follow-up review? | The owner has lost what to ask for. Review depends on retaining the reply and day. [Capability dependency](specs/capabilities/follow-up-review.capability.spec.md#depends-on) |
| `business-rule` | I owe delivery on Saturday and planned to chase on Thursday. Where does the errand belong on Friday morning? | Still waiting, with follow-up due. Saturday is the promise to the recipient. Thursday is my decision about when to act on silence. [Follow-up rule](specs/business-rules/follow-up.business-rule.spec.md#statement) |
| `use-case` | Two couriers reply. One asks for the door code. The other confirms delivery, which I verify. | Resume the first errand to supply the code. Complete the second. Receiving a message alone settles neither decision. [Return paths](specs/use-cases/resume-list.use-case.spec.md#main-success-scenario) |
| `concept` | Two entries have the same title. I revise the first wait, rename it, complete it and reopen it. | Its identity and revised title survive, its wait is cleared, and the second entry is untouched. [Checklist invariants and effects](specs/concepts/checklist.concept.spec.md#effects-table) |
| `component` | One tab still offers Complete in Open after another tab moved the parcel to Waiting. Which part decides whether the old form is valid? | The local service checks the current state at submission. It rejects the old form and leaves the parcel waiting. [Task service responsibilities](app/specs/task-service.component.spec.md#responsibilities) |
| `requirement` | I leave Waiting open on Thursday night with a Friday follow-up. When do I expect “Follow up now”? | When I reload on Friday. A freshly loaded page that still says “Upcoming” violates the criterion. [Requirement and conditions](specs/requirements/follow-up.requirement.spec.md#applies-to) |
| `quality-model` | An expense journal replaces a saved amount with its earlier value after restart. Is that the same kind of loss as forgetting a revised follow-up? | Both lose an acknowledged decision between sessions. Durability gives the two tools a shared term. [Personal data quality](specs/quality-models/personal-data.quality-model.spec.md#durability) |
| `quality-attribute` | After restart on Friday my Thursday entry keeps its identity and saved Friday date but shows a different cue. Is that a durability failure? | The cue reflects today. Losing the recorded date would be a failure even if the title and entry count matched. [To-do durability](specs/quality/durability/durability.quality-attribute.spec.md#definition) · [Restart measure](specs/requirements/restart.requirement.spec.md#measure) |

### Try changing a refinement into an ordinary citation

Use a disposable copy of the To-do project and a dotspec CLI on your PATH. In the copy’s root, run `dotspec model trace` and find the `refines` row from *Preserve tasks across restart* to *Keep recorded work available*. Change only the heading `## Refines` in the restart requirement to `## Related specifications`. Run `dotspec requirement check restart` and then `dotspec model trace` again. The citation still resolves and the focused check succeeds. The refinement row disappears. Restore the heading and rerun the trace to recover the row.

The example gate tests this specific relationship by applying the same edit in its own temporary copy. Its success output includes `Rejected relationship regression: missing restart refinement`. That message comes from comparing the exported relationships with the example’s expected refinement. [Example gate](../../scripts/check-example/check.py).

## Field inventory and remaining work

The detailed inventory compares each field in the generated [relationship reference](../../docs/reference/model/relationships.md) and structured content fields (including export containers) in the [model catalog](../../model/generated/catalog.json) with exported records and generated subject documents. It scans To-do and the temporary authoring and diagnostic project. Each populated field names a source or an export to inspect. Read the cited source alongside the typed value to assess its meaning. Presence alone establishes only that the mapping was exercised.

The owner is the only actor at the forms. Delivery deadlines, notifications and message collection are outside this list. A shared repair desk could justify authorization, a calendar contract and coordination between systems. A single checklist needs no synchronization between concepts. Leaving those subjects to a larger example keeps the decisions here reviewable in one sitting. Malformed sources, ambiguous references, failed scans and corrupt snapshots belong in temporary diagnostic or conformance cases. The [model tests](../../internal/model) exercise additional representations outside these example exports.

A human still needs to judge whether the requirement is sufficient and whether a person can resume their own work. The walkthrough provides the protocol. User observations are pending.

### Gaps in the exported examples

These are the review assignments for fields that currently have no populated witness in the inventory. The gate detects a previously demonstrated field disappearing. The detailed report identifies each gap by its model field.

| Missing fields | Where to exercise the distinction |
| --- | --- |
| `SystemQualityAttribute.derived_from`, `SystemQualityAttribute.relies_on`, `SystemQualityAttribute.supports` | Give a proposed quality attribute a research basis, identify an assumption it relies on, and trace the goal it supports. |
| `BusinessRule.refines`, `BusinessRule.refines_content` | In recovery, introduce a specific retention policy and check that its Refines endpoint is the general business rule. |
| `Component.concept_uses`, `Component.synchronizations` | In a scratch copy, declare a local feature and its capability under Task service. Link a pair of local concept uses through a synchronization and trace their owners. |
| `Component.interacts_with`, `Component.interacts_with_content`, `Component.required_interfaces`, `Component.out_of_scope`, `System.part_of`, `System.part_of_content` | Use recovery's external clock and an explicitly declared subsystem. Compare interaction, interface direction and structural membership in the overview. |
| `GlossaryTerm.related_terms`, `Specification.preamble`, `System.constraints` | In a temporary recovery spec, add a preamble, a system constraint and a glossary cross-reference using a Markdown reference definition. Inspect original text, occurrence and resolved target. |
| `ConceptAction.caller_assumptions` | Extend a scratch checklist concept with a caller assumption. Compare it with the enforced guards. |
| `Applicability.workload` | Specify the workload under which recovery must succeed. |
| `Document.diagnostics`, `SourceDocument.diagnostics`, `Resolution.candidates`, `Resolution.diagnostics`, `SourceLocation.source` | Compare a malformed specification and ambiguous target with a valid unresolved reference. Inspect the diagnostic's owning record, candidates and source span. The [diagnostic scenarios](../../internal/integration/testdata/diagnostics.txtar) exercise user-visible findings. |
| `Applicability.exclusions`, `Applicability.population` | In a scratch copy of the restart requirement, identify the population and exclusions that qualify the promise. |
| `Capability.purpose`, `DesiredOutcome.supporting_content` | Review a system that needs a supporting capability. Trace its providers and dependencies to a user outcome. |
| `Component.architecture`, `Component.boundary_description`, `Component.provided_interfaces`, `Component.purpose`, `Interface.counterpart`, `SystemContextSpecification.specification_system_content` | Describe an external service contract, its owning component and the system responsible for the specification. Inspect the system overview. |
| `ConceptAction.effects`, `ConceptAction.initiation`, `ConceptAction.outputs`, `ConceptUse.configuration`, `ConceptUse.dependencies`, `ConceptUse.dependencies_content` | Declare two concept uses in a scratch system. Inspect its state descriptions, action outputs and resolved dependencies. |
| `Synchronization.content`, `Synchronization.owner`, `Synchronization.semantics`, `Synchronization.then`, `Synchronization.when`, `Synchronization.where` | Describe two cooperating concepts and verify which action triggers which response, with what arguments and result. |
| `DomainElement.open_questions`, `GlossaryTerm.alternative_names`, `GlossaryTerm.definition`, `GlossaryTerm.owner`, `OpenQuestion.question` | Record a consequential open question in a scratch project. Link its research basis and identify the terms that need definitions. |
| `System.glossary`, `System.interacts_with`, `System.interacts_with_content`, `System.required_interfaces`, `System.synchronizations` | In a multi-system design, inspect the declared partner, required interfaces, capabilities, glossary and coordinated actions. |
| `Document.glossary_terms`, `GlossaryTerm.relationships`, `System.glossary_sources` | Inspect the [dotspec model vocabulary](../../docs/reference/model/glossary.md) and follow a term in the [system glossary](../../specification/glossary.md) back to its defining heading. Inspect Feature’s topics and relationship to Use case. In a scratch project import two vocabularies defining Account in different contexts and check the distinct entries. |
| `GlossaryTermSpecReport.documents` | Follow a generated glossary term to its authored declaration. |
| `SynchronizationSpecReport.documents` | Follow a coordination rule back to its source heading. |
