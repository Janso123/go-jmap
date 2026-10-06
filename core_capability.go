package jmap

func init() {
	RegisterCapability(&Core{})
}

// Core is the urn:ietf:params:jmap:core capability (RFC 8620 §2). It lives
// here rather than in package core so every decoded Session carries it.
type Core struct {
	// The maximum file size, in bytes
	MaxSizeUpload       UnsignedInt `json:"maxSizeUpload"`
	MaxConcurrentUpload UnsignedInt `json:"maxConcurrentUpload"`

	// The maximum size, in bytes, that the server will accept for a request
	MaxSizeRequest        UnsignedInt `json:"maxSizeRequest"`
	MaxConcurrentRequests UnsignedInt `json:"maxConcurrentRequests"`
	MaxCallsInRequest     UnsignedInt `json:"maxCallsInRequest"`

	// The maximum number of objects that the client may request in a single
	// /get type method call.
	MaxObjectsInGet     UnsignedInt     `json:"maxObjectsInGet"`
	MaxObjectsInSet     UnsignedInt     `json:"maxObjectsInSet"`
	CollationAlgorithms []CollationAlgo `json:"collationAlgorithms"`
}

func (c *Core) URI() URI { return CoreURI }

func (c *Core) New() Capability { return &Core{} }
