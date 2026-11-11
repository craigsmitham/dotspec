package cli

import (
	"os"
	"strings"
	"testing"
)

func TestJTBDAuthoringPresetsAndNeedsView(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	for _, preset := range []string{"odi", "progress", "job-story"} {
		code, out, errOut := invoke("job", "init", preset, "--preset", preset)
		if code != 0 {
			t.Fatal(out, errOut)
		}
		code, out, errOut = invoke("job", "check", preset, "--preset", preset)
		if code != 0 {
			t.Fatal(out, errOut)
		}
	}
	code, _, _ := invoke("job", "init", "bad", "--preset", "unknown")
	if code == 0 {
		t.Fatal("invalid preset accepted")
	}
	if _, err := os.Stat("specs/jobs/bad.job.spec.md"); !os.IsNotExist(err) {
		t.Fatal("invalid preset wrote a file")
	}
	code, out, errOut := invoke("model", "needs")
	if code != 0 || !strings.Contains(out, "# Needs and specification map") || !strings.Contains(out, "DesiredOutcome: Success criterion") {
		t.Fatal(code, out, errOut)
	}
}
