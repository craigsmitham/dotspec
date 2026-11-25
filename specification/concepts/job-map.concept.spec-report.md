---
generated_by: dotspec
type: concept-spec-report
source: "../../specs/concepts/job-map.concept.spec.md#job-map"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Concept definitions](<../dotspec.system.spec-report.md#job-map>)

# Job Map

## Purpose

Make an accomplishment's steps and success criteria inspectable, so people can
discuss what the work involves and how success is judged.

## State

A job description, steps in authored order, optional stage labels, outcomes declared by the job,
and outcomes declared by individual steps. Steps and outcomes can
have explicit IDs for reference and selection.

This state is authored in a job specification. The
[job mapping](<../../model/markdown-mapping.md#job>) defines its representation.
Declaration ownership is separate from explicit outcome applicability; shared
outcomes can apply to multiple jobs or steps without duplication.

## Actions

[Describe the job and its steps](<#describe-the-job-and-its-steps>) · [View the map](<#view-the-map>) · [Inspect a step's outcomes](<#inspect-a-steps-outcomes>)

### Describe the job and its steps

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/job-map.concept.spec.md#describe-the-job-and-its-steps"} -->

[View section source](<../../specs/concepts/job-map.concept.spec.md#describe-the-job-and-its-steps>)

Edit the job specification to describe the accomplishment, add or reorder steps,
and associate outcomes with the whole job or a particular step. These are file
edits; the CLI supplies a starting document and reads the authored state.

### View the map

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/job-map.concept.spec.md#view-the-map"} -->

[View section source](<../../specs/concepts/job-map.concept.spec.md#view-the-map>)

Display the job's steps in authored order, optionally including outcomes beside
their owning job or step. A job with no steps still has a view that explicitly
reports their absence. The order expresses the author's map; it does not require
that execution sequence or establish completeness.

### Inspect a step's outcomes

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/job-map.concept.spec.md#inspect-a-steps-outcomes"} -->

[View section source](<../../specs/concepts/job-map.concept.spec.md#inspect-a-steps-outcomes>)

Select a step by its explicit ID and display outcomes it declares. Outcomes
declared elsewhere are excluded from this selection; the needs map shows their
explicit applicability. An existing step without outcomes
returns an empty result; an unknown step is an error.

The command reference describes how to [view the map](<../../docs/reference/cli.md#dotspec-job-map>) and [inspect a step’s outcomes](<../../docs/reference/cli.md#dotspec-job-outcome-list>).

## Operational principles

Describe a job such as
[Develop Requirements](<../jobs/develop-requirements.job.spec-report.md>), then add steps
such as identifying stakeholders and evaluating requirements. Attach success
criteria to the relevant steps and view the map with outcomes. The map makes
those relationships visible for discussion. Selecting one step lets a reviewer
focus on its criteria without confusing them with success for the entire job.

---

[View source](<../../specs/concepts/job-map.concept.spec.md#job-map>)
