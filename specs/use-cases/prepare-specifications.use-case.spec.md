---
type: use-case
id: prepare-specifications
---

# Prepare specifications for review

## Goal

Give reviewers a readable account of proposed behavior and its related sources.

## System boundary

The contributor's specification preparation process, including editing tools and
the dotspec CLI. Code review and acceptance decisions follow this process.

## Primary actor

A contributor proposing a product change.

## Goal level

Summary (kite).

## Success criterion

The contributor has assembled the proposed behavioral account, related sources,
and known unresolved questions for review. Preparation does not establish approval.

## Main success scenario

1. The contributor [creates a specification draft](create-specification.use-case.spec.md)
   or edits an existing source to describe the proposed behavior.
2. The contributor links its governing rules and evidence, then
   [discovers specification sources](discover-specification-sources.use-case.spec.md).
3. The contributor [checks specifications](check-specifications.use-case.spec.md),
   resolves findings as appropriate, and identifies questions for reviewers.
4. The contributor makes the sources and remaining questions available for review.

## Extensions

### Behavior remains undecided

During steps 1–3, the team lacks an agreed response to an identified condition.
Record the condition and explicitly unresolved handling. Continue preparation for
review; the review must resolve it before a delivery relies on that promise.

### Inspection cannot finish

At step 3, the check follows its
[operational-failure branch](check-specifications.use-case.spec.md#unreadable-input-or-failed-discovery).
Repair the inspection scope or access and repeat step 3. If this cannot be done,
end with the gap visible rather than claiming completed inspection.

## Sequence

Representative. Drafting, discovery, and checking can repeat in any useful order.

## Authority

This is a summary of a contributor workflow, not a mandatory approval process or
an additional CLI contract. Linked user tasks explain existing command behavior.
The [authoring guide](../../docs/develop-use-cases.md#plan-growth-without-weakening-promises)
explains how delivery increments can span these narratives without changing their
behavioral authority.
