package model

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type ExportResult struct {
	Snapshot *Snapshot
	ExitCode int
}

type snapshotBuilder struct {
	snapshot    *Snapshot
	types       []check.Type
	byPath      map[string]*ParsedDocument
	sources     map[string]Record
	elements    map[string]Record
	rootByPath  map[string]string
	primary     map[string]bool
	follow      map[string]bool
	contextPath map[string]string
	contextOK   map[string]bool
	// directorySystem records, by root element, the nearest enclosing system
	// found by directory discovery. It supplies the default association only.
	directorySystem map[string]string
	glossaries      map[string]map[string]bool // glossary path to the system sources that declared it
	attributeOwners map[string]string          // validated quality source path to declaring system path
	readFailed      bool
	sequence        int
}

// Export observes the same selected specifications, enclosing systems and
// direct reference boundary as check. It records source text before mapping,
// including invalid specifications. It never rewrites an input.
func Export(path string, types ...check.Type) (ExportResult, error) {
	dir, err := specfile.Directory(path)
	if err != nil {
		return ExportResult{}, err
	}
	inspection := check.InspectSources(dir, types...)
	return ExportInspection(inspection, types...), nil
}

// ExportInspection maps a previously inspected project's precise discovery result.
func ExportInspection(inspection check.Inventory, types ...check.Type) ExportResult {
	report := inspection.Report
	dir := report.Scope.Path
	repo := inspection.Repository
	kinds := make([]string, 0, len(types))
	for _, kind := range types {
		kinds = append(kinds, specfile.AuthoredType(kind.Layout.Kind))
	}
	scope := Record{
		"roots": []string{filepath.ToSlash(dir)}, "selected_spec_types": kinds,
		"discovery_complete": inspection.DiscoveryComplete,
		"reads_complete":     true, "inspection_complete": report.Complete,
		"discovery_policy":    Record{"mode": report.Scope.Discovery, "rule": "dotspec-discovery/1", "follow_symlinks": false, "include_nested_repositories": false},
		"reference_traversal": Record{"mode": "direct", "include_system_context": true, "follow_external_resources": false, "repository_bounded": repo != "", "maximum_depth": 1, "rule": "dotspec-direct-references/1"},
	}
	b := &snapshotBuilder{
		snapshot: NewSnapshot(scope), types: types, byPath: map[string]*ParsedDocument{},
		sources: map[string]Record{}, elements: map[string]Record{}, rootByPath: map[string]string{},
		primary: map[string]bool{}, follow: map[string]bool{}, contextPath: map[string]string{}, contextOK: map[string]bool{},
		directorySystem: map[string]string{}, attributeOwners: map[string]string{}, glossaries: map[string]map[string]bool{},
	}
	for _, file := range inspection.Files {
		for _, ref := range file.ReferencedBy {
			if ref.Element == specfile.GlossarySourcesField && system.Layout.Match(ref.Path) {
				target := cleanPath(file.Path)
				if b.glossaries[target] == nil {
					b.glossaries[target] = map[string]bool{}
				}
				b.glossaries[target][cleanPath(ref.Path)] = true
			}
		}
	}
	for _, path := range inspection.Selected {
		path = cleanPath(path)
		b.primary[path], b.follow[path] = true, true
		b.load(path)["selected"] = true
	}
	for _, source := range inspection.Files {
		b.load(cleanPath(source.Path))
	}
	// Resolve context once for each selected specification owner. Referenced
	// documents do not recursively expand their own context or outgoing links.
	for _, path := range sortedKeys(b.primary) {
		kind := b.kind(path)
		if kind == nil || kind.Context == nil {
			b.contextOK[path] = true
			continue
		}
		owner := kind.Layout.ContextDir(filepath.FromSlash(path))
		_, explicit := b.elements[b.rootByPath[path]]["specification_system_content"]
		var contexts []string
		var contextErr error
		if explicit && kind.Context.ResolveSource != nil {
			contexts, contextErr = kind.Context.ResolveSource(owner, filepath.FromSlash(path))
		} else {
			contexts, contextErr = kind.Context.Resolve(owner)
		}
		b.contextOK[path] = contextErr == nil
		for _, context := range contexts {
			context = cleanPath(context)
			if system.AttributeLayout.Match(context) {
				b.load(context)
				continue
			}
			if system.ComponentLayout.Match(context) {
				b.load(context)
				b.follow[context] = true
				continue
			}
			if explicit {
				continue
			}
			b.contextPath[path], b.follow[context] = context, true
			b.load(context)
			b.context(context, "system_context", path, "")
		}
		subject := b.elements[b.rootByPath[path]]
		if subject != nil && designKind(String(subject, "kind"), "SystemContextSpecification") && String(subject, "kind") != "Component" {
			boundary, err := specfile.ContextRoot(owner)
			if err == nil {
				declaration, err := specfile.NearestDeclaration(owner, boundary)
				if ownerID := b.rootByPath[filepath.ToSlash(declaration)]; err == nil && ownerID != "" {
					subject["specification_owner"] = ownerID
				}
			}
		}
	}
	// Attribute declarations are part of each acquired system's classification vocabulary.
	for _, file := range sortedKeys(b.rootByPath) {
		if !system.Layout.Match(file) {
			continue
		}
		sys, err := system.Read(file)
		if err != nil {
			continue
		}
		for _, a := range sys.QualityAttributes {
			if system.AttributeLayout.Match(a.Location.Path) {
				source := cleanPath(a.Location.Path)
				b.attributeOwners[source] = file
				b.load(source)
				b.follow[source] = true
			}
		}
	}
	b.associateGlossaries()
	b.associateSystems()
	b.resolveReferences()
	b.resolveLabels()
	compositionErrors := b.resolveComposition()
	for _, finding := range report.Diagnostics {
		b.addFinding(finding)
	}
	if p := report.Scope.Project; p != nil {
		b.snapshot.Scope["project"] = Record{"name": p.Name, "root": cleanPath(p.Root)}
		b.context(cleanPath(p.Manifest()), "project_configuration", "", "")
	}
	for _, source := range b.snapshot.Sources {
		if String(source, "kind") != "SourceDocument" {
			continue
		}
		path := filepath.FromSlash(String(source, "locator"))
		p, err := specfile.FindProject(filepath.Dir(path))
		if err != nil {
			b.readFailed = true
			b.addError("scan.context", err, path)
		} else if p != nil {
			source["project"] = Record{"name": p.Name, "root": cleanPath(p.Root)}
		}
	}
	b.materializeAssertions()
	b.finishCoverage()
	b.snapshot.Scope["reads_complete"] = !b.readFailed
	if b.readFailed || !inspection.DiscoveryComplete {
		b.snapshot.Scope["inspection_complete"] = false
	}
	b.sort()
	code := report.ExitCode()
	if compositionErrors && code == 0 {
		code = 1
	}
	for _, finding := range b.snapshot.Diagnostics {
		if String(finding, "severity") == "error" && code == 0 {
			code = 1
		}
	}
	if b.readFailed || !inspection.DiscoveryComplete {
		code = 2
	}
	return ExportResult{b.snapshot, code}
}

