package sieve

import "github.com/Janso123/go-jmap"

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// Filter criteria for Sieve script queries.
// https://www.rfc-editor.org/rfc/rfc9661.html#section-2.5
type FilterCondition struct {
	Name string `json:"name,omitempty"`

	IsActive *bool `json:"isActive,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
