package specfile

import (
	"errors"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestIdentifiersRetainDefinitionsWithinOneDocument(t *testing.T) {
	var ids Identifiers
	first := diagnostic.Location{Line: 3, Column: 5, Element: "criterion"}
	if err := ids.Register("QR-1", "criterion", first); err != nil {
		t.Fatal(err)
	}
	if err := ids.Register("", "unidentified draft", diagnostic.Location{}); err != nil {
		t.Fatal(err)
	}
	second := diagnostic.Location{Line: 12, Column: 1, Element: "measurement"}
	err := ids.Register("QR-1", "measurement", second)
	var finding *diagnostic.Diagnostic
	if !errors.As(err, &finding) || finding.Rule != "identity.duplicate" || finding.Location != second || len(finding.Related) != 1 || finding.Related[0] != first {
		t.Fatalf("lost duplicate definitions: %+v (%v)", finding, err)
	}
	anchors := map[string]bool{"criterion": true}
	ids.AddAnchors(anchors)
	if len(anchors) != 2 || !anchors["QR-1"] || !anchors["criterion"] {
		t.Fatalf("lost heading or explicit identity: %v", anchors)
	}
	var anotherDocument Identifiers
	if err := anotherDocument.Register("QR-1", "criterion", first); err != nil {
		t.Fatalf("identity leaked across documents: %v", err)
	}
}
