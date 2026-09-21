package calendars

import "github.com/Janso123/go-jmap"

// Get calendar details.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-4.1
type Get struct {
	jmap.Get[Calendar]
}

// GetResponse is the result of Calendar/get.
type GetResponse = jmap.GetResponse[Calendar]
