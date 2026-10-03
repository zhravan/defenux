package firewall

import "testing"

func TestParseUFWStatus(t *testing.T) {
	if got := parseUFWStatus([]byte("Status: active
")); got != "active" {
		t.Fatalf("got %q", got)
	}
	if got := parseUFWStatus([]byte("Status: inactive
")); got != "inactive" {
		t.Fatalf("got %q", got)
	}
}

func TestParseState(t *testing.T) {
	if got := parseState([]byte("running
")); got != "active" {
		t.Fatalf("got %q", got)
	}
	if got := parseState([]byte("not running
")); got != "inactive" {
		t.Fatalf("got %q", got)
	}
}

func TestRulesState(t *testing.T) {
	if got := rulesState([]byte("")); got != "inactive" {
		t.Fatalf("got %q", got)
	}
	if got := rulesState([]byte("table inet filter")); got != "active" {
		t.Fatalf("got %q", got)
	}
}

func TestParseSpec(t *testing.T) {
	port, protocol, err := ParseSpec("22/tcp")
	if err != nil || port != 22 || protocol != "tcp" {
		t.Fatalf("got %d/%s, %v", port, protocol, err)
	}
}

func TestParseSpecRejectsInvalid(t *testing.T) {
	for _, input := range []string{"0/tcp", "65536/tcp", "22/http", "22"} {
		if _, _, err := ParseSpec(input); err == nil {
			t.Fatalf("expected error for %q", input)
		}
	}
}
