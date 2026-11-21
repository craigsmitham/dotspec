# Maintaining the handbook

Use this page when a change affects what dotspec teaches authors, what its
commands can establish, or how it interprets a source. The [handbook](../cli-handbook/README.md)
is the shared entry point for authors, agents and maintainers.

## Review a decision against dotspec's philosophy

| Principle | Question to ask in review |
| --- | --- |
| Optimize for human readability. | Can the author state the idea in ordinary prose? If a new field is proposed, which authoring or inspection task requires the structure? [Human-readable specifications](../../AGENTS.md#design-principles) |
| Keep each requirement with its authoritative source. | Does each requirement have one requirement file? Can a reader locate the owning use case, concept action, external provision or verification example? [Specification authority](../cli-handbook/explanations/requirements-layers.md#choose-the-authoritative-specification) |
| Let authors develop precision as decisions are made. | Can someone record a useful sketch, a consequential exception or an unresolved threshold, then improve it when they have evidence? [Use-case authoring](../cli-handbook/guides/use-cases.md#sketch-the-interaction-then-investigate-alternatives) |
| Say exactly what an inspection establishes. | Could someone read a passing check or a resolved link as proof of stakeholder agreement or satisfaction? Identify the observation that would support the stronger claim. [Assessment purposes](../cli-handbook/explanations/verification-and-validation.md#dotspec-commands-and-model-roles) |

## Trace a choice to its sources

Record source-dependent guidance in the guide where the reader makes the
decision. Cite the edition and section, explain what dotspec takes from it,
and say where the local meaning or authoring form differs. Multiple sources
can illuminate the same decision. The following are starting points for review.

| Source | Contribution and dotspec choice | Where to audit the interpretation |
| --- | --- | --- |
| Craig's product-engineering practice | The practice of giving a requirement an authoritative home informs how dotspec links prose, contracts and tests. The consuming project decides which expression it accepts. The sibling research sketch also contains proposals. Check current behavior before carrying one into a guide. | In a sibling checkout, see `../me/projects/dotspec/dotspec-sketch.md`, “Relationship to the product-engineering bundle” and “Method lineage and connected views”. In `@craigsmitham/knowledge/product-engineering` v2.2.0 see `solution/requirements/foundations/one-authority-many-witnesses.md`. Compare with [dotspec's authority policy](../../CONTRIBUTING.md#specification-tests-and-implementation-tests). |
| Bertrand Meyer, *Handbook of Requirements and Business Analysis* | PEGS supplies reading perspectives over existing specifications. The violation question helps authors make a promise assessable. dotspec explicitly distinguishes verification and validation where Meyer sometimes uses a collective shorthand. | [Review guidance and chapter references](../cli-handbook/guides/review.md#basis-and-limits) · [Requirement guidance](../cli-handbook/guides/requirements.md#basis) · [Terminology decision](../cli-handbook/explanations/verification-and-validation.md#other-established-usages) |
| Alistair Cockburn, *The Mini-Book on Use Cases*, v1.1b | Actor goals, boundaries, extensions and progressive precision inform the use-case guide. A sketch can grow as consequential alternatives become known. Markdown headings and reference behavior are dotspec choices. | [Use-case guide and source sections](../cli-handbook/guides/use-cases.md#review-and-check) |
| Daniel Jackson, concept design | Purpose, state, actions and operational principles explain a unit of behavior. dotspec distinguishes the concept definition from a named use in a system. | [Concept guide](../cli-handbook/guides/concepts.md) · [Concepts and use cases](../cli-handbook/explanations/requirements-layers.md#concepts-and-use-cases-are-complementary) |
| ISO/IEC/IEEE 29148:2018 §§5.2.4–5.2.6 | Construction and quality criteria guide review of individual requirements and sets. Each requirement is stated in its own specification. | [Requirement authoring basis](../cli-handbook/guides/requirements.md#basis) · [Specification authority](../cli-handbook/explanations/requirements-layers.md#choose-the-authoritative-specification) |
| ISO/IEC 25010:2023, 25030:2019 and 25020:2019 | The first supplies product-quality categories, the second informs quality requirements in context, and the third informs measurement. Selecting an attribute classifies a concern. An author still needs to establish what result is acceptable. | [Edition-specific starter](../cli-handbook/reference/quality-models/iso-25010-2023.md) · [Applicable clauses and quality-in-use distinction](../cli-handbook/guides/requirements.md#basis) |

The [requirements layers](../cli-handbook/explanations/requirements-layers.md#method-basis) and
[jobs guide](../cli-handbook/explanations/jobs-to-be-done.md#choose-a-useful-perspective) give the
basis for the business and user perspectives. Keep book and standards citations
usable by title, edition and passage so contributors can use their own copies. When a new edition changes a term, follow its row to the
guides and examples, compare the old and new meanings, and record which edition
dotspec adopts. For a deliberate departure, give the reason beside the local
choice so it is visible in the same review.

## Follow a change through the repository

1. **Find the owning text.** Edit the linked authoring guide for advice and
   rationale. Edit the [metamodel](../../model/metamodel.yaml) for a type or
   relationship contract and the [mapping](../../model/markdown-mapping.md) for
   its Markdown representation. Commands and starters have their own
   [implementation sources](../../CONTRIBUTING.md). The handbook indexes these
   sources so a reviewer can find the relevant diff.
2. **Review a case that exposes the distinction.** If a contributor proposes
   requiring a numeric target whenever a quality attribute is used, read the
   [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md).
   In a temporary copy of the To-do project, remove its Target subsection and run
   `dotspec requirement check restart`. Compare the advice with the
   proposed rule. Where would an author get the bound? How should a qualitative
   requirement or an unresolved draft be handled?
3. **Check what readers receive.** Update the index if a decision or route
   changed. Follow the link from the handbook to the example. For changes to
   commands, starters or model definitions, use the existing generation and
   drift checks in [Contributing](../../CONTRIBUTING.md). For a change of
   meaning, ask a reviewer to apply the revised guidance to the example and
   explain why the result fits the principles above.

IREB's glossary illustrates term-level attribution for adapted definitions.
Django's design philosophies show how a principle can guide future decisions,
and the Rust API Guidelines connect a compact review list to the reasoning for
each item. [IREB glossary, attribution appendix](https://hub.ireb.org/media/pages/resources/cpre-glossary/bb302c52fe-1787039949/ireb_cpre_glossary_en_2.2.pdf) · [Django design philosophies](https://docs.djangoproject.com/en/dev/misc/design-philosophies/) · [Rust review checklist](https://rust-lang.github.io/api-guidelines/checklist.html).

## Try the guidance from an installation

Give a reader the installed CLI and an empty directory outside the dotspec checkout. Ask them to find the To-do job via `dotspec docs job` and open its **Job** section offline. In a temporary project, ask them to create a job with a name they choose and inspect it with `dotspec job show`. From the capability guide, ask where they would run `dotspec model needs` for their work and what they would need to reproduce the To-do result. Record where they guessed a project name, path or prerequisite.

For the following exercises, start the reader in `examples/todo/` in a disposable copy of the dotspec checkout.

- Find the promise that governs a task's identity across restart. Which test would detect a change?
- Someone wants to pick up unfinished errands after a break. Identify the job and the interaction that supports it.
- Remove the restart criterion's Target in the copy. Use the CLI to find what the author still needs to decide.

Review the reader's changes. They should be able to distinguish the use case's
general promise from its restart refinement, and the evidence from the restart
test from observations of someone using the list. Record where they had to guess.

To review the Jobs topic, give a reader the draft “Resume an interrupted errand” and the occasion of a plumber asking for a photo after a call. Starting at `dotspec docs job`, can they identify who owns the commitment, explain where that responsibility ends, and name a criterion outside the return interaction? Ask what evidence would change their map. Then ask what question would warrant a switch interview and which capability might support one identified need.

The research basis includes the OneDrive library copies of Anthony Ulwick's *Jobs to Be Done: Theory to Practice* (2016, revised June 2023), chs. 2 and 4, Jim Kalbach's *The Jobs to Be Done Playbook* (2020), chs. 2–4, Bob Moesta's *Learning to Build*, ch. 3, and Richard Allen's *User Needs Mapping* (first edition, 2025), chs. 3–4. In the sibling checkout, the product-engineering bundle's `foundations/jobs-to-be-done.md` compares the research questions. See `../me/projects/trellis/widowed-household-job-map.md` for a map with an explicit boundary and open research questions, and `../me/projects/dotspec/dotspec-sketch.md`, “Jobs, job steps, and desired outcomes,” for the earlier commitments example. The published guides link the primary sources readers can access outside the checkout.

## Distribute the same guides

The CLI publishes the handbook landing page, topic catalog and Markdown beneath the documentation-form directories. Keep repository authoring instructions in the handbook root or maintainer documentation. Put a how-to in
`docs/cli-handbook/guides/`, an explanation in `docs/cli-handbook/explanations/`, a guided exercise in
`docs/cli-handbook/tutorials/`, and lookup material in `docs/cli-handbook/reference/`. Use the H1 for the
page title. A CLI page ID is its path under `docs/cli-handbook/` without the `.md` suffix.
Worked examples keep their repository paths. New handbook Markdown is available by page ID and search at build time. Add it to a topic to give readers a route from the index. Keep contributor procedures and plans in `docs/maintainers/`.

Review command context separately from link resolution. For each instruction, including those in a linked walkthrough, ask “Where do I run this, and which files must already exist?” Try commands whose prerequisites changed in a temporary consumer project. A fenced specimen can show paths or partial commands as authored.

The [publication test](../../internal/integration/handbook_test.go) reads local links from published prose and opens their targets with the built CLI from an empty directory. For a supplemental file that genuinely requires the checkout, put the purpose in the sentence and end the link label with `(dotspec checkout)`. This marks just that citation. The test checks its repository file and any linked heading. Contributor procedures belong here or in [Contributing](../../CONTRIBUTING.md#specification-tests-and-implementation-tests).

In the [topic catalog](../cli-handbook/topics.json), list each topic's pages in
reading order. The first is the direct destination after creating a draft.
Include synonyms people use for the topic. On a page, map each diagnostic
rule ID to a heading fragment or an empty string for the listed destination.
A page can serve several topics. Run `mise run generate` to refresh the
[handbook topic index](../cli-handbook/README.md#browse-by-topic) and reference, then `mise run check` to check
the installed commands, destinations, anchors, and package.

For AXM, `mise run knowledge-bundle` creates `dist/dotspec-knowledge.zip`.
Extract it and use `axm knowledge lint --path <extracted>/dotspec` to inspect
the package. In a consuming workspace, `axm knowledge install <extracted>/dotspec --knowledge dotspec`
installs it on request. Its index preserves each catalog topic and section link. Topics that share a page use one packaged copy. The bundle preserves fenced examples containing relative
links as downloadable text so link indexing treats their contents as sample
source. Release builds package the same sources with the release
version. Edit the linked handbook and authoring guides to change the knowledge itself.

## Changing a reference

Edit the linked format, layout, workspace, inspection, or mapping page directly to clarify its explanation. The following references are generated.

| Published material | Edit here | Regenerate with |
| --- | --- | --- |
| CLI reference and output gallery | [Cobra commands](../../internal/cli) and the [gallery scenarios](../../scripts/generate/main.go) | `mise run generate` |
| Quality-model starter catalogs | [Bundled starter specifications](../../internal/qualitymodel/specs/quality-models) | `mise run generate` |
| Core domain, fields, and relationships | [Metamodel](../../model/metamodel.yaml) and [core view definition](../../model/views/core-domain.yaml) | `mise run model-generate` |
| [To-do requirements specification](../../examples/todo/SPEC.md) and [generation guide](../cli-handbook/reference/generated-specification.md) | [To-do specifications](../../examples/todo/specs) for the example, [document generator](../../internal/specdoc) for presentation | `mise run generate` |
