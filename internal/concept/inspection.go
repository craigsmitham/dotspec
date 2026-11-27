package concept

import (
	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/speccollection"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

func CheckType() check.Type {
	return speccollection.CheckType(Layout, readFrom, func(spec *Concept) *specfile.ReferenceDocument { return spec.References }, Check)
}
