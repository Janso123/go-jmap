package principal

import "github.com/Janso123/go-jmap"

// Get principal details.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.1
type Get struct {
	jmap.Get[Principal]
}

// GetResponse is the result of Principal/get.
type GetResponse = jmap.GetResponse[Principal]
