---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../dotspec.system.spec-report.md#requirements>)

# Read and verify the generated requirements specification

## Statement

For a project with one root system and reports enabled, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and supporting references, and locate the editable source.

## Applies to

Whole-project builds and checks for projects declared with `dotspec.yaml`. Reports are enabled by default. A project can opt out by setting `report.enabled` to `false`.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. Standalone collection indexes sit beside the overview and use plural names such as `goals.md` or `assumptions.md`. An index entry links to the subject and includes its authored goal, statement, or description when available. Empty collections produce no index. From the overview, reach the same promise through its feature and through the system’s requirements index. Both routes lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

System and component overviews open with their subject name and authored description. Readers can expand On this page for the populated body sections in reading order, or Browse specifications for an alphabetical list of collections that require a separate browsing destination. Both panels start collapsed. Authored stakeholder needs stay in the body.

Every generated directory for choosing another specification presents one alphabetical list per specification type. Each entry links the title to its canonical detail report and includes a short excerpt of its authored preview when available. An ellipsis marks shortened wording. A title alone suffices when the preview is absent. This applies to inline directories on overviews, standalone indexes and directories on other subject reports.

A system directory includes specifications in its system and component hierarchy, stopping at a separately specified subsystem. A component directory covers its own component hierarchy. Explicit applicability and capability provision also establish membership. A specification reached through multiple paths appears once. Distinct specifications with the same title have separate entries, with stable ordering. Ownership, applicability and refinement are explained on the relevant subject reports. Authored citations, evidence and diagrams stay with the subject or passage they explain.

The system’s On this page groups links in body order. Intent covers Framing, Purpose, Goals, Stakeholders and Jobs. Context covers Scope, Environment, Assumptions, Constraints and Business rules. Behavior covers Features, Capabilities and Use cases. Quality Attributes and Requirements each have a direct link on its own line. Design covers Concepts, Concept uses, Synchronizations and Architecture. Each populated group has a plain label above a row of links that can wrap. Empty groups are omitted. The labels appear only in navigation. Additional populated top-level sections follow as direct links.

System reports present populated Business rules, Jobs, Requirements and Use cases directories as sections. A collection with a home on an overview uses that section as the destination for member breadcrumbs and collection citations. Standalone indexes provide the same direct member links. For To-do, the system’s Requirements entries are Acknowledge after commit, Keep recorded work available, Preserve tasks across restart and Show which replies need follow-up, in that order. Task service also leads to Acknowledge after commit. The restart report identifies the promise it refines.

Architecture shows immediate parts inside a labeled boundary and identifies declared external interaction partners. When both are present, separate structure and context diagrams explain decomposition and interactions. Subsystems, Components and Interfaces each have a directory when populated. Interface details identify the subject that provides or requires the interface. Authored component Architecture content leads the view.

A component lists scoped capabilities under Architecture. Each capability report lists its declared providers. Its Capabilities breadcrumb returns to its system’s section.

Between Stakeholders and Architecture, the system overview presents Scope, Environment, Assumptions and Constraints in that order as populated. Scope presents the authored scope text followed by an Out of scope subsection for explicit exclusions. When exclusions are the only boundary content, Scope opens directly with that subsection. When both fields are empty, the report omits Scope. On this page links once to each visible top-level section. Citations to exclusions lead directly to Out of scope, and local subheadings retain their hierarchy within each passage.

System and component overviews present Concepts, Concept uses and Synchronizations as separate directories. Concepts includes definitions declared in scope and external definitions referenced by scoped uses. External definitions without a local report retain their authored destination. A use awaiting a definition remains in Concept uses. The use report identifies its resolved definition under Uses.

The system overview has one Quality Attributes section when attributes are declared. Each attribute appears once, including an attribute with multiple parents. Its page links to parent and child attributes and classified requirements. A footer link opens the quality model catalog when present.

System and component overviews introduce named goals and assumptions
with a linked title and an authored statement preview when available. Each link
opens the exact entry in the owner's collection. On system reports these entry links provide access to the detail page. Readers can follow a relationship
in either direction from the passages it connects. The detail entry keeps its
qualifications and supporting references together. A citation in prose leads
back to the citing passage through a “Referenced by” link. Collection entries come from declared members.
The report selects authored fields, and inspection identifies misplaced collection prose and unrecognized sections at their source locations. The component’s
requirements index lists each explicitly applicable requirement once. Its overview
keeps authored citations beside the responsibilities they explain.

Across subject reports, body sections describe the project using authored content,
specific subjects, or declared relationships. Collection links that only serve
navigation belong in the section map or page navigation. Definitions of dotspec
categories belong in the handbook.

A reader can look up a term from an inline or imported glossary and follow it
to its defining heading. The topic index links to native Markdown headings. A term’s definition leads into a single line with its plain topic labels, model reference when supplied, and the source link when this is the first term from its file. When two
contexts define the same word, each entry names its context. A relationship
between terms links in both reading directions, with the original statement
and examples available in the declaring entry.

When any subject in the project records an open question, the root system’s
section map links to `specification/open-questions.md`. Each entry shows the
authored question, names its subject and links to the passage where it appears.
Resolving the last question removes the index and its navigation link on rebuild.

Each subject report identifies its authored file and heading once in the YAML
`source` field using a URI relative to the report. This field supplies the
source attribution for the page and its passages.
Each authored subject report ends with one **View source** link in the footer that opens that exact heading. This footer covers all embedded declarations and passages from its file. Each embedded declaration records its exact authored heading in its report metadata.
On every report or index, the first full passage from each additional authored file introduces one visible source link. Headings within the same resolved file share that link. Linked previews lead to their subject reports. Authored citations, evidence and relationship links remain beside the content they support.
Above each subject or collection title, the breadcrumb links to its owners in order. A subject’s final breadcrumb opens its collection. The root system is the starting point.
A use case with multiple extensions links to each named extension under Extensions. Concepts and quality models provide the same navigation for their named actions and attributes.
Readers can follow supporting references beside the content they supply.
The project check identifies duplicate page claims and missing or misdirected
source metadata.

## Evidence

Checks for this requirement include [To-do reader journey and regression probes](<../../scripts/check-example/check.py>), [generated specification guide](<../../docs/reference/generated-specification.md>).

## Rationale

Reviewers can inspect product intent in Git and follow a promise to its context and evidence without learning the repository layout.

## Related specifications

**Quality attribute**

- [Repeatability](<../quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>)
- [Usability](<../quality/usability/usability.quality-attribute.spec-report.md>)

---

[View source](<../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification>)
