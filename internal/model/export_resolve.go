package model

import (
	"errors"
	"io/fs"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func (b *snapshotBuilder) resolution(occurrence Record, status, target, rule string) Record {
	r := Record{"id": b.next("resolution"), "kind": "Resolution", "occurrence": occurrence["id"], "status": status, "rule": rule}
	if target != "" {
		r["target"] = target
	}
	occurrence["resolution"] = r["id"]
	b.snapshot.Resolutions = append(b.snapshot.Resolutions, r)
	return r
}

func (b *snapshotBuilder) resolveReferences() {
	// Loading direct targets can add references, but only the originally
	// selected sources and enclosing context may trigger filesystem traversal.
	for i := 0; i < len(b.snapshot.References); i++ {
		reference := b.snapshot.References[i]
		path, doc := b.documentForSource(String(reference, "source_document"))
		if doc == nil {
			b.resolution(reference, "not_attempted", "", "source-unavailable/1")
			continue
		}
		b.referenceContext(reference, doc)
		target := String(reference, "authored_target")
		node := doc.NodeByID[String(reference, "content_node")]
		form := String(reference, "form")
		if form == "footnote" {
			definition := String(reference, "definition_node")
			if definition == "" {
				b.resolution(reference, "undefined_label", "", "footnote-label/1")
			} else {
				b.resolution(reference, "resolved", b.fragmentForNode(doc, definition), "footnote-label/1")
			}
			continue
		}
		if form == "reference_link" || (form == "image" && String(node, "label") != "") {
			definition := doc.NodeByID[String(reference, "definition_node")]
			if definition == nil {
				b.resolution(reference, "undefined_label", "", "markdown-reference-label/1")
				continue
			}
			target = String(definition, "destination")
			if target == "" {
				target = String(node, "destination")
			}
		}
		u, err := url.Parse(target)
		if err != nil {
			b.resolution(reference, "unresolved", "", "invalid-reference/1")
			continue
		}
		if u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") {
			id := resourceID("external", target)
			if b.sources["external:"+target] == nil {
				r := Record{"id": id, "kind": "ExternalResource", "locator": target}
				b.sources["external:"+target] = r
				b.snapshot.Sources = append(b.snapshot.Sources, r)
			}
			b.resolution(reference, "external", "", "external-navigation/1")
			continue
		}
		if strings.Contains(u.Path, "\\") {
			b.resolution(reference, "unresolved", "", "invalid-reference-path/1")
			continue
		}
		targetPath := path
		if u.Path != "" {
			targetPath = cleanPath(filepath.Join(filepath.Dir(filepath.FromSlash(path)), filepath.FromSlash(u.Path)))
		}
		source := b.sources[targetPath]
		mayFollow := b.follow[path] && b.rootByPath[path] != ""
		if source == nil && targetPath != path {
			if !mayFollow {
				b.resolution(reference, "not_attempted", "", "direct-reference-boundary/1")
				continue
			}
			root, rootErr := specfile.RepositoryRoot(filepath.Dir(filepath.FromSlash(path)))
			if rootErr != nil {
				b.resolution(reference, "not_attempted", "", "repository-boundary-unavailable/1")
				b.readFailed = true
				b.addError("scan.reference", rootErr, path)
				continue
			}
			info, allowed, lookupErr := check.ReferenceTarget(filepath.FromSlash(targetPath), root)
			if !allowed {
				b.resolution(reference, "out_of_scope", "", "reference-boundary/1")
				continue
			}
			if lookupErr != nil {
				b.resolution(reference, "unresolved", "", "missing-reference-target/1")
				if !errors.Is(lookupErr, fs.ErrNotExist) {
					b.readFailed = true
					b.addError("scan.reference", lookupErr, targetPath)
				}
				continue
			}
			if info.IsDir() {
				source = Record{"id": resourceID("directory", targetPath), "kind": "SourceDirectory", "locator": targetPath}
				b.sources[targetPath] = source
				b.snapshot.Sources = append(b.snapshot.Sources, source)
			} else {
				source = b.load(targetPath)
			}
		}
		if source == nil || (String(source, "kind") == "SourceDocument" && String(source, "read_status") != "read") {
			b.resolution(reference, "not_attempted", "", "unreadable-reference-target/1")
			continue
		}
		if targetPath != path && mayFollow {
			b.context(targetPath, "reference_target", path, String(reference, "id"))
		}
		if u.Fragment == "" {
			b.resolution(reference, "resolved", String(source, "id"), "source-reference/1")
			continue
		}
		if String(source, "kind") == "SourceDirectory" {
			b.resolution(reference, "unresolved", "", "directory-has-no-fragment/1")
			continue
		}
		targetDoc := b.byPath[targetPath]
		if targetDoc == nil {
			b.resolution(reference, "not_attempted", "", "native-fragment-semantics/1")
			continue
		}
		matches := b.fragmentTargets(targetDoc, u.Fragment)
		switch len(matches) {
		case 0:
			b.resolution(reference, "unresolved", "", "markdown-fragment/1")
		case 1:
			b.resolution(reference, "resolved", matches[0], "markdown-fragment/1")
		default:
			r := b.resolution(reference, "ambiguous", "", "markdown-fragment/1")
			r["candidates"] = matches
		}
	}
}

