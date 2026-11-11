# Markdown mapping profile

`dotspec-markdown/12` maps specification Markdown to the
[canonical model](dotspec.yaml). This reference defines the recognized syntax
and domain roles. For reading or validating an exported snapshot, start with
[Specification model](../docs/specification-model.md).

The profile covers all eight specification families: quality model, business rule, system, job, use case,
concept, general requirement and quality requirement. Existing type parsers
remain authoritative for authoring structure and validity. The exporter retains
source and syntax before invoking those parsers. Rejection produces diagnostics
and failed domain decomposition without manufacturing a valid domain object.

## Shared mapping rules

| Authored structure | Mapping |
| --- | --- |
| Original file | `SourceDocument` with original bytes, digest and byte count; all spans address those bytes. |
| YAML frontmatter | Ordered YAML nodes with parsed scalar values, tags, aliases, anchors, comments and original spelling. Invalid YAML remains opaque content with a diagnostic. |
| `type` and `id` | Interpreted document metadata; an accepted specification's explicit ID also maps to its root domain record. System IDs are optional and document-local. |
| First level-one title | Root domain `name`, backed by its heading node. |
| Content before the first body section | Root `preamble`. |
| Recognized level-two section | Named field described below, with its whole content and a binding for its heading. |
| Recognized item heading | Addressable domain record with a source-backed name and explicit owner. |
| Rich narrative | Ordered source-node roots, `presence: authored`, and a separate decomposition state. |
| Present section without body content | `presence: empty` and an empty node array; the heading binding preserves the occurrence. |
| Absent section | The corresponding domain slot is omitted. |
| Unknown domain section | `ExtensionSection` with heading, full content, nearest addressable owner and a mapping-gap reason; nested extension sections retain their hierarchy. |
| Ordinary referenced Markdown | Source and content model, document-owned extension sections, and an explicit fallback binding; no invented specification type. |

Interpreted document `type` and `id` can follow a YAML alias to a string scalar;
the content tree retains the alias occurrence and its own source span.

Heading names match the names and aliases in this profile. Original spelling,
heading style, level and occurrence order remain in the content tree. Repeated
occurrences are retained even when the authoring parser rejects duplicate
structure. Repeated supplementary sections that the authoring profile accepts
append their content or items in source order instead of replacing earlier
occurrences. A section's full content includes its subordinate headings and
blocks; its introduction ends before its first subordinate heading.

Lists retain their nesting. Structured collections take the direct items of a
list at the mapped field's level; nested bullets remain part of their owning
item unless an explicit rule below introduces another structure. Source order
is representational. It does not, by itself, impose a required execution order.

Recognized explicit labels such as `**ID:** value` and `**Stage:** value` apply
to introductory paragraphs in the relevant item. Matching text inside a code
example or a deeper unrelated section is not metadata. Label values and retained
narrative point to source nodes, without rewritten copies of the prose.

## Shared specification content

These roles apply to each specification family where authored.

| Section | Domain field or structure |
| --- | --- |
| Business rules | `business_rules: ContextualReferences`; a citation does not assert enforcement or compliance |
| Related specifications | `related_specifications: ContextualReferences` |
| Related use cases | `related_use_cases: ContextualReferences` |
| Related concepts | `related_concepts: ContextualReferences` |
| Authority; Authority and evidence | `authority_and_evidence: RichContent` |
| Notes | `notes: RichContent` |
| Validation plan (TraceableElement) | `validation_plan: AssessmentPlan`; explicit subject and basis distinguish requirements/system validation and assumption assessment |
| Explanatory notes | `explanatory_notes: RichContent`, including on named domain records |
| Implementation suggestions; Implementation hints | `implementation_suggestions: RichContent`; nonbinding advice, not endorsement |
| Illustrative examples | `illustrative_examples: RichContent`, distinct from general obligations |
| Evidence | `EvidenceCollection` with full content, direct list-item statements and named `EvidenceStatement` records |
| Open question; Open questions | `OpenQuestionCollection.content` and ordered `questions` values; each question is distinct from its Evidence needed subsection |

