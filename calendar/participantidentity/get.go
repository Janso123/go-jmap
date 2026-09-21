package participantidentity

import "github.com/Janso123/go-jmap"

// Get participant identity details.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3.1
type Get struct {
	jmap.Get[ParticipantIdentity]
}

// GetResponse is the result of ParticipantIdentity/get.
type GetResponse = jmap.GetResponse[ParticipantIdentity]
