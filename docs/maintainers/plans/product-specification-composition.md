# Composing the product specification

**Implemented in model 0.32.0-draft and Markdown mapping 26.** Run `dotspec model product` to inspect membership, adoption, and verification evidence for a selected system. Model 0.35.0-draft and mapping 29 replace the **Verification procedures** section with `verifies` companions in frontmatter and retire the criterion-level result subject. [Product view](../../reference/specification-model.md#read-one-product-specification) · [Verification results](../../reference/generated-specification.md#record-verification-results)

The rule is **one specification file, one specification root**. A root records a decision about a particular product that its authors can maintain independently. It owns the elements declared within it. The roots in scope for a system, their elements, and their relationships together compose the full product specification. Generated pages present views of that specification.

The model distinguishes a specification root in `*.spec.md` from a reusable quality model in `*.quality-model.md`. This design carries that distinction through ownership, product membership, external authority, and verification. The acceptance cases below exercise that distinction. [Specification and reference roots](../../reference/model/core.md#specification-and-reference-roots) · [What belongs in a specification?](../../cli-handbook/explanations/what-belongs-in-a-specification.md)

## 1. Give each decision one editable home

| Relationship | Meaning | Task-list example |
| --- | --- | --- |
| Declaration containment | The parent owns the element’s declaration and lifecycle within one source. Following containment reaches exactly one document root. | The waiting-errand use case owns its extension for a failed save. |
| System membership | An independently authored specification belongs to an effective system. The export records how that association was established. | The system quality attribute Durability and the restart requirement each have their own file and belong to Personal task list. |
| Structural composition | A declared system or component forms part of another. | Task service is part of Personal task list. |
| Semantic relationship | An authored statement relates independently identifiable subjects. Each endpoint retains its declaring source. | A step illustrates a requirement. A requirement applies to My checklist and is classified by Durability. |

Reserve declaration containment for same-source parts. Preserve the distinction between the enclosing component or system used for source lookup, the effective specification system, and an explicit structural parent. A cross-file link can express any supported semantic relationship without transferring authorship. The membership validator checks same-source containment and validates the system quality vocabulary against each attribute’s effective system. [Current membership validator](../../../internal/model/validation_domain.go) · [Task-list use case](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#the-list-cannot-record-the-decision)

Make `QualityAttribute` an abstract common definition. A **quality model attribute** is a named element owned by a reusable quality model. A **system quality attribute** is a specification root with an effective system and its own definition. Give each concrete type parent attributes of its own kind within the same model or system. Applying a model copies the selected definitions into system declarations and retains their source citations.

A **system quality requirement** remains a `Requirement` classified by one or more system quality attributes of its effective system. Narrow the schema endpoint accordingly. Its statement and criterion stay in the requirement. An attribute may exist before the team defines any requirements for it.

Use the system quality attribute’s **Definition** for the local quality statement. For Personal task list, Durability means preserving the owner’s acknowledged decisions between sessions. That tells an author why a saved follow-up date matters. The restart requirement then specifies which fields must survive, under what conditions, and how many discrepancies are acceptable. A team can author a system quality attribute directly or adapt a definition from a model. A source link records where the definition came from. Editing the model leaves the local statement intact. [System definition of Durability](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md#definition)

## 2. Derive one product specification from a selected system

Use a product-specification view over the source-backed snapshot. The input is a project scope and a selected system. The output identifies the selected system, member specification roots with their membership basis, their owned elements, relationships, and the inspection scope and completeness. Each member has a route back to its authored declaration.

1. Include the selected system and its declared descendant systems within the project. Include specification roots associated with those systems, including components. Use explicit system selection where supplied and the established directory default otherwise.
2. Include each root once, then include its owned elements in authored order. Preserve cross-references, including links to roots in another product or project. A shared requirement may appear in several reader journeys with one canonical identity.
3. Expose separately the reference documents and evidence reached from those members. A project reference catalog may also list unused quality models. Its contents do not change membership.
4. Expose unresolved associations, malformed declarations, and inspection gaps with their source locations. A partial scan can establish observed membership only. With several independent systems, a consumer selects a system. A reference-only project can publish a catalog and has no product specification to assemble.

Have the generated system report and model reading views consume this selection. Give the report planner a separate navigation parent for reference pages. Project quality models have a navigation parent for their browse route. The reference catalog and breadcrumbs should work without implying system membership. No additional product manifest or authored membership list is needed. [Report placement](../../../internal/specdoc/plan.go) · [Project and system boundaries](../../reference/workspaces.md)

On the system summary, keep a single **Quality Attributes** section for that system’s declarations. List each attribute by its linked name with a short description drawn from its local definition. The attribute page carries the full definition, its source when supplied, and links to its system quality requirements. Each requirement page keeps the statement, conditions, and criterion together. Keep reusable models reachable through reference navigation. If the system declares no attributes, omit the section and its contents link.

## 3. Make adoption of an external provision explicit

An ordinary Markdown citation supplies context. For a normative dependency, use the optional **Adopts** section on the local specification. Each named entry identifies an exact provision and states how it applies. Record the declaration as an owned element with its source passage, target locator, available edition or revision, and local applicability. The external artifact retains authority for the provision’s wording. The local specification owns the decision to adopt it.

For example, an interface could adopt the named cancellation clause of a partner contract for requests already accepted by the partner. A link to that same contract in **Contracts** supplies a citation. Two adopted clauses with different applicability get two entries. Preserve the author’s choice of pinned edition or evolving reference. A changed or unresolved target calls for review and retains its locator. Resolving a link establishes addressability alone.

## 4. Report procedures and observations independently

A requirement owns its criterion for acceptance. A **verification procedure** describes how someone can assess it. A **verification result** records what happened on one occasion. List each procedure as a `verifies` companion in the requirement’s frontmatter. Model 0.35.0-draft implements this and retires the earlier **Verification procedures** section, including its subsection under a criterion. Keep ordinary **Evidence** prose available for rationale, witnesses, and observations. Give a procedure a linkable identity even when it lives in a test suite or manual protocol. [Companions](../../../model/markdown-mapping.md#companions)

Accept a result from its evidence source with the requirement, procedure, outcome, assessed product revision, relevant conditions, and a reference to the observation. Record the specification revision assessed when known. Keep results outside specification aggregates so recording a run does not revise product intent. Preserve multiple observations, including conflicting runs. A missing revision or mismatched condition limits the claim a report can make.

Expose two independent facts in reports: the procedures identified for a requirement, and the observations available for the selected revision and conditions. Show “procedure linked” when only a procedure is known and “no matching result recorded” when applicable evidence is absent. A pass or failure belongs to the identified run. Existing prose and links retain their original meaning during migration. Explicit authoring or evidence import supplies the stronger roles. dotspec can index that evidence without becoming the test runner. A result names the requirement’s project-relative URI, without a criterion fragment. [Restart requirement and its procedure](../../../examples/todo/specs/requirements/restart.requirement.spec.md)

## Acceptance cases

| Change or fixture | Expected observation |
| --- | --- |
| Edit the failed-save extension | The owning use case changes. The extension remains addressable. The root count stays constant. |
| Edit a reusable definition after applying its model | The reference page changes. The system quality attribute and its classified requirements retain their wording and identity. |
| Classify a requirement using only a model attribute, or an attribute in another system | Validation reports the unresolved or invalid classification at its source. Matching names cannot satisfy it. |
| Link two use cases to one requirement and cite another project’s contract | The product view contains the requirement once. Both use cases link to it. The contract remains an external endpoint. |
| Replace a citation with an explicit adoption | The view identifies the local adopter, exact provision, and applicability. An unrelated citation to the same source retains its original role. |
| Link the restart test, then record a failing run and a passing run for a later revision | The procedure is discoverable before either run. Each observation belongs to its revision and conditions. The authored promise and product membership stay stable. |
| Add an unused quality model or generate a separate page for an owned element | The reference catalog or navigation changes. Product membership stays stable. |
| Declare Durability before writing its first requirement, then classify the restart requirement | The system summary shows Durability and its local meaning at both stages. The attribute page gains a link to the requirement. An unused model attribute never appears in the summary. |
| Make a member unreadable or narrow the scan | The view identifies the inspection limit and reports membership outside the observed scope as unknown. |

## Delivery order

1. **Ownership and quality types.** Change the metamodel, Markdown mapping, export, and offline validator together. Regenerate the schema, catalog, and model reference. Require fixtures for same-file containment, cross-file membership, and same-system quality classification.
2. **Product composition.** Deliver the shared selection in the model layer, exercise it against a multi-system fixture, then use it in reports. Compare the exported member identities with the generated reader routes. Rebuilding a report must preserve source ownership.
3. **External adoption.** Add the optional authoring form, verification of its structure, and a report showing the adopted provision with its scope. Test citation and adoption against the same target.
4. **Verification evidence.** Prove the procedure/result distinction with one executable test and one manual assessment. Use those fixtures to settle the smallest result interchange format before adding an importer. Test a stale revision, conflicting observations, and a requirement with no recorded result.

Each stage should ship its source mapping, generated references, and reader-journey tests together. For breaking schema changes, version the model and mapping and regenerate saved snapshots from source. The acceptance cases can be reviewed independently of the proposed heading names or serialization.
