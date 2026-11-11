package specfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

const WorkspaceFilename = "dotspec.yaml"

var projectName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// Workspace partitions discovery and implicit context; it is not a domain System.
type Workspace struct {
	Path     string    `json:"manifest"`
	Projects []Project `json:"projects"`
}

type Project struct {
	Name string `json:"name"`
	Root string `json:"root"`
}

func Within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// FindWorkspace stops at the Git/worktree boundary; outside Git it searches ancestors.
func FindWorkspace(dir string) (*Workspace, error) {
	dir, err := Directory(dir)
	if err != nil {
		return nil, err
	}
	repo, err := RepositoryRoot(dir)
	if err != nil {
		return nil, err
	}
	for current := dir; ; current = filepath.Dir(current) {
		path := filepath.Join(current, WorkspaceFilename)
		info, err := os.Lstat(path)
		if err == nil {
			if !info.Mode().IsRegular() {
				return nil, fmt.Errorf("%s must be a regular file, not a symlink", path)
			}
			return readWorkspace(path)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if current == repo || filepath.Dir(current) == current {
			return nil, nil
		}
	}
}

func readWorkspace(path string) (*Workspace, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config struct {
		Projects map[string]string `yaml:"projects"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("%s must contain exactly one YAML document", path)
	}
	if len(config.Projects) == 0 {
		return nil, fmt.Errorf("%s: projects must not be empty", path)
	}
	w := &Workspace{Path: path}
	base := filepath.Dir(path)
	seen := map[string]bool{}
	names := make([]string, 0, len(config.Projects))
	for name := range config.Projects {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		relative := config.Projects[name]
		if !projectName.MatchString(name) {
			return nil, fmt.Errorf("%s: invalid project name %q", path, name)
		}
		if relative == "" || filepath.IsAbs(relative) || strings.Contains(relative, "\\") {
			return nil, fmt.Errorf("%s: project %s needs a relative directory", path, name)
		}
		root := filepath.Clean(filepath.Join(base, filepath.FromSlash(relative)))
		if !Within(base, root) {
			return nil, fmt.Errorf("%s: project %s escapes the workspace", path, name)
		}
		if seen[root] {
			return nil, fmt.Errorf("%s: duplicate project root %s", path, relative)
		}
		// Roots must exist and may not cross symlinks, metadata folders, or nested repos.
		for current := root; current != base; current = filepath.Dir(current) {
			if filepath.Base(current) == Folder || filepath.Base(current) == ".git" {
				return nil, fmt.Errorf("%s: project root cannot be inside metadata: %s", path, relative)
			}
			info, err := os.Lstat(current)
			if err != nil {
				return nil, fmt.Errorf("%s: project %s: %w", path, name, err)
			}
			if !info.IsDir() {
				return nil, fmt.Errorf("%s: project root must use real directories: %s", path, relative)
			}
			if _, err := os.Lstat(filepath.Join(current, ".git")); err == nil {
				return nil, fmt.Errorf("%s: project %s crosses a nested repository", path, name)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
		}
		seen[root] = true
		w.Projects = append(w.Projects, Project{Name: name, Root: root})
	}
	sort.Slice(w.Projects, func(i, j int) bool { return w.Projects[i].Name < w.Projects[j].Name })
	return w, nil
}

func (w *Workspace) Containing(path string) (Project, bool) {
	var result Project
	for _, p := range w.Projects {
		if Within(p.Root, path) && len(p.Root) > len(result.Root) {
			result = p
		}
	}
	return result, result.Root != ""
}

// ProjectContext returns the declared partition. No manifest preserves legacy scope.
func ProjectContext(dir string) (*Workspace, Project, error) {
	dir, err := Directory(dir)
	if err != nil {
		return nil, Project{}, err
	}
	w, err := FindWorkspace(dir)
	if err != nil || w == nil {
		return w, Project{}, err
	}
	p, ok := w.Containing(dir)
	if !ok {
		return w, p, fmt.Errorf("%s is outside registered projects; select --project NAME or --all", dir)
	}
	return w, p, nil
}

// ContextRoot bounds implicit lookup only. Explicit links keep repository rules.
func ContextRoot(dir string) (string, error) {
	w, p, err := ProjectContext(dir)
	if err != nil {
		return "", err
	}
	if w != nil {
		return p.Root, nil
	}
	return RepositoryRoot(dir)
}
