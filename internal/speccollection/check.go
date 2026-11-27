package speccollection

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

// CheckType inspects family members with their associated system as context.
// read validates placement and structure; advise runs only after it succeeds.
func CheckType[T any](layout specfile.Layout, read func(*specfile.SourceSet, string) (*T, error), references func(*T) *specfile.ReferenceDocument, advise func(*T) []diagnostic.Diagnostic) check.Type {
	kind := check.ReadType(layout, read, references, advise)
	kind.Context = system.CheckContext
	return kind
}
