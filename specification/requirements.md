---
generated_by: dotspec
---
# Requirements

[dotspec](<dotspec.system.spec-report.md>)

- [Advice permits a successful check](<requirements/advice-permits-success.requirement.spec-report.md>) — When a completed check reports only advisory findings, dotspec must return exit status zero.
- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>) — When a discovered specification contains a structural error, dotspec must continue inspecting the remaining discovered specifications.
- [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>) — Every diagnostic must identify its severity using a text label independently of terminal styling.
- [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>) — Inspection must leave source-file contents and paths unchanged, without creating, deleting, or renaming source files.
- [Isolate specification projects in a workspace](<requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>) — Each dotspec.yaml declares the name of the specification project rooted in its directory. A spec belongs to the nearest enclosing project. Discovery must exclude nested projects, and implicit System lookup must stop at the project root. Explicit references may acquire context across project boundaries.
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>) — For every file whose bytes are acquired during an export, dotspec must retain a payload from which those exact bytes can be recovered, together with its digest, byte count and source identity. Export must leave the authored file unchanged. An unsuccessful read must retain the source identity and failure diagnostic without claiming an acquired payload.
- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>) — For an opted-in project with one root system, contributors must be able to generate a deterministic Markdown requirements specification and check its fidelity, navigation, and freshness with the project-wide check. A reader must be able to reach each authored promise in scope, read its qualifications, follow its declared relationships and verification guidance, and locate the editable source.
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>) — When the target path already exists, initialization must terminate with an unsuccessful exit status.
- [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>) — Repeated inspections must produce equivalent structured reports, including the ordering of array elements.
