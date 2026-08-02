package router

import (
    "net/netip"
)

// FilterLocalRoutes returns a filtered copy of routes with any prefixes
// that overlap the provided localNetworks removed. It also excludes
// obvious local-only addresses such as loopback and link-local.
func FilterLocalRoutes(routes []netip.Prefix, localNetworks []netip.Prefix) []netip.Prefix {
    var out []netip.Prefix
    for _, r := range routes {
        if r.IsSingleIP() {
            if r.Addr().IsLoopback() || r.Addr().IsLinkLocalUnicast() {
                continue
            }
        }
        if overlapsWithLocalNetwork(r, localNetworks) {
            continue
        }
        out = append(out, r)
    }
    return out
}

// overlapsWithLocalNetwork reports whether pfx overlaps any prefix in
// localNetworks. It treats prefixes as canonical (masked) and detects
// overlap by testing containment of one prefix's base address in the
// other prefix.
func overlapsWithLocalNetwork(pfx netip.Prefix, localNetworks []netip.Prefix) bool {
    p := pfx.Masked()
    for _, l := range localNetworks {
        lp := l.Masked()
        if p.Bits() <= lp.Bits() {
            if lp.Contains(p.Addr()) {
                return true
            }
        } else {
            if p.Contains(lp.Addr()) {
                return true
            }
        }
    }
    return false
}
