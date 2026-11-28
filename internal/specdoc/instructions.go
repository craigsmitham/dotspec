package specdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	blockStart = "<!-- dotspec:specification -->"
	blockEnd   = "<!-- /dotspec:specification -->"
)

// Resolve each existing component, including directory aliases. Missing ordinary
// components are allowed because explicit targets may need to be created.
func instructionPath(root, rel string) (string, error) {
	local, err := filepath.Localize(rel)
	if err != nil || !filepath.IsLocal(local) {
		return "", fmt.Errorf("invalid instruction path %q", rel)
	}
	current := root
	for _, part := range strings.Split(rel, "/") {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			current, err = filepath.EvalSymlinks(current)
			if err != nil {
				return "", fmt.Errorf("resolve %s: %w", rel, err)
			}
		}
		// Inspect the physical ancestors, including ones skipped by a symlink.
		local, err = filepath.Rel(root, current)
		if err != nil || !filepath.IsLocal(local) {
			return "", fmt.Errorf("%s resolves outside the project", rel)
		}
		ancestor := root
		for _, name := range strings.Split(filepath.ToSlash(local), "/") {
			if strings.EqualFold(name, ".git") {
				return "", fmt.Errorf("%s enters Git metadata", rel)
			}
			ancestor = filepath.Join(ancestor, name)
			if ancestor == root {
				continue
			}
			info, err := os.Stat(ancestor)
			if err != nil {
				return "", err
			}
			if info.IsDir() {
				boundary, err := nestedBoundary(ancestor)
				if err != nil {
					return "", err
				}
				if boundary {
					return "", fmt.Errorf("%s enters a nested project or repository", rel)
				}
			}
		}
	}
	info, err := os.Stat(current)
	if err == nil && !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", rel)
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	local, err = filepath.Rel(root, current)
	if err != nil {
		return "", err
	}
	local = filepath.ToSlash(local)
	if slices.ContainsFunc(strings.Split(local, "/"), func(s string) bool { return strings.EqualFold(s, ".git") }) {
		return "", fmt.Errorf("%s enters Git metadata", rel)
	}
	return local, nil
}

