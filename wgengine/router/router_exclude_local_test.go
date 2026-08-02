// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package router

import (
	"net/netip"
	"testing"
)

func TestFilterLocalRoutes(t *testing.T) {
	localPrefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("192.168.0.0/16"),
		netip.MustParsePrefix("fe80::/10"),
	}

	tests := []struct {
		name     string
		routes   []netip.Prefix
		want     []netip.Prefix
	}{
		{
			name: "no overlapping routes",
			routes: []netip.Prefix{
				netip.MustParsePrefix("8.8.8.0/24"),
				netip.MustParsePrefix("1.1.1.0/24"),
			},
			want: []netip.Prefix{
				netip.MustParsePrefix("8.8.8.0/24"),
				netip.MustParsePrefix("1.1.1.0/24"),
			},
		},
		{
			name: "filter overlapping local routes",
			routes: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/8"),  // overlaps with local 10.0.0.0/8
				netip.MustParsePrefix("8.8.8.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"), // overlaps with local 192.168.0.0/16
				netip.MustParsePrefix("1.1.1.0/24"),
			},
			want: []netip.Prefix{
				netip.MustParsePrefix("8.8.8.0/24"),
				netip.MustParsePrefix("1.1.1.0/24"),
			},
		},
		{
			name: "empty routes",
			routes: []netip.Prefix{},
			want:   []netip.Prefix{},
		},
		{
			name: "all routes local",
			routes: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/8"),
				netip.MustParsePrefix("192.168.0.0/16"),
			},
			want: []netip.Prefix{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterLocalRoutes(tt.routes, localPrefixes)
			if len(got) != len(tt.want) {
				t.Errorf("FilterLocalRoutes() returned %d routes, want %d", len(got), len(tt.want))
				return
			}
			for i, route := range got {
				if route != tt.want[i] {
					t.Errorf("FilterLocalRoutes()[%d] = %v, want %v", i, route, tt.want[i])
				}
			}
		})
	}
}

func TestOverlapsWithLocalNetwork(t *testing.T) {
	localPrefixes := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("192.168.0.0/16"),
	}

	tests := []struct {
		name      string
		route     netip.Prefix
		wantLocal bool
	}{
		{
			name:      "exact match 10.0.0.0/8",
			route:     netip.MustParsePrefix("10.0.0.0/8"),
			wantLocal: true,
		},
		{
			name:      "subnet of local 10.1.0.0/16",
			route:     netip.MustParsePrefix("10.1.0.0/16"),
			wantLocal: true,
		},
		{
			name:      "subnet of local 192.168.1.0/24",
			route:     netip.MustParsePrefix("192.168.1.0/24"),
			wantLocal: true,
		},
		{
			name:      "no overlap 8.8.8.0/24",
			route:     netip.MustParsePrefix("8.8.8.0/24"),
			wantLocal: false,
		},
		{
			name:      "no overlap 172.16.0.0/12",
			route:     netip.MustParsePrefix("172.16.0.0/12"),
			wantLocal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := overlapsWithLocalNetwork(tt.route, localPrefixes)
			if got != tt.wantLocal {
				t.Errorf("overlapsWithLocalNetwork(%v) = %v, want %v", tt.route, got, tt.wantLocal)
			}
		})
	}
}
