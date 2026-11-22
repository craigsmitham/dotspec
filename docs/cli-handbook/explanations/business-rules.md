# Business rules

A business rule states a constraint an organization adopts or is obliged to follow. Suppose Finance adopts a policy requiring independent approval before a customer refund above USD 500 can be paid. The approver must be someone other than the requester. If the payment system went offline and staff issued refunds by hand, the same approval would still be due. Meyer uses the **source of the constraint** to distinguish a business rule from a physical limit or an engineering decision. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §1.7.2](https://doi.org/10.1007/978-3-031-06739-6).

## What does the rule let someone decide?

A **behavioral rule** governs conduct. Paying a USD 600 refund with only the requester’s approval would breach the finance rule. A **definitional rule** states what necessarily holds under the business’s definitions. One use is to decide which cases count. Suppose a club defines a qualifying visit as one the member completed during the membership year. When membership renews on January 1, a visit booked in December and completed in January counts toward the new year. If the club credits the visit to December, the member’s conduct is perfectly in order. The error is in the club’s count.

In dotspec, either decision can warrant its own rule spec when people need to cite or change it independently. A glossary entry can explain what the club means by *visit*. The renewal rule owns the decision about which visits count. SBVR distinguishes behavioral and definitional rules. ISO/IEC/IEEE 29148 also recognizes judgment criteria and formulas as operational rules. [SBVR 1.5, §§17.1.2 and 18.1.2](https://www.omg.org/spec/SBVR/1.5/PDF) · [ISO/IEC/IEEE 29148:2018, §9.3.11](https://www.iso.org/standard/72089.html).

## How does the rule reach the software?

The policy owner can change the refund threshold. The product team can change how an approver records a decision. Those changes affect different specifications.

| Subject | What it contributes |
| --- | --- |
| Business rule | Establishes which refunds require independent approval and when it is due. Its authority might be a finance policy, a recorded decision, or a binding external provision. |
| Requirement | Commits the payment service to hold a refund while approval required by the rule is outstanding. Its criterion says how to tell whether the service kept that commitment. |
| Use case | Describes how an approver reviews a request and what happens if they decline it. The scenario can cite the same rule wherever the decision arises. [Refund review paths](../guides/write-a-use-case.md#explore-what-can-change-the-outcome). |

A team could meet the policy through a manual review. Holding a payment in software is an additional product decision. Keeping one rule to cite makes it possible to review the consequences of a threshold change across every workflow that uses it. This follows the Business Rules Group’s principle of rule independence and the distinction ISO makes between business policies and the system requirements that address them. [Business Rules Manifesto, v2.0, arts. 2.2–2.3 and 4.5–4.6](https://www.businessrulesgroup.org/brmanifesto/BRManifesto.pdf) · [ISO/IEC/IEEE 29148:2018, §9.4.11](https://www.iso.org/standard/72089.html).

To develop a rule from an interview or policy document, see [Write a business rule](../guides/write-a-business-rule.md).
