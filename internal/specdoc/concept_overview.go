package specdoc

import (
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/model"
)

// Each concept lists its named uses beneath it. A use without a definition is
// listed on its own. Include a referenced concept even when its definition is
// authored outside the project.
func (p *planner) indexConceptDefinitions() {
	for _, use := range p.subjects {
		if use.kind() != "ConceptUse" {
			continue
		}
		named := entry{title: use.name.title, target: use.file, preview: subjectPreview(use)}
		var member entry
		key := attr(use.r, "definition")
		switch def := p.pages[key]; {
		case key == "":
		case def != nil:
			member = entry{title: def.name.title, target: def.file, preview: subjectPreview(def)}
		default:
			r := p.elements[key]
			for _, src := range p.snapshot.Sources {
				if id(src) != attr(r, "source_document") {
					continue
				}
				rel, err := filepath.Rel(p.root, filepath.FromSlash(attr(src, "locator")))
				if err != nil {
					continue
				}
				member = entry{title: attr(p.nodes[nameNode(r)], "decoded_text"), target: filepath.ToSlash(rel), preview: p.conceptPurpose(r, attr(src, "markdown"))}
				break
			}
		}
		if member.target == "" {
			named.note = "no definition"
			use.navigationParent.includeMember(collectionTitle("ConceptUse"), named)
			continue
		}
		member.children = []entry{named}
		member.childLabel = "Uses:"
		use.navigationParent.includeMember(collectionTitle("ConceptUse"), member)
	}
}

func (p *planner) conceptPurpose(r model.Record, data string) string {
	for _, n := range model.Strings(model.Object(r["purpose"])["content_nodes"]) {
		node := p.nodes[n]
		if attr(node, "node_kind") == "paragraph" {
			start, end := position(node, "start"), position(node, "end")
			if start >= 0 && end >= start && end <= len(data) {
				return paragraphPreview(data[start:end])
			}
		}
	}
	return ""
}
