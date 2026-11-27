package job

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func CheckType() check.Type { return CheckTypeWithPreset("general") }

func CheckTypeWithPreset(preset string) check.Type {
	return speccollection.CheckType(Layout, readFrom, func(job *Job) *specfile.ReferenceDocument { return job.References },
		func(job *Job) []diagnostic.Diagnostic { return CheckWithPreset(job, preset) })
}
