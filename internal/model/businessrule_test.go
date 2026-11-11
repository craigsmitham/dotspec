package model

import (
	"strings"
	"testing"
)

func TestBusinessRuleExportBindingsAndContextualReferences(t *testing.T) {
	dir := t.TempDir()
	body := "---\r\ntype: business-rule\r\nid: approval\r\n---\r\n# Approval\r\n## Statement\r\nManager approval is required.\r\n## Applies to\r\nAll purchases.\r\n### Conditions\r\nAbove the limit.\r\n## Rationale\r\nAccountability.\r\n## Authority\r\nBoard policy.\r\n## Exceptions\r\n\r\n## Related specifications\r\n[Purchase](../use-cases/purchase.use-case.spec.md).\r\n## Local custom\r\nPreserved extension.\r\n"
	conformanceWrite(t, dir, "specs/business-rules/approval.business-rule.spec.md", "\ufeff"+body)
	conformanceWrite(t, dir, "specs/use-cases/purchase.use-case.spec.md", "---\ntype: use-case\n---\n# Purchase\n## Business rules\n[Approval](../business-rules/approval.business-rule.spec.md#approval).\n")
	result := conformanceExport(t, dir)
	if result.ExitCode != 0 {
		t.Fatalf("export status %d", result.ExitCode)
	}
	conformancePayloads(t, result.Snapshot)
	var rule Record
	for _, e := range result.Snapshot.Elements {
		if String(e, "kind") == "BusinessRule" {
			rule = e
		}
	}
	if rule == nil || String(rule, "authored_id") != "approval" {
		t.Fatalf("rule identity: %+v", rule)
	}
	if String(Object(rule["exceptions"]), "presence") != "empty" {
		t.Fatalf("empty exceptions lost: %+v", rule)
	}
	for _, name := range []string{"statement", "applies_to", "rationale", "authority_and_evidence", "exceptions", "related_specifications"} {
		if _, ok := rule[name]; !ok {
			t.Errorf("missing named field %s", name)
		}
	}
	if len(Strings(rule["extension_sections"])) != 1 {
		t.Fatal("unknown section not represented")
	}
	found := false
	for _, r := range result.Snapshot.Resolutions {
		if String(r, "target") == String(rule, "id") && String(r, "status") == "resolved" {
			found = true
		}
	}
	if !found {
		t.Fatal("authored rule ID did not resolve")
	}
	for _, a := range result.Snapshot.Assertions {
		if strings.Contains(String(a, "kind"), "Enforce") {
			t.Fatal("citation inferred enforcement")
		}
	}
}
