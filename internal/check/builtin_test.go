package check_test

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

// Exercise the same built-in composition as the CLI against the existing
// behavior contracts, independently of tests for other type integrations.
func checkProject(path string) check.Report {
	return check.Project(path, system.CheckType(), job.CheckType(), usecase.CheckType())
}

func checkJob(path, slug string) check.Report {
	return check.Focused(path, slug, job.CheckType())
}

func checkUseCase(path, slug string) check.Report {
	return check.Focused(path, slug, usecase.CheckType())
}
