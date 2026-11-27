---
generated_by: dotspec
type: concept-spec-report
source: "../../specs/concepts/job-map.concept.spec.md#job-map"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Concepts](<../dotspec.system.spec-report.md#concepts>)

# Job Map

## Purpose

Make an accomplishment's steps and success criteria inspectable, so people can
discuss what the work involves and how success is judged.

## State

A job description, steps in authored order, optional stage labels, outcomes declared by the job,
and outcomes declared by individual steps. Steps and outcomes can
have explicit IDs for reference.

This state is authored in a job specification. The
[job mapping](<../../model/markdown-mapping.md#job>) defines its representation.
Declaration ownership is separate from explicit outcome applicability; shared
outcomes can apply to multiple jobs or steps without duplication.

## Actions

[Describe the job and its steps](<#describe-the-job-and-its-steps>) · [View the map](<#view-the-map>)

### Describe the job and its steps

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/job-map.concept.spec.md#describe-the-job-and-its-steps"} -->

Edit the job specification to describe the accomplishment, add or reorder steps,
and associate outcomes with the whole job or a particular step. These are file
edits; the CLI supplies a starting document and reads the authored state.

### View the map

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../specs/concepts/job-map.concept.spec.md#view-the-map"} -->

Read the job's generated page, which presents its steps in authored order with
outcomes beside their owning job or step. A job with no steps still has a page
that explicitly reports their absence. The order expresses the author's map; it
does not require that execution sequence or establish completeness.

## Operational principles

Describe a job such as
[Develop Requirements](<../jobs/develop-requirements.job.spec-report.md>), then add steps
such as identifying stakeholders and evaluating requirements. Attach success
criteria to the relevant steps and read the generated job page. The map makes
those relationships visible for discussion: each step's criteria appear beside
it, apart from the criteria for success of the entire job.

---

[View source](<../../specs/concepts/job-map.concept.spec.md#job-map>)
