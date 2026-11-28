# Working on dotspec

Use the [dotspec handbook](docs/cli-handbook/README.md) to choose an authoring guide or worked example. For changes to domain meaning or authoring advice, consult [Maintaining the handbook](docs/maintainers/handbook.md) for the principles, sources, and review questions.

Breaking changes are acceptable. Do not introduce backward compatibility mechanisms.

<!-- axm:knowledge -->
## Knowledge Bundles

Use `axm knowledge concepts --help` to search, read, and explore these bundles.

### @craigsmitham

| Bundle | Description |
| --- | --- |
| [docs](agent_extensions/registry.agentxm.ai/@craigsmitham/knowledge/docs/src/index.md) | Explainers, authoring guides, and patterns for documentation across domains, including knowledge bases, product guides, and reusable checklists |
<!-- /axm:knowledge -->

<!-- dotspec:specification -->
For behavior changes and reviews, use the [system specification](<specification/dotspec.system.spec-report.md>) to find the relevant requirement. Follow its links to the governing source and any verification procedures.

Record new or revised behavior in the source that defines it. After editing specification sources, run `dotspec build` then `dotspec check` from the directory containing `dotspec.yaml`. Use `dotspec docs` for authoring guidance.
<!-- /dotspec:specification -->
