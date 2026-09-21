package contacts

import "github.com/Janso123/go-jmap"

// URI is the JMAP Contacts capability (RFC 9610).
const URI jmap.URI = "urn:ietf:params:jmap:contacts"

const (
	// AddressBookEvent is the AddressBook event type.
	AddressBookEvent jmap.EventType = "AddressBook"

	// ContactCardEvent is the ContactCard event type.
	ContactCardEvent jmap.EventType = "ContactCard"
)

func init() {
	jmap.RegisterCapability(&Capability{})
}

// Capability describes the JMAP Contacts capability.
// The same type is also used for the empty session capability object.
type Capability struct {
	MaxAddressBooksPerCard *uint64 `json:"maxAddressBooksPerCard,omitzero"`
	MayCreateAddressBook   *bool   `json:"mayCreateAddressBook,omitzero"`
}

func (c *Capability) URI() jmap.URI { return URI }

func (c *Capability) New() jmap.Capability { return &Capability{} }
