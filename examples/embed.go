// Package examples carries the worked specifications shipped with dotspec.
package examples

import "embed"

// The spec prefix also matches the authored specs directory during a clean rebuild.
//
//go:embed todo/README.md todo/dotspec.yaml todo/[Ss][Pp][Ee][Cc]* todo/tests/acceptance.py
var Files embed.FS
