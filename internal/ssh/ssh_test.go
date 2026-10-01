package ssh

import "testing"

func TestParseConfig(t *testing.T) {
	config := parseConfig([]byte(`Port 22
Port 2222
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
PermitEmptyPasswords no
X11Forwarding yes
MaxAuthTries 4
AllowUsers alice bob
AllowGroups admins
`))

	if len(config.Ports) != 2 || config.Ports[0] != 22 || config.Ports[1] != 2222 {
		t.Fatalf("ports: %#v", config.Ports)
	}
	if config.PermitRootLogin != "no" {
		t.Fatalf("permitrootlogin: %q", config.PermitRootLogin)
	}
	if config.PasswordAuthentication != "no" {
		t.Fatalf("passwordauthentication: %q", config.PasswordAuthentication)
	}
	if config.PubkeyAuthentication != "yes" {
		t.Fatalf("pubkeyauthentication: %q", config.PubkeyAuthentication)
	}
	if config.PermitEmptyPasswords != "no" {
		t.Fatalf("permitemptypasswords: %q", config.PermitEmptyPasswords)
	}
	if config.X11Forwarding != "yes" {
		t.Fatalf("x11forwarding: %q", config.X11Forwarding)
	}
	if config.MaxAuthTries != "4" {
		t.Fatalf("maxauthtries: %q", config.MaxAuthTries)
	}
	if len(config.AllowUsers) != 2 || config.AllowUsers[0] != "alice" || config.AllowUsers[1] != "bob" {
		t.Fatalf("allowusers: %#v", config.AllowUsers)
	}
	if len(config.AllowGroups) != 1 || config.AllowGroups[0] != "admins" {
		t.Fatalf("allowgroups: %#v", config.AllowGroups)
	}
}

func TestParseConfigDefaultsPort(t *testing.T) {
	config := parseConfig([]byte("PasswordAuthentication yes\n"))
	if len(config.Ports) != 1 || config.Ports[0] != 22 {
		t.Fatalf("ports: %#v", config.Ports)
	}
}
