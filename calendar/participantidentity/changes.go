package participantidentity

import "github.com/Janso123/go-jmap"

// Changes gets participant identity changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3.2
type Changes struct {
	jmap.Changes[ParticipantIdentity]
}

// ChangesResponse is the result of ParticipantIdentity/changes.
// updatedProperties is ParticipantIdentity-specific.
type ChangesResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	OldState string `json:"oldState,omitzero"`

	NewState string `json:"newState,omitzero"`

	HasMoreChanges bool `json:"hasMoreChanges,omitzero"`

	Created []jmap.ID `json:"created,omitzero"`

	Updated []jmap.ID `json:"updated,omitzero"`

	Destroyed []jmap.ID `json:"destroyed,omitzero"`

	UpdatedProperties []string `json:"updatedProperties,omitzero"`
}
