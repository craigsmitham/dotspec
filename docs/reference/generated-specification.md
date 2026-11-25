# Generate a requirements specification

Use a generated spec when readers need one entry point for requirements authored across a project. Open the project’s `SPEC.md` and follow its link to read the system’s purpose and scope, choose a feature, and follow its promises to the source and supporting evidence. The To-do example includes a [complete generated spec](../../examples/todo/SPEC.md).

For a review, start with one stakeholder goal and use the generated links to follow the promises across their owners. [Review a set of requirements](../cli-handbook/guides/review-a-set-of-requirements.md).

## Set up a project

Work in the directory that owns your project's `dotspec.yaml`. If you are starting a project, create that file with your project name. The To-do example uses the following manifest. Generated reports are enabled by default.

```yaml
name: todo
```

For a source-only project, add this section to the manifest. The project check still validates authored files.

```yaml
report:
  enabled: false
```

For one root system, the project landing page opens its specification. For several independent systems, the landing page links to a systems index. Use **Part of** when a system really is a subsystem. A project containing only reusable quality models publishes their reference catalog. [Describe a system](../cli-handbook/guides/describe-a-system.md) · [Separate specification projects](workspaces.md).

Explicit specification-system selections and directory defaults place subjects in their system. Standalone quality models are reusable reference material. Each `*.quality-model.md` source has a `*.quality-model-report.md` reference page. Links into other projects retain their original targets. A locally authored specification that selects another project’s system is reachable from the project landing page under **Specifications for other systems**.

From the project root, build the reports and check them against the authored specifications.

```sh
dotspec spec build
dotspec check
```

Commit the authored changes, `SPEC.md` and the generated Markdown in `specification/` together. From `SPEC.md`, open the overview. If the project has requirements, follow one to its statement and conditions. If the wording needs to change, select **View source**, edit that file, rebuild, and check again. Build refreshes generator-owned files and removes obsolete generated pages. If it reports that an authored file occupies a planned destination, move the authored file to a permanent location, update any links to it, then rerun the build.

## Read the generated files

| Destination | Contents |
| --- | --- |
| `SPEC.md` | The system overview, or a project index for multiple systems, external-system specifications, or a reference catalog |
| `specification/<name>.system.spec-report.md` | Root system overview and links into its specification |
| `specification/glossary.md` | Alphabetical definitions linked to their sources, a topic index when terms declare topics, context names for distinct meanings, and links back to authored relationships |
| `specification/use-cases.md` | Use cases with their authored goals and links to the scenarios |
| `specification/quality/<attribute>/<attribute>.quality-attribute.spec-report.md` | One system attribute, its definition, parent and child attributes, and classified requirements |
| `specification/requirements.md` | Requirements in this system, each linked to its canonical page or passage |
| `specification/systems/<child>/<child>.system.spec-report.md` | Child system overview with its own local directories |
| `specification/open-questions.md` | Questions recorded across the project, each with its subject and a link to the passage. The index appears when at least one authored question exists. |

### Navigate an overview

System and component overviews open with the subject name, authored description, and a linked section map. Component overviews also show their responsibilities, contracts, architecture, and implementation where authored.
Expand **On this page** to jump to a section in reading order. Expand **Browse specifications** to find a collection by name. For example, choose **Requirements** to review individual promises. To edit the overview’s authored content, follow **View source** at the foot of the report.

