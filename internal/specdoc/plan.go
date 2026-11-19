package specdoc

import (
	"cmp"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"

	"github.com/craigsmitham/dotspec/internal/model"
)

// Plan binds source-backed expectations to their project-relative generated paths.
type Plan struct{ Pages map[string]*Page }

type Page struct {
	Type, SubjectKind, Title, Markdown                       string
	RequiredLinks, RequiredContent, Members, RequiredAnchors []string
	RequiredLabels, RequiredSectionLinks                     map[string][]string
	RequiredSectionContent                                   map[string]string
	RequiredResources                                        []string
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
}

type subject struct {
	r         model.Record
	source    *sourceFile
	name      *heading
	system    *subject
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
	collections     map[string]*collection
	routes          map[string][]route
}

type collection struct {
	file, label string
	system      *subject
	entries     []entry
}
type entry struct{ title, target string }

var families = map[string]string{
	"System": "systems", "Component": "components", "Feature": "features", "Capability": "capabilities",
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
	p := &planner{root: root, snapshot: snapshot, plan: &Plan{Pages: map[string]*Page{}}, sources: map[string]*sourceFile{}, nodes: map[string]model.Record{}, elements: map[string]model.Record{}, pages: map[string]*subject{}, targets: map[string]string{}, collections: map[string]*collection{}}
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
	for _, r := range p.snapshot.Sources {
		if selected, _ := r["selected"].(bool); !selected {
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
			switch attr(r, "kind") {
			case "ObligationClause":
				h.obligation = true
			case "QualityAttribute":
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
	slices.SortFunc(p.subjects, func(a, b *subject) int {
		if n := strings.Compare(a.source.path, b.source.path); n != 0 {
			return n
		}
		return cmp.Compare(a.name.start, b.name.start)
	})
	return nil
}

func (s *subject) slug() string {
	name := attr(s.r, "authored_id")
	if name == "" {
		name = s.name.title
	}
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			dash = false
		} else {
			dash = true
		}
	}
	if b.Len() == 0 {
		return "section"
	}
	return b.String()
}
func (s *subject) kind() string { return attr(s.r, "kind") }
func (s *subject) collectionDir() string {
	if s.kind() == "System" {
		return path.Join(s.system.dir, "systems")
	}
	return path.Join(s.system.dir, families[s.kind()])
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
		if parent == "" {
			parent = attr(s.r, "owner")
		}
		for parent != "" && attr(p.elements[parent], "kind") != "System" {
			parent = attr(p.elements[parent], "owner")
		}
		s.system = p.pages[parent]
		if s.system == nil {
			return fmt.Errorf("%s: cannot place %s in a selected system", s.source.path, s.name.title)
		}
	}
	if len(roots) != 1 {
		return fmt.Errorf("generated specification requires one root system; found %d (%s)", len(roots), strings.Join(roots, ", "))
	}
	p.governing.file = "SPEC.md"
	p.governing.dir = "specification"
	for _, s := range p.subjects {
		if s.kind() == "QualityModel" {
			s.system = p.governing
		}
	}
	// Disambiguate siblings by stable project-relative source and authored heading.
	names := map[string]int{}
	key := func(s *subject) string { return id(s.system.r) + "/" + families[s.kind()] + "/" + s.slug() }
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
		if err := assign(s.system); err != nil {
			return err
		}
		slug := s.slug()
		if slug == "index" || names[key(s)] > 1 {
			sum := sha256.Sum256([]byte(s.source.path + "#" + s.name.anchor))
			slug += fmt.Sprintf("-%x", sum[:5])
		}
		base := path.Join(s.collectionDir(), slug)
		if s.kind() == "System" {
			s.dir = base
			s.file = base + "/SPEC.md"
		} else {
			s.file = base + ".md"
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

func (p *planner) indexTargets() {
	for _, s := range p.subjects {
		p.targets[id(s.r)] = s.file
		if s.name.level == 1 {
			p.targets[attr(s.r, "source_document")] = s.file
		}
		if s != p.governing {
			p.addMember(s.system, families[s.kind()], entry{s.name.title, s.file})
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
		ids := parser.NewContext().IDs()
		for _, h := range readingOrder(s) {
			if h.owner != s {
				continue
			}
			p.targets[h.node] = s.file
			generated := string(ids.Generate([]byte(h.title), ast.KindHeading))
			if h != s.name {
				if strings.IndexFunc(h.title, func(r rune) bool { return r > 127 }) >= 0 || strings.ContainsAny(s.source.data[h.start:h.end], "[]*_~`") {
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
		if dest := p.targets[id(r)]; attr(r, "kind") == "ObligationClause" && dest != "" {
			page := p.ownerPage(r)
			if page != nil {
				sys := page.system
				if page.kind() == "System" {
					sys = page
				}
				n := p.nodes[nameNode(r)]
				p.addMember(sys, "requirements", entry{attr(n, "decoded_text"), dest})
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
func (p *planner) addMember(sys *subject, folder string, e entry) {
	file := path.Join(sys.dir, folder, "index.md")
	c := p.collections[file]
	if c == nil {
		c = &collection{file: file, label: strings.ReplaceAll(folder, "-", " "), system: sys}
		p.collections[file] = c
	}
	c.entries = append(c.entries, e)
}
