package calendareventnotification

import "github.com/Janso123/go-jmap"

// Changes gets notification changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.2
type Changes struct {
	jmap.Changes[CalendarEventNotification]
}

// ChangesResponse is the result of CalendarEventNotification/changes.
// updatedProperties is CalendarEventNotification-specific.
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

// Set destroys event notifications.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.3
type Set struct {
	jmap.Set[CalendarEventNotification]
}

// SetResponse is the result of CalendarEventNotification/set.
type SetResponse = jmap.SetResponse[CalendarEventNotification]

// QueryChanges gets changes to a notification query since a given state.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.5
type QueryChanges struct {
	jmap.QueryChanges[CalendarEventNotification]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of CalendarEventNotification/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
