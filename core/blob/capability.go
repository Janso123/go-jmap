package blob

import "github.com/Janso123/go-jmap"

// URI is the JMAP Blob Management capability (RFC 9404).
const URI jmap.URI = "urn:ietf:params:jmap:blob"

func init() {
	jmap.RegisterCapability(&AccountCapability{})
}

// AccountCapability describes blob management limits advertised by an account.
// The same type is also used for the empty session capability object.
type AccountCapability struct {
	MaxSizeBlobSet            *uint64  `json:"maxSizeBlobSet,omitzero"`
	MaxDataSources            uint64   `json:"maxDataSources,omitzero"`
	SupportedTypeNames        []string `json:"supportedTypeNames,omitzero"`
	SupportedDigestAlgorithms []string `json:"supportedDigestAlgorithms,omitzero"`
}

func (c *AccountCapability) URI() jmap.URI { return URI }

func (c *AccountCapability) New() jmap.Capability { return &AccountCapability{} }