func cleanPath(path string) string {
	abs, err := filepath.Abs(filepath.FromSlash(path))
	if err == nil {
		path = abs
	}
	return filepath.ToSlash(filepath.Clean(path))
}

func (b *snapshotBuilder) next(prefix string) string {
	b.sequence++
	return fmt.Sprintf("%s:%d", prefix, b.sequence)
}

func (b *snapshotBuilder) kind(path string) *check.Type {
	for i := range b.types {
		if b.types[i].Layout.Match(path) {
			return &b.types[i]
		}
	}
	return nil
}

func (b *snapshotBuilder) load(path string) Record {
	if source := b.sources[path]; source != nil {
		return source
	}
	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		source := Record{"id": SourceID(path), "kind": "SourceDocument", "locator": path, "read_status": "failed"}
		b.sources[path] = source
		b.snapshot.Sources = append(b.snapshot.Sources, source)
		b.readFailed = true
		b.addError("scan.read", err, path)
		return source
	}
	kind := b.kind(path)
	ext := strings.ToLower(filepath.Ext(path))
	if kind == nil && ext != ".md" && ext != ".markdown" {
		source := sourceRecord(path)
		source["media_type"] = "application/octet-stream"
		b.sources[path] = source
		b.snapshot.Sources = append(b.snapshot.Sources, source)
		return source
	}
	doc := ParseContent(path, data)
	b.byPath[path], b.sources[path] = doc, doc.Source
	kindName := ""
	if kind != nil {
		kindName = kind.Layout.Kind
	} else if b.glossaries[path] != nil {
		kindName = "glossary"
	}
	mapping := MapDomain(doc, kindName)
	b.rootByPath[path] = mapping.RootID
	b.snapshot.Sources = append(b.snapshot.Sources, doc.Source)
	b.snapshot.Documents = append(b.snapshot.Documents, doc.Document)
	b.snapshot.ContentNodes = append(b.snapshot.ContentNodes, doc.Nodes...)
	b.snapshot.Fragments = append(b.snapshot.Fragments, doc.Fragments...)
	b.snapshot.References = append(b.snapshot.References, doc.References...)
	b.snapshot.Elements = append(b.snapshot.Elements, mapping.Elements...)
	b.snapshot.LabelOccurrences = append(b.snapshot.LabelOccurrences, mapping.Labels...)
	b.snapshot.ContentBindings = append(b.snapshot.ContentBindings, mapping.Bindings...)
	for _, finding := range doc.Diagnostics {
		if kind == nil {
			finding["severity"] = "advice"
		}
		b.snapshot.Diagnostics = append(b.snapshot.Diagnostics, finding)
	}
	b.snapshot.Diagnostics = append(b.snapshot.Diagnostics, mapping.Diagnostics...)
	for _, element := range mapping.Elements {
		b.elements[String(element, "id")] = element
	}
	doc.Document["elements"] = recordIDs(mapping.Elements)
	return doc.Source
}

