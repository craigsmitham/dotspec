// Package specdoc generates and checks reader-facing requirements specifications.
package specdoc

import "github.com/craigsmitham/dotspec/internal/specfile"

// Document identifies a logical page or passage. Documents is a portable URI
// to the exact authored subject. The report has project scope and no subject.
type Document struct {
	Type      string `yaml:"type,omitempty" json:"type"`
	Documents string `yaml:"documents,omitempty" json:"documents,omitempty"`
}

const RequirementsReviewReport = "requirements-review-report"

func subjectDocument(file string, sf *sourceFile, h *heading) Document {
	if h.subjectKind == "" {
		return Document{} // A physical file can collect several logical documents.
	}
	return Document{specfile.AuthoredType(kindName(h.subjectKind)) + "-report", relativeURL(file, sf.path, h.anchor)}
}
