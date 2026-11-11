package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the mapping/content seam independently of filesystem discovery and
// reference resolution. Unattempted references are valid in this partial view.
func TestMappedCorpusGraphIntegrity(t *testing.T) {
	snapshot := NewSnapshot(Record{
		"roots": []string{"../.."}, "selected_spec_kinds": []string{"system", "job", "use-case", "concept", "requirement", "quality", "business-rule", "quality-model"},
		"discovery_complete": true, "reads_complete": true, "inspection_complete": true,
		"discovery_policy":    Record{"mode": "explicit", "rule": "mapping-test/1", "follow_symlinks": false, "include_nested_repositories": false},
		"reference_traversal": Record{"mode": "none", "include_system_context": false, "follow_external_resources": false, "repository_bounded": true, "rule": "mapping-test/1"},
	})
	snapshot.Coverage["status"] = "complete_for_profile"
	for _, root := range []string{"../../specs", "../../examples/starter/specs", "../../internal/integration/specs", "../../internal/qualitymodel/specs"} {
		if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".spec.md") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			doc := ParseContent(path, data)
			mapped := MapDomain(doc, String(doc.Document, "declared_type"))
			Object(doc.Document["coverage"])["status"] = "complete_for_profile"
			elementIDs := []string{}
			for _, element := range mapped.Elements {
				elementIDs = append(elementIDs, String(element, "id"))
			}
			doc.Document["elements"] = elementIDs
			snapshot.Sources = append(snapshot.Sources, doc.Source)
			snapshot.Documents = append(snapshot.Documents, doc.Document)
			snapshot.ContentNodes = append(snapshot.ContentNodes, doc.Nodes...)
			snapshot.Fragments = append(snapshot.Fragments, doc.Fragments...)
			snapshot.References = append(snapshot.References, doc.References...)
			snapshot.Elements = append(snapshot.Elements, mapped.Elements...)
			snapshot.LabelOccurrences = append(snapshot.LabelOccurrences, mapped.Labels...)
			snapshot.ContentBindings = append(snapshot.ContentBindings, mapped.Bindings...)
			snapshot.Diagnostics = append(snapshot.Diagnostics, mapped.Diagnostics...)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	findings, err := Validate(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range findings {
		t.Errorf("%s %s: %s", finding.Code, finding.Path, finding.Message)
	}
}
