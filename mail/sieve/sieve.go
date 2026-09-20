package sieve

import "git.sr.ht/~rockorager/go-jmap"

// URI is the JMAP Sieve capability (RFC 9661).
const URI jmap.URI = "urn:ietf:params:jmap:sieve"

const (
	// SieveScriptEvent is the SieveScript event type.
	SieveScriptEvent jmap.EventType = "SieveScript"
)

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterMethod("SieveScript/get", newGetResponse)
	jmap.RegisterMethod("SieveScript/query", newQueryResponse)
	jmap.RegisterMethod("SieveScript/set", newSetResponse)
	jmap.RegisterMethod("SieveScript/validate", newValidateResponse)
}

// Capability describes the JMAP Sieve capability.
//
// The same type is used for both the session capability object and the
// per-account capability object defined in RFC 9661 Section 1.2.1.
type Capability struct {
	Implementation      string   `json:"implementation,omitempty"`
	MaxSizeScriptName   uint64   `json:"maxSizeScriptName,omitempty"`
	MaxSizeScript       *uint64  `json:"maxSizeScript,omitempty"`
	MaxNumberScripts    *uint64  `json:"maxNumberScripts,omitempty"`
	MaxNumberRedirects  *uint64  `json:"maxNumberRedirects,omitempty"`
	SieveExtensions     []string `json:"sieveExtensions,omitempty"`
	NotificationMethods []string `json:"notificationMethods,omitempty"`
	ExternalLists       []string `json:"externalLists,omitempty"`
}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// SieveScript is a JMAP SieveScript object.
type SieveScript struct {
	ID       jmap.ID `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	BlobID   jmap.ID `json:"blobId,omitempty"`
	IsActive bool    `json:"isActive,omitempty"`
}
