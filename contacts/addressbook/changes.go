package addressbook

import "github.com/Janso123/go-jmap"

// Changes gets address book changes for the whole account.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-2.2
type Changes struct {
	jmap.Changes[AddressBook]
}

// ChangesResponse is the result of AddressBook/changes.
type ChangesResponse = jmap.ChangesResponse
