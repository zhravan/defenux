//go:build linux

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestUserList(t *testing.T) {
	if _, err := exec.LookPath("getent"); err != nil {
		t.Skip("getent is not installed")
	}

	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binary := filepath.Join(root, "..", "defenux")
	cmd := exec.Command(binary, "user", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("defenux user list failed: %v: %s", err, out)
	}

	if !strings.Contains(string(out), "root") {
		t.Fatalf("expected root account in output: %s", out)
	}
}
