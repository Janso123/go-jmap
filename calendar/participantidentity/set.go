package participantidentity

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
)

// Set creates, updates, and destroys participant identities.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-3.3
type Set struct {
	Account jmap.ID `json:"accountId,omitempty"`

	IfInState string `json:"ifInState,omitempty"`

	Create map[jmap.ID]*ParticipantIdentity `json:"create,omitempty"`

	Update map[jmap.ID]jmap.Patch `json:"update,omitempty"`

	Destroy []jmap.ID `json:"destroy,omitempty"`

	OnSuccessSetIsDefault jmap.ID `json:"onSuccessSetIsDefault,omitempty"`
}

func (m *Set) Name() string { return "ParticipantIdentity/set" }

func (m *Set) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type SetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	NewState string `json:"newState,omitempty"`

	Created map[jmap.ID]*ParticipantIdentity `json:"created,omitempty"`

	Updated map[jmap.ID]*ParticipantIdentity `json:"updated,omitempty"`

	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`

	NotUpdated map[jmap.ID]*jmap.SetError `json:"notUpdated,omitempty"`

	NotDestroyed map[jmap.ID]*jmap.SetError `json:"notDestroyed,omitempty"`
}

func newSetResponse() jmap.MethodResponse { return &SetResponse{} }
