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
	// AddressBook/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("AddressBook/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

// AddressBook is a named collection of ContactCards.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2
type AddressBook struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	Description string `json:"description,omitzero"`

	SortOrder uint64 `json:"sortOrder,omitzero"`

	IsDefault bool `json:"isDefault,omitzero"`

	IsSubscribed *bool `json:"isSubscribed,omitzero"`

	ShareWith map[jmap.ID]*Rights `json:"shareWith,omitzero"`

	MyRights *Rights `json:"myRights,omitzero"`
}

func (AddressBook) JMAPType() string { return "AddressBook" }

func (AddressBook) Requires() []jmap.URI { return []jmap.URI{contacts.URI} }

// Rights is the set of permissions the user has on an AddressBook.
type Rights struct {
	MayRead bool `json:"mayRead,omitzero"`

	MayWrite bool `json:"mayWrite,omitzero"`

	MayShare bool `json:"mayShare,omitzero"`

	MayDelete bool `json:"mayDelete,omitzero"`
}
