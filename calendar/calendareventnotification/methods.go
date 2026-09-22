package calendareventnotification

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Changes gets notification changes for the whole account.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.2
type Changes struct {
	jmap.Changes[CalendarEventNotification]
}

// ChangesResponse is the result of CalendarEventNotification/changes.
type ChangesResponse = jmap.ChangesResponse

// Set destroys event notifications.
// https://datatracker.ietf.org/doc/html/draft-ietf-jmap-calendars-29#section-7.3
// Create/update are not permitted.
type Set struct {
	Account          jmap.ID               `json:"accountId,omitzero"`
	IfInState        string                `json:"ifInState,omitzero"`
	Destroy          []jmap.ID             `json:"destroy,omitzero"`
	ReferenceDestroy *jmap.ResultReference `json:"#destroy,omitzero"`
}

func (m *Set) Name() string { return "CalendarEventNotification/set" }

func (m *Set) Requires() []jmap.URI {
	var zero CalendarEventNotification
	return zero.Requires()
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

type queryChangesAlias QueryChanges

func (q *QueryChanges) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s struct {
		queryChangesAlias
		Filter jsontext.Value `json:"filter"`
	}
	if err := jsonv2.UnmarshalDecode(dec, &s); err != nil {
		return err
	}
	f, err := jmap.UnmarshalFilter[FilterCondition](s.Filter)
	if err != nil {
		return err
	}
	*q = QueryChanges(s.queryChangesAlias)
	q.Filter = f
	return nil
}

// QueryChangesResponse is the result of CalendarEventNotification/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
