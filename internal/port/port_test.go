package port

import "testing"

func TestSplitAddress(t *testing.T) {
	for input, want := range map[string][2]string{
		"0.0.0.0:8000": {"0.0.0.0", "8000"},
		"[::]:22": {"::", "22"},
	} {
		address, port := splitAddress(input)
		if address != want[0] || port != want[1] {
			t.Fatalf("%q: got %q:%q, want %q:%q", input, address, port, want[0], want[1])
		}
	}
}
