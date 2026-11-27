package model

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/specfile"
)

// A component names its parent and a system states its purpose.
const (
	componentPart = "## Part of\n\n[Parent](../parent.system.spec.md)\n\n"
	systemPurpose = "## Purpose\n\nIt serves.\n\n"
)

// headingHosts place a heading where checkFields reads it for each class: at
// level two in its own specification, or at level four under a named
// collection entry. AdoptedProvision is absent because its entries admit
// exactly Edition and Applicability, a separate structural rule.
var headingHosts = []struct {
	class, kind, before string
	level               int
}{
	{"BusinessRule", "business-rule", "", 2},
	{"Capability", "capability", "", 2},
	{"Component", "component", componentPart, 2},
	{"Concept", "concept", "", 2},
	{"Feature", "feature", "", 2},
	{"Job", "job", "", 2},
	{"SystemQualityAttribute", "quality-attribute", "", 2},
	{"QualityModel", "quality-model", "", 2},
	{"Requirement", "requirement", "", 2},
	{"System", "system", systemPurpose, 2},
	{"UseCase", "use-case", "", 2},
	{"JobStep", "job", "## Steps\n\n### Step\n\nDo it.\n\n", 4},
	{"DesiredOutcome", "job", "## Outcomes\n\n### Outcome\n\nA result.\n\n", 4},
	{"Assumption", "system", systemPurpose + "## Assumptions\n\n### Premise\n\nIt holds.\n\n", 4},
	{"ConceptAction", "concept", "## Actions\n\n### act\n\nActs.\n\n", 4},
	{"ConceptUse", "system", systemPurpose + "## Concept uses\n\n### Use\n\nUses it.\n\n", 4},
	{"Synchronization", "system", systemPurpose + "## Synchronizations\n\n### Sync\n\nSyncs.\n\n", 4},
	{"Interface", "system", systemPurpose + "## Required interfaces\n\n### Port\n\nA port.\n\n", 4},
	{"Interface", "component", componentPart + "## Provided interfaces\n\n### Port\n\nA port.\n\n", 4},
	{"UseCaseExtension", "use-case", "## Extensions\n\n### 1a. Fails\n\nIt fails.\n\n", 4},
}

// The component parser rejects these accepted headings for reasons of its own.
var kindRejections = map[string]string{"Component\tSpecification system": "structure.part-of"}

// Check and export read the same heading vocabulary: every heading the
// specfile check accepts maps to a domain field, and every other heading is
// rejected as an unknown section.
func TestCheckAndExportAcceptTheSameHeadings(t *testing.T) {
	golden, err := os.ReadFile("../specfile/testdata/section-headings.tsv")
	if err != nil {
		t.Fatal(err)
	}
	accepted := map[string]bool{}
	var titles []string
	for _, line := range strings.Split(strings.TrimSpace(string(golden)), "\n") {
		class, title, _ := strings.Cut(line, "\t")
		accepted[class+"\t"+title] = true
		if !slices.Contains(titles, title) {
			titles = append(titles, title)
		}
	}
	for _, host := range headingHosts {
		for _, title := range titles {
			want := accepted[host.class+"\t"+title]
			if _, ok := specfile.SectionSlot(host.class, title); ok != want {
				t.Fatalf("SectionSlot(%s, %q) = %v, golden %v", host.class, title, ok, want)
			}
			source := "---\ntype: " + documentType(host.kind) + "\n---\n\n# Example\n\n" + host.before
			// A required section the host already declares is tested in place.
			if host.level > 2 || !strings.Contains(host.before, "## "+title+"\n") {
				source += strings.Repeat("#", host.level) + " " + title + "\n\n" + headingBody(title, host.level+1)
			}
			d := ParseContent("specs/example"+specfile.AuthoredSuffix(host.kind), []byte(source))
			m := MapDomain(d, host.kind)
			checked, exported := true, true
			for _, diag := range m.Diagnostics {
				code := String(diag, "code")
				switch {
				case code == kindRejections[host.class+"\t"+title]:
					checked, exported = want, want
				case code != "structure.unknown-section":
					t.Fatalf("%s %q: unexpected diagnostic %v", host.class, title, diag)
				case diag["record"] == nil:
					checked = false
				default:
					exported = false
				}
			}
			if checked != want {
				t.Errorf("%s %q: check accepts %v, golden %v", host.class, title, checked, want)
			}
			if kindRejections[host.class+"\t"+title] == "" {
				exported = checked && exported && headingMapped(d, d.Sections, title, host.level)
			}
			if checked != exported {
				t.Errorf("%s %q in %s: check accepts %v, export maps %v", host.class, title, host.kind, checked, exported)
			}
		}
	}
}

// A reusable quality model is a reference root rather than a specification.
func documentType(kind string) string {
	if kind == "quality-model" {
		return kind
	}
	return kind + "-spec"
}

func headingBody(title string, level int) string {
	entry := strings.Repeat("#", level) + " Entry\n\n"
	switch {
	case title == "Adopts":
		detail := strings.Repeat("#", level+1)
		return entry + "[Provision](https://example.com/spec#clause)\n\n" + detail + " Edition\n\nevolving\n\n" + detail + " Applicability\n\nAlways.\n"
	case specfile.IsCollectionHeading(title):
		return entry + "**ID:** entry\n\nText.\n"
	}
	// Structural sections such as Part of link to exactly one subject.
	return "See [the subject](../subject.system.spec.md).\n"
}

func headingMapped(d *ParsedDocument, sections []*Section, title string, level int) bool {
	for _, s := range sections {
		if s.Title == title && s.Level == level {
			return d.NodeByID[s.HeadingID]["domain_field"] != nil
		}
		if headingMapped(d, s.Children, title, level) {
			return true
		}
	}
	return false
}
