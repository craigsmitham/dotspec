package businessrule

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dotspec-md/dotspec/internal/speccollection"
	"github.com/dotspec-md/dotspec/internal/system"
)

func TestLifecycleDiscoveryAndContext(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%s: %v", out, err)
	}
	if _, err := system.InitNamed(root, "New system"); err != nil {
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
	if got, err := Read(path); err != nil || got.Statement == "" {
		t.Fatalf("read: %+v %v", got, err)
	}
	listed, err := speccollection.Members(root, Layout, "business rules", func(path string) (string, error) {
		rule, err := Read(path)
		if err != nil {
			return "", err
		}
		return rule.Name, nil
	})
	if err != nil || len(listed) != 1 || listed[0].Path != path || listed[0].System == nil {
		t.Fatalf("inventory and context: %+v %v", listed, err)
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
