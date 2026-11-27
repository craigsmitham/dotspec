// Package diagnostic defines findings shared by readers and check commands.
package diagnostic

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"syscall"
	"unicode/utf8"
)

type Location struct {
	Path    string `json:"path,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Element string `json:"element,omitempty"`
}

type Diagnostic struct {
	Rule     string     `json:"rule"`
	Severity string     `json:"severity"`
	Message  string     `json:"message"`
	Location Location   `json:"location"`
	Related  []Location `json:"related,omitempty"`
}

func (d Diagnostic) Error() string {
	message := "[" + d.Rule + "] " + d.Message
	if d.Location.Path != "" {
		message = d.Location.String() + ": " + message
	}
	for _, related := range d.Related {
		message += "\n  Related: " + related.String()
		if related.Element != "" {
			message += " (" + related.Element + ")"
		}
	}
	return message
}

func (l Location) String() string {
	s := l.Path
	if l.Line > 0 {
		s += fmt.Sprintf(":%d:%d", l.Line, l.Column)
	}
	return s
}

func Error(rule, message string, location Location, related ...Location) *Diagnostic {
	d := &Diagnostic{Rule: rule, Severity: "error", Message: message, Location: location}
	for _, location := range related {
		if location != (Location{}) {
			d.Related = append(d.Related, location)
		}
	}
	return d
}

func Advice(rule, message string, location Location) Diagnostic {
	return Diagnostic{Rule: rule, Severity: "advice", Message: message, Location: location}
}

// At uses one-based lines and Unicode-code-point columns in normalized source.
func At(source []byte, offset int) Location {
	if offset < 0 {
		offset = 0
	}
	if offset > len(source) {
		offset = len(source)
	}
	prefix := string(source[:offset])
	start := strings.LastIndexByte(prefix, '\n') + 1
	return Location{Line: strings.Count(prefix, "\n") + 1, Column: utf8.RuneCountInString(prefix[start:]) + 1}
}

// WithPath retains the underlying cause for I/O errors and adds source paths to
// validation findings, including both locations of a duplicate definition.
func WithPath(err error, path string) error {
	var d *Diagnostic
	if !errors.As(err, &d) {
		return fmt.Errorf("%s: %w", path, err)
	}
	copy := *d
	copy.Related = append([]Location(nil), d.Related...)
	if copy.Location.Path == "" {
		copy.Location.Path = path
	}
	for i := range copy.Related {
		if copy.Related[i].Path == "" {
			copy.Related[i].Path = path
		}
	}
	return &copy
}

// FileProblem states a filesystem failure in plain words. It omits the
// operating-system call that failed, and the path, which callers report as the
// finding's location or in front of the problem.
func FileProblem(err *fs.PathError) string {
	subject := "path"
	if kind, ok := strings.CutPrefix(err.Op, "find "); ok {
		subject = kind
	}
	switch {
	case errors.Is(err.Err, fs.ErrNotExist):
		return subject + " does not exist"
	case errors.Is(err.Err, fs.ErrPermission):
		return "permission denied"
	case errors.Is(err.Err, syscall.ENOTDIR):
		return "a parent of this path is not a directory"
	case errors.Is(err.Err, syscall.EISDIR):
		return "path is a directory"
	}
	return err.Err.Error()
}
