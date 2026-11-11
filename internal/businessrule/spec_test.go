package businessrule

import (
	"errors"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestBusinessRulePreservesBusinessMeaningAndLocations(t *testing.T) {
	source := "---\nid: approval\ntype: business-rule\n---\n# Purchase approval\n## Statement\nPurchases above $10,000 require manager approval.\n## Applies to\nAll departments.\n## Rationale\nControl expenditure.\n## Authority\nBoard policy.\n## Exceptions\nEmergency purchases follow [emergency policy](emergency.business-rule.spec.md).\n## Related specifications\n[Submit a purchase](../use-cases/submit.use-case.spec.md).\n"
	rule, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if rule.ID != "approval" || rule.Statement != "Purchases above $10,000 require manager approval." || rule.Authority != "Board policy." || !strings.Contains(rule.Exceptions, "emergency.business-rule.spec.md") || !strings.Contains(rule.Markdown, "# Purchase approval") {
		t.Fatalf("lost authored meaning: %+v", rule)
	}
	if rule.SectionLocations["Statement"].Line != 6 {
		t.Fatalf("statement location: %+v", rule.SectionLocations)
	}
	if len(Check(rule)) != 0 {
		t.Fatalf("complete rule advice: %+v", Check(rule))
	}
}

func TestBusinessRuleDraftAndStructuralErrors(t *testing.T) {
	prefix := "---\ntype: business-rule\n---\n# Rule\n"
	draft, err := Parse([]byte(prefix))
	if err != nil || len(Check(draft)) != 4 {
		t.Fatalf("draft: %+v %v", draft, err)
	}
	cases := []struct{ source, code string }{
		{prefix + "## Statement\nFirst\n## Statement\nSecond\n", "structure.duplicate-section"},
		{prefix + "### Statement\nRule\n", "structure.section-level"},
		{prefix + "## Authority\nFirst\n## Authority and evidence\nSecond\n", "structure.duplicate-section"},
	}
	for _, c := range cases {
		_, err := Parse([]byte(c.source))
		if err == nil {
			t.Fatalf("accepted %s", c.source)
		}
		var finding *diagnostic.Diagnostic
		if !errors.As(err, &finding) || finding.Rule != c.code {
			t.Fatalf("error: %v", err)
		}
	}
	for _, source := range []string{strings.Replace(prefix, "type: business-rule", "type: business-rule\nunknown: value", 1), strings.Replace(prefix, "type: business-rule", "type: job", 1), prefix + "[broken](#absent)\n"} {
		if _, err := Parse([]byte(source)); err == nil {
			t.Fatalf("accepted invalid source %s", source)
		}
	}
	rule, err := Parse([]byte(prefix + "## Statement\n\n## Authority and evidence\nPolicy.\n## Unrecognized\nRetain this.\n"))
	if err != nil || rule.Statement != "" || rule.Authority != "Policy." || !strings.Contains(rule.Markdown, "Retain this.") {
		t.Fatalf("empty/unknown content: %+v %v", rule, err)
	}
}

func FuzzParse(f *testing.F) {
	f.Add([]byte(Template))
	f.Add([]byte("---\ntype: business-rule\n---\n# Draft\n"))
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = Parse(data) })
}
