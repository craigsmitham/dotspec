package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemAttributeFiles(t *testing.T) {
	root := t.TempDir()
	sysFile, err := InitNamed(root, "Service")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse([]byte(Template + "\n## Quality attributes\n### Old declaration\n")); err == nil {
		t.Fatal("inline attribute accepted")
	}
	write := func(dir, name, meta string) {
		t.Helper()
		path := filepath.Join(root, "specs", "quality", dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
		src := "---\ntype: quality-attribute-spec\n" + meta + "---\n# " + name + "\n## Definition\nConcern.\n"
		if err := os.WriteFile(filepath.Join(path, filepath.Base(dir)+".quality-attribute.spec.md"), []byte(src), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("reliability", "Reliability", "id: dependable\n")
	write("serviceability", "Serviceability", "")
	write("reliability/recovery", "Recovery", "subattribute_of: [serviceability]\n")
	sys, err := Read(sysFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(sys.FindQualityAttribute("recovery").SubattributeOf, ","); got != "serviceability" {
		t.Fatal(got)
	}
	write("serviceability", "Serviceability", "subattribute_of: [recovery]\n")
	if _, err := Read(sysFile); err == nil || !strings.Contains(err.Error(), "cyclic") {
		t.Fatal(err)
	}
	write("serviceability", "Wrong", "")
	if _, err := Read(sysFile); err == nil {
		t.Fatal("accepted mismatched heading")
	}
}
