// Package examples carries the worked specifications shipped with dotspec.
package examples

import "embed"

//go:embed todo/README.md todo/COVERAGE.md todo/dotspec.yaml todo/specs todo/tests/acceptance.py todo/views
var Files embed.FS
