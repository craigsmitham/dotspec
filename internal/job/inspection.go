package job

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func CheckType() check.Type {
	return speccollection.CheckType(Layout, readFrom, func(job *Job) *specfile.ReferenceDocument { return job.References }, Check)
}
