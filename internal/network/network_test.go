package network

import "testing"

func TestParseInterfaces(t *testing.T) {
	input := []byte("2: eth0: <BROADCAST> mtu 1500 state UP group default qlen 1000 link/ether aa:bb:cc:dd:ee:ff\n")
	got := parseInterfaces(input)
	if len(got) != 1 || got[0].Name != "eth0" || got[0].State != "UP" || got[0].MAC != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("unexpected interface: %+v", got)
	}
}

func TestParseRoutes(t *testing.T) {
	input := []byte("default via 192.168.1.1 dev eth0 proto dhcp\n192.168.1.0/24 dev eth0 proto kernel\n")
	got := parseRoutes(input)
	if len(got) != 2 {
		t.Fatalf("got %d routes", len(got))
	}
	if got[0].Gateway != "192.168.1.1" || got[0].Device != "eth0" {
		t.Fatalf("unexpected default route: %+v", got[0])
	}
}
