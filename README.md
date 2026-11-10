# dotspec

A framework and CLI for managing software product requirements in `*.spec.md` files.

System specs describe a system's name, purpose, and scope. Subsystems use the
same format, placed alongside the code they describe.

## Run

With Go 1.27.3 or later installed:

```sh
go run . system show
go run . system list
```

## Build

```sh
go build -o dotspec .
./dotspec --help
```

## System commands

```sh
dotspec system init                     # Create a system spec here
dotspec system show                     # Describe the nearest enclosing system
dotspec system list                     # Discover systems in this subtree
dotspec system show --path internal/cli
dotspec system list --path internal
```

Use `./dotspec` for a local build, or replace `dotspec` with `go run .` when
running from the repository root.

Every command accepts `--path DIR`, defaulting to the current directory. The
directory must already exist. Use `dotspec system <command> --help` for usage.

- `init` creates `system.spec.md` with editable placeholders for the name,
  purpose, and scope. It never overwrites an existing file.
- `show` reads the selected directory's system spec, or the nearest ancestor's
  if there is no local spec. It displays the name, purpose, scope, nearest parent
  system, and source path. Ancestor lookup stops at the Git repository root
  (including a worktree root), or the filesystem root outside Git.
- `list` displays a tree of systems within the selected directory and its
  descendants, ordered by source path. It does not include ancestor systems
  outside that subtree. An empty result is successful and reports that no
  systems were found.

In Git repositories, discovery uses Git's ignore rules, including nested
`.gitignore` files and local exclusions. Tracked files remain visible even if
an ignore rule matches them. Git must be installed for discovery in a repository.
Outside Git, discovery walks the directory tree without applying ignore rules.
Discovery skips `.git`, nested repositories, and symlinked files and directories;
an explicitly selected symlink directory is resolved before the command runs.

Missing systems for `show`, invalid paths, malformed specs, and command errors
produce a diagnostic on standard error and a nonzero exit status. A malformed
discovered system makes `list` fail rather than silently omit it.

## System spec format

Each directory can contain one system spec, named `system.spec.md`:

```markdown
---
type: system
---

# Spec parser

## Purpose

Make specifications available to dotspec commands as structured information.

## Scope

Reading and interpreting specification files.
```

The frontmatter currently accepts only `type: system`. The Markdown body must
contain exactly one nonempty document-level H1 naming the system, followed by
nonempty `Purpose` and `Scope` H2 sections. Section names are case-sensitive.
Section bodies can contain Markdown, including lists, code blocks, and deeper
headings. Other sections are allowed but are not displayed by `system show`.
There is no requirements section or requirement inheritance behavior.

The nearest ancestor directory containing a system spec defines the parent:

```text
system.spec.md                  # dotspec
internal/parser/system.spec.md  # Spec parser, a subsystem of dotspec
```

Moving a spec can therefore change its parent. Source paths in command output
are relative to the current directory.

## Development

```sh
go test ./...
go vet ./...
```

Tests that exercise Git ignore rules and repository boundaries require Git.
