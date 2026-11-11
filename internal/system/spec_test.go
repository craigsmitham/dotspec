package system

import (
	"strings"
	"testing"
)

const validSpec = "---\ntype: system\n---\n\n# Parser\n\n## Purpose\n\nRead specs.\n\n## Scope\n\nMarkdown files.\n"

func TestParse(t *testing.T) {
	tests := []struct {
		name, source, title, purpose, scope string
	}{
		{"minimal", validSpec, "Parser", "Read specs.", "Markdown files."},
		{"CRLF and BOM", "\xef\xbb\xbf" + strings.ReplaceAll(validSpec, "\n", "\r\n"), "Parser", "Read specs.", "Markdown files."},
		{"quoted YAML", strings.Replace(validSpec, "type: system", "type: 'system' # comment", 1), "Parser", "Read specs.", "Markdown files."},
		{"formatted title", strings.Replace(validSpec, "# Parser", "# **Spec** `parser`", 1), "Spec parser", "Read specs.", "Markdown files."},
		{"setext headings", "---\ntype: system\n---\nParser\n======\n\nPurpose\n-------\nRead specs.\n\nScope\n-----\nMarkdown files.", "Parser", "Read specs.", "Markdown files."},
		{"closing heading markers", strings.ReplaceAll(strings.ReplaceAll(validSpec, "## Purpose", "## Purpose ##"), "## Scope", "## Scope ##"), "Parser", "Read specs.", "Markdown files."},
		{"scope first", "---\ntype: system\n---\n# Parser\n## Scope\nMarkdown files.\n## Purpose\nRead specs.", "Parser", "Read specs.", "Markdown files."},
		{"extra section", validSpec + "\n## Notes\nNot part of scope.\n", "Parser", "Read specs.", "Markdown files."},
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
		{"unclosed frontmatter", "---\ntype: system", "must end"},
		{"invalid YAML", strings.Replace(validSpec, "type: system", "type: [", 1), "invalid YAML"},
		{"duplicate YAML key", strings.Replace(validSpec, "type: system", "type: system\ntype: system", 1), "invalid YAML"},
		{"unknown metadata", strings.Replace(validSpec, "type: system", "type: system\nparent: root", 1), "invalid YAML"},
		{"wrong type", strings.Replace(validSpec, "type: system", "type: job", 1), "type: system"},
		{"missing type", strings.Replace(validSpec, "type: system", "{}", 1), "type: system"},
		{"no title", strings.Replace(validSpec, "# Parser", "Parser", 1), "name heading"},
		{"empty title", strings.Replace(validSpec, "# Parser", "#", 1), "nonempty top-level"},
		{"two titles", validSpec + "\n# Another\n", "exactly one"},
		{"title in code", strings.Replace(validSpec, "# Parser", "```\n# Parser\n```", 1), "name heading"},
		{"no purpose", strings.Replace(validSpec, "## Purpose", "## Rationale", 1), "Purpose"},
		{"empty purpose", strings.Replace(validSpec, "Read specs.", "", 1), "Purpose"},
		{"no scope", strings.Replace(validSpec, "## Scope", "### Scope", 1), "Scope"},
		{"empty scope", strings.Replace(validSpec, "Markdown files.", "", 1), "Scope"},
		{"duplicate section", validSpec + "\n## Scope\nAgain.\n", "duplicate Scope"},
		{"wrong heading case", strings.Replace(validSpec, "## Purpose", "## purpose", 1), "Purpose"},
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
	for _, seed := range []string{validSpec, validSpec + qualityExample, validSpec + "\n## Quality attributes\n### Cycle\n**ID:** cycle\n**Sub-attribute of:** cycle", Template, "", "---\ntype: system\n---\n#\n## Purpose\n## Scope"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Parse(data)
	})
}

func TestSystemAndDeclarationIdentities(t *testing.T) {
	source := strings.Replace(validSpec, "type: system", "type: system\nid: parser", 1) + "\n## Components\n### Scanner\n**ID:** scanner\n#### Part of\n[Parser](#parser)\n## Features\n### Scan files\n**ID:** scan-files\n#### Notes\n[Scanner heading](#scanner)\n"
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
