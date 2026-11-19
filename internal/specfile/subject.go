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

// LocalSubject finds the single declaration in a subject's own specs directory.
func (l Layout) LocalSubject(dir string) (string, error) {
	base := filepath.Join(dir, Folder)
	info, err := os.Lstat(base)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", diagnostic.Error("layout.directory", "expected a regular specs directory", diagnostic.Location{Path: base})
	}
	entries, err := os.ReadDir(base)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	var result string
	for _, e := range entries {
		if !l.Match(e.Name()) {
			continue
		}
		file := filepath.Join(base, e.Name())
		if e.Type()&os.ModeSymlink != 0 {
			return "", diagnostic.Error("layout.subject-link", "subject declaration must be a regular file, not a symlink", diagnostic.Location{Path: file})
		}
		if result != "" {
			return "", diagnostic.Error("layout.subject-ambiguous", fmt.Sprintf("multiple %s declarations in %s: %s and %s", l.Kind, base, filepath.Base(result), e.Name()), diagnostic.Location{Path: file})
		}
		result = file
	}
	return result, nil
}
