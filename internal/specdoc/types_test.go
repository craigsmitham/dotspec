package specdoc

import (
	"encoding/json"
	"slices"
	"testing"

	schema "github.com/craigsmitham/dotspec/model"
)

func TestDocumentRolesMatchSupportingModel(t *testing.T) {
	var catalog struct {
		Enums map[string][]string
	}
	if err := json.Unmarshal(schema.CatalogJSON, &catalog); err != nil {
		t.Fatal(err)
	}
	modelRoles := catalog.Enums["SpecificationDocumentType"]
	pageRoles := []string{RequirementsSpecification, TableOfContents, CollectionIndex, SpecificationSection, RequirementsReviewReport}
	slices.Sort(modelRoles)
	slices.Sort(pageRoles)
	if !slices.Equal(pageRoles, modelRoles) {
		t.Fatalf("document roles %v differ from model %v", pageRoles, modelRoles)
	}
}
