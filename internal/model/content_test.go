package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestContentPreservesOriginalAndOrderedTree(t *testing.T) {
	source := "\ufeff---\r\n# metadata comment\r\ntype: use-case\r\nid: sample\r\nlabels: [one, two]\r\nnested:\r\n  key: \"a # quoted\" # tail\r\n  enabled: true\r\n---\r\n# A *title*\r\n\r\nPreamble.\r\n\r\n## Scenario\r\n\r\n1. [x] **ID:** first\r\n   **Stage:** explain\r\n   - nested [link](other.md#part)\r\n\r\n### Arbitrary child\r\n\r\n| Head | Other |\r\n| :--- | ---: |\r\n| ~~strike~~ | ![picture](x.png) |\r\n\r\n> ## Quoted, not a section\r\n> next\r\n\r\n```mermaid custom\r\nflowchart LR\r\nA-->B\r\n## Not a section\r\n```\r\n\r\n<!-- author comment -->\r\n\r\n[shown][r] <https://example.test> [^note]\r\n\r\n[r]: other.md \"A title\"\r\n\r\n[^note]: Footnote text.\r\n\r\n## Empty\r\n"
	p := ParseContent("specs/a.use-case.spec.md", []byte(source))
	payload := Object(p.Source["payload"])
	decoded, err := base64.StdEncoding.DecodeString(String(payload, "bytes_base64"))
	if err != nil || !bytes.Equal(decoded, []byte(source)) {
		t.Fatalf("payload mismatch: %v", err)
	}
	if String(payload, "sha256") != fmt.Sprintf("%x", sha256.Sum256([]byte(source))) {
		t.Fatal("digest mismatch")
	}
	if p.Title != "A title" {
		t.Fatalf("title %q", p.Title)
	}
	if len(p.Sections) != 2 || len(p.Sections[0].Children) != 1 || p.Sections[0].Children[0].Title != "Arbitrary child" {
		t.Fatalf("sections %#v", p.Sections)
	}
	if len(p.Sections[1].Content) != 0 || String(p.Rich(p.Sections[1].Content), "presence") != "empty" {
		t.Fatal("blank section gained body")
	}
	kinds := map[string]int{}
	for _, n := range p.Nodes {
		kinds[String(n, "node_kind")]++
	}
	for _, kind := range []string{"frontmatter", "yaml_mapping", "yaml_sequence", "yaml_scalar", "yaml_entry", "heading", "list", "list_item", "strong", "table", "table_row", "table_cell", "strikethrough", "image", "block_quote", "fenced_code", "comment", "link", "autolink", "link_definition", "footnote_definition", "footnote_reference", "trivia", "soft_break"} {
		if kinds[kind] == 0 {
			t.Errorf("missing %s", kind)
		}
	}
	var foundFence bool
	for _, n := range p.Nodes {
		if String(n, "node_kind") == "fenced_code" {
			foundFence = true
			if String(n, "code_info") != "mermaid custom" || !strings.HasSuffix(p.Raw(String(n, "id")), "```\r\n") {
				t.Errorf("fence metadata/span: %#v %q", n, p.Raw(String(n, "id")))
			}
		}
	}
	if !foundFence {
		t.Fatal("missing fence")
	}
	if len(p.References) != 5 {
		t.Errorf("references %d: %#v", len(p.References), p.References)
	}
	checkContentTree(t, p)
	// Every normalized byte boundary maps back to exactly its original bytes;
	// the parser offset map records the BOM and each CRLF replacement explicitly.
	segments, ok := p.Document["normalization"].([]Record)
	if !ok || len(segments) < 2 {
		t.Fatal("normalization map missing")
	}
	for _, s := range segments {
		a, b := Number(s["original_start"]), Number(s["original_end"])
		if String(s, "transformation") == "CRLF to LF" && string(p.Data[a:b]) != "\r\n" {
			t.Fatal("bad CRLF map")
		}
	}
	first, _ := json.Marshal(p)
	second, _ := json.Marshal(ParseContent("specs/a.use-case.spec.md", []byte(source)))
	if !bytes.Equal(first, second) {
		t.Fatal("nondeterministic extraction")
	}
}
func checkContentTree(t *testing.T, p *ParsedDocument) {
	t.Helper()
	seen := map[string]bool{}
	var visit func(string)
	visit = func(id string) {
		if seen[id] {
			t.Errorf("node multiply owned: %s", id)
			return
		}
		seen[id] = true
		n := p.NodeByID[id]
		if n == nil {
			t.Fatalf("missing node %s", id)
		}
		span := Object(n["span"])
		a, b := Number(span["start"]), Number(span["end"])
		if a < 0 || b < a || b > len(p.Data) {
			t.Errorf("invalid span %s %d:%d", id, a, b)
		}
		at := a
		for _, childID := range Strings(n["children"]) {
			child := p.NodeByID[childID]
			cs := Object(child["span"])
			c, d := Number(cs["start"]), Number(cs["end"])
			if c < at || d > b {
				t.Errorf("child overlap/outside %s [%d,%d) parent %s [%d,%d) previous end %d raw=%q", childID, c, d, id, a, b, at, p.Raw(childID))
			}
			at = max(at, d)
			visit(childID)
		}
	}
	visit(String(p.Document, "root"))
	if len(seen) != len(p.Nodes) {
		for _, n := range p.Nodes {
			if !seen[String(n, "id")] {
				t.Errorf("orphan node %s", String(n, "id"))
			}
		}
	}
}
func TestContentMalformedYAMLRetainsMarkdown(t *testing.T) {
	p := ParseContent("broken.spec.md", []byte("---\ntype: [broken\n---\n# Kept\n\n## Goal\nStill readable.\n"))
	if String(p.Document, "frontmatter_parse_status") != "failed" || String(p.Document, "parse_status") != "partial" || len(p.Diagnostics) != 1 {
		t.Fatalf("parse result %#v", p.Document)
	}
	if p.Title != "Kept" || len(p.Sections) != 1 {
		t.Fatal("Markdown after bad metadata was discarded")
	}
	checkContentTree(t, p)
}
func TestContentMissingDelimiterAndInvalidUTF8(t *testing.T) {
	for _, source := range [][]byte{[]byte("---\ntype: job\n# never closed"), {'#', ' ', 0xff, '\n'}} {
		p := ParseContent("bad.md", source)
		if len(p.Diagnostics) == 0 {
			t.Fatal("missing diagnostic")
		}
		payload, _ := base64.StdEncoding.DecodeString(String(Object(p.Source["payload"]), "bytes_base64"))
		if !bytes.Equal(payload, source) {
			t.Fatal("original bytes changed")
		}
		checkContentTree(t, p)
	}
}
func TestContentSetextEscapesBreaksAndUndefinedReferences(t *testing.T) {
	p := ParseContent("native.md", []byte("Title\n=====\n\nBody\n----\n\nEscaped \\* &amp; next  \nline\\\nagain\nsoft\n\n[missing][nope] [^missing]\n\n`[not a link](x)`\n"))
	if p.Title != "Title" || len(p.Sections) != 1 || p.Sections[0].Title != "Body" {
		t.Fatal("setext headings not retained")
	}
	kinds := map[string]int{}
	for _, n := range p.Nodes {
		kinds[String(n, "node_kind")]++
	}
	for _, kind := range []string{"escape", "entity", "hard_break", "soft_break", "code_span"} {
		if kinds[kind] == 0 {
			t.Errorf("missing %s", kind)
		}
	}
	if len(p.References) != 2 {
		t.Fatalf("references: %#v", p.References)
	}
	checkContentTree(t, p)
}
func TestContentYAMLTypesAliasesAndDuplicateKeys(t *testing.T) {
	p := ParseContent("metadata.md", []byte("---\ntype: job\nitems: &items\n  - name: first\n    count: 2\n  - name: second\n    count: 3\ncopy: *items\ncustom: !tag value\nmultiline: |\n  first\n  second\nwhen: 2020-01-01\nnullish: null\ndup: one\ndup: two\n---\n# Job\n"))
	if String(p.Document, "frontmatter_parse_status") != "parsed" {
		t.Fatalf("failed metadata: %#v", p.Diagnostics)
	}
	seen := map[string]bool{}
	for _, n := range p.Nodes {
		if String(n, "scalar_kind") != "" {
			seen[String(n, "scalar_kind")] = true
		}
	}
	for _, kind := range []string{"string", "integer", "timestamp", "null", "custom"} {
		if !seen[kind] {
			t.Errorf("missing scalar kind %s", kind)
		}
	}
	checkContentTree(t, p)
}

