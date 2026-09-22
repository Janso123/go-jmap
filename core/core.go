package core

import "github.com/Janso123/go-jmap"

func init() {
	jmap.RegisterCapability(&Core{})
	jmap.RegisterMethod("Core/echo", newEcho)
}

type Core struct {
	// The maximum file size, in bytes
	MaxSizeUpload       jmap.UnsignedInt `json:"maxSizeUpload"`
	MaxConcurrentUpload jmap.UnsignedInt `json:"maxConcurrentUpload"`

	// The maximum size, in bytes, that the server will accept for a request
	MaxSizeRequest        jmap.UnsignedInt `json:"maxSizeRequest"`
	MaxConcurrentRequests jmap.UnsignedInt `json:"maxConcurrentRequests"`
	MaxCallsInRequest     jmap.UnsignedInt `json:"maxCallsInRequest"`

	// The maximum number of objects that the client may request in a single
	// /get type method call.
	MaxObjectsInGet     jmap.UnsignedInt     `json:"maxObjectsInGet"`
	MaxObjectsInSet     jmap.UnsignedInt     `json:"maxObjectsInSet"`
	CollationAlgorithms []jmap.CollationAlgo `json:"collationAlgorithms"`
}

func (c *Core) URI() jmap.URI { return jmap.CoreURI }

func (c *Core) New() jmap.Capability { return &Core{} }
