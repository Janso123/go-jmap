package sharing

import "github.com/Janso123/go-jmap"

// URI is the JMAP principals session capability (RFC 9670).
const URI jmap.URI = "urn:ietf:params:jmap:principals"

// OwnerURI is the JMAP principals owner account capability (RFC 9670).
const OwnerURI jmap.URI = "urn:ietf:params:jmap:principals:owner"

func init() {
	jmap.RegisterCapability(&AccountCapability{})
	jmap.RegisterCapability(&OwnerCapability{})
}

// AccountCapability is the urn:ietf:params:jmap:principals capability object.
// currentUserPrincipalId is Id|null on an account's accountCapabilities
// object (RFC 9670 §1.5.1). The session object for this URI is empty;
// the registry has one type per URI, so this type decodes both.
type AccountCapability struct {
	CurrentUserPrincipalID jmap.Optional[jmap.ID] `json:"currentUserPrincipalId,omitzero"`
}

func (c *AccountCapability) URI() jmap.URI { return URI }

func (c *AccountCapability) New() jmap.Capability { return &AccountCapability{} }

// OwnerCapability describes the principals owner account capability.
type OwnerCapability struct {
	AccountIDForPrincipal jmap.ID `json:"accountIdForPrincipal,omitzero"`
	PrincipalID           jmap.ID `json:"principalId,omitzero"`
}

func (c *OwnerCapability) URI() jmap.URI { return OwnerURI }

func (c *OwnerCapability) New() jmap.Capability { return &OwnerCapability{} }
