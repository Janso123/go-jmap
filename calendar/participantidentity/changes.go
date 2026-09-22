package participantidentity

import "github.com/Janso123/go-jmap"

// Changes gets participant identity changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3.2
type Changes struct {
	jmap.Changes[ParticipantIdentity]
}

// ChangesResponse is the result of ParticipantIdentity/changes.
type ChangesResponse = jmap.ChangesResponse
