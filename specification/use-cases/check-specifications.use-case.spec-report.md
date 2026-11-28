---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/check-specifications.use-case.spec.md#check-specifications-before-committing"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Use cases](<../dotspec.system.spec-report.md#use-cases>)

# Check specifications before committing

## Goal

Identify structural problems and missing-content advice in the specifications
affected by a proposed change before committing it.

## Primary actor

A [contributor](<../dotspec.system.spec-report.md#contributor>) preparing a change for review.

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

[Malformed specification](<#malformed-specification>) · [Unreadable input or failed discovery](<#unreadable-input-or-failed-discovery>) · [Empty scope](<#empty-scope>) · [Unfinished content](<#unfinished-content>)

### Malformed specification

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#malformed-specification"} -->

#### Handling

During steps 2–3, a specification is malformed. Report its structural error,
skip dependent checks for that file, and continue according to the
[continued-inspection requirement](<../requirements/continue-after-invalid-spec.requirement.spec-report.md>).
Resume reporting at step 4. A complete report with errors still achieves this goal.

### Unreadable input or failed discovery

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#unreadable-input-or-failed-discovery"} -->

#### Handling

During steps 2–3, a file cannot be read or discovery fails, including a required
Git invocation. Report the operational failure and mark the scan incomplete.
End without achieving the complete-report goal; the contributor can repair access
and start another check.

### Empty scope

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#empty-scope"} -->

#### Handling

At step 2, no supported specifications are found. Explicitly report an empty,
completed scan at step 4. This is an alternative successful outcome.

### Unfinished content

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/check-specifications.use-case.spec.md#unfinished-content"} -->

#### Handling

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

The dotspec CLI, reading specifications from the [selected project subtree](<../../docs/reference/layout.md#creating-and-finding-sources>),
enclosing system specifications, and direct reference targets needed as context.

## Goal level

User task (sea level).

## Supporting actors

Git supplies repository discovery information when the selected directory is
inside a Git repository.

## Stakeholder interests

[Contributors](<../dotspec.system.spec-report.md#contributor>) need actionable findings. [Reviewers](<../dotspec.system.spec-report.md#reviewer>) need an honest account of what
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

## Authority

This use case specifies the contributor’s inspection interaction. The linked report contract governs output and the linked requirements govern diagnostic severity, repeatability and preservation of the input.

## Related specifications

[Specification Inspection](<../concepts/specification-inspection.concept.spec-report.md>)
describes the report and inspection behavior used in this interaction.

Its [quality requirements](<../concepts/specification-inspection.concept.spec-report.md#state>)
cover diagnostic severity, report repeatability, and source preservation.

## Specification sources

<a id="specification.repository-discovery"></a>

### [Repository discovery](<../../internal/integration/testdata/discovery.txtar>)

Defines requirements.

<a id="specification.diagnostics-and-exit-statuses"></a>

### [Diagnostics and exit statuses](<../../internal/integration/testdata/diagnostics.txtar>)

Defines requirements.

<a id="specification.focused-and-project-wide-checks-of-draft-and-malformed-use-cases"></a>

### [Focused and project-wide checks of draft and malformed use cases](<../../internal/integration/testdata/use-cases.txtar>)

Defines requirements.

<a id="specification.advice-on-assumptions-without-basis-or-applicability-and-interfaces-without-a-contract-link"></a>

### [Advice on assumptions without Basis or Applicability and interfaces without a contract link](<../../internal/integration/testdata/boundary-advice.txtar>)

Defines requirements.

---

[View source](<../../specs/use-cases/check-specifications.use-case.spec.md#check-specifications-before-committing>)
