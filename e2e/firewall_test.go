//go:build linux

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFirewallStatus(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "..", "defenux")

	cmd := exec.Command(binary, "firewall", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("defenux firewall status failed: %v: %s", err, out)
	}
	text := string(out)
	if !strings.Contains(text, "backend:") || !strings.Contains(text, "state:") {
		t.Fatalf("unexpected output: %s", text)
	}
}
