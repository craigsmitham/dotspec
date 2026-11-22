# Markdown mapping profile

`dotspec-markdown/24` maps authored specifications to the domain classes and fields in the [model reference](../docs/reference/model/index.md). Each file declares its specification type in YAML frontmatter and names its subject in a level-one heading. The source text is normalized to UTF-8 with LF line endings and no leading byte-order mark.

## Shared mapping rules

Choose a level-two field belonging to the subject. Within a narrative field, Markdown paragraphs, lists, tables, code, examples and local subheadings retain their authored order. For example, a limitation on a quality attribute belongs in its Definition. A named declaration collection contains one immediately subordinate heading per member. For example, put a question about spring frost under `### Spring frost` inside a system’s `## Framing` to retain the question as narrative. To declare an agreed goal, put `### Warn before dawn` under `## Goals` and write the desired result beneath it. If `structure.collection-content` identifies prose above the first goal, move it under the goal it describes or to the narrative field it qualifies. From the project directory, rerun `dotspec model export` and inspect the system’s `goals` and the named goal’s `owner`. Unrecognized field headings also produce located diagnostics.

Omitting a narrative field leaves it absent. An explicitly empty field has an empty list of spans. Spans identify half-open byte intervals in the normalized source. The source also preserves spacing, punctuation, comments and any invalid text. Markdown linked as external documentation supplies source and heading anchors. A glossary explicitly selected by a system additionally supplies vocabulary terms.

## System

The optional `## Framing` section maps to `System.framing` and preserves authored Markdown, including local subheadings. Purpose, scope, environment, stakeholders and constraints hold authored Markdown. The Project references section maps to `System.project_references`: its content preserves the explanation and its references identify the cited project records. Goals, assumptions, interfaces, concept uses, synchronizations and inline glossary terms have named entries. Declare a system quality attribute in its own specification under the system’s quality directory. A reusable quality model declares its attributes by name under Quality attributes.

## Job

Use a named heading for each step and outcome. An outcome can belong to the job or an individual step. Its body holds its statement and an optional **ID** gives it a stable authored identifier. An Applies to link can identify a step independently of containment. See the [worked heading structure](../docs/cli-handbook/guides/write-a-job.md#record-steps-and-shared-outcomes).

## Use case

Goal, primary actor, supporting actors, goal level, guarantees and preconditions hold Markdown. The main success scenario retains the authored passage. Its ordered list items identify steps. An extension has a name and can describe its Condition and Handling in separate fields. Step numbers in prose carry their ordinary authored meaning.

## Concept

Purpose, state, operational principles and the effects table hold Markdown. A named action owns its description and authored fields such as Inputs, Outputs, Guards and Effects. Blank table cells and examples retain their original wording and layout.

## Business rule

Statement, rationale, exceptions and authority describe a rule. A rule can cite related specifications and explicitly identify applicability and traceability targets. An ordinary citation remains navigable wherever it occurs.

## Requirement

Only a `*.requirement.spec.md` file declares a Requirement. It keeps its statement, applicability, criterion and rationale. Conditions, workload, population and exclusions qualify applicability. Measure and target qualify a criterion. A quality classification is a label occurrence resolved within its governing vocabulary.

## Relationship contract

A relationship section preserves its authored Markdown, including qualifications and unresolved choices, in its named content field. A semantic link carries its authored destination, source span, owner, role and resolution. An ordinary citation has no semantic role. A successful resolution identifies one endpoint. Other statuses record missing, invalid, ambiguous, external or uninspected targets. Membership fields record declarations in source order. Reverse navigation follows the recorded memberships and resolved occurrences. See the [relationship reference](../docs/reference/model/relationships.md).
