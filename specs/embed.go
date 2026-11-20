// Package specs provides dotspec's own specifications as offline authoring examples.
package specs

import "embed"

//go:embed *.spec.md */*.spec.md quality
var Files embed.FS
