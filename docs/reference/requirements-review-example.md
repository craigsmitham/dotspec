# Requirements review

Specification project: **requirements-review**. Project membership does not imply a System relationship.

Generated from acquired specifications. Findings below are advisory review prompts, not completeness, acceptance or implementation claims. Missing links mean **no relationship recorded in this scope**. Drafts may conform to the authoring format while requirements verification and validation remain unfinished.

Model `0.16.0-draft`; mapping `dotspec-markdown/13`.

## Needs and intent

Requirements validation asks whether the specification expresses the right needs. System validation asks whether a proposed or realized solution serves intended use. The relationships below identify review questions, not either conclusion.

### Goal contribution paths

Paths follow explicit supports, derived-from and refinement links from obligations, concept actions or use cases. A recorded path identifies intent; reviewers must assess whether the behavior is sufficient.

- [Reduce lost work](<../../examples/requirements-review/specs/system.spec.md#reduce-lost-work>): Contribution path recorded; sufficiency requires review.
- [Reduce retention cost](<../../examples/requirements-review/specs/system.spec.md#reduce-retention-cost>): No contribution path recorded from specified behavior; review required.

### Scenario grounding

An Illustrates link identifies a rule or action behind a step. Participation links and ordinary citations do not establish this relationship; not every conversational step needs a separate obligation.

- [Main success scenario / ScenarioStep 1](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): no rule or action relationship recorded; review required.
- [Main success scenario / ScenarioStep 2](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): illustrates [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>), [Restore](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#restore>).
- [Main success scenario / ScenarioStep 3](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md>): no rule or action relationship recorded; review required.
## Environment and boundaries

Operating conditions and constraints describe what the solution must accommodate. Assumptions are assessed separately for evidence and justified reliance. Explicit exclusions delimit scope. External partners are separate from structural parents.

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

**ID:** time-readings

Read UTC time for retention comparisons.

##### Counterpart

[Platform clock](../../examples/requirements-review/clock/specs/system.spec.md).

##### Contracts

[Time API](../../examples/requirements-review/contracts/time.openapi.yaml) is the authoritative wire contract.

**Provided interfaces**

#### Restoration requests

**ID:** restoration-requests

Accept restoration requests from authorized clients. The native contract remains to be selected.

##### Obligations

###### Preserve request identity

**ID:** request-identity

A retried restoration request must retain its original request identifier.

**Applies to**

**Targets:** [Restoration requests](../../examples/requirements-review/specs/system.spec.md#restoration-requests).

**Verification plan**

Repeat a request after an interrupted response and inspect the retained identifier.

##### Notes

The interaction contract covers retries; the native protocol remains to be selected.

Review required: Provided interfaces Restoration requests has no native contract reference recorded.

## Assumptions: evidence and justified reliance

Assess both empirical support in the stated conditions and whether relying on the premise is appropriate for intended use. An observation supports a bounded claim, not a guarantee that the premise always holds.

### [Trusted time](<../../examples/requirements-review/specs/system.spec.md#trusted-time>)

**Statement**

The platform supplies a trustworthy UTC clock for retention calculations.

**Basis**

The [clock observation](../../examples/requirements-review/specs/system.spec.md#clock-observation) motivates this premise; it does not establish a permanent guarantee.

**Applicability**

Retention calculations in the managed production environment while the platform clock service meets its stated operating conditions.

**Assumption validation plan**

#### Subject

The premise that production clock offset stays within the agreed tolerance, and the justification for relying on it at the recovery boundary.

#### Basis

[Clock observation](../../examples/requirements-review/specs/system.spec.md#clock-observation) and the [recovery window](../../examples/requirements-review/specs/business-rules/retention.business-rule.spec.md).

#### Participants

Platform operators and the retention reviewer.

#### Conditions

The managed production environment, including clock resynchronization events.

#### Method

Compare platform time readings to an independent reference before deployment and during operation.

#### Success criteria

Offset remains within the retention boundary's agreed tolerance. The tolerance is still an open question.

## Validation planning

Identify the subject and the needs, intended use or other basis for the assessment. Requirements validation assesses proposed obligations; system validation assesses proposed or realized solutions. A method or supports link alone establishes neither purpose nor a result.

### System: [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### Concept: [Trash](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#trash>)

**Validation plan**

#### Subject

The proposed reversible-deletion concept; assess the operational principle against intended recovery use.

#### Basis

[Restore a document](../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md).

#### Participants

Workspace members and the recovery-policy stakeholder.

#### Method

Review examples and a prototype with the participants, including counterexamples and unsupported situations.

#### Success criteria

The participants can explain how the proposal addresses recovery needs, and unresolved mismatches are recorded for revision.

#### Conditions

Representative accidental-deletion situations and the stated recovery-window boundary.

#### Evidence to retain

Keep example responses and observations, the reviewed revision, the assessment rationale and any disagreements. Record acceptance separately.

### QualityRequirement: [Timely recovery feedback](<../../examples/requirements-review/specs/quality/recovery-feedback.quality.spec.md#timely-recovery-feedback>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### QualityRequirement: [Unassisted recovery](<../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#unassisted-recovery>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### GeneralRequirement: [Preserve document metadata on restoration](<../../examples/requirements-review/specs/requirements/restoration.requirement.spec.md#preserve-document-metadata-on-restoration>)

**Validation plan**

#### Subject

The proposed obligation to preserve all four metadata fields; assess whether it correctly expresses recovery needs.

#### Basis

[Recover lost work](../../examples/requirements-review/specs/jobs/recover.job.spec.md).

#### Participants

Workspace members and the recovery-policy stakeholder.

#### Method

Review examples and a prototype with the participants, including counterexamples and unsupported situations.

#### Success criteria

The participants can explain how the proposal addresses recovery needs, and unresolved mismatches are recorded for revision.

#### Conditions

Representative accidental-deletion situations and the stated recovery-window boundary.

#### Evidence to retain

Keep example responses and observations, the reviewed revision, the assessment rationale and any disagreements. Record acceptance separately.

### System: [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Validation plan**

#### Subject

The proposed workspace recovery workflow in a prototype, then the identified release build; this assesses intended use rather than only conformance to the written recovery rule.

#### Basis

[Recover lost work](../../examples/requirements-review/specs/jobs/recover.job.spec.md) and [Reduce lost work](../../examples/requirements-review/specs/system.spec.md#reduce-lost-work).

#### Participants

Representative workspace members, including mobile users, and a stakeholder who owns the recovery policy.

#### Method

Observe representative users attempting recovery with realistic document collections; discuss where the proposed workflow fails their intended progress.

#### Success criteria

Participants recover the intended document and can resume their interrupted task without specialist assistance. Record exceptions and limits rather than generalizing to every user.

#### Conditions

Representative desktop and mobile sessions, including reconnection after intermittent connectivity; record prototype or build revision and participant context.

#### Cases

##### Resume interrupted work

Ask a participant to recover a mistakenly deleted document and resume the task.

###### Expected result

The participant identifies the intended document, recovers it, and resumes the task.

###### Evidence to retain

Consented observation notes, task outcome, revision, conditions and observed obstacles.

#### Evidence to retain

Retain the protocol and observations separately from the assessment of intended-use fitness and any stakeholder acceptance decision. This example records a plan only.

### ObligationClause: [Preserve request identity](<../../examples/requirements-review/specs/system.spec.md#preserve-request-identity>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### UseCase: [Restore a document](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-a-document>)

**Validation plan**

#### Subject

The proposed interaction and alternatives; assess whether they support members in resuming interrupted work.

#### Basis

[Recover lost work](../../examples/requirements-review/specs/jobs/recover.job.spec.md).

#### Participants

Workspace members and the recovery-policy stakeholder.

#### Method

Review examples and a prototype with the participants, including counterexamples and unsupported situations.

#### Success criteria

The participants can explain how the proposal addresses recovery needs, and unresolved mismatches are recorded for revision.

#### Conditions

Representative accidental-deletion situations and the stated recovery-window boundary.

#### Evidence to retain

Keep example responses and observations, the reviewed revision, the assessment rationale and any disagreements. Record acceptance separately.

### ObligationClause: [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

## System verification planning

These plans assess realizations against obligations. Requirements verification instead assesses the quality of the requirement statements and set. A plan is distinct from an observation, an evidence-based assessment, and an acceptance decision.

### [Timely recovery feedback](<../../examples/requirements-review/specs/quality/recovery-feedback.quality.spec.md#timely-recovery-feedback>)

**Measurement plan**

#### Procedure

Measure elapsed time from submitting a recovery request to visible feedback.

#### Sampling

Collect observations for representative recovery requests.

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

Review required: A measure or measurement plan is recorded, but no separately identified criterion, conditions or applicability is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result.

### [Unassisted recovery](<../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#unassisted-recovery>)

**Verification plan**

#### Subject

The candidate recovery workflow, with implementation revision recorded before the study.

#### Basis

[Statement](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#statement), [Criterion](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#criterion), and [Applies to](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#applies-to).

#### Participants

Members matching the specified population, observed by a facilitator who does not
provide recovery instructions.

#### Method

Recruit 30 eligible members, each attempting one recovery in a seeded workspace.
Balance recruitment across supported desktop browsers and include members with
different levels of general workspace experience. Record all attempts and compute
the proportion using the criterion's definition. Report the sample size, uncertainty,
and recruitment limitations alongside the comparison; the sample is not proof of
success for every user or environment.

#### Success criteria

Compare the observed proportion with the required bound in [Criterion](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#criterion).
Retain the comparison and its limitations for the assessment decision.

#### Conditions

Use a working network and supported browsers, with normal in-product help enabled.
Give each member a recognizable document to delete before the recovery task.
Explain how the seeded task represents the context in [Applies to](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#applies-to)
and record deviations without silently changing the requirement's scope.

#### Evidence to retain

Retain the implementation revision, recruitment criteria, browser and task setup,
attempt outcomes, interventions, calculation and assessment limitations. Record
the assessment conclusion separately from this plan.

**Criterion**

#### Measure

The proportion of eligible task attempts in which the member restores the intended
document and opens the recovered contents without external assistance. Count one
attempt per member; an abandoned attempt or facilitator intervention is a failure.
Divide successful attempts by all eligible attempts, including failures.

#### Comparison

At least.

#### Target

90%.

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

### [Preserve document metadata on restoration](<../../examples/requirements-review/specs/requirements/restoration.requirement.spec.md#preserve-document-metadata-on-restoration>)

No system verification or measurement guidance recorded; review required.

### [Preserve request identity](<../../examples/requirements-review/specs/system.spec.md#preserve-request-identity>)

**Verification plan**

Repeat a request after an interrupted response and inspect the retained identifier.

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

### [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>)

**Verification plan**

Exercise restoration with retained content and metadata; compare both to their pre-deletion values. Include requests immediately before the recovery deadline.

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

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

## Requirements quality

Requirements verification examines individual statements and the set for clarity, consistency and assessability against declared quality criteria. Terminology checks and inspection findings below supply selected review prompts; passing them does not establish complete requirements verification.

### Terminology review

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

## Automated inspection findings

- **advice** `relationship.motivation`: No explicit motivation or derivation is recorded for this obligation; review its justification.
- **advice** `relationship.measurement_basis`: A measure or measurement plan is recorded, but no separately identified criterion, conditions or applicability is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result.
- **advice** `relationship.verification_plan`: No system verification or measurement plan is recorded for this obligation; review how a realization will be checked.
- **advice** `relationship.motivation`: No explicit motivation or derivation is recorded for this obligation; review its justification.
- **advice** `quality.applies-to`: Identify the subject and applicable conditions, workload, or population.
- **advice** `quality.measurement-basis`: A measure or measurement plan is recorded, but no separately identified criterion, conditions or applicability is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result.
- **advice** `requirement.verification`: Record assessment guidance for checking a realization against the obligation; a plan is not a verification result.
