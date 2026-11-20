# Generate a requirements specification

Use a generated spec when readers need one entry point for requirements authored across a project. Open the project’s `SPEC.md` and follow its link to read the system’s purpose and scope, choose a feature, and follow its promises to the source and planned verification. The To-do example includes a [complete generated spec](../../examples/todo/SPEC.md).

## Set up a project

Enable the spec in the project root's `dotspec.yaml`:

```yaml
name: todo
generated_spec: true
```

The project needs one root system. Declare child systems with `Part of`; directory nesting alone establishes discovery context. Explicit specification-system selections and directory defaults place subjects in their system. Standalone quality models are reusable reference material. Links into other projects retain their original targets.

Run the following in the project:

```sh
dotspec spec build
dotspec check
```

Commit the authored changes, `SPEC.md` and the generated Markdown in `specification/` together. Follow a subject report’s source link to edit the authored text, then rebuild. Build refreshes generator-owned files, removes obsolete generated pages, and reports a collision if an authored file occupies a planned destination.

## Read the generated files

| Destination | Contents |
| --- | --- |
| `SPEC.md` | System name and a single link to its full report |
| `specification/<name>.system.spec-report.md` | Root system overview and links into its specification |
| `specification/glossary.md` | Alphabetical definitions linked to their sources, a topic index when terms declare topics, context names for distinct meanings, and links back to authored relationships |
| `specification/use-cases.md` | Use cases with their authored goals and links to the scenarios |
| `specification/quality/<attribute>/<attribute>.quality-attribute.spec-report.md` | One system attribute, its definition, parent and child attributes, and classified obligations |
| `specification/requirements.md` | Requirements and embedded obligations in this system, each linked to its canonical page or passage |
| `specification/systems/<child>/<child>.system.spec-report.md` | Child system overview with its own local directories |
| `specification/requirements-review.report.md` | Review items with links to passages for follow-up |

### Navigate a system report

The system name and authored description lead into a linked section map.
Choose a section to read in order or open a collection index directly.

| Section | What you can find |
| --- | --- |
| Goals | Purpose, objectives, supporting relationships, and the business rules index |
| User needs | Stakeholders, jobs, and desired outcomes |
| Scope | Inclusions, exclusions, environment, assumptions, constraints, project context, and glossary |
| System topology | Immediate parts, boundary partners, and required or provided interfaces, with a Mermaid diagram and linked descriptions |
| Requirements | Features, capabilities, and use cases under **Behavior**. Definitions, named uses, and synchronizations under **Concepts**. **All requirements** includes embedded obligations and component routes. |
| Quality attributes | The expanded hierarchy of declared attributes with linked definitions, the attribute index, and project quality models |
| Validation | Intended-use and requirements validation plans, success criteria, and open questions |

The map includes populated sections. Additional authored sections remain
available under **Reference**. **Requirements review** opens the project-wide
assessment and findings report. A topology edge represents a declared relationship.
A child system has its own overview showing the next level of detail.

### Browse subject reports

System and component overviews keep the source basename and change `.spec.md`
to `.spec-report.md`. For example, the Task service report lives at
`components/task-service/task-service.component.spec-report.md`. Its directory also
holds its local requirements, interfaces, and child components. The optional source
`id` provides an anchor. The subject name determines the overview filename.

The system’s `quality-attributes.md` index introduces each declared attribute with its definition preview. Shared attributes appear under each declared parent. Their own subattributes expand under the first parent. Select an attribute to read its full definition. Attribute reports nest beneath their first declared parent in the system’s `quality/` directory. Requirement reports live under their owner’s `requirements/` folder. Each attribute page links to requirements that name it. On a requirement page, read the statement, applicability, criterion, and verification plan together. Links beside the requirement lead back to the precise attribute or use-case passage.

For example, [Pick up where I left off](../../examples/todo/specification/features/pick-up-where-i-left-off.feature.spec-report.md) opens with its authored description of the behavior it offers. Follow the interaction to see how the person resumes work, or follow the restart requirement to find what must survive and how the acceptance test checks it.

A system, component, feature, capability, or concept use can supply a short
Description. The system report uses the description as its lead. Other subject reports open with the Description section. The index entry reuses
the visible text of the first paragraph. An entry with no description shows
the subject name alone. Interfaces use their introductory description and quality
models use their Description section. For requirements and business rules, a
short Statement introduces the promise or rule. Use cases use their Goal and
jobs use their Job statement.

Each overview links to the populated collections in its directory. For example,
open `use-cases.md` to choose an interaction by its goal, or `components.md` to
choose a component by its description. A system’s `systems.md` lists its immediate
child systems. When a component has entries in a collection, **By component**
on the parent index leads to the component’s index of the same name.

Within a component, `interfaces.md` groups required and provided interfaces.
Its `requirements.md` lists local requirements and embedded obligations, plus
requirements owned elsewhere that explicitly apply to the component. Each
entry links to the canonical page or passage.

Each standalone authored specification in the selected project has exactly one standalone report page. Select **View source** in the page navigation to open the authored heading. On a standalone specification this is the H1. The owning system or component links to indexes of its features and capabilities. Inline subjects can also have a report page for their own heading.

Authored Markdown uses `<subject>-spec`. The corresponding generated Markdown
uses `<subject>-spec-report`. For example, `<name>.system.spec.md` declares
`type: system-spec`, and its generated `<name>.system.spec-report.md` declares `type: system-spec-report`.

On a subject page, the **View source** link carries the model’s `documents`
relationship. Its Markdown reference definition records the relative URI once,
for example `[subject-source]: <../../specs/personal-task-list.system.spec.md#my-checklist>`.
The reference name gains a numeric suffix if it would collide with authored text.
Frontmatter identifies the report type and generator.

A subject rendered as a section carries its type and `documents` URI in a
`dotspec-doc` comment beneath the heading. Follow the source link at the end of
an embedded obligation or a passage linked from elsewhere in the reports to
open that exact heading in the authored file. Each system quality attribute has a `quality-attribute-spec-report` file. Reusable quality model entries remain sections of the model’s report. The file’s `generated_by: dotspec` marker identifies generator
ownership. The project-wide `requirements-review-report` links to the passages
that need review. See the [document model](model/index.md#generateddocument).

## Check a workspace

Use `dotspec spec build --project NAME` to refresh a selected project. Run `dotspec check --all` on a clean checkout in CI, with a pinned dotspec version. Whole-project checks include the generated files for each opted-in project. A missing, obsolete or stale file is a finding. The check also inspects the visible content and navigation. A subtree or individual-subject check covers only its selected source scope.

For early feedback, add `dotspec check --all` to a pre-push hook. Use CI to verify the committed tree. A pre-commit hook that enforces synchronization must check staged contents when contributors stage only part of a change.
