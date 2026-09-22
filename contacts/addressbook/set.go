package addressbook

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys address books.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2.3
type Set struct {
	jmap.Set[AddressBook]

	OnDestroyRemoveContents bool `json:"onDestroyRemoveContents,omitzero"`

	OnSuccessSetIsDefault jmap.Optional[jmap.ID] `json:"onSuccessSetIsDefault,omitzero"`
}

// SetResponse is the result of AddressBook/set.
type SetResponse = jmap.SetResponse[AddressBook]