Contextual references retain full explanatory prose, reference occurrences and
explicit list-item descriptions. They do not create support, authority,
applicability or satisfaction assertions. Evidence collections recognize
Research basis and Limitations subsections. Named evidence statements recognize
Type, Research status and Limitations. An authored hypothesis remains a claim;
the mapper does not assign confidence or prove it.

An open-question section retains its complete narrative in the collection. A
singular heading maps its introductory question and explicit detail subsections
into one question. Raised by, Raised on, Resolve by, Importance and Evidence needed
map to corresponding rich-content fields; Tracking reference maps to
`tracking_reference: ContextualReferences`. Dates, milestones, conditions and
importance labels retain authored wording; no workflow state is inferred. A plural heading can contain named questions or
direct list items. Evidence requirements stay in `evidence_needed`; they do not
replace the question. An explicitly empty question remains empty, and introductory
prose survives even when the section also contains named entries.

## Business rule

A `business-rule` specification maps to `BusinessRule`, independently of any
software enforcement. Its optional ID remains document-local. The rule belongs
to the core business-requirements subset; subset metadata is not an authored
requirement-level field.

| Section | Domain field |
| --- | --- |
| Statement | `statement: RichContent` |
| Applies to | `applicability: Applicability`, including the same explicit dimensions as requirements |
| Rationale | `rationale: RichContent` |
| Authority; Authority and evidence | Shared `authority_and_evidence: RichContent` |
| Exceptions | `exceptions: RichContent`, retaining qualifications and references |
| Related specifications | Shared contextual-reference content |

Absent and empty exceptions remain distinct. Unknown sections remain visible
extensions. The existing generic Constraint record remains available for
constraint narrative in system documents; identifying a BusinessRule is an
explicit authoring decision.

## Quality model

A `quality-model` specification maps to a standalone `QualityModel`, with no
implicit enclosing-system association.

| Section or item | Mapping |
| --- | --- |
| Description | `description: RichContent` |
| Edition | `edition: RichContent` |
| Sources | `sources: ContextualReferences`, preserving prose and source links |
| Quality attributes | `quality_attributes_content` plus ordered `quality_attributes` |
| Named attribute subsection | `QualityAttribute` owned by the declaring `QualityModel` |
| Attribute ID, Sub-attribute of, definition, open questions | Same syntax and content accounting as system attributes below |

A parent label on a model definition resolves inside that model and records
`quality_model_context`. Requirement classification continues to resolve only
against a system's declarations. Both owner types project
`DeclaresQualityAttribute` assertions. Copy provenance uses ordinary source links;
no synchronization, equivalence, or satisfaction relationship is inferred.

## System

| Section or item | Mapping |
| --- | --- |
| Frontmatter `id` | Optional System `authored_id`, document-local anchor |
| Purpose; Scope | `purpose`; `scope` |
| System or Component: Out of scope | `out_of_scope: RichContent` |
| System or Component: Interacts with | Authored content and reference occurrences; resolved `interacts_with` System/Component targets, independent of `part_of` |
| System or Component: Required interfaces; Provided interfaces | Separate `InterfaceCollection` values with full content and ordered named `InterfaceDescription` entries |
| Named interface | Heading `name`, introductory `description`; Contracts and Counterpart subsections are contextual references |
| Project context | `project_context: ContextualReferences`; authoritative project artifacts remain external |
| Part of | `part_of_content`, explicit `part_of_references`, resolved scalar System `part_of`; no inferred parent |
| Components; Features; Capabilities | Paired section content and ordered declaration membership in `components`, `features` or `capabilities` |
| Named Component/Feature/Capability | `name`, complete `content`, declaring System `owner`, optional introductory `**ID:**`; Feature `system` is the declaring System |
| Component/Feature/Capability: Purpose | `purpose` |
| Component: Responsibilities; Boundary | `responsibilities`; `boundary_description` |
| Component: Part of | Narrative and explicit references; resolved scalar System or Component `part_of` |
| Component: Contracts; Architecture; Implementation | `ContextualReferences` preserving narrative and links to native artifacts |
| Component: Obligations | Shared named `ObligationClause` declarations; no implicit applicability |
| Feature: Groups | Narrative and explicit references; resolved non-owning `groups` membership in requirements, clauses, use cases, concepts, concept uses or synchronizations |
| Capability: Supports; Depends on; Provided by | Explicit semantic associations to user needs/goals, Capability dependencies, and System/Component providers; declaration does not imply provision |
| Capability: Enablers | `enablers` narrative, including human/procedural providers and external references |
| Environment | `EnvironmentDescription.content` |
| Goals | `GoalCollection.content`; direct list items and named subsections become `SystemGoal` records |
| Assumptions | `AssumptionCollection.content`; direct list items and named subsections become `Assumption` records |
| Assumption: Basis; Applicability; Validation plan | `basis`, prose `applicability`, and `validation_plan: AssessmentPlan`; statement contains the named assumption introduction only |
| Constraints | `ConstraintCollection.content`; direct list items and named subsections become `Constraint` records |
| Glossary | `Glossary.content`; named subsections become `GlossaryTerm` records with full definitions |
| Term: Alternative names; Aliases | List-item values, or the unsegmented subsection, populate `alternative_names` |
| Term: Related terms | `related_terms: ContextualReferences` |
| Quality attributes | `System.quality_attributes_content` preserves the section; `System.quality_attributes` lists its attributes in authored order |
| Named attribute subsection | `QualityAttribute` with title, definition and declaring `System` owner |
| Attribute `**ID:**` | Required system-local `authored_id` |
| Attribute `**Sub-attribute of:**` | Ordered `subattribute_of_labels` from comma-separated IDs; resolutions supply multi-valued `subattribute_of` targets |
| Attribute: Open question/question sections | Explicit `OpenQuestion` values |

