package model

import (
	"strings"
	"testing"
)

func mapped(t *testing.T, kind, body string) (*ParsedDocument, Mapping, Record) {
	t.Helper()
	d := ParseContent("specs/example."+kind+".spec.md", []byte("---\ntype: "+kind+"-spec\n---\n\n# Example\n\n"+body))
	m := MapDomain(d, kind)
	if len(m.Diagnostics) > 0 {
		t.Fatalf("unexpected mapping diagnostics: %#v", m.Diagnostics)
	}
	for _, r := range m.Elements {
		if String(r, "id") == m.RootID {
			return d, m, r
		}
	}
	t.Fatal("no mapped root")
	return nil, Mapping{}, nil
}

func recordsOf(m Mapping, kind string) []Record {
	var out []Record
	for _, r := range m.Elements {
		if String(r, "kind") == kind {
			out = append(out, r)
		}
	}
	return out
}

func richText(doc *ParsedDocument, value any) string {
	var out []string
	for _, id := range Strings(Object(value)["content_nodes"]) {
		out = append(out, doc.Text(id))
	}
	return strings.Join(out, "")
}

func TestNamedOutcomeAndNarrativeFields(t *testing.T) {
	d, m, root := mapped(t, "job", "## Outcomes\n### Earlier detection\n**ID:** early\nDetect failures before release.\n")
	outcome := recordsOf(m, "DesiredOutcome")[0]
	if String(outcome, "authored_id") != "early" || richText(d, outcome["statement"]) != "Detect failures before release." || len(Strings(root["outcomes"])) != 1 {
		t.Fatal(outcome)
	}
	d, _, root = mapped(t, "concept", "## State\n### Counter\nOne per member.\n## Effects table\n| Read | Result |\n| --- | --- |\n| get | unchanged |\n")
	if !strings.Contains(richText(d, root["abstract_state"]), "Counter") || !strings.Contains(richText(d, root["effects_table"]), "unchanged") {
		t.Fatal(root)
	}
}
func TestInvalidCollectionsAndUnknownFieldsHaveNoDomainMeaning(t *testing.T) {
	for _, body := range []string{"## Notes\nKeep me.", "## Outcomes\nLoose prose.\n### Outcome\nResult.", "## Outcomes\n- A bullet", "## Steps\n### First\n#### Typo\nText."} {
		d := ParseContent("draft.job.spec.md", []byte("---\ntype: job-spec\n---\n# Draft\n"+body))
		m := MapDomain(d, "job")
		if len(m.Diagnostics) == 0 || len(m.Elements) != 0 || !strings.Contains(String(d.Source, "markdown"), body) {
			t.Fatal(body, m)
		}
	}
}
