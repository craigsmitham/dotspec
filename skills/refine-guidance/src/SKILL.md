---
name: refine-guidance
description: "Compare dotspec example specs with worked examples from sources cited by the CLI handbook, then refine the guidance, specs, and example app behavior where the evidence warrants a change."
---

# Refine dotspec guidance

Use the requested topic, app, or spec path to bound the review. With no scope given, cover the authored specs across the apps under `examples/`, including component specs alongside implementation. Honor a request for findings only.

## Follow an authoring decision

Start with an authored `*.spec.md` file. Find the relevant explanation or guide through `docs/cli-handbook/README.md`. Frame the comparison as a decision an author faces. For example, which alternative paths in the To-do use case change what the system promises at the end? Follow the spec into the app and its acceptance tests when behavior helps answer the question.

Follow the handbook citations and research a worked example online from the cited author, publication, or standards body. Compare how the two examples handle the same decision, even across different domains or formats. Cite the URL, edition where applicable, and the section, page, or named example inspected. If only a definition or principle is accessible, record what it establishes and leave the worked-example comparison open.

## Decide where a change belongs

Apply the design principles in `AGENTS.md`, the adaptation questions in `docs/maintainers/handbook.md`, and the specification-authority guidance in `CONTRIBUTING.md`. Trace each local promise to its authoritative home. Use the external example to assess how a decision is expressed. Use local product evidence to establish the intended behavior. Record the reason for retaining a deliberate dotspec adaptation. If a product decision remains open, identify the evidence needed to settle it.

- **Handbook.** Where a reader would have to guess, revise the guide or explanation where they make the decision. Cite the relevant passage and explain what dotspec adopts or adapts.
- **Example spec.** Where an established decision is hard to understand or verify, improve its authored statement, conditions, criteria, or relationships. Regenerate affected reports from their inputs.
- **Example app.** Where behavior contradicts an established promise, reproduce the discrepancy when practical, fix the implementation, and update acceptance coverage if existing checks miss it.

## Verify and report

Apply revised advice to the case that exposed the gap. For a handbook edit, rebuild the CLI and open the published page with `dotspec docs` as described in the maintainer guide. For an app fix, exercise the same scenario through its normal entry point. Run the relevant example tests and repository checks.

Summarize the apps, specs, and topics reviewed. Connect each material finding to the local evidence, cited source, change or reason for retaining the current approach, and verification result. Identify unresolved product decisions and gaps in source evidence.
