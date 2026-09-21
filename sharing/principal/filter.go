package principal

import "github.com/Janso123/go-jmap"

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// Filter criteria for principal queries.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-2.3
type FilterCondition struct {
	AccountIDs []jmap.ID `json:"accountIds,omitempty"`

	Email string `json:"email,omitempty"`

	Name string `json:"name,omitempty"`

	Text string `json:"text,omitempty"`

	Type PrincipalType `json:"type,omitempty"`

	CalendarAddress string `json:"calendarAddress,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
