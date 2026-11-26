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

Explicit specification-system selections and directory defaults place subjects in their system. Standalone quality models are reusable reference material. Each `*.quality-model.md` source has a `*.quality-model-report.md` reference page. Links into other projects retain their original targets. A locally authored specification that selects another project’s system is reachable from the project landing page through the index for its type, such as **Requirements for other systems**.

From the project root, build the reports and check them against the authored specifications.

```sh
dotspec spec build
dotspec check
```

Commit the authored changes, `SPEC.md` and the generated Markdown in `specification/` together. From `SPEC.md`, open the system report. If the project has requirements, follow one to its statement and conditions. If the wording needs to change, select **View source**, edit that file, rebuild, and check again. Build refreshes generator-owned files and removes obsolete generated pages. If it reports that an authored file occupies a planned destination, move the authored file to a permanent location, update any links to it, then rerun the build.

## Read the generated files

| Destination | Contents |
| --- | --- |
| `SPEC.md` | A link to the system report, or a project index for multiple systems, external-system specifications, or a reference catalog |
| `specification/<name>.system.spec-report.md` | Root system report and links into its specification |
| `specification/glossary.md` | Alphabetical definitions linked to their sources, a topic index when terms declare topics, context names for distinct meanings, and links back to authored relationships |
| `specification/quality/<attribute>/<attribute>.quality-attribute.spec-report.md` | One system attribute, its definition, parent and child attributes, and classified requirements |
| `specification/systems/<child>/<child>.system.spec-report.md` | Child system report with its own local directories |
| `specification/companions.md` | Every companion listed in the system, grouped by role, each with a link to the subject that declares it. The index appears under **Browse specifications** when at least one specification lists a companion. |
| `specification/open-questions.md` | Questions recorded across the project, each with its subject and a link to the passage. The index appears when at least one authored question exists. |

### Navigate a system or component report

System and component reports open with the subject name, authored description, and a linked section map. Component reports also show their purpose, overview, responsibilities, contracts, architecture, implementation, and companions where authored.
Expand **On this page** to jump to a section in reading order. On a system report, **Intent**, **Context**, **Behavior** and **Design** bring related links together. Choose **Quality Attributes** or **Requirements** directly between Behavior and Design. Only populated sections appear. **Browse specifications** appears when a collection has a separate index. To edit the report’s authored content, follow **View source** at the foot of the report.

| Section | What you can find |
| --- | --- |
| Purpose | The system’s authored purpose |
| Overview | The authors’ explanation of what the system is and how its parts fit together |
| Goals | Authored goals. Select a named goal to read its full statement and see which specifications support it. |
| Stakeholders | Named stakeholder roles in source order, each with its description and a heading that jobs and use cases can link to |
| Jobs | Linked jobs with their authored job statements |
| Scope | The work, people and situations the system serves, followed by agreed exclusions under **Out of scope** |
| Environment | Conditions in the surrounding world that help explain the requirements |
| Assumptions | Named premises in declaration order. Follow a premise to its full statement, basis, applicability and the requirements that rely on it. |
| Constraints | Externally imposed limits on the solutions available for this system |
| Business rules | Linked local and applicable rules with statement excerpts |
| Features | Product features and their authored descriptions |
| Capabilities | Abilities declared or provided in scope. Open a capability to see who provides it. |
| Use cases | Linked interactions with their authored goals |
| Quality Attributes | Every attribute declared for this system, with a short definition excerpt. Follow an attribute to see its parent and child attributes and the requirements it classifies. |
| Requirements | Requirements in scope with excerpts from their statements |
| Concepts | Local definitions and definitions referenced by scoped uses, with their authored purpose. Each named use appears beneath its definition, after *Uses:*. Open a use for its configuration and the specifications that apply to it. |
| Synchronizations | Authored coordination between concept uses. |
| Architecture | The system boundary and its immediate subsystems and components, the interfaces it offers or uses, and its declared interactions |
| Companions | Artifacts the system specification itself lists, such as a formal model or an acceptance test, grouped by role with each author’s note |
| Open questions | Decisions that remain open in the system specification |
| Glossary | Terms defined for this system and links to their definitions |
| Project references | Authored explanations and links to delivery plans, responsibilities, risks and decisions, when supplied |

**On this page** lists **Scope** once when the system describes its scope or records exclusions. If only exclusions are supplied, **Scope** opens directly with the **Out of scope** subsection. Links to exclusions go directly to that subsection. **Environment**, **Assumptions** and **Constraints** each appear as a peer of **Scope** when populated. Local Markdown subheadings stay with their passage.