The attribute definition excludes its introductory ID and Sub-attribute of metadata lines
while retaining all other definition content. Subheadings that subdivide a
quality-attribute or glossary definition remain part of that narrative unless a
specific rule above assigns a more precise role. Unrecognized headings inside
structured detail fields, such as a glossary alias list, receive extension
records. An absent section leaves both system fields absent. An explicitly empty section retains empty content and an empty attribute list. Directory-derived enclosing system documents
and specification ownership are associations; they introduce no inheritance.

## Job

| Section or item | Mapping |
| --- | --- |
| Job | `statement` |
| Performer; Situation; Motivation | `performer`; `situation`; `motivation` |
| Intended progress | `intended_progress`; the old singular Desired outcome heading is extension content |
| Steps | `JobMap.content`; named subsections become ordered `JobStep` records |
| Step `**ID:**`; `**Stage:**` | Optional `authored_id`; source-backed `stage` |
| Step introduction | `description`, excluding its explicit ID and Stage metadata lines |
| Outcomes under the job or a step | `OutcomeCollection` and direct list-item `DesiredOutcome` records with that actual owner |
| DesiredOutcome `**identifier:**` prefix | Optional outcome `authored_id`; complete item remains available as its statement |
| Named outcome with introductory `**ID:**` | `name` from its heading, `authored_id`, and full `statement` from introductory prose excluding ID; headings without IDs remain extensions |
| Named outcome: Applies to | `applies_to_content`, authored reference occurrences, and resolved Job/JobStep `applies_to` targets, independently of declaration owner |
| Named outcome: Supporting content | `supporting_content` |
| Nested outcome support content | `supporting_content` on the outcome |
| Progress cases | `ProgressCaseCollection.content`; named subsections become `ProgressCase` records |
| Case: Circumstances; Struggling moment; Participants | `circumstances`; `struggling_moment`; `participants` |
| Case: Alternatives considered; Trade-offs and selection criteria | `alternatives_considered`; `tradeoffs_and_selection_criteria` |
| Case: Current approach; Desired progress | `current_approach`; `desired_progress` |
| Case: Forces | `ForceCollection` with full narrative and explicit list-item or headed force values |
| Case: Timeline | `Timeline.content`; direct list items or named events supply `TimelineEvent` values |
| Event: Time; Date | Authored `time` content, without inferred dates |
| Case: Evidence | The shared evidence structure |

A force receives `kind` only from an explicit push, pull, anxiety or habit label.
Other force descriptions retain their content without a fabricated category.
Stage labels remain open. An explicit empty Stage field remains distinguishable
from an omitted Stage field. Unknown detail headings under a step, progress case
or timeline event become extension sections on the nearest addressable owner.
A timeline saying “to be established” remains authored content. Overall outcomes belong to the job and are not copied into every step.

