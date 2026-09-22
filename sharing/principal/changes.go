package principal

import "github.com/Janso123/go-jmap"

// Changes gets principal changes for the whole account.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.2
type Changes struct {
	jmap.Changes[Principal]
}

// ChangesResponse is the result of Principal/changes.
// RFC 9670 §2.2 is a standard /changes method.
type ChangesResponse = jmap.ChangesResponse
