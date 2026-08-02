# Local Route Exclusion Feature Implementation

## Problem Statement

When using Tailscale with `--accept-routes`, users may experience routing conflicts when:
- Local machine is on a subnet (e.g., 10.0.0.x)
- A Tailscale peer advertises a broader route (e.g., 10.0.0.0/8)
- The advertised route takes priority, forcing traffic through the Tailscale tunnel unnecessarily

## Solution

Automatically exclude routes that overlap with local network interfaces from `accept-routes`, ensuring local connectivity is preserved.

## Implementation Details

### Files Added/Modified

1. **net/routemanager/local_routes.go** (NEW)
   - `isLocalPrefix()`: Check if a route overlaps with local networks
   - `getLocalNetworkPrefixes()`: Collect all local network prefixes
   - `isLocalRoute()`: Determine if a route is local
   - `isLinkLocal()`: Check for link-local addresses

2. **wgengine/router/router_exclude_local.go** (NEW)
   - `FilterLocalRoutes()`: Main filtering function
   - `overlapsWithLocalNetwork()`: Overlap detection

3. **.github/workflows/build-windows.yml** (NEW)
   - Automated Windows build for both amd64 and arm64
   - Test execution
   - Artifact upload for easy testing

### Route Types Excluded

The implementation excludes routes in these categories:
- **RFC 1918**: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
- **RFC 3927**: 169.254.0.0/16 (link-local)
- **RFC 4291**: fe80::/10 (IPv6 link-local)
- **Loopback**: 127.0.0.0/8, ::1/128

## Testing

### Manual Testing (Windows)

```powershell
# 1. Build the feature branch
git clone https://github.com/JMNSY/tailscale.git
cd tailscale
git checkout feature/exclude-local-routes
go build -o tailscale.exe ./cmd/tailscale
go build -o tailscaled.exe ./cmd/tailscaled

# 2. Test local network connectivity
# Before: observe routes being overridden
# After: verify local routes are preserved

# 3. Check route table
route print

# 4. Monitor Tailscale logs
.\tailscaled.exe -v=1 2>&1 | Select-String "exclude|local|route"
```

### Automated Testing (CI/CD)

The workflow runs:
- `go test ./net/routemanager/...`
- `go test ./wgengine/router/...`
- Builds both amd64 and arm64 Windows binaries
- Artifacts available for 30 days

## Integration Points

### RouteManager Integration
```go
// In deriveOSRoutes() or similar:
for pfx := range routes {
    if rm.isLocalPrefix(pfx) {
        continue  // Skip local routes
    }
    acceptedRoutes = append(acceptedRoutes, pfx)
}
```

### Config Integration
```go
// In ipn/prefs.go (optional):
// Add configuration option to control this behavior:
type Prefs struct {
    // ...
    ExcludeLocalRoutes bool  // New field
}
```

## Performance Considerations

- Minimal overhead: O(n*m) where n = advertised routes, m = local routes
- Typical case: 1-5 advertised routes, 1-3 local routes
- Local prefix detection cached on route table updates

## Backward Compatibility

- Enabled by default (safest behavior)
- Optional configuration field for fine-tuning
- No breaking changes to existing APIs
- Existing `--accept-routes` behavior preserved for non-overlapping routes

## Next Steps

1. **Code Integration**: Integrate local route detection into main routing loop
2. **Testing**: Validate with various network topologies
3. **Documentation**: Update user docs and release notes
4. **Upstream**: Consider PR to tailscale/tailscale for broader adoption

## References

- RFC 1918: Private IPv4 Address Space
- RFC 3927: Dynamic Host Configuration Protocol (DHCP) Link-Local Addresses
- RFC 4291: IP Version 6 Addressing Architecture
- Tailscale Issues: Route management and overlapping networks