func (b *snapshotBuilder) documentForSource(id string) (string, *ParsedDocument) {
	for path, doc := range b.byPath {
		if String(doc.Source, "id") == id {
			return path, doc
		}
	}
	return "", nil
}

func (b *snapshotBuilder) fragmentTargets(doc *ParsedDocument, anchor string) []string {
	var matches []string
	for _, element := range b.snapshot.Elements {
		if String(element, "source_document") == String(doc.Source, "id") && String(element, "authored_id") == anchor {
			matches = append(matches, String(element, "id"))
		}
	}
	if len(matches) > 0 {
		sort.Strings(matches)
		return matches
	}
	for _, fragment := range b.snapshot.Fragments {
		if String(fragment, "source") == String(doc.Source, "id") && (String(fragment, "generated_anchor") == anchor || String(fragment, "authored_anchor") == anchor) {
			matches = append(matches, String(fragment, "id"))
		}
	}
	sort.Strings(matches)
	return matches
}

func (b *snapshotBuilder) fragmentForNode(doc *ParsedDocument, node string) string {
	for _, fragment := range b.snapshot.Fragments {
		if String(fragment, "content_node") == node {
			return String(fragment, "id")
		}
	}
	id := resourceID("fragment", node)
	f := Record{"id": id, "kind": "SourceFragment", "source": doc.Source["id"], "content_node": node, "span": doc.NodeByID[node]["span"]}
	b.snapshot.Fragments = append(b.snapshot.Fragments, f)
	return id
}

func (b *snapshotBuilder) referenceContext(reference Record, doc *ParsedDocument) {
	span := Object(doc.NodeByID[String(reference, "content_node")]["span"])
	best := int(^uint(0) >> 1)
	for _, binding := range b.snapshot.ContentBindings {
		if b.elements[String(binding, "owner")] == nil {
			continue
		}
		node := doc.NodeByID[String(binding, "content_node")]
		if node == nil {
			continue
		}
		bound := Object(node["span"])
		width := Number(bound["end"]) - Number(bound["start"])
		if Number(bound["start"]) <= Number(span["start"]) && Number(bound["end"]) >= Number(span["end"]) && width < best {
			best = width
			reference["context_owner"], reference["context_field"] = binding["owner"], binding["field_path"]
		}
	}
}

func (b *snapshotBuilder) resolveLabels() {
	// Direct targets may be associated with already acquired system context,
	// but no additional context is read merely to complete their classifications.
	for _, path := range sortedKeys(b.rootByPath) {
		if b.elements[b.rootByPath[path]] == nil {
			continue
		}
		if b.primary[path] || b.contextOK[path] {
			continue
		}
		kind := b.kind(path)
		if kind == nil || kind.Context == nil {
			continue
		}
		owner, err := kind.Layout.Owner(filepath.FromSlash(path))
		if err != nil {
			continue
		}
		repo, _ := specfile.RepositoryRoot(owner)
		for current := owner; ; current = filepath.Dir(current) {
			candidate := cleanPath(filepath.Join(current, "specs", "system.spec.md"))
			if source := b.sources[candidate]; source != nil {
				b.contextPath[path] = candidate
				b.contextOK[path] = b.rootByPath[candidate] != ""
				if sys := b.elements[b.rootByPath[candidate]]; String(sys, "kind") == "System" {
					b.elements[b.rootByPath[path]]["owning_system"] = sys["id"]
				}
				break
			}
			// An unacquired nearer system leaves the label uninspected, not unbound.
			if _, _, err := check.ReferenceTarget(filepath.FromSlash(candidate), repo); !errors.Is(err, fs.ErrNotExist) {
				break
			}
			if current == repo || filepath.Dir(current) == current {
				b.contextOK[path] = true
				break
			}
		}
	}
	attributes := map[string]map[string][]string{}
	for _, element := range b.snapshot.Elements {
		if String(element, "kind") != "QualityAttribute" {
			continue
		}
		owner, label := String(element, "owner"), String(element, "authored_id")
		if attributes[owner] == nil {
			attributes[owner] = map[string][]string{}
		}
		attributes[owner][label] = append(attributes[owner][label], String(element, "id"))
	}
	for _, label := range b.snapshot.LabelOccurrences {
		owner := b.elements[String(label, "owner")]
		path, _ := b.documentForSource(String(owner, "source_document"))
		sys := ""
		if String(label, "role") == "subattribute_of" {
			sys = String(owner, "owner")
		} else {
			sys = String(owner, "owning_system")
			if sys == "" && !b.contextOK[path] {
				b.resolution(label, "not_attempted", "", "uninspected-quality-context/1")
				continue
			}
		}
		if sys != "" {
			if String(b.elements[sys], "kind") == "QualityModel" {
				label["quality_model_context"] = sys
			} else {
				label["system_context"] = sys
			}
		}
		if _, declared := b.elements[sys]["quality_attributes"]; !declared {
			b.resolution(label, "unbound_label", "", "no-quality-attributes/1")
			continue
		}
		matches := attributes[sys][String(label, "label")]
		switch len(matches) {
		case 0:
			b.resolution(label, "undefined_label", "", "quality-attributes/1")
		case 1:
			b.resolution(label, "resolved", matches[0], "quality-attributes/1")
			field := "attributes"
			if String(label, "role") == "subattribute_of" {
				field = "subattribute_of"
			}
			owner[field] = appendUnique(Strings(owner[field]), matches[0])
		default:
			r := b.resolution(label, "ambiguous", "", "quality-attributes/1")
			r["candidates"] = matches
		}
	}
}

