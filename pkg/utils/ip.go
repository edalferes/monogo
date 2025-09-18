package utils

import "net"

// ToIPv4 convert an IP to IPv4 string if possible, or returns the original
func ToIPv4(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ip
	}
	if ipv4 := parsed.To4(); ipv4 != nil {
		return ipv4.String()
	}
	if parsed.IsLoopback() {
		return "127.0.0.1"
	}
	return ip
}
