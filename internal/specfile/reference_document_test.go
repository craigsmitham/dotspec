package specfile

import "testing"

func TestReferenceExtractionUsesMarkdown(t *testing.T) {
	data := []byte("---\nlayout: example\n---\n## A heading\n## A heading\n[local][ref]\n\n[ref]: guide%20one.md#target\n\n`[not a link](absent.md)`\n\n```md\n# Not a heading\n[ignored](missing.md)\n<a id=\"example\"></a>\n```\n\n<a id=\"concept.checklist\"></a>\n")
	doc := MarkdownReferences(data)
	if len(doc.Links) != 1 || doc.Links[0].Destination != "guide%20one.md#target" || doc.Links[0].Location.Line != 6 || !doc.Anchors["a-heading"] || !doc.Anchors["a-heading-1"] || doc.Anchors["not-a-heading"] || !doc.Anchors["concept.checklist"] || doc.Anchors["example"] {
		t.Fatalf("extraction: %+v", doc)
	}
}

func TestNativeMarkdownWithOpeningThematicBreak(t *testing.T) {
	refs := NewMarkdown([]byte("---\n\n# Native heading\n")).AnchorReferences()
	if !refs.Anchors["native-heading"] {
		t.Fatal("a thematic break must preserve the native document's heading targets")
	}
}
