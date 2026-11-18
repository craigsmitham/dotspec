# Reference

| What do you want to look up? | Where to look |
| --- | --- |
| Where does a specification belong? Which project does a command select? | [File layout](layout.md) · [Workspaces](workspaces.md) |
| Which Markdown fields and sections can I write? | [Specification formats](formats/README.md) |
| What arguments does a command accept? What output can I expect? | [CLI reference](cli.md) · [Command output examples](command-examples.md) |
| What does a diagnostic, reference, or source inventory establish? | [Inspection](inspection.md) |
| How does Markdown become an exported record? | [Specification model](specification-model.md) · [Markdown mapping](../../model/markdown-mapping.md) |
| Which types and relationships does the model define? | [Core domain](model/core.md) · [All fields](model/index.md) · [Relationships](model/relationships.md) |
| Which attributes come with a bundled quality model? | [ISO 25010](quality-models/iso-25010-2023.md) · [FURPS](quality-models/furps.md) · [Q42](quality-models/q42.md) |

## Changing a reference

Edit the linked format, layout, workspace, inspection, or mapping page directly to clarify its explanation. The following references are generated.

| Published material | Edit here | Regenerate with |
| --- | --- | --- |
| CLI reference and output gallery | [Cobra commands](../../internal/cli) and the [gallery scenarios](../../scripts/generate/main.go) | `mise run generate` |
| Quality-model starter catalogs | [Bundled starter specifications](../../internal/qualitymodel/specs/quality-models) | `mise run generate` |
| Core domain, fields, and relationships | [Metamodel](../../model/metamodel.yaml) and [core view definition](../../model/views/core-domain.yaml) | `mise run model-generate` |
| [To-do model views](../../examples/todo/README.md#read-the-model-views) | [To-do specifications](../../examples/todo/specs) for the example, [model renderers](../../internal/model) for presentation | `mise run generate` |
