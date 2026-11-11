# Markdown mapping profile

`dotspec-markdown/3` maps specification Markdown to the
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
| `type` and `id` | Interpreted document metadata; an accepted specification's explicit ID also maps to its root domain record. Systems have no authored-ID field in their authoring format. |
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
| Evidence | `EvidenceCollection` with full content, direct list-item statements and named `EvidenceStatement` records |
| Open decision; Open decisions; Open question; Open questions | `OpenDecisionCollection.content` and ordered `decisions` values; each question is distinct from its Evidence needed subsection |

Contextual references retain full explanatory prose, reference occurrences and
explicit list-item descriptions. They do not create support, authority,
applicability or satisfaction assertions. Evidence collections recognize
Research basis and Limitations subsections. Named evidence statements recognize
Type, Research status and Limitations. An authored hypothesis remains a claim;
the mapper does not assign confidence or prove it.

An open-decision section retains its complete narrative in the collection. A
singular heading maps its introductory question and explicit Evidence needed
subsection into one decision. A plural heading can contain named questions or
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
| Applies to | `applies_to: Applicability`, including the same explicit dimensions as requirements |
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
| Attribute ID, Sub-attribute of, definition, open decisions | Same syntax and content accounting as system attributes below |

A parent label on a model definition resolves inside that model and records
`quality_model_context`. Requirement classification continues to resolve only
against a system's declarations. Both owner types project
`DeclaresQualityAttribute` assertions. Copy provenance uses ordinary source links;
no synchronization, equivalence, or satisfaction relationship is inferred.

## System

| Section or item | Mapping |
| --- | --- |
| Purpose; Scope | `purpose`; `scope` |
| Environment | `EnvironmentDescription.content` |
| Goals | `GoalCollection.content`; direct list items and named subsections become `SystemGoal` records |
| Assumptions | `AssumptionCollection.content`; direct list items and named subsections become `Assumption` records |
| Constraints | `ConstraintCollection.content`; direct list items and named subsections become `Constraint` records |
| Glossary | `Glossary.content`; named subsections become `GlossaryTerm` records with full definitions |
| Term: Alternative names; Aliases | List-item values, or the unsegmented subsection, populate `alternative_names` |
| Term: Related terms | `related_terms: ContextualReferences` |
| Quality attributes | `System.quality_attributes_content` preserves the section; `System.quality_attributes` lists its attributes in authored order |
| Named attribute subsection | `QualityAttribute` with title, definition and declaring `System` owner |
| Attribute `**ID:**` | Required system-local `authored_id` |
| Attribute `**Sub-attribute of:**` | Ordered `subattribute_of_labels` from comma-separated IDs; resolutions supply multi-valued `subattribute_of` targets |
| Attribute: Open decision/question sections | Explicit `OpenDecision` values |

The attribute definition excludes its introductory ID and Sub-attribute of metadata lines
while retaining all other definition content. Subheadings that subdivide a
quality-attribute or glossary definition remain part of that narrative unless a
specific rule above assigns a more precise role. Unrecognized headings inside
structured detail fields, such as a glossary alias list, receive extension
records. An absent section leaves both system fields absent. An explicitly empty section retains empty content and an empty attribute list. Directory-derived parent systems
and specification ownership are associations; they introduce no inheritance.

## Job

