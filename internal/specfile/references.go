package specfile

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

// ValidateReferences checks only same-document fragments. File and remote links
// remain navigation; this function does not infer cross-file identity.
func ValidateReferences(source, body []byte, doc ast.Node, anchors map[string]bool) error {
	return ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}
		destination := string(link.Destination)
		if !strings.HasPrefix(destination, "#") || destination == "#" {
			return ast.WalkContinue, nil
		}
		location := BodyLocation(source, body, link.Pos())
		location.Element = destination
		fragment, err := url.PathUnescape(destination[1:])
		if err != nil {
			return ast.WalkStop, diagnostic.Error("reference.invalid", fmt.Sprintf("invalid local reference %q", destination), location)
		}
		if !anchors[fragment] {
			return ast.WalkStop, diagnostic.Error("reference.unresolved", fmt.Sprintf("unresolved local reference %q; use a declared ID or a Markdown heading anchor", destination), location)
		}
		return ast.WalkContinue, nil
	})
}
