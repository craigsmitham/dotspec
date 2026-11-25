package system

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func TestOptionalSystemContext(t *testing.T) {
	source := Template + "\n## Goals\n### Review quality\nImprove reviews.\n## Environment\nGit repositories.\n## Assumptions\n### Stable sources\nSources do not change during checks.\n## Constraints\nPortable Markdown.\n## Glossary\n### Complete scan\nAll selected files were examined.\n"
	s, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(s.Goals, "Improve reviews.") || s.Environment == "" || s.Assumptions == "" || s.Constraints == "" || !strings.Contains(s.Glossary, "### Complete scan") || !s.References.Anchors["complete-scan"] {
		t.Fatalf("context: %+v", s)
	}
	for _, name := range []string{"Description", "Goals", "Environment", "Assumptions", "Constraints", "Glossary"} {
		for _, bad := range []string{source + "\n## " + name + "\n", strings.Replace(Template, "# New system", "# New system\n\n### "+name, 1)} {
			_, err := Parse([]byte(bad))
			var d *diagnostic.Diagnostic
			if !errors.As(err, &d) || !strings.HasPrefix(d.Rule, "structure.") {
				t.Errorf("accepted bad context %s: %v", name, err)
			}
		}
	}
}

func TestDirectoryContextUsesObservedSystem(t *testing.T) {
	root := t.TempDir()
	path := writeSpec(t, root, "Workspace")
	source := filepath.Join(root, "specs", "requirements", "save.requirement.spec.md")
	writeFile(t, source, "---\ntype: requirement-spec\n---\n# Save\n")
	files := &specfile.SourceSet{}
	if _, err := ReadFrom(files, path); err != nil {
		t.Fatal(err)
	}
	writeFile(t, path, "Incomplete edit")

	paths, err := CheckContext.ResolveSource(files, root, source)
	if err != nil || !slices.Equal(paths, []string{path}) {
		t.Fatalf("context from observed sources = %v, %v", paths, err)
	}
	_, err = CheckContext.ResolveSource(&specfile.SourceSet{}, root, source)
	var finding *diagnostic.Diagnostic
	if !errors.As(err, &finding) || finding.Rule != "metadata.frontmatter" || finding.Location.Path != path {
		t.Fatalf("fresh context should diagnose the edit at %s, got %v", path, err)
	}
}
