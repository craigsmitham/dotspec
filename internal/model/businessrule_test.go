package model

import "testing"

func TestBusinessRuleExportBindingsAndContextualReferences(t *testing.T) {
	dir := t.TempDir()
	body := "---\r\ntype: business-rule-spec\r\nid: approval\r\n---\r\n# Approval\r\n## Statement\r\nManager approval is required.\r\n## Applies to\r\nAll purchases.\r\n### Conditions\r\nAbove the limit.\r\n## Rationale\r\nAccountability.\r\n## Authority\r\nBoard policy.\r\n## Exceptions\r\n\r\n## Related specifications\r\n[Purchase](../use-cases/purchase.use-case.spec.md).\r\n### Local detail\r\nPreserved extension.\r\n"
	conformanceWrite(t, dir, "specs/business-rules/approval.business-rule.spec.md", "\ufeff"+body)
	conformanceWrite(t, dir, "specs/use-cases/purchase.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Purchase\n## Related specifications\n[Approval](../business-rules/approval.business-rule.spec.md#approval).\n")
	result := conformanceExport(t, dir)
	if result.ExitCode != 0 {
		t.Fatalf("export status %d", result.ExitCode)
	}
	checkSourceText(t, result.Snapshot)
	var rule Record
	for _, e := range result.Snapshot.Elements {
		if String(e, "kind") == "BusinessRule" {
			rule = e
		}
	}
	if rule == nil || String(rule, "authored_id") != "approval" {
		t.Fatalf("rule identity: %+v", rule)
	}
	if len(records(Object(rule["exceptions"])["spans"])) != 0 {
		t.Fatalf("empty exceptions lost: %+v", rule)
	}
	for _, name := range []string{"statement", "applicability", "rationale", "authority_and_evidence", "exceptions", "related_specifications"} {
		if _, ok := rule[name]; !ok {
			t.Errorf("missing named field %s", name)
		}
	}
	found := false
	for _, ref := range result.Snapshot.References {
		r := Object(ref["resolution"])
		if String(r, "target") == String(rule, "id") && String(r, "status") == "resolved" {
			found = true
		}
	}
	if !found {
		t.Fatal("authored rule ID did not resolve")
	}
}
