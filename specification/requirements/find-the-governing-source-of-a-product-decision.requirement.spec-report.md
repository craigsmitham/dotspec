---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/specification-catalog.requirement.spec.md#find-the-governing-source-of-a-product-decision"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../dotspec.system.spec-report.md#requirements>)

# Find the governing source of a product decision

## Statement

An author must be able to specify behavior in Markdown or name the source that governs it in its native form. A declaration on the relevant subject gives the source a name, states what it defines and addresses a file or an exact definition. One declaration may cover a suite of requirements. Other subjects can cite the declaration or a native definition. An author can associate verification procedures with either form of specification. Contextual links supply examples and explanation.

## Applies to

Specification frontmatter, inspection, source inventory, snapshot export and generated specifications.

## Criterion

From a feature and the requirements index a reader can reach its named executable suite. From the suite's inventory entry a reader can locate the feature that declared it. Adding another governing case in the suite requires no Markdown edit. Renaming a cited scenario produces a finding at the referring occurrence. A report distinguishes a procedure from an observation and can show a run of the suite itself for the recorded revision and conditions.

A whole-project check identifies candidate artifacts using recognized formats and project-selected patterns. The declaration establishes authority. Inspection establishes addressability within the supported formats. The snapshot preserves the declaration, its owner and source, verification procedures and contextual references.

## Rationale

Maintaining the same decision in prose and in executable cases creates two places to edit and review. Naming the governing source lets a reader find it from the product specification while the team changes each decision in one place.

## Verification

- [Catalog a suite and resolve its declarations](<../../internal/integration/testdata/catalog.txtar>)

No matching result recorded.

## Related specifications

**Quality attribute**

- [Usability](<../quality/usability/usability.quality-attribute.spec-report.md>)

---

[View source](<../../specs/requirements/specification-catalog.requirement.spec.md#find-the-governing-source-of-a-product-decision>)
