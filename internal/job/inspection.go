package job

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/system"
)

func CheckType() check.Type {
	return check.Type{Layout: Layout, ValidateSlug: ValidateSlug, Context: system.CheckContext,
		Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
			spec, err := Read(path)
			if err != nil {
				return nil, err
			}
			return Check(spec), nil
		}}
}
