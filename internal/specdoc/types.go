package specdoc

import "github.com/dotspec-md/dotspec/internal/specfile"

// Document identifies a logical page or passage. Documents is a portable URI
// to the exact authored subject.
type Document struct {
	Type      string `json:"type"`
	Documents string `json:"documents,omitempty"`
}

func subjectDocument(file string, sf *sourceFile, h *heading) Document {
	if h.subjectKind == "" {
		return Document{} // A physical file can collect several logical documents.
	}
	return Document{specfile.ReportType(h.subjectKind), relativeURL(file, sf.path, h.anchor)}
}
