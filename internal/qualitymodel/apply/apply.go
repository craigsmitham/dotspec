// Package apply copies reusable attributes into independent system declarations.
package apply

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Result struct {
	SchemaVersion int      `json:"schema_version"`
	Source        string   `json:"source"`
	Target        string   `json:"target"`
	Added         []string `json:"added"`
	DryRun        bool     `json:"dry_run"`
	Content       string   `json:"content"`
}

// Run validates the complete result before replacing the file. Dry-run returns
// the same bytes that a real run would write, including the original BOM/EOLs.
func Run(source *qualitymodel.QualityModel, target string, selected []string, dryRun bool) (Result, error) {
	result := Result{SchemaVersion: 1, Source: source.Path, DryRun: dryRun, Added: []string{}}
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
		return result, fmt.Errorf("target must be a regular file: %s", target)
	}
	owner, err := system.Layout.Owner(target)
	if err != nil {
		return result, err
	}
	if err := system.Layout.ValidateLocal(owner); err != nil {
		return result, err
	}
	if filepath.Base(target) != system.Filename {
		return result, fmt.Errorf("target must be a %s file", system.Filename)
	}
	original, err := os.ReadFile(target)
	if err != nil {
		return result, err
	}
	current, err := system.Parse(original)
	if err != nil {
		return result, err
	}
	chosen, err := selection(source, selected)
	if err != nil {
		return result, err
	}
	for _, a := range current.QualityAttributes {
		if chosen[a.ID] {
			return result, fmt.Errorf("attribute ID %q already exists in %s; no changes made", a.ID, target)
		}
	}
	// Relocate Markdown links before slicing sections; reference-style links become
	// self-contained inline links so definitions elsewhere in the model aren't lost.
	relocated, err := relocate(source, target)
	if err != nil {
		return result, err
	}
	var meta struct {
		specfile.Metadata `yaml:",inline"`
		ID                string `yaml:"id"`
	}
	doc, err := specfile.ParseOutline(relocated, "quality-model", &meta)
	if err != nil {
		return result, err
	}
	var sections []*specfile.Section
	for _, s := range doc.Sections {
		if s.Title == "Quality attributes" {
			sections = s.Children
		}
	}
	var addition strings.Builder
	for i, a := range source.QualityAttributes {
		if !chosen[a.ID] {
			continue
		}
		if i >= len(sections) {
			return result, fmt.Errorf("source attribute sections disagree with parsed definitions")
		}
		result.Added = append(result.Added, a.ID)
		section := sections[i]
		addition.WriteString(strings.TrimSpace(string(doc.Source[section.Start:section.End])))
		sourceLink, err := origin(source, target, a.ID)
		if err != nil {
			return result, err
		}
		fmt.Fprintf(&addition, "\n\nSource: [quality model definition](<%s>).\n\n", sourceLink)
	}
	targetDoc, err := specfile.ParseOutline(original, "system", &specfile.Metadata{})
	if err != nil {
		return result, err
	}
	insert := len(original)
	heading := "## Quality attributes\n\n"
	for _, s := range targetDoc.Sections {
		if s.Title == "Quality attributes" {
			insert = rawOffset(original, len(targetDoc.Original)-len(targetDoc.Source)+s.End)
			heading = ""
		}
	}
	eol := "\n"
	if bytes.Contains(original, []byte("\r\n")) {
		eol = "\r\n"
	}
	added := strings.ReplaceAll("\n\n"+heading+strings.TrimRight(addition.String(), "\n")+"\n\n", "\n", eol)
	next := append([]byte{}, original[:insert]...)
	next = append(next, []byte(added)...)
	next = append(next, original[insert:]...)
	parsed, err := system.Parse(next)
	if err != nil {
		return result, fmt.Errorf("copied attributes would make the system invalid: %w", err)
	}
	if len(parsed.QualityAttributes) != len(current.QualityAttributes)+len(result.Added) {
		return result, fmt.Errorf("target Markdown prevents insertion of the selected attributes; no changes made")
	}
	for _, id := range result.Added {
		if parsed.FindQualityAttribute(id) == nil {
			return result, fmt.Errorf("copied attribute %s was not parsed; no changes made", id)
		}
	}
	result.Content = string(next)
	if !dryRun {
		if err := replace(target, original, next, info.Mode()); err != nil {
			return result, err
		}
	}
	return result, nil
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

func rawOffset(raw []byte, normalized int) int {
	i := 0
	if bytes.HasPrefix(raw, []byte("\xef\xbb\xbf")) {
		i = 3
	}
	for n := 0; n < normalized && i < len(raw); n++ {
		if raw[i] == '\r' && i+1 < len(raw) && raw[i+1] == '\n' {
			i++
		}
		i++
	}
	return i
}

func replace(path string, before, after []byte, mode os.FileMode) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".dotspec-quality-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer func() { _ = os.Remove(name) }()
	if err = f.Chmod(mode.Perm()); err == nil {
		_, err = f.Write(after)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	now, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || !bytes.Equal(now, before) {
		return fmt.Errorf("target changed during population; no changes made")
	}
	return os.Rename(name, path)
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
