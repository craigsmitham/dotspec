package usecase

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// Check suggests missing content without judging whether the narrative is
// complete, the goal is achieved, or a linked rule is adequately expressed.
func Check(spec *UseCase) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	for _, field := range []struct{ rule, name, value, message string }{
		{"use-case.goal", "Goal", spec.Goal, "Describe what the primary actor wants to accomplish through this interaction."},
		{"use-case.system-boundary", "System boundary", spec.SystemBoundary, "Name the system whose behavior this use case describes; file ownership does not define this boundary."},
		{"use-case.primary-actor", "Primary actor", spec.PrimaryActor, "Identify the role pursuing the goal."},
		{"use-case.goal-level", "Goal level", spec.GoalLevel, "Identify the goal level, such as summary, user task, or subfunction."},
		{"use-case.success-criterion", "Success criterion", spec.SuccessCriterion, "State or link to the observable condition that completes the actor's goal."},
		{"use-case.main-success-scenario", "Main success scenario", spec.MainSuccessScenario, "Describe or reference the actor-system interaction that reaches the goal."},
		{"use-case.sequence", "Sequence", spec.Sequence, "State whether scenario order is representative or required, and link to the governing rule for required order."},
	} {
		if strings.TrimSpace(field.value) != "" || (field.name == "Sequence" && strings.TrimSpace(spec.MainSuccessScenario) == "") {
			continue
		}
		location, found := spec.SectionLocations[field.name]
		if !found {
			location = spec.Location
		}
		location.Path, location.Element = spec.Path, field.name
		suggestions = append(suggestions, diagnostic.Advice(field.rule, field.message, location))
	}
	return suggestions
}
