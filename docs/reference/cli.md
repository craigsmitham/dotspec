# CLI reference

Run project commands in the directory containing your specifications. Substitute your own names and paths in the examples. The root command's flags, including `-C`, apply to every command. Handbook reading commands work from any directory.

## dotspec

Specify what software must do and why, alongside its code

```sh
dotspec [flags]
```

Specify what software must do and why, alongside its code: author Markdown specifications, check them and the artifacts that carry their intent, and build one generated specification a reader can navigate.

Kinds: run dotspec new --help. Handbook: run dotspec docs.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec init "Booking desk"
dotspec new requirement fast-save
dotspec build
dotspec check
```

| Flag | Default | Description |
| --- | --- | --- |
| `-C, --chdir <dir>` | `` | run as if dotspec was started in `<dir>` |
| `-h, --help` | `false` | help for dotspec |
| `-v, --version` | `false` | version for dotspec |

## dotspec build

Build the generated specification

```sh
dotspec build [flags]
```

Build the generated specification for the enclosing project declared with dotspec.yaml: SPEC.md and the specification/ tree of pages that let a reader navigate the project's specifications as one. Nothing is written while a source has error findings; the findings are reported instead. Exits 1 when error findings are reported.

Learn more: dotspec docs docs/reference/inspection.md#check-the-generated-requirements-specification

```sh
dotspec build
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for build |

## dotspec check

Check specifications and reference documents in a project subtree

```sh
dotspec check [<path>] [flags]
```

Check authored specifications and the reference documents they cite, using the same discovery rules as list. With no `<path>`, check covers the enclosing project; a directory narrows it to that subtree, and a specification file checks that file and the context it requests. A whole-project check in a project declared with dotspec.yaml also checks that the generated specification is current, unless build.enabled is false. With --all, check covers every project discovered in the workspace separately and exits with the highest project exit code.

A passing check means no checked file has an error finding under dotspec's automated authoring rules. Exits 1 when error findings are reported.

Learn more: dotspec docs docs/reference/inspection.md#check-specifications

```sh
dotspec check
dotspec check internal --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | check every discovered project separately |
| `-h, --help` | `false` | help for check |
| `--json` | `false` | emit the result as JSON |

## dotspec completion

Generate a shell completion script

```sh
dotspec completion [flags]
```

Generate a completion script, then load it using your shell's completion setup.

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for completion |

## dotspec completion bash

Generate the bash completion script

```sh
dotspec completion bash [flags]
```

```sh
dotspec completion bash
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for bash |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion fish

Generate the fish completion script

```sh
dotspec completion fish [flags]
```

```sh
dotspec completion fish
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for fish |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion powershell

Generate the powershell completion script

```sh
dotspec completion powershell [flags]
```

```sh
dotspec completion powershell
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for powershell |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion zsh

Generate the zsh completion script

```sh
dotspec completion zsh [flags]
```

