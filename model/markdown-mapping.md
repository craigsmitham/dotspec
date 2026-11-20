# Markdown mapping profile

`dotspec-markdown/17` maps specification Markdown to the
[canonical model](metamodel.yaml). This reference defines the recognized syntax
and domain roles. For reading or validating an exported snapshot, start with
[Specification model](../docs/reference/specification-model.md).

The profile covers the following specification types: system, component, feature, capability, job, quality attribute, quality model,
business rule, use case, concept and requirement. A specification is one authored
`*.spec.md` document, and its root record is the one concrete class of its type:
`System`, `Component`, `Feature`, `Capability`, `SystemQualityAttribute`, `Job`, `QualityModel`, `BusinessRule`, `UseCase`, `Concept` or
`Requirement`. Existing type parsers
remain authoritative for authoring structure and validity. The exporter retains
source and syntax before invoking those parsers. Rejection produces diagnostics
and failed domain decomposition without manufacturing a valid domain object.

## Shared mapping rules

| Authored structure | Mapping |
| --- | --- |
| Original file | `SourceDocument` with original bytes, digest and byte count; all spans address those bytes. |
| YAML frontmatter | Ordered YAML nodes with parsed scalar values, tags, aliases, anchors, comments and original spelling. Invalid YAML remains opaque content with a diagnostic. |
| `type` and `id` | The authored type is `<subject>-spec` in kebab case (for example, `use-case-spec`). The exporter preserves it as `declared_type`; an accepted specification's explicit ID also maps to its root domain record. System IDs are optional and document-local. |
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

These roles apply where authored. Some are recognized only on a specification,
the root record of a document. The others are recognized beneath the heading of
any element that has the field, including elements declared inside a
specification.

