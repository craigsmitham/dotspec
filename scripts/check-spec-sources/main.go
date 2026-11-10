// Command check-spec-sources enforces this repository's specification-test
// convention using the public dotspec JSON inventory. It does not infer intent.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Keep this consumer's JSON contract independent of dotspec's implementation.
type inventory struct {
	SchemaVersion int  `json:"schema_version"`
	Complete      bool `json:"complete"`
	Scope         struct {
		Path string `json:"path"`
	} `json:"scope"`
	Files []struct {
		Path string `json:"path"`
	} `json:"files"`
	Diagnostics []struct {
		Severity string `json:"severity"`
	} `json:"diagnostics"`
}

func main() {
	count, err := run(".", runSources)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d specification-test files are discoverable through dotspec.\n", count)
}

func runSources(root string) ([]byte, error) {
	cmd := exec.Command("go", "run", ".", "sources", "--json")
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("dotspec sources failed: %w\n%s%s", err, stderr.String(), data)
	}
	return data, nil
}

func run(root string, load func(string) ([]byte, error)) (int, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return 0, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return 0, err
	}
	// These acceptance scripts define command-level examples. Other *_test.go
	// files are not opted in merely because they are tests.
	expected, err := filepath.Glob(filepath.Join(root, "internal/integration/testdata/*.txtar"))
	if err != nil {
		return 0, err
	}
	if len(expected) == 0 {
		return 0, fmt.Errorf("no acceptance scripts found; run from the dotspec repository root")
	}
	// This dedicated file verifies the prose quality requirements linked to it.
	expected = append(expected, filepath.Join(root, "internal/cli/quality_contracts_test.go"))
	for _, path := range expected {
		info, err := os.Lstat(path)
		if err != nil {
			return 0, err
		}
		if !info.Mode().IsRegular() {
			return 0, fmt.Errorf("expected a regular specification-test file: %s", path)
		}
	}
	data, err := load(root)
	if err != nil {
		return 0, err
	}
	var report inventory
	if err := json.Unmarshal(data, &report); err != nil {
		return 0, fmt.Errorf("decode dotspec sources: %w", err)
	}
	if report.SchemaVersion != 1 || !report.Complete || filepath.Clean(report.Scope.Path) != root {
		return 0, fmt.Errorf("dotspec sources did not provide a complete version-1 inventory of %s", root)
	}
	for _, d := range report.Diagnostics {
		if d.Severity == "error" {
			return 0, fmt.Errorf("dotspec sources reported errors; repair the source inventory first")
		}
	}
	found := make(map[string]bool)
	for _, file := range report.Files {
		if !filepath.IsAbs(file.Path) {
			return 0, fmt.Errorf("inventory source path is not absolute: %s", file.Path)
		}
		found[filepath.Clean(file.Path)] = true
	}
	var missing []string
	for _, path := range expected {
		if !found[path] {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return 0, err
			}
			missing = append(missing, filepath.ToSlash(rel))
		}
	}
	if len(missing) != 0 {
		sort.Strings(missing)
		return 0, fmt.Errorf("specification-test files missing from dotspec sources:\n  %s\nAdd ordinary file links from relevant specs; describe whether each file defines behavior or verifies an existing requirement", strings.Join(missing, "\n  "))
	}
	return len(expected), nil
}
