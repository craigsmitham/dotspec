---
type: requirement-spec
quality_attributes: [usability]
verification:
  - name: Catalog a suite and resolve its declarations
    source: ../../internal/integration/testdata/catalog.txtar
---
# Find the governing source of a product decision

## Statement

An author must be able to specify behavior in Markdown or name the source that governs it in its native form. A declaration on the relevant subject gives the source a name, states what it defines and addresses a file or an exact definition. One declaration may cover a suite of requirements. Other subjects can cite the declaration or a native definition. An author can associate verification procedures with either form of specification. Contextual links supply examples and explanation.

## Applies to

Specification frontmatter, inspection, source inventory, snapshot export and generated specifications.

## Criterion

A reader can explain the promised behavior from the declaring subject’s Markdown and follow its link to the native source that governs the decision. When a particular definition governs the decision, the link identifies that definition.

From a feature and the requirements index a reader can reach its named executable suite. From the suite's inventory entry a reader can locate the feature that declared it. Adding another governing case in the suite requires no Markdown edit. Renaming a cited scenario produces a finding at the referring occurrence. A report distinguishes a procedure from an observation and can show a run of the suite itself for the recorded revision and conditions.

A passing bounded model check supports a claim about the model at the recorded revision, within the checked bounds and conditions. A claim that the running implementation conforms to that model requires evidence connecting the implementation to the model.

A whole-project check identifies candidate artifacts using recognized formats and project-selected patterns. The declaration establishes authority. Inspection establishes addressability within the supported formats. The snapshot preserves the declaration, its owner and source, verification procedures and contextual references.

## Rationale

Maintaining the same decision in prose and in executable cases creates two places to edit and review. Naming the governing source lets a reader find it from the product specification while the team changes each decision in one place.
