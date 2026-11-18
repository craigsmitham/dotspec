// Package docs carries the versioned Markdown used by dotspec guide.
package docs

import "embed"

//go:embed *.md handbook/*.md handbook/topics.json reference
var Files embed.FS
