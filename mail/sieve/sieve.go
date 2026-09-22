package sieve

import "github.com/Janso123/go-jmap"

// URI is the JMAP Sieve capability (RFC 9661).
const URI jmap.URI = "urn:ietf:params:jmap:sieve"

const (
	// SieveScriptEvent is the SieveScript event type.
	SieveScriptEvent jmap.EventType = "SieveScript"
)

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterObject[SieveScript](
		jmap.MethodGet |
			jmap.MethodQuery |
			jmap.MethodSet,
	)
	// Manual method (not part of the standard kit).
	jmap.RegisterMethod("SieveScript/validate", newValidateResponse)
}

// Capability describes the JMAP Sieve capability.
//
// The same type is used for both the session capability object and the
// per-account capability object defined in RFC 9661 Section 1.2.1.
type Capability struct {
	Implementation      string                          `json:"implementation,omitzero"`
	MaxSizeScriptName   jmap.UnsignedInt                `json:"maxSizeScriptName"`
	MaxSizeScript       jmap.Optional[jmap.UnsignedInt] `json:"maxSizeScript,omitzero"`
	MaxNumberScripts    jmap.Optional[jmap.UnsignedInt] `json:"maxNumberScripts,omitzero"`
	MaxNumberRedirects  jmap.Optional[jmap.UnsignedInt] `json:"maxNumberRedirects,omitzero"`
	SieveExtensions     []string                        `json:"sieveExtensions,omitzero"`
	NotificationMethods jmap.Optional[[]string]         `json:"notificationMethods,omitzero"`
	ExternalLists       jmap.Optional[[]string]         `json:"externalLists,omitzero"`
}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// SieveScript is a JMAP SieveScript object.
type SieveScript struct {
	ID     jmap.ID               `json:"id,omitzero"`
	Name   jmap.Optional[string] `json:"name,omitzero"`
	BlobID jmap.ID               `json:"blobId,omitzero"`
	// IsActive is server-set: nil omits it on create, *false stays on the wire.
	IsActive *bool `json:"isActive,omitzero"`
}

func (SieveScript) JMAPType() string { return "SieveScript" }

func (SieveScript) JMAPCreatable() {}

func (SieveScript) Requires() []jmap.URI { return []jmap.URI{URI} }
