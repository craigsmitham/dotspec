---
type: requirement
---
# Isolate specification projects in a workspace

## Statement

When a workspace registers distinct specification project roots, dotspec must
limit a selected project's discovered specifications to its most specific root
partition and stop implicit System context lookup at that root. Explicit
references may acquire another project's context without recursively selecting
its specifications. Workspace-wide operations must preserve project identity and
separate results, attempt every selected project, and report the highest exit code.

## Rationale

A reference app or unrelated monorepo product must not become part of dotspec's own
specification collection because of directory nesting. System relationships remain
explicit domain assertions, independent of tooling partitions.

## Applies to

Workspace-aware discovery, System context, source inventories, model exports and
reading views, and CLI project selection. Repositories without a workspace
manifest retain their existing discovery behavior.

## Verification plan

Exercise nested roots under both Git and filesystem discovery; a workspace without
a root project; missing local Systems; explicit cross-project references; scoped
and workspace-wide commands; malformed manifests; and offline export validation.

The [discovery tests](../../internal/specfile/workspace_test.go) and
[CLI scenarios](../../internal/cli/workspace_test.go) provide executable cases.
Check the actual repository's isolated root export and all-project check as well.