| Section | Recognized on | Domain field or structure |
| --- | --- | --- |
| Related specifications | Specification | `related_specifications: ContextualReferences`; a citation does not assert support, refinement, enforcement or compliance |
| Authority; Authority and evidence | Specification | `authority_and_evidence: RichContent` |
| Evidence | Specification; progress case | `evidence: EvidenceCollection` with full content, direct list-item statements and named `EvidenceStatement` records |
| Notes | Any element | `notes: RichContent` |
| Implementation suggestions; Implementation hints | Any element | `implementation_suggestions: RichContent`; nonbinding advice, not endorsement |
| Illustrative examples | Any element | `illustrative_examples: RichContent`, distinct from general obligations |
| Validation plan | Any TraceableElement | `validation_plan: AssessmentPlan`; explicit subject and basis distinguish requirements/system validation and assumption assessment |
| Supports; Derived from; Relies on | Any TraceableElement | Trace sections; see [Layered traceability and obligation clauses](#layered-traceability-and-obligation-clauses) |
| Open question; Open questions | Specification; any named element whose subsections are mapped, such as a declared component, job step, obligation clause, concept action, quality attribute or named assumption (a goal's or glossary term's subsections stay in its statement or definition); an obligation's Criterion | `open_questions_content` and ordered `open_questions` values on that record (the obligation, for a Criterion); each question is distinct from its Evidence needed subsection |

A specification's root record also carries its `preamble` and, for every type
except system, the `slug` taken from its file name. Related specifications is
the single contextual-reference section: cite a business rule, job, use case or
concept there, or use Derived from for a typed basis link. Beneath a ConceptUse, Synchronization or Interface, which are
elements and not specifications, Related specifications, Authority and Evidence
headings are retained as extension sections.

Contextual references retain full explanatory prose, reference occurrences and
explicit list-item descriptions. They do not create support, authority,
applicability or satisfaction assertions. Evidence collections recognize
Research basis and Limitations subsections. Named evidence statements recognize
Type, Research status and Limitations. An authored hypothesis remains a claim;
the mapper does not assign confidence or prove it.

An open-question section retains its complete narrative in
`open_questions_content`. A
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
| Applies to | `applicability: Applicability`, including the same explicit dimensions and `**Targets:**` field as requirements |
| Rationale | `rationale: RichContent` |
| Authority; Authority and evidence | Shared `authority_and_evidence: RichContent` |
| Exceptions | `exceptions: RichContent`, retaining qualifications and references |
| Refines | `refines_content`, `refines_references`; resolved `refines` targets business rules only |
| Related specifications | Shared contextual-reference content |

Absent and empty exceptions remain distinct. Unknown sections remain visible
extensions. A requirement relates to a rule through Derived from, through
Refines on the requirement, or through a citation under Related specifications.
A business rule does not contain obligation clauses. Constraints in a system
specification remain narrative; identifying a BusinessRule is an explicit
authoring decision.

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
`governing_quality_model`. Obligation classification continues to resolve only
against a system's declarations. Assertions project both memberships, as
`QualityModel.quality_attributes` and `System.quality_attributes`. Copy provenance uses ordinary source links;
no synchronization, equivalence, or satisfaction relationship is inferred.

## System

| Section or item | Mapping |
| --- | --- |
| Frontmatter `id` | Optional System `authored_id`, document-local anchor |
| System, Component, Feature, Capability, or ConceptUse: direct Description section | Optional `description: RichContent` on the declaring subject, with bindings to the authored heading and content. |
| Purpose; Scope | `purpose`; `scope` |
| System or Component: Out of scope | `out_of_scope: RichContent` |
| System or Component: Interacts with | Authored content and reference occurrences; resolved `interacts_with` System/Component targets, independent of `part_of` |
| System or Component: Required interfaces; Provided interfaces | `required_interfaces_content` or `provided_interfaces_content` with the full section, plus ordered `required_interfaces` or `provided_interfaces` membership of owned `Interface` elements |
| Named interface | Heading `name`, complete `content`, introductory `description`, declaring System/Component `owner`, optional introductory `**ID:**`; Contracts and Counterpart subsections are contextual references; Notes, Open questions, Validation plan, trace sections and Obligations are supported |
| Project context | `project_context: ContextualReferences`; authoritative project artifacts remain external |
| Part of | `part_of_content`, explicit `part_of_references`, resolved scalar System `part_of`; no inferred parent |
| Feature or capability source | Root Feature or Capability with a name from its H1, optional frontmatter `id` and nearest `specification_owner`. |
| Component/Feature/Capability: Purpose | `purpose` |
| Component: Responsibilities; Boundary | `responsibilities`; `boundary_description` |
| Component: Part of | Narrative and explicit references; resolved scalar System or Component `part_of` |
| Component: Contracts; Architecture; Implementation | `ContextualReferences` preserving narrative and links to native artifacts |
| Component: Obligations | Shared named `ObligationClause` declarations; no implicit applicability |
| Feature: Groups | Narrative and explicit references; resolved non-owning `groups` membership in requirements, clauses, use cases, concepts, concept uses, synchronizations or interfaces |
| Capability: Supports; Depends on; Provided by | Explicit semantic associations to user needs/goals, Capability dependencies, and System/Component providers; declaration does not imply provision |
| Capability: Enablers | `enablers` narrative, including human/procedural providers and external references |
| Environment; Stakeholders; Constraints | Narrative `environment`, `stakeholders` and `constraints`; list items and subordinate headings remain part of the narrative |
| Goals | `goals_content`; direct list items and named subsections become `Goal` records in `goals` |
| Assumptions | `assumptions_content`; direct list items and named subsections become `Assumption` records in `assumptions` |
| Assumption: Basis; Applicability; Validation plan | `basis`, prose `applicability`, and `validation_plan: AssessmentPlan`; statement contains the named assumption introduction only |
| Glossary | `glossary_content`; named subsections become `GlossaryTerm` records in `glossary` with full definitions |
| Term: Alternative names; Aliases | List-item values, or the unsegmented subsection, populate `alternative_names` |
| Term: Related terms | `related_terms: ContextualReferences` |
| Quality attributes | `System.quality_attributes_content` preserves the section; `System.quality_attributes` lists the declarations in its quality tree |
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
and a specification's `specification_system` are associations; they introduce no inheritance.

A component is the root of a named `component-spec` document. Its required `Part of` link supplies its structural parent. Local specifications record their nearest system or component in `specification_owner`. An explicit classification system can differ from this owner. Systems and components can declare concept uses, synchronizations and interfaces. There
is no constraint record: an individually addressable constraint is a requirement
or a business rule.

## System quality attribute

A `quality-attribute-spec` root maps to `SystemQualityAttribute`, with a
`definition` from the Definition section. The nearest declared subject must be a system to own an attribute.
A `subattribute_of` list names each parent. A requirement’s `quality_attributes`
list names each classification. Both fields produce authored label occurrences.


## Job

| Section or item | Mapping |
| --- | --- |
| Job | `statement` |
| Performer; Situation; Motivation | `performer`; `situation`; `motivation` |
| Intended progress | `intended_progress`; the old singular Desired outcome heading is extension content |
| Steps | `steps_content`; named subsections become ordered `JobStep` records in `steps` |
| Step `**ID:**`; `**Stage:**` | Optional `authored_id`; source-backed `stage` |
| Step introduction | `description`, excluding its explicit ID and Stage metadata lines |
| Outcomes under the job or a step | `outcomes_content` and direct list-item `DesiredOutcome` records in `outcomes`, with that actual owner |
| DesiredOutcome `**identifier:**` prefix | Optional outcome `authored_id`; complete item remains available as its statement |
| Named outcome with introductory `**ID:**` | `name` from its heading, `authored_id`, and full `statement` from introductory prose excluding ID; headings without IDs remain extensions |
| Named outcome: Applies to | `applies_to_content`, authored reference occurrences, and resolved Job/JobStep `applies_to` targets, independently of declaration owner |
| Named outcome: Supporting content | `supporting_content` |
| Nested outcome support content | `supporting_content` on the outcome |
| Progress cases | `progress_cases_content`; named subsections become `ProgressCase` records in `progress_cases` |
| Case: Circumstances; Struggling moment; Participants | `circumstances`; `struggling_moment`; `participants` |
| Case: Alternatives considered; Trade-offs and selection criteria | `alternatives_considered`; `tradeoffs_and_selection_criteria` |
| Case: Current approach; Desired progress | `current_approach`; `desired_progress` |
| Case: Forces | `forces_content` with the full narrative; explicit list-item or headed `Force` values in `forces` |
| Case: Timeline | Narrative `timeline`; list items and subordinate headings remain part of it, without inferred dates |
| Case: Evidence | The shared evidence structure |

A force receives `kind` only from an explicit push, pull, anxiety or habit label.
Other force descriptions retain their content without a fabricated category.
Stage labels remain open. An explicit empty Stage field remains distinguishable
from an omitted Stage field. Unknown detail headings under a step or progress
case become extension sections on the nearest addressable owner.
A timeline saying “to be established” remains authored content. Overall outcomes belong to the job and are not copied into every step.
Jobs and job steps do not contain obligation clauses.

## Use case

| Section or item | Mapping |
| --- | --- |
| Goal; System boundary; Trigger | `goal`; `boundary_description`; `trigger` |
| Primary actor | `primary_actor: ActorMention` with full content and reference occurrences |
| Goal level | `goal_level_content`; recognized explicit wording also populates canonical `goal_level` |
| Preconditions; Minimal guarantees; Success guarantees | `preconditions_content`, `minimal_guarantees_content` or `success_guarantees_content` with the full section; list items or headed conditions become `Condition` values in `preconditions`, `minimal_guarantees` or `success_guarantees` |
| Supporting actors; Stakeholder interests; Sequence | Narrative `supporting_actors`, `stakeholder_interests` and `sequence` |
| Step: nested `Illustrates:` bullet | `illustrates_content`, `illustrates_references`; resolved `illustrates` Obligation, BusinessRule or ConceptAction targets; ordinary step links remain citations |
| Main success scenario | `Scenario` with full narrative; direct ordered-list items become `ScenarioStep` records |
| Step `**Actor:**`, `**Action:**`, `**Response:**`, `**Outcome:**` | Preserved in step `content`, without separate semantic fields |
| Extensions | `extensions_content` with the full section; named or explicitly labelled branch records in `extensions` |
| System boundary: explicit `**Subject system:**` link | `subject_system_references`; resolved `subject_system: System` |
| Supports | `supports_content`, `supports_references`; resolved `supports` (Goal, Job, JobStep, DesiredOutcome or UseCase) |

Scenario steps retain their authored numbers or labels, narrative, ordinary
references and source order. Ordered-list items and explicitly labelled handling
steps such as `2a1` or `.1` become `ScenarioStep` records. A paragraph-only scenario
remains valid without invented steps. Actor, action and result wording remains
in the narrative, including explicitly labelled fields from older documents.

Each extension retains its full `content` and ordinary `references`. An optional
Condition subsection or introductory field populates `condition`. Where the
condition arises stays in the branch label or narrative; “during steps 2–3” does
not resolve a step target.

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
`1a2a` or `.1a`) or appear under an Extensions subsection. A Notes subsection is
the branch's `notes`; other detail headings become `ExtensionSection` records. Explicitly branch-labelled list
items also create branches, including nested items under a handling step. Labels
include `2a`, `*a`, `1a2a` and `.1a`; `authored_label` preserves the label while
Markdown containment establishes ownership. Labels do not determine control-flow
edges, durable identity or return targets. Sibling branches retain source order.

