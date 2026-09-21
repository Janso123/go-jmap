package addressbook

import "github.com/Janso123/go-jmap"

// Get address book details.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2.1
type Get struct {
	jmap.Get[AddressBook]
}

// GetResponse is the result of AddressBook/get.
type GetResponse = jmap.GetResponse[AddressBook]
