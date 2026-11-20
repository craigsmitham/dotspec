package model

import (
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// folderClassification derives the directory edge only in inspected system
// context. Offline validation replays the same rule from source locators.
func folderClassification(r Record, file string, attrs map[string]Record) (string, Record) {
	folders := specfile.QualityFolders(file)
	if len(folders) == 0 {
		return "", nil
	}
	dir := filepath.Dir(file)
	field, owner := "quality_attributes", String(r, "specification_system")
	switch String(r, "kind") {
	case "Requirement":
	case "SystemQualityAttribute":
		if len(folders) < 2 {
			return "", nil
		}
		dir = filepath.Dir(dir)
		field, owner = "subattribute_of", String(r, "owner")
	default:
		return "", nil
	}
	a := attrs[filepath.ToSlash(specfile.AttributePath(dir))]
	if owner == "" || String(a, "owner") != owner {
		return "", nil
	}
	return field, a
}

func (b *snapshotBuilder) associateQualityTree() {
	attrs := map[string]Record{}
	for _, file := range sortedKeys(b.attributeOwners) {
		r := b.elements[b.rootByPath[file]]
		sys := b.elements[b.rootByPath[b.attributeOwners[file]]]
		if r == nil || sys == nil {
			continue
		}
		attrs[file] = r
		r["owner"] = sys["id"]
		sys["quality_attributes"] = appendUnique(Strings(sys["quality_attributes"]), String(r, "id"))
	}
	files := sortedKeys(b.rootByPath)
	for _, file := range files {
		r := b.elements[b.rootByPath[file]]
		if field, a := folderClassification(r, file, attrs); a != nil {
			r[field] = []string{String(a, "id")}
		}
	}
}
