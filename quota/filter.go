package quota

import "github.com/Janso123/go-jmap"

type Filter interface {
	implementsFilter()
}

type FilterOperator struct {
	Operator jmap.Operator `json:"operator,omitempty"`

	Conditions []Filter `json:"conditions,omitempty"`
}

func (fo *FilterOperator) implementsFilter() {}

// Filter criteria for quota queries.
type FilterCondition struct {
	Name string `json:"name,omitempty"`

	Scope string `json:"scope,omitempty"`

	ResourceType string `json:"resourceType,omitempty"`

	Type string `json:"type,omitempty"`
}

func (fc *FilterCondition) implementsFilter() {}
