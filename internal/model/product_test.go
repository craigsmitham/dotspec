package model

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestProductSourceOwnershipAndOfflineMembership(t *testing.T) {
	root := t.TempDir()
	put := func(path, body string) { t.Helper(); conformanceWrite(t, root, path, body) }
	put("dotspec.yaml", "name: example\n")
	put("specs/shop.system.spec.md", "---\ntype: system-spec\n---\n# Shop\n## Purpose\nPlace orders.\n## Scope\nOrders.\n")
	put("billing/specs/billing.system.spec.md", "---\ntype: system-spec\n---\n# Billing\n## Purpose\nBill.\n## Scope\nInvoices.\n## Part of\n[Shop](../../specs/shop.system.spec.md)\n")
	put("other/specs/other.system.spec.md", "---\ntype: system-spec\n---\n# Other\n## Purpose\nOther.\n## Scope\nOther.\n")
	put("specs/requirements/receipt.requirement.spec.md", "---\ntype: requirement-spec\n---\n# Receipt\n## Statement\nSend a receipt.\n## Adopts\n### Receipt fields\n[Fields](../../partner/contract.md#fields).\n#### Edition\n2026.\n#### Applicability\nPaid orders.\n")
	put("partner/dotspec.yaml", "name: partner\n")
	put("partner/contract.md", "# Contract\n\n## Fields\nOrder ID.\n")
	const interaction = "---\ntype: use-case-spec\n---\n# Pay\n## Goal\nPlace an order.\n## Related specifications\n[Receipt](../requirements/receipt.requirement.spec.md).\n## Extensions\n### Payment declined\nTry again.\n"
	put("specs/use-cases/pay.use-case.spec.md", interaction)
	put("specs/use-cases/retry.use-case.spec.md", strings.Replace(interaction, "# Pay\n", "# Retry\n", 1))
	put("specs/quality-models/team.quality-model.md", "---\ntype: quality-model\n---\n# Team\n## Quality attributes\n### Timeliness\n**ID:** timeliness\nRespond in time.\n")
	get := func() (*Snapshot, *ProductSpecification) {
		t.Helper()
		result := conformanceExport(t, root)
		if result.ExitCode != 0 {
			t.Fatal(result.Snapshot.Diagnostics)
		}
		view, err := result.Snapshot.Product("shop")
		if err != nil {
			t.Fatal(err)
		}
		return result.Snapshot, view
	}
	s, view := get()
	if !view.Complete || len(view.Members) != 5 || len(view.Adoptions) != 1 || len(view.ReferenceCatalog) != 1 {
		t.Fatal(view)
	}
	if _, err := s.Product(""); err == nil {
		t.Fatal("independent systems require a selection")
	}
	for _, m := range view.Members {
		if strings.Contains(m.Source, "other.system") || m.Kind == "QualityModel" {
			t.Fatal(m)
		}
	}
	citing := map[string]bool{}
	for _, ref := range view.References {
		if ref["authored_target"] == "../requirements/receipt.requirement.spec.md" {
			citing[String(ref, "source_document")] = true
		}
	}
	if len(citing) != 2 || len(view.ReferencedSources) != 1 || !strings.HasSuffix(String(view.ReferencedSources[0], "locator"), "partner/contract.md") {
		t.Fatal(citing, view.ReferencedSources)
	}
	offline, err := cloneSnapshot(t, s).Product("shop")
	if err != nil || !reflect.DeepEqual(view.Members, offline.Members) {
		t.Fatal(err, offline)
	}
	a, _ := json.Marshal(view.Adoptions)
	b, _ := json.Marshal(offline.Adoptions)
	if string(a) != string(b) {
		t.Fatalf("%s != %s", a, b)
	}
	put("specs/use-cases/pay.use-case.spec.md", strings.Replace(interaction, "Try again.", "Ask for another card.", 1))
	put("partner/contract.md", "# Contract\n\n## Fields\nOrder ID and currency.\n")
	put("specs/quality-models/extra.quality-model.md", "---\ntype: quality-model\n---\n# Extra\n## Quality attributes\n### Cost\n**ID:** cost\nPrice per order.\n")
	_, revised := get()
	if len(revised.ReferenceCatalog) != 2 {
		t.Fatal(revised.ReferenceCatalog)
	}
	if digest := view.Adoptions[0]["observed_source_sha256"]; digest == nil || digest == revised.Adoptions[0]["observed_source_sha256"] {
		t.Fatal("a revised provision must be visible for review")
	}
	if !reflect.DeepEqual(view.Members, revised.Members) {
		t.Fatal("an extension edit changed membership")
	}
	broken := cloneSnapshot(t, s)
	adoption := elementKind(t, broken, "AdoptedProvision")
	for _, ref := range broken.References {
		if ref["authored_target"] == "../requirements/receipt.requirement.spec.md" {
			adoption["provision"] = ref["id"]
			break
		}
	}
	requireFinding(t, broken, "adoption.provision")
	for _, doc := range broken.Documents {
		if doc["source"] == adoption["source_document"] {
			doc["elements"] = []string{String(adoption, "owner")}
		}
	}
	requireFinding(t, broken, "document.elements")
	put("specs/requirements/broken.requirement.spec.md", "invalid")
	incomplete := conformanceExport(t, root)
	v, err := incomplete.Snapshot.Product("shop")
	if err != nil || v.Complete || len(v.Diagnostics) == 0 || len(v.Members) != 5 {
		t.Fatal(err, v)
	}
	scoped := conformanceExport(t, filepath.Join(root, "specs", "requirements"))
	if partial, err := scoped.Snapshot.Product("shop"); err != nil || partial.Complete {
		t.Fatal(err, partial)
	}
}
