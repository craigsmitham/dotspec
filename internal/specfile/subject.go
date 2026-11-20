package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// NameSlug converts the visible subject name into its filesystem identifier.
func NameSlug(name string) string {
	var b strings.Builder
	separator := false
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if separator {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			separator = false
		} else {
			separator = b.Len() > 0
		}
	}
	return b.String()
}

func (l Layout) CheckName(path, name string) error {
	slug := NameSlug(name)
	if slug == "" {
		return diagnostic.Error("layout.subject-name", "subject heading must include a letter or digit", diagnostic.Location{Path: path})
	}
	if filepath.Base(path) != l.filename(slug) {
		return diagnostic.Error("layout.subject-name", fmt.Sprintf("heading %q requires filename %s; rename the file or update its heading", name, l.filename(slug)), diagnostic.Location{Path: path})
	}
	return nil
}

// LocalSubject looks in both locations that claim dir. Direct declarations in
// a specs directory belong to its parent and are inspected at that parent.
func (l Layout) LocalSubject(dir string) (string, error) {
	var bases []string
	if !specsAlias(dir) {
		bases = append(bases, dir)
	}
	if metadata := filepath.Join(dir, Folder); specsAlias(metadata) {
		bases = append(bases, metadata)
	}
	var result string
	for _, base := range bases {
		info, err := os.Lstat(base)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			continue
		}
		entries, err := os.ReadDir(base)
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if !l.Match(e.Name()) {
				continue
			}
			file := filepath.Join(base, e.Name())
			if !e.IsDir() && !e.Type().IsRegular() {
				return "", diagnostic.Error("layout.subject-link", "subject declaration must be a regular file", diagnostic.Location{Path: file})
			}
			if result != "" {
				return "", diagnostic.Error("layout.subject-ambiguous", fmt.Sprintf("multiple %s declarations in %s: %s and %s", l.Kind, dir, result, file), diagnostic.Location{Path: file}, diagnostic.Location{Path: result})
			}
			result = file
		}
	}
	return result, nil
}

var SystemLayout = Layout{Kind: "system", Suffix: ".system.spec.md", NamedSubject: true}
var ComponentLayout = Layout{Kind: "component", Suffix: ".component.spec.md", NamedSubject: true}

func LocalDeclaration(dir string) (string, error) {
	sys, err := SystemLayout.LocalSubject(dir)
	if err != nil {
		return "", err
	}
	comp, err := ComponentLayout.LocalSubject(dir)
	if err != nil {
		return "", err
	}
	if sys != "" && comp != "" {
		return "", diagnostic.Error("layout.subject-ambiguous", "directory contains both a system and a component declaration", diagnostic.Location{Path: comp}, diagnostic.Location{Path: sys})
	}
	if sys != "" {
		return sys, nil
	}
	return comp, nil
}

func NearestDeclaration(dir, boundary string) (string, error) {
	for {
		file, err := LocalDeclaration(dir)
		if file != "" || err != nil || dir == boundary || filepath.Dir(dir) == dir {
			return file, err
		}
		dir = filepath.Dir(dir)
	}
}
