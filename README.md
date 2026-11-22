# dotspec

dotspec is a framework and command-line tool for writing and checking software
requirements in human-readable Markdown alongside your code.

## Quick start

Start with the CLI on your PATH. If a [release archive](https://github.com/craigsmitham/dotspec/releases) is available for your platform, extract its executable. To build from a checkout, run `mise trust`, `mise install`, then `mise run build` at the repository root. Add the resulting `dist` directory to your PATH (on Windows, build with `mise exec -- go build -o dist/dotspec.exe .`). Run `dotspec --version` from the project where you want to keep the specifications.

To start the first specification, work from your project root.

1. Create a `dotspec.yaml` in the project root with a name, for example `name: my-app`.
2. Run `dotspec system init "Booking desk"`, substituting your system name. Open the file named in the output. If this directory already declares the system, edit that file.
3. Replace the starter prose with a description, the benefit someone seeks, and the work the system covers. [Describe a system and its boundary](docs/cli-handbook/guides/describe-a-system.md).
4. Run `dotspec spec build`. Open `SPEC.md` and follow the system link. Can a teammate tell which work your system takes responsibility for?
5. Run `dotspec check`. If the report is stale after a source edit, rebuild and check again. Commit the manifest, sources, `SPEC.md`, and `specification/` together.

Run `dotspec docs` to browse the offline handbook. Use `dotspec docs job` to choose guidance on jobs to be done, or `dotspec docs success criteria` to search.

## Documentation

- [Authoring handbook](docs/cli-handbook/README.md) and [runnable To-do example](examples/todo/README.md)
- [Generate a requirements specification](docs/reference/generated-specification.md) and [browse the To-do spec](examples/todo/SPEC.md)
- [CLI reference](docs/reference/cli.md) and [documentation index](docs/README.md)
- [dotspec's requirements specification](SPEC.md) and [contributing guide](CONTRIBUTING.md)
