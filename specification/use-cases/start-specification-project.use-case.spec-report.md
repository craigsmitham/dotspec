---
generated_by: dotspec
type: use-case-spec-report
source: "../../specs/use-cases/start-specification-project.use-case.spec.md#start-a-specification-project"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Use cases](<../dotspec.system.spec-report.md#use-cases>)

# Start a specification project

## Goal

Declare a directory as a specification project with a named root system, so
that later checks and builds know which specifications belong together.

## Primary actor

A [contributor](<../dotspec.system.spec-report.md#contributor>) beginning the
specifications for a product, or for a separately maintained part of one.

## Preconditions

The contributor works in the directory that will own the project.

## Main success scenario

<!-- dotspec-doc {"type":"scenario-spec-report","documents":"../../specs/use-cases/start-specification-project.use-case.spec.md#main-success-scenario"} -->

1. The contributor names the system and, optionally, the project, and requests
   a new project in the working directory.
2. Dotspec writes the manifest and the root system specification, then reports
   both paths with a direct command for the
   [authoring guidance](<../requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>).
   Without a chosen project name, the project takes the slug of the system name.
3. The contributor describes the system, then
   [creates other specification drafts](<create-specification.use-case.spec-report.md>)
   and [builds the generated specification](<build-generated-specification.use-case.spec-report.md>).

## Extensions

[Project or system already declared](<#project-or-system-already-declared>) · [Invalid project name](<#invalid-project-name>) · [Directory inside another project](<#directory-inside-another-project>)

### Project or system already declared

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/start-specification-project.use-case.spec.md#project-or-system-already-declared"} -->

#### Handling

At step 2, the directory already has a manifest or a root system specification.
Dotspec refuses to overwrite it, writes nothing and ends unsuccessfully. The
contributor edits the existing files instead.

### Invalid project name

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/start-specification-project.use-case.spec.md#invalid-project-name"} -->

#### Handling

At step 2, the chosen or derived project name is not a valid project name.
Dotspec writes nothing and reports the name. The contributor chooses another
name and repeats step 1.

### Directory inside another project

<!-- dotspec-doc {"type":"use-case-extension-spec-report","documents":"../../specs/use-cases/start-specification-project.use-case.spec.md#directory-inside-another-project"} -->

#### Handling

At step 2, the directory lies inside another project. Dotspec creates the new
project and names the enclosing one, whose discovery now stops at the new
manifest under the
[project isolation requirement](<../requirements/isolate-specification-projects-in-a-workspace.requirement.spec-report.md>).
Continue at step 3.

## Success guarantees

The directory declares a named project whose root system specification carries
the chosen system name. The contributor knows the path of each file and how to
read the guidance for describing the system.

## Minimal guarantees

Existing files are preserved according to the
[reject-existing-target requirement](<../requirements/reject-existing-target.requirement.spec-report.md>).
An unsuccessful attempt leaves no new manifest.

## System boundary

The dotspec CLI and the two files it writes in the working directory: the
project manifest, `dotspec.yaml`, and the root system specification.

## Goal level

User task (sea level).

## Sequence

Representative. A contributor can also write the manifest by hand; the project
isolation requirement defines what it declares.

## Authority

The [init command](<../../docs/reference/cli.md#dotspec-init>) and the linked
requirements define command behavior.

## Specification sources

<a id="specification.starting-a-project-refusing-to-overwrite-and-starting-one-inside-another"></a>

### [Starting a project, refusing to overwrite, and starting one inside another](<../../internal/integration/testdata/init.txtar>)

Defines requirements.

---

[View source](<../../specs/use-cases/start-specification-project.use-case.spec.md#start-a-specification-project>)
