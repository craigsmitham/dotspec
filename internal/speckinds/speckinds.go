// Package speckinds lists the built-in specification kinds. It is a leaf that
// imports only the check integration and the kind packages.
package speckinds

import (
	"github.com/dotspec-md/dotspec/internal/businessrule"
	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/concept"
	"github.com/dotspec-md/dotspec/internal/job"
	"github.com/dotspec-md/dotspec/internal/qualitymodel"
	"github.com/dotspec-md/dotspec/internal/requirement"
	"github.com/dotspec-md/dotspec/internal/system"
	"github.com/dotspec-md/dotspec/internal/usecase"
)

// All returns the check integration of every built-in specification kind.
// Independently registered types are composed by their callers.
func All() []check.Type {
	return []check.Type{qualitymodel.CheckType(), businessrule.CheckType(), system.CheckType(), system.ComponentCheckType(), system.ProductCheckType(system.FeatureLayout), system.ProductCheckType(system.CapabilityLayout), system.AttributeCheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType()}
}
