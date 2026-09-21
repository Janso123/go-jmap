package jmap

import "sync"

// A Capability broadcasts that the server supports underlying methods
type Capability interface {
	// The URI of the capability, eg "urn:ietf:params:jmap:core"
	URI() URI

	// Generates a pointer to a new Capability object
	New() Capability
}

var (
	capabilitiesMu sync.RWMutex
	capabilities   = make(map[URI]Capability)
)

// RegisterCapability registers a session/account capability decoder. It is
// safe for concurrent use with session decoding and other Register* calls.
func RegisterCapability(c Capability) {
	capabilitiesMu.Lock()
	capabilities[c.URI()] = c
	capabilitiesMu.Unlock()
}

func rangeCapabilities(fn func(URI, Capability)) {
	capabilitiesMu.RLock()
	defer capabilitiesMu.RUnlock()
	for key, cap := range capabilities {
		fn(key, cap)
	}
}
