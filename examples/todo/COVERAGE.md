# Capability coverage

Run `mise run example-check` from the repository root. For the individual public command paths and their evidence, run `mise exec -- python -I scripts/check-example/check.py --inventory`. The checker walks the freshly built binary's help tree and compares it with the commands exercised by the walkthrough, scratch authoring exercises and subprocess tests.

## CLI workflows

| Surface | Example and check |
| --- | --- |
| Project discovery, selection and scope | Select To-do by path, then by project name for the attribute filter. The gate checks the scope in the JSON inventory. The [workspace tests](../../internal/cli/workspace_test.go) exercise separate multi-project results with `--all`. The [discovery fixture](../../internal/integration/testdata/discovery.txtar) distinguishes tracked, ignored and nested repository files. |
| Generated specification | From [the overview](specification/personal-task-list.system.spec-report.md), follow a feature and the system’s [requirements index](specification/requirements.md) to the restart requirement, its edit link and verification method. The gate rebuilds from sources and deliberately breaks a condition and a navigation link in a temporary project. |
| Specification families | List and show each collection. Create a named `component` source and verify the Task service report links to its local acknowledgement requirement. Run each available focused check, then the shared project check. Show the restart requirement through the durability attribute filter. Follow the Durability report to both the colocated restart requirement and the component-owned acknowledgement requirement. |
| Authoring and diagnostics | In a temporary To-do copy, initialize a system and each collection. Exercise general, ODI, progress and job-story jobs, brief and expanded use cases, and general and quality requirements. Preview and apply the local quality model, check duplicate rejection leaves the target intact and break a link to obtain a structural finding. The [quality-model fixture](../../internal/integration/testdata/quality-models.txtar) also exercises bundled starters and ancestor copying. |
| Job navigation | Map the resume job in text and Mermaid. List its steps in authored order and find the named outcome in the job's outcome list. |
| Source inventory and model interchange | Find the acceptance file and its referring requirement, export a snapshot to a temporary file and validate it. Inspect the typed edges. Compare the scope with the [independent snapshot exercise](../../internal/integration/testdata/model.txtar). |
| Guides and diagnostics | Resolve each catalog topic and every cataloged rule to the same text as its bundled source. Open To-do through the offline example index. The [guide scenario](../../internal/integration/testdata/guide.txtar) starts in a consumer directory. |
| Help, completion and manual | Visit the command groups through public help. Run the existing [subprocess scenario](../../internal/integration/testdata/help.txtar) for version, plain output, four shell completions and the roff manual. |

For the model views, run `go run . model VIEW --path examples/todo` from the repository root.

| View | What to look for |
| --- | --- |
| `overview` | The system provides Task list forms. |
| `pegs` | The completion policy appears in the environment and boundaries perspective. |
| `needs` | Less reconstruction appears as a desired outcome. |
| `docs` | The checklist effects table preserves identity and title on Reopen. |
| `trace` | The restart requirement refines the standalone availability requirement. |
| `review` | The return scenario illustrates the Reopen action. Decide whether its contribution to the user's goal is adequate. |

## Handbook and authoring exercises

Use `dotspec guide TOPIC` for each topic below. The questions identify the decision a reader should be able to review.

