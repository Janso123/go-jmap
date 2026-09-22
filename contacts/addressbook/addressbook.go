package addressbook

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/contacts"
)

func init() {
	jmap.RegisterObject[AddressBook](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodSet,
	)
}

// AddressBook is a named collection of ContactCards.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2
type AddressBook struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	Description jmap.Optional[string] `json:"description,omitzero"`

	SortOrder uint64 `json:"sortOrder,omitzero"`

	IsDefault *bool `json:"isDefault,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`

	ShareWith jmap.Optional[map[jmap.ID]*Rights] `json:"shareWith,omitzero"`

	MyRights *Rights `json:"myRights,omitzero"`
}

func (AddressBook) JMAPType() string { return "AddressBook" }

func (AddressBook) JMAPCreatable() {}

func (AddressBook) Requires() []jmap.URI { return []jmap.URI{contacts.URI} }

// Rights is the set of permissions the user has on an AddressBook.
// Bools are bare so false stays on the wire.
type Rights struct {
	MayRead bool `json:"mayRead"`

	MayWrite bool `json:"mayWrite"`

	MayShare bool `json:"mayShare"`

	MayDelete bool `json:"mayDelete"`
}
