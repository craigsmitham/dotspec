# Capability coverage

Run `mise run example-check` from the repository root. For the individual public command paths and their evidence, run `mise exec -- python -I scripts/check-example/check.py --inventory`. The checker walks the freshly built binary's help tree and compares it with the commands exercised by the walkthrough, scratch authoring exercises and subprocess tests.

## CLI workflows

| Surface | Example and check |
| --- | --- |
| Project discovery, selection and scope | Select To-do by path, then by project name for the attribute filter. The gate checks the scope in the JSON inventory. The [workspace tests](../../internal/cli/workspace_test.go) exercise separate multi-project results with `--all`. The [discovery fixture](../../internal/integration/testdata/discovery.txtar) distinguishes tracked, ignored and nested repository files. |
| Generated specification | From [the overview](specification/personal-task-list.system.spec-report.md), follow a feature and the system’s Requirements section to the restart requirement, its edit link and verification method. The gate rebuilds from sources and deliberately breaks a condition and a navigation link in a temporary project. |
| Specification families | List and show each collection. Create a named `component` source and verify the Task service report links to its local acknowledgement requirement. Run each available focused check, then the shared project check. Show the restart requirement through the durability attribute filter. |
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
| `trace` | The restart requirement refines the use case's named obligation. |
| `review` | The return scenario illustrates the Reopen action. Decide whether its contribution to the user's goal is adequate. |

## Handbook and authoring exercises

Use `dotspec guide TOPIC` for each topic below. The questions identify the decision a reader should be able to review.

