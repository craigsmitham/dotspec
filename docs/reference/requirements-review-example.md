# Requirements review

Generated from acquired specifications. Findings below are advisory review prompts, not completeness, acceptance or implementation claims. Missing links mean **no relationship recorded in this scope**. Drafts remain valid.

Model `0.12.0-draft`; mapping `dotspec-markdown/11`.

## Goal contribution paths

Paths follow explicit supports, derived-from and refinement links from obligations, concept actions or use cases. A recorded path identifies intent; reviewers must assess whether the behavior is sufficient.

- [Reduce lost work](<../../examples/requirements-review/specs/system.spec.md#reduce-lost-work>): Contribution path recorded; sufficiency requires review.
- [Reduce retention cost](<../../examples/requirements-review/specs/system.spec.md#reduce-retention-cost>): No contribution path recorded from specified behavior; review required.

## Scenario grounding

An Illustrates link identifies a rule or action behind a step. Participation links and ordinary citations do not establish this relationship; not every conversational step needs a separate obligation.

- [Main success scenario / ScenarioStep 1](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): no rule or action relationship recorded; review required.
- [Main success scenario / ScenarioStep 2](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): illustrates [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>), [Restore](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#restore>).
- [Main success scenario / ScenarioStep 3](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): no rule or action relationship recorded; review required.
## Environment and boundaries

Operating conditions and constraints describe what the solution must accommodate. Assumptions are premises relied upon and require a justified basis and validation. Explicit exclusions delimit scope. External partners are separate from structural parents.

### System: [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>)

### BusinessRule: [Recovery window](<../../examples/requirements-review/specs/business-rules/retention.business-rule.spec.md#recovery-window>)

**Statement**

Workspace members may recover accidentally deleted documents during the 30 days following deletion.

**Applicability**

**Targets:** [Document workspace](../../examples/requirements-review/specs/system.spec.md)

Applies to workspace members and documents deleted within that workspace.

### System: [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Environment**

Mobile clients sometimes lose connectivity. The service must accommodate retries after reconnection; connectivity loss is an operating condition, not a premise that excuses restoration obligations.

**Out of scope**

Recovering documents after their retained copies expire. External backup restoration has a separate boundary.

Interacts with: [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>).

**Required interfaces**

#### Time readings

Read UTC time for retention comparisons.

##### Counterpart

[Platform clock](../../examples/requirements-review/clock/specs/system.spec.md).

##### Contracts

[Time API](../../examples/requirements-review/contracts/time.openapi.yaml) is the authoritative wire contract.

**Provided interfaces**

#### Restoration requests

Accept restoration requests from authorized clients. The native contract remains to be selected.

Review required: Provided interfaces Restoration requests has no native contract reference recorded.

### Assumption: [Trusted time](<../../examples/requirements-review/specs/system.spec.md#trusted-time>)

**Statement**

The platform supplies a trustworthy UTC clock for retention calculations.

**Basis**

The [clock observation](../../examples/requirements-review/specs/system.spec.md#clock-observation) motivates this premise; it does not establish a permanent guarantee.

**Applicability**

Retention calculations in the managed production environment while the platform clock service meets its stated operating conditions.

**Validation plan**

#### Method

Compare platform time readings to an independent reference before deployment and during operation.

#### Success criteria

Offset remains within the retention boundary's agreed tolerance. The tolerance is still an open question.

## Obligation verification

Plans are guidance for verification or measurement, not evidence of execution.

### [Preserve document metadata on restoration](<../../examples/requirements-review/specs/requirements/restoration.requirement.spec.md#preserve-document-metadata-on-restoration>)

No verification or measurement guidance recorded; review required.

### [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>)

**Verification plan**

Exercise restoration with retained content and metadata; compare both to their pre-deletion values. Include requests immediately before the recovery deadline.

## Open questions

Consolidated authored questions. External tracking references own resolution workflow. Omitted specification fields are not automatically questions; missing metadata below is advisory.

### Question 1 — [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Question**

#### Clock tolerance

What time offset can the recovery-window boundary tolerate?

**Raised by**

Workspace requirements reviewer.

**Raised on**

2026-11-11.

**Resolve by**

Before acceptance of retention-boundary tests.

**Importance**

Serious: affects decisions near the recovery deadline.

**Evidence needed**

Clock offset measurements and stakeholder agreement on the boundary policy.

**Tracking reference**

[Boundary decision](../../examples/requirements-review/project-plan.md#retention-boundary-decision) owns resolution.

### Question 2 — [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Question**

#### External backup handoff

Who owns the handoff after retained copies expire?

Review required: no Raised by, Raised on, Resolve by, Importance, Evidence needed recorded.

## Terminology review

Use one preferred term consistently in specification prose; declared alternative names and synonym entries remain useful. Collisions below are scoped to the declaring System and require human review. This report does not scan prose for undeclared synonyms.

- Review required for **recovery copy**: [Retained document](<../../examples/requirements-review/specs/system.spec.md#retained-document>) (alternative name), [Recovery copy](<../../examples/requirements-review/specs/system.spec.md#recovery-copy>) (preferred term).

## Explanations and suggestions

These roles distinguish explanation, nonbinding implementation advice and examples from specified obligations. They do not record stakeholder endorsement.

### [Trash](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#trash>) — Implementation suggestions

An index could accelerate expiration queries without changing the specified effects.

### [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>) — Explanatory notes

A recovery window limits restoration from retained state; it does not describe external backup guarantees.

### [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>) — Implementation suggestions

A background worker could expire retained copies. This suggestion does not prescribe a scheduler or database.

### [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>) — Illustrative examples

A member restores a document one day after deletion and resumes editing.

## Inspection findings

- **advice** `relationship.verification_plan`: No verification or measurement plan is recorded for this obligation; review how it will be checked.
- **advice** `requirement.verification`: Describe how to assess the statement and retain the evidence.
