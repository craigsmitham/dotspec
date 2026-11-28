package specdoc

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestInstructionLifecycle(t *testing.T) {
	root := t.TempDir()
	p := testPlan(t, root)
	p.InstructionFiles = nil
	const original = "<!-- axm:knowledge -->\r\nRules.\r\n<!-- /axm:knowledge -->\r\n\r\n"
	agents := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(agents, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	originalInfo, err := os.Stat(agents)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("AGENTS.md", filepath.Join(root, "CLAUDE.md")); err != nil {
		t.Skip(err)
	}
	read := func(rel string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	build := func() Changes {
		t.Helper()
		changes, err := Write(root, p)
		if err != nil {
			t.Fatal(err)
		}
		return changes
	}
	checkPending := func(rel string) {
		t.Helper()
		ds, err := Check(root, p)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range ds {
			if strings.HasSuffix(filepath.ToSlash(d.Location.Path), "/"+rel) {
				return
			}
		}
		t.Fatalf("missing finding for %s in %v", rel, ds)
	}
	checkPending("AGENTS.md")
	if read("AGENTS.md") != original {
		t.Fatal("check wrote instructions")
	}
	old := filepath.Join(root, "SPEC.md")
	if err := os.WriteFile(old, []byte("---\ngenerated_by: dotspec\n---\n# Old\n"), 0644); err != nil {
		t.Fatal(err)
	}
	build()
	if _, err := os.Lstat(old); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("retained old report", err)
	}
	if err := os.WriteFile(old, []byte("# Authored"), 0644); err != nil {
		t.Fatal(err)
	}
	body := read("AGENTS.md")
	if !strings.HasPrefix(body, original) || strings.Count(body, blockStart) != 1 || !strings.Contains(body, "[project specification index](<specification/README.md>)") || strings.Contains(strings.ReplaceAll(body, "\r\n", ""), "\n") {
		t.Fatal(body)
	}
	if info, err := os.Stat(agents); err != nil || info.Mode().Perm() != originalInfo.Mode().Perm() {
		t.Fatal(info, err)
	}
	if dest, err := os.Readlink(filepath.Join(root, "CLAUDE.md")); err != nil || dest != "AGENTS.md" {
		t.Fatal(dest, err)
	}
	if changes := build(); len(changes.Written)+len(changes.Removed) != 0 {
		t.Fatal(changes)
	}
	custom := ".github/copilot-instructions.md"
	p.InstructionFiles = []string{custom}
	changes := build()
	if !slices.Contains(changes.Written, "AGENTS.md") || read("AGENTS.md") != original || !strings.Contains(read(custom), "<../specification/README.md>") {
		t.Fatal(changes)
	}
	if err := os.Remove(filepath.Join(root, custom)); err != nil {
		t.Fatal(err)
	}
	checkPending(custom)
	build()
	p.InstructionFiles = []string{}
	checkPending(custom)
	build()
	if read(custom) != "" || read("SPEC.md") != "# Authored" {
		t.Fatal("opt-out or migration damaged an authored file")
	}
	if changes := build(); len(changes.Written)+len(changes.Removed) != 0 {
		t.Fatal(changes)
	}
}

func TestInstructionMarkers(t *testing.T) {
	fenced := "~~~md\n" + blockStart + "\n" + blockEnd + "\n~~~\n\n"
	first, err := updateBlock([]byte(fenced), "Read this.")
	if err != nil || !bytes.HasPrefix(first, []byte(fenced)) {
		t.Fatal(string(first), err)
	}
	for _, s := range []string{blockStart, blockEnd, blockEnd + "\n" + blockStart, blockStart + "\n" + blockStart + "\n" + blockEnd, string(first) + blockStart + "\n" + blockEnd} {
		if _, err := updateBlock([]byte(s), "Read this."); err == nil {
			t.Errorf("accepted malformed markers: %q", s)
		}
	}
	next, err := updateBlock(first, "Read this.")
	if err != nil || !bytes.Equal(next, first) {
		t.Fatal(string(next), err)
	}
}

func TestInstructionDiscovery(t *testing.T) {
	for _, test := range []struct {
		file, content string
		shared        bool
	}{
		{"CLAUDE.md", "@AGENTS.md\n\nClaude-specific text.", true},
		{".claude/CLAUDE.md", "@../AGENTS.md\n", true},
		{"CLAUDE.md", "See [instructions](AGENTS.md).", false},
		{"CLAUDE.md", "```text\n@AGENTS.md\n```", false},
	} {
		t.Run(test.content, func(t *testing.T) {
			root := t.TempDir()
			p := testPlan(t, root)
			p.InstructionFiles = nil
			for rel, body := range map[string]string{"AGENTS.md": "# Instructions", "AGENTS.override.md": "Override", test.file: test.content} {
				full := filepath.Join(root, rel)
				if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(full, []byte(body), 0644); err != nil {
					t.Fatal(err)
				}
			}
			changes, err := Write(root, p)
			if err != nil || len(changes.Advice) != 1 || !slices.Contains(changes.Written, "AGENTS.md") || slices.Contains(changes.Written, test.file) == test.shared {
				t.Fatal(changes, err)
			}
		})
	}
}

func TestUnsafeInstructionSelectionAbortsBuild(t *testing.T) {
	for _, kind := range []string{"dangling", "outside", "cycle", "nested", "malformed"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			p := testPlan(t, root)
			rel := "CLAUDE.md"
			target := filepath.Join(root, rel)
			var err error
			switch kind {
			case "dangling":
				err = os.Symlink("missing.md", target)
			case "outside":
				dest := filepath.Join(t.TempDir(), "shared.md")
				if err := os.WriteFile(dest, []byte("External."), 0644); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(dest, target)
			case "cycle":
				err = os.Symlink(rel, target)
			case "nested":
				rel = "specs/agent.md"
				err = os.WriteFile(filepath.Join(root, "specs/dotspec.yaml"), []byte("name: nested"), 0644)
			case "malformed":
				err = os.WriteFile(target, []byte(blockEnd), 0600)
			}
			if err != nil {
				t.Skip(err)
			}
			p.InstructionFiles = []string{"AGENTS.md", rel}
			if _, err := Write(root, p); err == nil {
				t.Fatal("accepted", kind)
			}
			for _, file := range []string{"AGENTS.md", projectIndex} {
				if _, err := os.Stat(filepath.Join(root, file)); !errors.Is(err, fs.ErrNotExist) {
					t.Fatal("partial build", file, err)
				}
			}
		})
	}
}

func TestGitSymlinkStub(t *testing.T) {
	root := t.TempDir()
	p := testPlan(t, root)
	p.InstructionFiles = nil
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatal(string(out), err)
		}
	}
	git("init", "-q")
	stub := filepath.Join(root, "CLAUDE.md")
	if err := os.Symlink("AGENTS.md", stub); err != nil {
		t.Skip(err)
	}
	git("add", "CLAUDE.md")
	if err := os.Remove(stub); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stub, []byte("AGENTS.md"), 0644); err != nil {
		t.Fatal(err)
	}
	changes, err := Write(root, p)
	if err != nil || len(changes.Advice) != 1 || slices.Contains(changes.Written, "CLAUDE.md") {
		t.Fatal(changes, err)
	}
}
