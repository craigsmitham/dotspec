---
type: requirement-spec
quality_attributes: [usability]
companions:
  - target: ../../internal/integration/testdata/companions.txtar
    role: illustrates
    note: Declare, check, report and export companions
---
# List the companions that carry a specification's intent

## Applies to

Every specification type, `dotspec check`, `dotspec sources`, `dotspec spec build` and `dotspec model export` for projects declared with `dotspec.yaml`.

## Statement

An author must be able to list, in a specification's frontmatter, each artifact that carries part of its intent in another form, such as a formal model, an executable example, a test or an explanatory document, and state its role: formalizes, illustrates, verifies or explains. A reader of the generated report must be able to reach every listed artifact from the subject it serves and from the companion index of the outermost system that contains it, with each role stated. Inspection must locate a malformed entry, an unknown role, a missing target and a missing named definition at the entry that declares it.

## Criterion

A companion target resolves with the same boundary rules as an authored link. A fragment on a Quint, Alloy, TLA+, Gherkin, Go, Python or txtar target names a declaration found by its leading keyword, and the report links it to the declaring line. Each subject report groups its companions under Formalized by, Illustrated by, Verified by and Explained by. The outermost system's section map links to `specification/companions.md`, which lists each companion with its declaring subject. A requirement's verifies companions are its verification procedures. In a whole-project check, a `.qnt`, `.als`, `.tla` or `.feature` file that no specification lists produces advice naming the file. The export records each entry as a Companion owned by its specification, with its role, note and frontmatter target reference.

## Rationale

Part of a product's intent can live in formal models, executable examples and tests. Readers need one entry point that shows each of them with the authority its author declared. A link alone establishes addressability; the role tells a reviewer whether the artifact constrains the product, illustrates it, assesses it or explains it.

## Derived from

[Catalog what must survive regeneration](../../AGENTS.md#the-intent-layer).
