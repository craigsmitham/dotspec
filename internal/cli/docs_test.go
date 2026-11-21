package cli

import (
	"testing"

	"github.com/craigsmitham/dotspec/internal/guide"
)

func TestHandbookTopicsMatchSpecificationTypes(t *testing.T) {
	missing := map[string]bool{}
	for _, typ := range supportedTypes() {
		missing[typ.Layout.Kind] = true
	}
	for _, topic := range guide.Topics() {
		if !missing[topic.Name] {
			t.Errorf("handbook topic %q has no specification type", topic.Name)
		}
		delete(missing, topic.Name)
		forms := map[string]bool{}
		for _, page := range topic.Pages {
			forms[page.Form()] = true
		}
		if !forms["Explanation"] || !forms["How-to"] {
			t.Errorf("topic %q needs an explanation and an authoring guide", topic.Name)
		}
	}
	for kind := range missing {
		t.Errorf("specification type %q has no handbook topic", kind)
	}
}
