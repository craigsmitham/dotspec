# Project, Environment, Goals, System

Generated PEGS reading perspective. This view complements business, user and software requirements subsets; it does not replace them or impose a dependency order. It covers only acquired specifications. Referenced project plans and native contracts remain authoritative.

Model `0.12.0-draft`; mapping `dotspec-markdown/11`.

## Project

### [Platform clock](<../../examples/requirements-review/clock/specs/system.spec.md#platform-clock>)

No project context recorded. Link schedules, responsibilities, risks and decisions where they are maintained.

### [Document workspace](<../../examples/requirements-review/specs/system.spec.md#document-workspace>)

**Project context**

[Delivery and review plan](../../examples/requirements-review/project-plan.md) owns milestones and stakeholder review responsibilities.

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
