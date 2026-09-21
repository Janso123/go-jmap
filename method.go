package jmap

import "sync"

// A JMAP method. The method object will be marshaled as the arguments to an
// invocation.
type Method interface {
	// The name of the method, ie "Core/echo"
	Name() string

	// The JMAP capabilities required for the method, ie "urn:ietf:params:jmap:core"
	Requires() []URI
}

// A response to a method call
type MethodResponse any

// A Factory function which produces a new MethodResponse object
type MethodResponseFactory func() MethodResponse

var (
	methodsMu sync.RWMutex
	methods   = map[string]MethodResponseFactory{}
)

// RegisterMethod registers a method response factory. It is safe for concurrent
// use with response decoding and other Register* calls (RFC 8620 §3.10).
// Prefer RegisterObject for standard Object methods (get/changes/query/queryChanges/set/copy).
func RegisterMethod(name string, factory MethodResponseFactory) {
	methodsMu.Lock()
	methods[name] = factory
	methodsMu.Unlock()
}

func lookupMethod(name string) (MethodResponseFactory, bool) {
	methodsMu.RLock()
	fn, ok := methods[name]
	methodsMu.RUnlock()
	return fn, ok
}
