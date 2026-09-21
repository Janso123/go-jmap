package email

import "github.com/Janso123/go-jmap"

// Copy messages from one account to another
// https://www.rfc-editor.org/rfc/rfc8621.html#section-4.7
type Copy struct {
	jmap.Copy[Email]
}

// CopyResponse is the result of Email/copy.
type CopyResponse = jmap.CopyResponse[Email]
