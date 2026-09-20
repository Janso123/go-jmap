package calendarevent

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/calendar"
)

// Changes gets calendar event changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.8
type Changes struct {
	Account jmap.ID `json:"accountId,omitempty"`

	SinceState string `json:"sinceState,omitempty"`

	MaxChanges uint64 `json:"maxChanges,omitempty"`
}

func (m *Changes) Name() string { return "CalendarEvent/changes" }

func (m *Changes) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type ChangesResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	NewState string `json:"newState,omitempty"`

	HasMoreChanges bool `json:"hasMoreChanges,omitempty"`

	Created []jmap.ID `json:"created,omitempty"`

	Updated []jmap.ID `json:"updated,omitempty"`

	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	UpdatedProperties []string `json:"updatedProperties,omitempty"`
}

func newChangesResponse() jmap.MethodResponse { return &ChangesResponse{} }

// Set creates, updates, and destroys calendar events.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.9
type Set struct {
	Account jmap.ID `json:"accountId,omitempty"`

	IfInState string `json:"ifInState,omitempty"`

	Create map[jmap.ID]*CalendarEvent `json:"create,omitempty"`

	Update map[jmap.ID]jmap.Patch `json:"update,omitempty"`

	Destroy []jmap.ID `json:"destroy,omitempty"`

	SendSchedulingMessages bool `json:"sendSchedulingMessages,omitempty"`
}

func (m *Set) Name() string { return "CalendarEvent/set" }

func (m *Set) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type SetResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	NewState string `json:"newState,omitempty"`

	Created map[jmap.ID]*CalendarEvent `json:"created,omitempty"`

	Updated map[jmap.ID]*CalendarEvent `json:"updated,omitempty"`

	Destroyed []jmap.ID `json:"destroyed,omitempty"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`

	NotUpdated map[jmap.ID]*jmap.SetError `json:"notUpdated,omitempty"`

	NotDestroyed map[jmap.ID]*jmap.SetError `json:"notDestroyed,omitempty"`
}

func newSetResponse() jmap.MethodResponse { return &SetResponse{} }

// Copy events from one account to another.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.10
type Copy struct {
	FromAccount jmap.ID `json:"fromAccountId,omitempty"`

	IfFromInState string `json:"ifFromInState,omitempty"`

	Account jmap.ID `json:"accountId,omitempty"`

	IfInState string `json:"ifInState,omitempty"`

	Create map[jmap.ID]*CalendarEvent `json:"create,omitempty"`

	OnSuccessDestroyOriginal bool `json:"onSuccessDestroyOriginal,omitempty"`

	DestroyFromIfInState string `json:"destroyFromIfInState,omitempty"`
}

func (m *Copy) Name() string { return "CalendarEvent/copy" }

func (m *Copy) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type CopyResponse struct {
	FromAccount jmap.ID `json:"fromAccountId,omitempty"`

	Account jmap.ID `json:"accountId,omitempty"`

	OldState string `json:"oldState,omitempty"`

	NewState string `json:"newState,omitempty"`

	Created map[jmap.ID]*CalendarEvent `json:"created,omitempty"`

	NotCreated map[jmap.ID]*jmap.SetError `json:"notCreated,omitempty"`
}

func newCopyResponse() jmap.MethodResponse { return &CopyResponse{} }

// QueryChanges gets changes to a calendar event query since a given state.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.12
type QueryChanges struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Filter Filter `json:"filter,omitempty"`

	Sort []*SortComparator `json:"sort,omitempty"`

	SinceQueryState string `json:"sinceQueryState,omitempty"`

	MaxChanges uint64 `json:"maxChanges,omitempty"`

	UpToID jmap.ID `json:"upToId,omitempty"`

	CalculateTotal bool `json:"calculateTotal,omitempty"`
}

func (m *QueryChanges) Name() string { return "CalendarEvent/queryChanges" }

func (m *QueryChanges) Requires() []jmap.URI { return []jmap.URI{calendar.URI} }

type QueryChangesResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	OldQueryState string `json:"oldQueryState,omitempty"`

	NewQueryState string `json:"newQueryState,omitempty"`

	Removed []jmap.ID `json:"removed,omitempty"`

	Added []jmap.AddedItem `json:"added,omitempty"`
}

func newQueryChangesResponse() jmap.MethodResponse { return &QueryChangesResponse{} }

// Parse blobs as iCalendar files to get CalendarEvent objects.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-5.13
type Parse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	BlobIDs []jmap.ID `json:"blobIds,omitempty"`

	Properties []string `json:"properties,omitempty"`
}

func (m *Parse) Name() string { return "CalendarEvent/parse" }

func (m *Parse) Requires() []jmap.URI { return []jmap.URI{calendar.ParseURI} }

type ParseResponse struct {
	Account jmap.ID `json:"accountId,omitempty"`

	Parsed map[jmap.ID][]*CalendarEvent `json:"parsed,omitempty"`

	NotParsable []jmap.ID `json:"notParsable,omitempty"`

	NotFound []jmap.ID `json:"notFound,omitempty"`
}

func newParseResponse() jmap.MethodResponse { return &ParseResponse{} }
