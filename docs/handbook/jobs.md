# Jobs and user needs

Use this guide when documenting user progress, success criteria, or the abilities
needed to support them. The [job format](../reference/formats/job.md#job-spec-format) defines the
Markdown syntax; the [core model reference](../reference/model/core.md) defines the
formal types and relationships.

## Choose a useful perspective

dotspec supports several JTBD approaches through one model. It does not assume
that their theories or uses of the word “outcome” are identical.

| Perspective | Focus | Authoring support |
| --- | --- | --- |
| ODI | Functional accomplishment and measurable success criteria | `job init --preset odi`; optional steps and measurable outcome prompts |
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
argues for distinguishing the underlying theories more sharply. dotspec preserves
the authored distinctions without declaring that debate settled.
[Intercom's job stories](https://www.intercom.com/blog/accidentally-invented-job-stories)
provide a concise presentation using fields already present on Job.

## Separate progress, criteria, commitments and evidence

- **Job:** what someone seeks to accomplish or change in circumstances, independently of the chosen solution.
- **JobStep:** a meaningful accomplishment within an optional functional map.
- **DesiredOutcome:** a desired result or criterion by which the person judges success.
- **Requirement:** a commitment about the specified subject.
- **Evidence:** authored observations, hypotheses, sources, research basis and limitations.

“Feel confident about arriving on time” is a valid qualitative outcome.
“Minimize time before discovering a disruption” is a candidate ODI criterion.
“Submit a notification within 30 seconds of receiving the event” is a software
commitment. Satisfying that commitment does not prove that either user outcome
has been achieved. Jobs and job steps do not contain obligation clauses; state the
commitment in a requirement, or in a clause of a use case or concept, and relate it
with Supports.

Keep importance, satisfaction and opportunity scores in dated research with a
defined population, method and source. Do not store them as enduring properties
of the outcome. Evidence can link external studies; dotspec does not implement a
survey repository, segmentation engine or scoring method. A separate assessment
record would be appropriate if comparison across studies becomes a product need.

## Declare a shared outcome once

An outcome has exactly one declaration owner, a Job or JobStep. Its optional
`applies_to` links identify jobs or steps for which it is a success criterion.
Use a named heading for each outcome and write the full statement underneath.
Add an introductory `**ID:**` when the outcome needs a stable identifier. Put
applicability links under `Applies to`. A job-level outcome uses H3 for its
name and H4 for its fields. A step-level outcome uses H5 and H6. Ordinary links in supporting prose are citations.

In [Resume interrupted work](../../examples/todo/specs/jobs/resume.job.spec.md),
the job owns “Less reconstruction”. Its Applies to link identifies the particular
step where the person recalls outstanding work.

`job outcome list --step ID` selects declaration ownership. `job map --outcomes`
shows the authored map and outcome declarations; it does not resolve shared
applicability. Use `model needs` for the relationship view across documents.

## Cite the basis for a need

Use Evidence to summarize the basis for a stated need and link to its source.
Say whether it came from an observation or remains a hypothesis. Keep interview
timelines, alternatives, switching forces, and study limitations in the research
record, where readers can evaluate them in context.

## Connect needs to capabilities and specifications

Richard Allen's [User Needs Mapping](https://leanpub.com/userneedsmapping) connects
user needs to enabling capabilities and their dependencies, then considers
technical systems and organizational design. dotspec adopts the needs and
capability relationships while leaving team design and Wardley evolution in
external artifacts. The distinction is informed by *User Needs Mapping*,
printed pp. 39–49; the mapping below is a dotspec design choice.

| Type | Role |
| --- | --- |
| Capability | An ability needed to deliver value, possibly combining people, procedures and technology |
| Feature | A recognizable grouping of specified product behavior |
| System / Component | The subject or structural part being specified |

Create a capability with `dotspec capability init local-task-storage`. The
`*.capability.spec.md` file starts with `type: capability-spec`, an optional
frontmatter `id`, and an H1 naming the ability. Use level-two headings for
Description, Purpose, Supports, Depends on, Provided by and Enablers as needed.
The nearest enclosing system or component owns the file. Record providers
explicitly, including when a different component provides the ability.
See [Local task storage](../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).

- `supports` links to goals, jobs, steps, outcomes or use cases; it expresses intended contribution.
- `depends_on` links to other capabilities; dependencies can be many-to-many and cyclic.
- `provided_by` links to explicitly identified systems or components; it is optional and non-exhaustive.
- `enablers` describes human, procedural or other enabling context and external references.

A dependency is not a structural `part_of` relationship, and a provider link is
not verified implementation. Capabilities can span systems. Human and procedural
providers do not need to be falsely represented as software components.

Job, JobStep and DesiredOutcome remain in the core
`user-requirements` subset. Capability participates in the core
`business-requirements` and `software-requirements` subsets; these memberships
organize model definitions, not the strategic importance of each capability.

## Generate a needs-to-specification view

```sh
dotspec model needs
```

The generated Markdown includes a Mermaid diagram, linked relationship table,
full statements and evidence. Short labels keep the graph readable; the
statements retain their meaning. Unresolved references remain in the authored
narrative and inspection findings rather than becoming invented edges.

In the [To-do job](../../examples/todo/specification/jobs/resume-interrupted-work.job.spec-report.md#less-reconstruction),
Less reconstruction links to the restart requirement that supports it and the job
step where it applies. To inspect contributions together in the diagram and table,
run `dotspec model needs --project todo` from the dotspec repository.
