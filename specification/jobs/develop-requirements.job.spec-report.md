---
generated_by: dotspec
type: job-spec-report
---
# Develop Requirements

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Job

Define an agreed set of requirements for a proposed software change, from
understanding the intended change to making an agreed baseline available.

## Performer

The person responsible for defining the change, working with affected stakeholders.

## Situation

I’m defining a new product or changing an existing capability.

## Motivation

I want to turn stakeholder needs and product intent into clear requirements.

## Intended progress

The team can agree on what to build and how to determine whether it meets
those needs.

## Outcomes

- **OUT-001:** Minimize the likelihood that an important stakeholder need is
  discovered only after implementation begins.
- **OUT-002:** Minimize the time needed to reach agreement on requirements for
  a proposed change.

## Steps

### Determine the scope

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#determine-the-scope"} -->

**ID:** determine-scope
**Stage:** define

Establish the change's boundaries and the decisions the requirements must support.

#### Outcomes

- **OUT-003:** Minimize the likelihood that stakeholders interpret the boundaries
  of the proposed change differently.

### Identify stakeholders and sources

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#identify-stakeholders-and-sources"} -->

**ID:** identify-stakeholders
**Stage:** locate

Identify affected people and sources of relevant needs, constraints, and evidence.

#### Outcomes

- **OUT-004:** Minimize the time needed to identify affected stakeholders.
- **OUT-005:** Minimize the likelihood of overlooking a stakeholder with a
  relevant constraint.

### Organize the available information

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#organize-the-available-information"} -->

**ID:** organize-information
**Stage:** prepare

Organize needs, evidence, constraints, and assumptions so they can be considered together.

#### Outcomes

- **OUT-006:** Minimize the time needed to find evidence supporting a stated need.

### Establish readiness to specify

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#establish-readiness-to-specify"} -->

**ID:** establish-readiness
**Stage:** confirm

Determine whether enough information is available to formulate requirements
and identify any assumptions that must remain explicit.

#### Outcomes

- **OUT-007:** Minimize the likelihood of treating an unresolved assumption as
  an agreed stakeholder need.

### Formulate the requirements

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#formulate-the-requirements"} -->

**ID:** formulate-requirements
**Stage:** execute

Express the intended behavior and constraints in terms the team can understand
and assess.

#### Outcomes

- **OUT-008:** Minimize the likelihood that a requirement admits conflicting
  interpretations of the intended behavior.

### Evaluate the requirements

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#evaluate-the-requirements"} -->

**ID:** evaluate-requirements
**Stage:** monitor

Assess whether the set accurately represents stakeholder needs and is consistent,
complete enough for the change, and possible to verify.

**Supported by** · [Check specifications before committing](<../use-cases/check-specifications.use-case.spec-report.md>)

#### Outcomes

- **OUT-009:** Minimize the likelihood that conflicting requirements remain undetected.

[View source for “Evaluate the requirements”](<../../specs/jobs/develop-requirements.job.spec.md#evaluate-the-requirements>)

### Resolve gaps and disagreements

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#resolve-gaps-and-disagreements"} -->

**ID:** resolve-gaps
**Stage:** modify

Work with stakeholders to resolve missing information, ambiguity, and conflicts.
Revisit earlier steps as needed.

#### Outcomes

- **OUT-010:** Minimize the time needed to resolve a disagreement about intended behavior.

### Establish the agreed baseline

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/develop-requirements.job.spec.md#establish-the-agreed-baseline"} -->

**ID:** establish-baseline
**Stage:** conclude

Make the agreed requirements and remaining assumptions available to the people
responsible for implementation and verification.

#### Outcomes

- **OUT-011:** Minimize the time needed to identify the agreed requirements
  applicable to the proposed change.

## Evidence

This job map and its outcomes are working hypotheses based on the intended
purpose of dotspec. They have not been validated through interviews or outcome
research, and the map's completeness has not been established. The proposed
steps include collaboration and judgment that dotspec may only support.

[subject-source]: <../../specs/jobs/develop-requirements.job.spec.md#develop-requirements>
