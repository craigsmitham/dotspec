# Generate a requirements specification

Use a generated spec when readers need one entry point for requirements authored across a project. Start at the project's `SPEC.md` to read its purpose and scope, choose a feature, and follow its promises to the source and planned verification. The To-do example includes a [complete generated spec](../../examples/todo/SPEC.md).

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
| `SPEC.md` | Root system overview and links into its specification |
| `specification/index.md` | Contents across the selected project by system and kind |
| `specification/<kind>/index.md` | Members of one collection in the root system |
| `specification/systems/<child>/SPEC.md` | Overview of a declared child system with collections in its directory |
| `specification/review.md` | Review items with links to passages for follow-up |

Generated Markdown frontmatter identifies a document role. A page about a requirement uses `type: specification-section` and `subject_kind: requirement`. The authoritative authored requirement retains `type: requirement`. See the [documentation model](model/index.md#specificationdocumenttype).

## Check a workspace

Use `dotspec spec build --project NAME` to refresh a selected project. Run `dotspec check --all` on a clean checkout in CI, with a pinned dotspec version. Whole-project checks include the generated files for each opted-in project. A missing, obsolete or stale file is a finding. The check also inspects the visible content and navigation. A subtree or individual-subject check covers only its selected source scope.

For early feedback, add `dotspec check --all` to a pre-push hook. Use CI to verify the committed tree. A pre-commit hook that enforces synchronization must check staged contents when contributors stage only part of a change.
