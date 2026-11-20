// Package apply copies reusable attributes into independent system declarations.
package apply

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/rogpeppe/go-internal/lockedfile"
)

type Result struct {
	SchemaVersion int      `json:"schema_version"`
	Source        string   `json:"source"`
	Target        string   `json:"target"`
	Added         []string `json:"added"`
	DryRun        bool     `json:"dry_run"`
	Content       string   `json:"content"`
	Files         []File   `json:"files"`
}

type File struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// Run stages independent attribute sources. The first declared parent supplies
// the directory. All parents are explicit in the new source.
func Run(source *qualitymodel.QualityModel, target string, selected []string, dryRun bool) (Result, error) {
	result := Result{SchemaVersion: 2, Source: source.Path, DryRun: dryRun}
	target, err := filepath.Abs(target)
	if err != nil {
		return result, err
	}
	result.Target = target
	info, err := os.Lstat(target)
	if err != nil {
		return result, err
	}
	if !info.Mode().IsRegular() {
		return result, fmt.Errorf("target must be a regular system file: %s", target)
	}
	owner := system.Layout.ContextDir(target)
	lock := lockedfile.MutexAt(filepath.Join(os.TempDir(), fmt.Sprintf("dotspec-quality-%x.lock", sha256.Sum256([]byte(target)))))
	unlock, err := lock.Lock()
	if err != nil {
		return result, err
	}
	defer unlock()
	current, err := system.Read(target)
	if err != nil {
		return result, err
	}
	chosen, err := selection(source, selected)
	if err != nil {
		return result, err
	}
	byID := map[string]system.Attribute{}
	for _, a := range source.QualityAttributes {
		byID[a.ID] = a
	}
	for _, a := range current.QualityAttributes {
		if chosen[a.ID] {
			return result, fmt.Errorf("attribute ID %q already exists in %s", a.ID, target)
		}
	}
	root := filepath.Join(owner, specfile.Folder, "quality")
	var folder func(system.Attribute) string
	folder = func(a system.Attribute) string {
		parent := root
		if len(a.SubattributeOf) > 0 {
			parent = folder(byID[a.SubattributeOf[0]])
		}
		return filepath.Join(parent, specfile.NameSlug(a.Name))
	}
	used := map[string]bool{}
	relocated := map[string][]system.Attribute{}
	for _, a := range source.QualityAttributes {
		if !chosen[a.ID] {
			continue
		}
		dir := folder(a)
		file := specfile.AttributePath(dir)
		if used[file] {
			return result, fmt.Errorf("two attribute names produce the same destination: %s", file)
		}
		used[file] = true
		cache := dir
		if source.Builtin {
			cache = "builtin"
		}
		definitions := relocated[cache]
		if definitions == nil {
			data, err := relocate(source, file)
			if err != nil {
				return result, err
			}
			parsed, err := qualitymodel.Parse(data)
			if err != nil {
				return result, err
			}
			definitions = parsed.QualityAttributes
			relocated[cache] = definitions
		}
		i := slices.IndexFunc(definitions, func(d system.Attribute) bool { return d.ID == a.ID })
		if i < 0 {
			return result, fmt.Errorf("missing model attribute %s", a.ID)
		}
		link, err := origin(source, file, a.ID)
		if err != nil {
			return result, err
		}
		var body strings.Builder
		fmt.Fprintf(&body, "---\ntype: quality-attribute-spec\nid: %s\n", strconv.Quote(a.ID))
		if len(a.SubattributeOf) > 0 {
			body.WriteString("subattribute_of:\n")
			for _, id := range a.SubattributeOf {
				fmt.Fprintf(&body, "  - %s\n", strconv.Quote(id))
			}
		}
		fmt.Fprintf(&body, "---\n\n# %s\n\n## Definition\n\n%s\n\nSource: [quality model definition](<%s>).\n", a.Name, definitions[i].Definition, link)
		if _, err := system.ParseQualityAttribute(file, []byte(body.String())); err != nil {
			return result, err
		}
		result.Files = append(result.Files, File{file, body.String()})
		result.Added = append(result.Added, a.ID)
	}
	var preview strings.Builder
	for _, f := range result.Files {
		rel, _ := filepath.Rel(owner, f.Path)
		fmt.Fprintf(&preview, "-- %s --\n%s", filepath.ToSlash(rel), f.Content)
	}
	result.Content = preview.String()
	if !dryRun {
		err = publish(root, result.Files)
	}
	return result, err
}

