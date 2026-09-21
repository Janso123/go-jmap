package sharenotification

import "github.com/Janso123/go-jmap"

// Get share notification details.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.1
type Get struct {
	jmap.Get[ShareNotification]
}

// GetResponse is the result of ShareNotification/get.
type GetResponse = jmap.GetResponse[ShareNotification]
