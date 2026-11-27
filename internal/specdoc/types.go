package specdoc

import "github.com/craigsmitham/dotspec/internal/specfile"

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
	tag := specfile.AuthoredType(kindName(h.subjectKind)) + "-report"
	if h.subjectKind == "SystemQualityAttribute" {
		tag = "system-quality-attribute-spec-report"
	}
	if h.subjectKind == "QualityModelAttribute" || h.subjectKind == "AdoptedProvision" {
		tag = kindName(h.subjectKind) + "-report"
	}
	return Document{tag, relativeURL(file, sf.path, h.anchor)}
}
