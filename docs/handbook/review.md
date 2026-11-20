# Review requirements and unresolved questions

Use this guide when preparing a specification for stakeholder review or checking
what an evolving draft leaves unresolved. Add only the structures that help a
reader make a decision, then generate a consolidated review. Missing information
produces advisory prompts; it does not make a draft that conforms to the authoring format structurally invalid.

## Generate the review

```sh
dotspec model review
dotspec model pegs
```

Both commands write the selected perspective to standard output. If you save one
to a file, pass its containing directory with `--relative-to` so readers can
follow source and contract links.

When a project records an open question, `dotspec spec build` includes **All open questions** in the root system’s section map. Follow a question to read it in context. To prepare a broader review, run `dotspec model review` for goal contributions, scenario steps, requirement criteria, and inspection findings.

`dotspec model pegs --project todo` arranges the example into Project,
Environment, Goals and System reading perspectives. Business, user and software
requirements remain the model's concern subsets. PEGS neither replaces those
subsets nor introduces a required authoring order.

## Make questions actionable

Record the decision that remains open and the consequence for the requirement.
A question can link to its issue, research, or decision record.

```markdown
## Open questions

- What time offset can the recovery boundary tolerate? The answer determines
  which near-deadline requests qualify. See the [clock study](../decisions/clock.md).
```

Use the [four review questions](verification-and-validation.md#assessment-purposes)
to distinguish wording, stakeholder intent, conformance, and intended use.

## Separate premises from operating conditions

Use Environment for conditions the solution must accommodate, Constraints for
externally imposed restrictions, and Assumptions for premises on which the
specification relies. Environment and Constraints are narrative, as is a
Stakeholders section naming the people and organizations with an interest in the
system. Write a restriction that needs its own identity as a requirement or
business rule. A fact being outside the software's control does not by itself
make it an assumption that excuses required behavior.

For example, intermittent connectivity is an operating condition if reconnecting
clients are supported. An assumption that the platform supplies sufficiently
accurate time is a premise that needs justification and validation:

```markdown
## Assumptions

### Trusted time

The platform clock stays within the agreed offset tolerance.

#### Basis

[Clock measurements](#clock-measurements) motivate this premise.

#### Applicability

Retention calculations in the managed production environment.

```

Basis can cite a study or link to supporting evidence. Use
Derived from for a typed derivation and Relies on to identify consumers of the
assumption. Applicability here is prose describing conditions, not a duplicate
requirement `applies_to` relationship. Review the evidence for the premise and whether reliance is justified.

## Delimit the boundary and its interfaces

Add Out of scope, Interacts with, Required interfaces and Provided interfaces to a
System, or at the appropriate subsection level within a Component declaration:

```markdown
## Out of scope

Restoring documents from external backups after their retained copies expire.

## Interacts with

[Platform clock](../clock/specs/platform-clock.system.spec.md).

## Required interfaces

### Time readings

**ID:** time-readings

Read UTC time for retention comparisons.

#### Counterpart

[Platform clock](../clock/specs/platform-clock.system.spec.md).

#### Contracts

[Time API](../contracts/time.openapi.yaml).

## Provided interfaces

### Restoration requests

Accept requests from authorized clients.
```

Interacts with resolves explicit System or Component links. It does not imply
Part of, and it need not be repeated in the partner's specification. The two
interface sections express direction relative to the owning subject; named
entries are independently addressable Interface elements, with heading
anchors and optional introductory `**ID:**` values. Each belongs to one required
or provided list under one System or Component. Entries carry descriptions,
optional counterpart context, contract links, open questions and trace sections.

Requirements can target an interface using
`**Targets:** [Time readings](../platform.system.spec.md#time-readings)` under Applies to, and features may
include it under Groups. Counterpart links are contextual; they do not establish
implementation or compatibility between declarations. Two declarations can cite
the same native contract without sharing identity. The [task list interface](../../examples/todo/specs/personal-task-list.system.spec.md#task-list-forms) links to its HTML form contract.

Keep wire formats, protocol definitions and detailed architecture in their native
artifacts. The review flags interfaces without contract references, but a link
alone cannot establish that the contract is adequate or current. The generated
system overview also includes exclusions, partners and interface descriptions.

## Relate scenarios to general behavior

A scenario illustrates particular interactions; general requirements, business
rules and concept actions describe behavior that should apply beyond that one
sequence. Identify a relevant requirement, business rule or action with an explicit
nested bullet:

```markdown
## Main success scenario

1. The member selects a retained document.
2. The workspace restores the document.
   - **Illustrates:** [Restore retained content](../requirements/restore-content.requirement.spec.md), [Restore](../concepts/trash.concept.spec.md#restore).
3. The member resumes work.
```

Only the labelled bullet creates `ScenarioStep.illustrates`. Ordinary links in
step prose remain citations; use-case Involves remains participation. The review
asks about steps with no recorded rule or action. Assess whether a relationship
would help; do not create a duplicate requirement for every conversational step.
Recorded links do not prove scenario completeness.

## Review a concept's effects

When interactions between operations and observable properties are easy to miss,
add one authored matrix under Effects table in a concept:

```markdown
## Effects table

| Operation | Available(item) | Retained(item) |
| --- | --- | --- |
| Delete | false | true |
| Restore | true | false |
| Expire | unchanged | |
```

The first row names columns and the first column names operations. Cells preserve
full authored content, including formulas or links. **A blank cell is unspecified;
it never means unchanged.** Say `unchanged` explicitly where preservation is
intended. The [Checklist effects table](../../examples/todo/specification/concepts/checklist.concept.spec-report.md#effects-table)
shows what reopening an entry preserves. `dotspec model docs` also includes the matrix
in its concept reference.

The report reproduces the authored effects table. Include qualifications and
examples alongside the entries they explain. General requirements and action
contracts own their behavioral statements.

## Distinguish explanation and review authority

Put an example or implementation suggestion inside the field it explains.
State which behavior is required and which design choices are suggestions.

Use one preferred glossary term consistently in prose. Use Alternative names
for recognized synonyms and interface labels. The review checks
preferred-term and alias collisions within each declaring System. It does not
infer bounded contexts or scan prose for undeclared synonyms.

A System's Project context can link schedules, responsibilities, risks and decision
records. Keep those artifacts authoritative outside dotspec. Stakeholder endorsement
also remains external: any future endorsement model would need the reviewed
revision, qualified reviewers, scope, qualifications and disagreements. A generic
`approved: true` field would lose that information.

## Basis and limits

These additions adapt Bertrand Meyer's
[Handbook of Requirements and Business Analysis](https://doi.org/10.1007/978-3-031-06739-6):
explicit exclusions (§4.9), supplementary explanation (§5.4), glossary aliases and
preferred names (§5.7.5), treatment of open questions (§5.8), and goal, scenario and
environment completeness reviews (Chapter 11). Goal completeness concerns the
support for known goals; it does not establish that every possible goal has been
discovered. The operation-by-observable-property table adapts §11.6; dotspec's
source-backed effects matrix does not claim the formal completeness of Meyer's
command-query analysis.

Meyer's [PEGS plan](https://bertrandmeyer.com/2021/09/13/standard-plan-modern-requirements/)
inspires the additional reader perspective. Environment premises and constraints
follow §§1.7 and 3.5. Revision-specific stakeholder endorsement is discussed in
§4.14 and remains deliberately outside the implemented review workflow.
