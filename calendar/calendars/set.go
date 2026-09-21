package calendars

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys calendars.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-4.3
type Set struct {
	jmap.Set[Calendar]

	OnDestroyRemoveEvents bool `json:"onDestroyRemoveEvents,omitzero"`

	OnSuccessSetIsDefault jmap.ID `json:"onSuccessSetIsDefault,omitzero"`
}

// SetResponse is the result of Calendar/set.
type SetResponse = jmap.SetResponse[Calendar]