## Use case

| Section or item | Mapping |
| --- | --- |
| Goal; System boundary; Trigger; Success criterion | `goal`; `boundary_description`; `trigger`; `success_criterion` |
| Primary actor | `ActorMention` with full content and reference occurrences |
| Goal level | `GoalLevel.content`; recognized explicit wording may also populate canonical `level` |
| Preconditions; Minimum guarantees | `ConditionCollection` with full content and explicitly separated list items or headed conditions |
| Supporting actors | `ActorCollection` with full content and list-item or named actor mentions; named Services subsections populate `services` |
| Stakeholder interests | Full collection content and explicit list-item or named stakeholder-interest entries |
| Step: nested `Illustrates:` bullet | `illustrates_content`, source occurrences, resolved `illustrates` Obligation/ConceptAction targets; ordinary step links remain citations |
| Main success scenario | `Scenario` with full narrative; direct ordered-list items become `ScenarioStep` records |
| Step `**Actor:**`, `**Action:**`, `**Response:**`, `**Outcome:**` | Preserved in step `content`, without separate semantic fields |
| Extensions | Full `ExtensionCollection.content` and named or explicitly labelled branch records |
| Sequence | `SequenceStatement.content`; explicit Mode and list-item constraints as described below |
| Related jobs | `related_jobs: ContextualReferences` |
| System boundary: explicit `**Subject system:**` link | `subject_system_references`; resolved `subject_system: System` |
| Supports | `supports_content`, `supports_references`; resolved `supports` (SystemGoal, Job, JobStep or DesiredOutcome) |

Scenario steps retain their authored numbers or labels, narrative, ordinary
references and source order. Ordered-list items and explicitly labelled handling
steps such as `2a1` or `.1` become `ScenarioStep` records. A paragraph-only scenario
remains valid without invented steps. Actor, action and result wording remains
in the narrative, including explicitly labelled fields from older documents.

Each extension retains its full `content` and ordinary `references`. Optional
Condition and Location subsections or introductory fields populate `condition`
and `location`. At step and Point of occurrence are aliases for Location. Location
wording remains authored content; “during steps 2–3” does not resolve a step target.

Handling subsections or `**Handling:**` fields supply `handling` narrative.
Response and Continuation fields or subsections, and Flow subsections, also
contribute to that narrative in source order. An Unresolved handling subsection
retains its heading and content; no handling-status enum is inferred. Direct
ordered or explicitly labelled steps in the branch introduction or handling
subsections populate `steps`. Unlabelled prose remains in full `content` without
inferring condition, handling or completion. References can appear anywhere;
no Governing references subsection is required.

Direct headings in the use case's Extensions section distinguish branches.
Inside a branch, a nested branch must have an explicit label (for example,
`1a2a` or `.1a`) or appear under an Extensions subsection. Other detail headings,
such as Notes, become `ExtensionSection` records. Explicitly branch-labelled list
items also create branches, including nested items under a handling step. Labels
include `2a`, `*a`, `1a2a` and `.1a`; `authored_label` preserves the label while
Markdown containment establishes ownership. Labels do not determine control-flow
edges, durable identity or return targets. Sibling branches retain source order.

A Sequence section's `**Mode:**` field recognizes `representative`, `required`
and `partially_constrained`. Direct list items can carry explicit sequence
statements and their references. Ordinary narrative remains complete content
without silently deriving ordering edges. Canonical goal levels recognize
the following wording after trimming whitespace, lowercasing and removing a
final period; other wording remains in the source-backed value.

| Canonical goal level | Recognized wording |
| --- | --- |
| `summary` | `summary`, `kite`, `summary (kite)` |
| `user_task` | `user task`, `user-task`, `sea`, `sea level`, `user task (sea level)` |
| `subfunction` | `subfunction`, `fish`, `subfunction (fish)` |

## Concept

