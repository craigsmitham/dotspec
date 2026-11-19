---
generated_by: dotspec
type: use-case-spec-report
documents: "../../specs/use-cases/check-specifications.use-case.spec.md#check-specifications-before-committing"
---
# Check specifications before committing

[dotspec](<../system.spec-report.md>)

## Goal

Identify structural problems and missing-content advice in the specifications
affected by a proposed change before committing it.

## Primary actor

A contributor preparing a change for review.

## Trigger

The contributor requests a specification check before committing a change.

## Preconditions

The contributor has selected an existing project directory. Reading files and
invoking Git during discovery can fail; those failures are handled below, not
assumed away by this precondition.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#main-success-scenario"} -->

1. The contributor selects the project subtree to inspect and requests a check.
2. Dotspec discovers supported specifications and reads the system context
   needed to interpret them.
3. Dotspec checks structure and references and gathers writing advice.
4. Dotspec presents the report, including source locations and scan completion,
   so the contributor can decide what to revise before committing.

## Extensions

### Malformed specification

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#malformed-specification"} -->

During steps 2–3, a specification is malformed. Report its structural error,
skip dependent checks for that file, and continue according to the
[continued-inspection requirement](<../requirements/continue-after-invalid-spec.requirement.spec-report.md>).
Resume reporting at step 4. A complete report with errors still achieves this goal.

### Unreadable input or failed discovery

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#unreadable-input-or-failed-discovery"} -->

During steps 2–3, a file cannot be read or discovery fails, including a required
Git invocation. Report the operational failure and mark the scan incomplete.
End without achieving the complete-report goal; the contributor can repair access
and start another check.

### Empty scope

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#empty-scope"} -->

At step 2, no supported specifications are found. Explicitly report an empty,
completed scan at step 4. This is an alternative successful outcome.

### Unfinished content

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#unfinished-content"} -->

At step 3, content is unfinished. Retain readable drafts, report applicable writing
advice, and continue to step 4. The
[advice-permits-success requirement](<../requirements/advice-permits-success.requirement.spec-report.md>)
governs exit status. This also completes the information-gathering goal.

## Success guarantees

The contributor receives a complete report identifying the inspected scope,
structural errors, advisory findings, and any additional context. The
report distinguishes a completed scan from an operational failure. A report
with structural errors still completes this information-gathering goal.

The [check report contract](<../../docs/reference/inspection.md#check-specifications>) defines the
report fields and exit statuses. A successful check does not establish that
requirements are complete, validated, or satisfied.

## Minimal guarantees

Inspection leaves authored sources unchanged, as defined by the
[source preservation requirement](<../requirements/inspection-preserves-source.requirement.spec-report.md>).
The [check report contract](<../../docs/reference/inspection.md#check-specifications>) requires the
report to distinguish incomplete scans from completed ones, even when inspection
cannot finish.

## System boundary

The dotspec CLI, reading specifications from the selected project subtree,
enclosing system specifications, and direct reference targets needed as context.

## Goal level

User task (sea level).

## Supporting actors

Git supplies repository discovery information when the selected directory is
inside a Git repository.

## Stakeholder interests

Contributors need actionable findings. Reviewers need an honest account of what
was inspected. Both need their authored sources preserved.

## Sequence

Representative. The narrative explains the interaction; the
[check report contract](<../../docs/reference/inspection.md#check-specifications>) governs error
handling and reporting behavior.

## Supports

This use case contributes to
[Evaluate the requirements](<../jobs/develop-requirements.job.spec-report.md#evaluate-the-requirements>)
in the Develop Requirements job, including its conflict-detection outcome
(`OUT-009`).
Structural findings can inform that evaluation; human review still determines
whether the requirements accurately represent stakeholder needs.

## Authority and evidence

This narrative explains obligations defined by the linked report contract and
requirements. Their definitions remain authoritative. The
[use-case command examples](<../../internal/integration/testdata/use-cases.txtar>)
define the illustrated CLI outcomes; they do not prove every narrative branch is
covered. A reference to [unreadable input](<#unreadable-input-or-failed-discovery>)
identifies a branch, not a successful execution of it.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
describes the report and inspection behavior used in this interaction.

Its [quality obligations](<../concepts/specification-inspection.concept.spec-report.md#state>)
cover diagnostic severity, report repeatability, and source preservation.

[Edit source](<../../specs/use-cases/check-specifications.use-case.spec.md#check-specifications-before-committing>)
