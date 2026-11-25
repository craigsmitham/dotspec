package system

import (
	"strings"
	"testing"
)

const validSpec = "---\ntype: system-spec\n---\n\n# Parser\n\n## Purpose\n\nRead specs.\n\n## Scope\n\nMarkdown files.\n"

func TestParse(t *testing.T) {
	tests := []struct {
		name, source, title, purpose, scope string
	}{
		{"minimal", validSpec, "Parser", "Read specs.", "Markdown files."},
		{"empty scope", strings.Replace(validSpec, "Markdown files.", "", 1), "Parser", "Read specs.", ""},
		{"omitted scope", strings.Split(validSpec, "## Scope")[0], "Parser", "Read specs.", ""},
		{"CRLF and BOM", "\xef\xbb\xbf" + strings.ReplaceAll(validSpec, "\n", "\r\n"), "Parser", "Read specs.", "Markdown files."},
		{"quoted YAML", strings.Replace(validSpec, "type: system-spec", "type: 'system-spec' # comment", 1), "Parser", "Read specs.", "Markdown files."},
		{"formatted title", strings.Replace(validSpec, "# Parser", "# **Spec** `parser`", 1), "Spec parser", "Read specs.", "Markdown files."},
		{"setext headings", "---\ntype: system-spec\n---\nParser\n======\n\nPurpose\n-------\nRead specs.\n\nScope\n-----\nMarkdown files.", "Parser", "Read specs.", "Markdown files."},
		{"closing heading markers", strings.ReplaceAll(strings.ReplaceAll(validSpec, "## Purpose", "## Purpose ##"), "## Scope", "## Scope ##"), "Parser", "Read specs.", "Markdown files."},
		{"scope first", "---\ntype: system-spec\n---\n# Parser\n## Scope\nMarkdown files.\n## Purpose\nRead specs.", "Parser", "Read specs.", "Markdown files."},
		{"extra section", validSpec + "\n## Evidence\nNot part of scope.\n", "Parser", "Read specs.", "Markdown files."},
		{"Markdown body", strings.Replace(validSpec, "Read specs.", "- Read specs.\n\n### Detail\n\nPreserve structure.\n\n```markdown\n# Fake title\n## Scope\n```\n\n> ## Scope\n> Quoted heading.", 1), "Parser", "- Read specs.\n\n### Detail\n\nPreserve structure.\n\n```markdown\n# Fake title\n## Scope\n```\n\n> ## Scope\n> Quoted heading.", "Markdown files."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := Parse([]byte(tt.source))
			if err != nil {
				t.Fatal(err)
			}
			if spec.Name != tt.title || spec.Purpose != tt.purpose || spec.Scope != tt.scope {
				t.Fatalf("unexpected parsed system: %#v", spec)
			}
		})
	}
}

func TestParseRejectsMalformedSpecs(t *testing.T) {
	tests := []struct{ name, source, want string }{
		{"empty", "", "frontmatter"},
		{"no frontmatter", "# Parser", "frontmatter"},
		{"unclosed frontmatter", "---\ntype: system-spec", "must end"},
		{"invalid YAML", strings.Replace(validSpec, "type: system-spec", "type: [", 1), "invalid YAML"},
		{"duplicate YAML key", strings.Replace(validSpec, "type: system-spec", "type: system-spec\ntype: system-spec", 1), "invalid YAML"},
		{"unknown metadata", strings.Replace(validSpec, "type: system-spec", "type: system-spec\nparent: root", 1), "invalid YAML"},
		{"wrong type", strings.Replace(validSpec, "type: system-spec", "type: job-spec", 1), "type: system-spec"},
		{"bare subject kind", strings.Replace(validSpec, "type: system-spec", "type: system", 1), "type: system-spec"},
		{"generated report type", strings.Replace(validSpec, "type: system-spec", "type: system-spec-report", 1), "type: system-spec"},
		{"missing type", strings.Replace(validSpec, "type: system-spec", "{}", 1), "type: system-spec"},
		{"no title", strings.Replace(validSpec, "# Parser", "Parser", 1), "name heading"},
		{"empty title", strings.Replace(validSpec, "# Parser", "#", 1), "nonempty top-level"},
		{"two titles", validSpec + "\n# Another\n", "exactly one"},
		{"title in code", strings.Replace(validSpec, "# Parser", "```\n# Parser\n```", 1), "name heading"},
		{"no purpose", strings.Replace(validSpec, "## Purpose", "## Rationale", 1), "Rationale"},
		{"empty purpose", strings.Replace(validSpec, "Read specs.", "", 1), "Purpose"},
		{"duplicate section", validSpec + "\n## Scope\nAgain.\n", "duplicate Scope"},
		{"wrong heading case", strings.Replace(validSpec, "## Purpose", "## purpose", 1), "purpose"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.source))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got error %v; want containing %q", err, tt.want)
			}
		})
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{validSpec, validSpec + "\n## Quality attributes\n### Cycle\n**ID:** cycle\n**Sub-attribute of:** cycle", Template, "", "---\ntype: system-spec\n---\n#\n## Purpose\n## Scope"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Parse(data)
	})
}

func TestSystemAndDeclarationIdentities(t *testing.T) {
	source := strings.Replace(validSpec, "type: system-spec", "type: system-spec\nid: parser", 1) + "\n## Concept uses\n### Scanner\n**ID:** scanner\n#### Description\n[Parser](#parser)\n## Synchronizations\n### Scan files\n**ID:** scan-files\n#### Semantics\n[Scanner heading](#scanner)\n"
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if spec.ID != "parser" || !spec.References.Anchors["scanner"] || !spec.References.Anchors["scan-files"] {
		t.Fatal(spec)
	}
	for _, bad := range []string{
		strings.Replace(source, "**ID:** scanner", "**ID:** parser", 1),
		strings.Replace(source, "**ID:** scanner", "**ID:**", 1),
		strings.Replace(source, "**ID:** scanner", "**ID:** bad.id", 1),
		strings.Replace(source, "**ID:** scanner", "**ID:** scanner\n**ID:** another", 1),
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Fatal("accepted invalid declaration identity", bad)
		}
	}
}

func TestFeaturesAndCapabilitiesRequireStandaloneSources(t *testing.T) {
	for _, collection := range []string{"Features", "Capabilities"} {
		section := "\n## " + collection + "\n\n### Resume work\n"
		if _, err := Parse([]byte(validSpec + section)); err == nil {
			t.Fatalf("accepted inline %s on a system", collection)
		}
		component := "---\ntype: component-spec\n---\n# Worker\n\n## Part of\n\n[System](example.system.spec.md)\n"
		if _, err := ParseComponent([]byte(component + section)); err == nil {
			t.Fatalf("accepted inline %s on a component", collection)
		}
	}
}