func instructionSelection(root string, configured []string) ([]string, error) {
	if configured != nil {
		return configured, nil
	}
	var selected []string
	for _, rel := range []string{"AGENTS.md", "CLAUDE.md", ".claude/CLAUDE.md"} {
		if _, err := os.Lstat(filepath.Join(root, rel)); err == nil {
			selected = append(selected, rel)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	if len(selected) == 0 {
		return []string{"AGENTS.md"}, nil
	}
	return selected, nil
}

func (i *inventory) selectedDestinations(root string, selected []string) (map[string]bool, error) {
	repo, err := specfile.RepositoryRoot(root)
	if err != nil {
		return nil, err
	}
	active := map[string]bool{}
	for _, rel := range selected {
		dest, err := instructionPath(root, rel)
		if err != nil {
			return nil, err
		}
		// At the physical destination, a tracked Git symlink must have been
		// checked out as a plain file, or the OS resolver would have followed it.
		if repo != "" {
			out, err := exec.Command("git", "-C", root, "--literal-pathspecs", "ls-files", "--stage", "--", dest).Output()
			if err != nil {
				return nil, fmt.Errorf("inspect Git mode of %s: %w", dest, err)
			}
			if bytes.HasPrefix(out, []byte("120000 ")) {
				i.changes.Advice = append(i.changes.Advice, diagnostic.Advice("spec.agent-instructions", "Git symlink checked out as plain text; restore the symlink to load the specification pointer", diagnostic.Location{Path: filepath.Join(root, dest)}))
				continue
			}
		}
		active[dest] = true
	}
	if slices.Contains(selected, "AGENTS.md") {
		if info, err := os.Stat(filepath.Join(root, "AGENTS.override.md")); err == nil && info.Mode().IsRegular() && info.Size() > 0 && !slices.Contains(selected, "AGENTS.override.md") {
			i.changes.Advice = append(i.changes.Advice, diagnostic.Advice("spec.agent-instructions", "Codex loads AGENTS.override.md ahead of AGENTS.md; select the override explicitly to include the pointer there", diagnostic.Location{Path: filepath.Join(root, "AGENTS.override.md")}))
		}
		agents, err := instructionPath(root, "AGENTS.md")
		if err != nil {
			return nil, err
		}
		for _, rel := range []string{"CLAUDE.md", ".claude/CLAUDE.md"} {
			if !slices.Contains(selected, rel) {
				continue
			}
			dest, err := instructionPath(root, rel)
			if err != nil {
				return nil, err
			}
			if dest == agents || !active[agents] {
				continue
			}
			data, err := os.ReadFile(filepath.Join(root, dest))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
			if importsAgents(data, path.Dir(rel)) {
				delete(active, dest)
			}
		}
	}
	return active, nil
}

func (i *inventory) planInstructions(root string, plan *Plan) error {
	selected, err := instructionSelection(root, plan.InstructionFiles)
	if err != nil {
		return err
	}
	active, err := i.selectedDestinations(root, selected)
	if err != nil {
		return err
	}
	// Store the last successful selection in the generated entry point. Reading
	// obsolete pages too supports changing the system name or project layout.
	targets := maps.Clone(active)
	for _, data := range i.files {
		header, _ := splitHeader(data)
		if header.GeneratedBy == "dotspec" {
			for _, rel := range header.InstructionFiles {
				targets[rel] = active[rel]
			}
		}
	}
	if len(active) > 0 {
		paths, err := json.Marshal(slices.Sorted(maps.Keys(active)))
		if err != nil {
			return err
		}
		i.candidate[plan.EntryPoint] = bytes.Replace(i.candidate[plan.EntryPoint], []byte("generated_by: dotspec\n"), []byte("generated_by: dotspec\nagent_instruction_files: "+string(paths)+"\n"), 1)
	}
	for _, rel := range slices.Sorted(maps.Keys(targets)) {
		dest, err := instructionPath(root, rel)
		if err != nil {
			return err
		}
		// A previously managed physical file may have been replaced by an alias.
		if !targets[rel] && dest != rel {
			continue
		}
		if plan.Pages[rel] != nil || slices.Contains(i.changes.Removed, rel) {
			return fmt.Errorf("%s is a generated report destination", rel)
		}
		file := filepath.Join(root, rel)
		data, err := os.ReadFile(file)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		label := "system specification"
		if plan.EntryPoint == projectIndex {
			label = "project specification index"
		}
		body := ""
		if targets[rel] {
			body = "Before changing or reviewing product behavior, read the " + link(label, relative(rel, plan.EntryPoint)) + " and follow the relevant requirement to its governing source."
		}
		next, err := updateBlock(data, body)
		if err != nil {
			i.collisions = append(i.collisions, *diagnostic.Error("spec.agent-instructions", err.Error(), diagnostic.Location{Path: file}))
			continue
		}
		if !bytes.Equal(data, next) {
			message := "specification pointer is stale"
			switch {
			case !targets[rel]:
				message = "obsolete specification pointer"
			case data == nil:
				message = "agent instruction file is missing"
			case !bytes.Contains(data, []byte(blockStart)):
				message = "specification pointer is missing"
			}
			i.instructions[rel] = instructionEdit{original: data, message: message}
			i.candidate[rel] = next
			i.changes.Written = append(i.changes.Written, rel)
		}
	}
	return nil
}

// A direct Claude import is a standalone paragraph line outside code examples.
func importsAgents(data []byte, dir string) bool {
	found := false
	doc := markdownParser.Parse(text.NewReader(data))
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindParagraph && n.Parent() == doc {
			for line := 0; line < n.Lines().Len(); line++ {
				seg := n.Lines().At(line)
				s := strings.TrimSpace(string(seg.Value(data)))
				if name, ok := strings.CutPrefix(s, "@"); ok && !path.IsAbs(name) && path.Join(dir, name) == "AGENTS.md" {
					found = true
				}
			}
		}
		return ast.WalkContinue, nil
	})
	return found
}

func blockRange(data []byte) (int, int, error) {
	start, end := -1, -1
	doc := markdownParser.Parse(text.NewReader(data))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHTMLBlock || n.Parent() != doc || n.Lines().Len() != 1 {
			return ast.WalkContinue, nil
		}
		seg := n.Lines().At(0)
		switch strings.TrimRight(string(seg.Value(data)), "\r\n") {
		case blockStart:
			if start < 0 {
				start = seg.Start
				return ast.WalkContinue, nil
			}
		case blockEnd:
			if start >= 0 && end < 0 {
				end = seg.Stop
				return ast.WalkContinue, nil
			}
		default:
			return ast.WalkContinue, nil
		}
		return ast.WalkStop, fmt.Errorf("ambiguous specification markers; keep exactly one ordered dotspec:specification pair")
	})
	if err != nil {
		return -1, -1, err
	}
	if start >= 0 && end < 0 {
		return -1, -1, fmt.Errorf("missing closing specification marker; add %s", blockEnd)
	}
	return start, end, nil
}

func updateBlock(data []byte, body string) ([]byte, error) {
	eol := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		eol = "\r\n"
	}
	start, end, err := blockRange(data)
	if err != nil {
		return nil, err
	}
	block := ""
	if body != "" {
		block = blockStart + eol + body + eol + blockEnd + eol
	}
	if start >= 0 {
		return slices.Concat(data[:start], []byte(block), data[end:]), nil
	}
	if body == "" {
		return data, nil
	}
	sep := ""
	if len(data) > 0 && !bytes.HasSuffix(data, []byte(eol+eol)) {
		sep = eol
		if !bytes.HasSuffix(data, []byte("\n")) {
			sep += eol
		}
	}
	next := slices.Concat(data, []byte(sep+block))
	// Appending inside an unterminated code fence or HTML comment would
	// hide the pointer and cause each subsequent build to append another.
	if begin, _, err := blockRange(next); err != nil || begin < 0 {
		return nil, fmt.Errorf("cannot append a visible specification block; close the open Markdown block first")
	}
	return next, nil
}
