package participantidentity

import "github.com/Janso123/go-jmap"

// Set creates, updates, and destroys participant identities.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3.3
type Set struct {
	jmap.Set[ParticipantIdentity]

	OnSuccessSetIsDefault jmap.Optional[jmap.ID] `json:"onSuccessSetIsDefault,omitzero"`
}

// SetResponse is the result of ParticipantIdentity/set.
type SetResponse = jmap.SetResponse[ParticipantIdentity]
