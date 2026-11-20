# Verification and validation

Each dotspec assessment plan names its **purpose** in its heading. Choose the
purpose by asking what you will assess (the **subject**) and what you will judge
it against (the **basis**). This vocabulary applies across
business, user and software requirements and is independent of the model's
core/supporting/generic classification.

## Assessment purposes

<!-- assessment-purposes -->

| Purpose | Subject | Assessment question |
| --- | --- | --- |
| Requirements verification | A requirement or set of requirements | Do the requirement statements and set meet the applicable criteria for clarity, consistency and verifiability? |
| Requirements validation | A proposed requirement or set of requirements | Do the proposed requirements faithfully express stakeholder needs, higher-level intent and constraints? |
| System verification | A system, component or process realization | Does the realization fulfill the specified obligations under the assessed conditions? |
| System validation | A proposed or realized solution | Can the intended users accomplish their work with the proposed or realized solution in its intended context? |
| Assumption validation | A premise on which the specification relies | Does the evidence support the premise under the stated conditions, and is reliance on it justified for the intended use? |

<!-- /assessment-purposes -->

The requirements distinction follows ISO/IEC/IEEE 24765:2017, entries 3.3452–3.3453,
which reference requirements-engineering definitions in ISO/IEC/IEEE 29148.
[SEBoK's requirements assessment](https://sebokwiki.org/wiki/System_Requirements_Definition#Assessing_the_Requirements)
also distinguishes requirements quality from correctness with respect to intended
needs. Requirements verification is broader than parsing or checking required
fields: human review is needed for qualities that automated checks cannot establish.

The system distinction follows
[NASA's product verification and validation definitions](https://nodis3.gsfc.nasa.gov/displayDir.cfm?Internal_ID=N_PR_7123_001D_&page_name=AppendixA).
A system can conform to a poor specification. Conversely, a promising prototype
can demonstrate useful behavior without conforming to all specified obligations.

For document recovery, these are four different questions:

- Is the recovery-window requirement unambiguous and consistent with related rules?
- Does that recovery window address the intended users' actual recovery needs?
- Does the identified build restore the required content and metadata within the window?
- Can representative users recover the intended work and resume their tasks?

## Methods do not determine purpose

Tests, analysis, inspection, demonstration and formal methods can contribute
verification evidence. Prototypes, stakeholder reviews, simulations and experiments
can support validation before production software exists. Both activities recur
throughout development and operation; validation is not confined to final acceptance.
The same activity can contribute evidence to different assessment purposes, but
the claims and bases must remain distinguishable.
[NASA product realization](https://www.nasa.gov/reference/5-0-product-realization/)
and [SEBoK system validation](https://sebokwiki.org/wiki/System_Validation) describe
multiple methods and lifecycle use.

Checking a specified response-time threshold is system verification even though it
is a quality requirement. Assessing whether the response time lets users complete
their work is system validation. Neither functional/quality classification nor
unit/acceptance test labels determine the purpose on their own.

**Verifiability** is a requirement's assessability: whether fulfillment can be
checked. It is not a completed verification activity or result. The
[IREB glossary](https://cpre.ireb.org/en/downloads-and-resources/glossary?first-letter=V&tag=cpre-glossary)
distinguishes verifiability from verification and defines requirements validation
in relation to stakeholder needs. A recorded plan can help review verifiability;
its mere presence does not establish it.

## Other established usages

Terminology is not uniform. NASA uses **requirements validation** broadly enough
to include clarity, consistency, verifiability and traceability as well as fit to
needs. dotspec distinguishes those quality and fitness questions explicitly.

Meyer's *Handbook of Requirements and Business Analysis*, terminology note on
printed page 20, deliberately abbreviates verification and validation together as
**verification**, including in its requirements-quality chapters. It also describes
the programming-research use of verification for correctness proofs. dotspec adopts
useful review practices from the handbook without adopting its umbrella shorthand.
[Handbook publication](https://doi.org/10.1007/978-3-031-06739-6).

**Snapshot validation** is conventional schema/data-validation terminology:
checking a representation against its constraints. It is not requirements or system
validation in the intended-use sense. The
[JSON Schema validation specification](https://json-schema.org/draft/2020-12/json-schema-validation)
uses validation for instance constraint checks.

## dotspec commands and model roles

| Surface | Scope of the claim |
| --- | --- |
| `model validate` | Snapshot conformance and integrity: schema, typed references, source fidelity and graph invariants. |
| `check` | Automated authoring and selected requirements-quality checks that support requirements verification. Passing does not complete requirements verification or establish fitness. |
| `model review` | Groups authored plans by purpose and surfaces review prompts about requirements quality, needs, assumptions and missing guidance. |
| `assessment_plans` on a traceable element | An ordered list of plans. Each purpose-qualified heading supplies the plan’s required `purpose`. Subject and Basis can each cite sources. |
| `criterion` | An optional explicit satisfaction rule (`ObligationCriterion`), defining or referencing required bounds for any obligation, with an optional Measure and Target. |

A well-formed snapshot can faithfully represent a malformed or incomplete source.
A source that passes authoring checks can still express the wrong requirements.
Use the qualified terms above instead of a blanket “specification validated.”

Use the purpose label from the table followed by “plan” as the heading, such as
**Requirements verification plan** or **System verification plan**. Each occurrence
starts a separate plan, including a second plan with the same purpose.

The [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md)
has both. Its requirements-verification plan cites the statement as its **Subject**
and gives criteria for reviewing the wording. Its system-verification plan names the
service as its **Subject** and cites the same statement as its **Basis**. Where the
plan is recorded tells readers where to find it. The purpose tells them which
question the author intends to answer.

Subject, Basis, Participants, Method, Conditions, Success criteria, Cases and
Evidence to retain are optional. Subject and Basis preserve both prose and links.
A case shares the purpose of its plan. `supports` records an intended contribution.

Define each required bound once in the statement or criterion, or reference its
authoritative source. Plan Success criteria explain interpretation against that
basis; they should refer to its bounds rather than redefine them. Applicability
describes the required operating context, while plan Conditions describe the
assessment setup and how it represents that context. See
[Develop requirements](requirements.md) for examples.

## Assumptions: evidence and justified reliance

Assumption validation has two questions:

1. Does evidence support the premise under the stated conditions?
2. Is relying on that premise appropriate for the intended use and consequences?

Measuring clock offset supplies empirical evidence. Deciding whether that offset
is acceptable for the recovery-window boundary assesses the justification for
reliance. The Assumption's Basis and Applicability explain the premise; its
Assumption validation plan explains how to assess it. An observation supports a bounded claim,
not a guarantee about every future environment. Stakeholder agreement alone does
not establish an empirical premise.

## Plans, observations, assessments and acceptance

Keep these roles separate, even when one external artifact contains several:

| Role | Meaning |
| --- | --- |
| Plan | Proposed subject, basis, method, conditions, criteria and evidence to retain. |
| Observation | What was observed or measured for an identified revision, build or prototype under recorded conditions. A runner pass is an observation against that runner's assertions. |
| Assessment | What the evidence supports about conformance, fitness or an assumption, including scope and limits. |
| Acceptance decision | A stakeholder or governing decision informed by evidence, assessments and other considerations. |

The current metamodel structures plans and preserves authored evidence; it does not
introduce an assessment-result or approval workflow. Revision-bound observations,
assessment rationale and acceptance records remain in referenced native artifacts.
A future evidence model must preserve their subject, basis, revision, context,
method, result and limits. No `verified` or `validated` boolean is inferred from
links, a plan, a passing test or a merge.

See [Review requirements](review.md) for authoring plans and reading advisory
findings. Compare the [restart verification plan](../../examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md#system-verification-plan)
with the [system validation exercise](../../examples/todo/specification/personal-task-list.system.spec-report.md#system-validation-plan).
One tests preservation of acknowledged tasks, and the other asks a participant
to resume interrupted work.
