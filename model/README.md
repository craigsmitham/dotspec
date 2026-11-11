# Specification model

[dotspec.yaml](dotspec.yaml) is the canonical LinkML model for experimental
specification snapshots. The model version is `0.4.0-draft`; the JSON
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

- `ancestors` includes the class itself followed by its base classes.
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
Associations use plain arrows, without inferred lifecycle or enforcement meaning.

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