func (b *snapshotBuilder) context(path, role, requester, reference string) {
	source := b.sources[path]
	if source == nil || String(source, "kind") != "SourceDocument" {
		return
	}
	for _, context := range b.snapshot.ContextSources {
		if String(context, "source") == String(source, "id") && String(context, "role") == role {
			if request := b.sources[requester]; request != nil {
				context["requested_by"] = appendUnique(Strings(context["requested_by"]), String(request, "id"))
			}
			if reference != "" {
				context["reference_occurrences"] = appendUnique(Strings(context["reference_occurrences"]), reference)
			}
			return
		}
	}
	context := Record{"source": source["id"], "role": role, "rule": "dotspec-context/1"}
	if request := b.sources[requester]; request != nil {
		context["requested_by"] = []string{String(request, "id")}
	}
	if reference != "" {
		context["reference_occurrences"] = []string{reference}
	}
	b.snapshot.ContextSources = append(b.snapshot.ContextSources, context)
}

func (b *snapshotBuilder) associateGlossaries() {
	for _, file := range sortedKeys(b.rootByPath) {
		sys := b.elements[b.rootByPath[file]]
		if String(sys, "kind") != "System" {
			continue
		}
		if sys["glossary_source_paths"] == nil {
			continue
		}
		sys["glossary_sources"] = []string{}
		for _, name := range Strings(sys["glossary_source_paths"]) {
			target := cleanPath(filepath.Join(filepath.Dir(file), filepath.FromSlash(name)))
			vocab := b.byPath[target]
			if !b.glossaries[target][file] || vocab == nil || vocab.Document["glossary_terms"] == nil {
				continue
			}
			appendID(sys, "glossary_sources", String(vocab.Document, "id"))
			b.context(target, "glossary_source", file, "")
		}
	}
}

