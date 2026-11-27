package usecase

import (
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

// Check suggests missing content without judging whether the narrative is
// complete, the goal is achieved, or a linked rule is adequately expressed.
func Check(spec *UseCase) []diagnostic.Diagnostic {
	blank := specfile.NewBlankSections(spec.Path, spec.Location, spec.SectionLocations)
	// Success guarantees, extensions and sequence describe a scenario, so they
	// wait until one is written.
	hasScenario := strings.TrimSpace(spec.MainSuccessScenario) != ""
	blank.Suggest("use-case.goal", "Goal", spec.Goal, "Describe what the primary actor wants to accomplish through this interaction.")
	blank.Suggest("use-case.system-boundary", "System boundary", spec.BoundaryDescription, "Name the system whose behavior this use case describes; file ownership does not define this boundary.")
	blank.Suggest("use-case.primary-actor", "Primary actor", spec.PrimaryActor, "Identify the role pursuing the goal.")
	blank.Suggest("use-case.goal-level", "Goal level", spec.GoalLevel, "Identify the goal level, such as summary, user task, or subfunction.")
	if hasScenario {
		blank.Suggest("use-case.success-guarantees", "Success guarantees", spec.SuccessGuarantees, "State or link to what holds when the actor's goal is achieved.")
	}
	blank.Suggest("use-case.main-success-scenario", "Main success scenario", spec.MainSuccessScenario, "Describe or reference the actor-system interaction that reaches the goal.")
	if hasScenario {
		blank.Suggest("use-case.extensions", "Extensions", spec.Extensions, "Consider alternative success paths and failure conditions. List conditions first; handling may remain explicitly unresolved.")
		blank.Suggest("use-case.sequence", "Sequence", spec.Sequence, "State whether scenario order is representative or required, and link to the governing rule for required order.")
	}
	return blank.Advice
}
