---
type: concept
id: job-map
---

# Job Map

## Purpose

Make an accomplishment's steps and success criteria inspectable, so people can
discuss what the work involves and how success is judged.

## State

A job description, steps in authored order, optional stage labels, overall job
outcomes, and outcomes associated with individual steps. Steps and outcomes can
have explicit IDs for reference and selection.

This state is authored in a job specification. The
[job format](../../README.md#job-spec-format) defines its representation.
An overall outcome belongs to the whole job and is not inherited by its steps.

## Actions

### Describe the job and its steps

Edit the job specification to describe the accomplishment, add or reorder steps,
and associate outcomes with the whole job or a particular step. These are file
edits; the CLI supplies a starting document and reads the authored state.

### View the map

Display the job's steps in authored order, optionally including outcomes beside
their owning job or step. A job with no steps still has a view that explicitly
reports their absence. The order expresses the author's map; it does not require
that execution sequence or establish completeness.

### Inspect a step's outcomes

Select a step by its explicit ID and display that step's outcomes. Overall job
outcomes are excluded from this selection. An existing step without outcomes
returns an empty result; an unknown step is an error.

The [job commands](../../README.md#job-commands) define the CLI operations and
output formats for these views.

## Operational principle

Describe a job such as
[Develop Requirements](../jobs/develop-requirements.job.spec.md), then add steps
such as identifying stakeholders and evaluating requirements. Attach success
criteria to the relevant steps and view the map with outcomes. The map makes
those relationships visible for discussion. Selecting one step lets a reviewer
focus on its criteria without confusing them with success for the entire job.
