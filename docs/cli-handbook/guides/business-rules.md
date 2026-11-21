# Business rules

Use a business-rule specification to state a business constraint that can be
maintained independently of the process or software enforcing it. A rule can
apply to manual activity and several applications. Its source or authority is
an authored claim; dotspec does not verify that authority or establish compliance.

For example, “Purchases above $10,000 require approval from a department manager”
is a business rule. “Submit a purchase for approval” describes a software use
case. A requirement preventing unapproved purchases from reaching a supplier
specifies software enforcement. Keep these meanings separate and connect them
with explicit references.

## Create and inspect a rule

```sh
dotspec business-rule init purchase-approval
dotspec business-rule show purchase-approval
dotspec business-rule list
dotspec business-rule check purchase-approval --json
```

The default destination is `specs/business-rules/purchase-approval.business-rule.spec.md`
under the selected project directory. Commands accept `--path`; read commands
leave the file unchanged, and `init` refuses to overwrite an existing file.
System association follows the same repository context rules as other specs.
It does not limit the business scope of the rule.

## Authoring format

```markdown
---
type: business-rule-spec
id: purchase-approval
---

# Manager approval for large purchases

## Statement

Purchases above $10,000 require approval from a department manager.

## Applies to

Purchases made by any department.

## Rationale

Make expenditure accountable.

## Authority

The board's purchasing policy establishes this rule.

## Exceptions

Emergency purchases follow a separately approved exception rule.

## Related specifications

[Submit a purchase](../use-cases/submit-purchase.use-case.spec.md) describes
one software interaction governed by this rule.
```

The `id` is optional and local to the document. A title-only draft is valid;
missing Statement, Applies to, Rationale, or Authority content produces advice.
Exceptions can be absent, empty, or described explicitly. Omission does not mean
there are no exceptions.
Repeated recognized sections, incorrect levels for those sections, invalid
metadata, and unresolved local fragments remain structural errors.

Use Related specifications for links and Evidence for supporting sources. Put qualifications
next to the statement or exception they explain.

## Refer to a rule from software specifications

Cite a rule under `Related specifications` in a use case, requirement, or another spec:

```markdown
## Related specifications

[Manager approval](../business-rules/purchase-approval.business-rule.spec.md#purchase-approval)
governs this interaction. The linked requirement describes automated enforcement.
```

The model preserves the explanation and reference occurrences in
`related_specifications: ContextualReferences`. Resolving a citation identifies its
target; it does not create a typed enforcement or compliance assertion.

Use a typed relationship when the connection itself matters. `Derived from` on
another element names the rule as its basis. `Refines` on a requirement adds precision to the rule. The `**Targets:**` field in the
rule's own `Applies to` section identifies what the rule governs, and an
`Illustrates:` bullet under a scenario step names a rule that the step illustrates.
See [Declare relationships](relationships.md). None of these establishes
enforcement or compliance.

## Model representation

`BusinessRule` is a `Specification` in the core `business-requirements`
subset. Statement, applicability, rationale, authority, and exceptions have
named source-backed model fields. The rule's constraint is independent of any
software requirement that implements it. A rule's own `Refines` section targets
only other business rules. Specify required software enforcement in a requirement.

The [core domain reference](../reference/model/core.md) shows its place among the
specification types. The [complete reference](../reference/model/index.md#businessrule)
and [mapping profile](../../../model/markdown-mapping.md#business-rule) define its model
shape and mapping. Use [model export and validation](../reference/specification-model.md) to
inspect source retention, references, and coverage.
