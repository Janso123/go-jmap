package email

import "github.com/Janso123/go-jmap"

// Get changes to emails on the whole account since a given state
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.3
type Changes struct {
	jmap.Changes[Email]
}

// ChangesResponse is the result of Email/changes.
type ChangesResponse = jmap.ChangesResponse
