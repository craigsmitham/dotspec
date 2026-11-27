package specfile

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestRootKindInvertsDocumentRootClass(t *testing.T) {
	for kind, class := range documentRootClasses {
		if got, ok := RootKind(class); !ok || got != kind {
			t.Errorf("RootKind(%s) = %q, %v; want %q", class, got, ok, kind)
		}
	}
	if kind, ok := RootKind("ConceptUse"); ok {
		t.Errorf("ConceptUse has no authoring form, got %q", kind)
	}
}

func TestReportTypeReadsTheCatalog(t *testing.T) {
	for kind, want := range map[string]string{
		"System":                 "system-spec-report",
		"SystemQualityAttribute": "system-quality-attribute-spec-report",
		"QualityModel":           "quality-model-report",
		"QualityModelAttribute":  "quality-model-attribute-report",
		"AdoptedProvision":       "adopted-provision-report",
		"DomainElement":          "",
	} {
		if got := ReportType(kind); got != want {
			t.Errorf("ReportType(%s) = %q, want %q", kind, got, want)
		}
	}
}

// The golden list pins every (class, heading) pair checkFields accepts, so a
// change to the heading convention or its exceptions shows up as a diff.
func TestSectionSlotAcceptsThePinnedHeadings(t *testing.T) {
	golden, err := os.ReadFile("testdata/section-headings.tsv")
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Split(strings.TrimSpace(string(golden)), "\n")
	titles := map[string]bool{"Name": true, "Preamble": true, "Statement": true, "Steps": true, "Quality attributes": true, "Specification system": true}
	for _, line := range want {
		titles[strings.SplitN(line, "\t", 2)[1]] = true
	}
	for title := range irregularHeadings {
		titles[title] = true
	}
	for _, slots := range classSlots() {
		for slot := range slots.Slots {
			title := sentenceCase(slot)
			titles[title], titles[slot], titles[strings.ToLower(title)], titles[strings.ToUpper(title)] = true, true, true, true
		}
	}
	var got []string
	for class := range classSlots() {
		for title := range titles {
			if _, ok := SectionSlot(class, title); ok {
				got = append(got, class+"\t"+title)
			}
		}
	}
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Errorf("accepted headings changed:\ngot  %q\nwant %q", got, want)
	}
}

func TestSectionSlotFollowsTheConvention(t *testing.T) {
	for _, c := range []struct{ class, title, slot string }{
		{"System", "Out of scope", "out_of_scope"},
		{"System", "Stakeholders", "stakeholder_roles"},
		{"UseCase", "System boundary", "boundary_description"},
		{"Job", "Job", "statement"},
	} {
		if slot, ok := SectionSlot(c.class, c.title); !ok || slot != c.slot {
			t.Errorf("SectionSlot(%s, %q) = %q, %v; want %q", c.class, c.title, slot, ok, c.slot)
		}
	}
	for _, c := range []struct{ class, title string }{
		{"System", "Name"}, {"System", "purpose"}, {"System", "Out Of Scope"}, {"System", "out_of_scope"}, {"Goal", "Statement"}, {"Requirement", "Quality attributes"}, {"Scenario", "Steps"}, {"Widget", "Description"},
	} {
		if slot, ok := SectionSlot(c.class, c.title); ok {
			t.Errorf("SectionSlot(%s, %q) accepted slot %q", c.class, c.title, slot)
		}
	}
}
