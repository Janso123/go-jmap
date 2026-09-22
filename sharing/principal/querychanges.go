package principal

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/calendar"
	"github.com/Janso123/go-jmap/sharing"
)

// Get changes on a principal query.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.4
type QueryChanges struct {
	jmap.QueryChanges[Principal]

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

func (q *QueryChanges) Requires() []jmap.URI {
	uris := []jmap.URI{sharing.URI}
	if filterHasCalendarAddress(q.Filter) {
		uris = append(uris, calendar.AvailabilityURI)
	}
	return uris
}

// QueryChangesResponse is the result of Principal/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