| Section or item | Mapping |
| --- | --- |
| Job | `statement` |
| Performer; Situation; Motivation | `performer`; `situation`; `motivation` |
| Intended progress; Desired outcome | `intended_progress`, preserving the authored heading and existing duplicate-section validation |
| Steps | `JobMap.content`; named subsections become ordered `JobStep` records |
| Step `**ID:**`; `**Stage:**` | Optional `authored_id`; source-backed `stage` |
| Step introduction | `description`, excluding its explicit ID and Stage metadata lines |
| Outcomes under the job or a step | `OutcomeCollection` and direct list-item `Outcome` records with that actual owner |
| Outcome `**identifier:**` prefix | Optional outcome `authored_id`; complete item remains available as its statement |
| Nested outcome support content | `supporting_content` on the outcome |
| Progress cases | `ProgressCaseCollection.content`; named subsections become `ProgressCase` records |
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
| Goal; System boundary; Trigger; Success criterion | `goal`; `system_boundary`; `trigger`; `success_criterion` |
| Primary actor | `ActorMention` with full content and reference occurrences |
| Goal level | `GoalLevel.content`; recognized explicit wording may also populate canonical `level` |
| Preconditions; Minimum guarantees | `ConditionCollection` with full content and explicitly separated list items or headed conditions |
| Supporting actors | `ActorCollection` with full content and list-item or named actor mentions; named Services subsections populate `services` |
| Stakeholder interests | Full collection content and explicit list-item or named stakeholder-interest entries |
| Main success scenario | `Scenario` with full narrative; direct ordered-list items become `ScenarioStep` records |
| Step `**Actor:**`, `**Action:**`, `**Response:**`, `**Outcome:**` | The corresponding explicit step fields, when authored |
| Extensions | Full `ExtensionCollection.content` and named or explicitly labelled branch records |
| Sequence | `SequenceStatement.content`; explicit Mode and list-item constraints as described below |
| Related jobs | `related_jobs: ContextualReferences` |

Scenario steps retain their authored numbers or labels and source order. A
paragraph-only scenario remains a valid scenario narrative with no invented
steps. Actor, action and result clauses are never inferred from the grammar of
a prose sentence.

Each extension retains its full content and contextual references. Named
subsections or introductory explicit fields can supply Condition, At step,
Response and Continuation. Location and Point of occurrence headings also map
to location content. Location wording remains authored content; phrases such as
“during steps 2–3” do not automatically become resolved step targets.

A Flow subsection supplies ordered alternate steps. A Handling subsection or
`**Handling:**` field can explicitly state `specified` or `unresolved`; an
Unresolved handling heading also records the latter state. Omission means no
explicit handling classification. It does not mean that the author declared
handling unresolved. Other named branch subsections and explicit branch-labelled
list items retain nested branch ownership. Branch labels include forms such as
`2a` and `*a`.

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
| State: Invariants | Explicit condition values in `invariants` |
| State component: Value; Values; Type | `value_description` |
| Actions | `ActionCollection.content`; named subsections become `ConceptAction` records |
| Action: Inputs; Outputs; Effects; Rejections; Rejection behavior | Corresponding action fields; both rejection headings map to `rejections` |
| Action: Caller assumptions; Assumptions | `caller_assumptions: ConditionCollection` |
| Action: Guards | `guards: ConditionCollection` |
| Action: Open decision/question sections | Explicit `OpenDecision` values |
| Operational principle | `operational_principle` |

An unsegmented Actions paragraph remains complete named actions content. Abstract
state is a domain description, without an inferred database schema. Caller
assumptions and enforced guards have separate fields.

## General and quality requirements

| Section or metadata | Mapping |
| --- | --- |
| Applies to | `Applicability.content` |
| Applies to: Conditions; Workload; Population; Exclusions | Corresponding applicability fields |
| Applies to: Subject; Subjects | `subjects: ContextualReferences` |
| Rationale | `rationale` |
| Frontmatter `attributes` | Every authored classification label occurrence; resolved quality attribute targets populate `attributes` |
| General requirement: Statement | `statement` |
| General requirement: Verification | `VerificationPlan.content` and recognized explicit plan fields |
| Quality requirement: Criterion | `QualityCriterion.content` and recognized explicit criterion fields |
| Quality requirement: Measurement | `MeasurementPlan.content` and recognized explicit measurement fields |

Verification recognizes Method or Procedure, Success criteria, Conditions, and
Evidence to retain or Evidence retention. Named Cases subsections supply
`VerificationCase` values; each can contain Conditions, Procedure, Expected
result and Evidence to retain.

Criterion recognizes Statement, Measure, Comparison, Target, Units, Aggregation,
Window or Observation window, Conditions, Workload, Population and Exclusions.
Explicit open-decision headings remain distinguishable from an omitted target.
Qualitative criteria remain valid; classification labels cannot supply missing
thresholds.

Measurement recognizes Sampling, Conditions, Procedure, Evidence to retain or
Evidence retention, Workload, Population and Exclusions. Full measurement
narrative remains available alongside those components, including procedural
qualifications, examples and evidence instructions.

## Bindings, fallback and completeness

Bindings record a rule under `dotspec-markdown/3`, an owner and field path, source
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
