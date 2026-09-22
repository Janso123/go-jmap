package sharenotification

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
)

// Get changes on a share notification query.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.5
type QueryChanges struct {
	jmap.QueryChanges[ShareNotification]

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

// QueryChangesResponse is the result of ShareNotification/queryChanges.
type QueryChangesResponse = jmap.QueryChangesResponse