| Section or item | Mapping |
| --- | --- |
| Purpose | `purpose` |
| State | `StateDescription.content`; direct list items and named state components retain descriptions |
| Type parameters; Introduced types | Full section content and named `ConceptType` records in `type_parameters` or `introduced_types` |
| State: Initial state | `abstract_state.initial_state` narrative |
| State: Invariants | Explicit condition values in `invariants` |
| State component: Value; Values; Type | `value_description` |
| State component: Derivation | `derivation` narrative; no computed value is inferred |
| Effects table | `EffectsTable.content` plus ordered rows and rich-content cells from the first table; first row is the header, first column identifies operations; blank cells remain unspecified |
| Actions | `ActionCollection.content`; named subsections become `ConceptAction` records |
| Action: Inputs; Outputs; Initiation; Effects; Rejections; Rejection behavior | Corresponding action fields; both rejection headings map to `rejections` |
| Action: Caller assumptions; Assumptions | `caller_assumptions: ConditionCollection` |
| Action: Guards | `guards: ConditionCollection` |
| Action: Open question/question sections | Explicit `OpenQuestion` values |
| Operational principles | `operational_principles`; named subsections may distinguish individual stories |
| Named Inputs/Outputs subsections | `ActionParameter` records in `input_parameters`/`output_parameters`, preserving full content |

An unsegmented Actions paragraph remains complete named actions content. Abstract
state is a domain description, without an inferred database schema. Caller
assumptions and enforced guards have separate fields.

An action parameter's introductory `**Type:**` field records reference occurrences
and, when resolvable, one concept-local `ConceptType` target. `**Primitive:**`
records an open primitive name. A nested Type subsection is `type_description`
prose, including compound expressions; links inside it do not assert an atomic
type. Type declarations share a name scope within their concept; actions have
another scope. Parameter names are unique within each input or output list.

## System concept composition

| Authored structure | Mapping |
| --- | --- |
| Concept uses | `System.concept_uses_content` and explicitly headed `ConceptUse` records |
| Use: Definition | `definition_references`; resolved `definition: Concept` when one valid reference identifies a definition |
| Use: Bindings | Full `bindings_content` and named `TypeBinding` records |
| Binding: Parameter; Bound type; Supplying concept use fields | Authored `parameter_references`, `bound_type_references`, `supplying_concept_use_references`; corresponding resolved targets when valid |
| Binding: Primitive field | `primitive` string; mutually exclusive with Bound type/Supplying concept use references |
| Use: Configuration | Authored `configuration` narrative |
| Use: Dependencies | Full `dependencies_content`, reference occurrences and resolved same-system `dependencies` |
| Synchronizations | `System.synchronizations_content` and explicitly headed `Synchronization` records |
| Synchronization: When; Where; Then; Semantics | Corresponding narrative fields |
| Named entries under When/Then with Concept use or Action fields | `SynchronizationAction` records in `when_actions` or `then_actions` |
| Entry: Concept use; Action fields | Reference occurrences and resolved `concept_use`/`action` targets |
| Entry: Arguments; Result subsections | Corresponding narrative fields |

Reference occurrences retain ordinary source resolution. A file target can identify
its root Concept; a heading fragment can identify its named domain record. This
promotion does not change the ordinary reference target or trigger further reads.
An authored external, missing, ambiguous or unacquired target never becomes an
invented domain relationship. A scalar relation requires one reference. Resolved
bindings must name a parameter of the used definition and an introduced type of
the supplying use's definition. Synchronization actions must belong to their use's
definition, and uses/dependencies must belong to the same system. Type bindings
are unique by parameter per use. Named uses and synchronizations have system-local
name scopes. Unknown headings remain extension sections on their nearest owner.

Composition errors are source diagnostics, not invalid snapshot shape. Export
returns status 1 while preserving them; offline validation checks that materialized
relationships agree with explicit references and scope. Incomplete declarations
produce advice. Full source content remains available regardless of resolution.
Dependencies may cycle. Neither bindings nor incidental links infer dependency
edges. No argument-expression checking, execution, timing, atomicity or behavior
is inferred. The [authoring guide](../docs/develop-concepts.md) supplies examples.

## General and quality requirements

