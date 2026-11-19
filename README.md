# dotspec

dotspec is a framework and command-line tool for writing and checking software
requirements in human-readable Markdown alongside your code.

## Quick start

Download and extract the archive for your platform from
[GitHub Releases](https://github.com/craigsmitham/dotspec/releases), then put the
binary on your PATH. From your project directory:

1. Run `dotspec system init` to create `specs/<name>.system.spec.md`.
2. Edit the system name, Description, Purpose, and Scope.
3. Run `dotspec check` to check structure and local references.

Run `dotspec guide` to read the built-in handbook offline.

## Documentation

- [Authoring handbook](docs/handbook/README.md) and [runnable To-do example](examples/todo/README.md)
- [Generate a requirements specification](docs/reference/generated-specification.md) and [browse the To-do spec](examples/todo/SPEC.md)
- [CLI reference](docs/reference/cli.md) and [documentation index](docs/README.md)
- [dotspec's requirements specification](SPEC.md) and [contributing guide](CONTRIBUTING.md)
