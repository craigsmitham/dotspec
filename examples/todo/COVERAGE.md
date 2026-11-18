# Capability coverage

Run `mise run example-check` from the repository root. For the individual public command paths and their evidence, run `mise exec -- python -I scripts/check-example/check.py --inventory`. The checker walks the freshly built binary's help tree and compares it with the commands exercised by the walkthrough, scratch authoring exercises and subprocess tests.

## CLI workflows

| Surface | Example and check |
| --- | --- |
| Project discovery, selection and scope | Select To-do by path, then by project name for the attribute filter. The gate checks the scope in the JSON inventory. The [workspace tests](../../internal/cli/workspace_test.go) exercise separate multi-project results with `--all`. The [discovery fixture](../../internal/integration/testdata/discovery.txtar) distinguishes tracked, ignored and nested repository files. |
| Seven spec families | List each type and show its actual To-do title. Run each available focused check, then the shared project check. Show the restart requirement through the durability attribute filter. |
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
| `job` | Read the [job](specs/jobs/resume.job.spec.md). Distinguish the job-owned outcome from the step to which it applies. For forces, timelines, research status and alternative approaches, examine the [commuter progress case](../user-needs/specs/jobs/reach-work.job.spec.md). |
| `system`, `business-rule` | Find the [task service's implementation](specs/system.spec.md#task-service). Explain why the completion policy would survive a change of database. Compare component containment with [inventory and reservations](../system-structure/specs/system.spec.md). |
| `concept`, `use-case` | Follow the [return scenario](specs/use-cases/resume-list.use-case.spec.md) to the Reopen action. Predict the effects on identity, title and state. Try the deletion alternative. For multiple concept uses, type bindings, dependencies and when/then synchronization, follow one expiration in the [link-service composition](../concept-design/specs/system.spec.md). |
| `requirement`, `obligation`, `relationships` | Starting at [restart](specs/requirements/restart.requirement.spec.md), locate its governing clause, underlying policy, local-storage assumption and job outcome. Turn Refines into Notes on a copy and observe the exact missing-edge diagnostic. |
| `quality-model`, `quality` | Compare the reusable definition with the system's selected attribute. Find the zero-discrepancy target and its operating conditions. Compare the [recovery feedback draft](../requirements-review/specs/requirements/recovery-feedback.requirement.spec.md), where the bound is still open. |
| `inspection`, `core-domain`, `mapping` | Compare sources JSON, the typed export and the trace view. Locate the two endpoints of the refinement. Inspect which source file contributed the link and which test the requirement designates. Use the [Markdown mapping](../../model/markdown-mapping.md) for alternative headings, metadata and optional forms. |
| `verification`, `review`, `maintaining` | Compare restart results with the [validation plan](specs/system.spec.md#validation-plan). Ask a person to perform the exercise and retain actual observations. For assumptions, open questions, glossary conflicts and assessment cases, review the [recovery example](../requirements-review/specs/system.spec.md). |

## Field inventory and remaining work

The detailed inventory compares each field in the generated [relationship reference](../../docs/reference/model/relationships.md) and structured content fields (including export containers) in the [model catalog](../../model/generated/catalog.json) with actual exported records. It prefers To-do, then inspects the other worked examples and the temporary authoring and diagnostic project. Each populated field names a source or an export to inspect. Read the cited source alongside the typed value to assess its meaning. Presence alone establishes only that the mapping was exercised.

The single-user list covers one actor and one database. The second implementation can give multi-actor authorization, timed retention, nested exception handling, policy refinement and inter-system coordination genuine product consequences. The recovery scenarios, concept composition examples and generated views stay available while those behaviors lack a runnable replacement. Malformed sources, ambiguous references, failed scans and corrupt snapshots belong in temporary diagnostic or conformance cases. The [model tests](../../internal/model) exercise additional representations outside these example exports.

A human still needs to judge whether the requirement is sufficient and whether a person can resume their own work. The walkthrough provides the protocol. User observations are pending.

### Gaps in the exported examples

These are the review assignments for fields that currently have no populated witness in the inventory. The gate detects a previously demonstrated field disappearing. The detailed report identifies each gap by its model field.

| Missing fields | Where to exercise the distinction |
| --- | --- |
| `BusinessRule.refines`, `BusinessRule.refines_references`, `BusinessRule.refines_content` | In recovery, introduce a specific retention policy and check that its Refines endpoint is the general business rule. |
| `Component.interacts_with`, `Component.interacts_with_references`, `Component.interacts_with_content`, `Component.required_interfaces`, `Component.required_interfaces_content`, `Component.out_of_scope`, `System.part_of`, `System.part_of_references`, `System.part_of_content` | Use recovery's external clock and an explicitly declared subsystem. Compare interaction, interface direction and structural membership in the overview. |
| `Component.obligations`, `Concept.obligations`, `ConceptAction.obligations`, `ConceptUse.obligations`, `Synchronization.obligations`, `System.obligations`, `ObligationClause.refines`, `ObligationClause.refines_references`, `ObligationClause.quality_attribute_labels`, `ObligationClause.quality_attributes` | Place a retention or authorization clause with the behavior that owns it. Check its owner, quality classification and refinement endpoint. The [obligation tests](../../internal/model/obligation_capabilities_test.go) compare embedded and standalone forms. |
| `ActorMention.references`, `Condition.references`, `UseCaseExtension.branches` | In the recovery walkthrough, branch on both actor authorization and expiry. Inspect branch containment and citations in actor and condition text. |
| `JobStep.outcomes`, `JobStep.outcomes_content`, `EvidenceStatement.references`, `EvidenceStatement.limitations`, `EvidenceStatement.research_status`, `ProgressCase.description` | Amend a copy of the commuter job. Give one step a local outcome and label the research status, basis and limits of a progress case. Compare job-owned and step-owned results. |
| `GlossaryTerm.references`, `GlossaryTerm.related_terms`, `Specification.preamble`, `System.constraints`, `SourceReference.definition_node` | In a temporary recovery spec, add a preamble, a system constraint and a glossary cross-reference using a Markdown reference definition. Inspect original text, occurrence and resolved target. |
| `ConceptAction.caller_assumptions`, `ConceptAction.caller_assumptions_content`, `ActionParameter.type_description` | Extend a copy of the link-service concept with a caller assumption and a compound input type description. Compare them with guards and explicit type links. |
| `Applicability.workload`, `AssessmentCase.conditions` | Give the recovery assessment an explicit workload and a case-specific setup. Verify which context belongs to the promise and which belongs to the assessment. |
| `Assertion.source_nodes`, `CoverageReport.unbound_content_nodes`, `CoverageReport.unmapped_known_roles`, `UnmappedRole.content_nodes`, `DomainElement.extension_sections`, `ExtensionSection.children`, `ExtensionSection.extension_sections` | Use a scratch document with nested unknown headings and a partly mapped known role. Inspect retained source, binding coverage and provenance. The [validation fixtures](../../internal/model/validation_test.go) test deliberately corrupted source and coverage claims. |
| `Document.diagnostics`, `SourceDocument.diagnostics`, `Resolution.candidates`, `Resolution.diagnostics`, `SourceLocation.source` | Compare a malformed specification and ambiguous target with a valid unresolved reference. Inspect the diagnostic's owning record, candidates and source span. The [diagnostic scenarios](../../internal/integration/testdata/diagnostics.txtar) exercise user-visible findings. |
