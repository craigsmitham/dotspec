---
type: requirement-spec
quality_attributes: [usability, repeatability]
---
# Read and verify the generated requirements specification

## Applies to

Whole-project builds and checks for projects declared with `dotspec.yaml`. Reports are enabled by default. A project can opt out by setting `report.enabled` to `false`.

## Statement

For a project with one root system and reports enabled, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and supporting references, and locate the editable source.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. Standalone collection indexes sit beside the overview and use plural names such as `use-cases.md` or `requirements.md`. An index entry links to the subject and includes its authored goal, statement, or description when available. Empty collections produce no index. From the overview, reach the same promise through its feature and through the system’s requirements index. Both routes explain why the promise is listed and lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

System and component overviews open with their subject name and authored description. Readers can expand On this page for the populated body sections in reading order, or Browse specifications for an alphabetical list of collections that require a separate browsing destination. Both lists start collapsed. Authored stakeholder needs stay in the body.

For features, components, interfaces, concept definitions, concept uses and system quality attributes, a system overview section that lists every member serves as the collection. The build omits the separate index and its Browse specifications link. Subject breadcrumbs and citations to the authored collection heading open the corresponding overview section. For interfaces, breadcrumbs and group citations open the matching offered or used list. When the overview covers only part of a collection, Browse specifications links to the complete index, including any By component routes.

Architecture shows immediate parts inside a labeled boundary and labels parts and interfaces with their project names. Declared external interaction partners appear outside the boundary. When both immediate parts and external partners exist, separate structure and context diagrams answer those questions. The linked text groups subsystems, components, interfaces and interactions where populated and preserves authored explanations. When all interfaces are offered, the heading is Interfaces. When a subject uses interfaces, headings identify which it uses and, if applicable, which it offers. A parent link or capabilities alone can be read without a diagram. Authored component Architecture content leads the view.

When a system has capabilities, On this page links to Capabilities. The section follows Architecture when both are present. Each capability declared or explicitly provided in the system or its component hierarchy appears once, with its linked name, an excerpt from its authored description when available, and links to every declared provider. The provider cell stays empty when the capability has no declared provider. A component overview lists its scoped capabilities under Architecture. Both overviews provide direct access to the capability reports and omit the separate capabilities index. A capability report’s Capabilities breadcrumb returns to its system’s section.

Between Stakeholders and Architecture, the system overview presents Scope, Environment, Assumptions and Constraints in that order as populated. Scope presents the authored scope text followed by an Out of scope subsection for explicit exclusions. When exclusions are the only boundary content, Scope opens directly with that subsection. When both fields are empty, the report omits Scope. On this page links once to each visible top-level section. Citations to exclusions lead directly to Out of scope, and local subheadings retain their hierarchy within each passage.

System and component overviews group named concept uses by their resolved definition under Concepts. Each definition appears once with its linked name and authored purpose when supplied. The section includes definitions declared within scope even before a use exists, and links to external definitions used within scope. Each use retains its name, authored description, declaring component where applicable, and link to its detail report. Uses awaiting a definition appear together under Concept uses. Breadcrumbs on definition and use reports return to their group in the owning overview. Synchronizations follow as a separate subsection. Each overview shows every substantive entry belonging in a body section, in index order, regardless of collection size.

The system overview has one Quality Attributes section when the system declares attributes. It shows the full hierarchy with linked attribute names and short excerpts from their authored definitions. On this page opens the section. When the system has a quality model index, a footer link opens it even if the system declares no attributes. A shared attribute appears under each declared parent, with its descendants expanded once. Each attribute page links to its classified requirements.

Index entries use the subject’s linked name and a short excerpt of its authored preview. An ellipsis marks shortened wording. Links to a component collection name both the component and the collection. A concept use identifies its definition with Uses. A requirement index entry identifies each requirement or business rule it refines with Refines. Adding a second refinement preserves the first relationship.

System and component overviews introduce named goals and assumptions
with a linked title and an authored statement preview when available. Each link
opens the exact entry in the owner's collection. On system reports these entry links provide access to the detail page. Browse specifications also links to it when it lists component collections. Readers can follow a relationship
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
Each authored subject report ends with one **View source** link in the footer that opens that exact heading. Embedded subjects link to their own declarations. Collection headings serve as navigation and omit source links.
Above each subject or collection title, the breadcrumb links to its owners in order. A subject’s final breadcrumb opens its collection. The root system is the starting point.
A use case with multiple extensions links to each named extension under Extensions. Concepts and quality models provide the same navigation for their named actions and attributes.
Readers can follow supporting references beside the content they supply.
The project check identifies duplicate page claims and missing or misdirected
source metadata.

## Evidence

Checks for this requirement include [To-do reader journey and regression probes](../../scripts/check-example/check.py), [generated specification guide](../../docs/reference/generated-specification.md).

## Rationale

Reviewers can inspect product intent in Git and follow a promise to its context and evidence without learning the repository layout.
