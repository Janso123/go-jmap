package sharenotification

import (
	"time"

	"git.sr.ht/~rockorager/go-jmap"
)

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// Filter criteria for share notification queries.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.4.1
type FilterCondition struct {
	After *time.Time `json:"after,omitempty"`

	Before *time.Time `json:"before,omitempty"`

	ObjectType string `json:"objectType,omitempty"`

	ObjectAccountID jmap.ID `json:"objectAccountId,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
