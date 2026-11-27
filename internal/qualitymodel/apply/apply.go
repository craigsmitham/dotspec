// Package apply copies reusable attributes into independent system declarations.
package apply

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
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

type Options struct {
	Attributes []string
	All        bool
	DryRun     bool
	// Revision is the dotspec source commit or release version that citations
	// of a bundled model link to. Empty cites the embedded bytes by digest.
	Revision string
}

type Entry struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Decision string `json:"decision"`
	Path     string `json:"path"`
	Content  string `json:"content,omitempty"`
}

type Result struct {
	Entries       []Entry  `json:"entries"`
	SchemaVersion int      `json:"schema_version"`
	Target        string   `json:"target"`
	Added         []string `json:"added"`
	DryRun        bool     `json:"dry_run"`
}

// Run stages independent attribute sources. The first declared parent supplies
// the directory. All parents are explicit in the new source.
func Run(source *qualitymodel.QualityModel, target string, opts Options) (Result, error) {
	result := Result{SchemaVersion: 3, DryRun: opts.DryRun, Entries: []Entry{}, Added: []string{}}
	dir, err := specfile.Directory(filepath.Dir(target))
	if err != nil {
		return result, err
	}
	target = filepath.Join(dir, filepath.Base(target))
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
	local := map[string]system.Attribute{}
	for _, a := range current.QualityAttributes {
		local[a.ID] = a
	}
	byID := map[string]system.Attribute{}
	for _, a := range source.QualityAttributes {
		byID[a.ID] = a
	}
	chosen, requested, err := selection(source, byID, local, opts)
	if err != nil {
		return result, err
	}
	root := filepath.Join(owner, specfile.Folder, "quality")
	var folder func(system.Attribute) string
	folder = func(a system.Attribute) string {
		if existing, ok := local[a.ID]; ok {
			return filepath.Dir(existing.Location.Path)
		}
		parent := root
		if len(a.SubattributeOf) > 0 {
			parent = folder(byID[a.SubattributeOf[0]])
		}
		return filepath.Join(parent, specfile.NameSlug(a.Name))
	}
	combined := slices.Clone(current.QualityAttributes)
	used := map[string]bool{}
	relocated := map[string][]system.Attribute{}
	for _, a := range source.QualityAttributes {
		if !chosen[a.ID] {
			continue
		}
		role := "ancestor"
		if requested[a.ID] {
			role = "requested"
		}
		if existing, ok := local[a.ID]; ok {
			result.Entries = append(result.Entries, Entry{ID: a.ID, Role: role, Decision: "reuse", Path: existing.Location.Path})
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
		link, err := origin(source, file, opts.Revision, a)
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
		newAttribute, err := system.ParseQualityAttribute(file, []byte(body.String()))
		if err != nil {
			return result, err
		}
		combined = append(combined, newAttribute.Attribute)
		result.Entries = append(result.Entries, Entry{ID: a.ID, Role: role, Decision: "create", Path: file, Content: body.String()})
		result.Added = append(result.Added, a.ID)
	}
	if err := system.ValidateQualityAttributes(combined); err != nil {
		return result, err
	}
	additions := slices.DeleteFunc(slices.Clone(result.Entries), func(e Entry) bool { return e.Decision != "create" })
	for _, file := range additions {
		if err := validateDestination(owner, target, file.Path); err != nil {
			return result, err
		}
	}
	if len(additions) == 0 {
		return result, nil
	}
	if err := checkDiscoverable(owner, additions); err != nil {
		return result, err
	}
	if !opts.DryRun {
		err = publish(owner, additions, os.Link)
	}
	return result, err
}

// A successful adoption must be visible to the next system read, which uses
// Git's discovery rules. Check the complete batch before creating directories.
func checkDiscoverable(owner string, additions []Entry) error {
	root, err := specfile.RepositoryRoot(owner)
	if err != nil || root == "" {
		return err
	}
	var paths strings.Builder
	for _, f := range additions {
		rel, err := filepath.Rel(owner, f.Path)
		if err != nil {
			return err
		}
		paths.WriteString(filepath.ToSlash(rel))
		paths.WriteByte(0)
	}
	cmd := exec.Command("git", "-C", owner, "check-ignore", "--stdin", "-z")
	cmd.Stdin = strings.NewReader(paths.String())
	output, err := cmd.Output()
	var status *exec.ExitError
	if errors.As(err, &status) && status.ExitCode() == 1 {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check attribute destinations against Git ignore rules: %w", err)
	}
	first, _, _ := strings.Cut(string(output), "\x00")
	return fmt.Errorf("attribute destination is ignored by Git: %s", first)
}

// Check every existing ancestor during preview as well as execution.
func validateDestination(owner, target, path string) error {
	rel, err := filepath.Rel(owner, path)
	if err != nil || !filepath.IsLocal(rel) {
		return fmt.Errorf("destination outside target system: %s", path)
	}
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("attribute destination already exists: %s", path)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	for dir := filepath.Dir(path); dir != owner; dir = filepath.Dir(dir) {
		info, err := os.Lstat(dir)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("attribute destination requires a directory: %s", dir)
		}
		if err := specfile.CheckProjectBoundary(dir); err != nil {
			return err
		}
		decl, err := specfile.LocalDeclaration(dir)
		if err != nil {
			return err
		}
		if decl != "" && decl != target {
			return fmt.Errorf("destination belongs to %s", decl)
		}
	}
	return nil
}

