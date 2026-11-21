# What makes something a job to be done?

A person calls the plumber and considers the errand finished. The next day the
plumber requests a photo. Opening the old entry is one action the list offers.
The person's problem would also arise with a paper notebook: they have to pick
up the same commitment, recover what they already did, and work out what remains.

| Question | In this example |
| --- | --- |
| What changed in the person's circumstances? | A request for follow-up arrived after the call seemed finished. |
| What progress do they seek? | Carry the errand forward from where they left it, with a clear next action. |
| How could they recognize success? | They can say what they have done and what still needs attention without retracing the whole errand. |
| What might the software contribute? | The list lets them reopen the original entry. Its availability across restarts is a separate software commitment. |

A job helps an author explain why someone would choose a way to proceed.
It can outlast the particular tool, action or interface. The scenario above
illustrates a possible need. Finding out whether people actually lose track of
finished errands, and what they consult when returning, takes observation.
[Worked job](../../../examples/todo/specs/jobs/resume.job.spec.md) ·
[Return scenario](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md)

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
has been achieved. State the commitment in a requirement and link it to the outcome with Supports.

Keep importance, satisfaction and opportunity scores in dated research with a
defined population, method and source. Do not store them as enduring properties
of the outcome. Evidence can link external studies; dotspec does not implement a
survey repository, segmentation engine or scoring method. A separate assessment
record would be appropriate if comparison across studies becomes a product need.

To develop a draft, [write a job from a concrete situation](../guides/write-a-job.md).
