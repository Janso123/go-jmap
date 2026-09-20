package calendareventnotification

import (
	"time"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar/jscalendar"
)

func init() {
	jmap.RegisterMethod("CalendarEventNotification/get", newGetResponse)
	jmap.RegisterMethod("CalendarEventNotification/changes", newChangesResponse)
	jmap.RegisterMethod("CalendarEventNotification/set", newSetResponse)
	jmap.RegisterMethod("CalendarEventNotification/query", newQueryResponse)
	jmap.RegisterMethod("CalendarEventNotification/queryChanges", newQueryChangesResponse)
}

const (
	// CalendarEventNotificationEvent is the event source type for notifications.
	CalendarEventNotificationEvent jmap.EventType = "CalendarEventNotification"
)

// Type identifies the kind of event change represented by a notification.
type Type string

const (
	TypeCreated   Type = "created"
	TypeUpdated   Type = "updated"
	TypeDestroyed Type = "destroyed"
)

// Person describes who caused the event change.
type Person struct {
	Name string `json:"name,omitempty"`

	Email *string `json:"email,omitempty"`

	PrincipalID *jmap.ID `json:"principalId,omitempty"`

	CalendarAddress *string `json:"calendarAddress,omitempty"`

	Comment *string `json:"comment,omitempty"`
}

// CalendarEventNotification tracks server-created event change notifications.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7
type CalendarEventNotification struct {
	ID jmap.ID `json:"id,omitempty"`

	Created *time.Time `json:"created,omitempty"`

	ChangedBy *Person `json:"changedBy,omitempty"`

	Type Type `json:"type,omitempty"`

	CalendarEventID jmap.ID `json:"calendarEventId,omitempty"`

	IsDraft *bool `json:"isDraft,omitempty"`

	Event *jscalendar.Event `json:"event,omitempty"`

	EventPatch jmap.Patch `json:"eventPatch,omitempty"`
}
