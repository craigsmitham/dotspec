# Project, Environment, Goals, System

Specification project: **requirements-review**. Project membership does not imply a System relationship.

Generated PEGS reading perspective. This view complements business, user and software requirements subsets; it does not replace them or impose a dependency order. It covers only acquired specifications. Referenced project plans and native contracts remain authoritative.

Model `0.16.0-draft`; mapping `dotspec-markdown/13`.

## Project

### [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>)

No project context recorded. Link schedules, responsibilities, risks and decisions where they are maintained.

### [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Project context**

[Delivery and review plan](../../examples/requirements-review/project-plan.md) owns milestones and stakeholder review responsibilities.

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

## Goals

### Job: [Recover lost work](<../../examples/requirements-review/specs/jobs/recover.job.spec.md#recover-lost-work>)

**Statement**

When valuable work is accidentally removed, I want to recover it so I can continue without recreating it.

Supports: [Reduce lost work](<../../examples/requirements-review/specs/system.spec.md#reduce-lost-work>).

### DesiredOutcome: [DesiredOutcome](<../../examples/requirements-review/specs/jobs/recover.job.spec.md>)

**Statement**

- Minimize the time required to resume work after recovering an accidentally deleted document.

### SystemGoal: [Reduce lost work](<../../examples/requirements-review/specs/system.spec.md#reduce-lost-work>)

**Statement**

Reduce the amount of valuable work lost through accidental deletion.

### SystemGoal: [Reduce retention cost](<../../examples/requirements-review/specs/system.spec.md#reduce-retention-cost>)

**Statement**

Reduce unnecessary storage cost without shortening the recovery commitment.

## System

Specifications below describe the solution's subjects, concepts, interactions and obligations. Declaration, structural parentage, applicability and behavioral relationships retain their distinct meanings.

### System: [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>)

**Purpose**

Provide time readings to platform clients.

**Scope**

The clock service boundary; consumers remain separate systems.

### Concept: [Trash](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#trash>)

**Purpose**

Make deletion reversible until a retained item expires.

### QualityRequirement: [Timely recovery feedback](<../../examples/requirements-review/specs/quality/recovery-feedback.quality.spec.md#timely-recovery-feedback>)

**Statement**

The workspace must provide recovery feedback promptly enough for members to continue their task.

### QualityRequirement: [Unassisted recovery](<../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#unassisted-recovery>)

**Statement**

Members in the stated context must be able to recover the intended document
without external assistance at the success rate defined in [Criterion](../../examples/requirements-review/specs/quality/unassisted-recovery.quality.spec.md#criterion).

Supports: [Recover lost work](<../../examples/requirements-review/specs/jobs/recover.job.spec.md#recover-lost-work>).

Applies to: [Restore a document](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-a-document>).

### GeneralRequirement: [Preserve document metadata on restoration](<../../examples/requirements-review/specs/requirements/restoration.requirement.spec.md#preserve-document-metadata-on-restoration>)

**Statement**

For the restoration required by Restore retained content, the restored metadata must preserve the original title, author identity, creation time and workspace identity.

Applies to: [Restore a document](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-a-document>), [DocumentTrash](<../../examples/requirements-review/specs/system.spec.md#documenttrash>).

### System: [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Purpose**

Help teams preserve and recover valuable work.

**Scope**

Document deletion and restoration within a workspace.

Interacts with: [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>).

### ObligationClause: [Preserve request identity](<../../examples/requirements-review/specs/system.spec.md#preserve-request-identity>)

**Statement**

A retried restoration request must retain its original request identifier.

Applies to: [Restoration requests](<../../examples/requirements-review/specs/system.spec.md#restoration-requests>).

### ConceptUse: [DocumentTrash](<../../examples/requirements-review/specs/system.spec.md#documenttrash>)

Definition: [Trash](<../../examples/requirements-review/specs/concepts/trash.concept.spec.md#trash>).

Supports: [Recover lost work](<../../examples/requirements-review/specs/jobs/recover.job.spec.md#recover-lost-work>).

### Interface: [Time readings](<../../examples/requirements-review/specs/system.spec.md#time-readings>)

**Description**

**ID:** time-readings

Read UTC time for retention comparisons.

**Contracts**

[Time API](../../examples/requirements-review/contracts/time.openapi.yaml) is the authoritative wire contract.

**Counterpart**

[Platform clock](../../examples/requirements-review/clock/specs/system.spec.md).

### Interface: [Restoration requests](<../../examples/requirements-review/specs/system.spec.md#restoration-requests>)

**Description**

**ID:** restoration-requests

Accept restoration requests from authorized clients. The native contract remains to be selected.

### UseCase: [Restore a document](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-a-document>)

**Goal**

Recover an accidentally deleted document and resume work.

Supports: [Recover lost work](<../../examples/requirements-review/specs/jobs/recover.job.spec.md#recover-lost-work>).

Involves: [DocumentTrash](<../../examples/requirements-review/specs/system.spec.md#documenttrash>).

### ObligationClause: [Restore retained content](<../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content>)

**Statement**

When an authorized workspace member requests restoration before the 30-day recovery period ends, the workspace must restore the retained document content and metadata.

Supports: [Recover lost work](<../../examples/requirements-review/specs/jobs/recover.job.spec.md#recover-lost-work>).

Applies to: [DocumentTrash](<../../examples/requirements-review/specs/system.spec.md#documenttrash>).

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
