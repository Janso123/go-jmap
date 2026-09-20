package participantidentity

import "git.sr.ht/~rockorager/go-jmap"

func init() {
	jmap.RegisterMethod("ParticipantIdentity/get", newGetResponse)
	jmap.RegisterMethod("ParticipantIdentity/changes", newChangesResponse)
	jmap.RegisterMethod("ParticipantIdentity/set", newSetResponse)
}

// ParticipantIdentity stores information about a URI that represents the user
// within an account in an event's participants.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3
type ParticipantIdentity struct {
	ID jmap.ID `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	CalendarAddress string `json:"calendarAddress,omitempty"`

	IsDefault bool `json:"isDefault,omitempty"`
}
