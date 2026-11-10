// Command check-format performs a portable, read-only formatting check.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("gofmt", "-l", "main.go", "internal", "scripts")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofmt: %v\n%s", err, output)
		os.Exit(1)
	}
	if len(bytes.TrimSpace(output)) != 0 {
		fmt.Fprintf(os.Stderr, "Run mise run fmt:\n%s", output)
		os.Exit(1)
	}
}
