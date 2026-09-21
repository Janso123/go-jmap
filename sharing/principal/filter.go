package principal

import "github.com/Janso123/go-jmap"

// Filter criteria for principal queries.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	AccountIDs []jmap.ID `json:"accountIds,omitzero"`

	Email string `json:"email,omitzero"`

	Name string `json:"name,omitzero"`

	Text string `json:"text,omitzero"`

	Type PrincipalType `json:"type,omitzero"`

	TimeZone string `json:"timeZone,omitzero"`

	CalendarAddress string `json:"calendarAddress,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
