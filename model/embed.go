// Package schema embeds the generated contract for specification snapshots.
// LinkML and Python are generation-time dependencies, not CLI dependencies.
package schema

import "embed"

// SchemaJSON is the generated JSON Schema. Callers must treat it as read-only.
//
//go:embed generated/schema.json
var SchemaJSON []byte

// CatalogJSON describes classes, inherited slots and typed reference ranges.
// Callers must treat it as read-only.
//
//go:embed generated/catalog.json
var CatalogJSON []byte

// GuideFiles contains the authored snapshot and Markdown mapping references.
//
//go:embed README.md markdown-mapping.md
var GuideFiles embed.FS
