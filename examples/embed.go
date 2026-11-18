// Package examples carries the worked specifications shipped with dotspec.
package examples

import "embed"

//go:embed concept-design jtbd-presets requirements-review requirements-traceability starter system-structure user-needs
var Files embed.FS
