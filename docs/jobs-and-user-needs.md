# Jobs and user needs

Use this guide when documenting user progress, success criteria, or the abilities
needed to support them. The [job format](../README.md#job-spec-format) defines the
Markdown syntax; the [core model reference](reference/model/core.md) defines the
formal types and relationships.

## Choose a useful perspective

Dotspec supports several JTBD approaches through one model. It does not assume
that their theories or uses of the word “outcome” are identical.

| Perspective | Focus | Authoring support |
| --- | --- | --- |
| ODI | Functional accomplishment and measurable success criteria | `job init --preset odi`; optional steps and measurable outcome prompts |
| Progress and switching | Circumstances, motivation, choices and impediments | `job init --preset progress`; progress cases, forces, timeline and evidence |
| Job story | Concise situation, motivation and intended progress | `job init --preset job-story`; the existing Job fields |

Use `--preset` with `job check` to request the same guidance. The default
`general` check accepts every perspective. A preset is an authoring aid, not a
job subtype or a claim that the source complies with a research methodology.

Ulwick's [ODI needs framework](https://strategyn.com/customer-needs-through-a-jobs-to-be-done-lens/)
focuses on functional jobs, measurable outcomes and research into importance and
satisfaction. The [Christensen Institute](https://www.christenseninstitute.org/theory/jobs-to-be-done/)
emphasizes progress in circumstances, with functional, emotional and social
dimensions. [Moesta's approach](https://therewiredgroup.com/learn/complete-guide-jobs-to-be-done/)
examines struggling moments, choices, forces and trade-offs.

[Kalbach](https://experiencinginformation.com/2020/10/23/drills-and-milkshakes-say-yes-to-the-mess/)
connects aspirations, functional jobs and steps in a practical synthesis.
[Klement](https://medium.com/the-job-to-be-done/know-the-two-very-different-interpretations-of-jobs-to-be-done-5a18b748bd89)
argues for distinguishing the underlying theories more sharply. Dotspec preserves
the authored distinctions without declaring that debate settled.
[Intercom's job stories](https://www.intercom.com/blog/accidentally-invented-job-stories)
provide a concise presentation using fields already present on Job.

## Separate progress, criteria, commitments and evidence

- **Job:** what someone seeks to accomplish or change in circumstances, independently of the chosen solution.
- **JobStep:** a meaningful accomplishment within an optional functional map.
- **DesiredOutcome:** a desired result or criterion by which the person judges success.
- **ProgressCase:** an observed or hypothesized case illustrating progress, decisions and impediments.
- **Requirement:** a commitment about the specified subject.
- **Evidence:** authored observations, hypotheses, sources, research basis and limitations.

“Feel confident about arriving on time” is a valid qualitative outcome.
“Minimize time before discovering a disruption” is a candidate ODI criterion.
“Submit a notification within 30 seconds of receiving the event” is a software
commitment. Satisfying that commitment does not prove that either user outcome
has been achieved.

Keep importance, satisfaction and opportunity scores in dated research with a
defined population, method and source. Do not store them as enduring properties
of the outcome. Evidence can link external studies; dotspec does not implement a
survey repository, segmentation engine or scoring method. A separate assessment
record would be appropriate if comparison across studies becomes a product need.

## Declare a shared outcome once

An outcome has exactly one declaration owner, a Job or JobStep. Its optional
`applies_to` links identify jobs or steps for which it is a success criterion.
No applicability is inferred from ownership. Existing outcome lists remain
useful narrative declarations; add explicit applicability when a typed scope is
needed, including for an entire job.

Use a named outcome with an introductory `**ID:**`, a concise heading, and the
full statement. Place links under `Applies to`. A named outcome under a job uses
H3 and H4 headings; under a step it uses H5 and H6. Headings without an ID remain
supporting extensions. Ordinary links in supporting prose are citations.

See [Reach work predictably](../examples/user-needs/specs/jobs/reach-work.job.spec.md):
“Timely awareness” applies to both monitoring and adapting a journey, while
“Confidence in arrival” applies to the whole job. Declaration ownership and
applicability are separate relationships in the generated map.

`job outcome list --step ID` selects declaration ownership. `job map --outcomes`
shows the authored map and outcome declarations; it does not resolve shared
applicability. Use `model needs` for the relationship view across documents.

## Describe a progress case

Under `## Progress cases`, each H3 names a case. Optional H4 sections include:

- `Participants`: distinguish performer, buyer and decision maker when relevant.
- `Circumstances` and `Struggling moment`: the case-specific context and trigger.
- `Current approach` and `Desired progress`: the present situation and sought change.
- `Alternatives considered`: include continuing the current approach or doing nothing where relevant.
- `Trade-offs and selection criteria`: what matters in choosing or rejecting an approach.
- `Forces`: explicitly labeled Push, Pull, Anxiety and Habit descriptions.
- `Timeline`: the decision/change narrative, distinct from a functional job map.
- `Evidence`: source statements, `Research basis`, and `Limitations`.

These sections are source-backed model fields, not automatically inferred facts.
A case can be hypothetical, dates can be unknown, and fields can remain empty.
Keep emotional and social motivations visible without forcing metric wording.

## Connect needs to capabilities and specifications

Richard Allen's [User Needs Mapping](https://leanpub.com/userneedsmapping) connects
user needs to enabling capabilities and their dependencies, then considers
technical systems and organizational design. Dotspec adopts the needs and
capability relationships while leaving team design and Wardley evolution in
external artifacts. The distinction is informed by *User Needs Mapping*,
printed pp. 39–49; the mapping below is a dotspec design choice.

| Type | Role |
| --- | --- |
| Capability | An ability needed to deliver value, possibly combining people, procedures and technology |
| Feature | A recognizable grouping of specified product behavior |
| System / Component | The subject or structural part being specified |

Declare optional named capabilities under `## Capabilities` in a system file.
Each H3 capability can have an `**ID:**`, `Purpose`, `Supports`, `Depends on`,
`Provided by`, and `Enablers` sections. Its declaring system is its documentation
home; declaration does not make that system a provider.

- `supports` links to goals, jobs, steps or outcomes; it expresses intended contribution.
- `depends_on` links to other capabilities; dependencies can be many-to-many and cyclic.
- `provided_by` links to explicitly identified systems or components; it is optional and non-exhaustive.
- `enablers` describes human, procedural or other enabling context and external references.

A dependency is not a structural `part_of` relationship, and a provider link is
not verified implementation. Capabilities can span systems. Human and procedural
providers do not need to be falsely represented as software components.

Job, JobStep, DesiredOutcome and ProgressCase remain in the core
`user-requirements` subset. Capability participates in the core
`business-requirements` and `software-requirements` subsets; these memberships
organize model definitions, not the strategic importance of each capability.

## Generate a needs-to-specification view

```sh
dotspec model needs --path . --relative-to docs/reference > docs/reference/needs.md
```

The generated Markdown includes a Mermaid diagram, linked relationship table,
full statements and evidence. Short labels keep the graph readable; the
statements retain their meaning. Unresolved references remain in the authored
narrative and inspection findings rather than becoming invented edges.

Inspect the [generated example](reference/needs-map-example.md) and its
[source system](../examples/user-needs/specs/system.spec.md). The view combines
capability dependencies and providers with existing feature groups, requirement
applicability and contribution links. It covers only acquired records within
normal model-export discovery and direct-reference bounds.
