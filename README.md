# dotspec

dotspec is a framework and command-line tool for writing and checking software
requirements in human-readable Markdown alongside your code.

## Quick start

Start with the CLI on your PATH. If a [release archive](https://github.com/craigsmitham/dotspec/releases) is available for your platform, extract its executable. To build from a checkout, run `mise trust`, `mise install`, then `mise run build` at the repository root. Add the resulting `dist` directory to your PATH (on Windows, build with `mise exec -- go build -o dist/dotspec.exe .`). Run `dotspec --version` from the project where you want to keep the specifications.

To start the first specification, work from your project root.

1. Run `dotspec init "Booking desk"`, substituting your system name. It writes `dotspec.yaml` and the root system specification, and refuses to replace either file. Open the system file named in the output.
2. Replace the starter prose with a description, the benefit someone seeks, and the work the system covers. [Describe a system and its boundary](docs/cli-handbook/guides/describe-a-system.md).
3. Run `dotspec build`. Open `SPEC.md` and follow the system link. Can a teammate tell which work your system takes responsibility for?
4. Run `dotspec check`. If the generated specification is stale after a source edit, rebuild and check again. Commit the manifest, sources, `SPEC.md`, and `specification/` together.

Run `dotspec new --help` to see the kinds of specification you can add, such as `dotspec new requirement fast-save`, and `dotspec list requirement` to find them again. Run `dotspec docs` to browse the offline handbook. Use `dotspec docs job` to choose guidance on jobs to be done, or `dotspec docs success criteria` to search.

## Documentation

- [Authoring handbook](docs/cli-handbook/README.md) and [runnable To-do example](examples/todo/README.md)
- [Generate a requirements specification](docs/reference/generated-specification.md) and [browse the To-do spec](examples/todo/SPEC.md)
- [CLI reference](docs/reference/cli.md) and [documentation index](docs/README.md)
- [dotspec's requirements specification](SPEC.md) and [contributing guide](CONTRIBUTING.md)

## License

dotspec is released under the [MIT License](LICENSE).
