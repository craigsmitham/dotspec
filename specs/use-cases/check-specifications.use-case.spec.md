---
type: use-case
id: check-specifications
---

# Check specifications before committing

## Goal

Identify structural problems and missing-content advice in the specifications
affected by a proposed change before committing it.

## System boundary

The dotspec CLI, reading specifications from the selected project subtree and
any enclosing system specifications needed as context.

## Primary actor

A contributor preparing a change for review.

## Goal level

User task (sea level).

## Trigger

The contributor requests a specification check before committing a change.

## Preconditions

The selected project directory exists. Dotspec can read its files, and Git is
available when discovery takes place inside a Git repository.

## Success criterion

The contributor receives a complete report identifying the inspected scope,
structural errors, advisory findings, and any additional system context. The
report distinguishes a completed scan from an operational failure. A report
with structural errors still completes this information-gathering goal.

The [check report contract](../../README.md#check-specifications) defines the
report fields and exit statuses. A successful check does not establish that
requirements are complete, validated, or satisfied.

## Main success scenario

1. The contributor selects the project subtree to inspect and requests a check.
2. Dotspec discovers supported specifications and reads the system context
   needed to interpret them.
3. Dotspec checks structure and local references and gathers writing advice.
4. Dotspec presents the report, including source locations and scan completion,
   so the contributor can decide what to revise before committing.

## Extensions

- A specification is malformed: report its structural error, skip dependent
  checks for that file, and continue inspecting other discovered files.
- A file cannot be read or discovery fails: report the operational failure and
  mark the scan incomplete. The contributor can repair access and retry.
- No supported specifications are found: explicitly report an empty, completed
  scan.
- Content is unfinished: retain readable drafts and report applicable writing
  advice without treating omissions as structural errors.

## Sequence

Representative. The narrative explains the interaction; the
[check report contract](../../README.md#check-specifications) governs error
handling and reporting behavior.

## Related jobs

This use case contributes to
[Evaluate the requirements](../jobs/develop-requirements.job.spec.md#evaluate-the-requirements)
in the Develop Requirements job, including its conflict-detection outcome
(`OUT-009`).
Structural findings can inform that evaluation; human review still determines
whether the requirements accurately represent stakeholder needs.
