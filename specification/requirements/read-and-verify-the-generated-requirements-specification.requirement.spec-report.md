---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification"
---
# Read and verify the generated requirements specification

[dotspec](<../dotspec.system.spec-report.md>)

## Statement

For an opted-in project with one root system, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and supporting references, and locate the editable source.

## Applies to

Whole-project builds and checks for projects that enable generated specifications.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. Each overview links to a plural-named index beside it for each populated collection, such as `use-cases.md` or `requirements.md`. An index entry links to the subject and includes its authored goal, statement, or description when available. Empty collections produce no index. From the overview, reach the same promise through its feature and through the system’s requirements index. Both routes explain why the promise is listed and lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

The system overview opens with the system name and authored description. Its
section map links directly to populated indexes. Readers can follow goals,
business rules, user needs, scope, system topology, features, capabilities,
use cases, concepts, requirements, quality attributes, and open questions in
that order when present. Features, capabilities, use cases, concepts, and
requirements have separate rows. Concept definitions, named uses, and
synchronizations share the Concepts row. The Requirements row leads to
standalone requirements and embedded obligations. The topology shows declared immediate
parts and boundary relationships in a Mermaid diagram and in linked text.
The overview introduces the top-level quality attributes with their definition
previews. The attribute index shows the full hierarchy. A shared attribute
appears under each declared parent, with its descendants expanded once. The
business rules index lists all rules. When a small set explicitly applies to the system, the overview names those rules and includes
their authored statement previews.

System and component overviews introduce named goals, assumptions, and constraints
with a linked title and an authored statement preview when available. Each link
opens the exact entry in the owner's collection. Readers can follow a relationship
in either direction from the passages it connects. The detail entry keeps its
qualifications and supporting references together. A citation in prose leads
back to the citing passage through a “Referenced by” link. Prose and lists
directly under the collection heading stay in the overview. The component’s
requirements index lists each explicitly applicable requirement once. Its overview
keeps authored citations beside the responsibilities they explain.

Across subject reports, body sections describe the project using authored content,
specific subjects, or declared relationships. Collection links that only serve
navigation belong in the section map or page navigation. Definitions of dotspec
categories belong in the handbook.

A reader can look up a term from an inline or imported glossary and follow it
to its defining heading. Declared topics link to that same definition. When two
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
