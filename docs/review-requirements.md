# Review requirements and unresolved questions

Use this guide when preparing a specification for stakeholder review or checking
what an evolving draft leaves unresolved. Add only the structures that help a
reader make a decision, then generate a consolidated review. Missing information
produces advisory prompts; it does not make an otherwise valid draft invalid.

## Generate the review

```sh
dotspec model review --path . --relative-to docs/reference > docs/reference/requirements-review.md
dotspec model pegs --path . --relative-to docs/reference > docs/reference/pegs.md
```

Create the output directory first if necessary. `--relative-to` identifies where
the generated Markdown will be read, so source and contract links work there.
Regenerate these files after changing specifications; edit their source documents.

The [review example](reference/requirements-review-example.md) deliberately
includes gaps. The report shows recorded goal contribution paths, scenario
relationships, environment and boundaries, verification guidance, open questions,
and glossary collisions. Its scope is the export's acquired specifications and
direct reference context, not everything that may exist in the organization.

The [PEGS example](reference/pegs-example.md) arranges the same model into Project,
Environment, Goals and System reading perspectives. Business, user and software
requirements remain the model's concern subsets. PEGS neither replaces those
subsets nor introduces a required authoring order.

## Make questions actionable

A short `## Open questions` list remains valid. Use named questions when the
resolution needs an owner, a deadline or evidence:

```markdown
## Open questions

### Clock tolerance

What time offset can the recovery boundary tolerate?

#### Raised by

Retention reviewer.

#### Raised on

2026-11-11.

#### Resolve by

Before acceptance of retention-boundary tests.

#### Importance

Serious: affects requests near the deadline.

#### Evidence needed

Offset measurements and stakeholder agreement on the boundary policy.

#### Tracking reference

[Boundary decision](../decisions/retention-boundary.md).
```

All metadata is optional. Resolve by accepts a date, milestone or condition;
Importance preserves the author's wording, including Meyer's show-stopper,
serious or desirable categories if useful. The review inventory identifies
missing resolution metadata without inventing questions from every omitted field.
Tracking reference is optional and points to the issue or decision that owns
resolution workflow. Dotspec does not copy that workflow into a question status.

## Separate premises from operating conditions

Use Environment for conditions the solution must accommodate, Constraints or
Business rules for relevant restrictions, and Assumptions for premises on which
the specification relies. A fact being outside the software's control does not
by itself make it an assumption that excuses required behavior.

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

#### Validation plan

##### Method

Compare readings with an independent reference before deployment and in operation.

##### Success criteria

The observed offset remains within the agreed tolerance.
```

Basis can cite existing Evidence statements or native research artifacts; use
Derived from for a typed derivation and Relies on to identify consumers of the
assumption. Applicability here is prose describing conditions, not a duplicate
requirement `applies_to` relationship. Validation plans describe proposed checks,
not proof that a premise holds. The review prompts for missing basis,
applicability and validation guidance.

## Delimit the boundary and its interfaces

Add Out of scope, Interacts with, Required interfaces and Provided interfaces to a
System, or at the appropriate subsection level within a Component declaration:

```markdown
## Out of scope

Restoring documents from external backups after their retained copies expire.

## Interacts with

[Platform clock](../clock/specs/system.spec.md).

## Required interfaces

### Time readings

Read UTC time for retention comparisons.

#### Counterpart

[Platform clock](../clock/specs/system.spec.md).

#### Contracts

[Time API](../contracts/time.openapi.yaml).

## Provided interfaces

### Restoration requests

Accept requests from authorized clients.
```

Interacts with resolves explicit System or Component links. It does not imply
Part of, and it need not be repeated in the partner's specification. The two
interface collections express direction relative to the owning subject; named
entries carry descriptions, optional counterpart context and contract links.
Counterpart links are contextual; they do not create a second partner relation.

Keep wire formats, protocol definitions and detailed architecture in their native
artifacts. The review flags interfaces without contract references, but a link
alone cannot establish that the contract is adequate or current. The generated
system overview also includes exclusions, partners and interface descriptions.

## Relate scenarios to general behavior

A scenario illustrates particular interactions; general obligations and concept
actions describe behavior that should apply beyond that one sequence. Identify
a relevant rule or action with an explicit nested bullet:

```markdown
## Main success scenario

1. The member selects a retained document.
2. The workspace restores the document.
   - **Illustrates:** [Restore retained content](#restore-content), [Restore](../concepts/trash.concept.spec.md#restore).
3. The member resumes work.
```

Only the labelled bullet creates `ScenarioStep.illustrates`. Ordinary links in
step prose remain citations; use-case Involves remains participation. The review
asks about steps with no recorded rule or action. Assess whether a relationship
would help; do not create a duplicate obligation for every conversational step.
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
intended. `dotspec model docs` renders the matrix, as in the
[effects-table example](reference/effects-table-example.md).

This is an optional review aid. Dotspec does not infer operations, query
classifications, formal expressions or matrix entries from prose. The current
mapping decomposes the first table in the section and preserves the whole section;
keep one matrix there. General obligations and action contracts still own their
behavioral statements.

## Distinguish explanation and review authority

Use Explanatory notes for explanation, Implementation suggestions (or
Implementation hints) for nonbinding design advice, and Illustrative examples
for examples. These are separate model fields, including on named domain records.
The review presents them separately; concept documentation also labels suggestions
as nonbinding. The heading classifies content role, not stakeholder acceptance.

Use one preferred glossary term consistently in prose. Alternative names or
Aliases can record recognized synonyms and interface labels. The review checks
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
preferred names (§5.7.5), TBD resolution information (§5.8), and goal, scenario and
environment completeness reviews (Chapter 11). Goal completeness concerns the
support for known goals; it does not establish that every possible goal has been
discovered. The operation-by-observable-property table adapts §11.6; dotspec's
source-backed effects matrix does not claim the formal completeness of Meyer's
command-query analysis.

Meyer's [PEGS plan](https://bertrandmeyer.com/2021/09/13/standard-plan-modern-requirements/)
inspires the additional reader perspective. Environment premises and constraints
follow §§1.7 and 3.5. Revision-specific stakeholder endorsement is discussed in
§4.14 and remains deliberately outside the implemented review workflow.
