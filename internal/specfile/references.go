package specfile

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

// ValidateReferences checks that each same-document fragment link resolves to
// one of the outline's anchors. Check commands resolve relative file links
// separately, with repository context.
func (o *Outline) ValidateReferences() error {
	return ast.Walk(o.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
		location := o.Location(link.Pos(), destination)
		fragment, err := url.PathUnescape(destination[1:])
		if err != nil {
			return ast.WalkStop, diagnostic.Error("reference.invalid", fmt.Sprintf("invalid local reference %q", destination), location)
		}
		if !o.Anchors[fragment] {
			return ast.WalkStop, diagnostic.Error("reference.unresolved", fmt.Sprintf("unresolved local reference %q; use a declared ID or a Markdown heading anchor", destination), location)
		}
		return ast.WalkContinue, nil
	})
}