A Sequence section states, as narrative, which ordering is representative or
required. Scenario step order alone creates no constraint, and no ordering edges
are derived. Canonical goal levels recognize
the following wording after trimming whitespace, lowercasing and removing a
final period; other wording remains in `goal_level_content` only.

| Canonical goal level | Recognized wording |
| --- | --- |
| `summary` | `summary`, `kite`, `summary (kite)` |
| `user_task` | `user task`, `user-task`, `sea`, `sea level`, `user task (sea level)` |
| `subfunction` | `subfunction`, `fish`, `subfunction (fish)` |

## Concept

| Section or item | Mapping |
| --- | --- |
| Purpose | `purpose` |
| State | `abstract_state: StateDescription` with full `content`; direct list items and named state components retain descriptions |
| Type parameters; Introduced types | `type_parameters_content` or `introduced_types_content` with the full section; named `ConceptType` records in `type_parameters` or `introduced_types` |
| State: Initial state | `abstract_state.initial_state` narrative |
| State: Invariants | Explicit condition values in `invariants` |
| State component: Value; Values; Type | `value_description` |
| State component: Derivation | `derivation` narrative; no computed value is inferred |
| Effects table | `effects_table_content` with the full section; `effects_table` holds the ordered rows and rich-content cells of the first table; first row is the header, first column identifies operations; blank cells remain unspecified |
| Actions | `actions_content`; named subsections become `ConceptAction` records in `actions` |
| Action: Inputs; Outputs; Initiation; Effects; Rejections; Rejection behavior | Corresponding action fields; both rejection headings map to `rejections` |
| Action: Caller assumptions; Assumptions | `caller_assumptions_content` with the full section; list items or headed conditions become `Condition` values in `caller_assumptions` |
| Action: Guards | `guards_content` with the full section; `Condition` values in `guards` |
| Action: Open question/question sections | Explicit `OpenQuestion` values |
| Operational principles | `operational_principles`; named subsections may distinguish individual stories |
| Named Inputs/Outputs subsections | `ActionParameter` records in `input_parameters`/`output_parameters`, preserving full content |

