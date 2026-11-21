---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/prepare-specifications.use-case.spec.md#prepare-specifications-for-review"
---
# Prepare specifications for review

[dotspec](<../dotspec.system.spec-report.md>)

## Goal

Give reviewers a readable account of proposed behavior and its related sources.

## Primary actor

A contributor proposing a product change.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/prepare-specifications.use-case.spec.md#main-success-scenario"} -->

1. The contributor [creates a specification draft](<create-specification.use-case.spec-report.md>)
   or edits an existing source to describe the proposed behavior.
2. The contributor links its governing rules and evidence, then
   [discovers specification sources](<discover-specification-sources.use-case.spec-report.md>).
3. The contributor [checks specifications](<check-specifications.use-case.spec-report.md>),
   resolves findings as appropriate, and identifies questions for reviewers.
4. The contributor makes the sources and remaining questions available for review.

## Extensions

### Behavior remains undecided

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/prepare-specifications.use-case.spec.md#behavior-remains-undecided"} -->

#### Handling

During steps 1–3, the team lacks an agreed response to an identified condition.
Record the condition and explicitly unresolved handling. Continue preparation for
review; the review must resolve it before a delivery relies on that promise.

### Inspection cannot finish

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/prepare-specifications.use-case.spec.md#inspection-cannot-finish"} -->

#### Handling

At step 3, the check follows its
[operational-failure branch](<check-specifications.use-case.spec-report.md#unreadable-input-or-failed-discovery>).
Repair the inspection scope or access and repeat step 3. If this cannot be done,
end with the gap visible rather than claiming completed inspection.

## Success guarantees

The contributor has assembled the proposed behavioral account, related sources,
and known unresolved questions for review. Preparation does not establish approval.

## System boundary

The contributor's specification preparation process, including editing tools and
the dotspec CLI. Code review and acceptance decisions follow this process.

## Goal level

Summary (kite).

## Sequence

Representative. Drafting, discovery, and checking can repeat in any useful order.

## Authority

This is a summary of a contributor workflow, not a mandatory approval process or
an additional CLI contract. Linked user tasks explain existing command behavior.
Delivery increments can span these narratives. The linked use cases retain authority over their interactions.
