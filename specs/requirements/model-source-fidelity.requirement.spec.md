---
type: requirement-spec
id: model-source-fidelity
quality_attributes: [source-preservation]
---

# Preserve source and represented content in specification snapshots

## Applies to

Snapshot exports and validation of saved snapshots, including
all supported specification types, malformed specifications, directly
consulted Markdown, and other acquired reference targets.

## Statement

Export must preserve the meaning and order of authored Markdown, including conditions,
criteria, reference destinations, and the locations of unsupported or malformed
content. Source locations refer to the normalized Markdown included in the snapshot.
The original files remain authoritative for exact bytes. A native contract, executable
specification, or other non-Markdown target is identified by its path and read status.

An unreadable source retains its identity and diagnostic. Inspection must continue for
the remaining sources and report an incomplete read. A malformed specification must
retain the available text and findings without acquiring a fabricated domain record.

## Criterion

Every valid mapped specification file contributes one primary record of its declared type. A named part retains the identity of its declaring source. Each reusable quality model contributes a QualityModel reference root and its owned definitions. The export distinguishes selected reference types from selected specification types. Only requirement files contribute Requirement records. Multiple references to one requirement resolve to that same record. Offline validation reports a mismatched primary record or multiple primary subjects for one source.

## Rationale

Consumers need both faithful source material and explicit domain meaning to
review, navigate and transform a specification safely. Lossy exports hide
qualifications and unresolved decisions. Source ranges let consumers trace each represented field back to its authored passage.

## Evidence

Checks for this requirement include [content tests](../../internal/model/content_test.go), [export tests](../../internal/model/export_test.go), [export conformance tests](../../internal/model/export_conformance_test.go), [CLI snapshot tests](../../internal/cli/snapshot_test.go).

## Related specifications

[Inspection preserves source material](../quality/reliability/source-preservation/inspection-preserves-source.requirement.spec.md)
defines the existing source-preservation quality requirement.
[Continue checking after a structural error](continue-after-invalid-spec.requirement.spec.md)
defines continued inspection. The
[export use case](../use-cases/export-specification-model.use-case.spec.md)
places these requirements in the contributor's interaction.

The [specification model documentation](../../docs/reference/specification-model.md) and
[Markdown mapping profile](../../model/markdown-mapping.md) describe the
representation and named mappings.
