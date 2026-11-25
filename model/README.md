# Specification model

[metamodel.yaml](metamodel.yaml) defines the LinkML domain classes, snapshot representation and report types. The [Markdown profile](markdown-mapping.md) describes the authoring vocabulary. The CLI embeds the generated schema and runtime catalog for offline export and validation.

## Regeneration

Work from the root of a dotspec checkout with mise and the pinned tools installed. Run `mise run model-setup` to install the pinned generator environment. Edit `model/metamodel.yaml` for a type, field, or relationship. For a change to the selection or presentation in the core diagram, edit `model/views/core-domain.yaml`. Then run `mise run model-generate`, `mise run model-test` and `mise run model-check`. If `model-check` reports drift, inspect the change to its input, regenerate, and review the affected class or diagram. Commit the regenerated schema, catalog and [reference pages](../docs/reference/model/index.md) with the model change.

## Runtime catalog contract

The catalog lists inherited slots, record classes, allowed endpoint kinds, declaration memberships and the relationship fields available to views. The schema validates the serialized snapshot. For a semantic link or classification, the slot’s `range` or `any_of` declares its allowed endpoints. A `trace_heading` names the Markdown heading that supplies links. When a core diagram includes the slot, check that its node selection in `views/core-domain.yaml` covers the revised endpoints. In exported JSON, the resolved link or label occurrence identifies the endpoint. The [core diagram](../docs/reference/model/core.md) describes the domain. The [consumer guide](../docs/reference/specification-model.md) describes reading an export.
