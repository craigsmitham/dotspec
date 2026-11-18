# dotspec

dotspec is a framework and CLI for developing software product requirements in
human-readable Markdown alongside the code. Describe what people need, the
behavior a system promises, and how to assess it. Link to policies, contracts
and tests in their native form.

## Install

Download the binary for your platform from the [releases](https://github.com/craigsmitham/dotspec/releases)
and optionally the AXM Knowledge bundle. You can also build from this checkout with `go build -o dotspec .`. Run `dotspec guide` for
the handbook included with your version. It works offline in any directory.

## Run

From this checkout with the Go version in [go.mod](go.mod), explore the
framework and check the specifications that describe dotspec itself:

```sh
go run . guide
go run . guide obligation
go run . use-case show check-specifications
go run . check
```

In a project of your own, run `dotspec system init`, edit the generated Purpose
and Scope, then run `dotspec check`. Each starter points to its authoring guide.

## Find guidance

- Choose what to write and follow a user need through the runnable To-do app and its restart test. [Handbook](docs/handbook/README.md)
- Record where your team keeps specifications and who owns a promise. [Project conventions](docs/handbook/project-conventions.md)
- Look up file placement, Markdown formats, command options and check reports. [Layout and formats](docs/reference/layout.md) · [CLI](docs/cli.md) · [Inspection](docs/reference/inspection.md)
- Inspect the types, relationships and mapping into an exported snapshot. [Core domain](docs/reference/model/core.md) · [Specification model](docs/specification-model.md)
- Audit a domain decision or change the toolkit. [Maintainer's guide](docs/handbook/maintaining.md) · [Contributing](CONTRIBUTING.md)

## Design principle: keep meaning in the source; keep connections simple

Connect specification sources through ordinary file references, preserving their
native form. dotspec provides discovery and reference checking. Consuming projects
determine which sources carry specification weight and ensure those sources are
discoverable; readers and agents interpret their meaning.

This principle guides dotspec's design to reduce duplicate maintenance and keep
sources useful to people and tools. Prefer ordinary files, Markdown links, and
explanatory prose. Add machine-readable structure only when a demonstrated
operation cannot reasonably work without it.

An executable specification can own its required examples: link to it rather
than copying assertions into Markdown. When a test checks a requirement already
stated in prose, retain that definition and explain the test's verification role.
A passing implementation test does not acquire authority merely by being linked.
File discovery makes sources available together; deciding whether they adequately
specify a system remains an authoring and review task.
