# Develop requirements

Use this guide when turning a need, rule, or quality concern into an obligation
that people can understand and assess. Start with a short statement, identify its
scope, and add assessment detail where it resolves ambiguity.

## Choose an authoring form

Use `dotspec requirement init <slug>` for a required result, rule, property, or
restriction. Add `--preset quality` when quality-focused prompts help; that
template also prompts for a Criterion. Both presets create the same requirement
type, with the same applicability, quality attributes, criteria, traceability,
verification and validation plans. A quality requirement is a requirement
classified by `quality_attributes`; it may also require behavior.

An obligation can also live inside a concept or use case. Keep one authoritative
definition and link to it from scenarios or other specifications. A separate file
is useful when the obligation needs independent navigation and discussion.

## State one obligation and its operating context

Identify the subject under **Applies to**, with relevant operating conditions,
workload, population and exclusions. Use explicit `**Targets:**` links when you
want semantic applicability to a system, feature, component, use case or other
domain element. Ordinary links remain contextual citations.

Write **Statement** as one coherent required result, rule, property or restriction.
Ask what observable situation would violate it. A general rule can apply across
many use-case paths without enumerating all of them. Prescribe ordering only when
the order is itself required.

For example, this requirement is assessable without a numeric target:

```markdown
## Applies to

Closed claims submitted through any supported editing interface.

## Statement

An attempt to modify a closed claim must leave its recorded contents unchanged.

## Rationale

Closure preserves the agreed claim record for later review.

## Verification plan

Attempt changes through each supported editing interface and compare the stored
claim before and after each attempt. Retain the attempted changes, interface and
implementation revision, and before-and-after records.
```

Conditions the system must accommodate are not assumptions that excuse its
behavior. Describe any reliance on an environmental premise separately; see
[requirements review](review-requirements.md#separate-premises-from-operating-conditions).

## Define or reference the criterion once

**Criterion** makes an obligation's satisfaction rule explicit when that helps.
It can be qualitative, quantitative, or a reference to a rule already defined in
Statement or another authoritative specification. There is no need to restate a
bound that already has an authoritative definition.

| Content | Authoritative role |
| --- | --- |
| Quality attributes | Classify the concern; they do not establish thresholds. |
| Applies to | Defines the subject and operating context covered by the obligation. |
| Statement and optional Criterion | State the obligation and define or reference its satisfaction rule. |
| Verification plan | Describes how to obtain and assess evidence against the obligation, including any measurement procedure and sampling. |
| Validation plan | Assesses the stated subject against stakeholder needs or intended use. |

For a measured property, identify the measure and the required bound. The measure
carries its units, aggregation and observation window where needed. Define the
population, numerator, denominator and success event for a rate. A label such as
"fast" or "95% successful" leaves too much open on its own.

Recognized Criterion subsections are Measure and Target, alongside prose and Open
questions. **Target** means the required value, bound or acceptable range,
including the comparison, such as "At least 90%". Operating context belongs only
under Applies to, which recognizes Conditions, Workload, Population and Exclusions
subsections. Record preferred targets in Rationale or Notes and label them
nonbinding. An unresolved bound belongs in Open questions with the evidence needed
to choose it; a draft can remain useful without pretending to be complete.

See the [unfinished recovery feedback requirement](../examples/requirements-review/specs/requirements/recovery-feedback.requirement.spec.md) for a draft that records a measure while the target and operating conditions are open.

## Describe context of use explicitly

Product quality concerns the product's properties. Quality in use concerns the
outcomes people achieve in a particular context. A fast API response can support
successful task completion, but does not establish it.

For outcomes of use, identify users, tasks, available assistance and environment
under Applies to. Use existing actor, job and use-case references where helpful;
the current model does not require a new context-of-use object or perspective
field. Population and Conditions subsections can make the context easier to find.

See the [unassisted recovery example](../examples/requirements-review/specs/requirements/unassisted-recovery.requirement.spec.md)
for a quality-in-use obligation with an explicit measure and bound, context of
use, and one verification plan. Its values illustrate the structure; they are
not recommended targets for other products.

Use **Supports** to relate an obligation to a job or desired outcome, **Derived
from** for its source, and **Refines** only when it adds precision while preserving
the referenced obligation. Meeting a latency bound alone does not refine or prove
a user-task completion requirement. See [requirements layers](requirements-layers.md).

## Plan the assessment without repeating the obligation

Both generated templates include **Verification plan**. Describe the realization
and revision, method, assessment conditions and evidence to retain. Put the
measurement procedure and sampling, where needed, under the plan's **Method**
subsection. Compare observations with the existing statement or criterion. A
**Success criteria** subsection may reference that rule rather than repeat its
bounds.

Keep operating context and assessment setup distinct. Applies to says when the
obligation holds; plan Conditions explains the test or study setup and how it
represents that context. A small test sample or an omitted failure mode must not
silently narrow the obligation. Record limitations of the resulting evidence.

Use **Validation plan** to assess whether the requirement captures the right need
or whether a proposed solution serves intended use. Identify the subject and
basis explicitly. Neither plan records execution evidence or a verdict;
see [verification and validation](verification-and-validation.md).

## Check the draft and review the set

Run `dotspec requirement check <slug>`, then `dotspec model review` for the
acquired set. Measurement-basis advice (`requirement.measurement-basis`) applies
when a Criterion records an explicit Measure. It prompts for a separately
identified target and operating context even if a Verification plan exists.
Introductory narrative under Criterion, including a reference, is sufficient to
record the target. Introductory narrative or a Conditions subsection under Applies
to is sufficient to record the operating context. Neither is interpreted or judged
for adequacy.

Plan Success criteria and assessment Conditions do not substitute for the
obligation's criterion and operating context. Advice can be resolved by identifying
or linking existing content; do not duplicate an authoritative target merely to
satisfy a prompt. Qualitative obligations do not need invented numeric measures.

Review individual requirements for necessity, clarity, feasibility and assessability.
Review the set for conflicting rules, inconsistent terminology or units, missing
conditions, unresolved questions and tradeoffs. Passing the automated checks does
not establish completeness, sufficiency or stakeholder agreement.

## Basis

This guidance adapts [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html)
§§5.2.4–5.2.6 (requirement construction and individual/set quality),
[ISO/IEC 25030:2019](https://www.iso.org/standard/72116.html) §§6.2–6.3 and 7.4.2
(quality perspectives, context, measures and criteria), and ISO/IEC 25020:2019
§§6.1 and 6.4 (measure definitions and measurement procedures). The current
[ISO 25010:2023](https://www.iso.org/standard/78176.html) model concerns product
quality; [ISO 25019:2023](https://www.iso.org/standard/78177.html) concerns quality
in use. Built-in model editions remain explicit.

Meyer's [Handbook of Requirements and Business Analysis](https://doi.org/10.1007/978-3-031-06739-6)
§4.12 motivates the violation question; §8.6 motivates general rules and contracts
without premature execution ordering. dotspec's field names and authoring forms
are design choices, not prescribed standard formats or conformance claims.
