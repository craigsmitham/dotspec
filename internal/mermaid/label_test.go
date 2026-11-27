package mermaid

import (
	"strings"
	"testing"
)

func TestLabelEncodesDiagramSyntax(t *testing.T) {
	got := Label("A \"quoted\" <tag> | [edge] (shape) &\nnext")
	for _, raw := range []string{`"`, "<", ">", "|", "[", "]", "(", ")", "&", "\n"} {
		if strings.Contains(got, raw) {
			t.Errorf("raw %q in label %q", raw, got)
		}
	}
	if want := "A #34;quoted#34; #60;tag#62; #124; #91;edge#93; #40;shape#41; #38; next"; got != want {
		t.Errorf("Label = %q, want %q", got, want)
	}
}
