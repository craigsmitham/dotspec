package specdoc

import (
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/model"
)

// Include a referenced concept even when its definition is authored outside the project.
func (p *planner) indexConceptDefinitions() {
	for _, use := range p.subjects {
		if use.kind() != "ConceptUse" {
			continue
		}
		key := attr(use.r, "definition")
		if key == "" {
			continue
		}
		var member entry
		if def := p.pages[key]; def != nil {
			member = entry{def.name.title, def.file, subjectPreview(def)}
		} else {
			r := p.elements[key]
			for _, src := range p.snapshot.Sources {
				if id(src) != attr(r, "source_document") {
					continue
				}
				rel, err := filepath.Rel(p.root, filepath.FromSlash(attr(src, "locator")))
				if err != nil {
					continue
				}
				member = entry{attr(p.nodes[nameNode(r)], "decoded_text"), filepath.ToSlash(rel), p.conceptPurpose(r, attr(src, "markdown"))}
				break
			}
		}
		if member.target != "" {
			use.navigationParent.includeMember("Concepts", member)
		}
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
