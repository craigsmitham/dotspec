package check_test

import (
	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/concept"
	"github.com/dotspec-md/dotspec/internal/job"
	"github.com/dotspec-md/dotspec/internal/requirement"
	"github.com/dotspec-md/dotspec/internal/system"
	"github.com/dotspec-md/dotspec/internal/usecase"
)

// Exercise the same built-in composition as the CLI against the existing
// behavior contracts, independently of tests for other type integrations.
func checkProject(path string) check.Report {
	return check.Project(path, system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType())
}

func checkJob(path, slug string) check.Report {
	return check.Focused(path, slug, job.CheckType())
}
