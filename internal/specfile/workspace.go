package specfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

const ProjectFilename = "dotspec.yaml"

var projectName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

type Workspace struct {
	Root     string
	Projects []Project
}

type Project struct {
	Name                  string      `json:"name" yaml:"name"`
	SpecificationPatterns []string    `json:"-" yaml:"specification_patterns"`
	Build                 BuildConfig `json:"-" yaml:"build"`
	Root                  string      `json:"root" yaml:"-"`
}

// BuildConfig controls whether the project has a generated specification.
type BuildConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ValidProjectName reports whether name can identify a project in dotspec.yaml.
func ValidProjectName(name string) bool { return projectName.MatchString(name) }

func (p Project) Manifest() string {
	return filepath.Join(p.Root, ProjectFilename)
}

// Only real files define boundaries. Never follow a symlink to a configuration.
func hasManifest(dir string) (bool, error) {
	info, err := os.Lstat(filepath.Join(dir, ProjectFilename))
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return err == nil && info.Mode().IsRegular(), err
}

func readProject(root string) (*Project, error) {
	p := &Project{Root: root, Build: BuildConfig{Enabled: true}}
	path := p.Manifest()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(p); err != nil {
		return nil, fmt.Errorf("%s: %s", path, manifestError(err))
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("%s must contain exactly one YAML document", path)
	}
	if !projectName.MatchString(p.Name) {
		return nil, fmt.Errorf("%s: invalid project name %q", path, p.Name)
	}
	for _, pattern := range p.SpecificationPatterns {
		if _, err := pathpkg.Match(pattern, ""); err != nil || pattern == "" || strings.HasPrefix(pattern, "/") {
			return nil, fmt.Errorf("%s: invalid specification pattern %q", path, pattern)
		}
	}
	return p, nil
}

// FindProject resolves local ownership without loading sibling or child manifests.
func FindProject(path string) (*Project, error) {
	dir, err := Directory(path)
	if err != nil {
		return nil, err
	}
	root, err := RepositoryRoot(dir)
	if err != nil {
		return nil, err
	}
	for current := dir; ; current = filepath.Dir(current) {
		found, err := hasManifest(current)
		if err != nil {
			return nil, err
		}
		if found {
			return readProject(current)
		}
		if current == root || current == filepath.Dir(current) {
			return nil, nil
		}
	}
}

// EnclosingProject returns the project that would contain dir if dir declared
// no project of its own. A repository root has none.
func EnclosingProject(dir string) (*Project, error) {
	dir, err := Directory(dir)
	if err != nil {
		return nil, err
	}
	repo, err := RepositoryRoot(dir)
	if err != nil || dir == repo || dir == filepath.Dir(dir) {
		return nil, err
	}
	return FindProject(filepath.Dir(dir))
}

// In Git the entire worktree is addressable. Outside Git, search from the
// outermost enclosing project, or from the supplied directory if it has none.
func FindWorkspace(dir string) (*Workspace, error) {
	dir, err := Directory(dir)
	if err != nil {
		return nil, err
	}
	repo, err := RepositoryRoot(dir)
	if err != nil {
		return nil, err
	}
	root := repo
	if root == "" {
		root = dir
		for current := dir; ; current = filepath.Dir(current) {
			found, err := hasManifest(current)
			if err != nil {
				return nil, err
			}
			if found {
				root = current
			}
			if current == filepath.Dir(current) {
				break
			}
		}
	}
	paths, err := manifestPaths(root, repo)
	if err != nil {
		return nil, err
	}
	w := &Workspace{Root: root}
	names := map[string]string{}
	for _, path := range paths {
		p, err := readProject(filepath.Dir(path))
		if err != nil {
			return nil, err
		}
		if other, ok := names[p.Name]; ok {
			return nil, fmt.Errorf("duplicate project name %q in %s and %s", p.Name, other, path)
		}
		names[p.Name] = path
		w.Projects = append(w.Projects, *p)
	}
	sort.Slice(w.Projects, func(i, j int) bool { return w.Projects[i].Name < w.Projects[j].Name })
	return w, nil
}

func manifestPaths(root, repo string) ([]string, error) {
	match := func(name string) bool { return name == ProjectFilename }
	skip := func(dir string) (bool, error) { return false, nil }
	if repo != "" {
		return gitPaths(root, repo, ":(glob)**/"+ProjectFilename, match, skip)
	}
	return walkSelectedPaths(root, match, skip)
}

// With no containing manifest, distinguish a workspace-only directory from
// an ordinary configuration-free repository. Direct reads also work without Git.
func HasProjects(dir string) (bool, error) {
	repo, err := RepositoryRoot(dir)
	if err != nil {
		return false, err
	}
	if repo != "" {
		dir = repo
		paths, err := manifestPaths(dir, repo)
		if err == nil {
			return len(paths) != 0, nil
		}
	}
	paths, err := manifestPaths(dir, "")
	return len(paths) != 0, err
}

func ContextRoot(dir string) (string, error) {
	p, err := FindProject(dir)
	if err != nil {
		return "", err
	}
	if p != nil {
		return p.Root, nil
	}
	return RepositoryRoot(dir)
}

var unknownManifestField = regexp.MustCompile(`line (\d+): field (\S+) not found in type \S+`)

// manifestError restates YAML decoding errors in manifest terms, so an unknown
// key names the key and the keys dotspec.yaml accepts instead of a Go type.
func manifestError(err error) string {
	var typeErr *yaml.TypeError
	if !errors.As(err, &typeErr) {
		return err.Error()
	}
	messages := make([]string, 0, len(typeErr.Errors))
	for _, message := range typeErr.Errors {
		if m := unknownManifestField.FindStringSubmatch(message); m != nil {
			message = fmt.Sprintf("line %s: unknown key %q; dotspec.yaml accepts name and build.enabled", m[1], m[2])
		}
		messages = append(messages, message)
	}
	return strings.Join(messages, "; ")
}
