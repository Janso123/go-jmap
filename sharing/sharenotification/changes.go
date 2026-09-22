package sharenotification

import "github.com/Janso123/go-jmap"

// Changes gets share notification changes for the whole account.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.2
type Changes struct {
	jmap.Changes[ShareNotification]
}

// ChangesResponse is the result of ShareNotification/changes.
// RFC 9670 §3.2 is a standard /changes method.
type ChangesResponse = jmap.ChangesResponse
