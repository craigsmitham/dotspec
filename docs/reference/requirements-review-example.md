# Requirements review

Specification project: **todo**. Project membership does not imply a System relationship.

Generated from acquired specifications. Findings below are advisory review prompts, not completeness, acceptance or implementation claims. Missing links mean **no relationship recorded in this scope**. Drafts may conform to the authoring format while requirements verification and validation remain unfinished.

Model `0.19.0-draft`; mapping `dotspec-markdown/14`.

## Needs and intent

Requirements validation asks whether the specification expresses the right needs. System validation asks whether a proposed or realized solution serves intended use. The relationships below identify review questions, not either conclusion.

### Goal contribution paths

Paths follow explicit supports, derived-from and refinement links from obligations, concept actions or use cases. A recorded path identifies intent; reviewers must assess whether the behavior is sufficient.

- [Resume unfinished work](<../../examples/todo/specs/system.spec.md#resume-unfinished-work>): Contribution path recorded; sufficiency requires review.

### Scenario grounding

An Illustrates link identifies a rule or action behind a step. Participation links and ordinary citations do not establish this relationship; not every conversational step needs a separate obligation.

- [Main success scenario / ScenarioStep 1](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md>): illustrates [Complete](<../../examples/todo/specs/concepts/checklist.concept.spec.md#complete>).
- [Main success scenario / ScenarioStep 2](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md>): illustrates [Keep recorded work available](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available>).
- [Main success scenario / ScenarioStep 3](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md>): illustrates [Reopen](<../../examples/todo/specs/concepts/checklist.concept.spec.md#reopen>).
- [The person intentionally deleted the errand / ScenarioStep 1](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md>): no rule or action relationship recorded; review required.
- [The person intentionally deleted the errand / ScenarioStep 2](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md>): no rule or action relationship recorded; review required.
## Environment and boundaries

Operating conditions and constraints describe what the solution must accommodate. Assumptions are assessed separately for evidence and justified reliance. Explicit exclusions delimit scope. External partners are separate from structural parents.

### BusinessRule: [Completion preserves a commitment](<../../examples/todo/specs/business-rules/completion.business-rule.spec.md#completion-preserves-a-commitment>)

**Statement**

A finished commitment stays on the person's list until they explicitly delete it.

**Applicability**

**Targets:** [Personal task list](../../examples/todo/specs/system.spec.md).

### System: [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Environment**

The person may close the browser or stop the service between work sessions.

**Out of scope**

Sharing lists, synchronizing devices and recovering from loss of the database file.

**Provided interfaces**

#### Task list forms

The browser presents open and completed tasks and accepts changes through HTML forms.

##### Contracts

[Form template](../../examples/todo/app/list.html) supplies the actions and fields.

### Component: [Task service](<../../examples/todo/specs/system.spec.md#task-service>)

## Assumptions: evidence and justified reliance

Assess both empirical support in the stated conditions and whether relying on the premise is appropriate for intended use. An observation supports a bounded claim, not a guarantee that the premise always holds.

### [Same local storage](<../../examples/todo/specs/system.spec.md#same-local-storage>)

**Statement**

The person returns to the same writable database file.

Review required: no Basis, Applicability, Validation plan recorded.

## Validation planning

Identify the subject and the needs, intended use or other basis for the assessment. Requirements validation assesses proposed obligations; system validation assesses proposed or realized solutions. A method or supports link alone establishes neither purpose nor a result.

### Concept: [Checklist](<../../examples/todo/specs/concepts/checklist.concept.spec.md#checklist>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### Requirement: [Preserve tasks across restart](<../../examples/todo/specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### System: [Personal task list](<../../examples/todo/specs/system.spec.md#personal-task-list>)

**Validation plan**

#### Subject

A person returning to unfinished work using this build.

#### Method

Invite someone to record two real errands and mark one finished. After a break,
restart the app with the same database. Ask what still needs attention. Then tell
them that the finished errand needs a follow-up and ask them to resume it.
Let them choose the controls without coaching.

#### Success criteria

The participant identifies the outstanding errand and resumes the other without
entering it again. Ask what they expected to find on returning.

#### Evidence to retain

With consent, record the build, task outcome, any help given and the participant's
explanation. Review mismatches before deciding whether the workflow meets the need.

Review separately identified plan details: Basis, Participants, Conditions. The narrative may already address them; no adequacy or completion is inferred.

### UseCase: [Return to my list](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#return-to-my-list>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

### ObligationClause: [Keep recorded work available](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available>)

No validation plan recorded. Review whether an explicit plan would help assess needs or intended use; it is optional.

## System verification planning

These plans assess realizations against obligations. Requirements verification instead assesses the quality of the requirement statements and set. A plan is distinct from an observation, an evidence-based assessment, and an acceptance decision.

### [Preserve tasks across restart](<../../examples/todo/specs/requirements/restart.requirement.spec.md#preserve-tasks-across-restart>)

**Verification plan**

#### Method

The [restart acceptance test](../../examples/todo/tests/acceptance.py) submits the public forms,
stops the process, starts a new process on the same file and compares both views.

**Criterion**

#### Measure

Count missing, unexpected or changed entries compared with the last acknowledged list.

#### Target

Zero discrepancies across the two views.

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

### [Keep recorded work available](<../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available>)

**Verification plan**

Record a task, complete it and return to both views. Reopen it and compare its
identity with the original. For an interruption between sessions, follow the
[restart verification plan](../../examples/todo/specs/requirements/restart.requirement.spec.md#verification-plan).

Assessment guidance recorded; its adequacy and any verification result are not evaluated.

## Open questions

Consolidated authored questions. External tracking references own resolution workflow. Omitted specification fields are not automatically questions; missing metadata below is advisory.

No explicit open questions recorded in this scope.

## Requirements quality

Requirements verification examines individual statements and the set for clarity, consistency and assessability against declared quality criteria. Terminology checks and inspection findings below supply selected review prompts; passing them does not establish complete requirements verification.

### Terminology review

Use one preferred term consistently in specification prose; declared alternative names and synonym entries remain useful. Collisions below are scoped to the declaring System and require human review. This report does not scan prose for undeclared synonyms.

No collisions found among recorded preferred terms and alternative names in each System. This does not establish terminology consistency.

## Explanations and suggestions

These roles distinguish explanation, nonbinding implementation advice and examples from specified obligations. They do not record stakeholder endorsement.