| Section | What you can find |
| --- | --- |
| Framing | How the authors understand the situation and what they are still exploring |
| Purpose | The system’s authored purpose |
| Goals | Authored goals. Select a named goal to read its full statement and see which specifications support it. |
| Business rules | **Business rules** opens the index of applicable and local rules. |
| Stakeholders | Named stakeholder roles in source order, each with its description and a heading that jobs and use cases can link to |
| Scope | The work, people and situations the system serves, followed by agreed exclusions under **Out of scope** |
| Environment | Conditions in the surrounding world that help explain the requirements |
| Assumptions | Named premises in declaration order. Follow a premise to its full statement, basis, applicability and the requirements that rely on it. |
| Constraints | Externally imposed limits on the solutions available for this system |
| System topology | Immediate parts, boundary partners, and required or provided interfaces, with a Mermaid diagram and linked descriptions |
| Features | Product features and their authored descriptions |
| Capabilities | Capabilities the system provides. When topology identifies a provider, its entry describes the capability. |
| Use cases | The use-case index, where each interaction is listed with its authored goal |
| Concepts | Links to concept definitions, uses and synchronizations. The body introduces named uses and synchronizations under their own headings. |
| Requirements | **Requirements** opens the index, including routes by component. |
| Quality Attributes | Every attribute declared for this system, with a short definition excerpt and its subattributes. Follow an attribute to read its full definition and the requirements it classifies. |
| Open questions | Decisions that remain open in the system specification |
| Glossary | Terms defined for this system and links to their definitions |
| Project references | Authored explanations and links to delivery plans, responsibilities, risks and decisions, when supplied |

**On this page** lists **Scope** once when the system describes its scope or records exclusions. If only exclusions are supplied, **Scope** opens directly with the **Out of scope** subsection. Links to exclusions go directly to that subsection. **Environment**, **Assumptions** and **Constraints** each appear as a peer of **Scope** when populated. Local Markdown subheadings stay with their passage.

Overview sections list all entries with project-specific detail in index order.
A collection with only navigation stays in the section map. The topology
and quality hierarchy show their declared structure in full.

Across overviews and indexes, each entry uses the subject’s linked name and authored wording. Indexes shorten long previews and mark an excerpt with an ellipsis. **Uses** beside a concept use leads to its concept definition.
**Refines** beside a requirement leads to each promise it refines, whether a requirement or business rule. Open the requirement to read its full
statement, conditions, and criteria together.

When the project has authored questions, **Open questions** opens their index from the root system’s section map. Run `dotspec model review` for advisory findings. A topology edge represents a declared relationship.
A child system has its own overview showing the next level of detail.

### Browse subject reports

System and component overviews keep the source basename and change `.spec.md`
to `.spec-report.md`. For example, the Task service report lives at
`components/task-service/task-service.component.spec-report.md`. Its directory also
holds its local requirements, interfaces, and child components. The optional source
`id` provides an anchor. The subject name determines the overview filename.

System overviews introduce named goals and assumptions with their title and first authored paragraph. Follow a title to read its qualifications and supporting references in `goals.md` or `assumptions.md`. Constraints appear in the owning system’s scope. For example, follow **Same local storage** to see which requirement relies on that assumption.
The requirement links back to the same entry. A **Referenced by** link returns
to an authored passage that cites the entry in prose. [Task-list assumptions](../../examples/todo/specification/assumptions.md)

Use **On this page → Quality Attributes** to see the system’s hierarchy. The breadcrumb on an attribute report opens the same hierarchy in `quality-attributes.md`. To browse reusable models, follow **Quality models** beside **View source** at the foot of the system report. This link is available even if the system has yet to declare an attribute. Shared attributes appear under each declared parent. Their own subattributes expand under the first parent. Select an attribute to read its full definition. Attribute reports nest beneath their first declared parent in the system’s `quality/` directory. Requirement reports live under their owner’s `requirements/` folder. Each attribute page links to requirements that name it. On a requirement page, read the statement, applicability and criterion together. Links beside the requirement lead back to the precise attribute or use-case passage.

For example, [Pick up where I left off](../../examples/todo/specification/features/pick-up-where-i-left-off.feature.spec-report.md) opens with its authored description of the behavior it offers. Follow the interaction to see how the person resumes work, or follow the restart requirement to find what must survive and how the acceptance test checks it.

A system, component, feature, capability, or concept use can supply a short
Description. System and component overviews use the description as their lead. Other subject reports open with the Description section. The index entry reuses
an excerpt from the first paragraph. An entry with no description shows
the subject name alone. Interfaces use their introductory description and quality
models use their Description section. For requirements and business rules, a
short Statement introduces the promise or rule. Use cases use their Goal and
jobs use their Job statement. Concept definitions use their Purpose.

The overview’s **Browse specifications** list links to collections such as use cases and components. For example,
open `use-cases.md` to choose an interaction by its goal, or `components.md` to
choose a component by its description. A system’s `systems.md` lists its immediate
child systems. When a component has entries in a collection, **By component**
on the parent index leads to the component’s index of the same name.

