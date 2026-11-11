package model

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestOverviewIsPortableAndSeparatesRelationships(t *testing.T) {
	render := func(root string) string {
		t.Helper()
		result := conformanceExport(t, root)
		if result.ExitCode != 0 {
			t.Fatal(result.Snapshot.Diagnostics)
		}
		before, _ := json.Marshal(result.Snapshot)
		text, err := OverviewMarkdown(result.Snapshot, filepath.Join(root, "docs"))
		if err != nil {
			t.Fatal(err)
		}
		after, _ := json.Marshal(result.Snapshot)
		if string(before) != string(after) {
			t.Fatal("renderer mutated snapshot")
		}
		return text
	}
	text := render(structureFixture(t))
	if text != render(structureFixture(t)) {
		t.Fatal("overview depends on checkout identity")
	}
	for _, want := range []string{"flowchart BT", "|part of|", "Components involved (derived)", "Reserve purchased stock", "Stock holds", "Submit concurrent orders", "../order-contract.md", "explicit", "No requirements are inherited"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q", want)
		}
	}
	diagram := strings.Split(strings.Split(text, "```mermaid\n")[1], "```")[0]
	if strings.Contains(diagram, "Feature:") || strings.Contains(diagram, "Place an order") {
		t.Fatal("feature became a structural part")
	}
	root := structureFixture(t)
	traceEdit(t, root, "specs/system.spec.md", "[Inventory](#inventory)", "[Missing](missing.md)")
	result := conformanceExport(t, root)
	text, err := OverviewMarkdown(result.Snapshot, root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Unresolved or undeclared") || !strings.Contains(text, "[Missing](specs/missing.md)") || !strings.Contains(text, "Inspection findings") {
		t.Fatal("unresolved authored relationship disappeared", text)
	}
}

func TestComponentClausesAndCrossDocumentStructure(t *testing.T) {
	root := structureFixture(t)
	traceEdit(t, root, "specs/system.spec.md", "### Inventory", `#### Obligations

##### Idempotent submission

**ID:** idempotent-submission

A repeated submission must not duplicate the order.

###### Applies to

**Targets:** [Checkout](#checkout)

###### Verification plan

Repeat a submission and count orders.

### Inventory`)
	conformanceWrite(t, root, "worker/specs/system.spec.md", `---
type: system
id: worker
---
# Worker
## Purpose
Do background work.
## Scope
Fulfilment work.
## Part of
[Store](../../specs/system.spec.md#store)
## Components
### Worker queue
**ID:** worker-queue
#### Part of
[Checkout](../../specs/system.spec.md#checkout)
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	clause := designRecord(t, s, "ObligationClause", "Idempotent submission")
	checkout := designRecord(t, s, "Component", "Checkout")
	if clause["owner"] != checkout["id"] || Strings(clause["applies_to"])[0] != checkout["id"] {
		t.Fatal(clause)
	}
	worker := designRecord(t, s, "System", "Worker")
	store := designRecord(t, s, "System", "Online Store")
	queue := designRecord(t, s, "Component", "Worker queue")
	if worker["part_of"] != store["id"] || queue["part_of"] != checkout["id"] || queue["owner"] != worker["id"] {
		t.Fatal(worker, queue)
	}
	// A parent component is allowed to be declared in another system document.
	conformancePayloads(t, s)
}
