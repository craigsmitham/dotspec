# Specification model

A model export contains the acquired source text, declared subjects, relationships and inspection findings for a project. The generated specification uses the same domain fields to present the project to readers. See the [model reference](model/index.md) and [generated specification guide](generated-specification.md).

## Export and validate a snapshot

```sh
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
dotspec model schema
```

Use `--project` to select a project or `--all` to export a workspace envelope. Each project retains its own acquisition scope and identities. Validation runs offline. A valid snapshot can describe invalid authored input. Read the diagnostics and the discovery, read and inspection completion flags before using the result.

## Specification types and declared elements

A specification declares its subject through its frontmatter type and title. Systems, components, jobs, use cases, concepts, requirements, business rules, quality models, system quality attributes, features and capabilities each have their own specification type. A named item such as an interface, job step or desired outcome identifies its declaring owner. Each subject report type identifies the subject it documents.

Model `0.30.0-draft` uses mapping profile `dotspec-markdown/24`. Regenerate saved snapshots after moving former inline requirements to their own files. For scenario outcomes, retain the prose in the use case. Update incoming links before rebuilding reports. [Migration guidance](formats/requirement.md#requirement-format)


## Source text and locations

Look up the document by source ID. Its Markdown has LF line endings and no leading byte-order mark. Each rich field contains an ordered list of ranges in that source. Slice UTF-8 bytes from the inclusive start to the exclusive end. An empty range list represents an authored empty field. Omission means the field was absent. An unreadable file retains its locator and read status. A malformed specification retains its available text and diagnostics. A directly referenced native artifact is identified by its locator. Referenced Markdown provides its text and navigable headings.

## References and semantic relationships

Start with `references`. For each link, `authored_target` preserves the destination and `span` locates the occurrence. A typed occurrence adds its `role` and `context_owner`. Its `resolution` holds the status and, on success, the target ID. To find a requirement’s supporting links, select occurrences with that requirement as owner and `supports` as role. To navigate backwards, select resolved occurrences whose target is the requirement. Read ordered declaration membership directly from the owner. The [relationship reference](model/relationships.md) lists the allowed endpoint classes and meanings.

### Quality attribute grouping

A label occurrence records its owner, authored label and governing vocabulary. Successful resolution identifies a specific quality attribute. A requirement can name several attributes. A quality attribute can have multiple declared parents. Each attribute page links to the requirements classified under it.

## Use-case narratives

A use case owns Markdown fields for actors, goal level, conditions and guarantees. Its main success scenario may have ordered steps. An extension owns its condition and handling. These passages keep the author’s wording in the report. See the [Markdown mapping profile](../../model/markdown-mapping.md#use-case).

## Concept definitions and composition

A concept declares its purpose, state, named actions and operational principles. A system or component declares concept uses and synchronizations. Definition and dependency links identify their targets. A synchronization retains its authored When, Where, Then and Semantics fields. See the [core domain diagram](model/core.md).
