package specdoc

import (
	"slices"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

func isDetailCollection(title string) bool {
	return title == "Goals" || title == "Assumptions"
}

// A named goal or assumption owns its qualifications on the collection page.
func (p *planner) collectOverviewDetails(s *subject) {
	var root *heading
	var idx *collectionIndex
	for _, h := range s.source.headings {
		if h.owner != s {
			continue
		}
		if root != nil && h.level <= root.level {
			root = nil
		}
		switch h.subjectKind {
		case "Goal", "Assumption":
			root = h
			idx = s.index(h.subjectKind + "s")
			idx.passages = append(idx.passages, nil)
		}
		if root != nil {
			h.detailEntry = root
			i := len(idx.passages) - 1
			idx.passages[i] = append(idx.passages[i], h)
		}
	}
}

func (p *planner) indexDetailHeadings(s *subject) {
	for title, idx := range s.indexes {
		if len(idx.passages) == 0 {
			continue
		}
		ids := parser.NewContext().IDs()
		ids.Generate([]byte(title), ast.KindHeading)
		// Rebasing a rich heading can change its automatic ID and the suffixes
		// of later duplicates. Earlier entries keep their natural fragments.
		richHeadings := false
		for _, passage := range idx.passages {
			for _, h := range passage {
				anchor := string(ids.Generate([]byte(h.title), ast.KindHeading))
				richHeadings = richHeadings || needsHeadingAnchor(s, h)
				if richHeadings || anchor != h.anchor {
					h.explicitAnchor = "detail." + h.anchor
					anchor = h.explicitAnchor
				}
				p.targets[h.node] = idx.file + "#" + anchor
			}
		}
		p.glossaryTargets[idx.file] = p.glossaryTargets[s.file]
	}
}

// A citation to a named entry leads back to the passage where the link occurs.
// When the model already supplies a relationship to its subject, that label
// describes the connection more precisely.
func (p *planner) indexDetailReferences() {
	entries := map[string]*heading{}
	for _, src := range p.sources {
		for _, h := range src.headings {
			if h.detailEntry != nil {
				entries[p.targets[h.node]] = h
			}
		}
	}
	for _, ref := range p.snapshot.References {
		h := entries[p.targets[attr(model.Object(ref["resolution"]), "target")]]
		src := p.sources[attr(ref, "source_document")]
		n := p.nodes[attr(ref, "content_node")]
		if h == nil || src == nil || attr(n, "node_kind") != "link" {
			continue
		}
		from := src.headingAt(position(n, "start"))
		if from == nil || from.owner == nil || from.detailEntry == h.detailEntry {
			continue
		}
		subject := p.targets[attr(ref, "context_owner")]
		relation, typed := relations[attr(ref, "role")]
		if typed && slices.ContainsFunc(p.routes[h.node], func(r route) bool {
			return r.label == relation.reverse && r.target == subject
		}) {
			continue
		}
		target := p.targets[from.node]
		title := from.owner.name.title
		if from != from.owner.name {
			title += " · " + from.title
		}
		if target != "" {
			p.routes[h.node] = append(p.routes[h.node], route{label: "Referenced by", title: title, target: target})
		}
	}
}