Above the title, breadcrumbs lead to the owning system or component and the subject’s collection. Under **Extensions**, **Actions**, and a quality model’s **Quality attributes**, select a named entry to jump to its section when several are present. Longer lists can be expanded.

On subject reports, collections that only provide navigation appear with the page’s index links. Sections with authored descriptions keep that text in the body. Within a component, `interfaces.md` groups required and provided interfaces.
A provider’s `capabilities.md` groups explicitly provided capabilities, including
those specified elsewhere. Its `requirements.md` lists local requirements, plus
requirements owned elsewhere that explicitly apply to the component. The
**Applicable requirements** group collects explicit applicability links, including
local requirements. Each entry links to the canonical page or passage.

Each standalone authored specification in the selected project has exactly one standalone report page. The owning system or component links to indexes of its features and capabilities. Inline subjects can also have a report page for their own heading.

Authored Markdown uses `<subject>-spec`. The corresponding generated Markdown
uses `<subject>-spec-report`. For example, `<name>.system.spec.md` declares
`type: system-spec`, and its generated `<name>.system.spec-report.md` declares `type: system-spec-report`.

Each subject report records its source in a YAML `source` field using a URI
relative to the report, including the fragment for the exact authored heading.
For example, the todo system report declares
`source: "../specs/personal-task-list.system.spec.md#personal-task-list"`.
The visible **View source** link opens this same heading. For a subject defined inside
another spec, the fragment identifies its nested heading.

A subject rendered as a section carries its type and exact authored heading in a
`dotspec-doc` comment. **View section source** opens the authored declaration. Links to supporting specifications and external evidence
stay beside the passages that use them. On the glossary index, **Term source** opens the defining passage. Plain topic lists and a trailing model reference share that line. The topic index jumps to term headings, and incoming relationship links return to the passage that describes the connection.

Each system quality attribute has a `system-quality-attribute-spec-report` page. Reusable quality model entries remain sections of the model’s report. The file’s `generated_by: dotspec` marker identifies generator
ownership. See the [document model](model/index.md#generateddocument).

## Check a workspace

Use `dotspec spec build --project NAME` to refresh a selected project. Run `dotspec check --all` on a clean checkout in CI, with a pinned dotspec version. Whole-project checks include the generated files for each project unless its manifest sets `report.enabled` to `false`. A missing, obsolete or stale file is a finding. The check also inspects the visible content and navigation. A subtree or individual-subject check covers only its selected source scope.

For early feedback, add `dotspec check --all` to a pre-push hook. Use CI to verify the committed tree. A pre-commit hook that enforces synchronization must check staged contents when contributors stage only part of a change.

## Record verification results

Link each test or manual protocol under **Verification procedures** in its requirement or criterion. After a run, an evidence producer can supply `verification-results.json` in the project root. Running `dotspec spec build` reads this file and publishes observations alongside the requirement. Running `dotspec check` checks the resulting report. The requirement file changes only when the authored decision changes.

This illustrative file requests the results of a restart assessment for revision `build-42`. Replace the revision, conditions, and observation address with the ones from your run.

```json
{
  "selection": {
    "product_revision": "build-42",
    "conditions": "local database, acknowledged writes"
  },
  "results": [
    {
      "subject": "specs/requirements/restart.requirement.spec.md#criterion",
      "procedure": "tests/acceptance.py",
      "product_revision": "build-42",
      "conditions": "local database, acknowledged writes",
      "outcome": "fail",
      "observation": "evidence/restart-42.log"
    }
  ]
}
```

Identify the subject by its project-relative requirement URI. Append `#criterion` when the procedure belongs to the criterion. Relative procedure and observation URIs start at the project root. Use `pass`, `fail`, or `inconclusive` for the outcome. A run matches when its procedure is linked for that subject and its product revision and conditions equal the selection. Either record can also carry `specification_revision`. Supplying it in the selection requires an exact match. An unknown product revision or unknown conditions leave an observation outside the matching set. Use an empty results array for a selection with no observations yet. A report preserves competing results and their individual observation links.
