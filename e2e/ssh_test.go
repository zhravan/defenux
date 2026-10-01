//go:build linux

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSSHStatus(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binary := filepath.Join(root, "..", "defenux")
	cmd := exec.Command(binary, "ssh", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("defenux ssh status failed: %v: %s", err, out)
	}

	text := string(out)
	for _, field := range []string{"service:", "state:", "config:", "ports:"} {
		if !strings.Contains(text, field) {
			t.Fatalf("missing %s in output: %s", field, text)
		}
	}
}
