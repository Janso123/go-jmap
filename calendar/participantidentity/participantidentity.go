package participantidentity

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

func init() {
	jmap.RegisterObject[ParticipantIdentity](
		jmap.MethodGet |
			jmap.MethodChanges |
			jmap.MethodSet,
	)
	// ParticipantIdentity/changes includes updatedProperties; override kit factory.
	jmap.RegisterMethod("ParticipantIdentity/changes", func() jmap.MethodResponse { return &ChangesResponse{} })
}

// ParticipantIdentity stores information about a URI that represents the user
// within an account in an event's participants.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3
type ParticipantIdentity struct {
	ID jmap.ID `json:"id,omitzero"`

	Name string `json:"name,omitzero"`

	CalendarAddress string `json:"calendarAddress,omitzero"`

	IsDefault bool `json:"isDefault,omitzero"`
}

func (ParticipantIdentity) JMAPType() string { return "ParticipantIdentity" }

func (ParticipantIdentity) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }
