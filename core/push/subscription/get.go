package subscription

import "github.com/Janso123/go-jmap"

// Get push subscription details
// https://www.rfc-editor.org/rfc/rfc8620.html#section-7.2.1
type Get struct {
	jmap.Get[PushSubscription]
}

// GetResponse is the result of PushSubscription/get.
type GetResponse = jmap.GetResponse[PushSubscription]
