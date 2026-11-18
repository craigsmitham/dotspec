// Package docs carries the versioned Markdown used by dotspec guide.
package docs

import "embed"

//go:embed README.md handbook maintainers reference
var Files embed.FS
