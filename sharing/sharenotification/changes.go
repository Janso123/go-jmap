package sharenotification

import "github.com/Janso123/go-jmap"

// Changes gets share notification changes for the whole account.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.2
type Changes struct {
	jmap.Changes[ShareNotification]
}

// ChangesResponse is the result of ShareNotification/changes.
// updatedProperties is ShareNotification-specific (not on the kit ChangesResponse).
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
