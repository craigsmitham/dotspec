---
type: system
---
# Document workspace

## Purpose

Help teams preserve and recover valuable work.

## Scope

Document deletion and restoration within a workspace.

## Goals

### Reduce lost work

Reduce the amount of valuable work lost through accidental deletion.

### Reduce retention cost

Reduce unnecessary storage cost without shortening the recovery commitment.

## Environment

Mobile clients sometimes lose connectivity. The service must accommodate retries after reconnection; connectivity loss is an operating condition, not a premise that excuses restoration obligations.

## Out of scope

Recovering documents after their retained copies expire. External backup restoration has a separate boundary.

## Project context

[Delivery and review plan](../project-plan.md) owns milestones and stakeholder review responsibilities.

## Interacts with

[Platform clock](../clock/specs/system.spec.md) supplies time readings across the workspace boundary.

## Required interfaces

### Time readings

Read UTC time for retention comparisons.

#### Counterpart

[Platform clock](../clock/specs/system.spec.md).

#### Contracts

[Time API](../contracts/time.openapi.yaml) is the authoritative wire contract.

## Provided interfaces

### Restoration requests

Accept restoration requests from authorized clients. The native contract remains to be selected.

## Assumptions

### Trusted time

The platform supplies a trustworthy UTC clock for retention calculations.

#### Basis

The [clock observation](#clock-observation) motivates this premise; it does not establish a permanent guarantee.

#### Applicability

Retention calculations in the managed production environment while the platform clock service meets its stated operating conditions.

#### Validation plan

##### Method

Compare platform time readings to an independent reference before deployment and during operation.

##### Success criteria

Offset remains within the retention boundary's agreed tolerance. The tolerance is still an open question.

## Evidence

### Clock observation

This example uses an invented observation to illustrate evidence linking. Replace it with actual measurements before relying on the premise.

#### Type

Illustrative hypothesis.

## Open questions

### Clock tolerance

What time offset can the recovery-window boundary tolerate?

#### Raised by

Workspace requirements reviewer.

#### Raised on

2026-11-11.

#### Resolve by

Before acceptance of retention-boundary tests.

#### Importance

Serious: affects decisions near the recovery deadline.

#### Evidence needed

Clock offset measurements and stakeholder agreement on the boundary policy.

#### Tracking reference

[Boundary decision](../project-plan.md#retention-boundary-decision) owns resolution.

### External backup handoff

Who owns the handoff after retained copies expire?

## Glossary

### Retained document

A deleted document kept inside its recovery window.

#### Alternative names

- Recovery copy

### Recovery copy

An external backup copy. This deliberate alias collision needs review.

## Explanatory notes

A recovery window limits restoration from retained state; it does not describe external backup guarantees.

## Implementation suggestions

A background worker could expire retained copies. This suggestion does not prescribe a scheduler or database.

## Illustrative examples

A member restores a document one day after deletion and resumes editing.

## Concept uses

### DocumentTrash

#### Definition

[Trash](concepts/trash.concept.spec.md)

#### Supports

[Recover lost work](jobs/recover.job.spec.md).

#### Relies on

[Trusted time](#trusted-time) bounds retention intervals.
