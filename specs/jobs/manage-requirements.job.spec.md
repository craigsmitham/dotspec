---
type: job
id: manage-requirements
---

# Manage Requirements

## Job

Keep an agreed requirements set aligned with a proposed product change, from
identifying the change's implications to communicating an updated baseline.

## Performer

The person responsible for maintaining the requirements, working with
stakeholders, implementers, and reviewers affected by the change.

## Situation

The product, its implementation, or stakeholder needs are changing.

## Motivation

I want to understand which requirements are affected and keep them consistent
and current.

## Intended progress

The team can make changes while preserving its commitments and continue
to trust the requirements.

## Outcomes

- **OUT-001:** Minimize the likelihood that a product change violates a
  commitment that remains applicable.
- **OUT-002:** Minimize the time needed to bring the agreed requirements into
  alignment with an approved product change.

## Steps

### Determine the change to assess

**ID:** determine-change
**Stage:** define

Establish what is changing and the decisions needed about existing requirements.

#### Outcomes

- **OUT-003:** Minimize the likelihood of assessing a different change from the
  one stakeholders intend.

### Locate affected requirements

**ID:** locate-requirements
**Stage:** locate

Identify relevant requirements, commitments, rationale, and supporting evidence.

#### Outcomes

- **OUT-004:** Minimize the time needed to find requirements affected by a proposed change.

### Assemble the change context

**ID:** assemble-context
**Stage:** prepare

Bring together the current baseline, proposed change, dependencies, and open questions.

#### Outcomes

- **OUT-005:** Minimize the likelihood of assessing impact using an outdated baseline.

### Confirm applicable commitments

**ID:** confirm-commitments
**Stage:** confirm

Establish which commitments remain applicable and who must participate in decisions.

#### Outcomes

- **OUT-006:** Minimize the likelihood that an affected stakeholder is omitted
  from a decision about changing a commitment.

### Assess impact and update requirements

**ID:** update-requirements
**Stage:** execute

Determine how the change affects the requirements and propose the necessary revisions.

#### Outcomes

- **OUT-007:** Minimize the likelihood of overlooking a dependency affected by
  a requirements revision.

### Check alignment

**ID:** check-alignment
**Stage:** monitor

Evaluate whether the revised set is internally consistent and reflects the
agreed product change.

#### Outcomes

- **OUT-008:** Minimize the likelihood that the revised requirements conflict
  with commitments that remain applicable.

### Resolve discrepancies

**ID:** resolve-discrepancies
**Stage:** modify

Resolve inconsistencies and open questions with the relevant stakeholders,
revisiting the impact assessment when needed.

#### Outcomes

- **OUT-009:** Minimize the time needed to resolve a discrepancy between the
  agreed change and the revised requirements.

### Communicate the updated baseline

**ID:** communicate-baseline
**Stage:** conclude

Make the updated requirements, decisions, and remaining assumptions available
to the affected people.

#### Outcomes

- **OUT-010:** Minimize the likelihood that affected people continue to use
  superseded requirements after the change is agreed.

## Progress cases

### Adopting requirements alongside code

This case is a hypothesis to investigate with teams considering dotspec.

#### Current approach

Requirements are distributed across issue trackers, documents, and conversations.

#### Desired progress

Keep requirements connected to the implementation as the product changes.

#### Forces

- **Push:** Requirements repeatedly become disconnected from code changes.
- **Pull:** Relevant intent could be easier to find and review.
- **Anxiety:** Maintaining another artifact may create additional work.
- **Habit:** The team already knows its existing tools and conventions.

#### Timeline

To be established through interviews.

## Evidence

The map, outcomes, and progress case are working hypotheses based on dotspec's
intended purpose. No interview or outcome research has validated them yet.
The map's completeness has not been established. It includes decisions and
collaboration performed outside the CLI.
