package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/testproject"
)

func structureFixture(t *testing.T) string {
	t.Helper()
	return testproject.Open(t, "system-structure")
}

func TestStructureAndFeatureMembership(t *testing.T) {
	root := structureFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
	store := designRecord(t, s, "System", "Online Store")
	inventory := designRecord(t, s, "Component", "Inventory")
	reservations := designRecord(t, s, "Component", "Reservations")
	if String(store, "authored_id") != "store" || reservations["specification_system"] != store["id"] || reservations["part_of"] != inventory["id"] {
		t.Fatal(store, reservations)
	}
	req := designRecord(t, s, "Requirement", "Reserve purchased stock")
	if targets := Strings(req["applies_to"]); len(targets) != 1 || targets[0] != reservations["id"] {
		t.Fatal(req)
	}
	for _, name := range []string{"Place an order", "Prevent overselling"} {
		f := designRecord(t, s, "Feature", name)
		if f["specification_owner"] != store["id"] || len(Strings(f["groups"])) != 1 || Strings(f["groups"])[0] != req["id"] {
			t.Fatal(f)
		}
	}
	if req["owner"] != nil {
		t.Fatal("feature membership invented declaration ownership")
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "Feature", "Place an order")["groups"] = []string{String(inventory, "id")}
	requireFinding(t, forged, "shape.invalid")
}

func TestStructuralErrorsPreserveSourceWithoutInventingParents(t *testing.T) {
	for _, tc := range []struct {
		name, path, old, replacement string
		error                        bool
	}{
		{"cycle", "components/inventory/specs/inventory.component.spec.md", "[Online Store](../../../specs/online-store.system.spec.md#store)", "[Reservations](../../reservations/specs/reservations.component.spec.md)", true},
		{"two parents", "components/reservations/specs/reservations.component.spec.md", "[Inventory](../../inventory/specs/inventory.component.spec.md)", "[Inventory](../../inventory/specs/inventory.component.spec.md) and [Online Store](../../../specs/online-store.system.spec.md)", true},
		{"wrong parent kind", "components/reservations/specs/reservations.component.spec.md", "[Inventory](../../inventory/specs/inventory.component.spec.md)", "[Feature](../../../specs/features/place-an-order.feature.spec.md)", true},
		{"wrong group kind", "specs/features/place-an-order.feature.spec.md", "[Reserve purchased stock](../requirements/reserve-stock.requirement.spec.md)", "[Inventory](../../components/inventory/specs/inventory.component.spec.md)", true},
		{"draft parent", "components/reservations/specs/reservations.component.spec.md", "[Inventory](../../inventory/specs/inventory.component.spec.md)", "To be decided.", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := structureFixture(t)
			path := filepath.Join(root, tc.path)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), tc.old) {
				t.Fatal("fixture missing", tc.old)
			}
			data = []byte(strings.ReplaceAll(string(data), tc.old, tc.replacement))
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.error {
				t.Fatal(result.ExitCode, result.Snapshot.Diagnostics)
			}
			checkSourceText(t, result.Snapshot)
		})
	}
}

func TestComponentMapsSpecificationSections(t *testing.T) {
	root := structureFixture(t)
	traceEdit(t, root, "components/checkout/specs/checkout.component.spec.md", "## Implementation", "## Evidence\n\nA component note.\n\n## Evidence\n\nEvidence to review.\n\n## Related specifications\n\n[Reservation](../../../specs/concepts/reservation.concept.spec.md)\n\n## Implementation")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	checkout := designRecord(t, result.Snapshot, "Component", "Checkout")
	if checkout["evidence"] == nil || checkout["related_specifications"] == nil {
		t.Fatal("component specification sections missing", checkout)
	}
	checkSourceText(t, result.Snapshot)
}

func TestComponentRequirementsAndCrossDocumentStructure(t *testing.T) {
	root := structureFixture(t)
	conformanceWrite(t, root, "components/checkout/specs/requirements/submission.requirement.spec.md", `---
type: requirement-spec
id: idempotent-submission
---
# Idempotent submission

## Statement

A repeated submission must not duplicate the order.

## Applies to

**Targets:** [Checkout](../checkout.component.spec.md)

## Rationale

Repeat a submission and count orders.
`)
	conformanceWrite(t, root, "worker/specs/worker.system.spec.md", `---
type: system-spec
id: worker
---
# Worker
## Purpose
Do background work.
## Scope
Fulfilment work.
## Part of
[Store](../../specs/online-store.system.spec.md#store)
`)
	conformanceWrite(t, root, "worker/queue/specs/worker-queue.component.spec.md", `---
type: component-spec
id: worker-queue
---
# Worker queue
## Part of
[Checkout](../../../components/checkout/specs/checkout.component.spec.md)
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	req := designRecord(t, s, "Requirement", "Idempotent submission")
	checkout := designRecord(t, s, "Component", "Checkout")
	if Strings(req["applies_to"])[0] != checkout["id"] {
		t.Fatal(req)
	}
	worker := designRecord(t, s, "System", "Worker")
	store := designRecord(t, s, "System", "Online Store")
	queue := designRecord(t, s, "Component", "Worker queue")
	if worker["part_of"] != store["id"] || queue["part_of"] != checkout["id"] || queue["specification_system"] != store["id"] {
		t.Fatal(worker, queue)
	}
	// A parent component is allowed to be declared in another system document.
	checkSourceText(t, s)
}
