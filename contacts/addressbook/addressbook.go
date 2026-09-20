package addressbook

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("AddressBook/get", newGetResponse)
	jmap.RegisterMethod("AddressBook/changes", newChangesResponse)
	jmap.RegisterMethod("AddressBook/set", newSetResponse)
}

// AddressBook is a named collection of ContactCards.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2
type AddressBook struct {
	ID jmap.ID `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	SortOrder uint64 `json:"sortOrder,omitempty"`

	IsDefault bool `json:"isDefault,omitempty"`

	IsSubscribed bool `json:"isSubscribed,omitempty"`

	ShareWith map[jmap.ID]*Rights `json:"shareWith,omitempty"`

	MyRights *Rights `json:"myRights,omitempty"`
}

// Rights is the set of permissions the user has on an AddressBook.
type Rights struct {
	MayRead bool `json:"mayRead,omitempty"`

	MayWrite bool `json:"mayWrite,omitempty"`

	MayShare bool `json:"mayShare,omitempty"`

	MayDelete bool `json:"mayDelete,omitempty"`
}
