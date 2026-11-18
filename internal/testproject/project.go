// Package testproject opens model and CLI test fixtures in isolated directories.
package testproject

import (
	"embed"
	"testing"

	"github.com/rogpeppe/go-internal/txtar"
)

//go:embed testdata/*.txtar
var fixtures embed.FS

// Open extracts a named project under the test's temporary directory.
func Open(t *testing.T, name string) string {
	t.Helper()
	data, err := fixtures.ReadFile("testdata/" + name + ".txtar")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := txtar.Write(txtar.Parse(data), dir); err != nil {
		t.Fatal(err)
	}
	return dir
}
