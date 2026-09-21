package sharenotification

import "github.com/Janso123/go-jmap"

// Filter criteria for share notification queries.
// https://www.rfc-editor.org/rfc/rfc9670.html#section-3.4.1
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	After *jmap.UTCDate `json:"after,omitzero"`

	Before *jmap.UTCDate `json:"before,omitzero"`

	ObjectType string `json:"objectType,omitzero"`

	ObjectAccountID jmap.ID `json:"objectAccountId,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
