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
		class, tag := kind+"SpecReport", kindName(kind)+"-spec-report"
		if kind == "QualityModel" {
			class, tag = "QualityModelReport", "quality-model-report"
		}
		if kind == "SystemQualityAttribute" {
			tag = "system-quality-attribute-spec-report"
		}
		fields := catalog.Classes[class].Slots
		if fields["type"].Fixed != tag || !slices.Equal(fields["documents"].Ranges, []string{kind}) {
			t.Errorf("%s must document one %s", class, kind)
		}
	}
	for kind := range families {
		check(kind)
	}
}
