package contactcard

import "github.com/Janso123/go-jmap"

// Changes gets contact card changes for the whole account.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.2
type Changes struct {
	jmap.Changes[ContactCard]
}

// ChangesResponse is the result of ContactCard/changes.
type ChangesResponse = jmap.ChangesResponse