func (b *snapshotBuilder) associateSystems() {
	for _, path := range sortedKeys(b.rootByPath) {
		root := b.elements[b.rootByPath[path]]
		if root == nil {
			continue
		}
		if context := b.elements[b.rootByPath[b.contextPath[path]]]; String(context, "kind") == "System" && String(root, "kind") != "System" && String(root, "kind") != "QualityModel" {
			b.directorySystem[String(root, "id")] = String(context, "id")
		}
		if String(root, "kind") != "System" {
			continue
		}
		owner := system.Layout.ContextDir(filepath.FromSlash(path))
		repository, boundaryErr := specfile.ContextRoot(owner)
		if boundaryErr != nil {
			continue
		}
		if owner == repository {
			continue
		}
		for current := filepath.Dir(owner); current != owner; current = filepath.Dir(current) {
			file, err := system.Layout.LocalSubject(current)
			if err != nil {
				break
			}
			candidate := cleanPath(file)
			if file != "" {
				if parent := b.elements[b.rootByPath[candidate]]; String(parent, "kind") == "System" {
					root["enclosing_system"] = parent["id"]
				}
				break // An uninspected nearer system must not be skipped.
			}
			if current == repository || filepath.Dir(current) == current {
				break
			}
		}
	}
}

func (b *snapshotBuilder) addFinding(finding diagnostic.Diagnostic) {
	path := cleanPath(finding.Location.Path)
	source := b.sources[path]
	for _, existing := range b.snapshot.Diagnostics {
		if String(existing, "code") == finding.Rule && String(existing, "message") == finding.Message && String(existing, "source") == String(source, "id") && Number(existing["line"]) == finding.Location.Line {
			return
		}
	}
	d := Record{"id": b.next("finding"), "kind": "Diagnostic", "code": finding.Rule, "severity": finding.Severity, "message": finding.Message, "phase": "inspect"}
	if String(source, "kind") == "SourceDocument" {
		d["source"] = source["id"]
	}
	if finding.Location.Line > 0 {
		d["line"] = finding.Location.Line
	}
	if finding.Location.Column > 0 {
		d["column"] = finding.Location.Column
	}
	if finding.Location.Element != "" {
		d["element"] = finding.Location.Element
	}
	if len(finding.Related) > 0 {
		related := []Record{}
		for _, location := range finding.Related {
			item := Record{}
			if location.Path != "" {
				item["locator"] = cleanPath(location.Path)
				if target := b.sources[cleanPath(location.Path)]; String(target, "kind") == "SourceDocument" {
					item["source"] = target["id"]
				}
			}
			if location.Line > 0 {
				item["line"] = location.Line
			}
			if location.Column > 0 {
				item["column"] = location.Column
			}
			if location.Element != "" {
				item["element"] = location.Element
			}
			related = append(related, item)
		}
		d["related"] = related
	}
	if source == nil && finding.Location.Path != "" {
		d["message"] = fmt.Sprintf("%s: %s", finding.Location.Path, finding.Message)
	}
	b.snapshot.Diagnostics = append(b.snapshot.Diagnostics, d)
}

func (b *snapshotBuilder) addError(code string, err error, path string) {
	b.addFinding(*diagnostic.Error(code, err.Error(), diagnostic.Location{Path: path}))
}

func recordIDs(records []Record) []string {
	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, String(record, "id"))
	}
	return ids
}

func appendUnique(items []string, item string) []string {
	for _, existing := range items {
		if existing == item {
			return items
		}
	}
	return append(items, item)
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func resourceID(prefix, locator string) string {
	hash := sha256.Sum256([]byte(locator))
	return fmt.Sprintf("%s:%x", prefix, hash[:10])
}

func (b *snapshotBuilder) sort() {
	for _, records := range [][]Record{b.snapshot.Sources, b.snapshot.Documents, b.snapshot.ContentNodes, b.snapshot.Fragments, b.snapshot.Elements, b.snapshot.References, b.snapshot.LabelOccurrences, b.snapshot.Resolutions, b.snapshot.Assertions, b.snapshot.ContentBindings, b.snapshot.Diagnostics} {
		sort.Slice(records, func(i, j int) bool { return String(records[i], "id") < String(records[j], "id") })
	}
	sort.Slice(b.snapshot.ContextSources, func(i, j int) bool {
		a, c := b.snapshot.ContextSources[i], b.snapshot.ContextSources[j]
		return String(a, "source")+String(a, "role") < String(c, "source")+String(c, "role")
	})
}
