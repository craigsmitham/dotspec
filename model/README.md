# Specification model

[metamodel.yaml](metamodel.yaml) defines the LinkML domain classes, snapshot representation and report types. The [Markdown profile](markdown-mapping.md) describes the authoring vocabulary. The CLI embeds the generated schema and runtime catalog for offline export and validation.

## Regeneration

Run `mise run model-setup` to install the pinned generator environment. After editing the model, run `mise run model-generate`, `mise run model-test` and `mise run model-check`. Commit the regenerated schema, catalog and [reference pages](../docs/cli-handbook/reference/model/index.md) with the model change.

## Runtime catalog contract

The catalog lists inherited slots, record classes, allowed endpoint kinds, declaration memberships and the relationship fields available to views. The schema validates the serialized snapshot. Semantic link and classification roles have endpoint contracts in the catalog. Their resolved occurrences supply the endpoint in JSON. The [core diagram](../docs/cli-handbook/reference/model/core.md) describes the domain. The [consumer guide](../docs/cli-handbook/reference/specification-model.md) describes reading an export.