// Publish complete top-level folders. Roll back only directories installed by
// this operation if a later rename fails.
func publish(root string, files []File) error {
	metadata := filepath.Dir(root)
	if err := os.Mkdir(metadata, 0755); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	for _, dir := range []string{metadata, root} {
		info, err := os.Lstat(dir)
		if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("%s must be a directory, not a file or symlink", dir)
			}
			if err := specfile.CheckProjectBoundary(dir); err != nil {
				return err
			}
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	stage, err := os.MkdirTemp(filepath.Dir(root), ".dotspec-quality-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(stage) }()
	for _, file := range files {
		rel, err := filepath.Rel(root, file.Path)
		if err != nil {
			return err
		}
		dest := filepath.Join(stage, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, []byte(file.Content), 0644); err != nil {
			return err
		}
	}
	_, err = os.Lstat(root)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.Chmod(stage, 0755); err != nil {
			return err
		}
		return os.Rename(stage, root)
	}
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(stage)
	if err != nil {
		return err
	}
	for _, e := range entries {
		dest := filepath.Join(root, e.Name())
		if _, err := os.Lstat(dest); err == nil {
			return fmt.Errorf("attribute folder already exists: %s", dest)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	var installed []string
	for _, e := range entries {
		from, to := filepath.Join(stage, e.Name()), filepath.Join(root, e.Name())
		if err := os.Rename(from, to); err != nil {
			for _, name := range installed {
				err = errors.Join(err, os.Rename(filepath.Join(root, name), filepath.Join(stage, name)))
			}
			return err
		}
		installed = append(installed, e.Name())
	}
	return nil
}

func selection(source *qualitymodel.QualityModel, selected []string) (map[string]bool, error) {
	byID := map[string]system.Attribute{}
	for _, a := range source.QualityAttributes {
		byID[a.ID] = a
	}
	if len(byID) == 0 {
		return nil, fmt.Errorf("quality model %q declares no attributes", source.Slug)
	}
	if len(selected) == 0 {
		for _, a := range source.QualityAttributes {
			selected = append(selected, a.ID)
		}
	}
	chosen := map[string]bool{}
	queue := append([]string(nil), selected...)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		a, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("unknown attribute ID %q in model %s", id, source.Slug)
		}
		if chosen[id] {
			continue
		}
		chosen[id] = true
		queue = append(queue, a.SubattributeOf...)
	}
	return chosen, nil
}

func origin(source *qualitymodel.QualityModel, target, id string) (string, error) {
	if source.Builtin {
		return "https://github.com/craigsmitham/dotspec/blob/main/internal/qualitymodel/specs/quality-models/" + source.Slug + qualitymodel.Suffix + "#" + url.PathEscape(id), nil
	}
	rel, err := filepath.Rel(filepath.Dir(target), source.Path)
	if err != nil {
		return "", err
	}
	return (&url.URL{Path: filepath.ToSlash(rel), Fragment: id}).String(), nil
}

// relocate rewrites only link suffix spans. Labels, formatting, code examples,
// and unrelated source bytes are untouched. Nested image/link suffixes are disjoint.
func relocate(source *qualitymodel.QualityModel, target string) ([]byte, error) {
	data := specfile.Normalize(source.Data)
	parsed := model.ParseContent(source.Path, data)
	type edit struct {
		start, end int
		text       string
	}
	edits := []edit{}
	for _, n := range parsed.Nodes {
		kind := model.String(n, "node_kind")
		if kind != "link" && kind != "image" {
			continue
		}
		dest, has := n["destination"].(string)
		if !has {
			continue
		}
		u, err := url.Parse(dest)
		if err != nil {
			return nil, fmt.Errorf("invalid source link %q: %w", dest, err)
		}
		if u.Scheme == "" && u.Host == "" {
			if source.Builtin {
				return nil, fmt.Errorf("built-in model contains a nonportable link %q", dest)
			}
			original := source.Path
			if u.Path != "" {
				original = filepath.Join(filepath.Dir(source.Path), filepath.FromSlash(u.Path))
				if filepath.IsAbs(filepath.FromSlash(u.Path)) {
					original = filepath.FromSlash(u.Path)
				}
			}
			rel, err := filepath.Rel(filepath.Dir(target), original)
			if err != nil {
				return nil, err
			}
			u.Path = filepath.ToSlash(rel)
			u.RawPath = ""
			dest = u.String()
		}
		span := model.Object(n["span"])
		start, end := model.Number(span["start"]), model.Number(span["end"])
		if start < 0 || end > len(data) || start >= end {
			return nil, fmt.Errorf("source link span is invalid")
		}
		raw := data[start:end]
		begin := 0
		if raw[0] == '!' {
			begin++
		}
		close := labelEnd(raw, begin)
		if close < 0 {
			return nil, fmt.Errorf("cannot relocate source link %q", raw)
		}
		suffix := "(<" + strings.NewReplacer("<", "%3C", ">", "%3E", "\n", "%0A", "\r", "%0D").Replace(dest) + ">"
		if title := model.String(n, "title"); title != "" {
			suffix += " \"" + strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(title) + "\""
		}
		suffix += ")"
		edits = append(edits, edit{start + close + 1, end, suffix})
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	limit := len(data)
	for _, e := range edits {
		if e.end > limit {
			return nil, fmt.Errorf("overlapping source links cannot be relocated")
		}
		data = append(append(append([]byte{}, data[:e.start]...), []byte(e.text)...), data[e.end:]...)
		limit = e.start
	}
	return data, nil
}

func labelEnd(raw []byte, start int) int {
	if start >= len(raw) || raw[start] != '[' {
		return -1
	}
	depth := 0
	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