func TestContentUnusedFootnotesRetainSourceOrderAndParent(t *testing.T) {
	p := ParseContent("notes.md", []byte("# Title\n\n[^unused]: Definition **with emphasis**.\n\n> [^nested]: Nested definition.\n\nParagraph.\n"))
	count := 0
	for _, n := range p.Nodes {
		if String(n, "node_kind") == "footnote_definition" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("definitions = %d", count)
	}
	checkContentTree(t, p)
}

func TestContentLexicalEdges(t *testing.T) {
	cases := []string{
		"# T\n\n` foo ` and `` a ` b `` and [`]`](url \"title)\")\n",
		"---\ntype: job\ntext: several words\n  continuing here\nempty:\nlist: [one, two, {x: 2, y: null}]\n---\n# T\n",
		"---\n? [one, two]\n: complex key\na: &value 1\nb: *value\nc: !!str 12\n---\n# T\n",
		"# T\n\n> ~~~go\n> sample\n> ~~~\n\n- ```js\n  x()\n  ```\n",
	}
	for i, s := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) { p := ParseContent("edge.md", []byte(s)); checkContentTree(t, p) })
	}
}

func FuzzParseContent(f *testing.F) {
	for _, seed := range []string{"", "# T\n\nText.\n", "0\n0\n-", "**0*", "0*~~0~ ", "[^a]:0\n0!0^a]", "*0***0*", "~00~0", "*\x1c**", "0[]( \")\")0", "\ufeff---\r\ntype: job\r\n---\r\n# T\r\n", "- [x] **a**\n  - b\n", "---\nx: [a, b]\n---\n", "[^a]: Definition\n\nUse[^a].\n", "[a](b \"c)\")\n"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		p := ParseContent("fuzz.md", data)
		checkContentTree(t, p)
		encoded, err := json.Marshal(p.Source)
		if err != nil || len(encoded) == 0 {
			t.Fatalf("source not JSON serializable: %v", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(String(Object(p.Source["payload"]), "bytes_base64"))
		if err != nil || !bytes.Equal(decoded, data) {
			t.Fatal("byte retention failed")
		}
	})
}

func TestContentYAMLScalarLexicalSpans(t *testing.T) {
	cases := []struct{ name, body, value, raw, comment string }{
		{"plain-final-continuation", "extra: first\n  next\n  last\n", "first next last", "first\n  next\n  last", ""},
		{"plain-before-key", "extra: first\n  next\n  last\nother: value\n", "first next last", "first\n  next\n  last", ""},
		{"plain-comma", "extra: one, two\n", "one, two", "one, two", ""},
		{"literal-clip", "extra: |\n  first\n\n  second\n# trailing\n", "first\n\nsecond\n", "|\n  first\n\n  second\n", "# trailing\n"},
		{"literal-keep", "extra: |+\n  first\n\n\n# trailing\n", "first\n\n\n", "|+\n  first\n\n\n", "# trailing\n"},
		{"folded-strip", "extra: >-\n  first\n  second\n\n  third\n\n# trailing\n", "first second\nthird", ">-\n  first\n  second\n\n  third\n\n", "# trailing\n"},
		{"explicit-indent-literal-hash", "extra: |2-\n  # literal\n  body\n# outside\n", "# literal\nbody", "|2-\n  # literal\n  body\n", "# outside\n"},
	}
	for _, test := range cases {
		for _, crlf := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/crlf=%v", test.name, crlf), func(t *testing.T) {
				source := "---\ntype: concept\n" + test.body + "---\n# T\n"
				wantRaw, wantComment := test.raw, test.comment
				if crlf {
					source = "\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")
					wantRaw = strings.ReplaceAll(wantRaw, "\n", "\r\n")
					wantComment = strings.ReplaceAll(wantComment, "\n", "\r\n")
				}
				p := ParseContent("scalar.md", []byte(source))
				if len(p.Diagnostics) != 0 {
					t.Fatalf("syntax diagnostics: %#v", p.Diagnostics)
				}
				found, comment := false, wantComment == ""
				for _, n := range p.Nodes {
					if String(n, "node_kind") == "yaml_scalar" && String(n, "decoded_text") == test.value {
						found = true
						if raw := p.Raw(String(n, "id")); raw != wantRaw {
							t.Errorf("scalar span raw %q, want %q", raw, wantRaw)
						}
					}
					if String(n, "node_kind") == "comment" && p.Raw(String(n, "id")) == wantComment {
						comment = true
					}
				}
				if !found || !comment {
					t.Errorf("scalar found=%v, separate dedented comment found=%v", found, comment)
				}
				checkContentTree(t, p)
			})
		}
	}
}

