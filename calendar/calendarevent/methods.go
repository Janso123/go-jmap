package calendarevent

import (
	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
)

// Changes gets calendar event changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.8
type Changes struct {
	jmap.Changes[CalendarEvent]
}

// ChangesResponse is the result of CalendarEvent/changes.
// updatedProperties is CalendarEvent-specific.
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

// Set creates, updates, and destroys calendar events.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.9
type Set struct {
	jmap.Set[CalendarEvent]

	SendSchedulingMessages bool `json:"sendSchedulingMessages,omitzero"`
}

// SetResponse is the result of CalendarEvent/set.
type SetResponse = jmap.SetResponse[CalendarEvent]

// Copy events from one account to another.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.10
type Copy struct {
	jmap.Copy[CalendarEvent]
}

// CopyResponse is the result of CalendarEvent/copy.
type CopyResponse = jmap.CopyResponse[CalendarEvent]

// QueryChanges gets changes to a calendar event query since a given state.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.12
type QueryChanges struct {
	jmap.QueryChanges[CalendarEvent]

	Filter jmap.Filter `json:"filter,omitzero"`

	Sort []*jmap.Comparator `json:"sort,omitzero"`
}

// QueryChangesResponse is the result of CalendarEvent/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse

// Parse blobs as iCalendar files to get CalendarEvent objects.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.13
type Parse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	BlobIDs []jmap.ID `json:"blobIds,omitzero"`

	Properties []string `json:"properties,omitzero"`
}

func (m *Parse) Name() string { return "CalendarEvent/parse" }

func (m *Parse) Requires() []jmap.URI { return []jmap.URI{calendar.ParseURI} }

type ParseResponse struct {
	Account jmap.ID `json:"accountId,omitzero"`

	Parsed map[jmap.ID][]*CalendarEvent `json:"parsed,omitzero"`

	NotParsable []jmap.ID `json:"notParsable,omitzero"`

	NotFound []jmap.ID `json:"notFound,omitzero"`
}

func newParseResponse() jmap.MethodResponse { return &ParseResponse{} }
