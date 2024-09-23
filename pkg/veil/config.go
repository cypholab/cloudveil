package veil

import "net/netip"

type Config struct {
	CheckScheme  string
	IpAddresses  []netip.Addr
	ExpectedBody string
	Hostname     string
	RateLimit    int
	Timeout      int
}
