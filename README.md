# dotspec

dotspec is a framework and command-line tool for writing and checking software
requirements in human-readable Markdown alongside your code.

## Quick start

Download and extract the archive for your platform from
[GitHub Releases](https://github.com/craigsmitham/dotspec/releases), then put the
binary on your PATH. From your project directory:

1. Create a `dotspec.yaml` in the project root with a name, for example `name: my-app`.
2. Run `dotspec system init` to create `specs/<name>.system.spec.md`.
3. Fill in Description, Purpose, and Scope.
4. Run `dotspec spec build` to generate `SPEC.md` and its supporting reports.
5. Run `dotspec check` to check the sources and report freshness. Commit the generated files with their sources.

Run `dotspec guide` to read the built-in handbook offline.

## Documentation

- [Authoring handbook](docs/handbook/README.md) and [runnable To-do example](examples/todo/README.md)
- [Generate a requirements specification](docs/reference/generated-specification.md) and [browse the To-do spec](examples/todo/SPEC.md)
- [CLI reference](docs/reference/cli.md) and [documentation index](docs/README.md)
- [dotspec's requirements specification](SPEC.md) and [contributing guide](CONTRIBUTING.md)
