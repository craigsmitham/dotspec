package specdoc

import (
	"testing"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

func TestEveryReportedKindHasAReportType(t *testing.T) {
	for kind := range families {
		if specfile.ReportType(kind) == "" {
			t.Errorf("the model declares no report class documenting %s", kind)
		}
	}
}
