// Package speckinds lists the built-in specification kinds. It is a leaf that
// imports only the check integration and the kind packages.
package speckinds

import (
	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

// All returns the check integration of every built-in specification kind.
// Independently registered types are composed by their callers.
func All() []check.Type {
	return []check.Type{qualitymodel.CheckType(), businessrule.CheckType(), system.CheckType(), system.ComponentCheckType(), system.ProductCheckType(system.FeatureLayout), system.ProductCheckType(system.CapabilityLayout), system.AttributeCheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType()}
}
