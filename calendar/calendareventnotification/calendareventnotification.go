package calendareventnotification

import (
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/calendar/jscalendar"
)

func init() {
	jmap.RegisterObject[CalendarEventNotification](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodQuery |
			jmap.MethodQueryChanges |
			jmap.MethodSet,
	)
	// CalendarEventNotification/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("CalendarEventNotification/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
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
	Name string `json:"name,omitzero"`

	Email *string `json:"email,omitzero"`

	PrincipalID *jmap.ID `json:"principalId,omitzero"`

	CalendarAddress *string `json:"calendarAddress,omitzero"`

	Comment *string `json:"comment,omitzero"`
}

// CalendarEventNotification tracks server-created event change notifications.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7
type CalendarEventNotification struct {
	ID jmap.ID `json:"id,omitzero"`

	Created *time.Time `json:"created,omitzero"`

	ChangedBy *Person `json:"changedBy,omitzero"`

	Type Type `json:"type,omitzero"`

	CalendarEventID jmap.ID `json:"calendarEventId,omitzero"`

	IsDraft *bool `json:"isDraft,omitzero"`

	Event *jscalendar.Event `json:"event,omitzero"`

	EventPatch jmap.Patch `json:"eventPatch,omitzero"`
}

func (CalendarEventNotification) JMAPType() string { return "CalendarEventNotification" }

func (CalendarEventNotification) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }
