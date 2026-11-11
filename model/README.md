# Specification model

[dotspec.yaml](dotspec.yaml) is the canonical LinkML model for experimental
specification snapshots. The model version is `0.11.0-draft`; the JSON
representation version is `0.1.0`. The model describes source bytes, document
content, domain records, references, coverage and diagnostics independently of
Markdown authoring checks. See the [consumer guide](../docs/specification-model.md)
for export and validation, and the [Markdown mapping profile](markdown-mapping.md)
for recognized domain roles.

Generated projections are committed with the model:

- [schema.json](generated/schema.json) validates local JSON shape and concrete
  subtype fields.
- [catalog.json](generated/catalog.json) supplies inherited class/slot definitions
  and typed reference ranges to the runtime graph checker.
- [Relationship reference](../docs/reference/model/relationships.md) inventories typed edges, scope, provenance, inverse navigation and cycle rules.
- [Core domain reference](../docs/reference/model/core.md) provides selected domain types and generated Mermaid diagrams.
- [Specification model reference](../docs/reference/model/index.md) provides a
  complete class, field and enumeration lookup.

The Go package at this path is named `schema`. It embeds `SchemaJSON` and
`CatalogJSON` as read-only byte slices. The CLI does not invoke Python or LinkML.
Go snapshot records retain heterogeneous subtype fields instead of relying on
LinkML-generated Go structs.

## Regeneration

From the repository root, use the pinned mise toolchain and the explicit
setup task:

```sh
mise run model-setup
mise run model-generate
mise run model-check
mise run model-test
```

`model-setup` creates `.venv-model` with the pinned Python 3.11 interpreter and
installs the committed dependency lock. The launcher selects the correct Python
executable on Unix and Windows. Setup is a deliberate provisioning step;
`model-check` and `model-test` never install missing packages. The full
`mise run check` gate includes both checks.

For a separately managed environment, install
[scripts/generate-model/requirements.lock](../scripts/generate-model/requirements.lock)
into an isolated Python 3.11 environment, then use that interpreter to run
`scripts/generate-model/generate.py`, its `--check` form, and
`-m unittest discover -s scripts/generate-model -p 'test_*.py'`.

The generator checks every locked package version, loads only local source and
installed packages, and omits timestamps. `--check` compares all projections
without writing them. Installing the generation environment may use the network;
generation and runtime validation do not require it.

## Runtime catalog contract

`catalog_version: 1` identifies the catalog representation. The catalog records
the canonical schema ID, model and representation versions, source digest,
annotations, classes and enumerations. Class entries include every induced slot,
so inherited constraints and subtype narrowing are available without a second
hand-maintained endpoint table.

For each class:

- `ancestors` includes the class itself, base classes and mixed-in capabilities.
- `mixins` records explicitly declared reusable capabilities; `is_a` records the
  primary superclass.
- `abstract`, `description` and `annotations` retain model declarations. Typed
  assertion classes use the `projection_of` annotation where a relation is a
  materialized view of canonical fields. `owned_memberships` identifies canonical
  membership paths whose members must name the containing record as their owner.
- `rules` retains inherited conditional rules. A rule's presence does not mean
  every generator enforces it; the model explicitly records projection limits.

For each slot:

- `range` is the declared range; `any_of` narrows it when alternatives exist.
- `reference_ranges` contains allowed class ranges for record-ID references.
  A target's class must descend from one of these ranges.
- `inline_ranges` contains allowed embedded class ranges. Embedded values
  require recursive checking rather than record-ID lookup.
- `required`, `multivalued`, `inlined`, `inlined_as_list`, `ordered` and
  `identifier` describe the field shape. Applicable fixed values, enumerations,
  patterns and cardinality/value bounds are retained too.

JSON Schema cannot establish cross-record identity, endpoint existence,
provenance, source fidelity or coverage honesty. The canonical model identifies
those graph/content checks. It also records a specific LinkML 1.12.1 projection
gap: the conditional `RichContent` rule for empty versus authored node counts
requires a content check because generated JSON Schema omits those conditional
array bounds.

## Subsets, classifications, and documentation views

The canonical model declares seven named subsets through LinkML `subsets`.
Each subset has a description and one `classification` annotation: `core`,
`supporting`, or `generic`. Classification describes the subset's role in
**dotspec**, not the strategic importance of an individual authored requirement.

Every class declares explicit `in_subset` membership. Membership can overlap:
System and shared domain values participate in more than one requirements
perspective. Membership does not change inheritance, move a definition between
files, or require a classification field in authored specifications. Embedded
values and collection wrappers belong with the domain concepts they serve.
The complete reference lists all memberships by classification and subset.

[core-domain.yaml](views/core-domain.yaml) controls the
[core domain reference](../docs/reference/model/core.md): selected subsets,
types, fields, diagram nodes, and canonical field paths. It defines presentation,
not relationship semantics. The generator resolves paths through embedded
collection wrappers and derives endpoints and inheritance from the model.
Diagrams distinguish declaration containment (diamond), structural context and citations (dashed), semantic associations (plain), and inheritance (triangle). Containment describes declarations, not runtime lifecycle or enforcement.

`model-check` validates all declared subset references, allowed classifications,
class coverage, view selectors and diagram endpoints before comparing generated
files. Changing the model or the view requires regenerating and committing all
affected projections. A focused diagram need not show every type in a subset;
the complete reference remains exhaustive.

The `0.3.0-draft` model adds BusinessRule and the contextual `business_rules` field.
It also places quality attributes directly on System, removing the former quality
vocabulary wrapper and its declaration assertion. `DeclaresQualityAttribute`
now connects System to QualityAttribute; label lookup retains only system context.
Mapping profile `dotspec-markdown/2` recognizes these mappings and the new family
and section. The CLI's system JSON likewise exposes `quality_attributes` and
`quality_attributes_content` directly, and classification output uses `system_path`
in place of `model_path`. Consumers of these experimental fields must update. Saved
snapshots retain their declared versions; validation does not silently migrate
older snapshots. Use the corresponding prior CLI for an older model version,
or export the source again with the current CLI.

The `0.4.0-draft` model and mapping profile `dotspec-markdown/3` add standalone
QualityModel definitions. QualityAttribute ownership and DeclaresQualityAttribute
support System or QualityModel. Sub-attribute labels on model definitions record
`quality_model_context`; requirement classifications continue to use system-owned
attributes exclusively. Model description, edition, sources, attribute narrative,
and unknown sections retain the same complete-content representation.

This draft also renames quality grouping fields to `subattribute_of_labels` and
`subattribute_of`, with `SubattributeOf` assertion projections. References remain
multi-valued and acyclic; the inverse sub-attributes and transitive ancestors are
derived. Markdown uses `Sub-attribute of`; the older `Broader` and provisional
`Parents` spellings produce a migration diagnostic. CLI classification JSON uses
`ancestors` instead of `broader`. See the [grouping rationale](../docs/specification-model.md#quality-attribute-grouping).

The `0.5.0-draft` model and mapping profile `dotspec-markdown/4` simplify use-case
narratives. UseCaseExtension `handling` replaces `response` and `continuation`, `steps`
replaces `flow`, and `at_step` and `handling_state` are removed along with the
HandlingState enum. ScenarioStep retains label, content and references; separate
actor, action, response and outcome fields are removed. Both types use ordinary
`references` instead of `governing_references`. Extensions also retain explicit
`authored_label` values. Re-export Markdown to migrate; saved JSON is not silently
rewritten. Existing Response/Flow/Continuation Markdown remains supported.

Nested branches now require explicit branch labels or an Extensions grouping;
ordinary detail headings remain extension sections. See the
[use-case rationale](../docs/specification-model.md#use-case-narratives) and
[mapping rules](markdown-mapping.md#use-case).

The `0.6.0-draft` model and mapping profile `dotspec-markdown/5` add explicit
concept types, named system uses, type bindings, synchronizations, qualified action
references, and application dependencies. Dependencies may cycle. Concepts retain
optional initial and derived state, named action parameters/results and initiation
narrative. Singular and plural Operational principles headings map to the same rich
content field; they cannot both appear in one document. Existing prose remains
valid and no synchronization semantics are inferred.

`dotspec model docs` generates concept references, system compositions and dependency
diagrams. `mise run generate` maintains the complete example at
[Concept design example](../docs/reference/concept-design-example.md). Re-export
older snapshots with the new CLI; saved JSON is not silently migrated. See
[concept authoring](../docs/develop-concepts.md) and the
[concept design rationale](../docs/specification-model.md#concept-definitions-and-composition).


The `0.7.0-draft` model and `dotspec-markdown/6` consolidate relationships:

- `system_context` replaces `owning_system` on SystemContextSpecification. System
  and QualityModel no longer inherit this association; System uses `parent_system`.
  `InSystemContext` replaces OwnedBySystem.
- Every concrete assertion is a projection with required endpoints and
  `projection_path`. Canonical typed fields remain authoritative, including when
  assertion views are omitted. Existing snapshots must be re-exported.
- BoundaryIs, Supports and AppliesTo have explicit Markdown authoring forms and
  canonical fields. AppliesTo accepts BusinessRule as well as Requirement subjects.
- Record-reference slots declare relationship categories, scope, inverse, ordering,
  cycle policy and provenance basis. These edge categories do not alter subset
  classifications. The generated relationship reference presents this contract.
- Reciprocal `paired_with` annotations govern sibling declaration/content fields;
  embedded collection wrappers continue to require their full content. Optional
  decomposition uses `requires_content` so prose does not require invented records. Unknown
  definition scope defers dependent relationship resolution.

The catalog generator rejects incomplete relationship metadata, missing projection
paths, optional assertion endpoints and invalid declaration pairs. Runtime checks
validate realized relationships and source fidelity. See
[Declare relationships](../docs/develop-relationships.md).

The `0.8.0-draft` model and `dotspec-markdown/7` add layered traceability:

- `TraceableElement` separates motivation, derivation and assumption dependence
  from document kind and subset classification. `Requirement` continues to mean
  standalone obligation specifications; use cases and concepts can contain
  authoritative obligations without inheriting that document family.
- `supports`, `derived_from`, `relies_on`, obligation `refines`, and UseCase
  `involves` retain explicit source references and narrative rationale. Refinement
  cycles are authoring errors; cyclic edges are withheld while their source and
  diagnostics remain in a valid snapshot. Other semantic associations may cycle.
- Optional `ObligationClause` records distinguish obligations inside existing
  specifications. Heading anchors and optional explicit IDs address them. This
  is deliberate decomposition, not inferred acceptance or NLP extraction.
- `dotspec model trace` generates motivation, behavioral participation,
  direct/transitive change-review candidates, and verification-guidance views.
  Missing motivation and check guidance are advisory. Execution evidence and
  satisfaction judgments remain outside this implementation.

Re-export source with this version; saved snapshots are not silently migrated.
See [Requirements layers](../docs/requirements-layers.md),
[authoring relationships](../docs/develop-relationships.md), and the
[generated traceability example](../docs/reference/requirements-traceability-example.md).


The `0.9.0-draft` model and `dotspec-markdown/8` unify obligation capabilities:

- `Obligation` gives both standalone requirement forms and embedded clauses a
  statement, quality classification, applicability, refinement and assessment plans.
  Criterion and Measurement plan are optional assessment detail; Statement carries
  the obligation itself.
- `ObligationContainer` explicitly grants clause containment, separately from
  `TraceableElement`. Goals, assumptions and desired outcomes retain trace links
  without gaining nested obligations.
- Canonical names distinguish meaning: `quality_attributes`, `applies_to` targets
  versus `applicability` narrative, `subject_system` versus `boundary_description`,
  `verification_plan`, `measurement_plan`, `operational_principles`, `DesiredOutcome`,
  `UseCaseExtension`, `OpenQuestion`, `bound_type` and `supplying_concept_use`.
- Use cases use the shared Supports heading and relationship. Type bindings use
  Bound type and Supplying concept use fields; synchronization entries use Concept use.
- Generated diagrams use schema-owned human labels and explicit cardinalities,
  with quality classification separated from behavior and source machinery.

There are no compatibility aliases. Re-export Markdown using the current profile.


The `0.10.0-draft` model and `dotspec-markdown/9` introduced explicit specification subjects:

- `System.part_of` and `Component.part_of` identify optional single structural
  parents. They must be acyclic. `enclosing_system` replaces `parent_system` for
  directory-discovery context.
- `System.components` and `System.features` own inline declarations, while
  `Feature.groups` is non-owning, many-to-many specification membership.
  A feature belongs to its declaring system. Components reuse `ObligationContainer`.
- `Specification system` explicitly selects the associated system. The effective
  `specification_system` and its `explicit` or `directory` origin are separate from
  `system_context`, which retains discovery provenance. Invalid selections never
  fall back; quality classification follows the effective association.
- Systems, components and features support authored IDs. Descriptions and links to
  native contracts, architecture and implementation artifacts supply context;
  normative commitments remain requirements or obligation clauses.
- `dotspec model overview` generates explicit structure, feature membership,
  directly applicable obligations, concept uses, verification guidance and native
  references. Component involvement in features is derived through grouped
  obligations' applicability. It does not assert implementation, allocation,
  inherited requirements or satisfaction.

See [authoring relationships](../docs/develop-relationships.md#describe-systems-components-and-features)
and the [generated overview](../docs/reference/system-overview-example.md).


## Jobs, user needs and capabilities

Version `0.11.0-draft` and mapping `dotspec-markdown/10` add:

- Shared `DesiredOutcome.applies_to` relationships to jobs and steps, independent
  of the single declaring owner. Named outcomes have short heading labels,
  explicit IDs and full statements; ordinary list outcomes remain supported.
- Case-specific circumstances, struggling moments, participants, alternatives,
  and trade-offs/selection criteria on `ProgressCase`, reusing existing evidence.
- Optional System-owned `Capability` declarations with explicit `supports`,
  `depends_on`, and `provided_by` relationships and narrative `enablers`.
  Capability dependencies are distinct from structural part-of relationships;
  declaring systems are not inferred providers.
- General, ODI, progress and job-story authoring/check presets. These are CLI
  guidance, not domain classes or a method classification stored on jobs.
- `dotspec model needs`, a generated needs-to-specification diagram and source
  reference with full statements and evidence.

The legacy singular `Desired outcome` job heading is extension content rather
than an alias for `Intended progress`. Research assessments remain dated evidence,
not outcome properties. See [Jobs and user needs](../docs/jobs-and-user-needs.md).