An unsegmented Actions paragraph remains complete `actions_content`. Abstract
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
is inferred. The [authoring guide](../docs/handbook/concepts.md) supplies examples.

## Requirement

A `requirement` specification maps to `Requirement` and states one obligation. Functional and quality requirements share this type and can both reference
attributes through `quality_attributes`. A requirement does not contain an Obligations section. Obligation
clauses declared inside another specification share these fields; see
[Layered traceability and obligation clauses](#layered-traceability-and-obligation-clauses).

| Section or metadata | Mapping |
| --- | --- |
| Applies to | `applicability: Applicability` with full `content` |
| Applies to: Conditions; Workload; Population; Exclusions | Corresponding applicability fields |
| Applies to: explicit `**Targets:**` links | `applies_to_references`; resolved `applies_to: DomainElement` (also on BusinessRule) |
| Rationale | `rationale` |
| Frontmatter `quality_attributes` | `quality_attribute_labels` retains every authored classification label occurrence; resolved quality attribute targets populate `quality_attributes` |
| Statement | `statement` |
| Verification plan | `verification_plan: AssessmentPlan` for planned realization conformance |
| Validation plan | `validation_plan: AssessmentPlan` for needs/intent assessment; also recognized on other traceable records |
| Criterion (optional explicit satisfaction rule) | `criterion: ObligationCriterion` with full `content` |
| Criterion: Measure; Target | `measure`; `target` |

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
result and Evidence to retain. Measurement procedure and sampling belong under
the Verification plan's Method.

Criterion recognizes Measure and Target. The measure carries units, aggregation
and observation window as needed; the target carries the comparison and bound,
as in “At least 90%”. Operating context belongs to Applies to, not Criterion.
Explicit open-question headings remain distinguishable from an omitted target.
Qualitative criteria remain valid; classification labels cannot supply missing
thresholds. Define each required bound once in Statement or Criterion, or reference
its authoritative source. Criterion may reference Statement instead of restating
it. Target is a required bound or acceptable range, not an aspirational preference.
Keep preferred targets explicitly nonbinding in Rationale or Notes.

Applicability describes required operating context.
For outcomes of use, record users, tasks, assistance and environment in Applies to;
population and conditions preserve this context without
a separate context-of-use type. Plan conditions describe the assessment setup and its
relationship to that context. Plan Success criteria explain how observations are
interpreted against the obligation; reference its required bounds rather than
introducing competing thresholds. The mapper preserves authored content even when
it does not follow this guidance; it does not relocate or infer requirements.

Measurement-basis advice considers an explicit Measure even
when Verification plan is present. A Measure without a separately identified
target, or without conditions or applicability, produces advice; prose may already supply
them. Plan Success criteria and Conditions do not substitute for obligation
criteria and applicability. See [Develop requirements](../docs/handbook/requirements.md).

## Bindings, fallback and completeness

Bindings record a rule under `dotspec-markdown/17`, an owner and field path, source
occurrence and role. A field path is relative to its owner and uses canonical
field names. A section's narrative binds to `<field>_content`, or to
`<field>.content` on an embedded value such as `criterion.content`. Items of an
embedded list are indexed from zero, as in `open_questions.0.question` or
`effects_table.1.cells.2`. Named-field, record-identity and structured-item bindings
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
and the type-specific `mapping_*.go` files beside it. The
[mapping tests](../internal/model/mapping_test.go) exercise every checked-in
specification type and explicit structures; the
[export conformance tests](../internal/model/export_conformance_test.go) check
whole-snapshot behavior. Model projections are regenerated and checked as
specified in the [model README](README.md#regeneration).

## Relationship contract

See [Declare relationships](../docs/handbook/relationships.md) for complete examples.
Subject system references are scalar; Supports and Targets may name multiple
records. Ordinary links outside these explicit forms remain contextual. Definition,
parameter, type, use and action fields likewise retain all authored occurrences,
including malformed multiple references; an invalid draft must remain representable.
Only successfully resolved, compatible relationships populate semantic targets.
Unknown declaring scope defers dependent resolution rather than inventing a scope
error. Requirements, business rules, jobs, use cases and concepts record their
system in `specification_system`. An explicit Specification system section selects it;
otherwise the nearest enclosing system found by directory discovery supplies it.
`specification_system_origin` records which: `explicit` or `directory`. Systems use `enclosing_system` for discovery context
and explicit `part_of` for structural parentage; QualityModels are standalone.

Typed fields are canonical. An `Assertion` is an optional derived view of one
canonical field, with a required `projection_path`, subject and object. There is
one concrete assertion class. `projection_path` names the projected field, such
as `UseCase.subject_system`, `TraceableElement.supports`, `Obligation.applies_to`
or `BusinessRule.applies_to`, and `mapping_rule` is
`dotspec-markdown/17/projection/<path>`. The path also carries a membership's
role, as with `Concept.type_parameters` and `Concept.introduced_types`. The
generated catalog lists the projectable paths in `projections`: every
declaration membership and every semantic-association or structural-context
field. Embedded values such as evidence, applicability and assessment plans
retain required content; sibling
declaration/content pairs use reciprocal `paired_with` schema annotations.

In the canonical model a relationship field states its `relationship_category`
and only what differs from that category. The `RelationshipCategory` enumeration
supplies the default scope, cycle policy and basis. An unstated inverse is none,
an unstated label is the field name, and ordering follows the field shape. The
generated catalog and
[relationship reference](../docs/reference/model/relationships.md) show the
completed contract.

Optional input/output parameter and synchronization action collections require the
corresponding complete Inputs/Outputs/When/Then content through `requires_content`.
The converse is not required: prose-only sections remain valid. Explicit boundary
and applicability occurrences likewise require their narrative sections.

## Layered traceability and obligation clauses

`TraceableElement` supplies explicit Supports, Derived from and Relies on sections.
It is the base of Specification and of named Goal, Assumption,
JobStep, DesiredOutcome, ConceptAction and ObligationClause records. It is also
the base of the ConceptUse, Synchronization and
Interface elements declared inside a System specification; an Interface can also
be declared inside a Component. These declared elements are not specifications.
Notes, Open questions, Validation plan and trace sections work beneath each of
them; Evidence, Authority and Related specifications headings there are retained
as extension sections.
Heading anchors and optional introductory IDs identify these records. Interface
ownership is exactly one required or provided declaration under one System or
Component; assertions project those memberships through the
`required_interfaces` and `provided_interfaces` paths of the declaring class.
Native Contracts and Counterpart citations do not establish typed
implementation or compatibility edges. The exact heading annotations `trace_heading` and
`trace_targets` in the canonical schema drive mapping and endpoint resolution.
The mapper also recognizes these sections directly beneath named nested records.
It never interprets arbitrary prose or ordinary citations as semantic links.

Each trace section maps to `<field>_content: RichContent` plus ordered
`<field>_references`. Resolved targets populate the multivalued canonical field.
Supports targets Goal/Job/JobStep/DesiredOutcome/UseCase; Derived from targets DomainElement
or SourceFragment; Relies on targets Assumption. UseCase Involves targets ConceptUse.
Requirement and ObligationClause Refines targets Requirement,
BusinessRule or ObligationClause; BusinessRule Refines targets BusinessRule only.
Every section link is a declared target; rationale
prose remains in that section, and contextual citations belong elsewhere.

A ScenarioStep is a plain element, not a TraceableElement, and has no trace
sections. Its only typed link is the nested `Illustrates:` bullet, which targets
Obligation, BusinessRule or ConceptAction records.

A file reference promotes to a specification's root domain record. A fragment
promotes to its named domain record where one exists; Derived from alone also
accepts a source fragment without domain interpretation. External and unacquired
references are preserved with advice. Incompatible or broken local references are
errors. Refinement cycles produce errors and withhold every cycle-participating
edge from the resolved graph; references and source content remain intact. Other
semantic relationships may cycle. Offline validation replays these projections.

An Obligations section on an ObligationContainer maps to `obligations_content` and
ordered owned `ObligationClause` records, one per immediate child heading.
System, Component, Interface, UseCase, Concept, ConceptAction, ConceptUse and
Synchronization are obligation containers. Job, JobStep, BusinessRule,
Requirement, Feature, Capability and QualityModel are not: an Obligations
section there produces an `obligation.container` error, declares no clauses, and is
retained as an extension section. The
clause's introductory content, excluding optional **ID:** and **Quality attributes:** fields, is its
`statement`; its heading is its `name`. IDs follow the common identifier grammar
and document-wide uniqueness rule. Same-file and cross-file reference checks
recognize them. Named clauses can contain trace sections, Rationale, Applies to
(with the same explicit Targets field as Requirement), and structured Verification plan and Criterion.
Quality attributes is a comma-separated list of system-declared IDs and resolves
through the clause’s declaring owner, using the same rules as requirement classification.
Unknown sections remain ExtensionSection records. No obligation is inferred from
a scenario step, state invariant, sentence, or evidence link.

Typed fields remain authoritative. Assertions projecting
`TraceableElement.supports`, `TraceableElement.derived_from`,
`TraceableElement.relies_on`, `Obligation.refines`, `BusinessRule.refines`,
`UseCase.involves`, `ScenarioStep.illustrates` and
`ObligationContainer.obligations` are optional views, and `Obligation.applies_to`
covers clauses as well as requirements. Required `projection_path` identifies the source field. There
are no stored semantic inverse fields. Missing explicit motivation and verification
guidance on requirements and clauses produce advisory findings, not
shape failures or claims of missing required behavior.

`model trace` reads the same snapshot to render motivation, participation,
change-review and verification-guidance views. It makes no implementation,
satisfaction, validation or execution-evidence claims. Type-level diagrams remain
in the generated core model reference.


## Explicit specification-system selection

On every specification type except system and quality model
(`SystemContextSpecification`), a `Specification system` H2 maps to
`specification_system_content` and `specification_system_references`. Exactly one
Markdown link must identify a repository-local System, using its file, title
heading or authored ID. Images and footnotes remain content references, not
selection links. Empty, ambiguous or invalid selections remain explicit and never
fall back. Resolved selections populate scalar `specification_system`; when the
section is absent, the nearest enclosing system acquired by directory discovery
supplies that value.
`specification_system_origin` is `explicit` when a selection was authored, otherwise
`directory` when a default was acquired, otherwise absent.

Directory context is recorded only when acquired within the scan boundary. An
explicit selection does not require reading an unselected directory default.
Quality labels follow the effective association. Embedded clauses follow their
declaration chain to the selected specification system or declaring System, not
`applies_to`, `part_of` or feature membership. Structural and feature projections
are replayed during offline validation; retained source occurrences remain
available even when invalid semantic edges are withheld.
