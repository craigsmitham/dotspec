// Package docs embeds the public authoring handbook shipped with dotspec.
package docs

import "embed"

// Root is the repository path used for handbook citations.
const Root = "docs/cli-handbook"

// Publish the landing page, catalog and form directories. The handbook root
// can also contain instructions for repository contributors.
//
//go:embed cli-handbook/README.md cli-handbook/topics.json cli-handbook/*/*
var Files embed.FS
