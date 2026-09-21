package emailsubmission

import "github.com/Janso123/go-jmap"

// Get email submission details
// https://www.rfc-editor.org/rfc/rfc8621.html#section-7.1
type Get struct {
	jmap.Get[EmailSubmission]
}

// GetResponse is the result of EmailSubmission/get.
type GetResponse = jmap.GetResponse[EmailSubmission]
