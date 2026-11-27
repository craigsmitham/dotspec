package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func invoke(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func writeFixture(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestHelpAndUsageErrors(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"new", "--help"}, {"new", "system", "--help"}, {"list", "--help"}, {"list", "system", "--help"}} {
		code, out, errOut := invoke(args...)
		if code != 0 || !strings.Contains(out, "Usage:") || errOut != "" {
			t.Errorf("help %v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{{"bogus"}, {"new"}, {"list"}, {"new", "bogus"}, {"list", "system", "unexpected"}, {"list", "system", "--unknown"}, {"-C"}, {"-C", "", "list", "system"}} {
		code, out, errOut := invoke(args...)
		if code == 0 || out != "" || !strings.HasPrefix(errOut, "dotspec:") {
			t.Errorf("usage error %v: %d %q %q", args, code, out, errOut)
		}
	}
}