func TestContentYAMLFlowCommentCannotCloseCollection(t *testing.T) {
	p := ParseContent("flow.md", []byte("---\ntype: concept\nextra: [one, # ] is comment text\n  two, {inside: value}]\n---\n# T\n"))
	if len(p.Diagnostics) != 0 {
		t.Fatalf("syntax diagnostics: %#v", p.Diagnostics)
	}
	var sequence Record
	comment, inside := false, false
	for _, n := range p.Nodes {
		switch String(n, "node_kind") {
		case "yaml_sequence":
			sequence = n
		case "comment":
			comment = comment || strings.Contains(p.Raw(String(n, "id")), "] is comment")
		case "yaml_scalar":
			inside = inside || String(n, "decoded_text") == "value"
		}
	}
	if sequence == nil || p.Raw(String(sequence, "id")) != "[one, # ] is comment text\n  two, {inside: value}]" || !comment || !inside {
		t.Fatalf("flow representation incomplete: %#v comment=%v inner scalar=%v", sequence, comment, inside)
	}
	checkContentTree(t, p)
}

func TestContentFenceSpanUsesActualClosingLine(t *testing.T) {
	cases := []struct{ name, source, raw, body string }{
		{"quoted-looking-literal", "```md\n> ```\nkept\n```\n", "```md\n> ```\nkept\n```\n", "> ```\nkept\n"},
		{"fence-looking-suffix", "```md\n``` trailing\nkept\n```\n", "```md\n``` trailing\nkept\n```\n", "``` trailing\nkept\n"},
		{"quoted-fence", "> ```go\n> x\n> ```\n", "```go\n> x\n> ```\n", "x\n"},
		{"list-fence", "- ```go\n  x\n  ```\n", "```go\n  x\n  ```\n", "x\n"},
		{"unclosed-quote", "> ```go\n> x\n\noutside\n", "```go\n> x\n", "x\n"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			p := ParseContent("fence.md", []byte(test.source))
			found := false
			for _, n := range p.Nodes {
				if String(n, "node_kind") == "fenced_code" {
					found = true
					if raw := p.Raw(String(n, "id")); raw != test.raw {
						t.Errorf("fence span raw %q, want %q", raw, test.raw)
					}
					if String(n, "decoded_text") != test.body {
						t.Errorf("body %q, want %q", String(n, "decoded_text"), test.body)
					}
				}
				if String(n, "node_kind") == "opaque" {
					t.Errorf("known fence content fell back to opaque: %q", p.Raw(String(n, "id")))
				}
			}
			if !found {
				t.Fatal("no code fence")
			}
			checkContentTree(t, p)
		})
	}
}