// Stage complete files on the target filesystem. Linking a staged file into
// place never replaces existing content. Remove only our own additions on error.
func publish(owner string, files []Entry, install func(string, string) error) (err error) {
	stage, err := os.MkdirTemp(owner, ".dotspec-quality-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(stage) }()
	for i, f := range files {
		if err := os.WriteFile(filepath.Join(stage, strconv.Itoa(i)), []byte(f.Content), 0644); err != nil {
			return err
		}
	}
	var created []string
	defer func() {
		if err != nil {
			for _, path := range slices.Backward(created) {
				err = errors.Join(err, os.Remove(path))
			}
		}
	}()
	for i, f := range files {
		rel, err := filepath.Rel(owner, filepath.Dir(f.Path))
		if err != nil {
			return err
		}
		dir := owner
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			dir = filepath.Join(dir, part)
			if err := os.Mkdir(dir, 0755); err == nil {
				created = append(created, dir)
			} else if !errors.Is(err, fs.ErrExist) {
				return err
			}
			info, err := os.Lstat(dir)
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return fmt.Errorf("expected directory: %s", dir)
			}
			if dir != owner {
				if err := specfile.CheckProjectBoundary(dir); err != nil {
					return err
				}
			}
		}
		if err := install(filepath.Join(stage, strconv.Itoa(i)), f.Path); err != nil {
			return err
		}
		created = append(created, f.Path)
	}
	return nil
}

func selection(source *qualitymodel.QualityModel, byID, local map[string]system.Attribute, opts Options) (map[string]bool, map[string]bool, error) {
	selected := opts.Attributes
	if opts.All {
		selected = nil
		for _, a := range source.QualityAttributes {
			selected = append(selected, a.ID)
		}
	}
	requested := map[string]bool{}
	for _, id := range selected {
		if _, ok := byID[id]; !ok {
			return nil, nil, qualitymodel.UnknownAttribute(source, id)
		}
		requested[id] = true
	}
	chosen := map[string]bool{}
	queue := slices.Clone(selected)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if chosen[id] {
			continue
		}
		chosen[id] = true
		if _, ok := local[id]; !ok {
			queue = append(queue, byID[id].SubattributeOf...)
		}
	}
	return chosen, requested, nil
}

func bundledOrigin(source *qualitymodel.QualityModel, revision string, attribute system.Attribute) string {
	if revision == "" {
		return fmt.Sprintf("urn:dotspec:quality-model:%s:sha256:%x#%s", source.Slug, sha256.Sum256(source.Data), url.PathEscape(attribute.ID))
	}
	return "https://github.com/craigsmitham/dotspec/blob/" + revision + "/internal/qualitymodel/specs/quality-models/" + source.Slug + qualitymodel.Layout.Suffix + "#" + url.PathEscape(specfile.NameSlug(attribute.Name))
}

func origin(source *qualitymodel.QualityModel, target, revision string, attribute system.Attribute) (string, error) {
	if source.Builtin {
		return bundledOrigin(source, revision, attribute), nil
	}
	rel, err := filepath.Rel(filepath.Dir(target), source.Path)
	if err != nil {
		return "", err
	}
	return (&url.URL{Path: filepath.ToSlash(rel), Fragment: attribute.ID}).String(), nil
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
		close := specfile.LinkLabelEnd(raw, begin)
		if close < 0 {
			return nil, fmt.Errorf("cannot relocate source link %q", raw)
		}
		suffix := "(" + specfile.FormatDestination(dest)
		if title := model.String(n, "title"); title != "" {
			suffix += " " + specfile.FormatTitle(title)
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
