package businessrule

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/craigsmitham/dotspec/internal/system"
)

func TestLifecycleDiscoveryAndContext(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	if _, err := system.Init(root); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "child")
	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatal(err)
	}
	path, err := Init(child, "approval")
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Init(child, "approval"); err == nil {
		t.Fatal("overwrote existing rule")
	}
	got, err := Show(child, "approval")
	if err != nil || got.System == nil || got.Statement == "" {
		t.Fatalf("read/context: %+v %v", got, err)
	}
	if _, err := Show(root, "approval"); !os.IsNotExist(err) {
		t.Fatalf("show searched descendants: %v", err)
	}
	groups, err := List(root)
	if err != nil || len(groups) != 1 || len(groups[0].BusinessRules) != 1 {
		t.Fatalf("inventory: %+v %v", groups, err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatal("inspection changed source")
	}
	for _, slug := range []string{"../escape", "Uppercase", "a/b", "a--b"} {
		if _, err := Init(root, slug); err == nil {
			t.Fatalf("accepted slug %s", slug)
		}
	}
}
