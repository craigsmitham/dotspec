package speccollection

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type member struct {
	Path   string
	System *system.System
}

var notes = Collection[member]{
	Layout: specfile.Layout{Kind: "note", Collection: "notes", Suffix: ".note.spec.md", SlugExample: "work"},
	Plural: "notes",
	Read: func(path string) (*member, error) {
		if _, err := os.Stat(path); err != nil {
			return nil, err
		}
		return &member{Path: path}, nil
	},
	SetSystem: func(m *member, s *system.System) { m.System = s },
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeNote(t *testing.T, dir, slug string) {
	t.Helper()
	write(t, filepath.Join(dir, "specs", "notes", slug+".note.spec.md"), "# Note\n")
}

func writeSystem(t *testing.T, dir, name string) {
	t.Helper()
	write(t, filepath.Join(dir, "specs", specfile.NameSlug(name)+".system.spec.md"), "---\ntype: system-spec\n---\n# "+name+"\n## Purpose\nHelp people.\n## Scope\nThe project.\n")
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
}

func groupOf(owner *system.System, members []*member) []*member { return members }

func TestShowDoesNotFallBackToAnAncestorMember(t *testing.T) {
	dir := t.TempDir()
	gitInit(t, dir)
	writeNote(t, dir, "work")
	child := filepath.Join(dir, "child")
	writeSystem(t, child, "Child")
	if _, err := notes.Show(child, "work"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lookup must not fall back to an ancestor member: %v", err)
	}
}

func TestListFromSubtreeFailsOnMalformedEnclosingSystem(t *testing.T) {
	dir := t.TempDir()
	gitInit(t, dir)
	owner := filepath.Join(dir, "a")
	writeSystem(t, owner, "Owner")
	deep := filepath.Join(owner, "deep")
	writeNote(t, deep, "deep-work")
	groups, err := List(deep, notes, groupOf)
	if err != nil || len(groups) != 1 || len(groups[0]) != 1 || groups[0][0].System == nil || groups[0][0].System.Name != "Owner" {
		t.Fatalf("subtree context: %+v, %v", groups, err)
	}
	write(t, filepath.Join(owner, "specs", "new-system.system.spec.md"), "broken")
	if _, err := List(deep, notes, groupOf); err == nil {
		t.Fatal("malformed enclosing system was ignored")
	}
}
