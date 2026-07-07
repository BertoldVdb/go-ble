package bleatt

import (
	"sync"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// InMemoryDiscoveryCache is a process-lifetime GATT discovery cache
// keyed by the peer's connection-time BLE address. Construct one per
// stack and wire its method values into the relevant hooks:
//
//	cache := bleatt.NewInMemoryDiscoveryCache()
//	attCfg.DiscoveryCacheGet = cache.Get
//	attCfg.DiscoveryCacheSet = cache.Set
//	smpCfg.OnRebond           = cache.Evict
//
// The cache stores entries by reference — bleatt does not mutate
// CachedGATT after handing it to Set, and Get returns the same pointer
// it received, so no defensive copy is performed.
//
// Caveats:
//   - No TTL. Entries live until evicted (validation failure, OnRebond,
//     or process exit). Wrap or replace if time-based expiry matters.
//   - No capacity cap. Sized by the number of distinct peers.
//   - Keyed by raw BLEAddr (MAC + AddrType). Peers that connect with a
//     Resolvable Private Address rotate their address per session and
//     will miss every time; identity resolution is the caller's problem.
//   - In-memory only. Process restart wipes the cache.
type InMemoryDiscoveryCache struct {
	mu      sync.Mutex
	entries map[bleutil.BLEAddr]*CachedGATT
}

// NewInMemoryDiscoveryCache returns an empty cache ready for use.
func NewInMemoryDiscoveryCache() *InMemoryDiscoveryCache {
	return &InMemoryDiscoveryCache{
		entries: make(map[bleutil.BLEAddr]*CachedGATT),
	}
}

// peerKey extracts the cache key from a GattDevice. Returns the zero
// BLEAddr (with a false ok) when the device has no underlying BLE
// connection or its remote address is not a BLEAddr — in that case
// the cache silently no-ops, since there is nothing meaningful to key
// on.
func peerKey(dev *GattDevice) (bleutil.BLEAddr, bool) {
	if dev == nil {
		return bleutil.BLEAddr{}, false
	}
	conn := dev.BLEConnection()
	if conn == nil {
		return bleutil.BLEAddr{}, false
	}
	addr, ok := conn.RemoteAddr().(bleutil.BLEAddr)
	if !ok {
		return bleutil.BLEAddr{}, false
	}
	return addr, true
}

// Get returns the cached entry for the device's peer, or nil if the
// peer is not cached. Matches the signature of DiscoveryCacheGet.
func (c *InMemoryDiscoveryCache) Get(dev *GattDevice) *CachedGATT {
	addr, ok := peerKey(dev)
	if !ok {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entries[addr]
}

// Set stores the entry for the device's peer. A nil cached argument
// removes any existing entry — bleatt uses this to signal that
// validation detected a stale cache and the entry should be evicted.
// Matches the signature of DiscoveryCacheSet.
func (c *InMemoryDiscoveryCache) Set(dev *GattDevice, cached *CachedGATT) {
	addr, ok := peerKey(dev)
	if !ok {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached == nil {
		delete(c.entries, addr)
		return
	}
	c.entries[addr] = cached
}

// Evict removes the cache entry for the given peer address. Matches
// the signature of SMPConfig.OnRebond, so it can be wired in directly:
// when SMP detects a fresh pairing replaced an existing bond, the
// peer may have forgotten our subscriptions and reshuffled its
// attribute layout, and the cached structure is no longer trustworthy.
func (c *InMemoryDiscoveryCache) Evict(addr bleutil.BLEAddr) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, addr)
}

// Len returns the number of cached peers. Useful for debugging and
// metrics; not load-bearing for cache behaviour.
func (c *InMemoryDiscoveryCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
