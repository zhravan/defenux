package service

import "testing"

func TestParseListUnits(t *testing.T) {
	data := []byte(`ssh.service loaded active running OpenSSH server daemon
cups.service loaded inactive dead CUPS Scheduler
`)

	services := parseListUnits(data, map[string]string{
		"ssh.service":  "enabled",
		"cups.service": "disabled",
	})

	if len(services) != 2 {
		t.Fatalf("got %d services", len(services))
	}
	if services[0].Name != "ssh.service" {
		t.Fatalf("name: %q", services[0].Name)
	}
	if services[0].State != "active" || services[0].SubState != "running" {
		t.Fatalf("state: %q/%q", services[0].State, services[0].SubState)
	}
	if services[0].Enabled != "enabled" {
		t.Fatalf("enabled: %q", services[0].Enabled)
	}
	if services[0].Description != "OpenSSH server daemon" {
		t.Fatalf("description: %q", services[0].Description)
	}
}

func TestParseShow(t *testing.T) {
	data := []byte(`Id=ssh.service
Description=OpenSSH server daemon
ActiveState=active
SubState=running
UnitFileState=enabled
MainPID=1234
Type=notify
`)

	item := parseShow(data)

	if item.Name != "ssh.service" || item.State != "active" || item.SubState != "running" {
		t.Fatalf("item: %#v", item)
	}
	if item.Enabled != "enabled" || item.MainPID != 1234 || item.Type != "notify" {
		t.Fatalf("details: %#v", item)
	}
}
