package specfile

import (
	"os"
	"path/filepath"
)

// SourceSet owns the files observed by one inspection. Context resolution,
// link checks and export share the same bytes and parsed Markdown.
type SourceSet struct{ files map[string]*sourceFile }

type sourceFile struct {
	data   []byte
	err    error
	syntax *Markdown
}

func (s *SourceSet) file(path string) *sourceFile {
	if absolute, err := filepath.Abs(path); err == nil {
		path = absolute
	}
	if s.files == nil {
		s.files = map[string]*sourceFile{}
	}
	if f := s.files[path]; f != nil {
		return f
	}
	data, err := os.ReadFile(path)
	f := &sourceFile{data: data, err: err}
	s.files[path] = f
	return f
}

func (s *SourceSet) Read(path string) ([]byte, error) {
	f := s.file(path)
	return f.data, f.err
}

func (s *SourceSet) Markdown(path string) (*Markdown, error) {
	f := s.file(path)
	if f.err != nil {
		return nil, f.err
	}
	if f.syntax == nil {
		f.syntax = NewMarkdown(f.data)
	}
	return f.syntax, nil
}
