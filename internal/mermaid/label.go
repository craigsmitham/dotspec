package mermaid

import (
	"fmt"
	"strings"
	"unicode"
)

// Mermaid's decimal entity notation prevents labels from becoming diagram syntax.
func Label(label string) string {
	var output strings.Builder
	for _, r := range strings.Join(strings.Fields(label), " ") {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '.' || r == ',' || r == ':' {
			output.WriteRune(r)
		} else {
			fmt.Fprintf(&output, "#%d;", r)
		}
	}
	return output.String()
}