```sh
dotspec completion zsh
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for zsh |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec docs

Browse the offline authoring handbook

```sh
dotspec docs [<topic> | <page>[#<section>] | <rule> | <words>...] [flags]
```

Browse topics, open a page or section, or search this version's handbook and bundled examples. Topic names and exact addresses open directly. Enter words to find matching pages.

```sh
dotspec docs
dotspec docs job
dotspec docs success criteria
dotspec docs guides/write-a-job
dotspec docs job.performer
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for docs |

## dotspec help

Help about any command

```sh
dotspec help [command] [flags]
```

Show the help for any command, such as dotspec help check.

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for help |

## dotspec init

Start a specification project

```sh
dotspec init <system-name> [flags]
```

Start a specification project in the working directory. Init writes dotspec.yaml, which declares the project, and the root system specification, whose heading is `<system-name>` and whose filename is derived from it. The project name defaults to the system name's slug. Init refuses to overwrite either file and exits 2 instead. When the directory lies inside another project, init says so: that project's discovery now stops at the new manifest.

Learn more: dotspec docs system

```sh
dotspec init Bookings
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for init |
| `--json` | `false` | emit the result as JSON |
| `--name <project>` | `` | project name written to dotspec.yaml; defaults to the system name's slug |

## dotspec list

List the specifications of one kind

```sh
dotspec list <kind> [flags]
```

List the specification files of one kind in the enclosing project, grouped by the system each belongs to. list system prints systems by directory, and list component prints components by their Part of links.

Kinds:
  system              Describe a system's purpose, scope, and context
                      dotspec docs system
  component           Describe a part of a system and its responsibilities
                      dotspec docs component
  feature             Group stakeholder-visible behavior
                      dotspec docs feature
  capability          Describe an ability and the parts that provide it
                      dotspec docs capability
  job                 Define jobs to be done (JTBDs)
                      dotspec docs job
  use-case            Describe actor goals and interactions with a system
                      dotspec docs use-case
  concept             Describe purposeful product behavior through state and actions
                      dotspec docs concept
  business-rule       Describe business constraints independently of their enforcement
                      dotspec docs business-rule
  requirement         Specify required results, rules, properties, and restrictions
                      dotspec docs requirement
  quality-attribute   Adopt a quality attribute for a system, blank or from a quality model
                      dotspec docs quality-attribute
  quality-model       Define reusable quality attributes for systems to adopt
                      dotspec docs quality-model

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for list |

## dotspec list business-rule

List business rules by system

```sh
dotspec list business-rule [flags]
```

List business rules by system. Describe business constraints independently of their enforcement.

Learn more: dotspec docs business-rule

```sh
dotspec list business-rule
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for business-rule |
| `--json` | `false` | emit the result as JSON |

## dotspec list capability

List capabilities by system

```sh
dotspec list capability [flags]
```

List capabilities by system. Describe an ability and the parts that provide it.

Learn more: dotspec docs capability

```sh
dotspec list capability
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for capability |
| `--json` | `false` | emit the result as JSON |

## dotspec list component

List components by Part of

```sh
dotspec list component [flags]
```

List components by Part of. Describe a part of a system and its responsibilities.

Learn more: dotspec docs component

```sh
dotspec list component
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for component |
| `--json` | `false` | emit the result as JSON |

## dotspec list concept

List concepts by system

```sh
dotspec list concept [flags]
```

List concepts by system. Describe purposeful product behavior through state and actions.

Learn more: dotspec docs concept

```sh
dotspec list concept
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for concept |
| `--json` | `false` | emit the result as JSON |

## dotspec list feature

List features by system

```sh
dotspec list feature [flags]
```

List features by system. Group stakeholder-visible behavior.

Learn more: dotspec docs feature

```sh
dotspec list feature
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for feature |
| `--json` | `false` | emit the result as JSON |

## dotspec list job

List jobs by system

```sh
dotspec list job [flags]
```

List jobs by system. Define jobs to be done (JTBDs).

Learn more: dotspec docs job

```sh
dotspec list job
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for job |
| `--json` | `false` | emit the result as JSON |

## dotspec list quality-attribute

List quality attributes by system

```sh
dotspec list quality-attribute [flags]
```

List quality attributes by system. Adopt a quality attribute for a system, blank or from a quality model.

Learn more: dotspec docs quality-attribute

```sh
dotspec list quality-attribute
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for quality-attribute |
| `--json` | `false` | emit the result as JSON |

## dotspec list quality-model

List quality models by system

```sh
dotspec list quality-model [flags]
```

List quality models by system. Define reusable quality attributes for systems to adopt.

Learn more: dotspec docs quality-model

```sh
dotspec list quality-model
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for quality-model |
| `--json` | `false` | emit the result as JSON |

## dotspec list requirement

List requirements by system

```sh
dotspec list requirement [flags]
```

List requirements by system. Specify required results, rules, properties, and restrictions.

Learn more: dotspec docs requirement

```sh
dotspec list requirement
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for requirement |
| `--json` | `false` | emit the result as JSON |

## dotspec list system

List systems by directory

```sh
dotspec list system [flags]
```

List systems by directory. Describe a system's purpose, scope, and context.

Learn more: dotspec docs system

```sh
dotspec list system
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for system |
| `--json` | `false` | emit the result as JSON |

## dotspec list use-case

List use cases by system

```sh
dotspec list use-case [flags]
```

List use cases by system. Describe actor goals and interactions with a system.

Learn more: dotspec docs use-case

```sh
dotspec list use-case
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for use-case |
| `--json` | `false` | emit the result as JSON |

## dotspec new

Create a specification draft of one kind

```sh
dotspec new <kind> [flags]
```

Create a specification draft of one kind below the working directory and report its path with the authoring guide to read next. New never overwrites: when the target already exists, it changes nothing and exits 2.

Kinds:
  system              Describe a system's purpose, scope, and context
                      dotspec docs system
  component           Describe a part of a system and its responsibilities
                      dotspec docs component
  feature             Group stakeholder-visible behavior
                      dotspec docs feature
  capability          Describe an ability and the parts that provide it
                      dotspec docs capability
  job                 Define jobs to be done (JTBDs)
                      dotspec docs job
  use-case            Describe actor goals and interactions with a system
                      dotspec docs use-case
  concept             Describe purposeful product behavior through state and actions
                      dotspec docs concept
  business-rule       Describe business constraints independently of their enforcement
                      dotspec docs business-rule
  requirement         Specify required results, rules, properties, and restrictions
                      dotspec docs requirement
  quality-attribute   Adopt a quality attribute for a system, blank or from a quality model
                      dotspec docs quality-attribute
  quality-model       Define reusable quality attributes for systems to adopt
                      dotspec docs quality-model

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for new |

## dotspec new business-rule

Create a business-rule specification

```sh
dotspec new business-rule <slug> [flags]
```

Describe business constraints independently of their enforcement. Writes `specs/business-rules/<slug>.business-rule.spec.md` below the working directory.

Learn more: dotspec docs business-rule

```sh
dotspec new business-rule purchase-approval
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for business-rule |
| `--json` | `false` | emit the result as JSON |

## dotspec new capability

Create a capability specification

```sh
dotspec new capability <slug> [flags]
```

Describe an ability and the parts that provide it. Writes `specs/capabilities/<slug>.capability.spec.md` below the working directory.

Learn more: dotspec docs capability

```sh
dotspec new capability local-storage
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for capability |
| `--json` | `false` | emit the result as JSON |

## dotspec new component

Create a component specification

```sh
dotspec new component <name> [flags]
```

Describe a part of a system and its responsibilities. The `<name>` becomes the heading, and the filename is derived from it. --part-of names the parent system or component, relative to the working directory; Part of is never inferred from the file's location.

Learn more: dotspec docs component

```sh
dotspec new component Intake -C internal --part-of ../specs/new-system.system.spec.md
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for component |
| `--json` | `false` | emit the result as JSON |
| `--part-of <path>` | `` | parent system or component source file, relative to the working directory |

## dotspec new concept

Create a concept specification

```sh
dotspec new concept <slug> [flags]
```

Describe purposeful product behavior through state and actions. Writes `specs/concepts/<slug>.concept.spec.md` below the working directory.

Learn more: dotspec docs concept

```sh
dotspec new concept job-map
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for concept |
| `--json` | `false` | emit the result as JSON |

## dotspec new feature

Create a feature specification

```sh
dotspec new feature <slug> [flags]
```

Group stakeholder-visible behavior. Writes `specs/features/<slug>.feature.spec.md` below the working directory.

Learn more: dotspec docs feature

```sh
dotspec new feature resume-work
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for feature |
| `--json` | `false` | emit the result as JSON |

## dotspec new job

Create a job specification

```sh
dotspec new job <slug> [flags]
```

Define jobs to be done (JTBDs). Writes `specs/jobs/<slug>.job.spec.md` below the working directory.

Learn more: dotspec docs job

```sh
dotspec new job develop-requirements
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for job |
| `--json` | `false` | emit the result as JSON |

## dotspec new quality-attribute

Create a quality-attribute specification

```sh
dotspec new quality-attribute <id>... [flags]
```

Adopt a quality attribute for a system, blank or from a quality model. Each `<id>` becomes the attribute's ID and slug in the nearest enclosing system. Without --from, each attribute starts with a blank definition in `specs/quality/<id>/.` With --from, the definition is copied from the model's term and cited, and the term's missing parents are adopted too; parents the system already declares are reused. --dry-run previews the decisions and new file contents without writing.

Learn more: dotspec docs quality-attribute

```sh
dotspec new quality-attribute recoverability --from iso-25010-2023 --dry-run
```

| Flag | Default | Description |
| --- | --- | --- |
| `--dry-run` | `false` | preview the decisions and new file contents without writing |
| `--from <model>` | `` | copy from a bundled model, local model slug, or model source path |
| `-h, --help` | `false` | help for quality-attribute |
| `--json` | `false` | emit the result as JSON |

## dotspec new quality-model

Create a quality-model specification

```sh
dotspec new quality-model <slug> [flags]
```

Define reusable quality attributes for systems to adopt. With --from, the new model starts as a copy of a bundled model.

Learn more: dotspec docs quality-model

```sh
dotspec new quality-model organization-quality
dotspec new quality-model team-quality --from iso-25010-2023
```

| Flag | Default | Description |
| --- | --- | --- |
| `--from <model>` | `` | start from a copy of this bundled model |
| `-h, --help` | `false` | help for quality-model |
| `--json` | `false` | emit the result as JSON |

## dotspec new requirement

Create a requirement specification

```sh
dotspec new requirement <slug> [flags]
```

Specify required results, rules, properties, and restrictions. Writes `specs/requirements/<slug>.requirement.spec.md` below the working directory.

Learn more: dotspec docs requirement

```sh
dotspec new requirement reject-existing-target
```

| Flag | Default | Description |
| --- | --- | --- |
| `--attribute <id>` | `` | classify the requirement by this quality attribute of its governing system |
| `-h, --help` | `false` | help for requirement |
| `--json` | `false` | emit the result as JSON |

## dotspec new system

Create a system specification

```sh
dotspec new system <name> [flags]
```

Describe a system's purpose, scope, and context. The `<name>` becomes the heading, and the filename is derived from it. --part-of names the parent system, relative to the working directory; a project's root system comes from dotspec init.

Learn more: dotspec docs system

```sh
dotspec new system Billing -C internal --part-of ../specs/new-system.system.spec.md
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for system |
| `--json` | `false` | emit the result as JSON |
| `--part-of <path>` | `` | parent system source file, relative to the working directory |

## dotspec new use-case

Create a use-case specification

```sh
dotspec new use-case <slug> [flags]
```

Describe actor goals and interactions with a system. Writes `specs/use-cases/<slug>.use-case.spec.md` below the working directory.

Learn more: dotspec docs use-case

```sh
dotspec new use-case check-specifications
dotspec new use-case review-change --template expanded
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for use-case |
| `--json` | `false` | emit the result as JSON |
| `--template brief\|expanded` | `brief` | writing prompts: brief or expanded (same document format) |

## dotspec snapshot

Export and validate specification snapshots

```sh
dotspec snapshot [flags]
```

Export a project's source text, named domain fields, located relationships and inspection findings as one versioned JSON snapshot, print the snapshot's JSON Schema, or validate a saved snapshot.

Learn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for snapshot |

## dotspec snapshot export

Export a source-backed specification snapshot as JSON

```sh
dotspec snapshot export [<path>] [flags]
```

Export selected specifications, reusable quality models, enclosing systems, and direct reference context as JSON. Retain readable Markdown and source locations, including diagnostics for invalid documents. With no `<path>`, export covers the enclosing project; a directory selects that subtree. Filesystem discovery and authoring findings follow check; external references are not fetched. Exits 1 when error findings are reported.

Learn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md

```sh
dotspec snapshot export
dotspec snapshot export internal
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for export |

## dotspec snapshot schema

Print the embedded snapshot JSON Schema

```sh
dotspec snapshot schema [flags]
```

```sh
dotspec snapshot schema
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for schema |

## dotspec snapshot validate

Check a snapshot's shape, source ranges, and relationships

```sh
dotspec snapshot validate <file>|- [flags]
```

Validate a saved snapshot file, or - to read the snapshot from standard input. Source paths and external targets are not fetched.

A valid snapshot conforms to its declared model and mapping profile, and its source ranges and relationships are consistent; authoring errors it records faithfully do not make it invalid. Exits 1 when the snapshot is invalid.

Learn more: dotspec docs specs/use-cases/export-specification-model.use-case.spec.md

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for validate |
| `--json` | `false` | emit the result as JSON |

## dotspec sources

List specifications, their companions, reference documents, system context, and referenced files

```sh
dotspec sources [<path>] [flags]
```

List the files check reads: specifications, reference documents, enclosing system context and directly referenced files, each with the references that reach it. A file a specification lists as a companion shows the role that specification states for it. With no `<path>`, sources covers the enclosing project; a directory narrows it to that subtree.

A listed reference means its target exists and can be addressed, not that the target is authoritative; a companion's role is the author's declaration, not a verified relationship. Exits 1 when error findings are reported.

Learn more: dotspec docs docs/reference/inspection.md#discover-specification-sources

```sh
dotspec sources
dotspec sources internal --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `-h, --help` | `false` | help for sources |
| `--json` | `false` | emit the result as JSON |

