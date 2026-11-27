package specdoc

import (
	"testing"

	"github.com/dotspec-md/dotspec/internal/specfile"
)

func TestEveryReportedKindHasAReportType(t *testing.T) {
	for kind := range families {
		if specfile.ReportType(kind) == "" {
			t.Errorf("the model declares no report class documenting %s", kind)
		}
	}
}
