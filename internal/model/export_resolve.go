package model

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

func (b *snapshotBuilder) resolution(occurrence Record, status, target, rule string) Record {
	r := Record{"status": status, "rule": rule}
	if target != "" {
		r["target"] = target
	}
	occurrence["resolution"] = r
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
			reference["authored_target"] = target
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
		if b.glossaries[path] != nil && String(node, "node_kind") == "link" {
			switch String(reference, "context_field") {
			case "relationships.references", "related_terms.references":
				mayFollow = true
			}
		}
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
	if reference["role"] != nil {
		return
	}
	span := Object(doc.NodeByID[String(reference, "content_node")]["span"])
	best := int(^uint(0) >> 1)
	link := doc.NodeByID[String(reference, "content_node")]
	if link["link_owner"] != nil {
		reference["context_owner"], reference["context_field"] = link["link_owner"], link["link_field"]
		return
	}
	for _, node := range doc.Nodes {
		if node["field_owner"] == nil {
			continue
		}
		bound := Object(node["span"])
		width := Number(bound["end"]) - Number(bound["start"])
		if Number(bound["start"]) <= Number(span["start"]) && Number(bound["end"]) >= Number(span["end"]) && width < best {
			best = width
			reference["context_owner"], reference["context_field"] = node["field_owner"], node["field_path"]
		}
	}
}

func (b *snapshotBuilder) associateAttributeOwners() {
	for _, file := range sortedKeys(b.attributeOwners) {
		r := b.elements[b.rootByPath[file]]
		sys := b.elements[b.rootByPath[b.attributeOwners[file]]]
		if r == nil || sys == nil {
			continue
		}
		r["specification_system"] = sys["id"]
		sys["quality_attributes"] = appendUnique(Strings(sys["quality_attributes"]), String(r, "id"))
	}
}

func (b *snapshotBuilder) resolveLabels() {
	// Direct targets may be associated with already acquired system context,
	// but no additional context is read merely to complete their classifications.
	for _, path := range sortedKeys(b.rootByPath) {
		if b.elements[b.rootByPath[path]] == nil {
			continue
		}
		if _, explicit := b.elements[b.rootByPath[path]]["specification_system_content"]; !explicit && (b.primary[path] || b.contextOK[path]) {
			continue
		}
		kind := b.kind(path)
		if kind == nil || kind.Context == nil {
			continue
		}
		owner := kind.Layout.ContextDir(filepath.FromSlash(path))
		repo, err := specfile.ContextRoot(owner)
		if err != nil {
			continue
		}
		for current := owner; ; current = filepath.Dir(current) {
			file, err := system.Layout.LocalSubject(current)
			if err != nil {
				break
			}
			candidate := cleanPath(file)
			if source := b.sources[candidate]; source != nil {
				b.contextPath[path] = candidate
				b.contextOK[path] = b.rootByPath[candidate] != ""
				if sys := b.elements[b.rootByPath[candidate]]; String(sys, "kind") == "System" {
					b.directorySystem[b.rootByPath[path]] = String(sys, "id")
				}
				break
			}
			// An unacquired nearer system leaves the label uninspected, not unbound.
			if file != "" {
				break
			}
			if current == repo || filepath.Dir(current) == current {
				b.contextOK[path] = true
				break
			}
		}
	}
	b.associateSpecificationSystems()
	b.associateAttributeOwners()
	attributes := map[string]map[string][]string{}
	for _, element := range b.snapshot.Elements {
		if !designKind(String(element, "kind"), "QualityAttribute") {
			continue
		}
		owner, label := qualityScope(element), String(element, "authored_id")
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
			sys = qualityScope(owner)
		} else {
			sys = governingSystem(owner, b.elements)
			if sys == "" && !b.contextOK[path] {
				b.resolution(label, "not_attempted", "", "uninspected-quality-context/1")
				continue
			}
		}
		if sys != "" {
			if String(b.elements[sys], "kind") == "QualityModel" {
				label["governing_quality_model"] = sys
			} else {
				label["governing_system"] = sys
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
			field := "quality_attributes"
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

// Relationships derives navigation from the resolved domain fields.
func (s *Snapshot) Relationships() []Record {
	s = s.navigation()
	var relationships []Record
	for _, projection := range compiledModel().catalog.Projections {
		parent, field, _ := strings.Cut(projection, ".")
		for _, element := range s.Elements {
			if designKind(String(element, "kind"), parent) {
				for _, target := range valuesAt(element, field) {
					relationships = append(relationships, Record{"projection_path": projection, "subject": element["id"], "object": target})
				}
			}
		}
	}
	return relationships
}

// Observe acquired documents only. The fingerprint concerns the entire cited
// source, so a change prompts review even when the clause locator still resolves.
func (b *snapshotBuilder) observeAdoptions() {
	refs := map[string]Record{}
	for _, ref := range b.snapshot.References {
		refs[String(ref, "id")] = ref
	}
	digests := map[string]string{}
	for _, r := range b.snapshot.Elements {
		if r["kind"] != "AdoptedProvision" {
			continue
		}
		ref := refs[String(r, "provision")]
		u, err := url.Parse(String(ref, "authored_target"))
		if err != nil || u.IsAbs() || u.Host != "" || u.Path == "" || filepath.IsAbs(u.Path) {
			continue
		}
		source, _ := b.documentForSource(String(r, "source_document"))
		target := cleanPath(filepath.Join(filepath.Dir(source), filepath.FromSlash(u.Path)))
		acquired := b.sources[target]
		if acquired["read_status"] != "read" {
			continue
		}
		if digests[target] == "" {
			var data []byte
			if md, ok := acquired["markdown"].(string); ok {
				data = []byte(md)
			} else {
				data, err = b.inputs.Read(filepath.FromSlash(target))
				if err != nil {
					continue
				}
			}
			digests[target] = fmt.Sprintf("%x", sha256.Sum256(data))
		}
		r["observed_source_sha256"] = digests[target]
	}
}
