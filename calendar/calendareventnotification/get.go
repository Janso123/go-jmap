package calendareventnotification

import "github.com/Janso123/go-jmap"

// Get notification details.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.1
type Get struct {
	jmap.Get[CalendarEventNotification]
}

// GetResponse is the result of CalendarEventNotification/get.
type GetResponse = jmap.GetResponse[CalendarEventNotification]
