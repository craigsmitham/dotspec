package check_test

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

// Exercise the same built-in composition as the CLI against the existing
// behavior contracts, independently of tests for other type integrations.
func checkProject(path string) check.Report {
	return check.Project(path, system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), quality.CheckType(), requirement.CheckType())
}

func checkJob(path, slug string) check.Report {
	return check.Focused(path, slug, job.CheckType())
}

func checkUseCase(path, slug string) check.Report {
	return check.Focused(path, slug, usecase.CheckType())
}

func checkConcept(path, slug string) check.Report {
	return check.Focused(path, slug, concept.CheckType())
}
