package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestReferencedMarkdownContributesOnlyTextAndAnchors(t *testing.T) {
	p := parseMarkdownContent("guide.md", []byte("\ufeff# Guide\r\n## Local heading\r\n[more](unread.md)\r\n"), false)
	if len(p.Fragments) != 2 || len(p.References) != 0 || len(p.Nodes) != 0 || !strings.Contains(String(p.Source, "markdown"), "[more](unread.md)") {
		t.Fatal(p)
	}
	span := Object(p.Fragments[1]["span"])
	if text := string(p.Data[Number(span["start"]):Number(span["end"])]); text != "## Local heading\n" {
		t.Fatal(text)
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
func TestContentMissingDelimiterAndInvalidUTF8(t *testing.T) {
	for _, source := range [][]byte{[]byte("---\ntype: job-spec\n# never closed"), {'#', ' ', 0xff, '\n'}} {
		p := ParseContent("bad.md", source)
		if len(p.Diagnostics) == 0 {
			t.Fatal("missing diagnostic")
		}
		checkContentTree(t, p)
	}
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
		"---\ntype: job-spec\ntext: several words\n  continuing here\nempty:\nlist: [one, two, {x: 2, y: null}]\n---\n# T\n",
		"---\n? [one, two]\n: complex key\na: &value 1\nb: *value\nc: !!str 12\n---\n# T\n",
		"# T\n\n> ~~~go\n> sample\n> ~~~\n\n- ```js\n  x()\n  ```\n",
	}
	for i, s := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) { p := ParseContent("edge.md", []byte(s)); checkContentTree(t, p) })
	}
}

func FuzzParseContent(f *testing.F) {
	for _, seed := range []string{"", "# T\n\nText.\n", "0\n0\n-", "**0*", "0*~~0~ ", "[^a]:0\n0!0^a]", "*0***0*", "~00~0", "*\x1c**", "0[]( \")\")0", "\ufeff---\r\ntype: job-spec\r\n---\r\n# T\r\n", "- [x] **a**\n  - b\n", "---\nx: [a, b]\n---\n", "[^a]: Definition\n\nUse[^a].\n", "[a](b \"c)\")\n"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		p := ParseContent("fuzz.md", data)
		checkContentTree(t, p)
		encoded, err := json.Marshal(p.Source)
		if err != nil || len(encoded) == 0 {
			t.Fatalf("source not JSON serializable: %v", err)
		}

	})
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
