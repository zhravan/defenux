//go:build linux

package e2e

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPortInspection(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	portText := fmt.Sprintf("%d", port)

	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "..", "defenux")

	cmd := exec.Command(binary, "port", "status", portText)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("defenux failed: %v: %s", err, out)
	}
	if !strings.Contains(string(out), portText) {
		t.Fatalf("expected port %d in output: %s", port, out)
	}
}
