package contactcard

import "github.com/Janso123/go-jmap"

// Copy copies contact cards from one account to another.
// https://www.rfc-editor.org/rfc/rfc9610.html#section-3.6
type Copy struct {
	jmap.Copy[ContactCard]
}

// CopyResponse is the result of ContactCard/copy.
type CopyResponse = jmap.CopyResponse[ContactCard]
