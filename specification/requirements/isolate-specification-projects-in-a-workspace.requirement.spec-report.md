---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/specification-project-isolation.requirement.spec.md#isolate-specification-projects-in-a-workspace"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../dotspec.system.spec-report.md#requirements>)

# Isolate specification projects in a workspace

## Statement

Each `dotspec.yaml` declares the name of the specification project rooted in
its directory. A spec belongs to the nearest enclosing project. Discovery must
exclude nested projects, and implicit System lookup must stop at the project
root. Explicit references may acquire context across project boundaries.

From anywhere in a Git worktree, contributors must be able to select a project
by naming a directory within it, and to check all discovered projects,
including when the worktree root has no manifest. A workspace-wide check must
preserve project identity, check every discovered project separately, and
return the highest exit code.

## Applies to

Workspace-aware discovery and System context in checks, source inventories and
snapshot exports; CLI project selection by directory (`dotspec -C <dir>`); and the
workspace-wide check (`dotspec check --all`), the only workspace-wide
inspection. Repositories without project manifests retain their existing
discovery behavior.

## Rationale

A reference app or unrelated monorepo product must not become part of dotspec's own
specification collection because of directory nesting. System relationships remain
explicit domain assertions, independent of tooling partitions.

## Evidence

Checks for this requirement include [discovery tests](<../../internal/specfile/workspace_test.go>), [CLI scenarios](<../../internal/cli/workspace_test.go>).

---

[View source](<../../specs/requirements/specification-project-isolation.requirement.spec.md#isolate-specification-projects-in-a-workspace>)
