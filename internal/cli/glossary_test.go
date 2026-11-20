package cli

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/model"
)

// Two systems share a definition. Only the root system also declares the
// imported term's alias as a preferred term.
func TestSharedGlossarySources(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: docs\ngenerated_spec: true\n")
	writeFixture(t, "specs/docs.system.spec.md", `---
type: system-spec
glossary_sources: [../words.md]
---
# Docs
## Purpose
Explain words.
## Scope
Examples. See [the local term](#draft).

[^note]: Local terminology note.
## Glossary
Terms used by the documentation team.
### Draft
An early version of a document.[^note]
`)
	writeFixture(t, "child/specs/child.system.spec.md", `---
type: system-spec
glossary_sources: [../../words.md]
---
# Child
## Purpose
Share definitions.
## Scope
Examples.
## Part of
[Docs](../../specs/docs.system.spec.md)
`)
	const vocab = "# Shared vocabulary\n\n[^note]: Vocabulary note.\n\n## Sketch\n\nA preliminary drawing.[^note]\n\n### Aliases\n\n- Draft\n"
	writeFixture(t, "words.md", vocab)
	run := func(code int, args ...string) string {
		t.Helper()
		got, out, err := invoke(args...)
		if got != code {
			t.Fatalf("%v: exit %d, expected %d\n%s%s", args, got, code, out, err)
		}
		return out + err
	}
	read := func(name string) string {
		t.Helper()
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	run(0, "spec", "build")
	run(0, "check")
	overview := read("specification/docs.system.spec-report.md")
	if !strings.Contains(overview, "Terms used by the documentation team.") || !strings.Contains(overview, "glossary.md#term.draft") || strings.Contains(overview, "An early version") {
		t.Fatal(overview)
	}
	index := read("specification/glossary.md")
	if !strings.Contains(index, "../words.md#sketch") || !strings.Contains(index, "../specs/docs.system.spec.md#draft") || strings.Index(index, "## Draft") >= strings.Index(index, "## Sketch") {
		t.Fatal(index)
	}
	for _, note := range []string{"document.[^source1-note]", "[^source1-note]: Local terminology note.", "drawing.[^source2-note]", "[^source2-note]: Vocabulary note."} {
		if !strings.Contains(index, note) {
			t.Fatalf("missing %q in %s", note, index)
		}
	}
	child := read("specification/systems/child/glossary.md")
	if !strings.Contains(child, "../../../words.md#sketch") || strings.Contains(child, "## Draft") {
		t.Fatal(child)
	}
	review := run(0, "model", "review")
	if strings.Count(review, "Review required for **draft**") != 1 || !strings.Contains(review, "words.md#sketch") || !strings.Contains(review, "docs.system.spec.md#draft") {
		t.Fatal(review)
	}
	var s model.Snapshot
	if err := json.Unmarshal([]byte(run(0, "model", "export")), &s); err != nil {
		t.Fatal(err)
	}
	var docs []string
	for _, r := range s.Elements {
		if model.String(r, "kind") == "System" {
			docs = append(docs, model.Strings(r["glossary_sources"])...)
		}
	}
	if len(docs) != 2 || docs[0] != docs[1] {
		t.Fatal(docs)
	}
	writeFixture(t, "words.md", strings.Replace(vocab, "preliminary", "rough", 1))
	run(1, "check")
	run(0, "spec", "build")
	if !strings.Contains(read("specification/glossary.md"), "A rough drawing.") || !strings.Contains(read("specification/systems/child/glossary.md"), "A rough drawing.") {
		t.Fatal("shared definition did not refresh")
	}
	if !strings.Contains(run(0, "spec", "build"), "0 written, 0 removed") {
		t.Fatal("rebuild changed files")
	}
	writeFixture(t, "words.md", "# Vocabulary\n\n### Unscoped term\n\nInvalid level.\n")
	if !strings.Contains(run(1, "check"), "glossary.term") {
		t.Fatal("invalid glossary lacks a diagnostic")
	}
	if err := os.Remove("words.md"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(run(1, "check"), "glossary.source") {
		t.Fatal("missing source lacks a diagnostic")
	}
}

func TestVocabularyContextsTopicsAndRelationships(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: ledger\ngenerated_spec: true\n")
	writeFixture(t, "specs/ledger.system.spec.md", `---
type: system-spec
glossary_sources: [../billing.md, ../collections.md, ../identity.md]
---
# Ledger
## Purpose
Track payments.
## Scope
Customer balances.
## Glossary
### Account
A customer record.
#### Topics
- Customers
`)
	writeFixture(t, "child/specs/child.system.spec.md", `---
type: system-spec
glossary_sources: [../../billing.md, ../../identity.md]
---
# Child
## Purpose
Share terminology.
## Scope
Billing and identity.
## Part of
[Ledger](../../specs/ledger.system.spec.md)
`)
	const billing = `---
context: Billing
---
# Billing vocabulary
## Account
A ledger of payments.
### Aliases
- Balance
### Topics
- Customers
- Payments
### Relationships
- Uses the [login][login].
### Example
This [account](identity.md#account) can have two invoices.
## Invoice
A request for payment.
### Related terms
[Account](identity.md#account)

[login]: identity.md#account
`
	writeFixture(t, "billing.md", billing)
	writeFixture(t, "collections.md", "---\ncontext: Billing\n---\n# Collections\n## Balance\nAn amount due.\n### Topics\n- Collections\n")
	writeFixture(t, "identity.md", "---\ncontext: Identity\n---\n# Identity\n## Account\nA login.\n### Aliases\n- Balance\n### Topics\n- Customers\n")
	run := func(code int, args ...string) string {
		t.Helper()
		got, stdout, stderr := invoke(args...)
		if got != code {
			t.Fatalf("%v exited %d\n%s%s", args, got, stdout, stderr)
		}
		return stdout + stderr
	}
	run(0, "spec", "build")
	run(0, "check")
	data, err := os.ReadFile("specification/glossary.md")
	if err != nil {
		t.Fatal(err)
	}
	page := string(data)
	for _, want := range []string{"## Account (Billing)", "## Account (Identity)", "## Account (Ledger)", "## Browse by topic", "- **Customers**", "- **Payments**", "A ledger of payments.", "#term.identity.account", "Relationships from [Account (Billing)](<#term.billing.account.relationships>)"} {
		if !strings.Contains(page, want) {
			t.Fatalf("missing %q\n%s", want, page)
		}
	}
	if strings.Count(page, "A ledger of payments.") != 1 || strings.Count(page, "Relationships from ") != 1 || strings.Contains(page, "Relationships from [Invoice]") {
		t.Fatal(page)
	}
	child, err := os.ReadFile("specification/systems/child/glossary.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(child), "Relationships from [Account (Billing)](<#term.billing.account.relationships>)") || strings.Contains(string(child), "Account (Ledger)") {
		t.Fatal(string(child))
	}
	review := run(0, "model", "review")
	if strings.Count(review, "Review required for") != 1 || !strings.Contains(review, "Review required for **balance** in Billing") {
		t.Fatal(review)
	}
	var snapshot model.Snapshot
	exported := run(0, "model", "export")
	writeFixture(t, "snapshot.json", exported)
	run(0, "model", "validate", "snapshot.json")
	if err := json.Unmarshal([]byte(exported), &snapshot); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, term := range snapshot.Elements {
		if len(model.Strings(model.Object(term["relationships"])["references"])) != 1 {
			continue
		}
		found = true
		if len(model.Strings(model.Object(term["illustrative_examples"])["content_nodes"])) == 0 || strings.Join(model.Strings(term["topics"]), ",") != "Customers,Payments" {
			t.Fatal(term)
		}
	}
	if !found {
		t.Fatal("no exported relationship")
	}
	writeFixture(t, "dotspec.yaml", "name: ledger\n")
	for _, bad := range []string{"#missing", "identity.md#missing", "absent.md#account"} {
		writeFixture(t, "billing.md", strings.Replace(billing, "[login]: identity.md#account", "[login]: "+bad, 1))
		if out := run(1, "check"); !strings.Contains(out, "reference.unresolved") || !strings.Contains(out, "billing.md") {
			t.Fatal(out)
		}
	}
	writeFixture(t, "billing.md", strings.Replace(billing, "context: Billing", "context: [Billing]", 1))
	if out := run(1, "check"); !strings.Contains(out, "glossary.context") {
		t.Fatal(out)
	}
}
