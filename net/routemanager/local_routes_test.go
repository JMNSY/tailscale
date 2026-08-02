// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package routemanager

import (
	"net/netip"
	"testing"
)

func TestIsLocalPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   netip.Prefix
		wantLocal bool
	}{
		{
			name:      "private RFC1918 10.0.0.0/8",
			prefix:    netip.MustParsePrefix("10.0.0.0/8"),
			wantLocal: true,
		},
		{
			name:      "private RFC1918 172.16.0.0/12",
			prefix:    netip.MustParsePrefix("172.16.0.0/12"),
			wantLocal: true,
		},
		{
			name:      "private RFC1918 192.168.0.0/16",
			prefix:    netip.MustParsePrefix("192.168.0.0/16"),
			wantLocal: true,
		},
		{
			name:      "link-local 169.254.0.0/16",
			prefix:    netip.MustParsePrefix("169.254.0.0/16"),
			wantLocal: true,
		},
		{
			name:      "loopback 127.0.0.1/32",
			prefix:    netip.MustParsePrefix("127.0.0.1/32"),
			wantLocal: true,
		},
		{
			name:      "IPv6 link-local fe80::/10",
			prefix:    netip.MustParsePrefix("fe80::/10"),
			wantLocal: true,
		},
		{
			name:      "public route 8.8.8.0/24",
			prefix:    netip.MustParsePrefix("8.8.8.0/24"),
			wantLocal: false,
		},
		{
			name:      "public route 1.1.1.0/24",
			prefix:    netip.MustParsePrefix("1.1.1.0/24"),
			wantLocal: false,
		},
	}

	rm := &RouteManager{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rm.isLocalPrefix(tt.prefix)
			if got != tt.wantLocal {
				t.Errorf("isLocalPrefix(%v) = %v, want %v", tt.prefix, got, tt.wantLocal)
			}
		})
	}
}

func TestIsLinkLocal(t *testing.T) {
	tests := []struct {
		name          string
		addr          netip.Addr
		wantLinkLocal bool
	}{
		{
			name:          "IPv4 link-local 169.254.1.1",
			addr:          netip.MustParseAddr("169.254.1.1"),
			wantLinkLocal: true,
		},
		{
			name:          "IPv4 link-local 169.254.255.255",
			addr:          netip.MustParseAddr("169.254.255.255"),
			wantLinkLocal: true,
		},
		{
			name:          "IPv6 link-local fe80::1",
			addr:          netip.MustParseAddr("fe80::1"),
			wantLinkLocal: true,
		},
		{
			name:          "IPv4 not link-local 169.255.0.0",
			addr:          netip.MustParseAddr("169.255.0.0"),
			wantLinkLocal: false,
		},
		{
			name:          "IPv4 not link-local 192.168.1.1",
			addr:          netip.MustParseAddr("192.168.1.1"),
			wantLinkLocal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLinkLocal(tt.addr)
			if got != tt.wantLinkLocal {
				t.Errorf("isLinkLocal(%v) = %v, want %v", tt.addr, got, tt.wantLinkLocal)
			}
		})
	}
}
