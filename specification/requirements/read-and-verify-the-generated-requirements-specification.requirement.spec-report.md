---
generated_by: dotspec
type: requirement-spec-report
---
# Read and verify the generated requirements specification

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Statement

For an opted-in project with one root system, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and verification guidance, and locate the editable source.

## Applies to

Whole-project builds and checks for projects that enable generated specifications.

## Criterion

The project's `SPEC.md` gives the system name and one link to `specification/<name>.system.spec-report.md`. Each system and component overview uses the basename of its named source with `.spec.md` changed to `.spec-report.md`. Each overview links to a plural-named index beside it for each populated collection, such as `use-cases.md` or `requirements.md`. An index entry links to the subject and includes its authored goal, statement, or description when available. Empty collections produce no index. From the overview, reach the same promise through its feature and through the system’s requirements index. Both routes explain why the promise is listed and lead to its canonical page or passage. The statement, conditions and criterion agree with the authored text. Editing a condition only in the generated spec or removing the feature link produces a check finding. A second build over identical inputs produces identical files.

The system overview opens with the system name and authored description. Its
section map links directly to populated indexes. Readers can follow goals,
user needs, scope, system topology, requirements, quality attributes, and
validation in that order when present. The topology shows declared immediate
parts and boundary relationships in a Mermaid diagram and in linked text.
The overview displays the full declared quality hierarchy with definition
previews. A shared attribute appears under each declared parent, with its
descendants expanded once. Business rule listings belong in the rules index.

A reader can look up a term from an inline or imported glossary and follow it
to its defining heading. Declared topics link to that same definition. When two
contexts define the same word, each entry names its context. A relationship
between terms links in both reading directions, with the original statement
and examples available in the declaring entry.

From a subject page’s navigation, one link opens the authored heading. The
reference definition carries the page’s source URI once. From an embedded
obligation or an authored passage reached by a report link, one link at the end
of the passage opens its own source heading. The project check identifies
duplicate page claims and missing or misdirected source links.

## System verification plan

Run the [To-do reader journey and regression probes](<../../scripts/check-example/check.py>). Inspect the artifact and diagnostics using the [generated specification guide](<../../docs/reference/generated-specification.md>).

[View source for “System verification plan”](<../../specs/requirements/generated-specification.requirement.spec.md#system-verification-plan>)

## Rationale

Reviewers can inspect product intent in Git and follow a promise to its context and evidence without learning the repository layout.

## Related specifications

**Quality attribute**

- [Repeatability](<../quality/reliability/repeatability/repeatability.quality-attribute.spec-report.md>)
- [Usability](<../quality/usability/usability.quality-attribute.spec-report.md>)

[subject-source]: <../../specs/requirements/generated-specification.requirement.spec.md#read-and-verify-the-generated-requirements-specification>
