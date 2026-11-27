package specdoc

import (
	"cmp"
	"crypto/sha256"
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
	otherSourceLabels map[sourceRoute]int // authored or navigational links that share a provenance label
	Document
	Title, Markdown                                 string
	SectionDocs                                     map[string]Document
	DirectoryMembers                                map[string]directoryContract
	RequiredLinks, RequiredContent, RequiredAnchors []string
	RequiredLabels, RequiredSectionLinks            map[string][]string
	RequiredSectionContent                          map[string][]string
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
	field                               string
	bodyHeadings                        []string // rendered headings inside quotes and lists in this passage
	start, end, bodyEnd, level          int
	owner                               *subject
	detailEntry                         *heading // root of a passage collected on a detail page
	qualityAttribute                    bool
	inGlossary                          bool
	subjectKind                         string
}

type subject struct {
	r                   model.Record
	source              *sourceFile
	name                *heading
	system              *subject
	navigationParent    *subject
	indexes             map[string]*collectionIndex
	glossary            []glossaryPassage
	passages            [][]*heading
	overview            []overviewSection
	externalSystem      bool
	emptyHeadingAnchors []string
	file, dir           string
}

type planner struct {
	evidence        model.VerificationEvidence
	assessments     map[string][]*model.ProcedureAssessment
	companions      map[string][]companionView // by subject record
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
	routes          map[string][]route
}

