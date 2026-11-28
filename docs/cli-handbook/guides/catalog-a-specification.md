# Catalog a specification

Choose the system, feature, component, or other specification whose decisions the source governs. In its body explain the promise in terms a reader of the product can understand. In its frontmatter give the native source a name readers will recognize.

```yaml
specifications:
  - name: Follow-up timing for waiting replies
    defines: requirements
    source: ../../acceptance/follow-up.feature
```

## Choose the source

Use a whole file when its scope matches the declaration. Use a fragment for one rule in a mixed suite, one invariant in a model, or one operation in a contract. Another subject can link directly to that native definition or to the catalog entry on the declaring spec, for example `waiting.feature.spec.md#specification.follow-up-timing-for-waiting-replies`. The declaration anchor is `specification.` followed by the heading-style fragment of the name.

The To-do project uses one Python function as the governing example for [saved follow-up dates across time zones](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md). Open its named declaration to find the native cases and their verification procedure.

Choose `defines` from `systems`, `components`, `jobs`, `capabilities`, `features`, `business-rules`, `use-cases`, `concepts`, `requirements`, `quality-attributes`, or `contracts`. A diagram can define allowed transitions, a CSV can define business rules, and a JSON schema can define a data contract. An example that helps explain the rule belongs in a contextual link. [Where does a product decision live?](../explanations/specification-sources.md)

## Name a definition

A local path is relative to the declaring Markdown file. A Markdown fragment can name a heading. For the following native formats dotspec can locate a named definition.

| Format | A fragment can name |
| --- | --- |
| Quint `.qnt` | A `module`, `const`, `var`, `type`, `val`, `def`, `pure val`, `pure def`, `action`, `run`, `temporal`, `assume` or `nondet` |
| Alloy `.als` | A `sig`, including each name in a comma-separated list and with an `abstract`, `one`, `lone`, `some` or `var` prefix, or a `module`, `enum`, `fact`, `pred`, `fun`, `assert`, `check` or `run`. For `pred Sig.name`, use `name`. |
| TLA+ `.tla` | The `MODULE`, each name after `CONSTANT`, `CONSTANTS`, `VARIABLE` or `VARIABLES`, a named `THEOREM`, `LEMMA` or similar, or a top-level `Name ==` or `Name(args) ==` |
| Gherkin `.feature` | A `Feature`, `Rule`, `Background`, `Scenario`, `Scenario Outline`, `Scenario Template` or `Example` title, written as a heading anchor. `Scenario: Resume a waiting entry` becomes `#resume-a-waiting-entry`, and a repeated title gains `-1`, `-2` and so on. |
| Go `.go` | A `func`, including a method, or a `type` |
| Python `.py` | A `def`, `async def` or `class` |
| txtar | A file name from a `-- name --` line |

Other native formats retain their own fragment semantics. The snapshot records those fragments as uninspected.

## Declare verification

List verification procedures on the specification they assess. For a native specification, nest them inside its entry. For a Markdown specification, put `verification` at the top level of frontmatter, as in the [restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md).

```yaml
verification:
  - name: Restart acceptance test
    source: ../../tests/acceptance.py#verify
```

After running a procedure, record its product revision, conditions, outcome and observation in the project's `verification-results.json`. Identify a Markdown subject by its project-relative path. For a native specification append the declaration anchor to the declaring Markdown path. The procedure is its own project-relative source URI including any fragment. When the native specification itself ran, use its governing source URI as the procedure. [Verification result format (dotspec checkout)](../../reference/generated-specification.md#record-verification-results)

## Inspect the catalog

From the project directory, run `dotspec build`. Open the system's **Specifications** index or the relevant collection, such as Requirements. Follow the named entry to its declaring subject and then to the governing source. Run `dotspec sources` to inspect the inventory. To start from the example suite, run `dotspec sources acceptance/follow-up.feature` from its project root.

Run `dotspec check` to locate missing files and recognized definitions. A whole-project scan suggests reviewing unlisted Quint, Alloy, TLA+ and Gherkin files. In `dotspec.yaml`, add project-relative patterns to include other candidates:

```yaml
specification_patterns:
  - "policy/*.csv"
```

Patterns use shell-style wildcards within a path segment. A match invites review. The author decides which source governs the behavior.

## Resolve a finding

- `catalog.invalid` identifies a malformed declaration. Use a list of mappings, nonempty single-line names and source paths, and a supported category for each specification.
- `catalog.duplicate` identifies names that produce the same anchor within a collection. Choose distinct names.
- `catalog.unlisted` identifies a candidate file for review during a whole-project check.
- `catalog.unchecked` identifies an external locator or a fragment whose native format dotspec does not inspect.
- `structure.verification-procedure` identifies a procedure listed under a body heading. Move it to frontmatter verification and keep the account of coverage under Evidence.
