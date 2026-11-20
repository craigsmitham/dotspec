# Requirements reference

Functional requirements live in their owner’s `specs/requirements/` folder.
Use `dotspec requirement init <slug> --attribute <id>` to record a related
quality in frontmatter and link the requirement from the attribute report.
For a quality obligation owned by the declaring system, use `--preset quality --attribute <id>` to create it beside the attribute definition. That folder
supplies its primary classification. Use `quality_attributes` for additional
IDs. Component-owned requirements stay under the component’s `requirements/`.

## Requirements

Use a requirement to state an obligation without first writing a concept
or use case. A requirement can express a response, postcondition, invariant,
relationship, or imposed restriction. Its type does not implicitly mean
“functional”; functional/non-functional classification is not required.
Functional and quality requirements share the same type, commands, and format. Both can reference quality attributes.

```sh
dotspec requirement init reject-existing-target
dotspec requirement init diagnostic-severity --preset quality
dotspec requirement show reject-existing-target --json
dotspec requirement list
dotspec requirement list --attribute reliability --json
dotspec requirement check reject-existing-target --json
```

Every command accepts `--path DIR`, selecting an existing project directory,
and `--json` or `--format json`. Paths, slugs, and output follow the same
conventions as other types.

- `init <slug>` creates `specs/requirements/<slug>.requirement.spec.md`
  without overwriting an existing file or following symlinked metadata directories.
  `--preset quality` selects a template that also prompts for a `Criterion`;
  `general` is the default. The preset is authoring guidance: either one creates
  an ordinary requirement, and the selection is not stored. At the declaring system, combining `--preset quality` with `--attribute` selects the attribute folder.
- `show <slug>` reads the requirement in that project, displays its authored
  Markdown and classifications, and reports its associated system and source.
  Slug lookup does not search ancestor projects.
- `list` discovers requirements in the selected subtree, grouped by associated
  system. `--attribute ID` includes exact labels and sub-attributes of
  that ID among the attributes each requirement's system declares. Each matching
  requirement appears once. When that system declares no quality attributes, only
  exact labels match. Same-named attributes in different systems remain separate
  definitions.
- `check <slug>` checks the requirement and its system context with the shared
  report format and exit behavior. Project-wide `dotspec check` includes requirements.

### Requirement format

A requirement spec states one obligation.

```markdown
---
type: requirement-spec
id: continue-after-invalid-spec
quality_attributes: [diagnosability]
---

# Continue checking after a structural error

## Applies to

Project-wide specification checks.

## Statement

When a discovered specification contains a structural error, dotspec must
continue inspecting the remaining discovered specifications.

## Rationale

One malformed specification must not hide problems in other files.

## Verification plan

Check two files with distinct structural errors and verify that both appear
in the report.
```

Only `type: requirement-spec` and one nonempty document H1 are required.
`id` is optional and follows the document-local ID rules. `quality_attributes` is an
optional list of unique, nonempty string IDs; an empty list is allowed. A scalar,
null, non-string element, invalid ID, or duplicate label is a structural error.
Unknown metadata fields are rejected.