func (b *snapshotBuilder) materializeAssertions() {
	runtime := compiledModel()
	if runtime.err != nil {
		return // CLI validation reports unavailable schema/catalog as operational failure.
	}
	for _, kind := range sortedKeys(runtime.catalog.Classes) {
		definition := runtime.catalog.Classes[kind]
		projection := definition.Annotations["projection_of"]
		if projection == "" {
			continue
		}
		for _, expression := range strings.Split(projection, " or ") {
			parent, path, found := strings.Cut(strings.TrimSpace(expression), ".")
			if !found {
				continue
			}
			for _, subject := range b.snapshot.Elements {
				concrete := runtime.catalog.Classes[String(subject, "kind")]
				matches := false
				for _, ancestor := range concrete.Ancestors {
					matches = matches || ancestor == parent
				}
				if !matches {
					continue
				}
				for _, object := range valuesAt(subject, path) {
					assertion := Record{"id": b.next("assertion"), "kind": kind, "subject": subject["id"], "object": object, "origin": "derived", "mapping_rule": MappingProfile + "/projection/" + kind, "inputs": []string{String(subject, "id"), object}}
					b.snapshot.Assertions = append(b.snapshot.Assertions, assertion)
				}
			}
		}
	}
}

func (b *snapshotBuilder) finishCoverage() {
	allUnbound, allExtensions, allOpaque := []string{}, []string{}, []string{}
	allUnmapped := []Record{}
	status := "complete_for_profile"
	if b.readFailed {
		status = "partial"
	}
	for _, path := range sortedKeys(b.byPath) {
		doc := b.byPath[path]
		covered := map[string]bool{}
		var cover func(string)
		cover = func(id string) {
			if covered[id] {
				return
			}
			covered[id] = true
			for _, child := range Strings(doc.NodeByID[id]["children"]) {
				cover(child)
			}
		}
		opaque, extensions, unbound := []string{}, []string{}, []string{}
		unmapped := []Record{}
		for _, binding := range b.snapshot.ContentBindings {
			id := String(binding, "content_node")
			node := doc.NodeByID[id]
			if node == nil {
				continue
			}
			role := String(binding, "role")
			if role == "fallback" {
				opaque = appendUnique(opaque, id)
			}
			if String(binding, "scope") == "subtree" && (role != "structural_wrapper" || String(node, "node_kind") == "heading") {
				cover(id)
			} else {
				covered[id] = true
			}
			if role == "fallback" && String(binding, "known_role") != "" {
				unmapped = append(unmapped, Record{"role": binding["known_role"], "content_nodes": []string{id}, "reason": binding["reason"]})
			}
		}
		for _, node := range doc.Nodes {
			id := String(node, "id")
			if !covered[id] {
				unbound = append(unbound, id)
			}
			if String(node, "node_kind") == "opaque" {
				opaque = appendUnique(opaque, id)
			}
		}
		for _, element := range b.snapshot.Elements {
			if String(element, "source_document") == String(doc.Source, "id") && String(element, "kind") == "ExtensionSection" {
				extensions = append(extensions, String(element, "id"))
			}
		}
		docStatus := "complete_for_profile"
		limitations := []string{}
		if String(doc.Document, "domain_decomposition") == "failed" {
			docStatus = "partial"
			limitations = append(limitations, "The authoring parser rejected this document; domain extraction is unavailable.")
		}
		if len(unbound) != 0 || len(unmapped) != 0 {
			docStatus = "partial"
		}
		if docStatus != "complete_for_profile" {
			status = "partial"
		}
		for _, list := range [][]string{unbound, opaque, extensions} {
			sort.Strings(list)
		}
		doc.Document["coverage"] = Record{"status": docStatus, "mapping_profile": MappingProfile, "unbound_content_nodes": unbound, "unmapped_known_roles": unmapped, "extension_sections": extensions, "uninterpreted_nodes": opaque, "limitations": limitations}
		allUnbound = append(allUnbound, unbound...)
		allUnmapped = append(allUnmapped, unmapped...)
		allExtensions = append(allExtensions, extensions...)
		allOpaque = append(allOpaque, opaque...)
	}
	for _, list := range [][]string{allUnbound, allExtensions, allOpaque} {
		sort.Strings(list)
	}
	b.snapshot.Coverage = Record{"status": status, "mapping_profile": MappingProfile, "unbound_content_nodes": allUnbound, "unmapped_known_roles": allUnmapped, "extension_sections": allExtensions, "uninterpreted_nodes": allOpaque}
}
