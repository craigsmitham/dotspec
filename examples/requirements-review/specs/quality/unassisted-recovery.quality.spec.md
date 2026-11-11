---
type: quality
id: unassisted-recovery
---
# Unassisted recovery

## Applies to

**Targets:** [Restore a document](../use-cases/restore.use-case.spec.md).

### Population

Workspace members who have completed normal onboarding but have never used recovery.

### Conditions

Each member has accidentally deleted a document they recognize and has permission
to restore it within its retention window. They use a supported desktop browser
with a working network connection. Ordinary in-product help is available; no
facilitator or colleague provides recovery instructions.

### Exclusions

Expired copies, unsupported browsers, and recovery during a network interruption.
These exclusions limit this requirement; other restoration obligations still apply.

## Statement

Members in the stated context must be able to recover the intended document
without external assistance at the success rate defined in [Criterion](#criterion).

## Criterion

### Measure

The proportion of eligible task attempts in which the member restores the intended
document and opens the recovered contents without external assistance. Count one
attempt per member; an abandoned attempt or facilitator intervention is a failure.
Divide successful attempts by all eligible attempts, including failures.

### Comparison

At least.

### Target

90%.

## Supports

[Recover valuable work](../jobs/recover.job.spec.md).

## Rationale

Recovery should let members resume work without waiting for support. This fictional
example uses an illustrative bound; a real product would justify it through user
research and stakeholder agreement. A response-time measurement alone would not
establish successful recovery.

## Verification plan

### Subject

The candidate recovery workflow, with implementation revision recorded before the study.

### Basis

[Statement](#statement), [Criterion](#criterion), and [Applies to](#applies-to).

### Participants

Members matching the specified population, observed by a facilitator who does not
provide recovery instructions.

### Method

Recruit 30 eligible members, each attempting one recovery in a seeded workspace.
Balance recruitment across supported desktop browsers and include members with
different levels of general workspace experience. Record all attempts and compute
the proportion using the criterion's definition. Report the sample size, uncertainty,
and recruitment limitations alongside the comparison; the sample is not proof of
success for every user or environment.

### Success criteria

Compare the observed proportion with the required bound in [Criterion](#criterion).
Retain the comparison and its limitations for the assessment decision.

### Conditions

Use a working network and supported browsers, with normal in-product help enabled.
Give each member a recognizable document to delete before the recovery task.
Explain how the seeded task represents the context in [Applies to](#applies-to)
and record deviations without silently changing the requirement's scope.

### Evidence to retain

Retain the implementation revision, recruitment criteria, browser and task setup,
attempt outcomes, interventions, calculation and assessment limitations. Record
the assessment conclusion separately from this plan.
