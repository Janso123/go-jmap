package sharenotification

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys share notifications.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.3
type Set struct {
	jmap.Set[ShareNotification]
}

// SetResponse is the result of ShareNotification/set.
type SetResponse = jmap.SetResponse[ShareNotification]
