package sieve

import "github.com/Janso123/go-jmap"

// Filter criteria for Sieve script queries.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.5
type FilterCondition struct {
	jmap.FilterBase `json:"-"`

	Name string `json:"name,omitzero"`

	IsActive *bool `json:"isActive,omitzero"`
}

func And(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.And(conds...) }
func Or(conds ...jmap.Filter) *jmap.FilterOperator  { return jmap.Or(conds...) }
func Not(conds ...jmap.Filter) *jmap.FilterOperator { return jmap.Not(conds...) }
