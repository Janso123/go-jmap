package emailsubmission

import "github.com/Janso123/go-jmap"

// Get email submission changes for the whole account
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.2
type Changes struct {
	jmap.Changes[EmailSubmission]
}

// ChangesResponse is the result of EmailSubmission/changes.
type ChangesResponse = jmap.ChangesResponse
