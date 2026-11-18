# dotspec — Software product requirements framework and toolkit

dotspec helps teams develop and manage software product requirements throughout their lifecycle: elicitation, analysis, specification, validation, verification, and maintenance.

The goal is to produce high-quality requirements that are easy and inexpensive to validate, verify, and change.

Requirements live primarily in Markdown files named `*.spec.md`, co-located with the relevant parts of the software system in the repository. Together, these files form the product’s specification layer, capturing intent across business, user, and quality requirements.

A companion `dotspec` CLI helps teams manage and inspect these files, supporting the specification, validation, and verification of requirements.

Use the [dotspec handbook](docs/handbook/README.md) to choose an authoring guide or worked example. For changes to domain meaning or authoring advice, consult [Maintaining the handbook](docs/maintainers/handbook.md) for the principles, sources, and review questions.


Breaking changes are acceptable. Do not introduce backward compatbility mechanisms.


# Design principles

* Optimize *.spec.md files for human readability
