package utils

import "testing"

func TestToIPv4(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"::1", "127.0.0.1"},             // IPv6 loopback
		{"127.0.0.1", "127.0.0.1"},       // IPv4 loopback
		{"192.168.1.10", "192.168.1.10"}, // IPv4
		{"2001:db8::1", "2001:db8::1"},   // IPv6 (not loopback)
		{"invalid", "invalid"},           // Invalid IP
	}
	for _, c := range cases {
		got := ToIPv4(c.input)
		if got != c.expected {
			t.Errorf("ToIPv4(%q) = %q; want %q", c.input, got, c.expected)
		}
	}
}
