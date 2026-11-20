# Verification and validation

Use these four questions when reviewing a requirement and the product it describes.
The subject is what you examine. The basis is what you judge it against.

## Assessment purposes

| Purpose | Subject | Assessment question |
| --- | --- | --- |
| Requirements verification | A requirement or set of requirements | Do the requirement statements and set meet the applicable criteria for clarity, consistency and verifiability? |
| Requirements validation | A proposed requirement or set of requirements | Do the proposed requirements faithfully express stakeholder needs, higher-level intent and constraints? |
| System verification | A system, component or process realization | Does the realization fulfill the specified requirements under the assessed conditions? |
| System validation | A proposed or realized solution | Can the intended users accomplish their work with the proposed or realized solution in its intended context? |

The requirements distinction follows ISO/IEC/IEEE 24765:2017, entries 3.3452–3.3453,
which reference requirements-engineering definitions in ISO/IEC/IEEE 29148.
[SEBoK's requirements assessment](https://sebokwiki.org/wiki/System_Requirements_Definition#Assessing_the_Requirements)
also distinguishes requirements quality from correctness with respect to intended
needs. Requirements verification is broader than parsing or checking required
fields: human review is needed for qualities that automated checks cannot establish.

The system distinction follows
[NASA's product verification and validation definitions](https://nodis3.gsfc.nasa.gov/displayDir.cfm?Internal_ID=N_PR_7123_001D_&page_name=AppendixA).
A system can conform to a poor specification. Conversely, a promising prototype
can demonstrate useful behavior without conforming to all specified requirements.

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

## Dotspec commands and model roles

`check` inspects structure and links. `model review` brings stated criteria,
assumptions, and open questions together for human review. `model validate`
checks a saved snapshot against its schema and graph rules.

Keep the statement, applicable conditions and satisfaction criteria together.
Link to a test, study or assessment where it helps readers evaluate a claim.
The [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md)
defines what survives a restart. Its linked acceptance test exercises that promise.
The [participant exercise](../../examples/todo/README.md#try-the-workflow-with-a-participant)
asks whether someone can resume interrupted work.

## Assumptions: evidence and justified reliance

Assumption validation has two questions:

1. Does evidence support the premise under the stated conditions?
2. Is relying on that premise appropriate for the intended use and consequences?

Measuring clock offset supplies empirical evidence. Deciding whether that offset
is acceptable for the recovery-window boundary assesses the justification for
reliance. Use the assumption’s Basis and Applicability to explain why and under which conditions it is credible.

## Plans, observations, assessments and acceptance

Keep these roles separate, even when one external artifact contains several:

| Role | Meaning |
| --- | --- |
| Plan | Proposed subject, basis, method, conditions, criteria and evidence to retain. |
| Observation | What was observed or measured for an identified revision, build or prototype under recorded conditions. A runner pass is an observation against that runner's assertions. |
| Assessment | What the evidence supports about conformance, fitness or an assumption, including scope and limits. |
| Acceptance decision | A stakeholder or governing decision informed by evidence, assessments and other considerations. |

Record observations and decisions in their native artifacts, identifying the revision
and conditions to which they apply. Link them from the specification where useful.
