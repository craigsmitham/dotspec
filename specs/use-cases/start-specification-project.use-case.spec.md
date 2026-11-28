---
type: use-case-spec
id: start-specification-project
specifications:
  - name: Starting a project, refusing to overwrite, and starting one inside another
    defines: requirements
    source: ../../internal/integration/testdata/init.txtar
---

# Start a specification project

## Goal

Declare a directory as a specification project with a named root system, so
that later checks and builds know which specifications belong together.

## System boundary

The dotspec CLI and the two files it writes in the working directory: the
project manifest, `dotspec.yaml`, and the root system specification.

## Primary actor

A [contributor](../dotspec.system.spec.md#contributor) beginning the
specifications for a product, or for a separately maintained part of one.

## Goal level

User task (sea level).

## Preconditions

The contributor works in the directory that will own the project.

## Minimal guarantees

Existing files are preserved according to the
[reject-existing-target requirement](../requirements/reject-existing-target.requirement.spec.md).
An unsuccessful attempt leaves no new manifest.

## Success guarantees

The directory declares a named project whose root system specification carries
the chosen system name. The contributor knows the path of each file and how to
read the guidance for describing the system.

## Main success scenario

1. The contributor names the system and, optionally, the project, and requests
   a new project in the working directory.
2. Dotspec writes the manifest and the root system specification, then reports
   both paths with a direct command for the
   [authoring guidance](../requirements/offline-authoring-guidance.requirement.spec.md).
   Without a chosen project name, the project takes the slug of the system name.
3. The contributor describes the system, then
   [creates other specification drafts](create-specification.use-case.spec.md)
   and [builds the generated specification](build-generated-specification.use-case.spec.md).

## Extensions

### Project or system already declared

#### Handling

At step 2, the directory already has a manifest or a root system specification.
Dotspec refuses to overwrite it, writes nothing and ends unsuccessfully. The
contributor edits the existing files instead.

### Invalid project name

#### Handling

At step 2, the chosen or derived project name is not a valid project name.
Dotspec writes nothing and reports the name. The contributor chooses another
name and repeats step 1.

### Directory inside another project

#### Handling

At step 2, the directory lies inside another project. Dotspec creates the new
project and names the enclosing one, whose discovery now stops at the new
manifest under the
[project isolation requirement](../requirements/specification-project-isolation.requirement.spec.md).
Continue at step 3.

## Sequence

Representative. A contributor can also write the manifest by hand; the project
isolation requirement defines what it declares.

## Authority

The [init command](../../docs/reference/cli.md#dotspec-init) and the linked
requirements define command behavior.
