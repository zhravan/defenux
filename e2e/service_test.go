//go:build linux

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestServiceList(t *testing.T) {
	if _, err := exec.LookPath("systemctl"); err != nil {
		t.Skip("systemctl is not installed")
	}

	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binary := filepath.Join(root, "..", "defenux")
	cmd := exec.Command(binary, "service", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("defenux service list failed: %v: %s", err, out)
	}
	if strings.TrimSpace(string(out)) == "" {
		t.Fatal("expected service output")
	}
}
