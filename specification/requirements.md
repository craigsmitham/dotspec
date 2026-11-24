---
generated_by: dotspec
---
[dotspec](<dotspec.system.spec-report.md>)

# Requirements

- [Advice permits a successful check](<requirements/advice-permits-success.requirement.spec-report.md>) — When a completed check reports only advisory findings, dotspec must return exit status zero.
- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>) — When a discovered specification contains a structural error, dotspec must continue inspecting the remaining discovered specifications.
- [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>) — Every diagnostic must identify its severity using a text label independently of terminal styling.
- [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>) — A contributor must be able to reach the relevant versioned guidance from a specification type, a question in their own words, or a mapped diagnostic, then read the selected page, section, or worked example offline. …
- [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>) — Inspection must leave source-file contents and paths unchanged, without creating, deleting, or renaming source files.
- [Isolate specification projects in a workspace](<requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>) — Each dotspec.yaml declares the name of the specification project rooted in its directory. …
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>) — Export must preserve the meaning and order of authored Markdown, including conditions, criteria, reference destinations, and the locations of unsupported or malformed content. …
- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>) — For a project with one root system and reports enabled, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. …
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>) — When the target path already exists, initialization must terminate with an unsuccessful exit status.
- [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>) — Repeated inspections must produce equivalent structured reports, including the ordering of array elements.