type entry struct {
	title, target, preview string
	note                   string  // italic aside after the preview
	children               []entry // nested under this entry in a directory list
	childLabel             string  // italic label introducing the children
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
	if err := p.readEvidence(); err != nil {
		return nil, err
	}
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
		data := []byte(attr(r, "markdown"))
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
	nested := map[*sourceFile][]model.Record{}
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
				nested[sf] = append(nested[sf], n)
				continue
			}
			h := &heading{node: id(n), field: attr(n, "domain_field"), title: attr(n, "decoded_text"), level: model.Number(n["heading_level"]), start: position(n, "start"), end: position(n, "end")}
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
		// Nested headings participate in Markdown fragment numbering even
		// though they do not declare domain sections.
		bodyHeadings := nested[sf]
		slices.SortFunc(bodyHeadings, func(a, b model.Record) int { return cmp.Compare(position(a, "start"), position(b, "start")) })
		for i, h := range sf.headings {
			h.bodyEnd = len(sf.data)
			if i+1 < len(sf.headings) {
				h.bodyEnd = sf.headings[i+1].start
			}
			for len(bodyHeadings) > 0 && position(bodyHeadings[0], "start") < h.bodyEnd {
				if position(bodyHeadings[0], "start") >= h.end {
					h.bodyHeadings = append(h.bodyHeadings, attr(bodyHeadings[0], "decoded_text"))
				}
				bodyHeadings = bodyHeadings[1:]
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
			h.qualityAttribute = h.subjectKind == "QualityModelAttribute" || h.subjectKind == "SystemQualityAttribute"
		}
		if sf == nil || families[attr(r, "kind")] == "" || h == nil {
			continue
		}
		s := &subject{r: r, source: sf, name: h}
		h.owner = s
		p.pages[id(r)] = s
		p.subjects = append(p.subjects, s)
	}
	// Every selected source has one independently authored subject. Its H1 must
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
	if s.isOverview() || s.kind() == "SystemQualityAttribute" {
		return strings.TrimSuffix(path.Base(s.source.path), specfile.AuthoredSuffix(kindName(s.kind())))
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
func (s *subject) isOverview() bool {
	return s.kind() == "System" || s.kind() == "Component"
}

func (s *subject) kind() string {
	return attr(s.r, "kind")
}

// A system or component owns a directory. Its overview mirrors its source filename,
// including when a sibling with the same name needs a distinct directory.
func (s *subject) placeOverview(dir string) {
	s.dir = dir
	stem := strings.TrimSuffix(path.Base(s.source.path), specfile.AuthoredSuffix(kindName(s.kind())))
	s.file = path.Join(dir, stem+reportSuffix(s.kind()))
}

func (s *subject) collectionDir() string {
	if s.kind() == "System" {
		return path.Join(s.navigationParent.dir, "systems")
	}
	return path.Join(s.navigationParent.dir, families[s.kind()])
}

func (p *planner) place() error {
	var roots []*subject
	members := map[string]model.ProductMember{}
	for _, s := range p.subjects {
		if s.kind() == "System" && p.pages[attr(s.r, "part_of")] == nil {
			roots = append(roots, s)
			product, err := p.snapshot.Product(id(s.r))
			if err != nil {
				return err
			}
			for _, member := range product.Members {
				members[member.ID] = member
			}
		}
	}
	if len(p.subjects) == 0 {
		return fmt.Errorf("no root system or quality model found in this project")
	}
	// A project can author a specification for a system in another project.
	// Keep its report reachable at project level without enrolling it locally.
	hasExternal := false
	for _, s := range p.subjects {
		system := attr(s.r, "specification_system")
		if system != "" && p.pages[system] == nil && attr(p.elements[system], "kind") == "System" {
			s.externalSystem = true
			hasExternal = true
		}
	}
	if len(roots) == 1 && !hasExternal {
		p.governing = roots[0]
		p.governing.placeOverview("specification")
	} else {
		title := attr(model.Object(p.snapshot.Scope["project"]), "name")
		if title == "" {
			title = "Specifications"
		}
		p.governing = &subject{name: &heading{title: title}, file: "SPEC.md", dir: "specification"}
	}
	for _, s := range p.subjects {
		if s == p.governing {
			continue
		}
		if s.kind() == "QualityModel" || s.externalSystem || slices.Contains(roots, s) {
			s.navigationParent = p.governing
			continue
		}
		member, found := members[id(s.r)]
		if !found {
			owner := p.elements[attr(s.r, "owner")]
			for owner != nil {
				if member, found = members[id(owner)]; found {
					break
				}
				owner = p.elements[attr(owner, "owner")]
			}
		}
		if !found {
			if parent := p.pages[attr(s.r, "owner")]; parent != nil && parent.externalSystem {
				s.navigationParent = parent
				continue
			}
			return fmt.Errorf("%s: cannot place %s in a selected system", s.source.path, s.name.title)
		}
		s.system = p.pages[member.System]
		switch s.kind() {
		case "System", "Component":
			parent := p.pages[attr(s.r, "part_of")]
			if parent != nil && (s.kind() == "System" || members[id(parent.r)].System == member.System) {
				s.navigationParent = parent
			}
		default:
			s.navigationParent = p.pages[attr(s.r, "owner")]
			if s.navigationParent == nil {
				owner := p.pages[attr(s.r, "specification_owner")]
				if owner != nil && members[id(owner.r)].System == member.System {
					s.navigationParent = owner
				}
			}
		}
		if s.navigationParent == nil {
			s.navigationParent = s.system
		}
	}
	// Disambiguate siblings by stable project-relative source and authored heading.
	names := map[string]int{}
	key := func(s *subject) string { return id(s.navigationParent.r) + "/" + families[s.kind()] + "/" + s.slug() }
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
		if err := assign(s.navigationParent); err != nil {
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
			s.file = path.Join(s.dir, slug+reportSuffix(s.kind()))
		} else if s.isOverview() {
			s.placeOverview(base)
		} else {
			s.file = base + reportSuffix(s.kind())
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
		p.targets[s.name.node] = s.file
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
		// Concept uses are listed under their concepts; see indexConceptDefinitions.
		if s != p.governing && s.kind() != "ConceptUse" {
			owner := s.navigationParent
			if s.kind() == "SystemQualityAttribute" {
				owner = s.system
			}
			owner.includeMember(s.collection(), entry{title: s.name.title, target: s.file, preview: subjectPreview(s)})
		}
		if s.isOverview() {
			p.collectOverviewDetails(s)
		}
	}
	for _, s := range p.subjects {
		switch s.kind() {
		case "BusinessRule", "Requirement":
			for _, target := range model.Strings(s.r["applies_to"]) {
				if owner := p.pages[target]; owner != nil && (owner.kind() == "System" || owner.kind() == "Component") {
					e := entry{title: s.name.title, target: s.file, preview: subjectPreview(s)}
					owner.includeMember(s.collectionType(), e)
				}
			}
		case "Capability":
			for _, target := range model.Strings(s.r["provided_by"]) {
				if provider := p.pages[target]; provider != nil {
					provider.includeMember(collectionTitle("Capability"), entry{title: s.name.title, target: s.file, preview: subjectPreview(s)})
				}
			}
		}
	}
	p.indexConceptDefinitions()
	p.indexOpenQuestions()
	p.readCompanions()
	p.indexCompanions()
	for _, s := range p.subjects {
		if s.kind() == "System" {
			p.planGlossary(s)
		}
		p.indexDetailHeadings(s)
		if s.isOverview() {
			continue
		}
		s.passages = p.planPassages(s, nil)
		ids := parser.NewContext().IDs()
		for _, h := range readingOrder(s) {
			if h.owner != s || h.detailEntry != nil {
				continue
			}
			p.targets[h.node] = s.file
			generated := string(ids.Generate([]byte(h.reportTitle()), ast.KindHeading))
			if h != s.name {
				if needsHeadingAnchor(s, h) {
					h.explicitAnchor = "section-" + h.anchor
					generated = h.explicitAnchor
				}
				p.targets[h.node] += "#" + generated
			}
		}
	}
	for _, s := range p.subjects {
		if s.isOverview() {
			s.overview = p.planOverview(s)
			p.indexOverviewHeadings(s)
			p.planCollectionHomes(s)
		}
	}
	for _, r := range p.snapshot.Elements {
		if dest := p.targets[nameNode(r)]; dest != "" {
			p.targets[id(r)] = dest
		}
	}
	// Question entries initially point at source headings. Their final fragments
	// are known once every subject and overview has been laid out.
	for _, s := range p.subjects {
		idx := s.indexes["Open questions"]
		if idx == nil {
			continue
		}
		for i, e := range idx.members {
			idx.members[i].target = p.targets[e.target]
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
