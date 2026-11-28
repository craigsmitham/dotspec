# dotspec

dotspec helps teams record what their software must do and why. Write requirements in Markdown alongside your code. The CLI builds one specification where readers can follow a promise to the model or contract that defines it, or to the test that checks it.

## What survives a restart?

After the To-do app confirms a save, restarting it on the same database must preserve a waiting task’s expected reply and follow-up date. The report links to the Python function that tests this promise and shows whether a matching result has been recorded. [Read the requirement](examples/todo/specs/requirements/restart.requirement.spec.md#statement) · [Inspect the generated report](examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md#verification) · [Run the restart example](https://dotspec.mintlify.site/examples/restart#run-the-two-checks).

## Quick start

Install the alpha CLI with `brew install dotspec-md/tap/dotspec` or `npm install --global dotspec@alpha`. [Standalone downloads and Go installation](https://github.com/dotspec-md/dotspec/releases).

In an empty project directory, run `dotspec init "Booking desk"`. Open the system file named in the output. Describe who the desk serves and which work it handles, then build and check the specification.

```sh
dotspec build
dotspec check
```

Open `SPEC.md` and follow the system link. Can a teammate tell whether the desk takes payments? [Add a requirement to prevent double booking](https://dotspec.mintlify.site/get-started/first-requirement#give-the-desk-one-promise).

## Documentation

- [Choose an authoring guide](docs/cli-handbook/README.md) or run `dotspec docs` offline.
- [Look up a CLI command](docs/reference/cli.md).
- [Read dotspec’s own specification](SPEC.md).
- [Report a problem](https://github.com/dotspec-md/dotspec/issues) or [contribute to dotspec](CONTRIBUTING.md).

## License

[MIT License](LICENSE)
