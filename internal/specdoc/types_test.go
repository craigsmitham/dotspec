package specdoc

import (
	"encoding/json"
	"slices"
	"testing"

	schema "github.com/craigsmitham/dotspec/model"
)

func TestSubjectDocumentsMirrorTheModel(t *testing.T) {
	type slot struct {
		Fixed  string   `json:"equals_string"`
		Ranges []string `json:"reference_ranges"`
	}
	var catalog struct {
		Classes map[string]struct{ Slots map[string]slot }
	}
	if err := json.Unmarshal(schema.CatalogJSON, &catalog); err != nil {
		t.Fatal(err)
	}
	check := func(kind string) {
		t.Helper()
		fields := catalog.Classes[kind+"SpecReport"].Slots
		if fields["type"].Fixed != kindName(kind)+"-spec-report" || !slices.Equal(fields["documents"].Ranges, []string{kind}) {
			t.Errorf("%sSpecReport must document one %s", kind, kind)
		}
	}
	for kind := range families {
		if kind != "SystemQualityAttribute" {
			check(kind)
		}
	}
	check("QualityAttribute")
}
