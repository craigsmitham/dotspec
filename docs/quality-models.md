# Populate system quality attributes from a model

Use a reusable quality model when you want a starting set of concerns to tailor
for a system. Dotspec copies selected definitions into the system's existing
`Quality attributes` section. The resulting declarations are ordinary Markdown;
editing the source model or upgrading dotspec does not change them.

## Choose a starter

All three starters ship with the CLI and work offline. They include attributes and their
groupings, with original authoring prompts and source citations.
They are starting points for authoring, not authoritative reproductions of the
referenced frameworks or compliance checklists.

| Slug | Attributes | Use |
| --- | --- | --- |
| [`iso-25010-2023`](reference/quality-models/iso-25010-2023.md) | 49 (9 characteristics + 40 sub-characteristics) | Default: systematic product-quality concerns based on ISO/IEC 25010:2023 |
| [`furps`](reference/quality-models/furps.md) | 36 (5 categories + 31 distinct subcategories) | RUP FURPS variant; accuracy is grouped under both reliability and performance |
| [`q42`](reference/quality-models/q42.md) | 200 (9 dimensions + 191 qualities) | Complete distinct-quality inventory and overlapping dimension assignments at the pinned Q42 revision |

```sh
dotspec quality-model list
dotspec quality-model show furps
dotspec quality-model check iso-25010-2023
```

The FURPS source lists 32 subcategory placements; one shared accuracy attribute
preserves both placements. This starter follows the cited RUP version and does
not add FURPS+ constraints as quality attributes.

The Q42 starter pins upstream revision `0185215bfc133bac6a8091029b79b7104e47e351`.
Its 37 alias records are alternative-name prose on the canonical qualities.
Dimension tags become grouping references; requirements, approaches, and
standards remain available through the upstream links. The adapted Q42 starter
and its generated reference use [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/);
each definition retains attribution when copied. Updating dotspec never silently
changes a system's previously copied definitions.

`list` includes the bundled starters and project models discovered below `--path`.
`show` and `check` select a reserved built-in slug or a model in the selected
project directory. Built-in names are reserved; use a different slug for a local
adaptation. `check` validates structure and references and suggests missing
content. It does not assess a model's completeness or suitability for a system.

## Preview and copy

Create the system first if necessary. System initialization stays minimal;
it does not automatically add quality attributes.

```sh
dotspec system init
dotspec quality-model apply --system specs/system.spec.md --dry-run
dotspec quality-model apply --system specs/system.spec.md
```

Omitting the model selects `iso-25010-2023`. Omitting `--attributes` selects all
attributes. The dry run prints the complete resulting document; the real run
writes those same bytes, provided the inputs have not changed.

To choose another model and only some attributes:

```sh
dotspec quality-model apply furps --system specs/system.spec.md \
  --attributes reliability,performance --dry-run
```

Remove `--dry-run` to write. Paths supplied to `--system` are relative to
`--path`, or may be absolute. All referenced grouping attributes and their ancestors are included
automatically, including shared dimensions only once. Selecting a grouping alone
does not select its sub-attributes; select their IDs or omit `--attributes` to copy
the complete model. For example, selecting ISO `availability` copies `reliability`
and `availability`; selecting Q42 `resilience` copies `reliable`, `secure`, and
`resilience`. Definitions are appended in model declaration order, preserving
existing system prose, attribute order, file permissions, BOM, and line endings.

Any selected ID already present in the system causes the operation to fail
without writing, including a second application of the same selection. Unknown
IDs and an invalid resulting document also leave the file unchanged. Dotspec
never merges or replaces an existing declaration automatically.

Each copied definition receives a source link. Relative Markdown links and
images are relocated to retain their original targets, and reference-style links
are made self-contained. Code examples are preserved. A dry run makes these
changes reviewable before writing.

After copying, edit definitions to explain their meaning for your system. Add
quality requirements separately when you know the actual obligations, conditions,
and measures. Copying does not create requirements, thresholds, or evidence.

## Define or customize a model

```sh
dotspec quality-model init organization-quality
dotspec quality-model check organization-quality
```

Alternatively, start with a built-in using
`dotspec quality-model init organization-quality --from furps`. Both forms create
`specs/quality-models/organization-quality.quality-model.spec.md` without replacing
an existing file. Edit the name, description, edition, sources, and attributes.

```markdown
---
type: quality-model
id: organization-quality
---

# Organization quality

## Description

Concerns we review for internal developer tools.

## Edition

Team edition 1

## Sources

Describe or link the sources informing these definitions.

## Quality attributes

### Reliability

**ID:** reliability

Consistent operation under the conditions our developers encounter.

### Repeatability

**ID:** repeatability
**Sub-attribute of:** reliability

The same inputs and environment produce the same observable results.
```

Only `type: quality-model` and one nonempty H1 are required. An authored `id` is
optional. Conventional sections are optional, but cannot be repeated or placed
at the wrong heading level. Unknown sections and deeper definition content are
preserved. Each attribute needs a document-unique ID. Sub-attribute references stay
within the same model, can name several grouping attributes, and must not form
cycles. Use a comma-separated list of distinct IDs; omit the field for a root.
All attributes use level-three headings, regardless of grouping depth.

```markdown
### Resilience

**ID:** resilience
**Sub-attribute of:** reliable, secure

Consider continued useful service under adversity.
```

Here `reliable` and `secure` must also be declared in the model.
See [the grouping rationale](specification-model.md#quality-attribute-grouping)
for the distinction between organization and inherited obligations.

Older `Broader` or provisional `Parents` fields must be renamed to
`Sub-attribute of`. Checks report these old fields explicitly so a relationship
cannot silently become definition prose. JSON uses `subattribute_of` arrays;
computed requirement navigation uses `classification.ancestors`.

Model definitions are independent of enclosing systems. A system's requirement
labels resolve only against that system's declarations; they cannot implicitly
use attributes from a nearby model. Project `check`, `sources`, and `model export`
include quality-model files alongside other specification kinds. Bundled resources
are not automatically inserted into a project's source inventory.

All commands support `--json`. `apply --json` returns source, target, copied IDs,
dry-run state, and the complete resulting content. See the generated
[CLI reference](cli.md#dotspec-quality-model) for flags and
[core domain reference](reference/model/core.md) for model relationships.
