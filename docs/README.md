# Documentation map

| What do you need? | Start here |
| --- | --- |
| Decide how to express a user need, rule, interaction, or obligation | [Authoring handbook](handbook/README.md) |
| Look up a file format, command, check, or model field | [Reference](reference/README.md) |
| Follow one specification from intent to a running app | [To-do walkthrough and model views](../examples/todo/README.md) |
| Change a domain decision, CLI convention, or release | [Handbook maintenance](maintainers/handbook.md) · [CLI design](maintainers/cli-design.md) · [Releasing](maintainers/releasing.md) · [Contributing](../CONTRIBUTING.md) |

## CLI guide topics

Run `dotspec guide TOPIC` to open the page or section below. Use `dotspec guide --list` for the same lookup offline. A diagnostic's rule ID also opens its assigned topic. The [topic catalog](handbook/topics.json) defines the names, descriptions, destinations, and rule assignments. Edit that file and run `mise run generate` to refresh the table.

<!-- Generated topic table -->

| Topic | Use it to | Source |
| --- | --- | --- |
| `business-rule` | State business constraints and link their enforcement. | [`docs/handbook/business-rules.md`](handbook/business-rules.md) |
| `capability` | Describe an ability, its dependencies and its providers. | [`docs/handbook/jobs.md#connect-needs-to-capabilities-and-specifications`](handbook/jobs.md#connect-needs-to-capabilities-and-specifications) |
| `component` | Name a component, declare its parent and organize its local specifications. | [`docs/reference/formats/system.md#component-spec-format`](reference/formats/system.md#component-spec-format) |
| `concept` | Explain purposeful behavior through state and actions. | [`docs/handbook/concepts.md`](handbook/concepts.md) |
| `conventions` | Identify a project’s local specification authorities. | [`docs/handbook/project-conventions.md`](handbook/project-conventions.md) |
| `core-domain` | Inspect the core domain types, relationships and diagrams. | [`docs/reference/model/core.md`](reference/model/core.md) |
| `examples` | Browse the files bundled from the To-do example. | Listing built from [embedded example files](../examples/embed.go) |
| `feature` | Group stakeholder-visible behavior across specifications. | [`docs/handbook/relationships.md#describe-systems-components-and-features`](handbook/relationships.md#describe-systems-components-and-features) |
| `handbook` | Choose what to write using the runnable To-do example. | [`docs/handbook/README.md`](handbook/README.md) |
| `inspection` | Interpret check results, source discovery and references. | [`docs/reference/inspection.md`](reference/inspection.md) |
| `job` | Distinguish jobs, steps, outcomes and commitments. | [`docs/handbook/jobs.md`](handbook/jobs.md) |
| `layout` | Find a specification by source path or owning directory. | [`docs/reference/layout.md`](reference/layout.md) |
| `maintaining` | Audit domain choices, their sources and the effects of a change. | [`docs/maintainers/handbook.md`](maintainers/handbook.md) |
| `mapping` | Look up how authored Markdown maps to model elements. | [`model/markdown-mapping.md`](../model/markdown-mapping.md) |
| `obligation` | Locate the authoritative expression of an obligation. | [`docs/handbook/requirements-layers.md#distinguish-an-obligation-from-its-document`](handbook/requirements-layers.md#distinguish-an-obligation-from-its-document) |
| `quality` | Define a criterion and distinguish a measure from its required bound. | [`docs/handbook/requirements.md#define-or-reference-the-criterion-once`](handbook/requirements.md#define-or-reference-the-criterion-once) |
| `quality-model` | Select quality attributes and apply them to a system. | [`docs/handbook/quality-models.md`](handbook/quality-models.md) |
| `relationships` | Connect specifications through explicit relationships. | [`docs/handbook/relationships.md`](handbook/relationships.md) |
| `requirement` | Author assessable obligations and satisfaction criteria. | [`docs/handbook/requirements.md`](handbook/requirements.md) |
| `review` | Review a requirement set and surface unresolved decisions. | [`docs/handbook/review.md`](handbook/review.md) |
| `system` | Describe a system, its boundary, goals and quality attributes. | [`docs/reference/formats/system.md`](reference/formats/system.md) |
| `use-case` | Describe actor goals, boundaries and alternative paths. | [`docs/handbook/use-cases.md`](handbook/use-cases.md) |
| `verification` | Distinguish the four review questions and the evidence each needs. | [`docs/handbook/verification-and-validation.md`](handbook/verification-and-validation.md) |