| Topics | Concrete exercise |
| --- | --- |
| `handbook`, `layout`, `conventions` | Start the [app](README.md). Find the project marker, the specs tree, the file that owns the restart promise and the file that tests it. |
| `job` | Read the [job](specs/jobs/resume.job.spec.md). Distinguish the job-owned outcome from the step to which it applies. |
| `system`, `component`, `business-rule` | Find the [task service's implementation](components/task-service/specs/task-service.component.spec.md#implementation). Explain why the completion policy would survive a change of database. |
| `concept`, `use-case` | Follow the [return scenario](specs/use-cases/resume-list.use-case.spec.md) to the Reopen action. Predict the effects on identity, title and state. Try the deletion alternative. |
| `requirement`, `obligation`, `relationships` | Starting at [restart](specs/requirements/restart.requirement.spec.md), locate its governing clause, underlying policy, local-storage assumption and job outcome. Turn Refines into Notes on a copy and observe the exact missing-edge diagnostic. |
| `quality-model`, `quality` | Compare the reusable definition with the system's selected attribute. Find the zero-discrepancy target and its operating conditions. |
| `inspection`, `core-domain`, `mapping` | Compare sources JSON, the typed export and the trace view. Locate the two endpoints of the refinement. Inspect which source file contributed the link and which test the requirement designates. Use the [Markdown mapping](../../model/markdown-mapping.md) for alternative headings, metadata and optional forms. |
| `verification`, `review`, `maintaining` | Compare restart results with the [validation plan](specs/personal-task-list.system.spec.md#validation-plan). Ask a person to perform the exercise and retain actual observations. |

## Field inventory and remaining work

The detailed inventory compares each field in the generated [relationship reference](../../docs/reference/model/relationships.md) and structured content fields (including export containers) in the [model catalog](../../model/generated/catalog.json) with exported records and generated subject documents. It scans To-do and the temporary authoring and diagnostic project. Each populated field names a source or an export to inspect. Read the cited source alongside the typed value to assess its meaning. Presence alone establishes only that the mapping was exercised.

The single-user list covers one actor and one database. The second implementation can give multi-actor authorization, timed retention, nested exception handling, policy refinement and inter-system coordination genuine product consequences. Malformed sources, ambiguous references, failed scans and corrupt snapshots belong in temporary diagnostic or conformance cases. The [model tests](../../internal/model) exercise additional representations outside these example exports.

A human still needs to judge whether the requirement is sufficient and whether a person can resume their own work. The walkthrough provides the protocol. User observations are pending.

### Gaps in the exported examples

These are the review assignments for fields that currently have no populated witness in the inventory. The gate detects a previously demonstrated field disappearing. The detailed report identifies each gap by its model field.

| Missing fields | Where to exercise the distinction |
| --- | --- |
| `BusinessRule.refines`, `BusinessRule.refines_references`, `BusinessRule.refines_content` | In recovery, introduce a specific retention policy and check that its Refines endpoint is the general business rule. |
| `Component.features`, `Component.features_content`, `Component.capabilities`, `Component.capabilities_content`, `Component.concept_uses`, `Component.concept_uses_content`, `Component.synchronizations`, `Component.synchronizations_content` | In a scratch copy, declare a local feature and its capability under Task service. Link a pair of local concept uses through a synchronization and trace their owners. |
| `Component.interacts_with`, `Component.interacts_with_references`, `Component.interacts_with_content`, `Component.required_interfaces`, `Component.required_interfaces_content`, `Component.out_of_scope`, `System.part_of`, `System.part_of_references`, `System.part_of_content` | Use recovery's external clock and an explicitly declared subsystem. Compare interaction, interface direction and structural membership in the overview. |
| `Component.obligations`, `Concept.obligations`, `ConceptAction.obligations`, `ConceptUse.obligations`, `Synchronization.obligations`, `System.obligations`, `ObligationClause.refines`, `ObligationClause.refines_references`, `ObligationClause.quality_attribute_labels`, `ObligationClause.quality_attributes` | Place a retention or authorization clause with the behavior that owns it. Check its owner, quality classification and refinement endpoint. The [obligation tests](../../internal/model/obligation_capabilities_test.go) compare embedded and standalone forms. |
| `ActorMention.references`, `Condition.references`, `UseCaseExtension.branches` | In a proposed recovery scenario, branch on both actor authorization and expiry. Inspect branch containment and citations in actor and condition text. |
| `JobStep.outcomes`, `JobStep.outcomes_content`, `EvidenceStatement.references`, `EvidenceStatement.limitations`, `EvidenceStatement.research_status`, `ProgressCase.description` | Extend a temporary copy of the resume job. Give one step a local outcome and label the research status, basis and limits of a progress case. Compare job-owned and step-owned results. |
| `GlossaryTerm.references`, `GlossaryTerm.related_terms`, `Specification.preamble`, `System.constraints`, `SourceReference.definition_node` | In a temporary recovery spec, add a preamble, a system constraint and a glossary cross-reference using a Markdown reference definition. Inspect original text, occurrence and resolved target. |
| `ConceptAction.caller_assumptions`, `ConceptAction.caller_assumptions_content`, `ActionParameter.type_description` | Extend a scratch checklist concept with a caller assumption and a compound input type description. Compare them with guards and explicit type links. |
| `Applicability.workload`, `AssessmentCase.conditions` | Give the recovery assessment an explicit workload and a case-specific setup. Verify which context belongs to the promise and which belongs to the assessment. |
| `Assertion.source_nodes`, `CoverageReport.unbound_content_nodes`, `CoverageReport.unmapped_known_roles`, `UnmappedRole.content_nodes`, `DomainElement.extension_sections`, `ExtensionSection.children`, `ExtensionSection.extension_sections` | Use a scratch document with nested unknown headings and a partly mapped known role. Inspect retained source, binding coverage and provenance. The [validation fixtures](../../internal/model/validation_test.go) test deliberately corrupted source and coverage claims. |
| `Document.diagnostics`, `SourceDocument.diagnostics`, `Resolution.candidates`, `Resolution.diagnostics`, `SourceLocation.source` | Compare a malformed specification and ambiguous target with a valid unresolved reference. Inspect the diagnostic's owning record, candidates and source span. The [diagnostic scenarios](../../internal/integration/testdata/diagnostics.txtar) exercise user-visible findings. |
| `Applicability.exclusions`, `Applicability.population`, `AssessmentCase.evidence_to_retain`, `AssessmentCase.expected_result`, `AssessmentCase.name`, `AssessmentCase.procedure`, `AssessmentPlan.basis`, `AssessmentPlan.cases`, `AssessmentPlan.conditions`, `AssessmentPlan.participants`, `Assumption.applicability`, `Assumption.basis` | In a scratch copy of the restart requirement, identify a test population, assumptions, assessment participants, conditions and expected results for a named case. |
| `Capability.depends_on`, `Capability.depends_on_content`, `Capability.depends_on_references`, `Capability.enablers`, `Capability.purpose`, `DesiredOutcome.supporting_content`, `Feature.purpose` | Review a system that needs a supporting capability. Trace its providers and dependencies to a user outcome. |
| `Component.architecture`, `Component.boundary_description`, `Component.provided_interfaces`, `Component.provided_interfaces_content`, `Component.purpose`, `Interface.counterpart`, `Interface.obligations`, `SystemContextSpecification.specification_system_content`, `SystemContextSpecification.specification_system_references` | Describe an external service contract, its owning component and the system responsible for the specification. Inspect the system overview. |
| `Concept.type_parameters`, `Concept.type_parameters_content`, `ConceptAction.effects`, `ConceptAction.initiation`, `ConceptAction.output_parameters`, `ConceptAction.outputs`, `ConceptUse.bindings`, `ConceptUse.bindings_content`, `ConceptUse.configuration`, `ConceptUse.dependencies`, `ConceptUse.dependencies_content`, `ConceptUse.dependencies_references`, `StateDescription.elements`, `StateElement.derivation`, `StateElement.description`, `StateElement.name`, `StateElement.value_description`, `TypeBinding.bound_type`, `TypeBinding.bound_type_references`, `TypeBinding.content`, `TypeBinding.owner`, `TypeBinding.parameter`, `TypeBinding.parameter_references`, `TypeBinding.supplying_concept_use`, `TypeBinding.supplying_concept_use_references` | Bind a parameterized concept to a supplying use in a scratch system. Inspect its state descriptions, action outputs and resolved dependencies. |
| `Synchronization.content`, `Synchronization.owner`, `Synchronization.semantics`, `Synchronization.then`, `Synchronization.then_actions`, `Synchronization.when`, `Synchronization.when_actions`, `Synchronization.where`, `SynchronizationAction.action`, `SynchronizationAction.action_references`, `SynchronizationAction.arguments`, `SynchronizationAction.concept_use`, `SynchronizationAction.concept_use_references`, `SynchronizationAction.content`, `SynchronizationAction.owner`, `SynchronizationAction.result` | Describe two cooperating concepts and verify which action triggers which response, with what arguments and result. |
| `DomainElement.illustrative_examples`, `DomainElement.implementation_suggestions`, `DomainElement.notes`, `DomainElement.open_questions`, `DomainElement.open_questions_content`, `EvidenceCollection.statements`, `EvidenceStatement.evidence_type`, `EvidenceStatement.owner`, `EvidenceStatement.statement`, `GlossaryTerm.alternative_names`, `GlossaryTerm.definition`, `GlossaryTerm.owner`, `OpenQuestion.evidence_needed`, `OpenQuestion.importance`, `OpenQuestion.question`, `OpenQuestion.raised_by`, `OpenQuestion.raised_on`, `OpenQuestion.resolve_by`, `OpenQuestion.tracking_reference`, `Specification.evidence` | Record a consequential open question in a scratch project. Identify its owner, evidence needed, research status and glossary terms. |
| `System.glossary`, `System.glossary_content`, `System.interacts_with`, `System.interacts_with_content`, `System.interacts_with_references`, `System.required_interfaces`, `System.required_interfaces_content`, `System.synchronizations`, `System.synchronizations_content` | In a multi-system design, inspect the declared partner, required interfaces, capabilities, glossary and coordinated actions. |
| `GlossaryTermSpecReport.documents` | Follow a generated glossary term to its authored declaration. |
| `EvidenceStatementSpecReport.documents`, `ProgressCaseSpecReport.documents` | Inspect the passage identities in a job with named evidence and progress cases. |
| `SynchronizationSpecReport.documents`, `SynchronizationActionSpecReport.documents`, `TypeBindingSpecReport.documents` | Follow a synchronization, its named action, and a concept use’s type binding to their definitions. |