| Section or metadata | Mapping |
| --- | --- |
| Applies to | `Applicability.content` |
| Applies to: Conditions; Workload; Population; Exclusions | Corresponding applicability fields |
| Applies to: Subject; Subjects | `subjects: ContextualReferences` |
| Applies to: explicit `**Targets:**` links | Root `applies_to_references`; resolved `applies_to: DomainElement` (also on BusinessRule) |
| Rationale | `rationale` |
| Frontmatter `quality_attributes` | Every authored classification label occurrence; resolved quality attribute targets populate `quality_attributes` |
| Statement | `statement` |
| Verification plan | `verification_plan: AssessmentPlan` for planned realization conformance |
| Validation plan | `validation_plan: AssessmentPlan` for needs/intent assessment; also recognized on other traceable records |
| Criterion (optional assessment details) | `QualityCriterion.content` and recognized explicit criterion fields |
| Measurement plan | `MeasurementPlan.content` and recognized explicit measurement fields |

Verification plan and Validation plan share `AssessmentPlan`. They recognize Subject,
Basis (`ContextualReferences`), Participants, Method or Procedure, Success criteria,
Conditions, and Evidence to retain or Evidence retention. Verification plan on an
Obligation assesses a realization against the obligation. Validation plan is available
on TraceableElement records, including requirements, systems, concepts, use cases,
needs and assumptions; the subject and basis determine the intended-use or premise
assessment. Assumptions distinguish empirical support from justified reliance.
Neither field records an execution, verdict or acceptance. Repeated supplementary
plan content and cases append in authored order. Named Cases subsections supply
`AssessmentCase` values; each can contain Conditions, Procedure, Expected
result and Evidence to retain.

Criterion recognizes Measure, Comparison, Target, Units, Aggregation,
Window or Observation window, Conditions, Workload, Population and Exclusions.
Explicit open-question headings remain distinguishable from an omitted target.
Qualitative criteria remain valid; classification labels cannot supply missing
thresholds.

Measurement plan recognizes Success criteria, Sampling, Conditions, Procedure, Evidence to retain or
Evidence retention, Workload, Population and Exclusions. Full measurement
narrative remains available alongside those components, including procedural
qualifications, examples and evidence instructions.

## Bindings, fallback and completeness

Bindings record a rule under `dotspec-markdown/12`, an owner and field path, source
occurrence and role. Named-field, record-identity and structured-item bindings
account for interpreted domain content. Metadata bindings account for frontmatter.
Structural-wrapper bindings cover headings and syntax support; a document-root
wrapper does not claim that its substantive descendants were mapped.

Extension bindings identify unknown domain sections and carry a reason. Fallback
bindings retain uninterpreted or invalid-domain content with an explanation.
A known role routed to fallback remains an unmapped known role. The coverage
report exposes unbound nodes, unmapped known roles, extension sections and opaque
syntax independently. `complete_for_profile` cannot conceal unbound nodes or
missing known mappings, and it does not claim interpretation of every prose
proposition or agreement with the specification.

Decomposition is a separate question. Whole-field narrative is a valid domain
value even when no finer structure is extracted. Source retention, domain
mapping, authoring validity, reference resolution, scan completion and snapshot
conformance each retain their own evidence and status.

