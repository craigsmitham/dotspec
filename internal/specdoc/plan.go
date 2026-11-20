package specdoc

import (
	"cmp"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Plan binds source-backed expectations to their project-relative generated paths.
type Plan struct{ Pages map[string]*Page }

type Page struct {
	Document
	Title, Markdown                                 string
	SectionDocs                                     map[string]Document
	DirectoryMembers                                map[string][]string
	RequiredLinks, RequiredContent, RequiredAnchors []string
	RequiredLabels, RequiredSectionLinks            map[string][]string
	RequiredSectionContent                          map[string]string
	RequiredResources                               []string
}

type sourceFile struct {
	path, data       string
	headings         []*heading
	links, footnotes []model.Record
}

// Each section's body ends at the next heading. Nested declarations can move
// together to a page while narrative and unfamiliar sections keep source order.
type heading struct {
	node, title, anchor, explicitAnchor string
	start, end, bodyEnd, level          int
	owner                               *subject
	obligation                          bool
	qualityAttribute                    bool
	inGlossary                          bool
	subjectKind                         string
}

type subject struct {
	r         model.Record
	source    *sourceFile
	name      *heading
	system    *subject
	parent    *subject
	indexes   map[string]*collectionIndex
	glossary  []glossaryPassage
	overview  []overviewSection
	file, dir string
}

type planner struct {
	root            string
	snapshot        *model.Snapshot
	plan            *Plan
	sources         map[string]*sourceFile
	nodes, elements map[string]model.Record
	pages           map[string]*subject
	subjects        []*subject
	targets         map[string]string // snapshot record, source, or fragment to canonical project-relative URL
	governing       *subject

	glossaryTargets map[string]map[string]string // scoped routes for terms shared between systems
	linkedPassages  map[string]bool
	routes          map[string][]route
}

type entry struct {
	title, target, preview string
}

var families = map[string]string{
	"SystemQualityAttribute": "quality", "System": "systems", "Component": "components", "Feature": "features", "Capability": "capabilities",
	"ConceptUse": "concept-uses", "Interface": "interfaces", "Synchronization": "synchronizations",
	"QualityModel": "quality-models", "Requirement": "requirements", "BusinessRule": "business-rules",
	"Job": "jobs", "UseCase": "use-cases", "Concept": "concepts",
}

func Generate(root string, snapshot *model.Snapshot) (*Plan, error) {
	// Export resolves directory aliases (e.g. /var versus /private/var).
	// All source locators are relative to that same canonical root.
	if roots := model.Strings(snapshot.Scope["roots"]); len(roots) == 1 {
		root = filepath.FromSlash(roots[0])
	}
	p := &planner{root: root, snapshot: snapshot, plan: &Plan{Pages: map[string]*Page{}}, sources: map[string]*sourceFile{}, nodes: map[string]model.Record{}, elements: map[string]model.Record{}, pages: map[string]*subject{}, targets: map[string]string{}, glossaryTargets: map[string]map[string]string{}}
	if err := p.read(); err != nil {
		return nil, err
	}
	if err := p.place(); err != nil {
		return nil, err
	}
	p.indexTargets()
	plan := p.render()
	return plan, addSourceContracts(p)
}

func id(r model.Record) string                { return model.String(r, "id") }
func attr(r model.Record, key string) string  { return model.String(r, key) }
func position(r model.Record, key string) int { return model.Number(model.Object(r["span"])[key]) }
func nameNode(r model.Record) string {
	nodes := model.Strings(model.Object(r["name"])["content_nodes"])
	if len(nodes) == 0 {
		return ""
	}
	return nodes[0]
}
func (p *planner) read() error {
	glossarySources := map[string]bool{}
	for _, doc := range p.snapshot.Documents {
		p.elements[id(doc)] = doc
		if doc["glossary_terms"] != nil {
			glossarySources[attr(doc, "source")] = true
		}
	}
	for _, r := range p.snapshot.Sources {
		if selected, _ := r["selected"].(bool); !selected && !glossarySources[id(r)] {
			continue
		}
		data, err := base64.StdEncoding.DecodeString(attr(model.Object(r["payload"]), "bytes_base64"))
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(p.root, filepath.FromSlash(attr(r, "locator")))
		if err != nil {
			return err
		}
		if glossarySources[id(r)] && (rel == "SPEC.md" || strings.HasPrefix(filepath.ToSlash(rel), "specification/")) {
			if header, _ := splitHeader(data); header.GeneratedBy == "dotspec" {
				return fmt.Errorf("glossary source %s is an output of this report", rel)
			}
		}
		p.sources[id(r)] = &sourceFile{path: filepath.ToSlash(rel), data: string(data)}
	}
	headingByNode := map[string]*heading{}
	// Only top-level headings delimit domain sections. Quoted or listed headings
	// are authored body content and must remain inside their original block.
	top := map[string]bool{}
	roots := map[string]bool{}
	for _, d := range p.snapshot.Documents {
		roots[attr(d, "root")] = true
	}
	for _, n := range p.snapshot.ContentNodes {
		if roots[id(n)] {
			for _, c := range model.Strings(n["children"]) {
				top[c] = true
			}
		}
	}
	for _, n := range p.snapshot.ContentNodes {
		p.nodes[id(n)] = n
		sf := p.sources[attr(model.Object(n["span"]), "source")]
		if sf == nil {
			continue
		}
		switch attr(n, "node_kind") {
		case "heading":
			if !top[id(n)] {
				continue
			}
			h := &heading{node: id(n), title: attr(n, "decoded_text"), level: model.Number(n["heading_level"]), start: position(n, "start"), end: position(n, "end")}
			sf.headings = append(sf.headings, h)
			headingByNode[id(n)] = h
		case "footnote_definition", "footnote_reference":
			sf.footnotes = append(sf.footnotes, n)
		case "link", "image", "link_definition":
			sf.links = append(sf.links, n)
		}
	}
	for _, f := range p.snapshot.Fragments {
		if h := headingByNode[attr(f, "content_node")]; h != nil && attr(f, "generated_anchor") != "" {
			h.anchor = attr(f, "generated_anchor")
		}
	}
	for _, sf := range p.sources {
		slices.SortFunc(sf.headings, func(a, b *heading) int { return cmp.Compare(a.start, b.start) })
		for i, h := range sf.headings {
			h.bodyEnd = len(sf.data)
			if i+1 < len(sf.headings) {
				h.bodyEnd = sf.headings[i+1].start
			}
		}
		slices.SortFunc(sf.links, func(a, b model.Record) int {
			if n := cmp.Compare(position(a, "start"), position(b, "start")); n != 0 {
				return n
			}
			return cmp.Compare(position(b, "end"), position(a, "end"))
		})
	}
	for _, r := range p.snapshot.Elements {
		p.elements[id(r)] = r
		sf := p.sources[attr(r, "source_document")]
		h := headingByNode[nameNode(r)]
		if h != nil {
			h.subjectKind = attr(r, "kind")
			switch attr(r, "kind") {
			case "ObligationClause":
				h.obligation = true
			case "QualityAttribute", "SystemQualityAttribute":
				h.qualityAttribute = true
			}
		}
		if sf == nil || families[attr(r, "kind")] == "" || h == nil {
			continue
		}
		s := &subject{r: r, source: sf, name: h}
		h.owner = s
		p.pages[id(r)] = s
		p.subjects = append(p.subjects, s)
	}
	// Every selected source is a discovered authored specification. Its H1 must
	// have exactly one page even when other headings also get their own pages.
	rootsBySource := make(map[*sourceFile]int)
	for _, subject := range p.subjects {
		if subject.name.level == 1 {
			rootsBySource[subject.source]++
		}
	}
	for _, key := range slices.Sorted(maps.Keys(p.sources)) {
		if glossarySources[key] {
			continue
		}
		source := p.sources[key]
		if count := rootsBySource[source]; count != 1 {
			return fmt.Errorf("%s: expected one standalone report for the source title, found %d", source.path, count)
		}
	}
	slices.SortFunc(p.subjects, func(a, b *subject) int {
		if n := strings.Compare(a.source.path, b.source.path); n != 0 {
			return n
		}
		return cmp.Compare(a.name.start, b.name.start)
	})
	return nil
}

func (s *subject) slug() string {
	if s.kind() == "System" || s.kind() == "Component" || s.kind() == "SystemQualityAttribute" {
		return strings.TrimSuffix(path.Base(s.source.path), "."+kindName(s.kind())+".spec.md")
	}
	name := attr(s.r, "authored_id")
	if name == "" {
		name = s.name.title
	}
	if slug := specfile.NameSlug(name); slug != "" {
		return slug
	}
	return "section"
}
func (s *subject) kind() string {
	return attr(s.r, "kind")
}

// A system or component owns a directory. Its overview mirrors its source filename,
// including when a sibling with the same name needs a distinct directory.
func (s *subject) placeOverview(dir string) {
	s.dir = dir
	s.file = path.Join(dir, strings.TrimSuffix(path.Base(s.source.path), ".md")+"-report.md")
}

func (s *subject) collectionDir() string {
	if s.kind() == "System" {
		return path.Join(s.parent.dir, "systems")
	}
	return path.Join(s.parent.dir, families[s.kind()])
}

func (p *planner) place() error {
	var roots []string
	for _, s := range p.subjects {
		parent := attr(s.r, "specification_system")
		switch s.kind() {
		case "System":
			parent = attr(s.r, "part_of")
			if parent == "" {
				p.governing = s
				roots = append(roots, s.source.path)
				continue
			}
		case "QualityModel":
			continue // project-wide reference definition
		}
		switch s.kind() {
		case "Component":
			s.parent = p.pages[attr(s.r, "part_of")]
		case "System":
			s.parent = p.pages[parent]
		default:
			s.parent = p.pages[attr(s.r, "specification_owner")]
		}
		if s.parent == nil {
			s.parent = p.pages[attr(s.r, "owner")]
		}
		if parent == "" {
			parent = attr(s.r, "owner")
		}
		for parent != "" && attr(p.elements[parent], "kind") != "System" {
			element := p.elements[parent]
			if attr(element, "kind") == "Component" {
				parent = attr(element, "specification_system")
			} else {
				parent = attr(element, "owner")
			}
		}
		s.system = p.pages[parent]
		if s.parent == nil {
			s.parent = s.system
		}
		if s.system == nil {
			return fmt.Errorf("%s: cannot place %s in a selected system", s.source.path, s.name.title)
		}
	}
	if len(roots) != 1 {
		return fmt.Errorf("generated specification requires one root system; found %d (%s)", len(roots), strings.Join(roots, ", "))
	}
	p.governing.placeOverview("specification")
	for _, s := range p.subjects {
		if s.kind() == "QualityModel" {
			s.system = p.governing
			s.parent = p.governing
		}
	}
	// Disambiguate siblings by stable project-relative source and authored heading.
	names := map[string]int{}
	key := func(s *subject) string { return id(s.parent.r) + "/" + families[s.kind()] + "/" + s.slug() }
	for _, s := range p.subjects {
		if s != p.governing {
			names[key(s)]++
		}
	}
	active := map[*subject]bool{}
	var assign func(*subject) error
	assign = func(s *subject) error {
		if s.file != "" {
			return nil
		}
		if active[s] {
			return fmt.Errorf("cycle in system parentage at %s", s.source.path)
		}
		active[s] = true
		if err := assign(s.parent); err != nil {
			return err
		}
		slug := s.slug()
		if names[key(s)] > 1 {
			sum := sha256.Sum256([]byte(s.source.path + "#" + s.name.anchor))
			slug += fmt.Sprintf("-%x", sum[:5])
		}
		base := path.Join(s.collectionDir(), slug)
		if s.kind() == "SystemQualityAttribute" {
			dir := path.Join(s.system.dir, "quality")
			if parents := model.Strings(s.r["subattribute_of"]); len(parents) > 0 {
				if parent := p.pages[parents[0]]; parent != nil {
					if err := assign(parent); err != nil {
						return err
					}
					dir = parent.dir
				}
			}
			s.dir = path.Join(dir, slug)
			s.file = path.Join(s.dir, slug+".quality-attribute.spec-report.md")
		} else if s.kind() == "System" || s.kind() == "Component" {
			s.placeOverview(base)
		} else {
			s.file = base + "." + kindName(s.kind()) + ".spec-report.md"
		}
		return nil
	}
	files := map[string]bool{}
	for _, s := range p.subjects {
		if err := assign(s); err != nil {
			return err
		}
		if files[s.file] {
			return fmt.Errorf("generated address collision at %s", s.file)
		}
		files[s.file] = true
	}
	return nil
}

func needsHeadingAnchor(s *subject, h *heading) bool {
	return strings.IndexFunc(h.title, func(r rune) bool { return r > 127 }) >= 0 || strings.ContainsAny(s.source.data[h.start:h.end], "[]*_~`<>&")
}

func (p *planner) indexTargets() {
	for _, s := range p.subjects {
		p.targets[id(s.r)] = s.file
		if s.name.level == 1 {
			p.targets[attr(s.r, "source_document")] = s.file
		}
	}
	// A source heading belongs to the nearest lifted declaration. Preserve the
	// original fragment identity when a section changes page or heading level.
	for _, sf := range p.sources {
		var stack []*subject

		for _, h := range sf.headings {
			for len(stack) > 0 && stack[len(stack)-1].name.level >= h.level {
				stack = stack[:len(stack)-1]
			}
			if h.owner != nil {
				stack = append(stack, h.owner)
			}
			if len(stack) == 0 {
				continue
			}
			h.owner = stack[len(stack)-1]
		}
	}

	for _, s := range p.subjects {
		if s != p.governing && !s.name.qualityAttribute {
			p.placeDirectoryEntry(s)
		}
		if s.kind() == "System" && len(model.Strings(s.r["quality_attributes"])) > 0 {
			s.index("Quality attributes")
		}
	}
	for _, rule := range p.subjects {
		if rule.kind() != "BusinessRule" {
			continue
		}
		for _, target := range model.Strings(rule.r["applies_to"]) {
			sys := p.pages[target]
			if sys == nil || sys.kind() != "System" {
				continue
			}
			idx := sys.index("Business rules")
			idx.groups[""] = slices.DeleteFunc(idx.groups[""], func(e entry) bool { return e.target == rule.file })
			if len(idx.groups[""]) == 0 {
				delete(idx.groups, "")
			}
			idx.add(applicableRules, entry{rule.name.title, rule.file, subjectPreview(rule)})
		}
	}
	// Embedded obligations get their final fragments after the overview headings are planned.
	clauses := map[string]bool{}
	for _, r := range p.snapshot.Elements {
		if attr(r, "kind") == "ObligationClause" {
			page := p.ownerPage(r)
			if page != nil {
				sys := page.system
				if page.parent != nil && page.parent.kind() == "Component" {
					sys = page.parent
				}
				if page.kind() == "System" || page.kind() == "Component" {
					sys = page
				}
				n := p.nodes[nameNode(r)]
				sys.index("Requirements").add("", entry{attr(n, "decoded_text"), id(r), p.requirementPreview(r)})
				clauses[id(r)] = true
			}
		}
	}
	for _, r := range p.snapshot.Elements {
		if attr(r, "kind") != "Requirement" && attr(r, "kind") != "ObligationClause" {
			continue
		}
		for _, target := range model.Strings(r["applies_to"]) {
			component := p.pages[target]
			dest := p.targets[id(r)]
			if clauses[id(r)] {
				dest = id(r)
			}
			if component != nil && component.kind() == "Component" && dest != "" {
				index := component.index("Requirements")
				if !index.contains("", dest) {
					index.add("Applicable requirements", entry{attr(p.nodes[nameNode(r)], "decoded_text"), dest, p.requirementPreview(r)})
				}
			}
		}
	}
	// A system index also leads to the same collection in its components.
	// Walk to the enclosing system so nested components remain discoverable
	// regardless of the order in which their sources were read.
	for _, s := range p.subjects {
		if s.kind() != "Component" {
			continue
		}
		for title := range s.indexes {
			if title == "Components" {
				continue
			}
			for child := s; child.kind() == "Component"; child = child.parent {
				child.parent.index(title).add("By component", entry{child.name.title, child.indexes[title].file, subjectPreview(child)})
			}
		}
	}
	for _, s := range p.subjects {
		ids := parser.NewContext().IDs()
		if s.kind() == "System" {
			p.planGlossary(s)
			s.overview = p.planOverview(s)
			p.indexOverviewHeadings(s, ids)
			continue
		}
		for _, h := range readingOrder(s) {
			if h.owner != s {
				continue
			}
			p.targets[h.node] = s.file
			generated := string(ids.Generate([]byte(h.title), ast.KindHeading))
			if h != s.name {
				if needsHeadingAnchor(s, h) {
					h.explicitAnchor = "section-" + h.anchor
					generated = h.explicitAnchor
				}
				p.targets[h.node] += "#" + generated
			}
		}
	}
	for _, r := range p.snapshot.Elements {
		if dest := p.targets[nameNode(r)]; dest != "" {
			p.targets[id(r)] = dest
		}
	}
	for _, s := range p.subjects {
		if idx := s.indexes["Requirements"]; idx != nil {
			for _, members := range idx.groups {
				for i, member := range members {
					if clauses[member.target] {
						members[i].target = p.targets[member.target]
					}
				}
			}
		}
	}
	// Map source paths and their authored or generated fragments as link targets.
	for _, f := range p.snapshot.Fragments {
		node := attr(f, "content_node")
		dest := p.targets[node]
		if dest == "" {
			for _, r := range p.snapshot.Elements {
				if attr(r, "source_fragment") == id(f) {
					dest = p.targets[id(r)]
					break
				}
			}
		}
		if sf := p.sources[attr(f, "source")]; dest != "" && sf != nil {
			anchor := attr(f, "generated_anchor")
			if anchor == "" {
				anchor = attr(f, "authored_anchor")
			}
			if anchor != "" {
				p.targets[sourceKey(sf.path, anchor)] = dest
			}
			p.targets[id(f)] = dest
		}
	}
	for id, sf := range p.sources {
		if dest := p.targets[id]; dest != "" {
			p.targets[sourceKey(sf.path, "")] = dest
		}
	}
}

// NUL cannot occur in an OS path or authored Markdown fragment. Using it
// as a separator distinguishes filenames containing '#' from fragment URLs.
func sourceKey(file, fragment string) string { return "path:" + file + "\x00" + fragment }
func (p *planner) ownerPage(r model.Record) *subject {
	for r != nil {
		if s := p.pages[id(r)]; s != nil {
			return s
		}
		r = p.elements[attr(r, "owner")]
	}
	return nil
}
