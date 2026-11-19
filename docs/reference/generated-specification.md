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

Commit the authored changes, `SPEC.md` and the generated Markdown in `specification/` together. Follow each page's source link to edit the authored text, then rebuild. Build refreshes generator-owned files, removes obsolete generated pages, and reports a collision if an authored file occupies a planned destination.

## Read the generated files

| Destination | Contents |
| --- | --- |
| `SPEC.md` | System name and a single link to its full report |
| `specification/system.spec-report.md` | Root system overview and links into its specification |
| `specification/quality-attributes.spec-report.md` | The system's declared quality attributes, with their full definitions, hierarchy, qualifications, and links to classified requirements |
| `specification/system.spec-report.md#requirements` | Requirements and embedded obligations in this system, each linked to its canonical page or passage |
| `specification/systems/<child>/system.spec-report.md` | Child system overview with its own local directories |
| `specification/requirements-review.report.md` | Review items with links to passages for follow-up |

Every system directory contains a `system.spec-report.md`. In a collection, the
report filename includes the subject slug and kind, for example
`components/task-service.component.spec-report.md`. Use the canonical report
path for bookmarks and links to headings.

In a system overview, select an attribute in the name-and-description table to read its full entry. A child system with declared attributes has its own detail page beside its `system.spec-report.md`. On a requirement page, read the statement, applicability, criterion, and verification plan together. Links beside the requirement lead back to the precise attribute or use-case passage.

For example, [Pick up where I left off](../../examples/todo/specification/features/pick-up-where-i-left-off.feature.spec-report.md) opens with the goal of its grouped use case. Follow the interaction to see how the person resumes work, or follow the restart requirement to find what must survive and how the acceptance test checks it.

Each system lists its immediate child systems and the specifications in its scope.
Each component lists its structural children, declared interfaces, and requirements
that explicitly apply to it. A promise applicable to two components appears in
both directories and leads to the same page or heading.

Authored Markdown uses `<subject>-spec`. The corresponding generated Markdown
uses `<subject>-spec-report`. For example, `system.spec.md` declares
`type: system-spec`, and its generated `system.spec-report.md` declares `type: system-spec-report`.

The report’s `documents` value is a relative URI to its subject’s authored
heading. A component declared inside the system spec gets a
`component-spec-report` pointing to that heading. For a report rendered as a
section, the same two fields appear in a `dotspec-doc` HTML comment beneath its
heading. A shared quality-attributes file holds one `quality-attribute-spec-report`
per entry. The file’s `generated_by: dotspec` marker identifies generator
ownership. The project-wide `requirements-review-report` links to the passages
that need review. See the [document model](model/index.md#generateddocument).

## Check a workspace

Use `dotspec spec build --project NAME` to refresh a selected project. Run `dotspec check --all` on a clean checkout in CI, with a pinned dotspec version. Whole-project checks include the generated files for each opted-in project. A missing, obsolete or stale file is a finding. The check also inspects the visible content and navigation. A subtree or individual-subject check covers only its selected source scope.

For early feedback, add `dotspec check --all` to a pre-push hook. Use CI to verify the committed tree. A pre-commit hook that enforces synchronization must check staged contents when contributors stage only part of a change.
