package diagnostic

import (
	"errors"
	"io/fs"
	"testing"
)

func TestLocationsUseUnicodeColumns(t *testing.T) {
	source := []byte("first\nαβ [link]")
	location := At(source, len("first\nαβ "))
	if location.Line != 2 || location.Column != 4 {
		t.Fatalf("location: %+v", location)
	}
}

func TestWithPathRetainsErrorIdentityAndRelatedLocations(t *testing.T) {
	d := Error("identity.duplicate", "duplicate", Location{Line: 8, Column: 1}, Location{Line: 3, Column: 5})
	var got *Diagnostic
	if !errors.As(WithPath(d, "specs/jobs/work.job.spec.md"), &got) || got.Location.Path == "" || got.Related[0].Path != got.Location.Path {
		t.Fatalf("location paths: %+v", got)
	}
	if d.Location.Path != "" || d.Related[0].Path != "" {
		t.Fatal("adding a file path mutated the original diagnostic")
	}
	if !errors.Is(WithPath(fs.ErrNotExist, "missing"), fs.ErrNotExist) {
		t.Fatal("wrapped I/O error lost its cause")
	}
}
