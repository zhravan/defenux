package user

import "testing"

func TestParsePasswd(t *testing.T) {
	data := []byte(`root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
alice:x:1000:1000:Alice:/home/alice:/bin/bash
`)

	users := parsePasswd(data, map[string]bool{"alice": true})

	if len(users) != 3 {
		t.Fatalf("got %d users", len(users))
	}
	if users[0].Class != "root" {
		t.Fatalf("root class: %q", users[0].Class)
	}
	if users[1].Class != "system" {
		t.Fatalf("system class: %q", users[1].Class)
	}
	if users[2].Class != "human" || users[2].Admin != "yes" {
		t.Fatalf("human/admin: %#v", users[2])
	}
}

func TestParsePasswordState(t *testing.T) {
	tests := map[string]string{
		"alice P 01-01-2026 0 99999 7 -1": "active",
		"alice L 01-01-2026 0 99999 7 -1": "locked",
		"alice NP 01-01-2026 0 99999 7 -1": "no-password",
		"alice LK 01-01-2026 0 99999 7 -1": "locked",
	}

	for input, want := range tests {
		if got := parsePasswordState([]byte(input)); got != want {
			t.Fatalf("%q: got %q want %q", input, got, want)
		}
	}
}

func TestParsePasswordAging(t *testing.T) {
	data := []byte(`Last password change                                    : Sep 02, 2026
Password expires                                        : never
Account expires                                         : never
`)

	changed, expires, account := parsePasswordAging(data)
	if changed != "Sep 02, 2026" || expires != "never" || account != "never" {
		t.Fatalf("got %q/%q/%q", changed, expires, account)
	}
}
