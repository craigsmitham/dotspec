package specfile

import (
	"fmt"
	"regexp"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// Identifiers tracks explicit identities within one document. The caller owns
// which elements define IDs and when to expose them as reference targets.
// The zero value is ready to use; absent optional IDs are ignored.
type Identifiers struct {
	definitions map[string]diagnostic.Location
}

func (ids *Identifiers) Register(id, element string, location diagnostic.Location) error {
	if id == "" {
		return nil
	}
	location.Element = element
	if !identifierPattern.MatchString(id) {
		return diagnostic.Error("identity.invalid", fmt.Sprintf("%s: invalid ID %q; use letters, digits, hyphens, or underscores", element, id), location)
	}
	if previous, exists := ids.definitions[id]; exists {
		return diagnostic.Error("identity.duplicate", fmt.Sprintf("duplicate ID %q in %s and %s; assign distinct IDs", id, previous.Element, element), location, previous)
	}
	if ids.definitions == nil {
		ids.definitions = make(map[string]diagnostic.Location)
	}
	ids.definitions[id] = location
	return nil
}

func (ids *Identifiers) AddAnchors(anchors map[string]bool) {
	for id := range ids.definitions {
		anchors[id] = true
	}
}
