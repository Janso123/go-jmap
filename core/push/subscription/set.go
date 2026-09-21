package subscription

import "github.com/Janso123/go-jmap"

// Modify push subscription details
// https://www.rfc-editor.org/rfc/rfc8620.html#section-7.2.2
type Set struct {
	jmap.Set[PushSubscription]
}

// SetResponse is the result of PushSubscription/set.
type SetResponse = jmap.SetResponse[PushSubscription]
