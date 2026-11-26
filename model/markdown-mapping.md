# Markdown mapping profile

`dotspec-markdown/28` maps authored declarations to the classes and fields in the [model reference](../docs/reference/model/index.md). A `*.<kind>.spec.md` file declares one specification root with `type: <kind>-spec` and names it in the level-one heading. The root owns its nested elements. A reusable quality model is an independent reference root in `*.quality-model.md` with `type: quality-model`. Each named attribute in that file belongs to the model. The source text is normalized to UTF-8 with LF line endings and no leading byte-order mark.

## Shared mapping rules

Choose a level-two field belonging to the subject. Within a narrative field, Markdown paragraphs, lists, tables, code, examples and local subheadings retain their authored order. For example, a limitation on a quality attribute belongs in its Definition. A named declaration collection contains one immediately subordinate heading per member. For example, put an explanation of how frost warnings reach growers under `### Spring frost` inside a system’s `## Overview` to keep it as narrative. To declare an agreed goal, put `### Warn before dawn` under `## Goals` and write the desired result beneath it. If `structure.collection-content` identifies prose above the first goal, move it under the goal it describes or to the narrative field it qualifies. From the project directory, rerun `dotspec model export` and inspect the system’s `goals` and the named goal’s `owner`. Unrecognized field headings also produce located diagnostics.

Omitting a narrative field leaves it absent. An explicitly empty field has an empty list of spans. Spans identify half-open byte intervals in the normalized source. The source also preserves spacing, punctuation, comments and any invalid text. Markdown linked as external documentation supplies source and heading anchors. A glossary explicitly selected by a system additionally supplies vocabulary terms.

## System

Purpose holds authored Markdown. The optional `## Overview` section maps to `System.overview` and preserves authored Markdown, including local subheadings. Components and features accept the same section as `Component.overview` and `Feature.overview`. Reports place it after Purpose wherever it is authored. The Project references section maps to `System.project_references`: its content preserves the explanation and its references identify the cited project records. Goals, interfaces, concept uses, synchronizations and inline glossary terms have named entries. Declare a system quality attribute in its own specification under the system’s quality directory. A reusable quality model declares its attributes by name under Quality attributes.

A System owns the following optional sections in its `*.system.spec.md` file.

| Authored heading | System field | Content |
| --- | --- | --- |
| `## Scope` | `scope` | Rich Markdown describing the work covered |
| `## Out of scope` | `out_of_scope` | Rich Markdown recording agreed exclusions |
| `## Environment` | `environment` | Rich Markdown describing relevant conditions in the surrounding world |
| `## Assumptions` | `assumptions` | Named `Assumption` elements in declaration order |
| `## Constraints` | `constraints` | Rich Markdown describing externally imposed limits |

Under `## Assumptions`, each level-three heading names one assumption. Its introductory Markdown supplies `Assumption.statement`. Optional level-four headings `#### Basis` and `#### Applicability` map to the corresponding fields. A requirement can link to the named assumption under **Relies on**. Local subheadings in the other four sections organize their narrative.

Under `## Stakeholders`, each level-three heading declares one `SystemStakeholderRole` owned by the System. The heading supplies its name and Markdown anchor. Its body supplies the rich-text `description`. The System lists the roles in `stakeholder_roles` in source order. Link to a role’s heading from the passage where the relationship matters.

## Job

Use a named heading for each step and outcome. An outcome can belong to the job or an individual step. Its body holds its statement and an optional **ID** gives it a stable authored identifier. An Applies to link can identify a step independently of containment. See the [worked heading structure](../docs/cli-handbook/guides/write-a-job.md#record-steps-and-shared-outcomes).

## Use case

Goal, primary actor, supporting actors, goal level, guarantees and preconditions hold Markdown. The main success scenario retains the authored passage. Its ordered list items identify steps. An extension has a name and can describe its Condition and Handling in separate fields. Step numbers in prose carry their ordinary authored meaning.

## Concept

Purpose, state, operational principles and the effects table hold Markdown. A named action owns its description and authored fields such as Inputs, Outputs, Guards and Effects. Blank table cells and examples retain their original wording and layout.

## Business rule

Statement, rationale, exceptions and authority describe a rule. A rule can cite related specifications and explicitly identify applicability and traceability targets. Under **Applies to**, a line beginning `**Targets:**` names the subjects a rule or a requirement applies to. An ordinary citation remains navigable wherever it occurs.

## Requirement

Only a `*.requirement.spec.md` file declares a Requirement. It keeps its statement, applicability, criterion and rationale. A `**Targets:**` line under **Applies to** names its subjects. Conditions, workload, population and exclusions qualify applicability. Measure and target qualify a criterion. A quality classification is a label occurrence resolved within its governing vocabulary.

## Relationship contract

A relationship section preserves its authored Markdown, including qualifications and unresolved choices, in its named content field. A semantic link carries its authored destination, source span, owner, role and resolution. An ordinary citation has no semantic role. A successful resolution identifies one endpoint. Other statuses record missing, invalid, ambiguous, external or uninspected targets. Membership fields record declarations in source order. Reverse navigation follows the recorded memberships and resolved occurrences. See the [relationship reference](../docs/reference/model/relationships.md).

## Product membership

A quality model owns its `QualityModelAttribute` records in declaration order. A `SystemQualityAttribute` is an independent specification whose `specification_system` identifies the system that includes its source in its quality directory. A requirement’s classification labels resolve only to attributes of its effective system. The two concrete attribute types can each group under several parents of their own type within the same scope.

## Adopting a provision

A specification or a named interface can have an **Adopts** collection. Each immediately subordinate heading names one adopted provision. In the entry’s introductory paragraph, link once to the exact clause, operation, or other identifiable provision. Give the entry two subsections. **Edition** identifies the revision or says `evolving`. **Applicability** states when the product adopts the provision. Each entry maps to a source-local `AdoptedProvision` with a pointer to the particular link occurrence. Ordinary citations elsewhere retain their own context. [Worked interface](../examples/todo/specs/personal-task-list.system.spec.md#redirect-after-an-accepted-form-submission)

## Verification procedures

In a requirement, **Verification procedures** can be a level-two field for the requirement or an immediate subsection of **Criterion**. Link to each executable test or manual protocol. The mapping retains the passage and indexes its link occurrences as `verification_procedures.references` on the respective subject. **Evidence** retains its authored narrative. Run observations have their own interchange document. [Recording verification results](../docs/reference/generated-specification.md#record-verification-results)