| Topics | Concrete exercise |
| --- | --- |
| `handbook`, `layout`, `conventions` | Start the [app](README.md). Find the project marker, the specs tree, the file that owns the restart promise and the file that tests it. |
| `job` | Read the [job](specs/jobs/resume.job.spec.md). Distinguish the job-owned outcome from the step to which it applies. |
| `system`, `component`, `business-rule` | Find the [task service's implementation](app/specs/task-service.component.spec.md#implementation). Explain why the completion policy would survive a change of database. |
| `feature`, `capability` | Follow [Pick up where I left off](specs/features/pick-up-where-i-left-off.feature.spec.md) to the restart guarantee and [Local task storage](specs/capabilities/local-task-storage.capability.spec.md) to its provider. |
| `concept`, `use-case` | Follow the [return scenario](specs/use-cases/resume-list.use-case.spec.md) to the Reopen action. Predict the effects on identity, title and state. Try the deletion alternative. |
| `requirement`, `authority`, `relationships` | Starting at [restart](specs/requirements/restart.requirement.spec.md), locate the availability requirement it refines, underlying policy, local-storage assumption and job outcome. Turn Refines into Related specifications on a copy and observe the exact missing-edge diagnostic. |
| `quality-model`, `quality` | Compare the reusable definition with the system's selected attribute. Find the zero-discrepancy target and its operating conditions. |
| `inspection`, `core-domain`, `mapping` | Compare sources JSON, the typed export and the trace view. Locate the two endpoints of the refinement. Inspect which source file contributed the link and which test the requirement designates. Use the [Markdown mapping](../../model/markdown-mapping.md) for declared headings, metadata and optional fields. |
| `verification`, `review`, `maintaining` | Compare restart results with the [participant exercise](README.md#try-the-workflow-with-a-participant). Ask a person to perform the exercise and retain actual observations. |

## Review the return scenario

1. Recording and completing the errand use Checklist’s Add and Complete actions. The completion rule keeps the commitment on the list.
2. On returning to the same database, the availability requirement puts each saved entry in the appropriate view. This step can succeed without a process restart.
3. Reopen acts on the original entry. Its effects preserve identity and title while changing the completion state.
4. In the deletion extension, Delete removes the entry and availability excludes it from both views. The minimal guarantee on a rejected rename follows from the Rename guard.

The [form acceptance test](tests/acceptance.py) checks the return interaction while the service is running, then starts a new process on the same database to exercise the additional restart condition.

## Field inventory and remaining work

The detailed inventory compares each field in the generated [relationship reference](../../docs/reference/model/relationships.md) and structured content fields (including export containers) in the [model catalog](../../model/generated/catalog.json) with exported records and generated subject documents. It scans To-do and the temporary authoring and diagnostic project. Each populated field names a source or an export to inspect. Read the cited source alongside the typed value to assess its meaning. Presence alone establishes only that the mapping was exercised.

The single-user list covers one actor and one database. The second implementation can give multi-actor authorization, timed retention, nested exception handling, policy refinement and inter-system coordination genuine product consequences. Malformed sources, ambiguous references, failed scans and corrupt snapshots belong in temporary diagnostic or conformance cases. The [model tests](../../internal/model) exercise additional representations outside these example exports.

A human still needs to judge whether the requirement is sufficient and whether a person can resume their own work. The walkthrough provides the protocol. User observations are pending.

### Gaps in the exported examples

These are the review assignments for fields that currently have no populated witness in the inventory. The gate detects a previously demonstrated field disappearing. The detailed report identifies each gap by its model field.

| Missing fields | Where to exercise the distinction |
| --- | --- |
| `SystemQualityAttribute.derived_from`, `SystemQualityAttribute.relies_on`, `SystemQualityAttribute.supports` | Give a proposed quality attribute a research basis, identify an assumption it relies on, and trace the goal it supports. |
| `BusinessRule.refines`, `BusinessRule.refines_content` | In recovery, introduce a specific retention policy and check that its Refines endpoint is the general business rule. |
| `Component.concept_uses`, `Component.synchronizations` | In a scratch copy, declare a local feature and its capability under Task service. Link a pair of local concept uses through a synchronization and trace their owners. |
| `Component.interacts_with`, `Component.interacts_with_content`, `Component.required_interfaces`, `Component.out_of_scope`, `System.part_of`, `System.part_of_content` | Use recovery's external clock and an explicitly declared subsystem. Compare interaction, interface direction and structural membership in the overview. |
| `JobStep.outcomes` | Extend a temporary copy of the resume job. Give one step a local outcome. Compare job-owned and step-owned results. |
| `GlossaryTerm.related_terms`, `Specification.preamble`, `System.constraints` | In a temporary recovery spec, add a preamble, a system constraint and a glossary cross-reference using a Markdown reference definition. Inspect original text, occurrence and resolved target. |
| `ConceptAction.caller_assumptions` | Extend a scratch checklist concept with a caller assumption. Compare it with the enforced guards. |
| `Applicability.workload` | Specify the workload under which recovery must succeed. |
| `Document.diagnostics`, `SourceDocument.diagnostics`, `Resolution.candidates`, `Resolution.diagnostics`, `SourceLocation.source` | Compare a malformed specification and ambiguous target with a valid unresolved reference. Inspect the diagnostic's owning record, candidates and source span. The [diagnostic scenarios](../../internal/integration/testdata/diagnostics.txtar) exercise user-visible findings. |
| `Applicability.exclusions`, `Applicability.population`, `Assumption.applicability`, `Assumption.basis` | In a scratch copy of the restart requirement, identify the population, exclusions and assumptions that qualify the promise. |
| `Capability.depends_on`, `Capability.depends_on_content`, `Capability.enablers`, `Capability.purpose`, `DesiredOutcome.supporting_content`, `Feature.purpose` | Review a system that needs a supporting capability. Trace its providers and dependencies to a user outcome. |
| `Component.architecture`, `Component.boundary_description`, `Component.provided_interfaces`, `Component.purpose`, `Interface.counterpart`, `SystemContextSpecification.specification_system_content` | Describe an external service contract, its owning component and the system responsible for the specification. Inspect the system overview. |
| `ConceptAction.effects`, `ConceptAction.initiation`, `ConceptAction.outputs`, `ConceptUse.configuration`, `ConceptUse.dependencies`, `ConceptUse.dependencies_content` | Declare two concept uses in a scratch system. Inspect its state descriptions, action outputs and resolved dependencies. |
| `Synchronization.content`, `Synchronization.owner`, `Synchronization.semantics`, `Synchronization.then`, `Synchronization.when`, `Synchronization.where` | Describe two cooperating concepts and verify which action triggers which response, with what arguments and result. |
| `DomainElement.open_questions`, `GlossaryTerm.alternative_names`, `GlossaryTerm.definition`, `GlossaryTerm.owner`, `OpenQuestion.question` | Record a consequential open question in a scratch project. Link its research basis and identify the terms that need definitions. |
| `System.glossary`, `System.interacts_with`, `System.interacts_with_content`, `System.required_interfaces`, `System.synchronizations` | In a multi-system design, inspect the declared partner, required interfaces, capabilities, glossary and coordinated actions. |
| `Document.glossary_terms`, `GlossaryTerm.relationships`, `System.glossary_sources` | Inspect the [dotspec model vocabulary](../../docs/reference/model/glossary.md) and follow a term in the [system glossary](../../specification/glossary.md) back to its defining heading. Inspect Feature’s topics and relationship to Use case. In a scratch project import two vocabularies defining Account in different contexts and check the distinct entries. |
| `GlossaryTermSpecReport.documents` | Follow a generated glossary term to its authored declaration. |
| `SynchronizationSpecReport.documents` | Follow a coordination rule back to its source heading. |
