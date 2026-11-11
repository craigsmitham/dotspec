package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func structureFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/system-structure")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestStructureAndFeatureMembership(t *testing.T) {
	root := structureFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	store := designRecord(t, s, "System", "Online Store")
	inventory := designRecord(t, s, "Component", "Inventory")
	reservations := designRecord(t, s, "Component", "Reservations")
	if String(store, "authored_id") != "store" || reservations["owner"] != store["id"] || reservations["part_of"] != inventory["id"] {
		t.Fatal(store, reservations)
	}
	req := designRecord(t, s, "GeneralRequirement", "Reserve purchased stock")
	if targets := Strings(req["applies_to"]); len(targets) != 1 || targets[0] != reservations["id"] {
		t.Fatal(req)
	}
	for _, name := range []string{"Place an order", "Prevent overselling"} {
		f := designRecord(t, s, "Feature", name)
		if f["system"] != store["id"] || len(Strings(f["groups"])) != 1 || Strings(f["groups"])[0] != req["id"] {
			t.Fatal(f)
		}
	}
	if req["owner"] != nil {
		t.Fatal("feature membership invented declaration ownership")
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "Feature", "Place an order")["groups"] = []string{String(inventory, "id")}
	requireFinding(t, forged, "reference.kind")
}

func TestStructuralErrorsPreserveSourceWithoutInventingParents(t *testing.T) {
	for _, tc := range []struct {
		name, old, replacement string
		error                  bool
	}{
		{"cycle", "[Online Store](#store)", "[Reservations](#reservations)", true},
		{"two parents", "[Inventory](#inventory)", "[Inventory](#inventory) and [Online Store](#store)", true},
		{"wrong parent kind", "[Inventory](#inventory)", "[Feature](#place-order)", true},
		{"wrong group kind", "[Reserve purchased stock](requirements/reserve-stock.requirement.spec.md)", "[Inventory](#inventory)", true},
		{"draft parent", "[Inventory](#inventory)", "To be decided.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := structureFixture(t)
			path := filepath.Join(root, "specs/system.spec.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			// Replace all system parents to make the cycle cross two component levels.
			data = []byte(strings.ReplaceAll(string(data), tc.old, tc.replacement))
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.error {
				t.Fatal(result.ExitCode, result.Snapshot.Diagnostics)
			}
			conformancePayloads(t, result.Snapshot)
			if tc.name == "draft parent" {
				c := designRecord(t, result.Snapshot, "Component", "Reservations")
				if c["part_of"] != nil {
					t.Fatal("inferred parent from declaration", c)
				}
			}
		})
	}
}

func TestInlineSpecificationsMapCommonContent(t *testing.T) {
	root := structureFixture(t)
	traceEdit(t, root, "specs/system.spec.md", "### Inventory", "#### Notes\n\nA component note.\n\n#### Evidence\n\nEvidence to review.\n\n#### Related specifications\n\n[Reservation](concepts/reservation.concept.spec.md)\n\n### Inventory")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	checkout := designRecord(t, result.Snapshot, "Component", "Checkout")
	if checkout["extension_sections"] != nil || checkout["notes"] == nil || checkout["evidence"] == nil || checkout["related_specifications"] == nil {
		t.Fatal("common specification detail not mapped", checkout)
	}
	conformancePayloads(t, result.Snapshot)
}
