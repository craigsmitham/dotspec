package cli

import (
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/guide"
	"github.com/dotspec-md/dotspec/internal/speckinds"
	"github.com/spf13/cobra"
)

// Companions are authored in the frontmatter of every specification type.
var crossCuttingTopics = map[string]bool{"companion": true}

func TestHandbookTopicsMatchAuthoredDocumentTypes(t *testing.T) {
	missing := map[string]bool{}
	for _, typ := range speckinds.All() {
		missing[typ.Layout.Kind] = true
	}
	for _, topic := range guide.Topics() {
		if !missing[topic.Name] && !crossCuttingTopics[topic.Name] {
			t.Errorf("handbook topic %q has no authored document type", topic.Name)
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
		t.Errorf("authored document type %q has no handbook topic", kind)
	}
}

// Each Learn more line in command help opens a handbook topic or page offline.
func TestLearnMoreAddressesResolve(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		for _, line := range strings.Split(cmd.Long, "\n") {
			if address, ok := strings.CutPrefix(line, "Learn more: dotspec docs "); ok {
				if _, err := guide.Read(address); err != nil {
					t.Errorf("%s: %v", cmd.CommandPath(), err)
				}
			}
		}
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
}
