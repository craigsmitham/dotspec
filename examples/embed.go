// Package examples carries the worked specifications shipped with dotspec.
package examples

import "embed"

//go:embed concept-design jtbd-presets requirements-review requirements-traceability starter system-structure user-needs todo/README.md todo/COVERAGE.md todo/dotspec.yaml todo/specs todo/tests/acceptance.py
var Files embed.FS