Each specification directory lists titles alphabetically within one type. Select a title to open its detail report. A short excerpt from the authored text helps you choose where to go. An ellipsis marks shortened wording. A specification without a preview appears by title alone.

The system directory includes specifications declared in the system and its nested components. A component directory covers that component and its descendants. A separately specified subsystem has its own directories. Explicit applicability and capability providers can bring a specification into scope. A specification reached several ways appears once. Two specifications with the same title keep their own entries and destinations.

The **Architecture** diagram shows immediate parts inside the boundary and declared external interaction partners outside it. When both are present, separate structure and context diagrams show decomposition and interactions. Select a component to read its responsibilities or an interface to see which subject provides or requires it. On a component report, an authored **Architecture** passage leads the view.

A component lists its capabilities under **Architecture**. Open a capability for its declared providers. Under **Concepts**, select a definition for its state and actions, or a use nested beneath it for how this component applies it. A use awaiting a definition appears in the list by its own title, marked *(no definition)*. A use can also link to a definition in another project.

For example, the To-do system’s **Requirements** section includes **Acknowledge after commit**, declared on Task service. Following that title from either the system or the component’s requirements opens the same report. Open **Preserve tasks across restart** to read its full statement, conditions and criterion. Its **Refines** section leads to **Keep recorded work available**. [To-do system report](../../examples/todo/specification/personal-task-list.system.spec-report.md#requirements)

When the project has authored questions, **Open questions** opens their index from the root system’s section map. Run `dotspec model review` for advisory findings.
A child system has its own report showing the next level of detail.

### Browse subject reports

System and component reports keep the source basename and change `.spec.md`
to `.spec-report.md`. For example, the Task service report lives at
`components/task-service/task-service.component.spec-report.md`. Its directory also
holds its local requirements, interfaces, and child components. The optional source
`id` provides an anchor. The subject name determines the report filename.

System reports introduce named goals and assumptions with their title and first authored paragraph. Follow a title to read its qualifications and supporting references in `goals.md` or `assumptions.md`. For example, follow **Same local storage** to see which requirement relies on that assumption.
The requirement links back to the same entry. A **Referenced by** link returns
to an authored passage that cites the entry in prose. [Task-list assumptions](../../examples/todo/specification/assumptions.md)

Use **On this page → Quality Attributes** to browse every system quality attribute. The **Quality attributes** breadcrumb on an attribute report returns to this directory. Follow **Quality models** beside **View source** to browse reusable models. This link is available when the system has a model catalog, even if it has yet to declare an attribute. Select an attribute to see its definition and relationships. Attribute reports nest beneath their first declared parent in the system’s `quality/` directory. Requirement reports live under their owner’s `requirements/` folder. Each attribute page links to requirements that name it. On a requirement page, read the statement, applicability and criterion together. Links beside the requirement lead back to the precise attribute or use-case passage.

For example, [Pick up where I left off](../../examples/todo/specification/features/pick-up-where-i-left-off.feature.spec-report.md) opens with its authored description of the behavior it offers. Follow the interaction to see how the person resumes work, or follow the restart requirement to find what must survive and how the acceptance test checks it.

A system, component, feature, capability, or concept use can supply a short
Description. System and component reports use the description as their lead. Other subject reports open with the Description section. A feature report follows it with Purpose and then Overview when they are authored. The index entry reuses
an excerpt from the first paragraph. An entry with no description shows
the subject name alone. Interfaces use their introductory description and quality
models use their Description section. For requirements and business rules, a
short Statement introduces the promise or rule. Use cases use their Goal and
jobs use their Job statement. Concept definitions use their Purpose.

On a system report, choose a feature under **Features** or a component under **Architecture → Components**. When a system or component report contains the directory, each member’s collection breadcrumb returns to that section. **Browse specifications** links to collections with a separate index. The same alphabetical listing applies in both locations.

Above the title, breadcrumbs lead to the owning system or component and the subject’s collection. Under **Extensions**, **Actions**, and a quality model’s **Quality attributes**, select a named entry to jump to its section when several are present. Longer lists can be expanded.

On subject reports, collections that only provide navigation appear with the page’s index links. Sections with authored descriptions keep that text in the body. A component’s `requirements.md` includes requirements declared in its component hierarchy and those that explicitly apply there. Each entry links directly to its requirement.

### Follow a companion

A subject report gains a **Companions** section when its specification lists companions in its frontmatter. Entries are grouped under **Formalized by**, **Illustrated by**, **Verified by** and **Explained by**. Each entry names the file, then ` › ` and the definition when the target names one, then the author’s note after ` — `. A named definition in a supported format, such as a Quint `val` or a Python `def`, links to its declaring line as `#L<n>`. A companion that is itself a Markdown specification links to its report page. System and component reports list their own companions, and the section appears in **On this page**.

For example, the restart requirement report lists `acceptance.py › verify` under **Verified by**, with the note that the test kills the service and compares each view’s identities after it starts again. The link opens the `verify` function in the acceptance test. [Restart requirement report](../../examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md#companions)

To review every artifact that carries part of the system’s intent, open **Browse specifications → Companions** on the outermost system report. `companions.md` groups every companion in the system by role and links each to the subject that declares it. [To-do companions](../../examples/todo/specification/companions.md) · [Companion mapping](../../model/markdown-mapping.md#companions)

Each standalone authored specification in the selected project has exactly one standalone report page. The owning system or component links to its features and capabilities. Inline subjects can also have a report page for their own heading.

Authored Markdown uses `<subject>-spec`. The corresponding generated Markdown
uses `<subject>-spec-report`. For example, `<name>.system.spec.md` declares
`type: system-spec`, and its generated `<name>.system.spec-report.md` declares `type: system-spec-report`.

Each subject report records its source in a YAML `source` field using a URI
relative to the report, including the fragment for the exact authored heading.
For example, the todo system report declares
`source: "../specs/personal-task-list.system.spec.md#personal-task-list"`.
Every subject report has one **View source** link in its footer. It opens the heading recorded in `source`. For a subject defined inside
another spec, the fragment identifies its nested heading.

The footer also covers embedded declarations from the same file, such as **List owner** in the To-do system report. Each declaration retains its type and exact source heading in a `dotspec-doc` comment.

On a page that presents full passages, the first passage from each additional file carries **View section source**, or **Term source** for a glossary entry. Later passages from that file share the link. Short previews link to their subject reports. Authored citations and evidence links stay with the passages they support. The topic index jumps to term headings, and incoming relationship links return to the passage that describes the connection.

Each system quality attribute has a `system-quality-attribute-spec-report` page. Reusable quality model entries remain sections of the model’s report. The file’s `generated_by: dotspec` marker identifies generator
ownership. See the [document model](model/index.md#generateddocument).

## Check a workspace

Use `dotspec spec build --project NAME` to refresh a selected project. Run `dotspec check --all` on a clean checkout in CI, with a pinned dotspec version. Whole-project checks include the generated files for each project unless its manifest sets `report.enabled` to `false`. A missing, obsolete or stale file is a finding. The check also inspects the visible content and navigation. A subtree or individual-subject check covers only its selected source scope.

For early feedback, add `dotspec check --all` to a pre-push hook. Use CI to verify the committed tree. A pre-commit hook that enforces synchronization must check staged contents when contributors stage only part of a change.

## Record verification results

List each test or manual protocol as a companion of its requirement with `role: verifies`. Keep any account of what the procedure exercises under **Evidence**. After a run, an evidence producer can supply `verification-results.json` in the project root. Running `dotspec spec build` reads this file and publishes observations in the requirement’s **Verification** section. Running `dotspec check` checks the resulting report. The requirement file changes only when the authored decision changes.

This illustrative file requests the results of a restart assessment for revision `build-42`. Replace the revision, conditions, and observation address with the ones from your run.

```json
{
  "selection": {
    "product_revision": "build-42",
    "conditions": "local database, acknowledged writes"
  },
  "results": [
    {
      "subject": "specs/requirements/restart.requirement.spec.md",
      "procedure": "tests/acceptance.py#verify",
      "product_revision": "build-42",
      "conditions": "local database, acknowledged writes",
      "outcome": "fail",
      "observation": "evidence/restart-42.log"
    }
  ]
}
```

Identify the subject by its project-relative requirement URI, without a fragment. Identify the procedure by the companion’s target resolved against the project root, including any fragment. The restart requirement lists `../../tests/acceptance.py#verify`, so its results name `tests/acceptance.py#verify`. Relative observation URIs also start at the project root. Use `pass`, `fail`, or `inconclusive` for the outcome. A run matches when its procedure is a `verifies` companion of that requirement and its product revision and conditions equal the selection. Either record can also carry `specification_revision`. Supplying it in the selection requires an exact match. An unknown product revision or unknown conditions leave an observation outside the matching set. Use an empty results array for a selection with no observations yet. A report preserves competing results and their individual observation links.
