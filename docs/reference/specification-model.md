# Specification model

A model export contains the acquired source text, declared subjects, relationships and inspection findings for a project. The generated specification uses the same domain fields to present the project to readers. See the [model reference](model/index.md) and [generated specification guide](generated-specification.md).

## Export and validate a snapshot

Why can an export report a broken link and the saved snapshot still pass validation? Suppose a use case links to a requirement that has moved. The export preserves the destination the author wrote and the finding that it could not be resolved. A tool reading the saved snapshot can now show the author which link needs repair. Dropping the link to make the export look clean would lose precisely the information that reader needs.

Snapshot validation checks the representation. To judge the inspected project, read its diagnostics and the discovery, read and inspection completion flags. An unreadable source can leave a gap even when the rest of the project was inspected. [Model source fidelity](../../specs/requirements/model-source-fidelity.requirement.spec.md) · [Interpret a check report](inspection.md#check-specifications).

In the project directory containing the specifications you want to inspect, save the export and check its exit status before validating the file. An export can write a snapshot and return 1 for source findings. Run the validation command on the saved file even in that case.

```sh
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
```

If validation succeeds, read the snapshot’s `diagnostics` and `scope.discovery_complete`, `scope.reads_complete` and `scope.inspection_complete`. For a broken link, correct the destination in the referring specification and re-export. If validation itself fails, use the reported path to locate the problem in the snapshot. An export interrupted while writing may need to be created again. Run `dotspec model schema` to inspect the single-project format.

Use `--project` to select a project or `--all` to export a workspace envelope. Each project retains its own acquisition scope and identities.

## Specification types and declared elements

A specification declares its subject through its frontmatter type and title. Systems, components, jobs, use cases, concepts, requirements, business rules, system quality attributes, features and capabilities each have their own specification type. A named item such as an interface, job step or desired outcome identifies its declaring owner. Each subject report type identifies the subject it documents.

Model `0.32.0-draft` uses mapping profile `dotspec-markdown/26`. A reusable quality model is a reference document. Export records its type in `scope.selected_reference_types` and specification types in `scope.selected_spec_types`. To migrate a model, rename `<slug>.quality-model.spec.md` to `<slug>.quality-model.md`, set `type: quality-model`, and update citations. Re-export saved snapshots and rebuild the reports. [Markdown mapping](../../model/markdown-mapping.md)

## Read one product specification

From the project root, run `dotspec model product`. For multiple independent systems, supply a system’s source basename (without `.system.spec.md`), authored ID, or snapshot record ID. For example, run `dotspec model product personal-task-list --project todo` in the dotspec checkout. Each member identifies its source, effective system, membership basis, and owned element IDs. A selected system includes its declared descendant systems. A link to a requirement in another project retains its destination. Follow the reference catalog to reusable model definitions.

`complete` concerns observed membership in the selected product. A narrowed scan or an unreadable or malformed source leaves it false. Inspect the reported scope and diagnostics before treating an empty set as an absence. An independent reference catalog can have zero specification members. The `verification` list identifies each linked procedure and the observations for that requirement or criterion. When an evidence file supplies a selection, `verification_selection` gives the requested revision and conditions.

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

For an adoption, compare the provision locator, edition, applicability, and resolution between exports. When dotspec acquires the cited document, it includes `observed_source_sha256`. A different digest flags the whole document for review. Check the cited clause to decide whether its meaning changed. Remote references retain their locator for review with the provider.
