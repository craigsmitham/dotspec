# Maintaining the handbook

Use this page when a change affects what dotspec teaches authors, what its
commands can establish, or how it interprets a source. The [handbook](../handbook/README.md)
is the shared entry point for authors, agents and maintainers.

## Review a decision against dotspec's philosophy

| Principle | Question to ask in review |
| --- | --- |
| Optimize for human readability. | Can the author state the idea in ordinary prose? If a new field is proposed, which authoring or inspection task requires the structure? [Human-readable specifications](../../specs/dotspec.system.spec.md#human-readable-specifications) |
| Keep each obligation with its authoritative source. | If a test, contract or use-case clause already owns the promise, can another document cite it? Who decides which source carries authority here? [Obligations and documents](../handbook/requirements-layers.md#distinguish-an-obligation-from-its-document) |
| Let authors develop precision as decisions are made. | Can someone record a useful sketch, a consequential exception or an unresolved threshold, then improve it when they have evidence? [Use-case authoring](../handbook/use-cases.md#sketch-the-interaction-then-investigate-alternatives) |
| Say exactly what an inspection establishes. | Could someone read a passing check or a resolved link as proof of stakeholder agreement or satisfaction? Identify the observation that would support the stronger claim. [Assessment purposes](../handbook/verification-and-validation.md#dotspec-commands-and-model-roles) |

## Trace a choice to its sources

Record source-dependent guidance in the guide where the reader makes the
decision. Cite the edition and section, explain what dotspec takes from it,
and say where the local meaning or authoring form differs. Multiple sources
can illuminate the same decision. The following are starting points for review.

| Source | Contribution and dotspec choice | Where to audit the interpretation |
| --- | --- | --- |
| Craig's product-engineering practice | The practice of giving an obligation an authoritative home informs how dotspec links prose, contracts and tests. The consuming project decides which expression it accepts. The sibling research sketch also contains proposals. Check current behavior before carrying one into a guide. | In a sibling checkout, see `../me/projects/dotspec/dotspec-sketch.md`, “Relationship to the product-engineering bundle” and “Method lineage and connected views”. In `@craigsmitham/knowledge/product-engineering` v2.2.0 see `solution/requirements/foundations/one-authority-many-witnesses.md`. Compare with [dotspec's authority policy](../../CONTRIBUTING.md#specification-tests-and-implementation-tests). |
| Bertrand Meyer, *Handbook of Requirements and Business Analysis* | PEGS supplies reading perspectives over existing specifications. The violation question helps authors make a promise assessable. dotspec explicitly distinguishes verification and validation where Meyer sometimes uses a collective shorthand. | [Review guidance and chapter references](../handbook/review.md#basis-and-limits) · [Requirement guidance](../handbook/requirements.md#basis) · [Terminology decision](../handbook/verification-and-validation.md#other-established-usages) |
| Alistair Cockburn, *The Mini-Book on Use Cases*, v1.1b | Actor goals, boundaries, extensions and progressive precision inform the use-case guide. A sketch can grow as consequential alternatives become known. Markdown headings and reference behavior are dotspec choices. | [Use-case guide and source sections](../handbook/use-cases.md#review-and-check) |
| Daniel Jackson, concept design | Purpose, state, actions and operational principles explain a unit of behavior. dotspec distinguishes the concept definition from a named use in a system. | [Concept guide](../handbook/concepts.md) · [Concepts and use cases](../handbook/requirements-layers.md#concepts-and-use-cases-are-complementary) |
| ISO/IEC/IEEE 29148:2018 §§5.2.4–5.2.6 | Construction and quality criteria guide review of individual requirements and sets. dotspec uses *obligation* for the promise and *Requirement* for the specification that gives one promise its own file. | [Requirement authoring basis](../handbook/requirements.md#basis) · [Obligation and document](../handbook/requirements-layers.md#distinguish-an-obligation-from-its-document) |
| ISO/IEC 25010:2023, 25030:2019 and 25020:2019 | The first supplies product-quality categories, the second informs quality requirements in context, and the third informs measurement. Selecting an attribute classifies a concern. An author still needs to establish what result is acceptable. | [Edition-specific starter](../reference/quality-models/iso-25010-2023.md) · [Applicable clauses and quality-in-use distinction](../handbook/requirements.md#basis) |

The [requirements layers](../handbook/requirements-layers.md#method-basis) and
[jobs guide](../handbook/jobs.md#choose-a-useful-perspective) give the
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
   obligation or an unresolved draft be handled?
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

Give a reader a copy of `examples/todo/` and the installed CLI outside the
checkout. Let them discover the guides from `--help` and try these prompts:

- Find the promise that governs a task's identity across restart. Which test would detect a change?
- Someone wants to pick up unfinished errands after a break. Identify the job and the interaction that supports it.
- Remove the restart criterion's Target in the copy. Use the CLI to find what the author still needs to decide.

Review the reader's changes. They should be able to distinguish the use case's
general promise from its restart refinement, and the evidence from the restart
test from observations of someone using the list. Record where they had to guess.

## Distribute the same guides

The CLI embeds the reviewed Markdown at build time. The [topic catalog](../handbook/topics.json) owns topic names,
source pages and sections, and topic descriptions for both distributions. It also maps CLI diagnostic IDs to topics. When changing a heading used by a topic,
update the catalog, run `mise run generate` to refresh the [documentation map](../README.md), and run the installed-CLI discovery tests in
`internal/integration/testdata/guide.txtar`.

For AXM, `mise run knowledge-bundle` creates `dist/dotspec-knowledge.zip`.
Extract it and use `axm knowledge lint --path <extracted>/dotspec` to inspect
the package. In a consuming workspace, `axm knowledge install <extracted>/dotspec --knowledge dotspec`
installs it on request. Its index preserves each catalog topic and section link. Topics that share a page use one packaged copy. The bundle preserves fenced examples containing relative
links as downloadable text so link indexing treats their contents as sample
source. Release builds package the same sources with the release
version. Edit the linked handbook and authoring guides to change the knowledge itself.
