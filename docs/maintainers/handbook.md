# Maintaining the handbook

Use this page when a change affects what dotspec teaches authors, what its
commands can establish, or how it interprets a source. The [handbook](../cli-handbook/README.md)
is the shared entry point for authors, agents and maintainers.

## Review a decision against dotspec's philosophy

| Principle | Question to ask in review |
| --- | --- |
| Optimize for human readability. | Can the author state the idea in ordinary prose? If a new field is proposed, which authoring or inspection task requires the structure? [Human-readable specifications](../../AGENTS.md#design-principles) |
| Keep each requirement with its authoritative source. | Does each requirement have one requirement file? Can a reader locate the owning use case, concept action, external provision or verification example? [Specification authority](../../CONTRIBUTING.md#specification-tests-and-implementation-tests) |
| Let authors develop precision as decisions are made. | Can someone record a useful sketch, a consequential exception or an unresolved threshold, then improve it when they have evidence? |
| Say exactly what an inspection establishes. | Could someone read a passing check or a resolved link as proof of stakeholder agreement or satisfaction? Identify the observation that would support the stronger claim. |

## Trace a choice to its sources

Record source-dependent guidance in the guide where the reader makes the
decision. Cite the edition and section, explain what dotspec takes from it,
and say where the local meaning or authoring form differs. Multiple sources
can illuminate the same decision. The following are starting points for review.

| Source | Contribution and dotspec choice | Where to audit the interpretation |
| --- | --- | --- |
| Craig's product-engineering practice | The practice of giving a requirement an authoritative home informs how dotspec links prose, contracts and tests. The consuming project decides which expression it accepts. The sibling research sketch also contains proposals. Check current behavior before carrying one into a guide. | In a sibling checkout, see `../me/projects/dotspec/dotspec-sketch.md`, “Relationship to the product-engineering bundle” and “Method lineage and connected views”. In `@craigsmitham/knowledge/product-engineering` v2.2.0 see `solution/requirements/foundations/one-authority-many-witnesses.md`. Compare with [dotspec's authority policy](../../CONTRIBUTING.md#specification-tests-and-implementation-tests). |
| Bertrand Meyer, *Handbook of Requirements and Business Analysis* | PEGS supplies reading perspectives over existing specifications. The violation question helps authors make a promise assessable. dotspec explicitly distinguishes verification and validation where Meyer sometimes uses a collective shorthand. | [Meyer (2022), §1.1.1 (PEGS), §4.12 (verifiability) and the terminology note on p. 20](https://doi.org/10.1007/978-3-031-06739-6) · [§§3.4–3.5, Goals and Environment books](https://se.inf.ethz.ch/requirements/plan.pdf) |
| Alistair Cockburn, *The Mini-Book on Use Cases*, v1.1b, and *Writing Effective Use Cases*, 1999 draft; Ivar Jacobson and Alistair Cockburn, *Use-Case Foundation*, v1.1 | Actor goals, boundaries, extensions and progressive precision inform the use-case guide. A sketch can grow as consequential alternatives become known. Cockburn keeps a specific system-in-use story apart from a generically named use case. dotspec lets a step stay specific when it links to the general promise it illustrates. Markdown headings and reference behavior are dotspec choices. | [Cockburn’s books](https://alistaircockburn.com/Books) · [Cockburn (1999 draft), “The System-in-use story,” p. 17](https://jpaulgibson.synology.me/~jpaulgibson/TSP/Teaching/Teaching-ReadingMaterial/Cockburn00.pdf#page=17) · [Jacobson and Cockburn, “Browse and shop,” p. 5](https://alistaircockburn.com/Use%20Case%20Foundation.pdf#page=5) |
| Daniel Jackson, concept design | Purpose, state, actions and operational principles explain a unit of behavior. dotspec distinguishes the concept definition from a named use in a system. | [Jackson’s concept criteria](https://essenceofsoftware.com/tutorials/concept-basics/criteria/) |
| ISO/IEC/IEEE 29148:2018 §§5.2.4–5.2.6 | Construction and quality criteria guide review of individual requirements and sets. Each requirement is stated in its own specification. | [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html) · [Specification authority](../../CONTRIBUTING.md#specification-tests-and-implementation-tests) |
| OMG Essence 1.2, §8.3.2.1 and Annex A.3.2.1 | The Requirements alpha and the development extension’s Requirement Item distinguish shared understanding, readiness for development and evidence of delivery. The handbook uses their questions to guide review. | [OMG Essence 1.2, Tables 8.4 and A.5](https://www.omg.org/spec/Essence/1.2/PDF) · [Review one promise](../cli-handbook/guides/write-an-assessable-requirement.md#check-the-promise-against-its-source) · [Review the set](../cli-handbook/guides/review-a-set-of-requirements.md) |
| ISO/IEC 25010:2023, 25030:2019 and 25020:2019 | The first supplies product-quality categories, the second informs quality requirements in context, and the third informs measurement. Selecting an attribute classifies a concern. An author still needs to establish what result is acceptable. | [Edition-specific starter](../reference/quality-models/iso-25010-2023.md) |

The [Jobs explainer](../cli-handbook/explanations/jobs-to-be-done.md#choose-a-useful-perspective) identifies the research behind each JTBD perspective. Keep book and standards citations
usable by title, edition and passage so contributors can use their own copies. When a new edition changes a term, follow its row to the
guides and examples, compare the old and new meanings, and record which edition
dotspec adopts. For a deliberate departure, give the reason beside the local
choice so it is visible in the same review.

## Follow a change through the repository

1. **Find the owning text.** Edit the linked authoring guide for advice and
   rationale. Edit the [metamodel](../../model/metamodel.yaml) for a type or
   relationship contract and the [mapping](../../model/markdown-mapping.md) for
   its Markdown representation. Commands and starters have their own
   [implementation sources](../../CONTRIBUTING.md).
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

From the checkout root, run `mise run build` and put the absolute path to its `dist` directory on PATH. Use a new project directory outside the checkout. From there, run `dotspec --version` and open the changed page with `dotspec docs PAGE`, using the page ID shown in the topic menu. Check for your revised wording before trying its commands. A binary built before the edit still contains the earlier handbook.

Give a reader this CLI and the prerequisites named by the guide. If you are reviewing alone, follow the directions on a fresh case and inspect what the next reader would see. Vary one decision that matters, such as which system owns a quality definition. Record where each case led and revise any instruction that required an unstated path or input. Ask them to find the To-do job via `dotspec docs job` and open its **Job** section offline. In a temporary project, ask them to create a job with a name they choose and inspect it with `dotspec job show`. From the needs-to-capabilities guide, ask where they would run `dotspec model needs` for their work and what they would need to reproduce the To-do result. Record where they guessed a project name, path or prerequisite.

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
`docs/cli-handbook/guides/`, an explanation in `docs/cli-handbook/explanations/`, and supporting material with its topic. Publish technical contracts and generated lookup material in `docs/reference/`. Use the H1 for the
page title. A CLI page ID is its path under `docs/cli-handbook/` without the `.md` suffix.
Worked examples keep their repository paths. New handbook Markdown is available by page ID and search at build time. Add it to a topic to give readers a route from the index. Keep contributor procedures and plans in `docs/maintainers/`.

Review command context separately from link resolution. For each instruction, including those in a linked walkthrough, ask “Where do I run this, and which files must already exist?” Try commands whose prerequisites changed in a temporary consumer project. A fenced specimen can show paths or partial commands as authored.

The [publication test](../../internal/integration/handbook_test.go) reads local links from published prose and opens their targets with the built CLI from an empty directory. For a supplemental file that genuinely requires the checkout, put the purpose in the sentence and end the link label with `(dotspec checkout)`. This marks just that citation. The test checks its repository file and any linked heading. Contributor procedures belong here or in [Contributing](../../CONTRIBUTING.md#specification-tests-and-implementation-tests).

In the [topic catalog](../cli-handbook/topics.json), list each topic's pages in
reading order. The first is the direct destination after creating a draft.
Include synonyms people use for the topic. On a page, map each diagnostic
rule ID to a heading fragment or an empty string for the listed destination.
A page can serve several topics. Map cross-type inspection and layout findings in [the technical diagnostic destinations](../reference/diagnostics.json). Run `mise run generate` to refresh the
[handbook topic index](../cli-handbook/README.md#what-am-i-trying-to-specify) and reference, then `mise run check` to check
the installed commands, destinations, anchors, and package.

To distribute the handbook for AXM, run `mise run knowledge-bundle` in the dotspec checkout. Extract `dist/dotspec-knowledge.zip` into an empty directory. With AXM installed, run `axm knowledge lint --path /absolute/path/to/extracted/dotspec`, substituting the directory that contains `knowledge.json`. A reachability finding identifies a page to link from the package index or another reachable page.

To try the package, change to a disposable AXM workspace configured for your agent. Run `axm knowledge install /absolute/path/to/extracted/dotspec --knowledge dotspec` with that same bundle path. Open `axm knowledge concepts get '@craigsmitham/knowledge/dotspec#docs/cli-handbook/guides/write-a-job'` and compare the published guidance with its source. Its index preserves each catalog topic and section link. Topics that share a page use one packaged copy. The bundle preserves fenced examples containing relative
links as downloadable text so link indexing treats their contents as sample
source. Release builds package the same sources with the release
version. Edit the linked handbook and authoring guides to change the knowledge itself.

## Changing a reference

Edit the technical contract in `docs/reference/` or the Markdown mapping in `model/` to clarify its explanation. The [guide to generating a specification](../reference/generated-specification.md) is hand-authored. The following references are generated.

| Published material | Edit here | Regenerate with |
| --- | --- | --- |
| CLI reference and output gallery | [Cobra commands](../../internal/cli) and the [gallery scenarios](../../scripts/generate/main.go) | `mise run generate` |
| Quality-model starter catalogs | [Bundled starter specifications](../../internal/qualitymodel/specs/quality-models) | `mise run generate` |
| Core domain, fields, and relationships | [Metamodel](../../model/metamodel.yaml) and [core view definition](../../model/views/core-domain.yaml) | `mise run model-generate` |
| [To-do requirements specification](../../examples/todo/SPEC.md) | [To-do specifications](../../examples/todo/specs) for the example, [document generator](../../internal/specdoc) for presentation | `mise run generate` |
