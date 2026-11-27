package system

import (
	"slices"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func adviceRules(findings []diagnostic.Diagnostic) []string {
	rules := []string{}
	for _, f := range findings {
		rules = append(rules, f.Rule+"@"+f.Location.Element)
	}
	return rules
}

func TestAssumptionAndInterfaceAdvice(t *testing.T) {
	const head = "---\ntype: system-spec\n---\n# Example\n## Purpose\nHelp.\n"
	for _, tt := range []struct {
		name, body string
		want       []string
	}{
		{"supported assumption", "## Assumptions\n### Same storage\nThe file persists.\n#### Basis\nThe owner names it.\n#### Applicability\nEach restart.\n", []string{}},
		{"bare assumption", "## Assumptions\n### Same storage\nThe file persists.\n", []string{"assumption.basis@Assumptions / Same storage", "assumption.applicability@Assumptions / Same storage"}},
		{"blank basis", "## Assumptions\n### Same storage\nThe file persists.\n#### Basis\n#### Applicability\nEach restart.\n", []string{"assumption.basis@Assumptions / Same storage / Basis"}},
		{"linked contract", "## Provided interfaces\n### Forms\nThe forms.\n#### Contracts\n[Template](list.html).\n", []string{}},
		{"contract prose without a link", "## Required interfaces\n### Clock\nThe time source.\n#### Contracts\nAgreed informally.\n", []string{"interface.contract@Required interfaces / Clock"}},
		{"no contract", "## Provided interfaces\n### Forms\nThe forms.\n", []string{"interface.contract@Provided interfaces / Forms"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := Parse([]byte(head + tt.body))
			if err != nil {
				t.Fatal(err)
			}
			spec.Path = "example.system.spec.md"
			findings := Check(spec)
			if got := adviceRules(findings); !slices.Equal(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for _, f := range findings {
				if f.Severity != "advice" || f.Location.Path != spec.Path || f.Location.Line == 0 {
					t.Fatalf("finding is not located advice: %+v", f)
				}
			}
		})
	}
}

func TestComponentInterfaceAdvice(t *testing.T) {
	source := "---\ntype: component-spec\n---\n# Service\n## Part of\n[Example](example.system.spec.md)\n## Required interfaces\n### Clock\nThe time source.\n"
	c, err := ParseComponent([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if got := adviceRules(c.advice); !slices.Equal(got, []string{"interface.contract@Required interfaces / Clock"}) {
		t.Fatal(got)
	}
}
