---
type: requirement
id: model-source-fidelity
attributes:
  - source-preservation
---

# Preserve source and represented content in model exports

## Applies to

Specification-model exports and validation of their saved snapshots, including
all eight supported specification families, malformed specifications, directly
consulted Markdown, and other acquired reference targets.

## Statement

For every file whose bytes are acquired during an export, dotspec must retain a
payload from which those exact bytes can be recovered, together with its digest,
byte count and source identity. Export must leave the authored file unchanged.
An unsuccessful read must retain the source identity and failure diagnostic
without claiming an acquired payload.

For acquired Markdown, dotspec must represent ordered document content and map
every supported authored domain construct to its named model role under the
declared mapping profile. Original source spans, occurrence order, nesting,
metadata, comments, empty sections, repeated sections, code, references and
unsupported regions must survive export. Retaining only an undifferentiated
source string does not satisfy the domain-mapping obligation.

Malformed input must retain the available source and content with honest parse,
mapping and diagnostic states. An invalid domain structure must not become a
fabricated valid domain record. Inspection must continue for other discovered
files. An unreadable file or failed discovery must prevent a claim that the
selected operation completed.

Unknown domain sections and uninterpreted syntax must remain visible in coverage
results. Missing known mappings or unbound content must prevent a complete
mapping claim. A complete mapping or a conforming snapshot must not be presented
as proof of stakeholder agreement, requirement satisfaction, or complete domain
knowledge.

## Rationale

Consumers need both faithful source material and explicit domain meaning to
review, navigate and transform a specification safely. Lossy exports hide
qualifications and unresolved decisions. Combining source acquisition, mapping
coverage and semantic adequacy into one success flag overstates what inspection
established.

## Verification

Export fixtures from every supported specification family and validate the
resulting snapshots. Decode retained payloads and compare them with the inputs,
including BOMs, CRLF, comments, trailing spaces and malformed metadata. Check
ordered syntax and source spans, subtype fields, named domain mappings,
references, and the distinct states for absent, empty and unresolved content.

Exercise unknown sections, malformed documents, inaccessible sources and failed
discovery. Verify that retained results and diagnostics describe what was
available, that other files are still processed, and that completion and coverage
claims remain accurate. Tamper with payloads, source spans, references, typed
memberships and coverage claims to verify that snapshot validation rejects the
inconsistencies.

The [content tests](../../internal/model/content_test.go),
[export tests](../../internal/model/export_test.go),
[export conformance tests](../../internal/model/export_conformance_test.go), and
[CLI model tests](../../internal/cli/model_test.go) provide verification evidence
for the cases they exercise. The statements above define the obligation; a test
passing does not establish semantic adequacy of arbitrary specifications.

## Related specifications

[Inspection preserves source material](../quality/inspection-preserves-source.quality.spec.md)
defines the existing source-preservation quality obligation.
[Continue checking after a structural error](continue-after-invalid-spec.requirement.spec.md)
defines continued inspection. The
[export use case](../use-cases/export-specification-model.use-case.spec.md)
places these obligations in the contributor's interaction.

The [specification model documentation](../../docs/specification-model.md) and
[Markdown mapping profile](../../model/markdown-mapping.md) describe the
representation and named mappings.
