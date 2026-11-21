# dotspec — Software product requirements framework and toolkit

dotspec helps teams develop and manage software product requirements throughout their lifecycle: elicitation, analysis, specification, validation, verification, and maintenance.

The goal is to produce high-quality requirements that are easy and inexpensive to validate, verify, and change.




A companion `dotspec` CLI helps teams manage and inspect these files, supporting the specification, validation, and verification of requirements.

Use the [dotspec handbook](docs/cli-handbook/README.md) to choose an authoring guide or worked example. For changes to domain meaning or authoring advice, consult [Maintaining the handbook](docs/maintainers/handbook.md) for the principles, sources, and review questions.


Breaking changes are acceptable. Do not introduce backward compatbility mechanisms.


# Design principles

* Optimize *.spec.md files for human readability
* **Write each generated page for the reader’s task.** Give an overview enough substance to explain the system, an index enough context to choose a destination, and a subject page enough detail to understand its subject. Keep each requirement’s statement, conditions, and criteria together.
* **Describe the system in the project’s own terms.** A report section should tell readers something specific about the system. Under *Business rules*, for example, name an applicable rule and give its authored statement. Definitions of dotspec categories belong in the handbook. When the only information a section would add is where to browse a collection, put the link in the section map.
* **Give each detail a clear home and link to it in context.** Use a short, source-faithful description when introducing a related subject. Repeat a link wherever the relationship helps the reader, and point it directly to the section that carries the full detail.

<!-- axm:knowledge -->
## Knowledge Bundles

Use `axm knowledge concepts --help` to search, read, and explore these bundles.

### @craigsmitham

| Bundle | Description |
| --- | --- |
| [docs](agent_extensions/registry.agentxm.ai/@craigsmitham/knowledge/docs/src/index.md) | Documentation craft and patterns for authoring, reviewing, organizing, and maintaining documents, reusable checklists, and pattern collections across domains |
<!-- /axm:knowledge -->
