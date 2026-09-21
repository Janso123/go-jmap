package quota

import "github.com/Janso123/go-jmap"

// Changes gets quota changes for the whole account.
type Changes struct {
	jmap.Changes[Quota]
}

// ChangesResponse is the result of Quota/changes.
// updatedProperties is Quota-specific (not on the kit ChangesResponse).
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
