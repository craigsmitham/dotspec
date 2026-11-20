# Markdown mapping profile

`dotspec-markdown/23` maps authored specifications to the domain classes and fields in the [model reference](../docs/reference/model/index.md). Each file declares its specification type in YAML frontmatter and names its subject in a level-one heading. The source text is normalized to UTF-8 with LF line endings and no leading byte-order mark.

## Shared mapping rules

Choose a level-two field belonging to the subject. Within a narrative field, Markdown paragraphs, lists, tables, code, examples and local subheadings retain their authored order. For example, a limitation on a quality attribute belongs in its Definition. A named declaration collection contains one immediately subordinate heading per member. Place each statement under the member it describes. Unknown field headings and prose before a collection’s first member produce located diagnostics.

Omitting a narrative field leaves it absent. An explicitly empty field has an empty list of spans. Spans identify half-open byte intervals in the normalized source. The source also preserves spacing, punctuation, comments and any invalid text. Markdown linked as external documentation supplies source and heading anchors. A glossary explicitly selected by a system additionally supplies vocabulary terms.

## System

Purpose, scope, environment, stakeholders and constraints hold authored Markdown. Goals, assumptions, interfaces, concept uses, synchronizations and inline glossary terms have named entries. Declare a system quality attribute in its own specification under the system’s quality directory. A reusable quality model declares its attributes by name under Quality attributes. Find the [system format](../docs/reference/formats/system.md).

## Job

Use a named heading for each step and outcome. An outcome can belong to the job or an individual step. Its body holds its statement and an optional **ID** gives it a stable authored identifier. An Applies to link can identify a step independently of containment. See the [job format](../docs/reference/formats/job.md).

## Use case

Goal, primary actor, supporting actors, goal level, guarantees and preconditions hold Markdown. The main success scenario retains the authored passage. Its ordered list items identify steps. An extension has a name and can describe its Condition and Handling in separate fields. Step numbers in prose carry their ordinary authored meaning. See the [use-case format](../docs/reference/formats/use-case.md).

## Concept

Purpose, state, operational principles and the effects table hold Markdown. A named action owns its description and authored fields such as Inputs, Outputs, Guards and Effects. Blank table cells and examples retain their original wording and layout.

## Business rule

Statement, rationale, exceptions and authority describe a rule. A rule can cite related specifications and explicitly identify applicability and traceability targets. An ordinary citation remains navigable wherever it occurs.

## Requirement

Only a `*.requirement.spec.md` file declares a Requirement. It keeps its statement, applicability, criterion and rationale. Conditions, workload, population and exclusions qualify applicability. Measure and target qualify a criterion. A quality classification is a label occurrence resolved within its governing vocabulary.

## Relationship contract

A relationship section preserves its authored Markdown, including qualifications and unresolved choices, in its named content field. A semantic link carries its authored destination, source span, owner, role and resolution. An ordinary citation has no semantic role. A successful resolution identifies one endpoint. Other statuses record missing, invalid, ambiguous, external or uninspected targets. Membership fields record declarations in source order. Reverse navigation follows the recorded memberships and resolved occurrences. See the [relationship reference](../docs/reference/model/relationships.md).
