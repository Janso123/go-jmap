package sharing

import "github.com/Janso123/go-jmap"

// URI is the JMAP principals session capability (RFC 9670).
const URI jmap.URI = "urn:ietf:params:jmap:principals"

// OwnerURI is the JMAP principals owner account capability (RFC 9670).
const OwnerURI jmap.URI = "urn:ietf:params:jmap:principals:owner"

func init() {
	jmap.RegisterCapability(&Capability{})
	jmap.RegisterCapability(&OwnerCapability{})
}

// Capability describes the principals session capability.
type Capability struct {
	CurrentUserPrincipalID *jmap.ID `json:"currentUserPrincipalId,omitempty"`
}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }

// OwnerCapability describes the principals owner account capability.
type OwnerCapability struct {
	AccountIDForPrincipal jmap.ID `json:"accountIdForPrincipal,omitempty"`
	PrincipalID           jmap.ID `json:"principalId,omitempty"`
}

func (c *OwnerCapability) URI() jmap.URI { return OwnerURI }

func (c *OwnerCapability) New() jmap.Capability { return &OwnerCapability{} }
