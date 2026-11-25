---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../requirements.md>)

# Read and verify the generated requirements specification

[View source](<../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification>)

## Statement

For a project with one root system and reports enabled, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and supporting references, and locate the editable source.

## Applies to

Whole-project builds and checks for projects declared with `dotspec.yaml`. Reports are enabled by default. A project can opt out by setting `report.enabled` to `false`.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. Each overview links to a plural-named index beside it for each populated collection, such as `use-cases.md` or `requirements.md`. An index entry links to the subject and includes its authored goal, statement, or description when available. Empty collections produce no index. From the overview, reach the same promise through its feature and through the system’s requirements index. Both routes explain why the promise is listed and lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

System and component overviews open with their subject name and authored description. Readers can expand On this page for the populated body sections in reading order, or Browse specifications for an alphabetical list of the scope’s collection indexes. Both lists start collapsed. The overview’s View source link follows the body. Authored stakeholder needs stay in the body.

Architecture shows immediate parts inside a labeled boundary and labels parts and interfaces with their project names. Declared external interaction partners appear outside the boundary. When both immediate parts and external partners exist, separate structure and context diagrams answer those questions. The linked text groups subsystems, components, interfaces and interactions where populated and preserves authored explanations. Readers see the capabilities of each part together as linked names under one Capabilities label. A subject’s own capabilities appear under Capabilities. When all interfaces are offered, the heading is Interfaces. When a subject uses interfaces, headings identify which it uses and, if applicable, which it offers. A parent link or capabilities alone can be read without a diagram. Authored component Architecture content leads the view.

Between Stakeholders and Architecture, the system overview presents Scope, Environment, Assumptions and Constraints in that order as populated. Scope presents the authored scope text followed by an Out of scope subsection for explicit exclusions. When exclusions are the only boundary content, Scope opens directly with that subsection. When both fields are empty, the report omits Scope. On this page links once to each visible top-level section. Citations to exclusions lead directly to Out of scope, and local subheadings retain their hierarchy within each passage.

The Concepts body appears when the overview introduces a named concept use or synchronization. Its populated groups have the stable headings Concept uses and Synchronizations in that order. The definitions index is titled Concept definitions. Each overview shows every substantive entry belonging in a body section, in index order, regardless of collection size.

The system overview has one Quality Attributes section when the system declares attributes. It shows the full hierarchy with linked attribute names and short excerpts from their authored definitions. On this page opens the section. Browse specifications omits the quality attribute and quality model indexes. When the system has a quality model index, a footer link opens it even if the system declares no attributes. The attribute reports’ breadcrumbs lead to the attribute index. A shared attribute appears under each declared parent, with its descendants expanded once. Each attribute page links to its classified requirements.

Index entries use the subject’s linked name and a short excerpt of its authored preview. An ellipsis marks shortened wording. Links to a component collection name both the component and the collection. A concept use identifies its definition with Uses. A requirement index entry identifies each requirement or business rule it refines with Refines. Adding a second refinement preserves the first relationship.

System and component overviews introduce named goals and assumptions
with a linked title and an authored statement preview when available. Each link
opens the exact entry in the owner's collection. Readers can follow a relationship
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
to its defining heading. The topic index links to native Markdown headings. A term’s definition leads into a single line with its plain topic labels, model reference when supplied, and term source. When two
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
The visible View source link opens that same heading. Each embedded subject links to its own authored heading.
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
