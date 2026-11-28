package specdoc

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/dotspec-md/dotspec/internal/model"
)

// Index declarations beside Markdown specifications. A catalog entry leads to
// the owning subject, where its native source and verification can be read together.
func (p *planner) indexSpecifications() {
	p.catalogReferences = map[string]model.Record{}
	for _, ref := range p.snapshot.References {
		p.catalogReferences[id(ref)] = ref
		resolution := model.Object(ref["resolution"])
		if resolution["rule"] == "native-definition/1" && resolution["status"] == "resolved" {
			doc := p.sources[attr(ref, "source_document")]
			u, err := url.Parse(attr(ref, "authored_target"))
			if doc != nil && err == nil {
				target := path.Join(path.Dir(doc.path), u.Path)
				p.targets[sourceKey(target, u.Fragment)] = target + "#L" + strconv.Itoa(model.Number(resolution["line"]))
			}
		}
	}
	for _, s := range p.subjects {
		if s.name.level != 1 || s.kind() == "QualityModel" {
			continue
		}
		system := s.system
		if s.kind() == "System" {
			system = s
		}
		if system == nil {
			system = p.governing
		}
		for system.navigationParent != nil && system.navigationParent.kind() == "System" {
			system = system.navigationParent
		}
		system.index("Specifications").add(entry{title: s.name.title, target: s.file, note: s.collectionType()})
		for _, key := range model.Strings(s.r["specifications"]) {
			r := p.elements[key]
			title := attr(p.nodes[nameNode(r)], "decoded_text")
			anchor := attr(r, "anchor")
			dest := s.file + "#" + anchor
			p.targets[key] = dest
			p.targets[sourceKey(s.source.path, anchor)] = dest
			category := sentenceCase(strings.ReplaceAll(attr(r, "defines"), "-", " "))
			if s.externalSystem {
				category += " for other systems"
			}
			member := entry{title: title, target: dest, note: s.name.title}
			owner := s
			if !owner.isOverview() {
				owner = s.navigationParent
			}
			owner.includeMember(category, member)
			system.index("Specifications").add(member)
		}
	}
}

// Named native definitions resolve to their declaring line on code hosts.
func (p *planner) catalogLink(w *writer, s *subject, r model.Record) string {
	ref := p.catalogReferences[attr(r, "source")]
	raw := attr(ref, "authored_target")
	return link(inlineText(attr(p.nodes[nameNode(r)], "decoded_text")), p.documentLink(w.file, s.source.path, raw))
}

func (p *planner) renderSpecifications(w *writer, s *subject) {
	ids := model.Strings(s.r["specifications"])
	if len(ids) == 0 {
		return
	}
	w.heading("Specification sources")
	for _, key := range ids {
		r := p.elements[key]
		w.text(`<a id="` + attr(r, "anchor") + `"></a>`)
		w.text("### " + p.catalogLink(w, s, r))
		w.text("Defines " + strings.ReplaceAll(attr(r, "defines"), "-", " ") + ".")
		p.renderAssessment(w, s, r, "#### Verification")
	}
}
