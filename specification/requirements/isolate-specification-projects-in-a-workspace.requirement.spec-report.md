---
generated_by: dotspec
type: requirement-spec-report
---
# Isolate specification projects in a workspace

[dotspec](<../dotspec.system.spec-report.md>) · [View source][subject-source]

## Statement

Each `dotspec.yaml` declares the name of the specification project rooted in
its directory. A spec belongs to the nearest enclosing project. Discovery must
exclude nested projects, and implicit System lookup must stop at the project
root. Explicit references may acquire context across project boundaries.

From anywhere in a Git worktree, contributors must be able to select a project
by name or inspect all discovered projects, including when the worktree root
has no manifest. Workspace-wide operations must preserve project identity,
attempt every selected project, and return the highest exit code.

## Applies to

Workspace-aware discovery, System context, source inventories, model exports and
reading views, and CLI project selection. Repositories without project
manifests retain their existing discovery behavior.

## System verification plan

Exercise nested roots under both Git and filesystem discovery; a workspace without
a root project; missing local Systems; explicit cross-project references; scoped
and workspace-wide commands; malformed manifests; and offline export validation.

The [discovery tests](<../../internal/specfile/workspace_test.go>) and
[CLI scenarios](<../../internal/cli/workspace_test.go>) provide executable cases.
Check the actual repository's isolated root export and all-project check as well.

[View source for “System verification plan”](<../../specs/requirements/specification-project-isolation.requirement.spec.md#system-verification-plan>)

## Rationale

A reference app or unrelated monorepo product must not become part of dotspec's own
specification collection because of directory nesting. System relationships remain
explicit domain assertions, independent of tooling partitions.

[subject-source]: <../../specs/requirements/specification-project-isolation.requirement.spec.md#isolate-specification-projects-in-a-workspace>
