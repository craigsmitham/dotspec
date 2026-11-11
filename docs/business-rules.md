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

The file lives at `specs/business-rules/purchase-approval.business-rule.spec.md`
under the selected project directory. Commands accept `--path`; read commands
leave the file unchanged, and `init` refuses to overwrite an existing file.
System association follows the same repository context rules as other specs.
It does not limit the business scope of the rule.

## Authoring format

```markdown
---
type: business-rule
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
there are no exceptions. `Authority and evidence` is an alias for Authority.
Repeated recognized sections, incorrect levels for those sections, invalid
metadata, and unresolved local fragments remain structural errors.

Related specifications, Notes, Evidence, and the other shared model sections
retain their source content. Unknown sections survive as explicit extensions.

## Refer to a rule from software specifications

Use a `Business rules` section in a use case, requirement, or another spec:

```markdown
## Business rules

[Manager approval](../business-rules/purchase-approval.business-rule.spec.md#purchase-approval)
governs this interaction. The linked requirement describes automated enforcement.
```

The model preserves the explanation and reference occurrences in
`business_rules: ContextualReferences`. Resolving a citation identifies its
target; it does not create a typed enforcement or compliance assertion.
Stronger relationships require a separately defined authoring convention.

## Model representation

`BusinessRule` is a `SpecificationElement` in the core `business-requirements`
subset. Statement, applicability, rationale, authority, and exceptions have
named source-backed model fields. The rule's constraint is independent of any
software requirement that implements it.

The [core domain reference](reference/model/core.md) shows its place among the
specification kinds. The [complete reference](reference/model/index.md#businessrule)
and [mapping profile](../model/markdown-mapping.md#business-rule) define its model
shape and mapping. Use [model export and validation](specification-model.md) to
inspect source retention, references, and coverage.
