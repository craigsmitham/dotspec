# dotspec — Software product requirements framework and toolkit

dotspec helps teams develop and manage software product requirements throughout their lifecycle: elicitation, analysis, specification, validation, verification, and maintenance.

The goal is to produce high-quality requirements that are easy and inexpensive to validate, verify, and change.

A product's intent layer is everything that says what the software must do and why: the information that would have to survive if the implementation were discarded and rebuilt. Some of it is written as dotspec Markdown specifications. Some lives in other forms: formal models such as Alloy, Quint or TLA+, executable specifications and acceptance tests, schemas and interface contracts, and diagrams. dotspec catalogs that layer, records which source is authoritative for each decision, and generates a report that lets a reader understand and navigate it as one specification. A specification lists the artifacts that carry part of its intent as companions in its frontmatter.

A companion `dotspec` CLI helps teams manage and inspect these files, supporting the specification, validation, and verification of requirements.

Use the [dotspec handbook](docs/cli-handbook/README.md) to choose an authoring guide or worked example. For changes to domain meaning or authoring advice, consult [Maintaining the handbook](docs/maintainers/handbook.md) for the principles, sources, and review questions.

Breaking changes are acceptable. Do not introduce backward compatibility mechanisms.

# Design principles

## The intent layer

* **Catalog what must survive regeneration.** Ask of any artifact: if the implementation were thrown away and rebuilt, would this have to carry over? If so, it belongs in the product specification, whatever its format. The generated report should reach every authoritative specification in the project.
* **Make authority explicit.** Each decision has one authoritative source, which may be a dotspec specification, a formal model, an executable specification or a native contract. A link, a file's location or a matching name establishes a reference. Authority comes from a declaration that names the subject and the exact definition that governs it.
* **Classify by role, not by format.** The author states whether an artifact formalizes, illustrates, verifies or explains a specification. A Quint model can state an invariant and also check it. An acceptance script can be the governing example of a behavior or one witness among several. A Mermaid state diagram can define the allowed transitions or illustrate one path.
* **Use Markdown as the common entry point.** Every authoritative specification is reachable from a Markdown declaration that names its subject and states its claim in prose. A reader who cannot read Alloy still learns what is promised; a reader who can follows the link to the exact assertion. The native artifact stays in its own format and is checked with its own tools. dotspec records where it lives, how it is assessed, and what was observed.
* **Specify the acceptable implementations, at the precision each decision needs.** A specification is sufficient when every implementation that satisfies it would be acceptable. It need not describe the whole system. Let authors record a sketch, specify one component or state machine at a time, and add a formal model where a property such as an ordering or concurrency guarantee warrants it.
* **Say exactly what a check establishes.** A bounded Alloy check, a Quint simulation and a passing acceptance test each support a specific claim about a specific revision, scope and set of conditions. A verified model is evidence about the model; conformance of the code to that model needs evidence of its own. A resolved link establishes addressability, not agreement or satisfaction.

## Readers

* **Optimize Markdown specifications for human readability.** Use a diagram, table or quoted formal definition where it states something more clearly than prose. Prefer text-based diagrams such as Mermaid, which are reviewed and diffed like the rest of the source. Keep frontmatter to small, checked facts and put explanation in the body.
* **Write each generated page for the reader’s task.** Give an overview enough substance to explain the system, an index enough context to choose a destination, and a subject page enough detail to understand its subject. Keep each requirement’s statement, conditions, and criteria together, with the companions that formalize or verify it.
* **Describe the system in the project’s own terms.** A report section should tell readers something specific about the system. Under *Business rules*, for example, name an applicable rule and give its authored statement. Under *Companions*, name each artifact and what it covers. Definitions of dotspec categories and notations belong in the handbook. When the only information a section would add is where to browse a collection, put the link in the section map.
* **Give each detail a clear home and link to it in context.** Use a short, source-faithful description when introducing a related subject. Repeat a link wherever the relationship helps the reader, and point it directly to the section that carries the full detail. A detail's home may be a native artifact: link to the specific definition, not just the file.

<!-- axm:knowledge -->
## Knowledge Bundles

Use `axm knowledge concepts --help` to search, read, and explore these bundles.

### @craigsmitham

| Bundle | Description |
| --- | --- |
| [docs](agent_extensions/registry.agentxm.ai/@craigsmitham/knowledge/docs/src/index.md) | Explainers, authoring guides, and patterns for documentation across domains, including knowledge bases, product guides, and reusable checklists |
<!-- /axm:knowledge -->
