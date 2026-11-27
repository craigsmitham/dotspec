package cli

import (
	"os"
	"strings"
	"testing"
)

func TestUseCaseBriefAndOverview(t *testing.T) {
	prepareExample(t, "overview")
	if code, out, errOut := invoke("use-case", "init", "sketch"); code != 0 || errOut != "" {
		t.Fatalf("init: %d %s %s", code, out, errOut)
	}
	data, err := os.ReadFile("specs/use-cases/sketch.use-case.spec.md")
	if err != nil || strings.Contains(string(data), "## Preconditions") || !strings.Contains(string(data), "## System boundary") {
		t.Fatalf("brief: %s %v", data, err)
	}
	writeFixture(t, "specs/use-cases/sketch.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Reserve equipment\n## Primary actor\nCustomer\n## Goal level\nUser task\n## System boundary\nRental business\nincluding depot handover.\n")
	code, out, errOut := invoke("use-case", "list", "--overview")
	for _, want := range []string{"System: New system", "Customer", "Reserve equipment", "User task", "Rental business including depot handover.", "sketch.use-case.spec.md"} {
		if code != 0 || errOut != "" || !strings.Contains(out, want) {
			t.Fatalf("overview missing %q: %d %s %s", want, code, out, errOut)
		}
	}
	code, ordinary, errOut := invoke("use-case", "list", "--json")
	overviewCode, overview, overviewErr := invoke("use-case", "list", "--overview", "--json")
	if code != 0 || overviewCode != 0 || errOut != "" || overviewErr != "" || ordinary != overview {
		t.Fatal("overview changed JSON records")
	}
	writeFixture(t, "specs/use-cases/sketch.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Reserve equipment\n")
	_, out, _ = invoke("use-case", "list", "--overview")
	if strings.Count(out, "(unspecified)") != 3 {
		t.Fatalf("missing fields were inferred: %s", out)
	}
}