The implementation is in the [mapping source](../internal/model/mapping.go)
and the family-specific `mapping_*.go` files beside it. The
[mapping tests](../internal/model/mapping_test.go) exercise every checked-in
specification family and explicit structures; the
[export conformance tests](../internal/model/export_conformance_test.go) check
whole-snapshot behavior. Model projections are regenerated and checked as
specified in the [model README](README.md#regeneration).

## Relationship contract

See [Declare relationships](../docs/develop-relationships.md) for complete examples.
Subject system references are scalar; Supports and Targets may name multiple
records. Ordinary links outside these explicit forms remain contextual. Definition,
parameter, type, use and action fields likewise retain all authored occurrences,
including malformed multiple references; an invalid draft must remain representable.
Only successfully resolved, compatible relationships populate semantic targets.
Unknown declaring scope defers dependent resolution rather than inventing a scope
error. Requirements, business rules, jobs, use cases and concepts may have a
`system_context` derived from directory discovery. An explicit Specification system
section overrides that default through `specification_system`; its origin is
`explicit` or `directory`. Systems use `enclosing_system` for discovery context
and explicit `part_of` for structural parentage; QualityModels are standalone.

Typed fields are canonical. Assertions are optional derived views with a required
`projection_path`, subject and object. HasSubjectSystem, Supports and AppliesTo now
project the explicit fields above; AppliesTo supports both requirements and business
rules. Membership projections retain their originating role, including the two
HasConceptType paths. Section collection wrappers retain required content; sibling
declaration/content pairs use reciprocal `paired_with` schema annotations.

Optional input/output parameter and synchronization action collections require the
corresponding complete Inputs/Outputs/When/Then content through `requires_content`.
The converse is not required: prose-only sections remain valid. Explicit boundary
and applicability occurrences likewise require their narrative sections.

## Layered traceability and obligation clauses

`TraceableElement` supplies explicit Supports, Derived from and Relies on sections.
It is the base of SpecificationElement and of named SystemGoal, Assumption,
Constraint, JobStep, DesiredOutcome, ConceptUse, ConceptAction, Synchronization and
ObligationClause records. The exact heading annotations `trace_heading` and
`trace_targets` in the canonical schema drive mapping and endpoint resolution.
The mapper also recognizes these sections directly beneath named nested records.
It never interprets arbitrary prose or ordinary citations as semantic links.

Each trace section maps to `<field>_content: RichContent` plus ordered
`<field>_references`. Resolved targets populate the multivalued canonical field.
Supports targets SystemGoal/Job/JobStep/DesiredOutcome; Derived from targets DomainElement
or SourceFragment; Relies on targets Assumption. UseCase Involves targets ConceptUse.
Requirement, BusinessRule and ObligationClause Refines targets Requirement,
BusinessRule or ObligationClause. Every section link is a declared target; rationale
prose remains in that section, and contextual citations belong elsewhere.

A file reference promotes to a specification's root domain record. A fragment
promotes to its named domain record where one exists; Derived from alone also
accepts a source fragment without domain interpretation. External and unacquired
references are preserved with advice. Incompatible or broken local references are
errors. Refinement cycles produce errors and withhold every cycle-participating
edge from the resolved graph; references and source content remain intact. Other
semantic relationships may cycle. Offline validation replays these projections.

An Obligations section on an ObligationContainer maps to `obligations_content` and
ordered owned `ObligationClause` records, one per immediate child heading. The
clause's introductory content, excluding optional **ID:** and **Quality attributes:** fields, is its
`statement`; its heading is its `name`. IDs follow the common identifier grammar
and document-wide uniqueness rule. Same-file and cross-file reference checks
recognize them. Named clauses can contain trace sections, Rationale, Applies to
(with the same explicit Targets field as Requirement), and structured Verification plan, Criterion and Measurement plan.
Quality attributes is a comma-separated list of system-declared IDs and resolves
through the clause’s declaring owner, using the same rules as standalone classification.
Unknown sections remain ExtensionSection records. No obligation is inferred from
a scenario step, state invariant, sentence, or evidence link.

Typed fields remain authoritative. Supports, DerivedFrom, Refines, ReliesOn,
Involves and HasObligation assertions are optional projections, and AppliesTo also
supports clauses. Required `projection_path` identifies the source field. There
are no stored semantic inverse fields. Missing explicit motivation and verification
guidance on standalone requirements and clauses produce advisory findings, not
shape failures or claims of missing required behavior.

`model trace` reads the same snapshot to render motivation, participation,
change-review and verification-guidance views. It makes no implementation,
satisfaction, validation or execution-evidence claims. Type-level diagrams remain
in the generated core model reference.


## Explicit specification-system selection

On SystemContextSpecification document kinds, a `Specification system` H2 maps to
`specification_system_content` and `specification_system_references`. Exactly one
Markdown link must identify a repository-local System, using its file, title
heading or authored ID. Images and footnotes remain content references, not
selection links. Empty, ambiguous or invalid selections remain explicit and never
fall back. Resolved selections populate scalar `specification_system`; when the
section is absent, an acquired directory `system_context` supplies that value.
`specification_system_origin` is `explicit` when a selection was authored, otherwise
`directory` when a default was acquired, otherwise absent.

Directory context is recorded only when acquired within the scan boundary. An
explicit selection does not require reading an unselected directory default.
Quality labels follow the effective association. Embedded clauses follow their
declaration chain to the selected specification system or declaring System, not
`applies_to`, `part_of` or feature membership. Structural and feature projections
are replayed during offline validation; retained source occurrences remain
available even when invalid semantic edges are withheld.