`Applies to`, `Statement`, `Criterion`, `Rationale`, `Verification plan`, and
`Validation plan` are optional H2 sections. Missing statement, applicability,
rationale, or verification plan produces advice; `Criterion` and `Validation plan`
provide optional detail. Duplicate sections or invalid levels produce structural
errors. Deeper headings belong to their enclosing sections. Contract notation,
tables, additional sections, links, and other authored Markdown are preserved.
No sentence pattern or particular verification technique is required. Check
commands validate relative file links using the shared [reference rules](../inspection.md#references).
Model export reports an `Obligations` section inside a requirement as
`obligation.container`.

`Criterion` defines or references what counts as satisfaction. It can be prose,
or it can separate what is measured from the required bound:

```markdown
## Criterion

### Measure

The proportion of eligible task attempts in which the member restores the
intended document without external assistance.

### Target

At least 90%.
```

The measure carries units, aggregation, and observation window as needed; the
target carries the comparison and bound. Operating context belongs in
`Applies to`, which can use `Conditions`, `Workload`, `Population`, and
`Exclusions` subsections. Describe the measurement procedure and sampling in the
verification plan. When a criterion records a `Measure`,
`requirement.measurement-basis` advice appears if the criterion has neither
introductory prose nor a `Target`, or if `Applies to` has neither introductory
prose nor a `Conditions` subsection. If a target is unknown, explicitly record
the open decision and the evidence needed to choose it. The checker does not
evaluate prose, invent thresholds, assess test results, or decide requirement
satisfaction.

### Classify requirements by quality attribute

The associated system supplies quality attributes for label resolution.
This is independent of the subject and conditions stated in `Applies to`. If that
system has neither an attribute tree nor an explicit quality overview section, labels remain unbound without advice.
When the section is present, each undefined label produces
`requirement.undefined-attribute` advice at its source location. An explicitly empty
section consequently reports every used label as undefined. Statement and
rationale advice is worded for a quality obligation when a requirement has
quality attributes.

Use multiple attributes when they classify the same coherent obligation. Split
obligations that need independent management, even if they share an attribute.
Choose the most specific applicable labels; broad attributes can also be used.
Grouping beneath an attribute does not establish that the attribute is satisfied.

The show envelope uses `requirement`; the inventory uses `groups` containing
`requirements`, with `attribute` when one is selected. Both include
`schema_version: 1`. Each requirement exposes its
`statement`, `applicability`, `rationale`, `verification_plan`, `quality_attributes`,
`classification`, `markdown`, source identity, and associated system; `criterion`
and `validation_plan` appear when authored. Empty collections are arrays.

Requirement JSON lists the containing attribute followed by additional frontmatter IDs, once each. The snapshot separately preserves the authored label occurrences. The separate `classification`
object exposes `system_path` for the system supplying attribute declarations,
`resolved` attribute definitions, `ancestors` IDs in declaration order,
and `undefined` IDs in authored order. Without a system quality vocabulary,
`system_path` is empty and all three arrays are empty. Ancestor attributes are
computed navigation, not additional authored labels or inherited requirements.

Keep one authoritative statement of each obligation. Concepts organize related
rules and quality requirements; use cases illustrate interactions and link to
the rules that govern them. Preserve required logical dependencies without
turning a representative sequence into an obligation. Distinguish caller
assumptions from conditions the system must enforce or reject.

Examples are [existing-target rejection](../../../specs/requirements/reject-existing-target.requirement.spec.md),
[continued inspection](../../../specs/requirements/continue-after-invalid-spec.requirement.spec.md),
and [successful advisory checks](../../../specs/requirements/advice-permits-success.requirement.spec.md).
[Diagnostic severity](../../../specs/quality/usability/diagnosability/diagnostic-severity.requirement.spec.md),
[repeatable inspection](../../../specs/quality/reliability/repeatability/repeatable-inspection.requirement.spec.md),
and [source preservation](../../../specs/quality/reliability/source-preservation/inspection-preserves-source.requirement.spec.md)
also state explicit criteria. They own their criteria; the report-format
reference and inspection concept link to them. Responsiveness is a declared
attribute with no adopted threshold. A requirement does not establish fulfillment
merely by being present or passing a structural check.

Additional rules for requirements and system quality attributes:

| Rule ID | Meaning |
| --- | --- |
| `layout.requirement-location`, `layout.requirement-name` | Required layout and slug |
| `quality.attribute-source` | Give a system attribute a named source in its own quality folder |
| `identity.required` | Attribute definition needs an explicit ID |
| `quality.subattribute` | Sub-attribute list is empty or contains an undefined reference |
| `quality.cycle` | Sub-attribute graph contains a cycle |
| `quality.legacy-hierarchy` | Rename `Broader` or `Parents` to `Sub-attribute of` |
| `quality.definition` | Attribute definition is missing (advice) |
| `requirement.attributes` | Attribute list has an invalid shape or element |
| `requirement.undefined-attribute` | Label is absent from the associated system's attributes (advice) |
| `requirement.applies-to`, `requirement.statement`, `requirement.verification`, `requirement.rationale` | Missing section content (advice) |
| `requirement.measurement-basis` | A measure is recorded without a separately identified target or applicability conditions (advice) |

Existing identity, duplicate-field, section, and reference rules also apply.