func TestContentYAMLBlockHeaderCommentIsAnOccurrence(t *testing.T) {
	p := ParseContent("header.md", []byte("---\ntype: job\nextra: | # header\n  body\n# trailing\n---\n# Job\n"))
	comments := map[string]bool{}
	for _, n := range p.Nodes {
		if String(n, "node_kind") == "comment" {
			comments[p.Raw(String(n, "id"))] = true
		}
	}
	if !comments["# header\n"] || !comments["# trailing\n"] {
		t.Fatalf("comments: %#v", comments)
	}
	checkContentTree(t, p)
}

func TestContentFrontmatterDelimiterMustBeUnindented(t *testing.T) {
	source := "---\ntype: concept\nextra: |\n  ---\n  retained\n  ...\n---\n# Real title\n"
	p := ParseContent("separator.md", []byte(source))
	if p.Title != "Real title" || len(p.Diagnostics) != 0 {
		t.Fatalf("title=%q diagnostics=%#v", p.Title, p.Diagnostics)
	}
	found := false
	for _, n := range p.Nodes {
		if String(n, "node_kind") == "yaml_scalar" && String(n, "decoded_text") == "---\nretained\n...\n" {
			found = true
			if raw := p.Raw(String(n, "id")); raw != "|\n  ---\n  retained\n  ...\n" {
				t.Errorf("scalar raw %q", raw)
			}
		}
	}
	if !found {
		t.Fatal("indented delimiter split the frontmatter scalar")
	}
	checkContentTree(t, p)
}

func TestContentMetadataAliasesRetainSyntaxAndInterpretation(t *testing.T) {
	for _, source := range []string{"---\ntype: &kind concept\nid: *kind\n---\n# T\n", "---\nvalue: &kind concept\ntype: *kind\nid: *kind\n---\n# T\n"} {
		p := ParseContent("alias.md", []byte(source))
		if String(p.Document, "declared_type") != "concept" || String(p.Document, "authored_id") != "concept" {
			t.Fatalf("interpreted metadata: %#v", p.Document)
		}
		aliases := 0
		for _, n := range p.Nodes {
			if String(n, "node_kind") == "yaml_alias" {
				aliases++
				if String(n, "alias_name") != "kind" || p.Raw(String(n, "id")) != "*kind" {
					t.Errorf("alias occurrence rewritten: %#v raw=%q", n, p.Raw(String(n, "id")))
				}
			}
		}
		if aliases == 0 {
			t.Fatal("aliases lost")
		}
		checkContentTree(t, p)
	}
}
