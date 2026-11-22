# Write a business rule

Start with a policy provision or a decision someone in the organization makes. Find out who has authority to settle its scope.

“We usually get a manager for large refunds.” That interview remark is enough to start a draft. It leaves three decisions open. How large is *large*? Who may approve? What does *usually* allow? Use a recent refund to ask the person who owns the policy how they decided whether it could be paid. [What a business rule governs](../explanations/business-rules.md).

## Establish the authority and rationale

Under **Authority**, cite the decision that establishes the rule. Give enough detail to find the provision, such as its title, edition, and section. For an unwritten decision, identify who made it and when. If staff disagree about manager approval, record the disputed point and who can settle it. Observing a manager approve one refund tells you how that case was handled.

For a rule that classifies a case, ask which facts determine membership. If a club counts visits in the year they were completed, try a December booking completed in January. Record the completion criterion and the authority that chose it. [Definitional and behavioral rules](../explanations/business-rules.md#what-does-the-rule-let-someone-decide).

Check whether the condition settles every case the reader must classify. RuleSpeak’s high-risk customer example identifies a balance history that qualifies a customer. Ask the policy owner whether other histories can qualify. Then ask whether the classification obliges anyone to act. RuleSpeak illustrates an obligation with a notice that must be sent when an order ships. In To-do, a Wednesday follow-up becomes due on Wednesday. The owner can leave the errand waiting or choose Friday and see it marked upcoming again. Write the distinction in the owner’s terms. [Ross, *RuleSpeak Sentence Forms*, v2.2, Table 2, §§1.3–1.4, p. 9](https://www.rulespeak.com/en/RuleSpeak%20Sentence%20Forms%20v2-2-3.pdf#page=9) · [When a reply needs chasing](../../../examples/todo/specs/business-rules/follow-up.business-rule.spec.md#statement).

Under **Rationale**, explain why the distinction matters. In the fictional example below, Finance wants a second person to catch errors before money leaves the business. That explains the independent review. The threshold itself comes from the cited policy.

## State one decision

Write the **Statement** in terms the people applying the rule use. Suppose the policy owner confirms that independence from the requester is what matters and that a supervisor’s title is incidental. Say who may approve and when approval is due. Wording such as *may be paid only after* makes approval a necessary condition and leaves room for other reasons to withhold payment. [Ross, *RuleSpeak Sentence Forms*, v2.2, “Allowing something conditionally”](https://www.rulespeak.com/en/RuleSpeak%20Sentence%20Forms%20v2-2-3.pdf).

From your project directory, run `dotspec business-rule init refund-approval` to create `specs/business-rules/refund-approval.business-rule.spec.md`. Choose a slug for your own rule. If you already have the rule file, edit it in place. Replace the starter prompts with what the policy owner confirms. This worked draft uses a fictional policy and provision.

```markdown
---
type: business-rule-spec
---
# Independent approval for refunds

## Statement

A refund above USD 500 may be paid only after an employee other than the
requester has approved it.

## Applies to

Customer refunds denominated in US dollars, whether requested online or by
phone. The threshold applies to each individual refund.

## Authority

Finance, Refund Policy v3, §4 (Independent review).

## Rationale

Independent review gives a second employee a chance to catch an erroneous
refund before payment.
```

## Try the boundaries and exceptions

Use **Applies to** to settle which transactions, people, and circumstances the statement governs. Ask someone who handles refunds and the policy owner to decide the same cases. Where their answers differ, find the term, boundary, or missing decision that explains the difference. This adapts Matt Wynne’s Example Mapping technique of working from a rule to concrete cases and open questions. [Introducing Example Mapping](https://cucumber.io/blog/bdd/example-mapping-introduction/).

| Case | What the draft tells the reviewer |
| --- | --- |
| A customer requests exactly USD 500 by phone. | This threshold adds no approval condition. Other payment controls may still apply. |
| The requester approves their own USD 500.01 refund. | Independent approval is still outstanding. Using a different login would leave the identity of the approver unchanged. |
| Staff paid an urgent USD 600 refund before anyone else approved it. | On the stated facts the payment breached the rule. Ask whether an authorized waiver covered it and who could grant one. |

Put a confirmed qualification under **Exceptions**. If a waiver has its own authority or applies to several rules, give it a rule spec and link to it. An unresolved claim that “urgent refunds are different” belongs in the discussion with the policy owner. The Business Rules Manifesto expresses exceptions as rules in their own right. In dotspec, an author can also keep a short qualification with the rule it modifies. [Business Rules Manifesto, v2.0, art. 4.7](https://www.businessrulesgroup.org/brmanifesto/BRManifesto.pdf).

## Connect the requirement and check the draft

If the team commits the payment service to holding refunds pending approval, write that commitment and its criterion in a requirement file. From a requirement in `specs/requirements/`, add a relative Markdown link under **Derived from**, such as `[Independent approval](../business-rules/refund-approval.business-rule.spec.md)`. The link records the reason for the commitment. In the rule’s **Related specifications**, you can point a reviewer to the use case where an employee decides the request.

From the project directory containing the draft above, run `dotspec business-rule check refund-approval`. The CLI checks structure and references and suggests adding a missing statement, scope, rationale, or authority. Even the untouched starter can return zero findings because its prompts occupy those fields. The substantive check is whether the policy owner and the people applying the rule agree on the boundary cases and can point to the decision that settles each one.
